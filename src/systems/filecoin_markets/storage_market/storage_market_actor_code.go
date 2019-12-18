package storage_market

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, entryAddr addr.Address, amountRequested actor.TokenAmount) {
	IMPL_FINISH() // BigInt arithmetic
	amountSlashedTotal := actor.TokenAmount(0)

	if amountRequested < 0 {
		rt.AbortArgMsg("Negative amount.")
	}

	recipientAddr := _rtValidateImmediateCallerDetermineRecipient(rt, entryAddr)

	h, st := a.State(rt)
	st._rtAbortIfAddressEntryDoesNotExist(rt, entryAddr)

	// Before any operations that check the balance tables for funds, execute all deferred
	// deal state updates.
	//
	// Note: as an optimization, implementations may cache efficient data structures indicating
	// which of the following set of updates are redundant and can be skipped.
	amountSlashedTotal += st._rtUpdatePendingDealStatesForParty(rt, entryAddr)

	minBalance := st._getLockedReqBalanceInternal(entryAddr)
	newTable, amountExtracted, ok := actor.BalanceTable_WithExtractPartial(
		st.EscrowTable(), entryAddr, amountRequested, minBalance)
	Assert(ok)
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
	rt.SendFunds(recipientAddr, amountExtracted)
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime, entryAddr addr.Address) {
	_rtValidateImmediateCallerDetermineRecipient(rt, entryAddr)

	h, st := a.State(rt)
	st._rtAbortIfAddressEntryDoesNotExist(rt, entryAddr)

	msgValue := rt.ValueReceived()
	newTable, ok := actor.BalanceTable_WithAdd(st.EscrowTable(), entryAddr, msgValue)
	if !ok {
		// Entry not found; create implicitly.
		newTable, ok = actor.BalanceTable_WithNewAddressEntry(st.EscrowTable(), entryAddr, msgValue)
		Assert(ok)
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) PublishStorageDeals(rt Runtime, newStorageDeals []deal.StorageDeal) {
	IMPL_FINISH() // BigInt arithmetic
	amountSlashedTotal := actor.TokenAmount(0)

	// Deals may be submitted by any party (but are signed by their client and provider).
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	h, st := a.State(rt)

	// All storage deals will be added in an atomic transaction; this operation will be unrolled if any of them fails.
	for _, newDeal := range newStorageDeals {
		p := newDeal.Proposal()

		_rtAbortIfNewDealInvalid(rt, newDeal)

		// Before any operations that check the balance tables for funds, execute all deferred
		// deal state updates.
		//
		// Note: as an optimization, implementations may cache efficient data structures indicating
		// which of the following set of updates are redundant and can be skipped.
		amountSlashedTotal += st._rtUpdatePendingDealStatesForParty(rt, p.Client())
		amountSlashedTotal += st._rtUpdatePendingDealStatesForParty(rt, p.Provider())

		st._rtLockBalanceOrAbort(rt, p.Client(), p.ClientBalanceRequirement())
		st._rtLockBalanceOrAbort(rt, p.Provider(), p.ProviderBalanceRequirement())

		id := st._generateStorageDealID(newDeal)

		onchainDeal := &deal.OnChainDeal_I{
			ID_:               id,
			Deal_:             newDeal,
			SectorStartEpoch_: block.ChainEpoch_None,
		}
		st.Deals()[id] = onchainDeal
	}

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
}

// Note: in the case of a capacity-commitment sector (one with zero deals), this function should succeed vacuously.
func (a *StorageMarketActorCode_I) VerifyDealsOnSectorPreCommit(rt Runtime, dealIDs deal.DealIDs, sectorInfo sector.SectorPreCommitInfo) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		deal, _ := st._rtGetOnChainDealOrAbort(rt, dealID)
		_rtAbortIfDealInvalidForNewSectorSeal(rt, minerAddr, sectorInfo.Expiration(), deal)
	}

	Release(rt, h, st)
}

