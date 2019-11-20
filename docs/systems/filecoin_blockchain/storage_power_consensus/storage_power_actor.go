package storage_power_consensus

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

const (
	Method_StoragePowerActor_ProcessPowerReport = actor.MethodPlaceholder
	Method_StoragePowerActor_ProcessFaultReport = actor.MethodPlaceholder
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = msg.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type State = StoragePowerActorState

func (a *StoragePowerActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.Abort("IPLD lookup error")
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

func (r *FaultReport_I) GetDeclaredFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetDetectedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (r *FaultReport_I) GetTerminatedFaultSlash() actor.TokenAmount {
	return actor.TokenAmount(0)
}

////////////////////////////////////////////////////////////////////////////////

func (st *StoragePowerActorState_I) _slashPledgeCollateral(rt Runtime, address addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := st.PowerTable()[minerID]
	if !found {
		rt.Abort("minerID not found.")
	}

	amountToSlash := amount

	if currEntry.Impl().LockedPledgeCollateral() < amount {
		amountToSlash = currEntry.Impl().LockedPledgeCollateral_
		currEntry.Impl().LockedPledgeCollateral_ = 0
		// TODO: extra handling of not having enough pledgecollateral to be slashed
	} else {
		currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() - amount
	}

	// TODO: send amountToSlash to TreasuryActor
	panic(amountToSlash)
	st.Impl().PowerTable_[minerID] = currEntry

	// TODO: commit state change
}

// TODO: batch process this if possible
func (st *StoragePowerActorState_I) _lockPledgeCollateral(rt Runtime, address addr.Address, amount actor.TokenAmount) {
	// AvailableBalance -> LockedPledgeCollateral
	// TODO: potentially unnecessary check
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := st.PowerTable()[minerID]
	if !found {
		rt.Abort("minerID not found.")
	}

	if currEntry.Impl().AvailableBalance() < amount {
		rt.Abort("insufficient available balance.")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() + amount
	st.Impl().PowerTable_[minerID] = currEntry
}

func (st *StoragePowerActorState_I) _unlockPledgeCollateral(rt Runtime, address addr.Address, amount actor.TokenAmount) {
	// lockedPledgeCollateral -> AvailableBalance
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := st.PowerTable()[minerID]
	if !found {
		rt.Abort("minerID not found.")
	}

	if currEntry.Impl().LockedPledgeCollateral() < amount {
		rt.Abort("insufficient locked balance.")
	}

	currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() - amount
	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + amount
	st.Impl().PowerTable_[minerID] = currEntry

}

func (st *StoragePowerActorState_I) _getPledgeCollateralReq(rt Runtime, power block.StoragePower) actor.TokenAmount {

	// TODO: Implement
	pcRequired := actor.TokenAmount(0)

	return pcRequired
}

func (st *StoragePowerActorState_I) _sampleMinersToSurprise(rt Runtime, challengeCount int) []addr.Address {
	// this wont quite work -- a.PowerTable() is a HAMT by actor address, doesn't
	// support enumerating by int index. maybe we need that as an interface too,
	// or something similar to an iterator (or iterator over the keys)
	// or even a seeded random call directly in the HAMT: myhamt.GetRandomElement(seed []byte, idx int)

	allMiners := make([]addr.Address, len(st.PowerTable()))
	index := 0

	for address, _ := range st.PowerTable() {
		allMiners[index] = address
		index++
	}

	return postSurpriseSample(rt, allMiners, challengeCount)
}

// postSurpriseSample implements the PoSt-Surprise sampling algorithm
func postSurpriseSample(rt Runtime, allMiners []addr.Address, challengeCount int) []addr.Address {

	sm := make([]addr.Address, challengeCount)
	for i := 0; i < challengeCount; i++ {
		// rInt := rt.NextRandomInt() // we need something like this in the runtime
		rInt := 4 // xkcd prng for now.
		miner := allMiners[rInt%len(allMiners)]
		sm = append(sm, miner)
	}

	return sm
}

////////////////////////////////////////////////////////////////////////////////

func (a *StoragePowerActorCode_I) AddBalance(rt Runtime) {

	var msgValue actor.TokenAmount

	// TODO: this should be enforced somewhere else
	if msgValue < 0 {
		rt.Abort("negative message value.")
	}

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	h, st := a.State(rt)

	currEntry, found := st.PowerTable()[minerID]

	if !found {
		// AddBalance will just fail if miner is not created before hand
		rt.Abort("minerID not found.")
	}
	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + msgValue
	st.Impl().PowerTable_[minerID] = currEntry

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) WithdrawBalance(rt Runtime, amount actor.TokenAmount) {

	if amount < 0 {
		rt.Abort("negative amount.")
	}

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	h, st := a.State(rt)

	currEntry, found := st.PowerTable()[minerID]
	if !found {
		rt.Abort("minerID not found.")
	}

	if currEntry.AvailableBalance() < amount {
		rt.Abort("insufficient balance.")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	st.Impl().PowerTable_[minerID] = currEntry

	UpdateRelease(rt, h, st)

	// TODO: send funds to msgSender
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

	// TODO: minerID should be a MinerActorID
	// which is smaller than MinerAddress
	var minerID addr.Address

	h, st := a.State(rt)

	st.PowerTable()[minerID] = newMiner

	UpdateRelease(rt, h, st)

	return minerID

}

func (a *StoragePowerActorCode_I) RemoveStorageMiner(rt Runtime, address addr.Address) {

	// TODO: make explicit address type
	var minerID addr.Address

	h, st := a.State(rt)

	if (st.PowerTable()[minerID].ActivePower() + st.PowerTable()[minerID].InactivePower()) > 0 {
		rt.Abort("power still remains.")
	}

	delete(st.PowerTable(), minerID)

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

func (a *StoragePowerActorCode_I) EnsurePledgeCollateralSatisfied(rt Runtime) bool {

	ret := false

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	h, st := a.State(rt)

	powerEntry, found := st.PowerTable()[minerID]

	if !found {
		rt.Abort("miner not found.")
	}

	pledgeCollateralRequired := st._getPledgeCollateralReq(rt, powerEntry.ActivePower()+powerEntry.InactivePower())

	if pledgeCollateralRequired < powerEntry.LockedPledgeCollateral() {
		ret = true
	} else if pledgeCollateralRequired < (powerEntry.LockedPledgeCollateral() + powerEntry.AvailableBalance()) {
		st._lockPledgeCollateral(rt, minerID, (pledgeCollateralRequired - powerEntry.LockedPledgeCollateral()))
		ret = true
	}

	UpdateRelease(rt, h, st)

	return ret
}

func (a *StoragePowerActorCode_I) ProcessFaultReport(rt Runtime, report FaultReport) {

	var msgSender addr.Address // TODO replace this

	h, st := a.State(rt)

	declaredFaultSlash := report.GetDeclaredFaultSlash()
	detectedFaultSlash := report.GetDetectedFaultSlash()
	terminatedFaultSlash := report.GetTerminatedFaultSlash()

	st._slashPledgeCollateral(rt, msgSender, (declaredFaultSlash + detectedFaultSlash + terminatedFaultSlash))

	UpdateRelease(rt, h, st)
}

func (a *StoragePowerActorCode_I) ProcessPowerReport(rt Runtime, report PowerReport) {

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	h, st := a.State(rt)

	powerEntry, found := st.PowerTable()[minerID]

	if !found {
		rt.Abort("miner not found.")
	}
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

// TODO: add Surprise to the cron actor
func (a *StoragePowerActorCode_I) Surprise(rt Runtime, ticket block.Ticket) {

	// The number of blocks that a challenged miner has to respond
	// TODO: this should be set in.. spa?
	// var postChallengeTime util.UInt

	var provingPeriod uint // TODO

	// sample the actor addresses
	h, st := a.State(rt)

	challengeCount := len(st.PowerTable()) / int(provingPeriod)
	surprisedMiners := st._sampleMinersToSurprise(rt, challengeCount)

	UpdateRelease(rt, h, st)

	// now send the messages
	for _, addr := range surprisedMiners {
		// TODO: rt.SendMessage(addr, ...)
		panic(addr)
	}
}

func (a *StoragePowerActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
