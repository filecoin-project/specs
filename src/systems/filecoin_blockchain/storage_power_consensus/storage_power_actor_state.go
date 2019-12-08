package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (st *StoragePowerActorState_I) ActivePowerMeetsConsensusMinimum(minerPower block.StoragePower) bool {
	totPower := st._getActivePower()

	// TODO import from consts
	MIN_MINER_SIZE_STOR := block.StoragePower(0)
	MIN_MINER_SIZE_PERC := 0

	// if miner smaller than both min size in bytes and min percentage
	if (int(minerPower)*MIN_MINER_SIZE_PERC < int(totPower)*100) && minerPower < MIN_MINER_SIZE_STOR {
		return false
	}
	return true
}

func (st *StoragePowerActorState_I) _getActivePower() block.StoragePower {
	activePower := block.StoragePower(0)

	for _, miner := range st.PowerTable() {
		activePower = activePower + miner.ActivePower()
	}

	return activePower
}

func (st *StoragePowerActorState_I) _slashPledgeCollateral(rt Runtime, minerID addr.Address, amount actor.TokenAmount) actor.TokenAmount {
	if amount < 0 {
		rt.AbortArgMsg("negative amount.")
	}

	currEntry := st._safeGetPowerEntry(rt, minerID)

	amountToSlash := amount

	if currEntry.Impl().LockedPledgeCollateral() < amount {
		amountToSlash = currEntry.Impl().LockedPledgeCollateral_
		currEntry.Impl().LockedPledgeCollateral_ = 0
		// TODO: extra handling of not having enough pledge collateral to be slashed
	} else {
		currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() - amountToSlash
	}

	st.Impl().PowerTable_[minerID] = currEntry

	return amountToSlash

}

// TODO: batch process this if possible
func (st *StoragePowerActorState_I) _lockPledgeCollateral(rt Runtime, address addr.Address, amount actor.TokenAmount) {
	// AvailableBalance -> LockedPledgeCollateral
	if amount < 0 {
		rt.AbortArgMsg("negative amount.")
	}

	minerID := rt.ImmediateCaller()
	currEntry := st._safeGetPowerEntry(rt, minerID)

	if currEntry.Impl().AvailableBalance() < amount {
		rt.AbortFundsMsg("insufficient available balance.")
	}

	currEntry.Impl().AvailableBalance_ = currEntry.AvailableBalance() - amount
	currEntry.Impl().LockedPledgeCollateral_ = currEntry.LockedPledgeCollateral() + amount
	st.Impl().PowerTable_[minerID] = currEntry
}

func (st *StoragePowerActorState_I) _unlockPledgeCollateral(rt Runtime, address addr.Address, amount actor.TokenAmount) {
	// lockedPledgeCollateral -> AvailableBalance
	if amount < 0 {
		rt.AbortArgMsg("negative amount.")
	}

	minerID := rt.ImmediateCaller()
	panic("TODO: fix minerID usage and assert caller is miner worker")

	currEntry := st._safeGetPowerEntry(rt, minerID)
	if currEntry.Impl().LockedPledgeCollateral() < amount {
		rt.AbortFundsMsg("insufficient locked balance.")
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

func addrInArray(a addr.Address, list []addr.Address) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// _selectMinersToSurprise implements the PoSt-Surprise sampling algorithm
func (st *StoragePowerActorState_I) _selectMinersToSurprise(rt Runtime, challengeCount int, randomness util.Randomness) []addr.Address {
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

func (st *StoragePowerActorState_I) _ensurePledgeCollateralSatisfied(rt Runtime) bool {

	minerID := rt.ImmediateCaller()

	powerEntry := st._safeGetPowerEntry(rt, minerID)
	pledgeCollateralRequired := st._getPledgeCollateralReq(rt, powerEntry.ActivePower()+powerEntry.InactivePower())

	if pledgeCollateralRequired < powerEntry.LockedPledgeCollateral() {
		extraLockedFund := powerEntry.LockedPledgeCollateral() - pledgeCollateralRequired
		st._unlockPledgeCollateral(rt, minerID, extraLockedFund)
		return true
	} else if pledgeCollateralRequired < (powerEntry.LockedPledgeCollateral() + powerEntry.AvailableBalance()) {
		fundToLock := pledgeCollateralRequired - powerEntry.LockedPledgeCollateral()
		st._lockPledgeCollateral(rt, minerID, fundToLock)
		return true
	}

	return false
}

func (st *StoragePowerActorState_I) _getAffectedPledge(rt Runtime, minerID addr.Address, affectedPower block.StoragePower) actor.TokenAmount {

	// TODO: revisit this calculation
	powerEntry := st._safeGetPowerEntry(rt, minerID)
	totalPower := powerEntry.ActivePower() + powerEntry.InactivePower()
	pledgeRequired := st._getPledgeCollateralReq(rt, totalPower)
	affectedPledge := actor.TokenAmount(uint64(pledgeRequired) * uint64(affectedPower) / uint64(totalPower))

	return affectedPledge
}
