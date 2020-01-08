package storage_market

import (
	actor_util "github.com/filecoin-project/specs/actors/util"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	indices "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
)

////////////////////////////////////////////////////////////////////////////////
// Deal state operations
////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _updatePendingDealStates(dealIDs []deal.DealID, epoch block.ChainEpoch) (
	amountSlashedTotal actor.TokenAmount) {

	IMPL_FINISH() // BigInt arithmetic
	amountSlashedTotal = 0

	for _, dealID := range dealIDs {
		amountSlashedCurr := st._updatePendingDealState(dealID, epoch)
		amountSlashedTotal += amountSlashedCurr
	}

	return
}

func (st *StorageMarketActorState_I) _updatePendingDealState(dealID deal.DealID, epoch block.ChainEpoch) (
	amountSlashed actor.TokenAmount) {

	IMPL_FINISH() // BigInt arithmetic
	amountSlashed = 0

	deal, dealP := st._getOnChainDealAssert(dealID)

	everUpdated := (deal.LastUpdatedEpoch() != block.ChainEpoch_None)
	everSlashed := (deal.SlashEpoch() != block.ChainEpoch_None)

	Assert(!everUpdated || (deal.LastUpdatedEpoch() <= epoch))
	if deal.LastUpdatedEpoch() == epoch {
		return
	}

	if deal.SectorStartEpoch() == block.ChainEpoch_None {
		// Not yet appeared in proven sector; check for timeout.
		if dealP.StartEpoch() >= epoch {
			st._processDealInitTimedOut(dealID)
		}
		return
	}

	Assert(dealP.StartEpoch() <= epoch)

	dealEnd := dealP.EndEpoch()
	if everSlashed {
		Assert(deal.SlashEpoch() <= dealEnd)
		dealEnd = deal.SlashEpoch()
	}

	elapsedStart := dealP.StartEpoch()
	if everUpdated && deal.LastUpdatedEpoch() > elapsedStart {
		elapsedStart = deal.LastUpdatedEpoch()
	}

	elapsedEnd := dealEnd
	if epoch < elapsedEnd {
		elapsedEnd = epoch
	}

	numEpochsElapsed := elapsedEnd - elapsedStart
	st._processDealPaymentEpochsElapsed(dealID, numEpochsElapsed)

	if everSlashed {
		amountSlashed = st._processDealSlashed(dealID)
		return
	}

	if epoch >= dealP.EndEpoch() {
		st._processDealExpired(dealID)
		return
	}

	st.Deals()[dealID].Impl().LastUpdatedEpoch_ = epoch
	return
}

func (st *StorageMarketActorState_I) _deleteDeal(dealID deal.DealID) {
	_, dealP := st._getOnChainDealAssert(dealID)
	delete(st.Deals(), dealID)
	delete(st.CachedDealIDsByParty()[dealP.Provider()], dealID)
	delete(st.CachedDealIDsByParty()[dealP.Client()], dealID)
}

// Note: only processes deal payments, not deal expiration (even if the deal has expired).
func (st *StorageMarketActorState_I) _processDealPaymentEpochsElapsed(dealID deal.DealID, numEpochsElapsed block.ChainEpoch) {
	deal, dealP := st._getOnChainDealAssert(dealID)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	// Process deal payment for the elapsed epochs.
	IMPL_FINISH() // BigInt arithmetic
	totalPayment := int(numEpochsElapsed) * int(dealP.StoragePricePerEpoch())
	st._transferBalance(dealP.Client(), dealP.Provider(), actor.TokenAmount(totalPayment))
}

func (st *StorageMarketActorState_I) _processDealSlashed(dealID deal.DealID) (amountSlashed actor.TokenAmount) {
	deal, dealP := st._getOnChainDealAssert(dealID)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	slashEpoch := deal.SlashEpoch()
	Assert(slashEpoch != block.ChainEpoch_None)

	// unlock client collateral and locked storage fee
	clientCollateral := dealP.ClientCollateral()
	paymentRemaining := _dealGetPaymentRemaining(deal, slashEpoch)
	st._unlockBalance(dealP.Client(), clientCollateral+paymentRemaining)

	// slash provider collateral
	amountSlashed = dealP.ProviderCollateral()
	st._slashBalance(dealP.Provider(), amountSlashed)

	st._deleteDeal(dealID)
	return
}

