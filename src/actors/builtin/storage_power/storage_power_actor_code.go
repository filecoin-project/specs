package storage_power

import (
	"math"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	indices "github.com/filecoin-project/specs/actors/runtime/indices"
	serde "github.com/filecoin-project/specs/actors/serde"
	autil "github.com/filecoin-project/specs/actors/util"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	peer "github.com/libp2p/go-libp2p-core/peer"
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

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, minerAddr addr.Address, amountRequested abi.TokenAmount) {
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

func (a *StoragePowerActorCode_I) CreateMiner(rt Runtime, workerAddr addr.Address, sectorSize sector.SectorSize, peerId peer.ID) addr.Address {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	ownerAddr := rt.ImmediateCaller()

	newMinerAddr, err := addr.NewFromBytes(
		rt.Send(
			builtin.InitActorAddr,
			builtin.Method_InitActor_Exec,
			serde.MustSerializeParams(
				builtin.StorageMinerActorCodeID,
				ownerAddr,
				workerAddr,
				sectorSize,
				peerId,
			),
			abi.TokenAmount(0),
		).ReturnValue(),
	)
	autil.Assert(err == nil)

	h, st := a.State(rt)
	newTable, ok := autil.BalanceTable_WithNewAddressEntry(st.EscrowTable(), newMinerAddr, rt.ValueReceived())
	Assert(ok)
	st.Impl().EscrowTable_ = newTable
	st.PowerTable()[newMinerAddr] = abi.StoragePower(0)
	st.ClaimedPower()[newMinerAddr] = abi.StoragePower(0)
	st.NominalPower()[newMinerAddr] = abi.StoragePower(0)
	UpdateRelease(rt, h, st)

	return newMinerAddr
}

func (a *StoragePowerActorCode_I) DeleteMiner(rt Runtime, minerAddr addr.Address) {
	h, st := a.State(rt)

	minerPledgeBalance, ok := autil.BalanceTable_GetEntry(st.EscrowTable(), minerAddr)
	if !ok {
		rt.AbortArgMsg("Miner address not found")
	}

	if minerPledgeBalance > abi.TokenAmount(0) {
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
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorTerminate(
	rt Runtime, storageWeightDesc SectorStorageWeightDesc, terminationType SectorTerminationType) {

	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	minerAddr := rt.ImmediateCaller()
	a._rtDeductClaimedPowerForSectorAssert(rt, minerAddr, storageWeightDesc)

	if terminationType != SectorTerminationType_NormalExpiration {
		amountToSlash := rt.CurrIndices().StoragePower_PledgeSlashForSectorTermination(storageWeightDesc, terminationType)
		a._rtSlashPledgeCollateral(rt, minerAddr, amountToSlash)
	}
}

func (a *StoragePowerActorCode_I) OnSectorTemporaryFaultEffectiveBegin(rt Runtime, storageWeightDesc SectorStorageWeightDesc) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	a._rtDeductClaimedPowerForSectorAssert(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorTemporaryFaultEffectiveEnd(rt Runtime, storageWeightDesc SectorStorageWeightDesc) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDesc)
}

func (a *StoragePowerActorCode_I) OnSectorModifyWeightDesc(
	rt Runtime, storageWeightDescPrev SectorStorageWeightDesc, storageWeightDescNew SectorStorageWeightDesc) {

	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	a._rtDeductClaimedPowerForSectorAssert(rt, rt.ImmediateCaller(), storageWeightDescPrev)
	a._rtAddPowerForSector(rt, rt.ImmediateCaller(), storageWeightDescNew)
}

func (a *StoragePowerActorCode_I) OnMinerSurprisePoStSuccess(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)
	delete(st.Impl().PoStDetectedFaultMiners_, minerAddr)
	st._updatePowerEntriesFromClaimedPower(minerAddr)
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) OnMinerSurprisePoStFailure(rt Runtime, numConsecutiveFailures int) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
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

func (a *StoragePowerActorCode_I) OnMinerEnrollCronEvent(rt Runtime, eventEpoch abi.ChainEpoch, sectorNumbers []sector.SectorNumber) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	minerAddr := rt.ImmediateCaller()
	minerEvent := &autil.MinerEvent_I{
		MinerAddr_: minerAddr,
		Sectors_:   sectorNumbers,
	}

	h, st := a.State(rt)
	if _, found := st.Impl().CachedDeferredCronEvents_[eventEpoch]; !found {
		st.Impl().CachedDeferredCronEvents_[eventEpoch] = autil.MinerEventSetHAMT_Empty()
	}
	st.Impl().CachedDeferredCronEvents_[eventEpoch][minerEvent] = true
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) ReportVerifiedConsensusFault(rt Runtime, slasheeAddr addr.Address, faultEpoch abi.ChainEpoch, faultType ConsensusFaultType) {
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
	// validation checks to be done in runtime before calling this method
	// - there should be exactly two block headers in proof
	// - both blocks are mined by the same miner
	// - first block is of the same or lower block height as the second block
	//
	// Use EC's IsValidConsensusFault method to validate the proof

	// this method assumes that ConsensusFault has been checked in runtime
	slasherAddr := rt.ImmediateCaller()
	h, st := a.State(rt)

	claimedPower, powerOk := st.ClaimedPower()[slasheeAddr]
	if !powerOk {
		rt.AbortArgMsg("spa.ReportConsensusFault: miner to slash has been slashed")
	}
	Assert(claimedPower > 0)

	currPledge, pledgeOk := st._getCurrPledgeForMiner(slasheeAddr)
	if !pledgeOk {
		rt.AbortArgMsg("spa.ReportConsensusFault: miner to slash has no pledge")
	}
	Assert(currPledge > 0)

	// elapsed epoch from the latter block which committed the fault
	elapsedEpoch := rt.CurrEpoch() - faultEpoch
	if elapsedEpoch <= 0 {
		rt.AbortArgMsg("spa.ReportConsensusFault: invalid block")
	}

	collateralToSlash := st._getPledgeSlashForConsensusFault(currPledge, faultType)
	slasherReward := _getConsensusFaultSlasherReward(elapsedEpoch, collateralToSlash)

	// request slasherReward to be deducted from EscrowTable
	amountToSlasher := st._slashPledgeCollateral(slasherAddr, slasherReward)
	Assert(slasherReward == amountToSlasher)

	UpdateRelease(rt, h, st)

	// reward slasher
	rt.SendFunds(slasherAddr, amountToSlasher)

	// burn the rest of pledge collateral
	// delete miner from power table
	a._rtDeleteMinerActor(rt, slasheeAddr)
}

