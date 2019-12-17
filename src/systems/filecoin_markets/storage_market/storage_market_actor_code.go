package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type State = StorageMarketActorState

var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

func (a *StorageMarketActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state := DeserializeState(stateBytes.As_Bytes())
	return h, state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st State) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st State) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StorageMarketActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) State {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func _rtGetMinerAccountsAssert(rt Runtime, minerAddr addr.Address) (ownerAddr addr.Address, workerAddr addr.Address) {
	ownerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetOwnerAddr, []util.Serialization{}))

	workerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetWorkerAddr, []util.Serialization{}))

	return
}

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, minerAddr addr.Address, amountRequested actor.TokenAmount) {
	TODO() // must also handle case of non-miner (client) sender

	ownerAddr, workerAddr := _rtGetMinerAccountsAssert(rt, minerAddr)
	rt.ValidateImmediateCallerInSet([]addr.Address{ownerAddr, workerAddr})
	// TODO: should workerAddr be permitted here? Depends on actor principal security assumptions.

	callerAddr := rt.ImmediateCaller()

	if amountRequested < 0 {
		rt.AbortArgMsg("sma.WithdrawBalance: negative amount.")
	}

	h, st := a.State(rt)

	if !st._addressEntryExists(minerAddr) {
		rt.AbortArgMsg("sma.WithdrawBalance: address entry does not exist")
	}

	// Before any operations that check the balance tables for funds, execute all deferred
	// deal state updates.
	//
	// Note: as an optimization, implementations may cache efficient data structures indicating
	// which of the following set of updates are redundant and can be skipped.
	st._rtUpdatePendingDealStatesForParty(rt, minerAddr)

	minBalance := st._getLockedReqBalanceInternal(minerAddr)
	newTable, amountExtracted, ok := actor.BalanceTable_WithExtractPartial(
		st.EscrowTable(), minerAddr, amountRequested, minBalance)
	Assert(ok)
	st.Impl().EscrowTable_ = newTable

	UpdateRelease(rt, h, st)

	// send funds to miner
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    callerAddr,
		Value_: amountExtracted,
	})
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime, minerAddr addr.Address) {
	TODO() // must also handle case of non-miner (client) sender

	ownerAddr, workerAddr := _rtGetMinerAccountsAssert(rt, minerAddr)
	rt.ValidateImmediateCallerInSet([]addr.Address{ownerAddr, workerAddr})
	// TODO: should workerAddr be permitted here? Depends on actor principal security assumptions.

	msgValue := rt.ValueReceived()

	h, st := a.State(rt)

	if !st._addressEntryExists(minerAddr) {
		rt.AbortArgMsg("sma.WithdrawBalance: address entry does not exist")
	}

	newTable, ok := actor.BalanceTable_WithAdd(st.EscrowTable(), minerAddr, msgValue)
	if !ok {
		// Entry not found; create implicitly.
		newTable, ok = actor.BalanceTable_WithNewAddressEntry(st.EscrowTable(), minerAddr, msgValue)
		Assert(ok)
	}
	st.Impl().EscrowTable_ = newTable
	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) PublishStorageDeals(rt Runtime, newStorageDeals []deal.StorageDeal) {
	// Caller should be one of the parties involved in the deals (every deal?).
	TODO() // TODO: validate this
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_Account)

	h, st := a.State(rt)

	// all storage deals will be added in an atomic transaction
	// _validateNewStorageDeal will throw if
	//     - deal started/expired before it is signed
	//     - deal hits the chain after StartEpoch
	//     - incorrect client and provider addresses
	//     - insufficient balance lock up
	// this operation will be unrolled if any of the above triggers a throw
	for _, newDeal := range newStorageDeals {
		p := newDeal.Proposal()

		st._rtAbortIfNewDealInvalid(rt, newDeal)

		// Before any operations that check the balance tables for funds, execute all deferred
		// deal state updates.
		//
		// Note: as an optimization, implementations may cache efficient data structures indicating
		// which of the following set of updates are redundant and can be skipped.
		st._rtUpdatePendingDealStatesForParty(rt, p.Client())
		st._rtUpdatePendingDealStatesForParty(rt, p.Provider())

		st._rtLockBalanceUntrusted(rt, p.Client(), p.ClientBalanceRequirement())
		st._rtLockBalanceUntrusted(rt, p.Provider(), p.ProviderBalanceRequirement())

		id := st._generateStorageDealID(rt, newDeal)

		onchainDeal := &deal.OnChainDeal_I{
			ID_:               id,
			Deal_:             newDeal,
			SectorStartEpoch_: block.ChainEpoch_None,
		}
		st.Deals()[id] = onchainDeal
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) VerifyDealsOnSectorPreCommit(rt Runtime, dealIDs deal.DealIDs, sectorInfo sector.SectorPreCommitInfo) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		deal, _ := st._rtGetOnChainDealOrAbort(rt, dealID)
		_rtAbortIfDealInvalidForNewSectorSeal(rt, minerAddr, sectorInfo.Expiration(), deal)

		// deal must not expire before the maximum allowable epoch between pre and prove commits
		// we do not have to check if the deal has expired at ProveCommit
		// if the MAX_PROVE_COMMIT_SECTOR_EPOCH constraint is not violated
		// if dealP.EndEpoch() <= (rt.CurrEpoch() + sector.MAX_PROVE_COMMIT_SECTOR_EPOCH) {
		// 	rt.AbortStateMsg("Deal might expire before prove commit.")
		// }
		TODO() // TODO: confirm that the above is no longer required given the new deal expiration logic
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

	// TODO potentially refund clients for started deals
	//  ^-- TODO: omit this (would require pro-rating deal payment)?
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