// Deal start deadline elapsed without appearing in a proven sector.
// Delete deal, slash a portion of provider's collateral, and unlock remaining collaterals
// for both provider and client.
func (st *StorageMarketActorState_I) _processDealInitTimedOut(dealID deal.DealID) (amountSlashed actor.TokenAmount) {
	deal, dealP := st._getOnChainDealAssert(dealID)
	Assert(deal.SectorStartEpoch() == block.ChainEpoch_None)

	st._unlockBalance(dealP.Client(), dealP.ClientBalanceRequirement())

	amountSlashed = indices.StorageDeal_ProviderInitTimedOutSlashAmount(deal)
	amountRemaining := dealP.ProviderBalanceRequirement() - amountSlashed

	st._slashBalance(dealP.Provider(), amountSlashed)
	st._unlockBalance(dealP.Provider(), amountRemaining)

	st._deleteDeal(dealID)
	return
}

// Normal expiration. Delete deal and unlock collaterals for both miner and client.
func (st *StorageMarketActorState_I) _processDealExpired(dealID deal.DealID) {
	deal, dealP := st._getOnChainDealAssert(dealID)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	// Note: payment has already been completed at this point (_rtProcessDealPaymentEpochsElapsed)
	st._unlockBalance(dealP.Provider(), dealP.ProviderCollateral())
	st._unlockBalance(dealP.Client(), dealP.ClientCollateral())

	st._deleteDeal(dealID)
}

func (st *StorageMarketActorState_I) _generateStorageDealID(storageDeal deal.StorageDeal) deal.DealID {
	ret := st.NextID()
	st.NextID_ = st.NextID_ + deal.DealID(1)
	return ret
}

////////////////////////////////////////////////////////////////////////////////
// Balance table operations
////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _addressEntryExists(address addr.Address) bool {
	_, foundEscrow := actor_util.BalanceTable_GetEntry(st.EscrowTable(), address)
	_, foundLocked := actor_util.BalanceTable_GetEntry(st.LockedReqTable(), address)
	// Check that the tables are consistent (i.e. the address is found in one
	// if and only if it is found in the other).
	Assert(foundEscrow == foundLocked)
	return foundEscrow
}

