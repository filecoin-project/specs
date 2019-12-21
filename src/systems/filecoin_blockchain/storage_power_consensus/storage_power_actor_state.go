package storage_power_consensus

import (
	"sort"

	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	indices "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
	util "github.com/filecoin-project/specs/util"
)

func (st *StoragePowerActorState_I) ActivePowerMeetsConsensusMinimum(minerPower block.StoragePower) bool {
	// TODO import from consts
	MIN_MINER_SIZE_STOR := block.StoragePower(0)
	MIN_MINER_SIZE_TARG := 0

	// if miner is larger than min power requirement, we're set
	if minerPower >= MIN_MINER_SIZE_STOR {
		return true
	}

	// otherwise, if another miner meets min power requirement, return false
	if st.NumMinersMeetingMinPower() > 0 {
		return false
	}

	// else if none do, check whether in MIN_MINER_SIZE_TARG miners
	if len(st.PowerTable()) <= MIN_MINER_SIZE_TARG {
		// miner should pass
		return true
	}

	// get size of MIN_MINER_SIZE_TARGth largest miner
	minerSizes := make([]block.StoragePower, 0, len(st.PowerTable()))
	for _, v := range st.PowerTable() {
		minerSizes = append(minerSizes, v)
	}
	sort.Slice(minerSizes, func(i, j int) bool { return int(i) > int(j) })
	return minerPower >= minerSizes[MIN_MINER_SIZE_TARG-1]
}

func (st *StoragePowerActorState_I) _getActivePowerForConsensus() block.StoragePower {
	activePower := block.StoragePower(0)

	for _, minerPower := range st.PowerTable() {
		// only count miner power if they are larger than MIN_MINER_SIZE
		// (need to use either condition) in case no one meets MIN_MINER_SIZE_STOR
		if st.ActivePowerMeetsConsensusMinimum(minerPower) {
			activePower = activePower + minerPower
		}
	}

	return activePower
}

func (st *StoragePowerActorState_I) _slashPledgeCollateral(
	minerAddr addr.Address, slashAmountRequested actor.TokenAmount) actor.TokenAmount {

	Assert(slashAmountRequested >= actor.TokenAmount(0))

	newTable, amountSlashed, ok := actor_util.BalanceTable_WithSubtractPreservingNonnegative(
		st.EscrowTable(), minerAddr, slashAmountRequested)
	Assert(ok)
	st.Impl().EscrowTable_ = newTable

	// TODO: extra handling of not having enough pledge collateral to be slashed?

	return amountSlashed
}

func addrInArray(a addr.Address, list []addr.Address) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// _selectMinersToSurprise implements the PoSt-Surprise sampling algorithm
func (st *StoragePowerActorState_I) _selectMinersToSurprise(challengeCount int, randomness util.Randomness) []addr.Address {
	// this wont quite work -- a.PowerTable() is a HAMT by actor address, doesn't
	// support enumerating by int index. maybe we need that as an interface too,
	// or something similar to an iterator (or iterator over the keys)
	// or even a seeded random call directly in the HAMT: myhamt.GetRandomElement(seed []byte, idx int) using the ticket as a seed

	ptSize := len(st.PowerTable())
	allMiners := make([]addr.Address, len(st.PowerTable()))
	index := 0

	for address, _ := range st.PowerTable() {
		allMiners[index] = address
		index++
	}

	selectedMiners := make([]addr.Address, 0)
	for chall := 0; chall < challengeCount; chall++ {
		minerIndex := filcrypto.RandomInt(randomness, chall, ptSize)
		potentialChallengee := allMiners[minerIndex]
		// skip dups
		for addrInArray(potentialChallengee, selectedMiners) {
			minerIndex := filcrypto.RandomInt(randomness, chall, ptSize)
			potentialChallengee = allMiners[minerIndex]
		}
		selectedMiners = append(selectedMiners, potentialChallengee)
	}

	return selectedMiners
}

func (st *StoragePowerActorState_I) _getPowerTotalForMiner(minerAddr addr.Address) (
	power block.StoragePower, ok bool) {

	minerPower, found := st.PowerTable()[minerAddr]
	if !found {
		return block.StoragePower(0), found
	}

	return minerPower, true
}

func (st *StoragePowerActorState_I) _getCurrPledgeForMiner(minerAddr addr.Address) (currPledge actor.TokenAmount, ok bool) {
	return actor_util.BalanceTable_GetEntry(st.EscrowTable(), minerAddr)
}

func (st *StoragePowerActorState_I) _addPowerForSector(minerAddr addr.Address, storageWeightDesc SectorStorageWeightDesc) {
	// Note: The following computation does not use any of the dynamic information from CurrIndices();
	// it depends only on storageWeightDesc. This means that the power of a given storageWeightDesc
	// does not vary over time, so we can avoid continually updating it for each sector every epoch.
	//
	// The function is located in the indices module temporarily, until we find a better place for
	// global parameterization functions.
	sectorPower := indices.ConsensusPowerForStorageWeight(storageWeightDesc)

	currEntry, ok := st.NominalPower()[minerAddr]
	Assert(ok)
	st.Impl().NominalPower_[minerAddr] = currEntry + sectorPower

	st._updatePowerFromNominalPower(minerAddr)
}

func (st *StoragePowerActorState_I) _deductPowerForSectorAssert(minerAddr addr.Address, storageWeightDesc SectorStorageWeightDesc) {
	// Note: The following computation does not use any of the dynamic information from CurrIndices();
	// it depends only on storageWeightDesc. This means that the power of a given storageWeightDesc
	// does not vary over time, so we can avoid continually updating it for each sector every epoch.
	//
	// The function is located in the indices module temporarily, until we find a better place for
	// global parameterization functions.
	sectorPower := indices.ConsensusPowerForStorageWeight(storageWeightDesc)

	currentPower, ok := st.NominalPower()[minerAddr]
	Assert(ok)
	Assert(currentPower >= sectorPower)
	st.Impl().NominalPower_[minerAddr] = currentPower - sectorPower

	st._updatePowerFromNominalPower(minerAddr)
}

func (st *StoragePowerActorState_I) _updatePowerFromNominalPower(minerAddr addr.Address) {
	nominalPower, ok := st.NominalPower()[minerAddr]
	Assert(ok)

	power := nominalPower

	if nominalPower < indices.StoragePower_MinMinerPower() {
		power = 0
	}

	isPoStFailing, ok := st.PoStFailingMiners()[minerAddr]
	if isPoStFailing {
		power = 0
	}

	st._updatePowerEntryInternal(minerAddr, power)
}

func (st *StoragePowerActorState_I) _updatePowerEntryInternal(minerAddr addr.Address, updatedMinerPower block.StoragePower) {
	prevMinerPower, ok := st.PowerTable()[minerAddr]
	Assert(ok)
	st.Impl().PowerTable_[minerAddr] = updatedMinerPower
	st.Impl().TotalNetworkPower_ += (updatedMinerPower - prevMinerPower)

	consensusMinMinerPower := indices.StoragePower_MinMinerPower()
	if updatedMinerPower >= consensusMinMinerPower && prevMinerPower < consensusMinMinerPower {
		st.Impl().NumMinersMeetingMinPower_ += 1
	} else if updatedMinerPower < consensusMinMinerPower && prevMinerPower >= consensusMinMinerPower {
		st.Impl().NumMinersMeetingMinPower_ -= 1
	}
}
