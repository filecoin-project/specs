package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

const (
	MethodGetUnsealedCIDForDealIDs = actor.MethodPlaceholder
	MethodProcessDealExpiration    = actor.MethodPlaceholder
	MethodProcessDealSlash         = actor.MethodPlaceholder
)

const LastPaymentEpochNone = 0

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type State = StorageMarketActorState

var TODO = util.TODO

func (a *StorageMarketActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.Abort("IPLD lookup error")
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

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, amount actor.TokenAmount) {

	msgSender := rt.ImmediateCaller()
	panic("TODO: assert caller is miner worker")

	h, st := a.State(rt)

	if amount <= 0 {
		rt.Abort("non-positive balance to withdraw.")
	}

	senderBalance, found := st.Balances()[msgSender]
	if !found {
		rt.Abort("sender address not found.")
	}

	if senderBalance.Available() < amount {
		rt.Abort("insufficient balance.")
	}

	senderBalance.Impl().Available_ = senderBalance.Available() - amount
	st.Balances()[msgSender] = senderBalance

	UpdateRelease(rt, h, st)

	// send funds to miner
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    msgSender,
		Value_: amount,
	})
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime) {

	msgSender := rt.ImmediateCaller()
	msgValue := rt.ValueReceived()

	h, st := a.State(rt)

	senderBalance, found := st.Balances()[msgSender]
	if found {
		senderBalance.Impl().Available_ = senderBalance.Available() + msgValue
		st.Balances()[msgSender] = senderBalance
	} else {
		st.Balances()[msgSender] = &StorageParticipantBalance_I{
			Locked_:    0,
			Available_: msgValue,
		}
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) PublishStorageDeals(rt Runtime, newStorageDeals []deal.StorageDeal) []PublishStorageDealResponse {

	TODO() // verify StorageMinerActor

	h, st := a.State(rt)

	l := len(newStorageDeals)
	response := make([]PublishStorageDealResponse, l)

	// all storage deals will be added in an atomic transaction
	// _validateNewStorageDeal will throw if
	//     - deal started/expired before it is signed
	//     - deal hits the chain after StartEpoch
	//     - incorrect client and provider addresses
	//     - insufficient balance lock up
	// this operation will be unrolled if any of the above triggers a throw
	for i, newDeal := range newStorageDeals {
		if st._validateNewStorageDeal(rt, newDeal) {
			st._lockFundsForStorageDeal(rt, newDeal)
			id := st._generateStorageDealID(rt, newDeal)

			onchainDeal := &deal.OnChainDeal_I{
				ID_:               id,
				Deal_:             newDeal,
				LastPaymentEpoch_: block.ChainEpoch(LastPaymentEpochNone), // -1 = inactive
			}
			st.Deals()[id] = onchainDeal
			response[i] = PublishStorageDealSuccess(id)
		} else {
			response[i] = PublishStorageDealError()
		}
	}

	UpdateRelease(rt, h, st)

	return response
}

func (a *StorageMarketActorCode_I) VerifyPublishedDealIDs(rt Runtime, dealIDs []deal.DealID) {

	TODO() // verify StorageMinerActor

	h, st := a.State(rt)

	for _, dealID := range dealIDs {

		publishedDeal := st._getOnChainDeal(rt, dealID)
		st._assertPublishedDealState(rt, dealID)

		dealP := publishedDeal.Deal().Proposal()
		st._assertDealStartAfterCurrEpoch(rt, dealP)

		// deal must not expire before the maximum allowable epoch between pre and prove commits
		// we do not have to check if the deal has expired at ProveCommit
		// if the MAX_PROVE_COMMIT_SECTOR_EPOCH constraint is not violated
		st._assertDealExpireAfterMaxProveCommitWindow(rt, dealP)

	}

	Release(rt, h, st)
}

func (a *StorageMarketActorCode_I) ActivateDeals(rt Runtime, dealIDs []deal.DealID) []deal.OnChainDeal {

	TODO() // verify StorageMinerActor

	h, st := a.State(rt)
	ret := make([]deal.OnChainDeal, len(dealIDs))

	for _, dealID := range dealIDs {
		publishedDeal := st._getOnChainDeal(rt, dealID)
		st._assertPublishedDealState(rt, dealID)

		dealP := publishedDeal.Deal().Proposal()
		st._assertDealNotYetExpired(rt, dealP)

		onchainDeal := st._activateDeal(rt, publishedDeal)
		ret = append(ret, onchainDeal)
	}

	UpdateRelease(rt, h, st)

	return ret

}

func (a *StorageMarketActorCode_I) ProcessDealSlash(rt Runtime, dealIDs []deal.DealID, faultType sector.StorageFaultType) {

	TODO() // only call by StorageMinerActor

	h, st := a.State(rt)

	// only terminated fault will result in slashing of deal collateral
	amountSlashed := actor.TokenAmount(0)
	switch faultType {
	case sector.TerminatedFault:
		for _, dealID := range dealIDs {
			amountSlashed += st._terminateDeal(rt, dealID)
		}
	default:
		// do nothing
	}

	UpdateRelease(rt, h, st)

	// send funds to BurntFundsActor
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    addr.BurntFundsActorAddr,
		Value_: amountSlashed,
	})

}

func (a *StorageMarketActorCode_I) ProcessDealPayment(rt Runtime, dealIDs []deal.DealID, newPaymentEpoch block.ChainEpoch) {
	h, st := a.State(rt)

	for _, dealID := range dealIDs {
		deal := st._getOnChainDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)

		fee := st._getStorageFeeSinceLastPayment(rt, deal, newPaymentEpoch)

		dealP := deal.Deal().Proposal()
		st._transferBalance(rt, dealP.Client(), dealP.Provider(), fee)

		// update LastPaymentEpoch in deal
		deal.Impl().LastPaymentEpoch_ = newPaymentEpoch
		st.Deals()[dealID] = deal
	}

	UpdateRelease(rt, h, st)
}

// unlock remaining payments and return all UnlockedStorageFee to Provider
// remove deals from ActiveDeals
// return collaterals to both miner and client
func (a *StorageMarketActorCode_I) ProcessDealExpiration(rt Runtime, dealIDs []deal.DealID) {
	h, st := a.State(rt)

	for _, dealID := range dealIDs {

		expiredDeal := st._getOnChainDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)

		dealP := expiredDeal.Deal().Proposal()
		fee := st._getStorageFeeSinceLastPayment(rt, expiredDeal, dealP.EndEpoch())

		delete(st.Deals(), dealID)
		st._transferBalance(rt, dealP.Client(), dealP.Provider(), fee)

		// return storage deal collaterals to both miners and client
		st._unlockBalance(rt, dealP.Provider(), dealP.ProviderBalanceRequirement())
		st._unlockBalance(rt, dealP.Client(), dealP.TotalClientCollateral())

	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) GetPieceInfosForDealIDs(rt Runtime, sectorSize util.UVarint, dealIDs []deal.DealID) []sector.PieceInfo_I {
	pieceInfos := make([]sector.PieceInfo_I, len(dealIDs))

	h, st := a.State(rt)

	for index, deal := range st.Deals() {
		proposal := deal.Deal().Proposal()
		pieceSize := util.UInt(proposal.PieceSize().Total())

		var pieceInfo sector.PieceInfo_I
		pieceInfo.PieceCID_ = proposal.PieceCID()
		pieceInfo.Size_ = pieceSize

		pieceInfos[index] = pieceInfo
	}

	Release(rt, h, st)

	return pieceInfos
}

func (a *StorageMarketActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
