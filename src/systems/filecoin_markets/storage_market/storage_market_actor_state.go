package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	// TODO
	var dealID deal.DealID
	return dealID
}

func (st *StorageMarketActorState_I) _isBalanceAvailable(a addr.Address, amount actor.TokenAmount) bool {
	bal := st.Balances()[a]
	return bal.Available() >= amount
}

func (st *StorageMarketActorState_I) _assertDealStartAfterCurrEpoch(rt Runtime, p deal.StorageDealProposal) {

	currEpoch := rt.CurrEpoch()

	// deal has started before or in current epoch
	if p.StartEpoch() <= currEpoch {
		rt.Abort("sma._assertDealStartAfterCurrEpoch: deal started before or in CurrEpoch.")
	}

}

func (st *StorageMarketActorState_I) _assertDealNotYetExpired(rt Runtime, p deal.StorageDealProposal) {
	currEpoch := rt.CurrEpoch()

	if p.EndEpoch() <= currEpoch {
		rt.Abort("st._assertDealNotYetExpired: deal has expired.")
	}
}

func (st *StorageMarketActorState_I) _assertValidDealTimingAtPublish(rt Runtime, p deal.StorageDealProposal) {

	// TODO: verify deal did not expire when it is signed

	st._assertDealStartAfterCurrEpoch(rt, p)

	// deal ends before it starts
	if p.EndEpoch() <= p.StartEpoch() {
		rt.Abort("sma._assertValidDealTimingAtPublish: deal ends before it starts.")
	}

}

func (st *StorageMarketActorState_I) _assertValidDealMinimum(rt Runtime, p deal.StorageDealProposal) {

	// minimum deal duration
	if p.Duration() < deal.MIN_DEAL_DURATION {
		rt.Abort("sma._assertValidDealMinimum: deal duration shorter than minimum.")
	}

	if p.StoragePricePerEpoch() <= deal.MIN_DEAL_PRICE {
		rt.Abort("sma._assertValidDealMinimum: storage price less than minimum.")
	}

	// verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
	if p.ProviderCollateralPerEpoch() < deal.MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH ||
		p.ClientCollateralPerEpoch() < deal.MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH {
		rt.Abort("sma._assertValidDealMinimum: deal collaterals less than minimum.")
	}

}

func (st *StorageMarketActorState_I) _assertSufficientBalanceAvailForDeal(rt Runtime, p deal.StorageDealProposal) {

	// verify client and provider has sufficient balance
	isClientBalAvailable := st._isBalanceAvailable(p.Client(), p.ClientBalanceRequirement())
	isProviderBalAvailable := st._isBalanceAvailable(p.Provider(), p.ProviderBalanceRequirement())

	if !isClientBalAvailable || !isProviderBalAvailable {
		rt.Abort("sma._validateNewStorageDeal: client or provider insufficient balance.")
	}

}

// Call by PublishStorageDeals and GetInitialUtilization (consider remove this)
// This is the check before a StorageDeal appears onchain
// It checks the following:
//   - verify deal did not expire when it is signed
//   - verify deal hits the chain before StartEpoch
//   - verify client and provider address and signature are correct (TODO may not be needed)
//   - verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
//   - verify client and provider has sufficient balance
func (st *StorageMarketActorState_I) _validateNewStorageDeal(rt Runtime, d deal.StorageDeal) bool {
	// TODO verify client and provider signature
	p := d.Proposal()

	st._assertValidDealTimingAtPublish(rt, p)
	st._assertValidDealMinimum(rt, p)
	st._assertSufficientBalanceAvailForDeal(rt, p)

	return true
}

// TODO: consider returning a boolean
func (st *StorageMarketActorState_I) _lockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("sma._lockBalance: negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("sma._lockBalance: addr not found.")
	}

	currBalance.Impl().Available_ -= amount
	currBalance.Impl().Locked_ += amount
}

