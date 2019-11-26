package storage_power_consensus

import (
	"math/big"

	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

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

// _sampleMinersToSurprise implements the PoSt-Surprise sampling algorithm
func (st *StoragePowerActorState_I) _sampleMinersToSurprise(rt Runtime, challengeCount int, randomness util.Randomness) []addr.Address {
	// this wont quite work -- a.PowerTable() is a HAMT by actor address, doesn't
	// support enumerating by int index. maybe we need that as an interface too,
	// or something similar to an iterator (or iterator over the keys)
	// or even a seeded random call directly in the HAMT: myhamt.GetRandomElement(seed []byte, idx int) using the ticket as a seed

	ptSize := big.NewInt(int64(len(st.PowerTable())))
	allMiners := make([]addr.Address, len(st.PowerTable()))
	index := 0

	for address, _ := range st.PowerTable() {
		allMiners[index] = address
		index++
	}

	sampledMiners := make([]addr.Address, 0)

	for chall := 0; chall < challengeCount; chall++ {
		minerIndex := filproofs.RandomInt(randomness, chall, ptSize)
		panic(minerIndex)
		// hack to turn bigint into int
		minerIndexInt := 0
		potentialChallengee := allMiners[minerIndexInt]
		// call to storage miner actor:
		// if should_challenge(lookupMinerActorStateByAddr(potentialChallengee).ShouldChallenge(rt, SURPRISE_NO_CHALLENGE_PERIOD)){
		// hack below TODO fix
		if true {
			sampledMiners = append(sampledMiners, potentialChallengee)
		}
	}

	return sampledMiners
}
