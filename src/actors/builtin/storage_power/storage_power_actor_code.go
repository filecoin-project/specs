package storage_power

import (
	"math"

	actors "github.com/filecoin-project/specs/actors"
	serde "github.com/filecoin-project/specs/actors/serde"
	autil "github.com/filecoin-project/specs/actors/util"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	indices "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
)

type ConsensusFaultType int

const (
	UncommittedPowerFault ConsensusFaultType = 0
	DoubleForkMiningFault ConsensusFaultType = 1
	ParentGrindingFault   ConsensusFaultType = 2
	TimeOffsetMiningFault ConsensusFaultType = 3
)

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) AddBalance(rt Runtime, minerAddr addr.Address) {
	RT_MinerEntry_ValidateCaller_DetermineFundsLocation(rt, minerAddr, vmr.MinerEntrySpec_MinerOnly)

	msgValue := rt.ValueReceived()

	h, st := a.State(rt)
	newTable, ok := autil.BalanceTable_WithAdd(st.EscrowTable(), minerAddr, msgValue)
	if !ok {
		rt.AbortStateMsg("Escrow operation failed")
	}
	st.Impl().EscrowTable_ = newTable
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, minerAddr addr.Address, amountRequested actors.TokenAmount) {
	if amountRequested < 0 {
		rt.AbortArgMsg("Amount to withdraw must be nonnegative")
	}

	recipientAddr := RT_MinerEntry_ValidateCaller_DetermineFundsLocation(rt, minerAddr, vmr.MinerEntrySpec_MinerOnly)

	minBalanceMaintainRequired := a._rtGetPledgeCollateralReqForMinerOrAbort(rt, minerAddr)

	h, st := a.State(rt)
	newTable, amountExtracted, ok := autil.BalanceTable_WithExtractPartial(
		st.EscrowTable(), minerAddr, amountRequested, minBalanceMaintainRequired)
	if !ok {
		rt.AbortStateMsg("Escrow operation failed")
	}
	st.Impl().EscrowTable_ = newTable
	UpdateRelease(rt, h, st)

	rt.SendFunds(recipientAddr, amountExtracted)
}

func (a *StoragePowerActorCode_I) CreateMiner(rt Runtime, workerAddr addr.Address, sectorSize sector.SectorSize, peerId libp2p.PeerID) addr.Address {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	ownerAddr := rt.ImmediateCaller()

	newMinerAddr := addr.Deserialize_Address_Compact_Assert(
		rt.Send(
			addr.InitActorAddr,
			ai.Method_InitActor_Exec,
			serde.MustSerializeParams(
				actor.CodeID_Make_Builtin(actor.BuiltinActorID_StorageMiner),
				ownerAddr,
				workerAddr,
				sectorSize,
				peerId,
			),
			actors.TokenAmount(0),
		).ReturnValue(),
	)

	h, st := a.State(rt)
	newTable, ok := autil.BalanceTable_WithNewAddressEntry(st.EscrowTable(), newMinerAddr, rt.ValueReceived())
	Assert(ok)
	st.Impl().EscrowTable_ = newTable
	st.PowerTable()[newMinerAddr] = block.StoragePower(0)
	st.ClaimedPower()[newMinerAddr] = block.StoragePower(0)
	st.NominalPower()[newMinerAddr] = block.StoragePower(0)
	UpdateRelease(rt, h, st)

	return newMinerAddr
}

func (a *StoragePowerActorCode_I) DeleteMiner(rt Runtime, minerAddr addr.Address) {
	h, st := a.State(rt)

	minerPledgeBalance, ok := autil.BalanceTable_GetEntry(st.EscrowTable(), minerAddr)
	if !ok {
		rt.AbortArgMsg("Miner address not found")
	}

	if minerPledgeBalance > actors.TokenAmount(0) {
		rt.AbortStateMsg("Deletion requested for miner with pledge balance still remaining")
	}

	minerPower, ok := st.PowerTable()[minerAddr]
	Assert(ok)
	if minerPower > 0 {
		rt.AbortStateMsg("Deletion requested for miner with power still remaining")
	}

	Release(rt, h, st)

	ownerAddr, workerAddr := vmr.RT_GetMinerAccountsAssert(rt, minerAddr)
	rt.ValidateImmediateCallerInSet([]addr.Address{ownerAddr, workerAddr})

	a._rtDeleteMinerActor(rt, minerAddr)
}

