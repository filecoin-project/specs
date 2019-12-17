package storage_market

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert

func _rtAbortIfDealAlreadyProven(rt Runtime, deal deal.OnChainDeal) {
	if deal.SectorStartEpoch() != block.ChainEpoch_None {
		rt.AbortStateMsg("Deal has already appeared in proven sector.")
	}
}

func _rtAbortIfDealNotFromProvider(rt Runtime, dealP deal.StorageDealProposal, minerAddr addr.Address) {
	if !dealP.Provider().Equals(minerAddr) {
		rt.AbortStateMsg("Deal has incorrect miner as its provider.")
	}
}

func _rtAbortIfDealStartElapsed(rt Runtime, dealP deal.StorageDealProposal) {
	if rt.CurrEpoch() > dealP.StartEpoch() {
		rt.AbortStateMsg("Deal start epoch has already elapsed.")
	}
}

func _rtAbortIfDealEndElapsed(rt Runtime, dealP deal.StorageDealProposal) {
	if dealP.EndEpoch() > rt.CurrEpoch() {
		rt.AbortStateMsg("Deal end epoch has already elapsed.")
	}
}

func _rtAbortIfDealExceedsSectorLifetime(rt Runtime, dealP deal.StorageDealProposal, sectorExpiration block.ChainEpoch) {
	if dealP.EndEpoch() > sectorExpiration {
		rt.AbortStateMsg("Deal would outlive its containing sector.")
	}
}

func _rtAbortIfDealInvalidForNewSectorSeal(
	rt Runtime, minerAddr addr.Address, sectorExpiration block.ChainEpoch, deal deal.OnChainDeal) {

	dealP := deal.Deal().Proposal()

	_rtAbortIfDealNotFromProvider(rt, dealP, minerAddr)
	_rtAbortIfDealAlreadyProven(rt, deal)
	_rtAbortIfDealStartElapsed(rt, dealP)
	_rtAbortIfDealExceedsSectorLifetime(rt, dealP, sectorExpiration)
}

func _dealProposalIsInternallyValid(dealP deal.StorageDealProposal) bool {
	if dealP.EndEpoch() <= dealP.StartEpoch() {
		return false
	}
	if dealP.Duration() != dealP.EndEpoch()-dealP.StartEpoch() {
		return false
	}
	TODO() // validate client and provider signatures
	return true
}

func _dealGetPaymentRemaining(deal deal.OnChainDeal, epoch block.ChainEpoch) actor.TokenAmount {
	dealP := deal.Deal().Proposal()
	Assert(epoch <= dealP.EndEpoch())

	durationRemaining := dealP.EndEpoch() - (epoch - 1)
	Assert(durationRemaining > 0)
	// TODO: BigInt arithmetic
	return actor.TokenAmount(int(durationRemaining) * int(dealP.StoragePricePerEpoch()))
}

func (st *StorageMarketActorState_I) _rtAbortIfNewDealInvalid(rt Runtime, deal deal.StorageDeal) {
	dealP := deal.Proposal()

	if !_dealProposalIsInternallyValid(dealP) {
		rt.AbortStateMsg("Invalid deal proposal.")
	}

	_rtAbortIfDealStartElapsed(rt, dealP)
	st._rtAbortIfDealFailsParamBounds(rt, dealP)
}

func (st *StorageMarketActorState_I) _rtAbortIfDealFailsParamBounds(rt Runtime, dealP deal.StorageDealProposal) {
	TODO() // Parameterize the following bounds by global statistics (rt.Indices?)

	// minimum deal duration
	if dealP.Duration() < deal.MIN_DEAL_DURATION {
		rt.AbortStateMsg("sma._assertValidDealMinimum: deal duration shorter than minimum.")
	}

	if dealP.StoragePricePerEpoch() <= deal.MIN_DEAL_PRICE {
		rt.AbortStateMsg("sma._assertValidDealMinimum: storage price less than minimum.")
	}

	// verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
	if dealP.ProviderCollateral() < deal.MIN_PROVIDER_DEAL_COLLATERAL ||
		dealP.ClientCollateral() < deal.MIN_CLIENT_DEAL_COLLATERAL {
		rt.AbortStateMsg("sma._assertValidDealMinimum: deal collaterals less than minimum.")
	}
}

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	TODO() // use pair (minerAddr.ID, sequence number)?
	panic("")
}

func (st *StorageMarketActorState_I) _getOnchainDeal(dealID deal.DealID) (
	deal deal.OnChainDeal, dealP deal.StorageDealProposal, ok bool) {

	var found bool
	deal, found = st.Deals()[dealID]
	if found {
		dealP = deal.Deal().Proposal()
	} else {
		deal = nil
		dealP = nil
	}
	return
}

func (st *StorageMarketActorState_I) _rtGetOnChainDealOrAbort(rt Runtime, dealID deal.DealID) (deal deal.OnChainDeal, dealP deal.StorageDealProposal) {
	var found bool
	deal, dealP, found = st._getOnchainDeal(dealID)
	if !found {
		rt.AbortStateMsg("dealID not found in Deals.")
	}
	return
}

////////////////////////////////////////////////////////////////////////////////
// Balance table operations
////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _addressEntryExists(address addr.Address) bool {
	_, foundEscrow := actor.BalanceTable_GetEntry(st.EscrowTable(), address)
	_, foundLocked := actor.BalanceTable_GetEntry(st.LockedReqTable(), address)
	// Check that the tables are consistent (i.e. the address is found in one
	// if and only if it is found in the other).
	Assert(foundEscrow == foundLocked)
	return foundEscrow
}

func (st *StorageMarketActorState_I) _getTotalEscrowBalanceInternal(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor.BalanceTable_GetEntry(st.EscrowTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _getLockedReqBalanceInternal(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor.BalanceTable_GetEntry(st.LockedReqTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _rtLockBalanceUntrusted(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.AbortArgMsg("Negative amount")
	}

	if !st._addressEntryExists(addr) {
		rt.AbortArgMsg("Address does not exist in escrow table")
	}

	prevLocked := st._getLockedReqBalanceInternal(addr)
	if prevLocked+amount > st._getTotalEscrowBalanceInternal(addr) {
		rt.AbortFundsMsg("Insufficient funds available to lock.")
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
		rt.AbortFundsMsg("Attempt to deduct amount greater than present in table")
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

func (st *StorageMarketActorState_I) _rtSlashBalance(rt Runtime, addr addr.Address, slashAmount actor.TokenAmount) {
	Assert(st._addressEntryExists(addr))
	Assert(slashAmount >= 0)

	st.Impl().EscrowTable_ = st._rtTableWithDeductBalanceExact(rt, st.EscrowTable(), addr, slashAmount)
	st.Impl().LockedReqTable_ = st._rtTableWithDeductBalanceExact(rt, st.LockedReqTable(), addr, slashAmount)
}