// Called by Cron.
func (a *StoragePowerActorCode_I) OnEpochTickEnd(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.CronActorAddr)

	a._rtInitiateNewSurprisePoStChallenges(rt)
	a._rtProcessDeferredCronEvents(rt)
}

func (a *StoragePowerActorCode_I) Constructor(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	h := rt.AcquireState()

	st := &StoragePowerActorState_I{
		TotalNetworkPower_:        abi.StoragePower(0),
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

func (a *StoragePowerActorCode_I) _rtDeductClaimedPowerForSectorAssert(rt Runtime, minerAddr addr.Address, storageWeightDesc SectorStorageWeightDesc) {
	h, st := a.State(rt)
	st._deductClaimedPowerForSectorAssert(minerAddr, storageWeightDesc)
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) _rtInitiateNewSurprisePoStChallenges(rt Runtime) {
	provingPeriod := indices.StorageMining_SurprisePoStProvingPeriod()

	h, st := a.State(rt)

	// sample the actor addresses
	minerSelectionSeed := rt.GetRandomness(rt.CurrEpoch())
	randomness := filcrypto.DeriveRandWithEpoch(filcrypto.DomainSeparationTag_SurprisePoStSelectMiners, minerSelectionSeed, int(rt.CurrEpoch()))

	IMPL_FINISH() // BigInt arithmetic (not floating-point)
	challengeCount := math.Ceil(float64(len(st.PowerTable())) / float64(provingPeriod))
	surprisedMiners := st._selectMinersToSurprise(int(challengeCount), randomness)

	UpdateRelease(rt, h, st)

	for _, addr := range surprisedMiners {
		rt.Send(
			addr,
			builtin.Method_StorageMinerActor_OnSurprisePoStChallenge,
			nil,
			abi.TokenAmount(0))
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
			builtin.Method_StorageMinerActor_OnDeferredCronEvent,
			serde.MustSerializeParams(
				minerEvent.Sectors(),
			),
			abi.TokenAmount(0),
		)
	}
}

func (a *StoragePowerActorCode_I) _rtGetPledgeCollateralReqForMinerOrAbort(rt Runtime, minerAddr addr.Address) abi.TokenAmount {
	h, st := a.State(rt)
	minerNominalPower, found := st.NominalPower()[minerAddr]
	if !found {
		rt.AbortArgMsg("Miner not found")
	}
	Release(rt, h, st)
	return rt.CurrIndices().PledgeCollateralReq(minerNominalPower)
}

func (a *StoragePowerActorCode_I) _rtSlashPledgeCollateral(rt Runtime, minerAddr addr.Address, amountToSlash abi.TokenAmount) {
	h, st := a.State(rt)
	amountSlashed := st._slashPledgeCollateral(minerAddr, amountToSlash)
	UpdateRelease(rt, h, st)

	rt.SendFunds(builtin.BurntFundsActorAddr, amountSlashed)
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

	rt.Send(
		minerAddr,
		builtin.Method_StorageMinerActor_OnDeleteMiner,
		serde.MustSerializeParams(),
		abi.TokenAmount(0),
	)

	rt.SendFunds(builtin.BurntFundsActorAddr, amountSlashed)
}