func (a *StoragePowerActorCode_I) OnSectorProveCommit(rt Runtime, storageWeightDesc SectorStorageWeightDesc) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorTerminate(
	rt Runtime, storageWeightDesc SectorStorageWeightDesc, terminationType SectorTerminationType) {

	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()
	a._rtdeductClaimedPowerForSectorAssert(rt, minerAddr, storageWeightDesc)

	if terminationType != SectorTerminationType_NormalExpiration {
		amountToSlash := rt.CurrIndices().StoragePower_PledgeSlashForSectorTermination(storageWeightDesc, terminationType)
		a._rtSlashPledgeCollateral(rt, minerAddr, amountToSlash)
	}
}

func (a *StoragePowerActorCode_I) OnSectorTemporaryFaultEffectiveBegin(rt Runtime, storageWeightDesc SectorStorageWeightDesc) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	a._rtdeductClaimedPowerForSectorAssert(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorTemporaryFaultEffectiveEnd(rt Runtime, storageWeightDesc SectorStorageWeightDesc) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorModifyWeightDesc(
	rt Runtime, storageWeightDescPrev SectorStorageWeightDesc, storageWeightDescNew SectorStorageWeightDesc) {

	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	a._rtdeductClaimedPowerForSectorAssert(rt, rt.ImmediateCaller(), storageWeightDescPrev)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDescNew)
}

func (a *StoragePowerActorCode_I) OnMinerSurprisePoStSuccess(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)
	delete(st.Impl().PoStDetectedFaultMiners_, minerAddr)
	st._updatePowerEntriesFromClaimedPower(minerAddr)
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) OnMinerSurprisePoStFailure(rt Runtime, numConsecutiveFailures int) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	st.Impl().PoStDetectedFaultMiners_[minerAddr] = true
	st._updatePowerEntriesFromClaimedPower(minerAddr)

	minerClaimedPower, ok := st.ClaimedPower()[minerAddr]
	Assert(ok)

	UpdateRelease(rt, h, st)

	if numConsecutiveFailures > indices.StoragePower_SurprisePoStMaxConsecutiveFailures() {
		a._rtDeleteMinerActor(rt, minerAddr)
	} else {
		amountToSlash := rt.CurrIndices().StoragePower_PledgeSlashForSurprisePoStFailure(minerClaimedPower, numConsecutiveFailures)
		a._rtSlashPledgeCollateral(rt, minerAddr, amountToSlash)
	}
}

func (a *StoragePowerActorCode_I) ReportConsensusFault(rt Runtime, slasherAddr addr.Address, faultType ConsensusFaultType, proof []block.Block) {
	TODO()
	panic("")
	// TODO: The semantics here are quite delicate:
	//
	// - (proof []block.Block) can't be validated in isolation; we must query the runtime to confirm
	//   that at least one of the blocks provided actually appeared in the current chain.
	// - We must prevent duplicate slashes on the same offense, taking into account that the blocks
	//   may appear in different orders.
	// - We must determine how to reward multiple reporters of the same fault within a single epoch.
	//
	// Deferring to followup after these security/mechanism design questions have been resolved.
	// Previous notes:
	//
	// Use EC's IsValidConsensusFault method to validate the proof
	// slash block miner's pledge collateral
	// reward slasher
	//
	// include ReportUncommittedPowerFault(cheaterAddr addr.Address, numSectors util.UVarint) as case
	// Quite a bit more straightforward since only called by the cron actor (ie publicly verified)
	// slash cheater pledge collateral accordingly based on num sectors faulted
}

// Called by Cron.
func (a *StoragePowerActorCode_I) OnEpochTickEnd(rt Runtime) {
	rt.ValidateImmediateCallerIs(addr.CronActorAddr)

	a._rtInitiateNewSurprisePoStChallenges(rt)
	a._rtProcessDeferredCronEvents(rt)
}