func (a *StorageMarketActorCode_I) UpdateDealsOnSectorProveCommit(rt Runtime, dealIDs deal.DealIDs, sectorInfo sector.SectorProveCommitInfo) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		deal, _ := st._rtGetOnChainDealOrAbort(rt, dealID)
		_rtAbortIfDealInvalidForNewSectorSeal(rt, minerAddr, sectorInfo.Expiration(), deal)
		st.Deals()[dealID].Impl().SectorStartEpoch_ = rt.CurrEpoch()
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) GetPieceInfosForDealIDs(rt Runtime, dealIDs deal.DealIDs) sector.PieceInfos {
	ret := []sector.PieceInfo{}

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._rtGetOnChainDealOrAbort(rt, dealID)
		ret = append(ret, sector.PieceInfo_I{
			PieceCID_: dealP.PieceCID(),
			Size_:     util.UInt(dealP.PieceSize().Total()),
		}.Ref())
	}

	Release(rt, h, st)

	return &sector.PieceInfos_I{Items_: ret}
}

func (a *StorageMarketActorCode_I) GetWeightForDealSet(rt Runtime, dealIDs deal.DealIDs) deal.DealWeight {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	IMPL_FINISH() // BigInt arithmetic
	ret := 0

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._getOnChainDealAssert(dealID)
		Assert(dealP.Provider().Equals(minerAddr))

		IMPL_FINISH() // BigInt arithmetic
		ret += int(dealP.Duration()) * int(dealP.PieceSize().Total())
	}

	UpdateRelease(rt, h, st)

	return deal.DealWeight(ret)
}

func (a *StorageMarketActorCode_I) TerminateDealsOnSlashProviderSector(rt Runtime, dealIDs deal.DealIDs) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._rtGetOnChainDealOrAbort(rt, dealID)
		Assert(dealP.Provider().Equals(minerAddr))

		// Note: we do not perform the balance transfers here, but rather simply record the flag
		// to indicate that _processDealSlashed should be called when the deferred state computation
		// is performed.
		st.Deals()[dealID].Impl().SlashEpoch_ = rt.CurrEpoch()
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) OnEpochTickEnd(rt Runtime) {
	h, st := a.State(rt)

	// Some deals may never be affected by the normal calls to _rtUpdatePendingDealStatesForParty
	// (notably, if the relevant party never checks its balance).
	// Without some cleanup mechanism, these deals may gradually accumulate and cause
	// the StorageMarketActor state to grow without bound.
	// To prevent this, we amortize the cost of this cleanup by processing a relatively
	// small number of deals every epoch, independent of the calls above.
	var cleanupDealIDs []deal.DealID
	// Populate with the N oldest deals (e.g., by a priority queue on EndEpoch).
	// N is a system parameter TBD, which may be a function of the global statistics
	// (including the number of deals in each prior epoch).
	IMPL_TODO()

	amountSlashedTotal := st._updatePendingDealStates(cleanupDealIDs, rt.CurrEpoch())

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
}

func (a *StorageMarketActorCode_I) Constructor(rt Runtime) {
	h := rt.AcquireState()

	IMPL_FINISH() // Initialize these structures
	var dealsAMTEmpty DealsAMT
	var balanceTableEmpty actor.BalanceTableHAMT

	st := &StorageMarketActorState_I{
		Deals_:          dealsAMTEmpty,
		EscrowTable_:    balanceTableEmpty,
		LockedReqTable_: balanceTableEmpty,
		NextID_:         deal.DealID(0),
	}
	UpdateRelease(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _rtUpdatePendingDealStatesForParty(rt Runtime, addr addr.Address) (
	amountSlashedTotal actor.TokenAmount) {

	// For consistency with OnEpochTickEnd, only process updates up to the end of the _previous_ epoch.
	epoch := rt.CurrEpoch() - 1

	var relevantDealIDs []deal.DealID
	// Populate with the set of all elements in st.Deals() in which addr is one of the parties
	// (either client or provider).
	//
	// Note: as an optimization, implementations may cache efficient data structures to maintain
	// this index.
	IMPL_TODO()

	amountSlashedTotal = st._updatePendingDealStates(relevantDealIDs, epoch)
	return
}

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

func (st *StorageMarketActorState_I) _rtAbortIfAddressEntryDoesNotExist(rt Runtime, entryAddr addr.Address) {
	if !st._addressEntryExists(entryAddr) {
		rt.AbortArgMsg("Address entry does not exist")
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

func _rtValidateImmediateCallerDetermineRecipient(rt Runtime, entryAddr addr.Address) addr.Address {
	if _rtIsStorageMiner(rt, entryAddr) {
		// Storage miner actor; implied funds recipient is the associated owner address.
		ownerAddr, workerAddr := _rtGetMinerAccountsAssert(rt, entryAddr)
		rt.ValidateImmediateCallerInSet([]addr.Address{ownerAddr, workerAddr})
		return ownerAddr
	} else {
		// Account actor (client); funds recipient is just the entry address itself.
		rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)
		return entryAddr
	}
}

func _rtIsStorageMiner(rt Runtime, minerAddr addr.Address) bool {
	codeID, ok := rt.GetActorCodeID(minerAddr)
	Assert(ok)
	if !codeID.IsBuiltin() {
		return false
	}
	return (codeID.As_Builtin() == actor.BuiltinActorID_StorageMiner)
}

func _rtGetMinerAccountsAssert(rt Runtime, minerAddr addr.Address) (ownerAddr addr.Address, workerAddr addr.Address) {
	ownerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetOwnerAddr, []util.Serialization{}))

	workerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetWorkerAddr, []util.Serialization{}))

	return
}

