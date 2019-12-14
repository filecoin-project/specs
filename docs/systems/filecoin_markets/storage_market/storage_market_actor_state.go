package storage_market

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	TODO()
	var dealID deal.DealID
	return dealID
}

func (st *StorageMarketActorState_I) _getTotalEscrowBalance(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor.BalanceTable_GetEntry(st.EscrowTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _getLockedReqBalance(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor.BalanceTable_GetEntry(st.LockedReqTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _getAvailableBalance(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	escrowBalance := st._getTotalEscrowBalance(a)
	lockedReqBalance := st._getLockedReqBalance(a)
	ret := escrowBalance - lockedReqBalance
	Assert(ret >= 0)
	return ret
}

func (st *StorageMarketActorState_I) _isBalanceAvailable(a addr.Address, amount actor.TokenAmount) bool {
	Assert(amount >= 0)
	Assert(st._addressEntryExists(a))
	availableBalance := st._getAvailableBalance(a)
	return (availableBalance >= amount)
}

func (st *StorageMarketActorState_I) _assertValidClientSignature(rt Runtime, dealP deal.StorageDealProposal) {
	// TODO: verify if we need to check provider signature
	// or it is implicit in the message

	// Optimization: make this a batch verification
	panic("TODO")
}

func (st *StorageMarketActorState_I) _assertDealStartAfterCurrEpoch(rt Runtime, p deal.StorageDealProposal) {

	currEpoch := rt.CurrEpoch()

	// deal has started before or in current epoch
	if p.StartEpoch() <= currEpoch {
		rt.AbortStateMsg("sma._assertDealStartAfterCurrEpoch: deal started before or in CurrEpoch.")
	}

}

func (st *StorageMarketActorState_I) _assertDealNotYetExpired(rt Runtime, p deal.StorageDealProposal) {
	currEpoch := rt.CurrEpoch()

	if p.EndEpoch() <= currEpoch {
		rt.AbortStateMsg("st._assertDealNotYetExpired: deal has expired.")
	}
}

func (st *StorageMarketActorState_I) _assertValidDealTimingAtPublish(rt Runtime, p deal.StorageDealProposal) {

	// TODO: verify deal did not expire when it is signed

	st._assertDealStartAfterCurrEpoch(rt, p)

	// deal ends before it starts
	if p.EndEpoch() <= p.StartEpoch() {
		rt.AbortStateMsg("sma._assertValidDealTimingAtPublish: deal ends before it starts.")
	}

	// duration validation
	if p.Duration() != p.EndEpoch()-p.StartEpoch() {
		rt.AbortStateMsg("sma._assertValidDealTimingAtPublish: deal duration does not match end - start.")
	}
}

func (st *StorageMarketActorState_I) _assertValidDealMinimum(rt Runtime, p deal.StorageDealProposal) {

	// minimum deal duration
	if p.Duration() < deal.MIN_DEAL_DURATION {
		rt.AbortStateMsg("sma._assertValidDealMinimum: deal duration shorter than minimum.")
	}

	if p.StoragePricePerEpoch() <= deal.MIN_DEAL_PRICE {
		rt.AbortStateMsg("sma._assertValidDealMinimum: storage price less than minimum.")
	}

	// verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
	if p.ProviderCollateralPerEpoch() < deal.MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH ||
		p.ClientCollateralPerEpoch() < deal.MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH {
		rt.AbortStateMsg("sma._assertValidDealMinimum: deal collaterals less than minimum.")
	}

}

func (st *StorageMarketActorState_I) _assertSufficientBalanceAvailForDeal(rt Runtime, p deal.StorageDealProposal) {

	// verify client and provider has sufficient balance
	isClientBalAvailable := st._isBalanceAvailable(p.Client(), p.ClientBalanceRequirement())
	isProviderBalAvailable := st._isBalanceAvailable(p.Provider(), p.ProviderBalanceRequirement())

	if !isClientBalAvailable || !isProviderBalAvailable {
		rt.AbortFundsMsg("sma._validateNewStorageDeal: client or provider insufficient balance.")
	}

}

func (st *StorageMarketActorState_I) _assertDealExpireAfterMaxProveCommitWindow(rt Runtime, dealP deal.StorageDealProposal) {

	currEpoch := rt.CurrEpoch()
	dealExpiration := dealP.EndEpoch()

	if dealExpiration <= (currEpoch + sector.MAX_PROVE_COMMIT_SECTOR_EPOCH) {
		rt.AbortStateMsg("sma._assertDealExpireAfterMaxProveCommitWindow: deal might expire before prove commit.")
	}

}

// Call by PublishStorageDeals
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

func (st *StorageMarketActorState_I) _rtLockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	Assert(amount >= 0)
	Assert(st._addressEntryExists(addr))

	prevLocked := st._getLockedReqBalance(addr)
	if prevLocked+amount > st._getTotalEscrowBalance(addr) {
		rt.AbortFundsMsg("sma._lockBalance: insufficient funds available to lock.")
	}

	newLockedReqTable, ok := actor.BalanceTable_WithAdd(st.LockedReqTable(), addr, amount)
	Assert(ok)
	st.Impl().LockedReqTable_ = newLockedReqTable
}

func (st *StorageMarketActorState_I) _rtUnlockBalance(
	rt Runtime, addr addr.Address, unlockAmountRequested actor.TokenAmount) {

	Assert(unlockAmountRequested >= 0)
	Assert(st._addressEntryExists(addr))

	st.Impl().LockedReqTable_ = st._rtTableWithDeductBalanceExact(rt, st.LockedReqTable(), addr, unlockAmountRequested)
}

func (st *StorageMarketActorState_I) _rtTableWithAddBalance(
	rt Runtime, table actor.BalanceTableHAMT,
	toAddr addr.Address, amountToAdd actor.TokenAmount) actor.BalanceTableHAMT {

	Assert(amountToAdd >= 0)

	newTable, ok := actor.BalanceTable_WithAdd(table, toAddr, amountToAdd)
	Assert(ok)
	return newTable
}

func (st *StorageMarketActorState_I) _rtTableWithDeductBalanceExact(
	rt Runtime, table actor.BalanceTableHAMT,
	fromAddr addr.Address, amountRequested actor.TokenAmount) actor.BalanceTableHAMT {

	Assert(amountRequested >= 0)

	newTable, amountDeducted, ok := actor.BalanceTable_WithSubtractPreservingNonnegative(
		table, fromAddr, amountRequested)
	Assert(ok)
	if amountDeducted != amountRequested {
		TODO() // Should be Assert(), as an invariant violation in SMA?
		rt.AbortFundsMsg("sma._rtDeduct: attempt to deduct amount greater than present in table")
	}
	return newTable
}

// move funds from locked in client to available in provider
func (st *StorageMarketActorState_I) _rtTransferBalance(
	rt Runtime, fromAddr addr.Address, toAddr addr.Address, transferAmountRequested actor.TokenAmount) {

	Assert(transferAmountRequested >= 0)
	Assert(st._addressEntryExists(fromAddr))
	Assert(st._addressEntryExists(toAddr))

	st.Impl().EscrowTable_ = st._rtTableWithDeductBalanceExact(rt, st.EscrowTable(), fromAddr, transferAmountRequested)
	st.Impl().LockedReqTable_ = st._rtTableWithDeductBalanceExact(rt, st.LockedReqTable(), fromAddr, transferAmountRequested)
	st.Impl().EscrowTable_ = st._rtTableWithAddBalance(rt, st.EscrowTable(), toAddr, transferAmountRequested)
}

func (st *StorageMarketActorState_I) _rtLockFundsForStorageDeal(rt Runtime, deal deal.StorageDeal) {
	p := deal.Proposal()

	st._rtLockBalance(rt, p.Client(), p.ClientBalanceRequirement())
	st._rtLockBalance(rt, p.Provider(), p.ProviderBalanceRequirement())
}

func (st *StorageMarketActorState_I) _getOnchainDeal(dealID deal.DealID) (deal deal.OnChainDeal, ok bool) {
	deal, found := st.Deals()[dealID]
	if !found {
		return nil, false
	}
	return deal, true
}

func (st *StorageMarketActorState_I) _safeGetOnChainDeal(rt Runtime, dealID deal.DealID) deal.OnChainDeal {
	deal, found := st._getOnchainDeal(dealID)
	if !found {
		rt.AbortStateMsg("sm._safeGetOnChainDeal: dealID not found in Deals.")
	}

	return deal
}

func (st *StorageMarketActorState_I) _assertPublishedDealState(rt Runtime, dealID deal.DealID) {

	// if returns then it is on chain
	deal := st._safeGetOnChainDeal(rt, dealID)

	// must not be active
	if deal.LastPaymentEpoch() != block.ChainEpoch(LastPaymentEpochNone) {
		rt.AbortStateMsg("sma._assertPublishedDealState: deal is not in PublishedDealState.")
	}

}

func (st *StorageMarketActorState_I) _assertActiveDealState(rt Runtime, dealID deal.DealID) {

	deal := st._safeGetOnChainDeal(rt, dealID)

	if deal.LastPaymentEpoch() == block.ChainEpoch(LastPaymentEpochNone) {
		rt.AbortStateMsg("sma._assertActiveDealState: deal is not in ActiveDealState.")
	}
}

func (st *StorageMarketActorState_I) _addressEntryExists(address addr.Address) bool {
	_, foundEscrow := actor.BalanceTable_GetEntry(st.EscrowTable(), address)
	_, foundLocked := actor.BalanceTable_GetEntry(st.LockedReqTable(), address)
	// Check that the tables are consistent (i.e. the address is found in one
	// if and only if it is found in the other).
	Assert(foundEscrow == foundLocked)
	return foundEscrow
}

func (st *StorageMarketActorState_I) _getStorageFeeSinceLastPayment(rt Runtime, deal deal.OnChainDeal, newPaymentEpoch block.ChainEpoch) actor.TokenAmount {

	duration := newPaymentEpoch - deal.LastPaymentEpoch()
	dealP := deal.Deal().Proposal()
	fee := actor.TokenAmount(0)

	if duration > 0 {
		unitPrice := dealP.StoragePricePerEpoch()
		fee := actor.TokenAmount(uint64(duration) * uint64(unitPrice))

		clientLockedBalance := st._getLockedReqBalance(dealP.Client())

		if fee > clientLockedBalance {
			rt.AbortFundsMsg("sma._getStorageFeeSinceLastPayment: fee cannot exceed client LockedBalance.")
		}

	} else {
		rt.AbortStateMsg("sma._getStorageFeeSinceLastPayment: no new payment since last payment.")
	}

	return fee

}

func (st *StorageMarketActorState_I) _rtSlashDealCollateral(rt Runtime, dealP deal.StorageDealProposal) actor.TokenAmount {
	Assert(st._addressEntryExists(dealP.Provider()))

	slashAmount := dealP.ProviderBalanceRequirement()
	Assert(slashAmount >= 0)

	st.Impl().EscrowTable_ = st._rtTableWithDeductBalanceExact(rt, st.EscrowTable(), dealP.Provider(), slashAmount)
	st.Impl().LockedReqTable_ = st._rtTableWithDeductBalanceExact(rt, st.LockedReqTable(), dealP.Provider(), slashAmount)

	return slashAmount
}

// delete deal from active deals
// send deal collateral to BurntFundsActor
// return locked storage fee to client
// return client collateral
func (st *StorageMarketActorState_I) _terminateDeal(rt Runtime, dealID deal.DealID) actor.TokenAmount {

	deal := st._safeGetOnChainDeal(rt, dealID)
	st._assertActiveDealState(rt, dealID)

	dealP := deal.Deal().Proposal()
	delete(st.Deals(), dealID)

	// return client collateral and locked storage fee
	clientCollateral := dealP.TotalClientCollateral()
	lockedFee := st._getStorageFeeSinceLastPayment(rt, deal, dealP.EndEpoch())
	st._rtUnlockBalance(rt, dealP.Client(), clientCollateral+lockedFee)

	return st._rtSlashDealCollateral(rt, dealP)
}

func (st *StorageMarketActorState_I) _assertEpochEqual(rt Runtime, epoch1 block.ChainEpoch, epoch2 block.ChainEpoch) {
	if epoch1 != epoch2 {
		rt.AbortArgMsg("sm._assertEpochEqual: different epochs")
	}
}

func (st *StorageMarketActorState_I) _getSectorPowerFromDeals(sectorDuration block.ChainEpoch, sectorSize block.StoragePower, dealPs []deal.StorageDealProposal) block.StoragePower {
	TODO()

	ret := block.StoragePower(0)
	return ret
}
