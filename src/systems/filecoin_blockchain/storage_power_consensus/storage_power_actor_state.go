package storage_power_consensus

import (
	"sort"

	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
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
	if st._minersLargerThanMin() > util.UVarint(0) {
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

func _getStorageFaultSlashPledgePercent(faultType sector.StorageFaultType) int {
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
