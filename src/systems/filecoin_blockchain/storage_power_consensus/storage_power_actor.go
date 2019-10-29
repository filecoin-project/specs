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
	StoragePowerActor_ProcessPowerReport actor.MethodNum = 1
	StoragePowerActor_ProcessFaultReport actor.MethodNum = 2
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

// TODO: add Surprise to the chron actor
func (a *StoragePowerActorCode_I) Surprise(rt Runtime, ticket block.Ticket) []addr.Address {

	surprisedMiners := []addr.Address{}

	// The number of blocks that a challenged miner has to respond
	// TODO: this should be set in.. spa?
	var provingPeriod uint
	// The number of blocks that a challenged miner has to respond
	// TODO: this should be set in.. spa?
	// var postChallengeTime util.UInt

	// The current currBlockHeight
	// TODO: should be found in vm context
	// var currBlockHeight util.UInt

	h, st := a.State(rt)

	// The number of miners that are challenged at this block
	challengeCount := uint(len(st.PowerTable())) / provingPeriod

	// TODO: seem inefficient but spa.PowerTable() is now a map from address to power
	minerAddresses := make([]addr.Address, len(st.PowerTable()))

	// {
	// NOTE and TODO:
	// this picks challengeCount consecutive miners, starting from random offset
	// good for sampling w/o replacement, but may be bad as it would be the same cohort always challenged at the same time
	// if we want to pick properly random set, can:
	// 	make a list of indices first, shuffle it, then pick consecutive from those
	//  or pick at random every time, but skip ones we've seen already (sample w/ replacement, but skip doubles)
	// table := spa.PowerTable()
	// totalMinerCount := len(table)
	// random := 4 // TODO: get randomness from chain
	// offset := random % totalMinerCount
	// for i := 0; i < challengeCount; i++ {
	//   j := (offset + i) % totalMinerCount
	//   minerAddresses[i] = table[j]
	// }

	// return minerAddresses
	// }

	index := 0
	for address, _ := range st.PowerTable() {
		minerAddresses[index] = address
		index++
	}

	for i := uint(0); i < challengeCount; i++ {
		// TODO: randomNumber := hash(ticket, i)
		var randomNumber uint
		minerIndex := randomNumber % uint(len(st.PowerTable()))
		minerAddress := minerAddresses[minerIndex]
		surprisedMiners = append(surprisedMiners, minerAddress)
		// TODO: minerActor := GetActorFromID(actor).(storage_mining.StorageMinerActor)

		// TODO: send message to StorageMinerActor to update ProvingPeriod
		// TODO: should this update be done after surprisedMiners respond with a successful PoSt?
		// var minerActor storage_mining.StorageMinerActor
		// minerActor.ProvingPeriodEnd_ = currBlockHeight + postChallengeTime
		// SendMessage(sm.ExtendProvingPeriod)
	}

	UpdateRelease(rt, h, st)

	return surprisedMiners

}

func (a *StoragePowerActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
