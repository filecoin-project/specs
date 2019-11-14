package storage_market

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
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

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	// TODO
	var dealID deal.DealID
	return dealID
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
	// TODO: verify deal did not expire when it is signed

	currEpoch := rt.CurrEpoch()
	p := d.Proposal()

	// deal ends before it starts
	if p.EndEpoch() <= p.StartEpoch() {
		rt.Abort("sma._validateNewStorageDeal: deal ends before it starts.")
	}

	// deal has started before publish
	if p.StartEpoch() < currEpoch {
		rt.Abort("sma._validateNewStorageDeal: deal started before publish.")
	}

	// minimum deal duration
	if p.Duration() < deal.MIN_DEAL_DURATION {
		rt.Abort("sma._validateNewStorageDeal: deal duration shorter than minimum.")
	}

	// TODO: verify client and provider address and signature are correct (may not be needed)

	if p.StoragePricePerEpoch() <= deal.MIN_DEAL_PRICE {
		rt.Abort("sma._validateNewStorageDeal: storage price less than minimum.")
	}

	// verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
	if p.ProviderCollateralPerEpoch() < deal.MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH ||
		p.ClientCollateralPerEpoch() < deal.MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH {
		rt.Abort("sma._validateNewStorageDeal: deal collaterals do not meet minimum requirement.")
	}

	// verify client and provider has sufficient balance
	isClientBalAvailable := st._isBalanceAvailable(p.Client(), p.ClientBalanceRequirement())
	isProviderBalAvailable := st._isBalanceAvailable(p.Provider(), p.ProviderBalanceRequirement())

	if !isClientBalAvailable || !isProviderBalAvailable {
		rt.Abort("sma._validateNewStorageDeal: client or provider insufficient balance.")
	}

	return true
}

// TODO: consider returning a boolean
func (st *StorageMarketActorState_I) _lockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("addr not found.")
	}

	currBalance.Impl().Available_ -= amount
	currBalance.Impl().Locked_ += amount
}

func (st *StorageMarketActorState_I) _unlockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("addr not found.")
	}

	currBalance.Impl().Locked_ -= amount
	currBalance.Impl().Available_ += amount
}

// move funds from locked in client to available in provider
func (st *StorageMarketActorState_I) _transferBalance(rt Runtime, fromLocked addr.Address, toAvailable addr.Address, amount actor.TokenAmount) {
	fromB := st.Balances()[fromLocked]
	toB := st.Balances()[toAvailable]

	if fromB.Locked() < amount {
		rt.Abort("attempt to lock funds greater than actor has")
		return
	}

	fromB.Impl().Locked_ -= amount
	toB.Impl().Available_ += amount
}

func (st *StorageMarketActorState_I) _isBalanceAvailable(a addr.Address, amount actor.TokenAmount) bool {
	bal := st.Balances()[a]
	return bal.Available() >= amount
}

func (st *StorageMarketActorState_I) _lockFundsForStorageDeal(rt Runtime, deal deal.StorageDeal) {
	p := deal.Proposal()

	st._lockBalance(rt, p.Client(), p.ClientBalanceRequirement())
	st._lockBalance(rt, p.Provider(), p.ProviderBalanceRequirement())
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
			st.PublishedDeals()[id] = newDeal
			response[i] = PublishStorageDealSuccess(id)
		} else {
			response[i] = PublishStorageDealError()
		}
	}

	UpdateRelease(rt, h, st)

	return response
}

func (a *StorageMarketActorCode_I) VerifyPublishedDealIDs(rt Runtime, dealIDs []deal.DealID) bool {
	h, st := a.State(rt)

	ret := true

	for _, dealID := range dealIDs {
		publishedDeal, publishedFound := st.PublishedDeals()[dealID]
		if !publishedFound {
			ret = false
		}

		currEpoch := rt.CurrEpoch()

		if publishedDeal.Proposal().StartEpoch() < currEpoch {
			delete(st.PublishedDeals(), dealID)
			ret = false
		}

		dealExpiration := publishedDeal.Proposal().EndEpoch()

		// deal must not expire before the maximum allowable epoch between pre and prove commits
		// we do not have to check if the deal has expired at ProveCommit
		// if the MAX_PROVE_COMMIT_SECTOR_EPOCH constraint is not violated
		if dealExpiration <= (currEpoch + sector.MAX_PROVE_COMMIT_SECTOR_EPOCH) {
			delete(st.PublishedDeals(), dealID)
			ret = false
		}
	}

	UpdateRelease(rt, h, st)

	return ret
}

