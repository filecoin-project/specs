package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
	storage_mining "github.com/filecoin-project/specs/systems/filecoin_mining/storage_mining"

)

const PLEDGE_COLLATERAL_PER_GB = -1 // TODO define

// Actor
func (spa *StoragePowerActor_I) ReportConsensusFault(slasherAddr addr.Address, faultType ConsensusFaultType, proof []block.Block) {
	panic("TODO")

	// Use EC's IsValidConsensusFault method to validate the proof
	// slash block miner's pledge collateral
	// reward slasher
}

func (spa *StoragePowerActor_I) ReportUncommittedPowerFault(cheaterAddr addr.Address, numSectors util.UVarint) {
	panic("TODO")
	// Quite a bit more straightforward since only called by the cron actor (ie publicly verified)

	// slash cheater pledge collateral accordingly based on num sectors faulted
}

func (spa *StoragePowerActor_I) CommitPledgeCollateral(deals []deal.StorageDeal) {

	panic("TODO")
	// check that based on deals (ie sector sizes and num sectors) miner has enough associated balance in the storage miner wallet
	// pledge and associate
}

func (spa *StoragePowerActor_I) DecommitPledgeCollateral(deals []deal.StorageDeal) {
	panic("TODO")
	// must check more than finality post deal expiration
	// return appropriate amount to storage market based on deals
}

// TODO: add Surprise to the chron actor
func (spa *StoragePowerActor_I) Surprise(ticket block.Ticket) []addr.Address {
	surprisedMiners := []addr.Address{}

	// The number of blocks that a challenged miner has to respond
	// TODO: this should be set in.. spa?
	provingPeriod := 42
	// The number of blocks that a challenged miner has to respond
	// TODO: this should be set in.. spa?
	var postChallengeTime util.UInt
	// The current currBlockHeight
	// TODO: should be found in vm context
	var currBlockHeight util.UInt

	// The number of miners that are challenged at this block
	numSurprised := len(spa.Miners()) / provingPeriod

	for i := 0; i < numSurprised; i++ {
		// TODO: randomNumber := hash(ticket, i)
		minerIndex := 42 % len(spa.Miners())
		minerAddress := spa.Miners()[minerIndex]
		surprisedMiners = append(surprisedMiners, minerAddress)
		// TODO: minerActor := GetActorFromID(actor).(storage_mining.StorageMinerActor)

		// TODO: the following creates a cycle
		var minerActor storage_mining.StorageMinerActor
		minerActor.ProvingPeriodEnd_ = currBlockHeight + postChallengeTime
	}

	return surprisedMiners
}

// Power Table

func (pt *PowerTable_I) RegisterMiner(addr addr.Address, pk filcrypto.PubKey, sectorSize sector.SectorSize) {
	panic("")
	// newMiner := &StorageMiner_I{
	// 	MinerAddress_:        addr,
	// 	MinerStoragePower_:   0,
	// 	MinerSuspendedPower_: 0,
	// 	MinerPK_:             pk,
	// 	MinerSectorSize_:     sectorSize,
	// }
	// pt.miners[&addr] = *newMiner
}

func (pt *PowerTable_I) GetMinerPower(addr addr.Address) block.StoragePower {
	panic("")
	// return pt.miners[addr].MinerStoragePower
}

func (pt *PowerTable_I) GetTotalPower() block.StoragePower {
	panic("")
	// totalPower := 0
	// for _, miner := range pt.miners {
	// 	totalPower += miner.MinerStoragePower
	// }
	// return totalPower
}

func (pt *PowerTable_I) GetMinerPublicKey(addr addr.Address) filcrypto.PubKey {
	panic("")
	// return pt.miners[addr].MinerPK
}

func (pt *PowerTable_I) IncrementPower(addr addr.Address, numSectors util.UVarint) {
	panic("")
	// pt.miners[addr].MinerStoragePower += numSectors * pt.miners[addr].minerSectorSize
}

// must be atomic
func (pt *PowerTable_I) SuspendPower(addr addr.Address, numSectors util.UVarint) {
	panic("")
	// pt.miners[addr].MinerStoragePower -= numSectors * pt.miners[addr].minerSectorSize
	// pt.miners[addr].MinerSuspendedPower += numSectors * pt.miners[addr].minerSectorSize
}

// must be atomic
func (pt *PowerTable_I) UnsuspendPower(addr addr.Address, numSectors util.UVarint) {
	panic("")
	// pt.miners[addr].MinerSuspendedPower -= numSectors * pt.miners[addr].minerSectorSize
	// pt.miners[addr].MinerStoragePower += numSectors * pt.miners[addr].minerSectorSize
}

func (pt *PowerTable_I) RemovePower(addr addr.Address, numSectors util.UVarint) {
	panic("")
	// pt.miners[addr].MinerSuspendedPower -= numSectors * pt.miners[addr].minerSectorSize
}

func (pt *PowerTable_I) RemoveAllPower(addr addr.Address, numSectors util.UVarint) {
	panic("")
	// pt.miners[addr].MinerStoragePower = 0
}
