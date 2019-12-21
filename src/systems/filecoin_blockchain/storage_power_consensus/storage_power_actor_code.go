package storage_power_consensus

import (
	"math"

	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) AddBalance(rt Runtime, minerAddr addr.Address) {
	TODO() // update

	// caller verification
	TODO()

	msgValue := rt.ValueReceived()

	h, st := a.State(rt)
	newTable, ok := actor_util.BalanceTable_WithAdd(st.EscrowTable(), minerAddr, msgValue)
	if !ok {
		rt.AbortStateMsg("spa.AddBalance: Escrow operation failed.")
	}
	st.Impl().EscrowTable_ = newTable
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, minerAddr addr.Address, amountRequested actor.TokenAmount) {
	TODO() // update

	if amountRequested < 0 {
		rt.AbortArgMsg("spa.WithdrawBalance: negative amount.")
	}

	// caller verification
	TODO()
	msgSender := rt.ImmediateCaller()
	minBalance := a._rtGetPledgeCollateralReqForMinerOrAbort(rt, minerAddr)

	h, st := a.State(rt)
	newTable, amountExtracted, ok := actor_util.BalanceTable_WithExtractPartial(
		st.EscrowTable(), minerAddr, amountRequested, minBalance)
	if !ok {
		rt.AbortStateMsg("spa.WithdrawBalance: Escrow operation failed.")
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	// send funds to sender (pledge collateral)
	rt.SendFunds(msgSender, amountExtracted)
}

func (a *StoragePowerActorCode_I) CreateStorageMiner(
	rt Runtime, workerAddr addr.Address, peerId libp2p.PeerID) addr.Address {

	TODO() // update

	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	// ownerAddr := rt.ImmediateCaller()
	msgValue := rt.ValueReceived()

	var newMinerAddr addr.Address
	TODO() // TODO: call InitActor::Exec to create the StorageMinerActor
	panic("")

	h, st := a.State(rt)
	newTable, ok := actor_util.BalanceTable_WithNewAddressEntry(st.EscrowTable(), newMinerAddr, msgValue)
	if !ok {
		panic("Internal error: newMinerAddr (result of InitActor::Exec) already exists in escrow table")
	}
	st.Impl().EscrowTable_ = newTable
	st.PowerTable()[newMinerAddr] = block.StoragePower(0)
	UpdateRelease(rt, h, st)

	return newMinerAddr
}

func (a *StoragePowerActorCode_I) RemoveStorageMiner(rt Runtime) {
	TODO() // update

	minerAddr := rt.ImmediateCaller()

	// caller verification
	TODO()

	h, st := a.State(rt)

	activePower, inactivePower, ok := st._getPowerTotalForMiner(minerAddr)
	if !ok {
		rt.AbortArgMsg("spa.RemoveStorageMiner: miner entry not found in PowerTable.")
	}
	if activePower > 0 || inactivePower > 0 {
		rt.AbortStateMsg("spa.RemoveStorageMiner: power still remains.")
	}

	minerPledgeBalance, ok := actor_util.BalanceTable_GetEntry(st.EscrowTable(), minerAddr)
	if !ok {
		rt.AbortArgMsg("spa.RemoveStorageMiner: miner entry not found in escrow.")
	}

	if minerPledgeBalance > 0 {
		rt.AbortStateMsg("spa.RemoveStorageMiner: pledge collateral still remains.")
	}

	delete(st.PowerTable(), minerAddr)

	newTable, ok := actor_util.BalanceTable_WithDeletedAddressEntry(st.EscrowTable(), minerAddr)
	if !ok {
		panic("Internal error: miner entry in escrow table does not exist")
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)
}

// slash pledge collateral for Declared, Detected and Terminated faults
func (a *StoragePowerActorCode_I) SlashPledgeForStorageFault(rt Runtime, affectedPower block.StoragePower, faultType sector.StorageFaultType) {
	TODO() // update

	minerAddr := rt.ImmediateCaller()
	inds := rt.CurrIndices()

	// caller verification
	TODO()

	h, st := a.State(rt)

	// getAffectedPledge
	activePower, inactivePower, ok := st._getPowerTotalForMiner(minerAddr)
	Assert(ok)
	currPledge, ok := st._getCurrPledgeForMiner(minerAddr)
	Assert(ok)
	affectedPledge := inds.BlockReward_GetPledgeSlashForStorageFault(affectedPower, activePower, inactivePower, currPledge)

	TODO() // BigInt arithmetic
	amountToSlash := actor.TokenAmount(
		st._getStorageFaultSlashPledgePercent(faultType) * int(affectedPledge) / 100)

	amountSlashed := st._slashPledgeCollateral(minerAddr, amountToSlash)
	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashed)

}