func (a *StorageMarketActorCode_I) ActivateDeals(rt Runtime, dealIDs []deal.DealID) bool {
	h, st := a.State(rt)

	ret := true

	// TODO: verify rt.Abort is sufficient and no need for return false
	// might need to return false but that may not be correct either, need to unroll changes
	for _, dealID := range dealIDs {
		publishedDeal, publishedFound := st.PublishedDeals()[dealID]
		if !publishedFound {
			delete(st.PublishedDeals(), dealID)
			rt.Abort("sma.ActivateSectorDealIDs: unpublished deal ID.")
		}

		if publishedDeal.Proposal().EndEpoch() <= rt.CurrEpoch() {
			delete(st.PublishedDeals(), dealID)
			rt.Abort("sma.ActivateSectorDealIDs: storage deal has expired.")
		}

		// should only go through if all deals satisfy the above invariant
		activeDeal := &deal.ActiveStorageDeal_I{
			Deal_:                        publishedDeal,
			ProviderCollateralRemaining_: publishedDeal.Proposal().ProviderBalanceRequirement(),
			LockedStorageFee_:            publishedDeal.Proposal().TotalStorageFee(),
			UnlockedStorageFee_:          actor.TokenAmount(0),
			LastPaymentEpoch_:            rt.CurrEpoch(),
		}

		delete(st.PublishedDeals(), dealID)
		st.ActiveDeals()[dealID] = activeDeal
	}

	UpdateRelease(rt, h, st)

	return ret
}