func (st *StorageMarketActorState_I) _unlockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("sma._unlockBalance: negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("sma._unlockBalance: addr not found.")
	}

	currBalance.Impl().Locked_ -= amount
	currBalance.Impl().Available_ += amount
}

// move funds from locked in client to available in provider
func (st *StorageMarketActorState_I) _transferBalance(rt Runtime, fromLocked addr.Address, toAvailable addr.Address, amount actor.TokenAmount) {
	fromB := st.Balances()[fromLocked]
	toB := st.Balances()[toAvailable]

	if fromB.Locked() < amount {
		rt.Abort("sma._transferBalance: attempt to lock funds greater than actor has")
		return
	}

	fromB.Impl().Locked_ -= amount
	toB.Impl().Available_ += amount
}

func (st *StorageMarketActorState_I) _lockFundsForStorageDeal(rt Runtime, deal deal.StorageDeal) {
	p := deal.Proposal()

	st._lockBalance(rt, p.Client(), p.ClientBalanceRequirement())
	st._lockBalance(rt, p.Provider(), p.ProviderBalanceRequirement())
}

func (st *StorageMarketActorState_I) _getDeal(rt Runtime, dealID deal.DealID) deal.StorageDeal {
	deal, found := st.Deals()[dealID]
	if !found {
		rt.Abort("sm._getDeal: dealID not found in Deals.")
	}

	return deal
}

func (st *StorageMarketActorState_I) _getDealTally(rt Runtime, dealID deal.DealID) deal.StorageDealTally {
	dealTally, tallyFound := st.DealTally()[dealID]
	if !tallyFound {
		rt.Abort("sma._getDealTally: dealID not found in DealTally.")
	}

	return dealTally
}

func (st *StorageMarketActorState_I) _assertPublishedDealState(rt Runtime, dealID deal.DealID) {
	dealState := st.Impl().DealStates_.GetDealState(dealID)
	if dealState != deal.PublishedDealState {
		rt.Abort("sma._assertPublishedDealState: deal is not in PublishedDealState.")
	}
}

func (st *StorageMarketActorState_I) _assertActiveDealState(rt Runtime, dealID deal.DealID) {
	dealState := st.Impl().DealStates_.GetDealState(dealID)
	if dealState != deal.ActiveDealState {
		rt.Abort("sma._assertActiveDealState: deal is not in ActiveDealState.")
	}
}

func (st *StorageMarketActorState_I) _getStorageFeeSinceLastPayment(rt Runtime, tally deal.StorageDealTally, p deal.StorageDealProposal, lastChallenge block.ChainEpoch) actor.TokenAmount {

	duration := lastChallenge - tally.LastPaymentEpoch()
	fee := actor.TokenAmount(0)

	if duration > 0 {
		unitPrice := p.StoragePricePerEpoch()
		fee := actor.TokenAmount(uint64(duration) * uint64(unitPrice))

		if fee > tally.LockedStorageFee() {
			rt.Abort("sma._getStorageFeeSinceLastPayment: fee cannot exceed LockedStorageFee.")
		}

	} else {
		rt.Abort("sma._getStorageFeeSinceLastPayment: no new payment since last payment.")
	}

	return fee

}

