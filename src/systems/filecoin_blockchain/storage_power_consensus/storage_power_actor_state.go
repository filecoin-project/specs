package storage_power_consensus

import (
	"sort"

	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
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

	// otherwise, if another miner meets min power requirement check whether this one does
	if st._minersLargerThanMin() > util.UVarint(0) {
		return minerPower >= MIN_MINER_SIZE_STOR
	}

	// else check whether in MIN_MINER_SIZE_TARG miners
	if len(st.PowerTable()) <= MIN_MINER_SIZE_TARG {
		// miner should pass
		return true
	}

	// get size of MIN_MINER_SIZE_TARGth largest miner
	minerSizes := make([]block.StoragePower, 0, len(st.PowerTable()))
	for _, v := range st.PowerTable() {
		minerSizes = append(minerSizes, v.ActivePower())
	}
	sort.Slice(minerSizes, func(i, j int) bool { return int(i) > int(j) })
	return minerPower >= minerSizes[MIN_MINER_SIZE_TARG-1]
}

func (st *StoragePowerActorState_I) _getActivePower() block.StoragePower {
	activePower := block.StoragePower(0)

	for _, miner := range st.PowerTable() {
		// only count miner power if they are larger than MIN_MINER_SIZE
		if st.ActivePowerMeetsConsensusMinimum(miner.ActivePower()) {
			activePower = activePower + miner.ActivePower()
		}
	}

	return activePower
}

func (st *StoragePowerActorState_I) _slashPledgeCollateral(
	minerAddr addr.Address, slashAmountRequested actor.TokenAmount) actor.TokenAmount {

	if slashAmountRequested < 0 {
		panic("_slashPledgeCollateral: error: negative amount specified")
	}

	newTable, amountSlashed, ok := actor.BalanceTable_WithSubtractPreservingNonnegative(
		st.EscrowTable(), minerAddr, slashAmountRequested)
	// TODO: extra handling of not having enough pledge collateral to be slashed?
	if !ok {
		panic("_slashPledgeCollateral: error: miner address not found")
	}

	st.Impl().EscrowTable_ = newTable

	return amountSlashed
}

func (st *StoragePowerActorState_I) _getPledgeCollateralReq(power block.StoragePower) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (st *StoragePowerActorState_I) _getPledgeCollateralReqForMiner(minerAddr addr.Address) actor.TokenAmount {
	minerPowerTotal, ok := st._getPowerTotalForMiner(minerAddr)
	if !ok {
		panic("Power entry not found for miner")
	}

	return st._getPledgeCollateralReq(minerPowerTotal)
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

func (st *StoragePowerActorState_I) _safeGetPowerEntry(rt Runtime, minerID addr.Address) PowerTableEntry {
	powerEntry, found := st.PowerTable()[minerID]

	if !found {
		rt.AbortStateMsg("sm._safeGetPowerEntry: miner not found in power table.")
	}

	return powerEntry
}

func (st *StoragePowerActorState_I) _getTotalPower() block.StoragePower {
	// TODO (optimization): cache this as a counter in the actor state,
	// and update it for relevant operations.

	totalPower := block.StoragePower(0)
	for _, minerEntry := range st.PowerTable() {
		totalPower = totalPower + minerEntry.ActivePower() + minerEntry.InactivePower()
	}
	return totalPower
}

func (st *StoragePowerActorState_I) _getPowerTotalForMiner(minerAddr addr.Address) (
	power block.StoragePower, ok bool) {

	IMPL_FINISH()
	panic("")
}

func (st *StoragePowerActorState_I) _getAffectedPledge(
	rt Runtime, minerAddr addr.Address, affectedPower block.StoragePower) actor.TokenAmount {

	// TODO: revisit this calculation
	minerPowerTotal, ok := st._getPowerTotalForMiner(minerAddr)
	Assert(ok)
	pledgeRequired := st._getPledgeCollateralReq(minerPowerTotal)
	affectedPledge := actor.TokenAmount(uint64(pledgeRequired) * uint64(affectedPower) / uint64(minerPowerTotal))

	return affectedPledge
}