// Call by StorageMinerActor at CommitSector
func (a *StorageMarketActorCode_I) GetInitialUtilizationInfo(rt Runtime, dealIDs []deal.DealID) sector.SectorUtilizationInfo {
	h, st := a.State(rt)

	var dealExpirationQueue deal.DealExpirationQueue
	var maxUtilization block.StoragePower
	var lastExpiration block.ChainEpoch
	activeDealIDs := deal.CompactDealSet(make([]byte, len(dealIDs)))

	for _, dealID := range dealIDs {
		d, found := st.ActiveDeals()[dealID]
		if !found {
			rt.Abort("sm.GetInitialUtilizationInfo: dealID not found in ActiveDeals.")
		}

		// TODO: more checks or be convinced that it's enough to assume deals are still valid
		// consider calling _validateNewStorageDeal

		dealExpiration := d.Deal().Proposal().EndEpoch()

		if dealExpiration > lastExpiration {
			lastExpiration = dealExpiration
		}

		// TODO: verify what counts towards power here
		// There is PayloadSize, OverheadSize, and Total, see piece.id
		dealPayloadPower := block.StoragePower(d.Deal().Proposal().PieceSize().PayloadSize())

		queueItem := &deal.DealExpirationQueueItem_I{
			DealID_:       dealID,
			PayloadPower_: dealPayloadPower,
			Expiration_:   dealExpiration,
		}
		dealExpirationQueue.Add(queueItem)
		activeDealIDs.Add(dealID)
		maxUtilization += dealPayloadPower

	}

	initialUtilizationInfo := &sector.SectorUtilizationInfo_I{
		DealExpirationQueue_: dealExpirationQueue,
		MaxUtilization_:      maxUtilization,
		CurrUtilization_:     maxUtilization,
	}

	Release(rt, h, st)

	return initialUtilizationInfo
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

// unlock remaining payments and return all UnlockedStorageFee to Provider
// remove deals from ActiveDeals
// return collaterals to both miner and client
func (st *StorageMarketActorState_I) _expireStorageDeals(rt Runtime, dealIDs []deal.DealID, lastChallengeEndEpoch block.ChainEpoch) {
	for _, dealID := range dealIDs {

		expiredDeal, found := st.ActiveDeals()[dealID]

		if !found {
			rt.Abort("sma._expireStorageDeals: dealID not found")
		}

		duration := lastChallengeEndEpoch - expiredDeal.LastPaymentEpoch()
		if duration > 0 {
			fee := actor.TokenAmount(uint64(duration) * uint64(expiredDeal.Deal().Proposal().StoragePricePerEpoch()))

			if fee > expiredDeal.LockedStorageFee() {
				rt.Abort("sma._expireStorageDeals: cannot unlock more than already locked.")
			}

			expiredDeal.Impl().LastPaymentEpoch_ = lastChallengeEndEpoch
			expiredDeal.Impl().UnlockStorageFee(fee)
		}

		delete(st.ActiveDeals(), dealID)

		dealP := expiredDeal.Deal().Proposal()
		// should we check if LockedStorageFee is now 0?
		// credit UnlockedStorageFee to miner
		st._unlockBalance(rt, dealP.Provider(), expiredDeal.UnlockedStorageFee())

		// return storage deal collaterals to both miners and client
		st._unlockBalance(rt, dealP.Provider(), expiredDeal.ProviderCollateralRemaining())
		st._unlockBalance(rt, dealP.Client(), dealP.TotalClientCollateral())

	}
}

func dealStorageFeeForEpochs(dealP deal.StorageDealProposal, epochs block.ChainEpoch) actor.TokenAmount {
	return actor.TokenAmount(uint64(epochs) * uint64(dealP.StoragePricePerEpoch()))
}

func (st *StorageMarketActorState_I) _creditStorageDeals(rt Runtime, dealIDs []deal.DealID, lastChallengeEndEpoch block.ChainEpoch) {

	for _, dealID := range dealIDs {
		activeDeal, found := st.ActiveDeals()[dealID]

		if !found {
			rt.Abort("sma._creditStorageDeals: dealID not found.")
		}

		duration := lastChallengeEndEpoch - activeDeal.LastPaymentEpoch()

		if duration <= 0 {
			rt.Abort("sma._creditStorageDeals: no new payment to be credited.")
		}

		dealP := activeDeal.Deal().Proposal()
		fee := dealStorageFeeForEpochs(dealP, duration)

		if fee > activeDeal.LockedStorageFee() {
			rt.Abort("sma._creditStorageDeals: cannot unlock more than already locked.")
		}

		activeDeal.Impl().LastPaymentEpoch_ = lastChallengeEndEpoch

		// potentially unnecessary as we can unlock funds for provider directly
		// unless the protocol plans on refunding client for MaxFaultCount consecutive fails
		activeDeal.Impl().UnlockStorageFee(fee)

		// TODO: align on provider storage fee withdrawal
	}

}

func (st *StorageMarketActorState_I) _slashDeclaredFaults(rt Runtime, dealIDs []deal.DealID) {

	for _, dealID := range dealIDs {

		deal, found := st.ActiveDeals()[dealID]

		if !found {
			rt.Abort("sm._slashDeclaredFaults: deal not found in ActiveDeals.")
		}

		// TODO: the exact slash amount is up for change
		// TODO: check if provider runs out of collateral to slash here
		deal.Impl().ProviderCollateralRemaining_ -= deal.Deal().Proposal().ProviderCollateralPerEpoch()

	}
}

func (st *StorageMarketActorState_I) _slashDetectedFaults(rt Runtime, dealIDs []deal.DealID) {

	for _, dealID := range dealIDs {

		deal, found := st.ActiveDeals()[dealID]

		if !found {
			rt.Abort("sm._slashDetectedFaults: deal not found in ActiveDeals.")
		}

		// TODO: the exact slash amount is up for change
		// TODO: check if provider runs out of collateral to slash here
		// TODO: more sever slashing for detected faults vs declared faults
		// maybe for the duration since the last Epoch that the sector was proven
		deal.Impl().ProviderCollateralRemaining_ -= deal.Deal().Proposal().ProviderCollateralPerEpoch()

	}
}

func (st *StorageMarketActorState_I) _terminateDeal(rt Runtime, dealID deal.DealID) {

	deal, found := st.ActiveDeals()[dealID]

	if !found {
		rt.Abort("sm._slashTerminatedFaults: deal not found in ActiveDeals.")
	}

	delete(st.ActiveDeals(), dealID)

	// return client collateral and locked storage fee
	clientAddr := deal.Deal().Proposal().Client()
	clientCollateral := deal.Deal().Proposal().TotalClientCollateral()
	st._unlockBalance(rt, clientAddr, clientCollateral+deal.LockedStorageFee())

	// burn all deal.ProviderCollateralRemaining by sending them to TreasuryActor
	// TODO: Send(deal.ProviderCollateralRemaining)
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

func (a *StorageMarketActorCode_I) ProcessDealSlash(rt Runtime, info deal.BatchDealSlashInfo) {
	TODO() // TODO: only call by StorageMinerActor

	h, st := a.State(rt)

	switch info.Action() {
	case deal.SlashDeclaredFaults:
		st._slashDeclaredFaults(rt, info.DealIDs())
	case deal.SlashDetectedFaults:
		st._slashDetectedFaults(rt, info.DealIDs())
	case deal.SlashTerminatedFaults:
		st._slashTerminatedFaults(rt, info.DealIDs())
	default:
		rt.Abort("sma.ProcessDealSlash: invalid action type")
	}

	UpdateRelease(rt, h, st)

}

func (a *StorageMarketActorCode_I) ProcessDealPayment(rt Runtime, info deal.BatchDealPaymentInfo) {
	h, st := a.State(rt)

	switch info.Action() {
	case deal.ExpireStorageDeals:
		st._expireStorageDeals(rt, info.DealIDs(), info.LastChallengeEndEpoch())
	case deal.CreditStorageDeals:
		st._creditStorageDeals(rt, info.DealIDs(), info.LastChallengeEndEpoch())
	default:
		rt.Abort("sma.ProcessDealPayment: invalid deal payment action.")
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
