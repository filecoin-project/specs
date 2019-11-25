package storage_market

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
)

var TreasuryAddr addr.Address

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	// TODO
	var dealID deal.DealID
	return dealID
}

func (st *StorageMarketActorState_I) _isBalanceAvailable(a addr.Address, amount actor.TokenAmount) bool {
	bal := st.Balances()[a]
	return bal.Available() >= amount
}

func (st *StorageMarketActorState_I) _assertValidClienSignature(rt Runtime, dealP deal.StorageDealProposal) {
	// TODO: verify if we need to check provider signature
	// or it is implicit in the message

	// Optimization: make this a batch verification
	panic("TODO")
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

	// duration validation
	if p.Duration() != p.EndEpoch()-p.StartEpoch() {
		rt.Abort("sma._assertValidDealTimingAtPublish: deal duration does not match end - start.")
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

func (st *StorageMarketActorState_I) _assertDealExpireAfterMaxProveCommitWindow(rt Runtime, dealP deal.StorageDealProposal) {

	currEpoch := rt.CurrEpoch()
	dealExpiration := dealP.EndEpoch()

	if dealExpiration < (currEpoch + sector.MAX_PROVE_COMMIT_SECTOR_EPOCH) {
		rt.Abort("sma._assertDealExpireAfterMaxProveCommitWindow: deal might expire before prove commit.")
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

	p := d.Proposal()

	st._assertValidClientSignature(rt, p)
	st._assertValidDealTimingAtPublish(rt, p)
	st._assertValidDealMinimum(rt, p)
	st._assertSufficientBalanceAvailForDeal(rt, p)

	return true
}

func (st *StorageMarketActorState_I) _activateDeal(rt Runtime, deal deal.OnChainDeal) deal.OnChainDeal {

	dealP := deal.Deal().Proposal()
	deal.Impl().LastPaymentEpoch_ = dealP.StartEpoch()
	st.Deals()[deal.ID()] = deal

	return deal
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

	if currBalance.Impl().Available() < amount {
		rt.Abort("sma._lockBalance: insufficient funds available to lock.")
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

	if currBalance.Impl().Locked < amount {
		rt.Abort("sma._unlockBalance: insufficient funds to unlock.")
	}

	currBalance.Impl().Locked_ -= amount
	currBalance.Impl().Available_ += amount
}

// move funds from locked in client to available in provider
func (st *StorageMarketActorState_I) _transferBalance(rt Runtime, fromLocked addr.Address, toAvailable addr.Address, amount actor.TokenAmount) {
	if fromB == TreasuryAddr {
		toB.Impl().Available_ += amount
		return
	}

	fromB := st.Balances()[fromLocked]
	toB := st.Balances()[toAvailable]

	if fromB.Locked() < amount {
		rt.Abort("sma._transferBalance: attempt to unlock funds greater than actor has")
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

func (st *StorageMarketActorState_I) _getOnChainDeal(rt Runtime, dealID deal.DealID) deal.OnChainDeal {
	deal, found := st.Deals()[dealID]
	if !found {
		rt.Abort("sm._getOnChainDeal: dealID not found in Deals.")
	}

	return deal
}

func (st *StorageMarketActorState_I) _assertPublishedDealState(rt Runtime, dealID deal.DealID) {

	// if returns then it is on chain
	deal := st._getOnChainDeal(rt, dealID)

	// must not be active
	if deal.LastPaymentEpoch() != block.ChainEpoch(LastPaymentEpochNone) {
		rt.Abort("sma._assertPublishedDealState: deal is not in PublishedDealState.")
	}

}

func (st *StorageMarketActorState_I) _assertActiveDealState(rt Runtime, dealID deal.DealID) {

	deal := st._getOnChainDeal(rt, dealID)

	if deal.LastPaymentEpoch() == block.ChainEpoch(LastPaymentEpochNone) {
		rt.Abort("sma._assertActiveDealState: deal is not in ActiveDealState.")
	}
}

func (st *StorageMarketActorState_I) _getParticipantBalance(rt Runtime, participant addr.Address) StorageParticipantBalance {
	balance, found := st.Balances()[participant]

	if !found {
		rt.Abort("sma._getParticipantBalance: participant balance not found.")
	}

	return balance
}

func (st *StorageMarketActorState_I) _getStorageFeeSinceLastPayment(rt Runtime, deal deal.OnChainDeal, newPaymentEpoch block.ChainEpoch) actor.TokenAmount {

	duration := newPaymentEpoch - deal.LastPaymentEpoch()
	dealP := deal.Deal().Proposal()
	fee := actor.TokenAmount(0)

	if duration > 0 {
		unitPrice := dealP.StoragePricePerEpoch()
		fee := actor.TokenAmount(uint64(duration) * uint64(unitPrice))

		clientBalance := st._getParticipantBalance(rt, dealP.Client())

		if fee > clientBalance.Locked() {
			rt.Abort("sma._getStorageFeeSinceLastPayment: fee cannot exceed client LockedBalance.")
		}

	} else {
		rt.Abort("sma._getStorageFeeSinceLastPayment: no new payment since last payment.")
	}

	return fee

}

func (st *StorageMarketActorState_I) _slashDealCollateral(rt Runtime, dealP deal.StorageDealProposal) {
	amountToSlash := dealP.ProviderBalanceRequirement()

	providerBal, found := st.Balances()[dealP.Provider()]
	if !found {
		rt.Abort("sma._slashDealCollateral: provider not found in balances.")
	}

	if providerBal.Locked() < amountToSlash {
		amountToSlash = providerBal.Locked()
		// TODO: decide on error handling here
		panic("TODO")
	}

	st.Balances()[dealP.Provider()].Impl().Locked_ -= amountToSlash
	st.DealCollateralSlashed_ += amountToSlash

}

func (st *StorageMarketActorState_I) _terminateDeal(rt Runtime, dealID deal.DealID) {

	deal := st._getOnChainDeal(rt, dealID)
	st._assertActiveDealState(rt, dealID)

	dealP := deal.Deal().Proposal()
	delete(st.Deals(), dealID)

	// return client collateral and locked storage fee
	clientCollateral := dealP.TotalClientCollateral()
	lockedFee := st._getStorageFeeSinceLastPayment(rt, deal, dealP.EndEpoch())
	st._unlockBalance(rt, dealP.Client(), clientCollateral+lockedFee)

	st._slashDealCollateral(rt, dealP)
}

// delete deal from active deals
// send deal collateral to TreasuryActor
// return locked storage fee to client
// return client collateral
func (st *StorageMarketActorState_I) _slashTerminatedFaults(rt Runtime, dealIDs []deal.DealID) {

	for _, dealID := range dealIDs {
		st._terminateDeal(rt, dealID)
	}

}
