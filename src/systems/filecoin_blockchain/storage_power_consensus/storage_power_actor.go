package storage_power_consensus

import (
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

////////////////////////////////////////////////////////////////////////////////

func (spa *StoragePowerActor_I) CreateStorageMiner(
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

	// TODO: commit state

	// TODO: call constructor of StorageMinerActor
	// store ownerAddr and workerAddr there
	// and return StorageMinerActor address

	// TODO: minerID should be a MinerActorID
	// which is smaller than MinerAddress
	var minerID addr.Address
	spa.PowerTable()[minerID] = newMiner
	return minerID
}

func (spa *StoragePowerActor_I) RemoveStorageMiner(address addr.Address) {
	isMinerVerified := spa.verifyMiner(address)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: make explicit address type
	// TODO: decide if verifyMiner takes in an Address or ActorID and if perform conversion
	var minerID addr.Address

	if (spa.PowerTable()[minerID].ActivePower() + spa.PowerTable()[minerID].InactivePower()) > 0 {
		// TODO: proper throw
		panic("TODO")
	}

	delete(spa.PowerTable(), minerID)

	// TODO: commit state
}

func (spa *StoragePowerActor_I) verifyMiner(address addr.Address) bool {
	// TODO: anything else to check?
	// TODO: check miner pledge collateral balances?
	// TODO: decide on what should be checked here
	// TODO: convert address to MinerActorID

	var minerID addr.Address
	_, found := spa.PowerTable()[minerID]
	if !found {
		return false
	}
	return true
}

func (spa *StoragePowerActor_I) GetTotalPower() block.StoragePower {
	totalPower := block.StoragePower(0)
	for _, miner := range spa.PowerTable() {
		totalPower = totalPower + miner.ActivePower() + miner.InactivePower()
	}
	return totalPower
}

func (spa *StoragePowerActor_I) GetPledgeCollateralReq(power block.StoragePower) actor.TokenAmount {
	// TODO: Implement
	return actor.TokenAmount(0)
}

func (spa *StoragePowerActor_I) EnsurePledgeCollateralSatisfied() bool {
	// var msgSender addr.Address // TODO replace this
	// TODO: convert msgSender to minerID
	var minerID addr.Address

	powerEntry, found := spa.PowerTable()[minerID]

	if !found {
		// TODO: proper throw
		panic("TODO")
	}

	pledgeCollateralRequired := spa.GetPledgeCollateralReq(powerEntry.ActivePower() + powerEntry.InactivePower())

	if pledgeCollateralRequired < powerEntry.LockedPledgeCollateral() {
		return true
	}

	if pledgeCollateralRequired < (powerEntry.LockedPledgeCollateral() + powerEntry.AvailableBalance()) {
		spa.lockPledgeCollateral(minerID, (pledgeCollateralRequired - powerEntry.LockedPledgeCollateral()))

		// TODO: commit state change
		return true
	}

	return false
}

func (spa *StoragePowerActor_I) AddBalance() {
	var msgSender addr.Address // TODO replace this
	var msgValue actor.TokenAmount

	isMinerVerified := spa.verifyMiner(msgSender)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	if msgValue < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert msgSender to MinerActorID
	// if not possible, MinerActorID needs to be passed in
	var minerID addr.Address

	currEntry, found := spa.PowerTable()[minerID]

	if !found {
		// AddBalance will just fail if miner is not created before hand
		// TODO: proper throw
		panic("TODO")
	}
	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + msgValue
	spa.PowerTable()[minerID] = currEntry

	// TODO: commit state change
}

func (spa *StoragePowerActor_I) WithdrawBalance(amount actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this

	isMinerVerified := spa.verifyMiner(msgSender)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	if amount < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	currEntry, found := spa.PowerTable()[minerID]
	if !found {
		// TODO: proper throw
		panic("TODO")
	}

	if currEntry.AvailableBalance() < amount {
		// TODO: proper throw
		panic("TODO")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	spa.PowerTable_[minerID] = currEntry

	// TODO: send funds to msgSender
	// TODO: commit state change

}

func (spa *StoragePowerActor_I) slashPledgeCollateral(address addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := spa.PowerTable()[minerID]
	if !found {
		// TODO: proper throw
		panic("TODO")
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
	spa.PowerTable_[minerID] = currEntry

	// TODO: commit state change
}

// TODO: batch process this if possible
func (spa *StoragePowerActor_I) lockPledgeCollateral(address addr.Address, amount actor.TokenAmount) {
	// AvailableBalance -> LockedPledgeCollateral
	// TODO: potentially unnecessary check
	if amount < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := spa.PowerTable()[minerID]
	if !found {
		// TODO: proper throw
		panic("TODO")
	}

	if currEntry.Impl().AvailableBalance() < amount {
		// TODO: proper throw cannot lock more than one has available
		panic("TODO")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() + amount
	spa.PowerTable_[minerID] = currEntry

	// TODO: commit state change
}

func (spa *StoragePowerActor_I) unlockPledgeCollateral(address addr.Address, amount actor.TokenAmount) {
	// lockedPledgeCollateral -> AvailableBalance
	if amount < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert address to MinerActorID
	var minerID addr.Address

	currEntry, found := spa.PowerTable()[minerID]
	if !found {
		// TODO: proper throw
		panic("TODO")
	}

	if currEntry.Impl().LockedPledgeCollateral() < amount {
		// TODO: proper throw cannot unlock more than one has locked
		panic("TODO")
	}

	currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() - amount
	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + amount
	spa.PowerTable_[minerID] = currEntry

	// TODO: commit state change
}

func (spa *StoragePowerActor_I) getDeclaredFaultSlash(util.UVarint) actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (spa *StoragePowerActor_I) getDetectedFaultSlash(util.UVarint) actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (spa *StoragePowerActor_I) getTerminatedFaultSlash(util.UVarint) actor.TokenAmount {
	return actor.TokenAmount(0)
}

func (spa *StoragePowerActor_I) ProcessFaultReport(report FaultReport) {
	var msgSender addr.Address // TODO replace this

	declaredFaultSlash := spa.getDeclaredFaultSlash(report.NewDeclaredFaults())
	detectedFaultSlash := spa.getDetectedFaultSlash(report.NewDetectedFaults())
	terminatedFaultSlash := spa.getTerminatedFaultSlash(report.NewTerminatedFaults())

	spa.slashPledgeCollateral(msgSender, (declaredFaultSlash + detectedFaultSlash + terminatedFaultSlash))

	// TODO: commit state change
}

func (spa *StoragePowerActor_I) ProcessPowerReport(report PowerReport) {
	var msgSender addr.Address // TODO replace this
	isMinerVerified := spa.verifyMiner(msgSender)

	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	// TODO: convert msgSender to MinerActorID
	var minerID addr.Address

	powerEntry, found := spa.PowerTable()[minerID]

	if !found {
		// TODO: proper throw
		panic("TODO")
	}
	powerEntry.Impl().ActivePower_ = report.ActivePower()
	powerEntry.Impl().InactivePower_ = report.InactivePower()
	spa.PowerTable_[minerID] = powerEntry

	// TODO: commit state change
}

func (spa *StoragePowerActor_I) ReportConsensusFault(slasherAddr addr.Address, faultType ConsensusFaultType, proof []block.Block) {
	panic("TODO")

	// Use EC's IsValidConsensusFault method to validate the proof
	// slash block miner's pledge collateral
	// reward slasher

	// include ReportUncommittedPowerFault(cheaterAddr addr.Address, numSectors util.UVarint) as case
	// Quite a bit more straightforward since only called by the cron actor (ie publicly verified)
	// slash cheater pledge collateral accordingly based on num sectors faulted

}

// TODO: add Surprise to the chron actor
func (spa *StoragePowerActor_I) Surprise(ticket block.Ticket) []addr.Address {
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

	// The number of miners that are challenged at this block
	challengeCount := uint(len(spa.PowerTable())) / provingPeriod

	// TODO: seem inefficient but spa.PowerTable() is now a map from address to power
	minerAddresses := make([]addr.Address, len(spa.PowerTable()))

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
	for address, _ := range spa.PowerTable() {
		minerAddresses[index] = address
		index++
	}

	for i := uint(0); i < challengeCount; i++ {
		// TODO: randomNumber := hash(ticket, i)
		var randomNumber uint
		minerIndex := randomNumber % uint(len(spa.PowerTable()))
		minerAddress := minerAddresses[minerIndex]
		surprisedMiners = append(surprisedMiners, minerAddress)
		// TODO: minerActor := GetActorFromID(actor).(storage_mining.StorageMinerActor)

		// TODO: send message to StorageMinerActor to update ProvingPeriod
		// TODO: should this update be done after surprisedMiners respond with a successful PoSt?
		// var minerActor storage_mining.StorageMinerActor
		// minerActor.ProvingPeriodEnd_ = currBlockHeight + postChallengeTime
		// SendMessage(sm.ExtendProvingPeriod)
	}

	return surprisedMiners

}

func (a *StoragePowerActor_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
