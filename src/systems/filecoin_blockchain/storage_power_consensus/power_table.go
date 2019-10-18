package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

// TODO: standardize uint type

func (pt *PowerTable_I) RegisterMiner(addr addr.Address, pk filcrypto.PubKey, sectorSize sector.SectorSize) {
	// TODO: anything else to check here?
	// redundant if sectorSize is uint
	if sectorSize <= 0 {
		// TODO: proper throw
		panic("TODO")
	}
	newMiner := &StorageMinerInfo_I{
		MinerAddress_:      addr,
		MinerStoragePower_: 0,
		MinerPK_:           pk,
		MinerSectorSize_:   sectorSize,
	}
	pt.Miners()[addr] = newMiner

	// TODO: commit state
}

func (pt *PowerTable_I) GetTotalPower() block.StoragePower {
	totalPower := block.StoragePower(0)
	for _, miner := range pt.Miners() {
		totalPower += miner.MinerStoragePower()
	}
	return totalPower
}

func (pt *PowerTable_I) IncrementPower(addr addr.Address, numSectors util.UVarint) {
	if numSectors < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	isMinerVerified := pt.verifyMiner(addr)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	powerDelta := block.StoragePower(numSectors * uint64(pt.Miners()[addr].MinerSectorSize()))
	pt.Miners()[addr].Impl().MinerStoragePower_ += powerDelta

	// TODO: commit state
}

func (pt *PowerTable_I) DecrementPower(addr addr.Address, numSectors util.UVarint) {
	if numSectors < 0 {
		// TODO: proper throw
		panic("TODO")
	}

	isMinerVerified := pt.verifyMiner(addr)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	powerDelta := block.StoragePower(numSectors * uint64(pt.Miners()[addr].MinerSectorSize()))
	if pt.Miners()[addr].Impl().MinerStoragePower_ < powerDelta {
		// TODO: proper throw
		panic("TODO")
	}

	pt.Miners()[addr].Impl().MinerStoragePower_ -= powerDelta

	// TODO: commit state
}

func (pt *PowerTable_I) RemoveMiner(addr addr.Address) {
	isMinerVerified := pt.verifyMiner(addr)
	if !isMinerVerified {
		// TODO: proper throw
		panic("TODO")
	}

	delete(pt.Miners(), addr)

	// TODO: commit state
}

func (pt *PowerTable_I) verifyMiner(addr addr.Address) bool {
	// TODO: anything else to check?
	_, found := pt.Miners()[addr]
	if !found {
		return false
	}
	return true
}

// func (pt *PowerTable_I) GetMinerPower(addr addr.Address) block.StoragePower {
// 	return pt.miners()[addr].MinerStoragePower()
// }

// func (pt *PowerTable_I) GetMinerPublicKey(addr addr.Address) filcrypto.PubKey {
// 	return pt.miners[addr].MinerPK()
// }

// must be atomic
// func (pt *PowerTable_I) SuspendPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// pt.miners[addr].MinerStoragePower -= numSectors * pt.miners[addr].minerSectorSize
// 	// pt.miners[addr].MinerSuspendedPower += numSectors * pt.miners[addr].minerSectorSize
// }

// must be atomic
// func (pt *PowerTable_I) UnsuspendPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// pt.miners[addr].MinerSuspendedPower -= numSectors * pt.miners[addr].minerSectorSize
// 	// pt.miners[addr].MinerStoragePower += numSectors * pt.miners[addr].minerSectorSize
// }

// func (pt *PowerTable_I) RemovePower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// pt.miners[addr].MinerSuspendedPower -= numSectors * pt.miners[addr].minerSectorSize
// }

// func (pt *PowerTable_I) RemoveAllPower(addr addr.Address, numSectors util.UVarint) {
// 	panic("")
// 	// pt.miners[addr].MinerStoragePower = 0
// }