func (a *StoragePowerActorCode_I) Constructor(rt Runtime) {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	h := rt.AcquireState()

	st := &StoragePowerActorState_I{
		TotalNetworkPower_:        block.StoragePower(0),
		PowerTable_:               PowerTableHAMT_Empty(),
		EscrowTable_:              autil.BalanceTableHAMT_Empty(),
		CachedDeferredCronEvents_: MinerEventsHAMT_Empty(),
		PoStDetectedFaultMiners_:  autil.MinerSetHAMT_Empty(),
		ClaimedPower_:             PowerTableHAMT_Empty(),
		NominalPower_:             PowerTableHAMT_Empty(),
		NumMinersMeetingMinPower_: 0,
	}

	UpdateRelease(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) _rtAddPowerForSector(rt Runtime, minerAddr addr.Address, storageWeightDesc SectorStorageWeightDesc) {
	h, st := a.State(rt)
	st._addClaimedPowerForSector(minerAddr, storageWeightDesc)
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) _rtdeductClaimedPowerForSectorAssert(rt Runtime, minerAddr addr.Address, storageWeightDesc SectorStorageWeightDesc) {
	h, st := a.State(rt)
	st._deductClaimedPowerForSectorAssert(minerAddr, storageWeightDesc)
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) _rtInitiateNewSurprisePoStChallenges(rt Runtime) {
	provingPeriod := indices.StorageMining_SurprisePoStProvingPeriod()

	h, st := a.State(rt)

	// sample the actor addresses
	randomness := rt.Randomness(filcrypto.DomainSeparationTag_SurprisePoStSelectMiners, rt.CurrEpoch())
	IMPL_FINISH() // BigInt arithmetic (not floating-point)
	challengeCount := math.Ceil(float64(len(st.PowerTable())) / float64(provingPeriod))
	surprisedMiners := st._selectMinersToSurprise(int(challengeCount), randomness)

	UpdateRelease(rt, h, st)

	for _, addr := range surprisedMiners {
		rt.Send(
			addr,
			ai.Method_StorageMinerActor_OnSurprisePoStChallenge,
			nil,
			actors.TokenAmount(0))
	}
}

func (a *StoragePowerActorCode_I) _rtProcessDeferredCronEvents(rt Runtime) {
	epoch := rt.CurrEpoch()

	h, st := a.State(rt)
	minerEvents, found := st.CachedDeferredCronEvents()[epoch]
	if !found {
		return
	}
	delete(st.CachedDeferredCronEvents(), epoch)
	UpdateRelease(rt, h, st)

	minerEventsRetain := []autil.MinerEvent{}
	for minerEvent := range minerEvents {
		if _, found := st.PowerTable()[minerEvent.MinerAddr()]; found {
			minerEventsRetain = append(minerEventsRetain, minerEvent)
		}
	}

	for _, minerEvent := range minerEventsRetain {
		rt.Send(
			minerEvent.MinerAddr(),
			ai.Method_StorageMinerActor_OnDeferredCronEvent,
			serde.MustSerializeParams(
				minerEvent.Sectors(),
			),
			actors.TokenAmount(0),
		)
	}
}

func (a *StoragePowerActorCode_I) _rtGetPledgeCollateralReqForMinerOrAbort(rt Runtime, minerAddr addr.Address) actors.TokenAmount {
	h, st := a.State(rt)
	minerNominalPower, found := st.NominalPower()[minerAddr]
	if !found {
		rt.AbortArgMsg("Miner not found")
	}
	Release(rt, h, st)
	return rt.CurrIndices().PledgeCollateralReq(minerNominalPower)
}

func (a *StoragePowerActorCode_I) _rtSlashPledgeCollateral(rt Runtime, minerAddr addr.Address, amountToSlash actors.TokenAmount) {
	h, st := a.State(rt)
	amountSlashed := st._slashPledgeCollateral(minerAddr, amountToSlash)
	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashed)
}

func (a *StoragePowerActorCode_I) _rtDeleteMinerActor(rt Runtime, minerAddr addr.Address) {
	h, st := a.State(rt)

	delete(st.PowerTable(), minerAddr)
	delete(st.ClaimedPower(), minerAddr)
	delete(st.NominalPower(), minerAddr)
	delete(st.PoStDetectedFaultMiners(), minerAddr)

	newTable, amountSlashed, ok := autil.BalanceTable_WithExtractAll(st.EscrowTable(), minerAddr)
	Assert(ok)
	newTable, ok = autil.BalanceTable_WithDeletedAddressEntry(newTable, minerAddr)
	Assert(ok)
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	rt.DeleteActor(minerAddr)
	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashed)
}
