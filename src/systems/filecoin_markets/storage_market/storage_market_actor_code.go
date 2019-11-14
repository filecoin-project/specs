package storage_market

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import util "github.com/filecoin-project/specs/util"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

const (
	MethodGetUnsealedCIDForDealIDs = actor.MethodNum(3)
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = msg.InvocOutput
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

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, balance actor.TokenAmount) {
	h, st := a.State(rt)

	var msgSender addr.Address // TODO replace this from VM runtime

	if balance < 0 {
		rt.Abort("negative balance to withdraw.")
	}

	senderBalance, found := st.Balances()[msgSender]
	if !found {
		rt.Abort("sender address not found.")
	}

	if senderBalance.Available() < balance {
		rt.Abort("insufficient balance.")
	}

	senderBalance.Impl().Available_ = senderBalance.Available() - balance
	st.Balances()[msgSender] = senderBalance

	// TODO send funds to msgSender with `transferBalance` in VM runtime

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime) {
	h, st := a.State(rt)

	var msgSender addr.Address    // TODO replace this
	var balance actor.TokenAmount // TODO replace this

	// TODO subtract balance from msgSender
	// TODO add balance to StorageMarketActor
	if balance < 0 {
		rt.Abort("negative balance to add.")
	}

	senderBalance, found := st.Balances()[msgSender]
	if found {
		senderBalance.Impl().Available_ = senderBalance.Available() + balance
		st.Balances()[msgSender] = senderBalance
	} else {
		st.Balances()[msgSender] = &StorageParticipantBalance_I{
			Locked_:    0,
			Available_: balance,
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
			st.Deals()[id] = newDeal
			st.Impl().DealStates_.Publish(id)
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

		publishedDeal := st._getDeal(rt, dealID)
		st._assertPublishedDealState(rt, dealID)

		dealP := publishedDeal.Proposal()
		st._assertDealStartAfterCurrEpoch(rt, dealP)

		// deal must not expire before the maximum allowable epoch between pre and prove commits
		// we do not have to check if the deal has expired at ProveCommit
		// if the MAX_PROVE_COMMIT_SECTOR_EPOCH constraint is not violated
		st._assertDealExpireAfterMaxProveCommitWindow(rt, dealP)

	}

	Release(rt, h, st)
}

func (a *StorageMarketActorCode_I) ActivateDeals(rt Runtime, dealIDs []deal.DealID) {

	TODO() // verify StorageMinerActor

	h, st := a.State(rt)

	for _, dealID := range dealIDs {
		publishedDeal := st._getDeal(rt, dealID)
		st._assertPublishedDealState(rt, dealID)

		dealP := publishedDeal.Proposal()
		st._assertDealNotYetExpired(rt, dealP)

		// should only go through if all deals satisfy the above invariant
		dealTally := &deal.StorageDealTally_I{
			ProviderCollateralRemaining_: publishedDeal.Proposal().ProviderBalanceRequirement(),
			LockedStorageFee_:            publishedDeal.Proposal().TotalStorageFee(),
			UnlockedStorageFee_:          actor.TokenAmount(0),
			LastPaymentEpoch_:            rt.CurrEpoch(),
		}

		st.DealTally()[dealID] = dealTally
		st.Impl().DealStates_.Activate(dealID)
	}

	UpdateRelease(rt, h, st)

}

func (a *StorageMarketActorCode_I) GetDeals(rt Runtime, dealIDs []deal.DealID) []deal.StorageDeal {

	TODO() // verify StorageMinerActor

	h, st := a.State(rt)

	ret := make([]deal.StorageDeal, len(dealIDs))

	for _, dealID := range dealIDs {

		activeDeal := st._getDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)

		ret = append(ret, activeDeal)
	}

	Release(rt, h, st)

	return ret

}

func (a *StorageMarketActorCode_I) ProcessDealSlash(rt Runtime, info deal.BatchDealSlashInfo) {

	TODO() // TODO: only call by StorageMinerActor

	h, st := a.State(rt)

	// only terminated fault will result in slashing of deal collateral
	switch info.Action() {
	case deal.SlashTerminatedFaults:
		st._slashTerminatedFaults(rt, info.DealIDs())
	default:
		rt.Abort("sma.ProcessDealSlash: invalid action type")
	}

	UpdateRelease(rt, h, st)

}

func (a *StorageMarketActorCode_I) CreditUnlockedFees(rt Runtime, dealIDs []deal.DealID) {
	TODO() // verify StorageMinerActor

	h, st := a.State(rt)

	for _, dealID := range dealIDs {
		activeDeal := st._getDeal(rt, dealID)
		st._assertActiveDealState(rt, dealID)
		dealTally := st._getDealTally(rt, dealID)
		dealP := activeDeal.Proposal()

		st._creditUnlockedFeeForProvider(rt, dealP, dealTally)
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) ProcessDealPayment(rt Runtime, info deal.BatchDealPaymentInfo) {
	h, st := a.State(rt)

	switch info.Action() {
	case deal.ExpireStorageDeals:
		st._expireStorageDeals(rt, info.DealIDs(), info.LastChallengeEndEpoch())
	case deal.TallyStorageDeals:
		st._tallyStorageDeals(rt, info.DealIDs(), info.LastChallengeEndEpoch())
	default:
		rt.Abort("sma.ProcessDealPayment: invalid deal payment action.")
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) GetPieceInfosForDealIDs(rt Runtime, sectorSize util.UVarint, dealIDs []deal.DealID) []sector.PieceInfo_I {
	pieceInfos := make([]sector.PieceInfo_I, len(dealIDs))

	h, st := a.State(rt)

	for index, deal := range st.Deals() {
		proposal := deal.Proposal()
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