func (a *StorageMarketActorCode_I) GetWeightForDealSet(rt Runtime, dealIDs deal.DealIDs) int {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	ret := 0 // TODO: BigInt arithmetic

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)
		_rtAbortIfDealNotFromProvider(rt, dealP, minerAddr) // TODO: Assert (not Abort)

		ret += int(dealP.Duration()) * int(dealP.PieceSize().Total()) // TODO: BigInt arithmetic
	}

	UpdateRelease(rt, h, st)

	return ret
}

func (a *StorageMarketActorCode_I) TerminateDealsOnSlashProviderSector(rt Runtime, dealIDs deal.DealIDs) {
	rt.ValidateImmediateCallerAcceptAnyOfType(actor.BuiltinActorID_StorageMiner)
	minerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	for _, dealID := range dealIDs.Items() {
		_, dealP := st._rtGetOnChainDealOrAbort(rt, dealID)
		Assert(dealP.Provider().Equals(minerAddr))
		// Note: we do not perform the balance transfers here, but rather simply record the flag
		// to indicate that _rtProcessDealSlashed should be called when the deferred state computation
		// is performed.
		st.Deals()[dealID].Impl().SlashEpoch_ = rt.CurrEpoch()
	}

	UpdateRelease(rt, h, st)
}

func (st *StorageMarketActorState_I) _rtUpdatePendingDealStatesForParty(rt Runtime, addr addr.Address) {
	// For consistency with OnEpochTickEnd, only process updates up to the end of the _previous_ epoch.
	epoch := rt.CurrEpoch() - 1

	var relevantDealIDs []deal.DealID
	// Populate with the set of all elements in st.Deals() in which addr is one of the parties
	// (either client or provider).
	IMPL_FINISH()

	st._rtUpdatePendingDealStates(rt, relevantDealIDs, epoch)
}

func (st *StorageMarketActorState_I) _rtUpdatePendingDealStates(rt Runtime, dealIDs []deal.DealID, epoch block.ChainEpoch) {
	for _, dealID := range dealIDs {
		st._rtUpdatePendingDealState(rt, dealID, epoch)
	}
}

