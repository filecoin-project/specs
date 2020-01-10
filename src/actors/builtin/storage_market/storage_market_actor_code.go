package storage_market

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	actor_util "github.com/filecoin-project/specs/actors/util"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
)

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, entryAddr addr.Address, amountRequested abi.TokenAmount) {
	IMPL_FINISH() // BigInt arithmetic
	amountSlashedTotal := abi.TokenAmount(0)

	if amountRequested < 0 {
		rt.AbortArgMsg("Negative amount.")
	}

	recipientAddr := RT_MinerEntry_ValidateCaller_DetermineFundsLocation(rt, entryAddr, vmr.MinerEntrySpec_MinerOrSignable)

	h, st := a.State(rt)
	st._rtAbortIfAddressEntryDoesNotExist(rt, entryAddr)

	// Before any operations that check the balance tables for funds, execute all deferred
	// deal state updates.
	//
	// Note: as an optimization, implementations may cache efficient data structures indicating
	// which of the following set of updates are redundant and can be skipped.
	amountSlashedTotal += st._rtUpdatePendingDealStatesForParty(rt, entryAddr)

	minBalance := st._getLockedReqBalanceInternal(entryAddr)
	newTable, amountExtracted, ok := actor_util.BalanceTable_WithExtractPartial(
		st.EscrowTable(), entryAddr, amountRequested, minBalance)
	Assert(ok)
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
	rt.SendFunds(recipientAddr, amountExtracted)
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime, entryAddr addr.Address) {
	RT_MinerEntry_ValidateCaller_DetermineFundsLocation(rt, entryAddr, vmr.MinerEntrySpec_MinerOrSignable)

	h, st := a.State(rt)
	st._rtAbortIfAddressEntryDoesNotExist(rt, entryAddr)

	msgValue := rt.ValueReceived()
	newTable, ok := actor_util.BalanceTable_WithAdd(st.EscrowTable(), entryAddr, msgValue)
	if !ok {
		// Entry not found; create implicitly.
		newTable, ok = actor_util.BalanceTable_WithNewAddressEntry(st.EscrowTable(), entryAddr, msgValue)
		Assert(ok)
	}
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) PublishStorageDeals(rt Runtime, newStorageDeals []deal.StorageDeal) {
	IMPL_FINISH() // BigInt arithmetic
	amountSlashedTotal := abi.TokenAmount(0)

	// Deal message must have a From field identical to the provider of all the deals.
	// This allows us to retain and verify only the client's signature in each deal proposal itself.
	RT_ValidateImmediateCallerIsSignable(rt)
	providerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	// All storage deals will be added in an atomic transaction; this operation will be unrolled if any of them fails.
	for _, newDeal := range newStorageDeals {
		p := newDeal.Proposal()

		if !p.Provider().Equals(providerAddr) {
			rt.AbortArgMsg("Incorrect provider listed in deal")
		}

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
			SectorStartEpoch_: epochUndefined,
		}

		st.Deals()[id] = onchainDeal

		if _, found := st.CachedExpirationsPending()[p.EndEpoch()]; !found {
			st.CachedExpirationsPending()[p.EndEpoch()] = actor_util.DealIDQueue_Empty()
		}
		st.CachedExpirationsPending()[p.EndEpoch()].Enqueue(id)
	}

	st.Impl().CurrEpochNumDealsPublished_ += len(newStorageDeals)

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
}

// Note: in the case of a capacity-commitment sector (one with zero deals), this function should succeed vacuously.
func (a *StorageMarketActorCode_I) VerifyDealsOnSectorPreCommit(rt Runtime, dealIDs deal.DealIDs, sectorInfo sector.SectorPreCommitInfo) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		deal, _ := st._rtGetOnChainDealOrAbort(rt, dealID)
		_rtAbortIfDealInvalidForNewSectorSeal(rt, minerAddr, sectorInfo.Expiration(), deal)
	}

	Release(rt, h, st)
}

func (a *StorageMarketActorCode_I) UpdateDealsOnSectorProveCommit(rt Runtime, dealIDs deal.DealIDs, sectorInfo sector.SectorProveCommitInfo) {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
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
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)

	ret := []sector.PieceInfo{}

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._rtGetOnChainDealOrAbort(rt, dealID)
		ret = append(ret, sector.PieceInfo_I{
			PieceCID_: dealP.PieceCID(),
			Size_:     uint64(dealP.PieceSize().Total()),
		}.Ref())
	}

	Release(rt, h, st)

	return &sector.PieceInfos_I{Items_: ret}
}