func (st *StorageMarketActorState_I) _getTotalEscrowBalanceInternal(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor_util.BalanceTable_GetEntry(st.EscrowTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _getLockedReqBalanceInternal(a addr.Address) actor.TokenAmount {
	Assert(st._addressEntryExists(a))
	ret, ok := actor_util.BalanceTable_GetEntry(st.LockedReqTable(), a)
	Assert(ok)
	return ret
}

func (st *StorageMarketActorState_I) _lockBalanceMaybe(addr addr.Address, amount actor.TokenAmount) (
	lockBalanceOK bool) {

	Assert(amount >= 0)
	Assert(st._addressEntryExists(addr))

	prevLocked := st._getLockedReqBalanceInternal(addr)
	if prevLocked+amount > st._getTotalEscrowBalanceInternal(addr) {
		lockBalanceOK = false
		return
	}

	newLockedReqTable, ok := actor_util.BalanceTable_WithAdd(st.LockedReqTable(), addr, amount)
	Assert(ok)
	st.Impl().LockedReqTable_ = newLockedReqTable

	lockBalanceOK = true
	return
}

func (st *StorageMarketActorState_I) _unlockBalance(
	addr addr.Address, unlockAmountRequested actor.TokenAmount) {

	Assert(unlockAmountRequested >= 0)
	Assert(st._addressEntryExists(addr))

	st.Impl().LockedReqTable_ = st._tableWithDeductBalanceExact(st.LockedReqTable(), addr, unlockAmountRequested)
}

func (st *StorageMarketActorState_I) _tableWithAddBalance(
	table actor_util.BalanceTableHAMT, toAddr addr.Address, amountToAdd actor.TokenAmount) actor_util.BalanceTableHAMT {

	Assert(amountToAdd >= 0)

	newTable, ok := actor_util.BalanceTable_WithAdd(table, toAddr, amountToAdd)
	Assert(ok)
	return newTable
}

func (st *StorageMarketActorState_I) _tableWithDeductBalanceExact(
	table actor_util.BalanceTableHAMT, fromAddr addr.Address, amountRequested actor.TokenAmount) actor_util.BalanceTableHAMT {

	Assert(amountRequested >= 0)

	newTable, amountDeducted, ok := actor_util.BalanceTable_WithSubtractPreservingNonnegative(
		table, fromAddr, amountRequested)
	Assert(ok)
	Assert(amountDeducted == amountRequested)
	return newTable
}

// move funds from locked in client to available in provider
func (st *StorageMarketActorState_I) _transferBalance(
	fromAddr addr.Address, toAddr addr.Address, transferAmountRequested actor.TokenAmount) {

	Assert(transferAmountRequested >= 0)
	Assert(st._addressEntryExists(fromAddr))
	Assert(st._addressEntryExists(toAddr))

	st.Impl().EscrowTable_ = st._tableWithDeductBalanceExact(st.EscrowTable(), fromAddr, transferAmountRequested)
	st.Impl().LockedReqTable_ = st._tableWithDeductBalanceExact(st.LockedReqTable(), fromAddr, transferAmountRequested)
	st.Impl().EscrowTable_ = st._tableWithAddBalance(st.EscrowTable(), toAddr, transferAmountRequested)
}

func (st *StorageMarketActorState_I) _slashBalance(addr addr.Address, slashAmount actor.TokenAmount) {
	Assert(st._addressEntryExists(addr))
	Assert(slashAmount >= 0)

	st.Impl().EscrowTable_ = st._tableWithDeductBalanceExact(st.EscrowTable(), addr, slashAmount)
	st.Impl().LockedReqTable_ = st._tableWithDeductBalanceExact(st.LockedReqTable(), addr, slashAmount)
}

////////////////////////////////////////////////////////////////////////////////
// State utility functions
////////////////////////////////////////////////////////////////////////////////

func _rtDealProposalIsInternallyValid(rt Runtime, dealP deal.StorageDealProposal) bool {
	if dealP.EndEpoch() <= dealP.StartEpoch() {
		return false
	}

	if dealP.Duration() != dealP.EndEpoch()-dealP.StartEpoch() {
		return false
	}

	IMPL_FINISH()
	// Get signature public key of client account actor.
	var pk filcrypto.PublicKey

	IMPL_FINISH()
	// Determine which subset of DealProposal to use as the message to be signed by the client.
	var m filcrypto.Message

	// Note: we do not verify the provider signature here, since this is implicit in the
	// authenticity of the on-chain message publishing the deal.
	sig := dealP.ClientSignature()
	sigVerified := vmr.RT_VerifySignature(rt, pk, sig, m)
	if !sigVerified {
		return false
	}

	return true
}

func _dealGetPaymentRemaining(deal deal.OnChainDeal, epoch block.ChainEpoch) actor.TokenAmount {
	dealP := deal.Deal().Proposal()
	Assert(epoch <= dealP.EndEpoch())

	durationRemaining := dealP.EndEpoch() - (epoch - 1)
	Assert(durationRemaining > 0)

	IMPL_FINISH() // BigInt arithmetic
	return actor.TokenAmount(int(durationRemaining) * int(dealP.StoragePricePerEpoch()))
}

func (st *StorageMarketActorState_I) _getOnChainDeal(dealID deal.DealID) (
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

func (st *StorageMarketActorState_I) _getOnChainDealAssert(dealID deal.DealID) (
	deal deal.OnChainDeal, dealP deal.StorageDealProposal) {

	var ok bool
	deal, dealP, ok = st._getOnChainDeal(dealID)
	Assert(ok)
	return
}