func _rtAbortIfNewDealInvalid(rt Runtime, deal deal.StorageDeal) {
	dealP := deal.Proposal()

	if !_dealProposalIsInternallyValid(dealP) {
		rt.AbortStateMsg("Invalid deal proposal.")
	}

	_rtAbortIfDealStartElapsed(rt, dealP)
	_rtAbortIfDealFailsParamBounds(rt, dealP)
}

func _rtAbortIfDealFailsParamBounds(rt Runtime, dealP deal.StorageDealProposal) {
	inds := rt.CurrIndices()

	minDuration, maxDuration := inds.StorageDeal_DurationBounds(dealP.PieceSize(), dealP.StartEpoch())
	if dealP.Duration() < minDuration || dealP.Duration() > maxDuration {
		rt.AbortStateMsg("Deal duration out of bounds.")
	}

	minPrice, maxPrice := inds.StorageDeal_StoragePricePerEpochBounds(dealP.PieceSize(), dealP.StartEpoch(), dealP.EndEpoch())
	if dealP.StoragePricePerEpoch() < minPrice || dealP.StoragePricePerEpoch() > maxPrice {
		rt.AbortStateMsg("Storage price out of bounds.")
	}

	minProviderCollateral, maxProviderCollateral := inds.StorageDeal_ProviderCollateralBounds(
		dealP.PieceSize(), dealP.StartEpoch(), dealP.EndEpoch())
	if dealP.ProviderCollateral() < minProviderCollateral || dealP.ProviderCollateral() > maxProviderCollateral {
		rt.AbortStateMsg("Provider collateral out of bounds.")
	}

	minClientCollateral, maxClientCollateral := inds.StorageDeal_ClientCollateralBounds(
		dealP.PieceSize(), dealP.StartEpoch(), dealP.EndEpoch())
	if dealP.ClientCollateral() < minClientCollateral || dealP.ClientCollateral() > maxClientCollateral {
		rt.AbortStateMsg("Client collateral out of bounds.")
	}
}

func (st *StorageMarketActorState_I) _rtGetOnChainDealOrAbort(rt Runtime, dealID deal.DealID) (deal deal.OnChainDeal, dealP deal.StorageDealProposal) {
	var found bool
	deal, dealP, found = st._getOnChainDeal(dealID)
	if !found {
		rt.AbortStateMsg("dealID not found in Deals.")
	}
	return
}

func (st *StorageMarketActorState_I) _rtLockBalanceOrAbort(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.AbortArgMsg("Negative amount")
	}

	st._rtAbortIfAddressEntryDoesNotExist(rt, addr)

	ok := st._lockBalanceMaybe(addr, amount)

	if !ok {
		rt.AbortFundsMsg("Insufficient funds available to lock.")
	}
}

////////////////////////////////////////////////////////////////////////////////
// Dispatch table
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMarketActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	IMPL_FINISH()
	panic("")
}
