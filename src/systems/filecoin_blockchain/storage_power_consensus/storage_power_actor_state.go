package storage_power_consensus

import (
	"sort"

	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	util "github.com/filecoin-project/specs/util"
)

func (st *StoragePowerActorState_I) MinerPowerMeetsConsensusMinimum(minerPower block.StoragePower) bool {
	// if miner is larger than min power requirement, we're set
	if minerPower >= node_base.MIN_MINER_SIZE_STOR {
		return true
	}

	// otherwise, if another miner meets min power requirement, return false
	if st._minersLargerThanMin() > util.UVarint(0) {
		return false
	}

	// else if none do, check whether in MIN_MINER_SIZE_TARG miners
	if len(st.PowerTable()) <= node_base.MIN_MINER_SIZE_TARG {
		// miner should pass
		return true
	}

	// get size of MIN_MINER_SIZE_TARGth largest miner
	minerSizes := make([]block.StoragePower, 0, len(st.PowerTable()))
	for _, v := range st.PowerTable() {
		minerSizes = append(minerSizes, v.Power())
	}
	sort.Slice(minerSizes, func(i, j int) bool { return int(i) > int(j) })
	return minerPower >= minerSizes[node_base.MIN_MINER_SIZE_TARG-1]
}

func (st *StoragePowerActorState_I) _slashPledgeCollateral(
	minerAddr addr.Address, slashAmountRequested actor.TokenAmount) actor.TokenAmount {

	if slashAmountRequested < 0 {
		panic("_slashPledgeCollateral: error: negative amount specified")
	}

	newTable, amountSlashed, ok := actor_util.BalanceTable_WithSubtractPreservingNonnegative(
		st.EscrowTable(), minerAddr, slashAmountRequested)
	// TODO: extra handling of not having enough pledge collateral to be slashed?
	if !ok {
		panic("_slashPledgeCollateral: error: miner address not found")
	}

	st.Impl().EscrowTable_ = newTable

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

func (st *StoragePowerActorState_I) GetSectorWeightForMiner(minerAddr addr.Address) (
	activeSectorWeight block.SectorWeight, inactiveSectorWeight block.SectorWeight, ok bool) {

	powerEntry, found := st.PowerTable()[minerAddr]
	if !found {
		return block.SectorWeight(0), block.SectorWeight(0), found
	}

	return powerEntry.ActiveSectorWeight(), powerEntry.InactiveSectorWeight(), true

}

func (st *StoragePowerActorState_I) GetCurrPledgeForMiner(minerAddr addr.Address) (currPledge actor.TokenAmount, ok bool) {
	return actor_util.BalanceTable_GetEntry(st.EscrowTable(), minerAddr)
}

func _getStorageFaultSlashPledgePercent(faultType sector.StorageFaultType) int {
	PARAM_FINISH() // TODO: instantiate these placeholders
	panic("")

	// these are the scaling constants for percentage pledge collateral to slash
	// given a miner's affected power and its total power
	switch faultType {
	case sector.DeclaredFault:
		return node_base.FAULT_SLASH_PERC_DECLARED
	case sector.DetectedFault:
		return node_base.FAULT_SLASH_PERC_DETECTED
	case sector.TerminatedFault:
		return node_base.FAULT_SLASH_PERC_TERMINATED
	default:
		panic("Case not supported")
	}
}