// unlock remaining payments and return all UnlockedStorageFee to Provider
// remove deals from ActiveDeals
// return collaterals to both miner and client
func (st *StorageMarketActorState_I) _expireStorageDeals(rt Runtime, dealIDs []deal.DealID, lastChallengeEndEpoch block.ChainEpoch) {

	for _, dealID := range dealIDs {

		expiredDeal := st._getDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)
		dealTally := st._getDealTally(rt, dealID)
		dealP := expiredDeal.Proposal()

		fee := st._getStorageFeeSinceLastPayment(rt, dealTally, dealP, lastChallengeEndEpoch)
		dealTally.Impl().LastPaymentEpoch_ = lastChallengeEndEpoch

		// move fee from locked to unlocked in tally
		dealTally.Impl().UnlockStorageFee(fee)

		unlockedFee := dealTally.UnlockedStorageFee()
		lockedFee := dealTally.LockedStorageFee() // extra storage fee remaining

		delete(st.Deals(), dealID)
		delete(st.DealTally(), dealID)
		st.Impl().DealStates_.Clear(dealID)

		// credit UnlockedStorageFee to miner
		st._unlockBalance(rt, dealP.Provider(), unlockedFee)

		// return extra storage fee to client
		st._unlockBalance(rt, dealP.Client(), lockedFee)

		// return storage deal collaterals to both miners and client
		st._unlockBalance(rt, dealP.Provider(), dealTally.ProviderCollateralRemaining())
		st._unlockBalance(rt, dealP.Client(), dealP.TotalClientCollateral())

	}
}

func (st *StorageMarketActorState_I) _creditUnlockedFeeForProvider(rt Runtime, dealP deal.StorageDealProposal, dealTally deal.StorageDealTally) {
	// credit fund for provider
	unlockedFee := dealTally.UnlockedStorageFee()
	dealTally.Impl().UnlockedStorageFee_ = actor.TokenAmount(0)
	st._unlockBalance(rt, dealP.Provider(), unlockedFee)
}

func (st *StorageMarketActorState_I) _creditStorageDeals(rt Runtime, dealIDs []deal.DealID, lastChallengeEndEpoch block.ChainEpoch) {

	for _, dealID := range dealIDs {

		activeDeal := st._getDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)
		dealTally := st._getDealTally(rt, dealID)
		dealP := activeDeal.Proposal()

		fee := st._getStorageFeeSinceLastPayment(rt, dealTally, dealP, lastChallengeEndEpoch)
		dealTally.Impl().LastPaymentEpoch_ = lastChallengeEndEpoch

		// move fee from locked to unlocked in tally
		dealTally.Impl().UnlockStorageFee(fee)

		st._creditUnlockedFeeForProvider(rt, dealP, dealTally)
	}

}

func (st *StorageMarketActorState_I) _slashDealCollateral(rt Runtime, dealTally deal.StorageDealTally, amount actor.TokenAmount) {
	amountToSlash := amount
	if amount > dealTally.ProviderCollateralRemaining() {
		amountToSlash = dealTally.ProviderCollateralRemaining()
	}

	st.DealCollateralSlashed_ += amountToSlash
	dealTally.Impl().ProviderCollateralRemaining_ -= amountToSlash
}

func (st *StorageMarketActorState_I) _terminateDeal(rt Runtime, dealID deal.DealID) {

	deal := st._getDeal(rt, dealID)
	st._assertActiveDealState(rt, dealID)
	dealTally := st._getDealTally(rt, dealID)
	dealP := deal.Proposal()

	delete(st.Deals(), dealID)
	delete(st.DealTally(), dealID)
	st.Impl().DealStates_.Clear(dealID)

	// return client collateral and locked storage fee
	clientCollateral := dealP.TotalClientCollateral()
	lockedFee := dealTally.LockedStorageFee()
	st._unlockBalance(rt, dealP.Client(), clientCollateral+lockedFee)

	// return unlocked storage fee to provider
	unlockedFee := dealTally.UnlockedStorageFee()
	st._unlockBalance(rt, dealP.Provider(), unlockedFee)

	// slash all deal collateral
	collateralToSlash := dealTally.ProviderCollateralRemaining()
	st._slashDealCollateral(rt, dealTally, collateralToSlash)
}

// delete deal from active deals
// send deal collateral to TreasuryActor
// return locked storage fee to client
// return client collateral
// TODO: decide what to do with unlocked storage fee here
func (st *StorageMarketActorState_I) _slashTerminatedFaults(rt Runtime, dealIDs []deal.DealID) {

	for _, dealID := range dealIDs {
		st._terminateDeal(rt, dealID)
	}

}
