package storage_power_consensus

import (
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (spa *StoragePowerActor_I) CreateStorageMiner(
	ownerAddr addr.Address,
	workerAddr addr.Address,
	peerId libp2p.PeerID,
) addr.Address {
	var msgSender addr.Address // TODO replace this

	// TODO: anything to check here?
	newMiner := &PowerTableEntry_I{
		ActivePower_:            block.StoragePower(0),
		InactivePower_:          block.StoragePower(0),
		AvailableBalance_:       actor.TokenAmount(0),
		LockedPledgeCollateral_: actor.TokenAmount(0),
	}
	spa.PowerTable()[msgSender] = newMiner

	// TODO: commit state

	// TODO: call constructor of StorageMinerActor
	// store ownerAddr and workerAddr there?
	// and return StorageMinerActor address

	var smAddress addr.Address
	return smAddress
}

// PowerTable Operation - consider remove
func (spa *StoragePowerActor_I) IncrementPower(numSectors util.UVarint, sectorSize sector.SectorSize) {
	var msgSender addr.Address // TODO replace this

	// redundant if numSectors is unsigned
	if numSectors < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	isMinerVerified := spa.verifyMiner(msgSender)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	powerDelta := block.StoragePower(numSectors * uint64(sectorSize))
	spa.PowerTable()[msgSender].Impl().ActivePower_ += powerDelta

	// TODO: commit state
}
func (spa *StoragePowerActor_I) DecrementPower(numSectors util.UVarint, sectorSize sector.SectorSize) {
	var msgSender addr.Address // TODO replace this

	if numSectors < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	isMinerVerified := spa.verifyMiner(msgSender)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	powerDelta := block.StoragePower(numSectors * uint64(sectorSize))
	if spa.PowerTable()[msgSender].Impl().ActivePower_ < powerDelta {
		// TODO: proper throw
		panic("TODO")
	}

	spa.PowerTable()[msgSender].Impl().ActivePower_ -= powerDelta

	// TODO: commit state
}

func (spa *StoragePowerActor_I) RemoveMiner(addr addr.Address) {
	isMinerVerified := spa.verifyMiner(addr)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	delete(spa.PowerTable(), addr)

	// TODO: commit state
}

func (spa *StoragePowerActor_I) verifyMiner(addr addr.Address) bool {
	// TODO: anything else to check?
	// TODO: check miner pledge collateral balances?
	// TODO: decide on what should be checked here
	_, found := spa.PowerTable()[addr]
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

func (spa *StoragePowerActor_I) GetPledgeCollateralReq(newPower block.StoragePower) actor.TokenAmount {
	// TODO: Implement
	return actor.TokenAmount(0)
}

func (spa *StoragePowerActor_I) AddBalance() {
	var msgSender addr.Address // TODO replace this
	var msgValue actor.TokenAmount

	isMinerVerified := spa.verifyMiner(msgSender)
	if !isMinerVerified {
		// TODO: proper throw
		// TODO: this might be okay, create new miner
		panic("TODO")
	}

	if msgValue < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	currEntry, found := spa.PowerTable()[msgSender]
	if found {
		currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() + msgValue
		spa.PowerTable()[msgSender] = currEntry

		// TODO: commit state change
	} else {
		newEntry := &PowerTableEntry_I{
			ActivePower_:            block.StoragePower(0),
			InactivePower_:          block.StoragePower(0),
			AvailableBalance_:       msgValue,
			LockedPledgeCollateral_: actor.TokenAmount(0),
		}
		spa.PowerTable_[msgSender] = newEntry

		// TODO: commit state change
	}

	// check that based on deals (ie sector sizes and num sectors) miner has enough associated balance in the storage miner wallet
	// pledge and associate
}

func (spa *StoragePowerActor_I) WithdrawBalance(amount actor.TokenAmount) {
	panic("TODO")
	// must check more than finality post deal expiration
	// return appropriate amount to storage market based on deals
}

func (spa *StoragePowerActor_I) lockPledgeCollateral(amount actor.TokenAmount) {
	// AvailableBalance -> LockedPledgeCollateral
	panic("TODO")
}

func (spa *StoragePowerActor_I) unlockPledgeCollateral(amount actor.TokenAmount) {
	// lockedPledgeCollateral -> AvailableBalance
	panic("TODO")
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
	numSurprised := uint(len(spa.PowerTable())) / provingPeriod

	// TODO: seem inefficient but spa.PowerTable() is now a map from address to power
	minerAddresses := make([]addr.Address, len(spa.PowerTable()))

	index := 0
	for address, _ := range spa.PowerTable() {
		minerAddresses[index] = address
		index++
	}

	for i := uint(0); i < numSurprised; i++ {
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

// func (pt *PowerTable_I) GetMinerPower(addr addr.Address) block.StoragePower {
// 	return spa.PowerTable()[addr].MinerStoragePower()
// }

// func (pt *PowerTable_I) GetMinerPublicKey(addr addr.Address) filcrypto.PubKey {
// 	return spa.PowerTable[addr].MinerPK()
// }

// must be atomic
// func (pt *PowerTable_I) SuspendPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// spa.PowerTable[addr].MinerStoragePower -= numSectors * spa.PowerTable[addr].minerSectorSize
// 	// spa.PowerTable[addr].MinerSuspendedPower += numSectors * spa.PowerTable[addr].minerSectorSize
// }

// must be atomic
// func (pt *PowerTable_I) UnsuspendPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// spa.PowerTable[addr].MinerSuspendedPower -= numSectors * spa.PowerTable[addr].minerSectorSize
// 	// spa.PowerTable[addr].MinerStoragePower += numSectors * spa.PowerTable[addr].minerSectorSize
// }

// func (pt *PowerTable_I) RemovePower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// spa.PowerTable[addr].MinerSuspendedPower -= numSectors * spa.PowerTable[addr].minerSectorSize
// }

// func (pt *PowerTable_I) RemoveAllPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// spa.PowerTable[addr].MinerStoragePower = 0
// }
