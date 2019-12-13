package storage_power_consensus

import (
	"math"

	ipld "github.com/filecoin-project/specs/libraries/ipld"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH
var PARAM_FINISH = util.PARAM_FINISH
var TODO = util.TODO

const (
	Method_StoragePowerActor_EpochTick = actor.MethodPlaceholder + iota
	Method_StoragePowerActor_ProcessPowerReport
	Method_StoragePowerActor_ProcessFaultReport
	Method_StoragePowerActor_SlashPledgeForStorageFault
	Method_StoragePowerActor_EnsurePledgeCollateralSatisfied
)

func _storageFaultSlashPledgePercent(faultType sector.StorageFaultType) int {
	PARAM_FINISH() // TODO: instantiate these placeholders
	panic("")

	// these are the scaling constants for percentage pledge collateral to slash
	// given a miner's affected power and its total power
	switch faultType {
	case sector.DeclaredFault:
		return 1 // placeholder
	case sector.DetectedFault:
		return 10 // placeholder
	case sector.TerminatedFault:
		return 100 // placeholder
	default:
		panic("Case not supported")
	}
}

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type State = StoragePowerActorState

func (a *StoragePowerActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state := DeserializeState(stateBytes.As_Bytes())
	return h, state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st State) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st State) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StoragePowerActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) State {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) AddBalance(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	ownerAddr := rt.ImmediateCaller()
	msgValue := rt.ValueReceived()

	h, st := a.State(rt)
	newTable, ok := actor.BalanceTable_WithAdd(st.EscrowTable(), ownerAddr, msgValue)
	if !ok {
		rt.AbortStateMsg("spa.AddBalance: Escrow operation failed.")
	}
	st.Impl().EscrowTable_ = newTable
	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, amountRequested actor.TokenAmount) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	if amountRequested < 0 {
		rt.AbortArgMsg("spa.WithdrawBalance: negative amount.")
	}

	minerAddr := rt.ImmediateCaller()

	var ownerAddr addr.Address
	TODO() // Determine owner address from miner

	h, st := a.State(rt)

	minerPowerTotal, ok := st._getPowerTotalForMiner(minerAddr)
	if !ok {
		rt.AbortArgMsg("spa.WithdrawBalance: Miner not found.")
	}

	minBalance := st._getPledgeCollateralReq(minerPowerTotal)
	newTable, amountExtracted, ok := actor.BalanceTable_WithExtractPartial(
		st.EscrowTable(), minerAddr, amountRequested, minBalance)
	if !ok {
		rt.AbortStateMsg("spa.WithdrawBalance: Escrow operation failed.")
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	// send funds to miner
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    ownerAddr,
		Value_: amountExtracted,
	})
}

func (a *StoragePowerActorCode_I) CreateStorageMiner(
	rt Runtime, workerAddr addr.Address, peerId libp2p.PeerID) addr.Address {

	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	// ownerAddr := rt.ImmediateCaller()
	msgValue := rt.ValueReceived()

	var newMinerAddr addr.Address
	TODO() // TODO: call InitActor::Exec to create the StorageMinerActor
	panic("")

	// TODO: anything to check here?
	newMinerEntry := &PowerTableEntry_I{
		ActivePower_:   block.StoragePower(0),
		InactivePower_: block.StoragePower(0),
	}

	h, st := a.State(rt)
	newTable, ok := actor.BalanceTable_WithNewAddressEntry(st.EscrowTable(), newMinerAddr, msgValue)
	if !ok {
		panic("Internal error: newMinerAddr (result of InitActor::Exec) already exists in escrow table")
	}
	st.Impl().EscrowTable_ = newTable
	st.PowerTable()[newMinerAddr] = newMinerEntry
	UpdateRelease(rt, h, st)

	return newMinerAddr
}

func (a *StoragePowerActorCode_I) RemoveStorageMiner(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)

	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	minerPowerTotal, ok := st._getPowerTotalForMiner(minerAddr)
	if !ok {
		rt.AbortArgMsg("spa.RemoveStorageMiner: miner entry not found")
	}
	if minerPowerTotal > 0 {
		// TODO: manually remove the power entries here (and update relevant counters),
		// instead of throwing a runtime error?
		rt.AbortStateMsg("power still remains.")

		TODO()
		// TODO: also fail if funds still remaining in escrow
	}

	delete(st.PowerTable(), minerAddr)

	newTable, ok := actor.BalanceTable_WithDeletedAddressEntry(st.EscrowTable(), minerAddr)
	if !ok {
		panic("Internal error: miner entry in escrow table does not exist")
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) EnsurePledgeCollateralSatisfied(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)

	TODO() // TODO: principal access control
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)
	pledgeReq := st._getPledgeCollateralReqForMiner(minerAddr)
	UpdateRelease(rt, h, st)

	balanceSufficient, ok := actor.BalanceTable_IsEntrySufficient(st.EscrowTable(), minerAddr, pledgeReq)
	Assert(ok)
	if !balanceSufficient {
		rt.AbortFundsMsg("exitcode.InsufficientPledgeCollateral")
	}

}

// slash pledge collateral for Declared, Detected and Terminated faults
func (a *StoragePowerActorCode_I) SlashPledgeForStorageFault(rt Runtime, affectedPower block.StoragePower, faultType sector.StorageFaultType) {

	minerID := rt.ImmediateCaller()

	h, st := a.State(rt)

	affectedPledge := st._getAffectedPledge(rt, minerID, affectedPower)

	TODO() // BigInt arithmetic
	amountToSlash := actor.TokenAmount(
		_storageFaultSlashPledgePercent(faultType) * int(affectedPledge) / 100)

	amountSlashed := st._slashPledgeCollateral(rt, minerID, amountToSlash)
	UpdateRelease(rt, h, st)

	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    addr.BurntFundsActorAddr,
		Value_: amountSlashed,
	})

}

// @param PowerReport with ActivePower and InactivePower
// update miner power based on the power report
func (a *StoragePowerActorCode_I) ProcessPowerReport(rt Runtime, report PowerReport) {

	minerID := rt.ImmediateCaller()

	h, st := a.State(rt)

	powerEntry := st._safeGetPowerEntry(rt, minerID)
	powerEntry.Impl().ActivePower_ = report.ActivePower()
	powerEntry.Impl().InactivePower_ = report.InactivePower()
	st.Impl().PowerTable_[minerID] = powerEntry

	UpdateRelease(rt, h, st)
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

// Surprise is in the storage power actor because it is a singleton actor and surprise helps miners maintain power
// TODO: add Surprise to the cron actor
func (a *StoragePowerActorCode_I) Surprise(rt Runtime) {

	PROVING_PERIOD := 0 // defined in storage_mining, TODO: move constants somewhere else

	// sample the actor addresses
	h, st := a.State(rt)

	randomness := rt.Randomness(rt.CurrEpoch(), 0)
	challengeCount := math.Ceil(float64(len(st.PowerTable())) / float64(PROVING_PERIOD))
	surprisedMiners := st._selectMinersToSurprise(int(challengeCount), randomness)

	UpdateRelease(rt, h, st)

	// now send the messages
	for _, addr := range surprisedMiners {
		// For each miner here send message
		panic(addr) // hack coz of import cycle
		// rt.SendPropagatingErrors(&vmr.InvocInput_I{
		// 	To_:     addr,
		// 	Method_: sms.Method_StorageMinerActor_NotifyOfSurprisePoStChallenge,
		// })
	}
}

func (a *StoragePowerActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
