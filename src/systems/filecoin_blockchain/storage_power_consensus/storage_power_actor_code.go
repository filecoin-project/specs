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

const (
	Method_StoragePowerActor_EpochTick = actor.MethodPlaceholder + iota
	Method_StoragePowerActor_ProcessPowerReport
	Method_StoragePowerActor_ProcessFaultReport
	Method_StoragePowerActor_SlashPledgeForStorageFault
	Method_StoragePowerActor_EnsurePledgeCollateralSatisfied
)

// placeholder values
// these are the scaling constants for percentage pledge collateral to slash
// given a miner's affected power and its total power
const (
	DeclaredFaultSlashPercent   = 1   // placeholder
	DetectedFaultSlashPercent   = 10  // placeholder
	TerminatedFaultSlashPercent = 100 // placeholder
)

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

	msgValue := rt.ValueReceived()
	minerID := rt.ImmediateCaller()
	panic("TODO: fix minerID usage")

	h, st := a.State(rt)

	currEntry, found := st.PowerTable()[minerID]

	if !found {
		// AddBalance will just fail if miner is not created before hand
		rt.AbortArgMsg("minerID not found.")
	}
	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + msgValue
	st.Impl().PowerTable_[minerID] = currEntry

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, amount actor.TokenAmount) {

	if amount < 0 {
		rt.AbortArgMsg("spa.WithdrawBalance: negative amount.")
	}

	minerID := rt.ImmediateCaller()
	panic("TODO: fix minerID usage and assert caller is miner worker")

	h, st := a.State(rt)

	ret := st._ensurePledgeCollateralSatisfied(rt)
	if !ret {
		rt.AbortFundsMsg("spa.WithdrawBalance: insufficient pledge collateral.")
	}

	currEntry := st._safeGetPowerEntry(rt, minerID)
	if currEntry.AvailableBalance() < amount {
		rt.AbortFundsMsg("spa.WithdrawBalance: insufficient available balance.")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	st.Impl().PowerTable_[minerID] = currEntry

	UpdateRelease(rt, h, st)

	// send funds to miner
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    minerID,
		Value_: amount,
	})
}

func (a *StoragePowerActorCode_I) CreateStorageMiner(
	rt Runtime,
	ownerAddr addr.Address,
	workerAddr addr.Address,
	peerId libp2p.PeerID,
) addr.Address {

	// TODO: anything to check here?
	newMiner := &PowerTableEntry_I{
		ActivePower_:            block.StoragePower(0),
		InactivePower_:          block.StoragePower(0),
		AvailableBalance_:       actor.TokenAmount(0),
		LockedPledgeCollateral_: actor.TokenAmount(0),
	}

	// TODO: call constructor of StorageMinerActor
	// store ownerAddr and workerAddr there
	// and return StorageMinerActor address

	minerID := rt.ImmediateCaller()
	panic("TODO: caller is not miner actor")

	h, st := a.State(rt)

	st.PowerTable()[minerID] = newMiner

	UpdateRelease(rt, h, st)

	return minerID

}

func (a *StoragePowerActorCode_I) RemoveStorageMiner(rt Runtime, address addr.Address) {

	minerID := rt.ImmediateCaller()
	panic("TODO: use address and verify address is the caller")
	panic(minerID)

	h, st := a.State(rt)

	if (st.PowerTable()[address].ActivePower() + st.PowerTable()[address].InactivePower()) > 0 {
		rt.AbortStateMsg("power still remains.")
	}

	delete(st.PowerTable(), address)

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) GetTotalPower(rt Runtime) block.StoragePower {

	totalPower := block.StoragePower(0)

	h, st := a.State(rt)

	for _, miner := range st.PowerTable() {
		totalPower = totalPower + miner.ActivePower() + miner.InactivePower()
	}

	Release(rt, h, st)

	return totalPower
}

func (a *StoragePowerActorCode_I) EnsurePledgeCollateralSatisfied(rt Runtime) {

	h, st := a.State(rt)
	ret := st._ensurePledgeCollateralSatisfied(rt)
	UpdateRelease(rt, h, st)

	if !ret {
		rt.AbortFundsMsg("exitcode.InsufficientPledgeCollateral")
	}

}

// slash pledge collateral for Declared, Detected and Terminated faults
func (a *StoragePowerActorCode_I) SlashPledgeForStorageFault(rt Runtime, affectedPower block.StoragePower, faultType sector.StorageFaultType) {

	minerID := rt.ImmediateCaller()

	h, st := a.State(rt)

	affectedPledge := st._getAffectedPledge(rt, minerID, affectedPower)
	amountToSlash := actor.TokenAmount(0)

	switch faultType {
	case sector.DeclaredFault:
		amountToSlash = actor.TokenAmount(DeclaredFaultSlashPercent * uint64(affectedPledge) / 100)
	case sector.DetectedFault:
		amountToSlash = actor.TokenAmount(DetectedFaultSlashPercent * uint64(affectedPledge) / 100)
	case sector.TerminatedFault:
		amountToSlash = actor.TokenAmount(TerminatedFaultSlashPercent * uint64(affectedPledge) / 100)
	}

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
	SURPRISE_CHALLENGE_FREQUENCY := 0

	// sample the actor addresses
	h, st := a.State(rt)

	randomness := rt.Randomness(rt.CurrEpoch(), 0)
	challengeCount := math.Ceil(float64(SURPRISE_CHALLENGE_FREQUENCY*len(st.PowerTable())) / float64(PROVING_PERIOD))
	surprisedMiners := st._sampleMinersToSurprise(rt, int(challengeCount), randomness)

	UpdateRelease(rt, h, st)

	// now send the messages
	for _, addr := range surprisedMiners {
		// For each miner here check if they should be challenged and send message
		// rt.SendMessage(addr, ...)
		panic(addr)
	}
}

func (a *StoragePowerActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