func (st *StorageMarketActorState_I) _rtUpdatePendingDealState(rt Runtime, dealID deal.DealID, epoch block.ChainEpoch) {
	deal, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)

	everUpdated := (deal.LastUpdatedEpoch() != block.ChainEpoch_None)
	everSlashed := (deal.SlashEpoch() != block.ChainEpoch_None)

	Assert(!everUpdated || (deal.LastUpdatedEpoch() <= epoch))
	if deal.LastUpdatedEpoch() == epoch {
		return
	}

	if deal.SectorStartEpoch() == block.ChainEpoch_None {
		// Not yet appeared in proven sector; check for timeout.
		if dealP.StartEpoch() >= epoch {
			st._rtProcessDealInitTimedOut(rt, dealID)
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
	st._rtProcessDealEpochsElapsed(rt, dealID, numEpochsElapsed)

	if everSlashed {
		st._rtProcessDealSlashed(rt, dealID)
		return
	}

	if epoch >= dealP.EndEpoch() {
		st._rtProcessDealExpired(rt, dealID)
		return
	}

	st.Deals()[dealID].Impl().LastUpdatedEpoch_ = epoch
}

func (st *StorageMarketActorState_I) _rtProcessDealEpochsElapsed(rt Runtime, dealID deal.DealID, numEpochsElapsed block.ChainEpoch) {
	deal, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	// Process deal payment for the current epoch.
	totalPayment := int(numEpochsElapsed) * int(dealP.StoragePricePerEpoch()) // TODO: BigInt arithmetic
	st._rtTransferBalance(rt, dealP.Client(), dealP.Provider(), actor.TokenAmount(totalPayment))
}

func (st *StorageMarketActorState_I) _rtProcessDealSlashed(rt Runtime, dealID deal.DealID) {
	deal, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	slashEpoch := deal.SlashEpoch()
	Assert(slashEpoch != block.ChainEpoch_None)

	// unlock client collateral and locked storage fee
	clientCollateral := dealP.ClientCollateral()
	paymentRemaining := _dealGetPaymentRemaining(deal, rt.CurrEpoch())
	st._rtUnlockBalance(rt, dealP.Client(), clientCollateral+paymentRemaining)

	// slash provider collateral
	slashAmount := dealP.ProviderCollateral()
	st._rtSlashBalance(rt, dealP.Provider(), slashAmount)

	TODO() // TODO: tabulate for later send (can't send while holding state lock)
	rt.SendFunds(addr.BurntFundsActorAddr, slashAmount)

	delete(st.Deals(), dealID)
}

// Deal start deadline elapsed without appearing in a proven sector.
// (TODO: slash portion of provider's collateral TBD)
// Delete deal and unlock collaterals for both provider and client.
func (st *StorageMarketActorState_I) _rtProcessDealInitTimedOut(rt Runtime, dealID deal.DealID) {
	deal, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)
	Assert(deal.SectorStartEpoch() == block.ChainEpoch_None)

	st._rtUnlockBalance(rt, dealP.Client(), dealP.ClientBalanceRequirement())

	TODO() // slash portion of provider's collateral TBD
	st._rtUnlockBalance(rt, dealP.Provider(), dealP.ProviderBalanceRequirement())

	delete(st.Deals(), dealID)
}

// Normal expiration. Delete deal and unlock collaterals for both miner and client.
func (st *StorageMarketActorState_I) _rtProcessDealExpired(rt Runtime, dealID deal.DealID) {
	deal, dealP := st._rtGetOnChainDealOrAbort(rt, dealID) // TODO: Assert (not Abort)
	Assert(deal.SectorStartEpoch() != block.ChainEpoch_None)

	st._rtUnlockBalance(rt, dealP.Provider(), dealP.ProviderCollateral())
	st._rtUnlockBalance(rt, dealP.Client(), dealP.ClientCollateral())

	delete(st.Deals(), dealID)
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
	IMPL_FINISH()

	st._rtUpdatePendingDealStates(rt, cleanupDealIDs, rt.CurrEpoch())

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	IMPL_FINISH()
	panic("")
}