func (a *StorageMarketActorCode_I) GetWeightForDealSet(rt Runtime, dealIDs deal.DealIDs) deal.DealWeight {
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
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
	rt.ValidateImmediateCallerAcceptAnyOfType(builtin.StorageMinerActorCodeID)
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
	rt.ValidateImmediateCallerIs(addr.CronActorAddr)

	h, st := a.State(rt)

	// Some deals may never be affected by the normal calls to _rtUpdatePendingDealStatesForParty
	// (notably, if the relevant party never checks its balance).
	// Without some cleanup mechanism, these deals may gradually accumulate and cause
	// the StorageMarketActor state to grow without bound.
	// To prevent this, we amortize the cost of this cleanup by processing a relatively
	// small number of deals every epoch, independent of the calls above.
	//
	// More specifically, we process deals:
	//   (a) In priority order of expiration epoch, up until the current epoch
	//   (b) Within a given expiration epoch, in order of original publishing.
	//
	// We stop once we have exhausted this valid set, or when we have hit a certain target
	// (DEAL_PROC_AMORTIZED_SCALE_FACTOR times the number of deals freshly published in the
	// current epoch) of deals dequeued, whichever comes first.

	const DEAL_PROC_AMORTIZED_SCALE_FACTOR = 2
	numDequeuedTarget := st.CurrEpochNumDealsPublished() * DEAL_PROC_AMORTIZED_SCALE_FACTOR

	numDequeued := 0
	extractedDealIDs := []deal.DealID{}

	for {
		if st.CachedExpirationsNextProcEpoch() > rt.CurrEpoch() {
			break
		}

		if numDequeued >= numDequeuedTarget {
			break
		}

		queue, found := st.CachedExpirationsPending()[st.CachedExpirationsNextProcEpoch()]
		if !found {
			st.Impl().CachedExpirationsNextProcEpoch_ += 1
			continue
		}

		queueDepleted := false
		for {
			dealID, ok := queue.Dequeue()
			if !ok {
				queueDepleted = true
				break
			}
			numDequeued += 1
			if _, found := st.Deals()[dealID]; found {
				// May have already processed expiration, independently, via _rtUpdatePendingDealStatesForParty.
				// If not, add it to the list to be processed.
				extractedDealIDs = append(extractedDealIDs, dealID)
			}
		}

		if !queueDepleted {
			Assert(numDequeued >= numDequeuedTarget)
			break
		}

		delete(st.CachedExpirationsPending(), st.CachedExpirationsNextProcEpoch())
		st.Impl().CachedExpirationsNextProcEpoch_ += 1
	}

	amountSlashedTotal := st._updatePendingDealStates(extractedDealIDs, rt.CurrEpoch())

	// Reset for next epoch.
	st.Impl().CurrEpochNumDealsPublished_ = 0

	UpdateRelease(rt, h, st)

	rt.SendFunds(addr.BurntFundsActorAddr, amountSlashedTotal)
}

func (a *StorageMarketActorCode_I) Constructor(rt Runtime) {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	h := rt.AcquireState()

	st := &StorageMarketActorState_I{
		Deals_:                          DealsAMT_Empty(),
		EscrowTable_:                    actor_util.BalanceTableHAMT_Empty(),
		LockedReqTable_:                 actor_util.BalanceTableHAMT_Empty(),
		NextID_:                         deal.DealID(0),
		CachedDealIDsByParty_:           CachedDealIDsByPartyHAMT_Empty(),
		CachedExpirationsPending_:       CachedExpirationsPendingHAMT_Empty(),
		CachedExpirationsNextProcEpoch_: abi.ChainEpoch(0),
	}

	UpdateRelease(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _rtUpdatePendingDealStatesForParty(rt Runtime, addr addr.Address) (
	amountSlashedTotal abi.TokenAmount) {

	// For consistency with OnEpochTickEnd, only process updates up to the end of the _previous_ epoch.
	epoch := rt.CurrEpoch() - 1

	cachedRes, ok := st.CachedDealIDsByParty()[addr]
	Assert(ok)
	extractedDealIDs := []deal.DealID{}
	for cachedDealID := range cachedRes {
		extractedDealIDs = append(extractedDealIDs, cachedDealID)
	}

	amountSlashedTotal = st._updatePendingDealStates(extractedDealIDs, epoch)
	return
}

func _rtAbortIfDealAlreadyProven(rt Runtime, deal deal.OnChainDeal) {
	if deal.SectorStartEpoch() != epochUndefined {
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

func _rtAbortIfDealExceedsSectorLifetime(rt Runtime, dealP deal.StorageDealProposal, sectorExpiration abi.ChainEpoch) {
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
	rt Runtime, minerAddr addr.Address, sectorExpiration abi.ChainEpoch, deal deal.OnChainDeal) {

	dealP := deal.Deal().Proposal()

	_rtAbortIfDealNotFromProvider(rt, dealP, minerAddr)
	_rtAbortIfDealAlreadyProven(rt, deal)
	_rtAbortIfDealStartElapsed(rt, dealP)
	_rtAbortIfDealExceedsSectorLifetime(rt, dealP, sectorExpiration)
}

func _rtAbortIfNewDealInvalid(rt Runtime, deal deal.StorageDeal) {
	dealP := deal.Proposal()

	if !_rtDealProposalIsInternallyValid(rt, dealP) {
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

func (st *StorageMarketActorState_I) _rtLockBalanceOrAbort(rt Runtime, addr addr.Address, amount abi.TokenAmount) {
	if amount < 0 {
		rt.AbortArgMsg("Negative amount")
	}

	st._rtAbortIfAddressEntryDoesNotExist(rt, addr)

	ok := st._lockBalanceMaybe(addr, amount)

	if !ok {
		rt.AbortFundsMsg("Insufficient funds available to lock.")
	}
}