func (a *StoragePowerActorCode_I) ProcessPowerReport(rt Runtime) {

	// TODO: update
	TODO()

	// minerAddr := rt.ImmediateCaller()

	// // caller verification
	// TODO()

	// powerEntry := a._rtGetPowerEntryOrAbort(rt, minerAddr)

	// h, st := a.State(rt)

	// // keep track of miners larger than minimum miner size before updating the PT
	// MIN_MINER_SIZE_STOR := block.StoragePower(0) // TODO: pull in from consts
	// if powerEntry.ActivePower() >= MIN_MINER_SIZE_STOR && report.ActivePower() < MIN_MINER_SIZE_STOR {
	// 	st.Impl()._minersLargerThanMin_ -= 1
	// }
	// if powerEntry.ActivePower() < MIN_MINER_SIZE_STOR && report.ActivePower() >= MIN_MINER_SIZE_STOR {
	// 	st.Impl()._minersLargerThanMin_ += 1
	// }

	// powerEntry.Impl().ActivePower_ = report.ActivePower()
	// powerEntry.Impl().InactivePower_ = report.InactivePower()
	// st.Impl().PowerTable_[minerAddr] = powerEntry

	// UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) ReportConsensusFault(rt Runtime, slasherAddr addr.Address, faultType ConsensusFaultType, proof []block.Block) {
	panic("TODO")

	// Use EC's IsValidConsensusFault method to validate the proof
	// slash block miner's pledge collateral
	// reward slasher

	// include ReportUncommittedPowerFault(cheaterAddr addr.Address, numSectors util.UVarint) as case
	// Quite a bit more straightforward since only called by the cron actor (ie publicly verified)
	// slash cheater pledge collateral accordingly based on num sectors faulted

}

// Called by Cron.
func (a *StoragePowerActorCode_I) OnEpochTickEnd(rt Runtime) {
	a._rtInitiateNewSurprisePoStChallenges(rt)
	a._rtProcessDeferredCronEvents(rt)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) _rtInitiateNewSurprisePoStChallenges(rt Runtime) {
	var PROVING_PERIOD block.ChainEpoch = 0 // defined in storage_mining, TODO: move constants somewhere else

	h, st := a.State(rt)

	// sample the actor addresses
	IMPL_TODO() // use randomness APIs
	randomness := rt.Randomness(rt.CurrEpoch(), 0)
	IMPL_FINISH() // BigInt arithmetic (not floating-point)
	challengeCount := math.Ceil(float64(len(st.PowerTable())) / float64(PROVING_PERIOD))
	surprisedMiners := st._selectMinersToSurprise(int(challengeCount), randomness)

	UpdateRelease(rt, h, st)

	for _, addr := range surprisedMiners {
		rt.Send(
			addr,
			ai.Method_StorageMinerActor_OnSurprisePoStChallenge,
			[]util.Serialization{},
			actor.TokenAmount(0))
	}
}

func (a *StoragePowerActorCode_I) _rtProcessDeferredCronEvents(rt Runtime) {
	epoch := rt.CurrEpoch()

	h, st := a.State(rt)

	minerEvents, found := st.CachedDeferredCronEvents()[epoch]
	if !found {
		return
	}

	for minerEvent := range minerEvents {
		rt.Send(
			minerEvent.MinerAddr(),
			ai.Method_StorageMinerActor_OnDeferredCronEvent,
			[]util.Serialization{
				sector.Serialize_SectorNumber_Array(minerEvent.Sectors()),
			},
			actor.TokenAmount(0),
		)
	}

	delete(st.CachedDeferredCronEvents(), epoch)

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) _rtGetPledgeCollateralReqForMinerOrAbort(rt Runtime, minerAddr addr.Address) actor.TokenAmount {
	inds := rt.CurrIndices()
	h, st := a.State(rt)

	activePower, inactivePower, ok := st._getPowerTotalForMiner(minerAddr)
	if !ok {
		rt.AbortArgMsg("spa._rtGetPledgeCollateralReqForMinerOrAbort: miner not in PowerTable.")
	}

	currPledge, ok := st._getCurrPledgeForMiner(minerAddr)
	if !ok {
		rt.AbortArgMsg("spa._rtGetPledgeCollateralReqForMinerOrAbort: miner not in EscrowTable.")
	}

	minBalance := inds.BlockReward_PledgeCollateralReq(activePower, inactivePower, currPledge)

	Release(rt, h, st)
	return minBalance
}

////////////////////////////////////////////////////////////////////////////////
// Dispatch table
////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	IMPL_FINISH()
	panic("TODO")
}
