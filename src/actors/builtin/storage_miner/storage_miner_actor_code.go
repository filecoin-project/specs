package storage_miner

import (
	"bytes"
	"math/big"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	indices "github.com/filecoin-project/specs/actors/runtime/indices"
	serde "github.com/filecoin-project/specs/actors/serde"
	autil "github.com/filecoin-project/specs/actors/util"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

const epochUndefined = abi.ChainEpoch(-1)

//////////////////
// SurprisePoSt //
//////////////////

// Called by StoragePowerActor to notify StorageMiner of SurprisePoSt Challenge.
func (a *StorageMinerActorCode_I) OnSurprisePoStChallenge(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.StoragePowerActorAddr)

	h, st := a.State(rt)

	// If already challenged, do not challenge again.
	// Failed PoSt will automatically reset the state to not-challenged.
	if st.PoStState().Is_Challenged() {
		Release(rt, h, st)
		return
	}

	// Do not challenge if the last successful PoSt was recent enough.
	noChallengePeriod := indices.StorageMining_PoStNoChallengePeriod()
	if st.PoStState().Is_OK() && st.PoStState().As_OK().LastSuccessfulPoSt() >= rt.CurrEpoch()-noChallengePeriod {
		Release(rt, h, st)
		return
	}

	numConsecutiveFailures := 0
	if st.PoStState().Is_DetectedFault() {
		numConsecutiveFailures = st.PoStState().As_DetectedFault().NumConsecutiveFailures()
	}

	var curRecBuf bytes.Buffer
	err := rt.CurrReceiver().MarshalCBOR(&curRecBuf)
	autil.Assert(err == nil)

	randomnessK := rt.GetRandomness(rt.CurrEpoch() - node_base.SPC_LOOKBACK_POST)
	challengedSectorsRandomness := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_SurprisePoStSampleSectors, randomnessK, rt.CurrReceiver())

	challengedSectors := _surprisePoStSampleChallengedSectors(
		challengedSectorsRandomness,
		SectorNumberSetHAMT_Items(st.ProvingSet()),
	)

	st.Impl().PoStState_ = MinerPoStState_New_Challenged(rt.CurrEpoch(), challengedSectors, numConsecutiveFailures)

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for SurprisePoSt challenge expiry.
	provingPeriod := indices.StorageMining_SurprisePoStProvingPeriod()
	a._rtEnrollCronEvent(rt, rt.CurrEpoch()+provingPeriod, []sector.SectorNumber{})
}

// Invoked by miner's worker address to submit a response to a pending SurprisePoSt challenge.
func (a *StorageMinerActorCode_I) SubmitSurprisePoStResponse(rt Runtime, onChainInfo sector.OnChainSurprisePoStVerifyInfo) {
	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())

	if !st.PoStState().Is_Challenged() {
		rt.AbortStateMsg("Not currently challenged")
	}

	Release(rt, h, st)

	a._rtVerifySurprisePoStOrAbort(rt, onChainInfo)
	a._rtUpdatePoStState(rt, MinerPoStState_New_OK(rt.CurrEpoch()))

	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerSurprisePoStSuccess,
		nil,
		abi.TokenAmount(0),
	)
}

// Called by StoragePowerActor.
func (a *StorageMinerActorCode_I) OnDeleteMiner(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.StoragePowerActorAddr)
	minerAddr := rt.CurrReceiver()
	rt.DeleteActor(minerAddr)
}

//////////////////
// ElectionPoSt //
//////////////////

// Called by the VM interpreter once an ElectionPoSt has been verified.
func (a *StorageMinerActorCode_I) OnVerifiedElectionPoSt(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

	// The receiver must be the miner who produced the block for which this message is created.
	Assert(rt.ToplevelBlockWinner() == rt.CurrReceiver())

	h, st := a.State(rt)
	updateSuccessEpoch := st.PoStState().Is_OK()
	Release(rt, h, st)

	// Advance the timestamp of the most recent PoSt success, provided the miner is currently
	// in normal state. (Cannot do this if SurprisePoSt mechanism already underway.)
	if updateSuccessEpoch {
		a._rtUpdatePoStState(rt, MinerPoStState_New_OK(rt.CurrEpoch()))
	}
}

///////////////////////
// Sector Commitment //
///////////////////////

// Deals must be posted on chain via sma.PublishStorageDeals before PreCommitSector.
// Optimization: PreCommitSector could contain a list of deals that are not published yet.
func (a *StorageMinerActorCode_I) PreCommitSector(rt Runtime, info sector.SectorPreCommitInfo) {
	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())

	if _, found := st.Sectors()[info.SectorNumber()]; found {
		rt.AbortStateMsg("Sector number already exists in table")
	}

	Release(rt, h, st)

	depositReq := rt.CurrIndices().StorageMining_PreCommitDeposit(st.Info().SectorSize(), info.Expiration())
	RT_ConfirmFundsReceiptOrAbort_RefundRemainder(rt, depositReq)

	// Verify deals with StorageMarketActor; abort if this fails.
	// (Note: committed-capacity sectors contain no deals, so in that case verification will pass trivially.)
	rt.Send(
		builtin.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_OnMinerSectorPreCommit_VerifyDealsOrAbort,
		serde.MustSerializeParams(
			info.DealIDs(),
			info,
		),
		abi.TokenAmount(0),
	)

	h, st = a.State(rt)

	newSectorInfo := &SectorOnChainInfo_I{
		State_:            SectorState_PreCommit,
		Info_:             info,
		PreCommitDeposit_: depositReq,
		PreCommitEpoch_:   rt.CurrEpoch(),
		ActivationEpoch_:  epochUndefined,
		DealWeight_:       *big.NewInt(-1),
	}
	st.Sectors()[info.SectorNumber()] = newSectorInfo

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for PreCommit expiry check.
	expiryBound := rt.CurrEpoch() + node_base.MAX_PROVE_COMMIT_SECTOR_EPOCH + 1
	a._rtEnrollCronEvent(rt, expiryBound, []sector.SectorNumber{info.SectorNumber()})

	if info.Expiration() <= rt.CurrEpoch() {
		rt.AbortArgMsg("PreCommit sector must have positive lifetime")
	}

	a._rtEnrollCronEvent(rt, info.Expiration(), []sector.SectorNumber{info.SectorNumber()})
}

func (a *StorageMinerActorCode_I) ProveCommitSector(rt Runtime, info sector.SectorProveCommitInfo) {
	h, st := a.State(rt)
	workerAddr := st.Info().Worker()
	rt.ValidateImmediateCallerIs(workerAddr)

	preCommitSector, found := st.Sectors()[info.SectorNumber()]
	if !found || preCommitSector.State() != SectorState_PreCommit {
		rt.AbortArgMsg("Sector not valid or not in PreCommit state")
	}

	if rt.CurrEpoch() > preCommitSector.PreCommitEpoch()+node_base.MAX_PROVE_COMMIT_SECTOR_EPOCH || rt.CurrEpoch() < preCommitSector.PreCommitEpoch()+node_base.MIN_PROVE_COMMIT_SECTOR_EPOCH {
		rt.AbortStateMsg("Invalid ProveCommitSector epoch")
	}

	TODO()
	// TODO: How are SealEpoch, InteractiveEpoch determined (and intended to be used)?
	// Presumably they cannot be derived from the SectorProveCommitInfo provided by an untrusted party.

	a._rtVerifySealOrAbort(rt, &sector.OnChainSealVerifyInfo_I{
		SealedCID_:        preCommitSector.Info().SealedCID(),
		SealEpoch_:        preCommitSector.Info().SealEpoch(),
		InteractiveEpoch_: info.InteractiveEpoch(),
		Proof_:            info.Proof(),
		DealIDs_:          preCommitSector.Info().DealIDs(),
		SectorNumber_:     preCommitSector.Info().SectorNumber(),
	})

	UpdateRelease(rt, h, st)

	// Check (and activate) storage deals associated to sector. Abort if checks failed.
	rt.Send(
		builtin.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_OnMinerSectorProveCommit_VerifyDealsOrAbort,
		serde.MustSerializeParams(
			preCommitSector.Info().DealIDs(),
			info,
		),
		abi.TokenAmount(0),
	)

	res := rt.SendQuery(
		builtin.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_GetWeightForDealSet,
		serde.MustSerializeParams(
			preCommitSector.Info().DealIDs(),
		),
	)
	var dealWeight *big.Int
	err := serde.Deserialize(res, dealWeight)
	Assert(err == nil)

	h, st = a.State(rt)

	st.Sectors()[info.SectorNumber()] = &SectorOnChainInfo_I{
		State_:           SectorState_Active,
		Info_:            preCommitSector.Info(),
		PreCommitEpoch_:  preCommitSector.PreCommitEpoch(),
		ActivationEpoch_: rt.CurrEpoch(),
		DealWeight_:      *dealWeight,
	}

	st.ProvingSet()[info.SectorNumber()] = true

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for sector expiry.
	a._rtEnrollCronEvent(
		rt, preCommitSector.Info().Expiration(), []sector.SectorNumber{info.SectorNumber()})

	// Notify SPA to update power associated to newly activated sector.
	storageWeightDesc := a._rtGetStorageWeightDescForSector(rt, info.SectorNumber())
	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnSectorProveCommit,
		serde.MustSerializeParams(
			storageWeightDesc,
		),
		abi.TokenAmount(0),
	)

	// Return PreCommit deposit to worker upon successful ProveCommit.
	rt.SendFunds(workerAddr, preCommitSector.PreCommitDeposit())
}

/////////////////////////
// Sector Modification //
/////////////////////////

func (a *StorageMinerActorCode_I) ExtendSectorExpiration(rt Runtime, sectorNumber sector.SectorNumber, newExpiration abi.ChainEpoch) {
	storageWeightDescPrev := a._rtGetStorageWeightDescForSector(rt, sectorNumber)

	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())

	sectorInfo, found := st.Sectors()[sectorNumber]
	if !found {
		rt.AbortStateMsg("Sector not found")
	}

	extensionLength := newExpiration - sectorInfo.Info().Expiration()
	if extensionLength < 0 {
		rt.AbortStateMsg("Cannot reduce sector expiration")
	}

	sectorInfo.Info().Impl().Expiration_ = newExpiration
	st.Sectors()[sectorNumber] = sectorInfo
	UpdateRelease(rt, h, st)

	storageWeightDescNew := storageWeightDescPrev
	storageWeightDescNew.Impl().Duration_ = storageWeightDescPrev.Duration() + extensionLength

	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnSectorModifyWeightDesc,
		serde.MustSerializeParams(
			storageWeightDescPrev,
			storageWeightDescNew,
		),
		abi.TokenAmount(0),
	)
}

func (a *StorageMinerActorCode_I) TerminateSector(rt Runtime, sectorNumber sector.SectorNumber) {
	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())
	Release(rt, h, st)

	a._rtTerminateSector(rt, sectorNumber, autil.SectorTerminationType_UserTermination)
}

////////////
// Faults //
////////////

func (a *StorageMinerActorCode_I) DeclareTemporaryFaults(rt Runtime, sectorNumbers []sector.SectorNumber, duration abi.ChainEpoch) {
	if duration <= abi.ChainEpoch(0) {
		rt.AbortArgMsg("Temporary fault duration must be positive")
	}

	storageWeightDescs := a._rtGetStorageWeightDescsForSectors(rt, sectorNumbers)
	requiredFee := rt.CurrIndices().StorageMining_TemporaryFaultFee(storageWeightDescs, duration)

	RT_ConfirmFundsReceiptOrAbort_RefundRemainder(rt, requiredFee)
	rt.SendFunds(builtin.BurntFundsActorAddr, requiredFee)

	effectiveBeginEpoch := rt.CurrEpoch() + indices.StorageMining_DeclaredFaultEffectiveDelay()
	effectiveEndEpoch := effectiveBeginEpoch + duration

	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())

	for _, sectorNumber := range sectorNumbers {
		sectorInfo, found := st.Sectors()[sectorNumber]
		if !found || sectorInfo.State() != SectorState_Active {
			continue
		}

		sectorInfo.Impl().State_ = SectorState_TemporaryFault
		sectorInfo.Impl().DeclaredFaultEpoch_ = rt.CurrEpoch()
		sectorInfo.Impl().DeclaredFaultDuration_ = duration
		st.Sectors()[sectorNumber] = sectorInfo
	}

	UpdateRelease(rt, h, st)

	// Request deferred Cron invocation to update temporary fault state.
	a._rtEnrollCronEvent(rt, effectiveBeginEpoch, sectorNumbers)
	a._rtEnrollCronEvent(rt, effectiveEndEpoch, sectorNumbers)
}

//////////
// Cron //
//////////

func (a *StorageMinerActorCode_I) OnDeferredCronEvent(rt Runtime, sectorNumbers []sector.SectorNumber) {
	rt.ValidateImmediateCallerIs(builtin.StoragePowerActorAddr)

	for _, sectorNumber := range sectorNumbers {
		a._rtCheckTemporaryFaultEvents(rt, sectorNumber)
		a._rtCheckSectorExpiry(rt, sectorNumber)
	}

	a._rtCheckSurprisePoStExpiry(rt)
}

/////////////////
// Constructor //
/////////////////

func (a *StorageMinerActorCode_I) Constructor(
	rt Runtime, ownerAddr addr.Address, workerAddr addr.Address, sectorSize sector.SectorSize, peerId peer.ID) {

	rt.ValidateImmediateCallerIs(builtin.StoragePowerActorAddr)
	h := rt.AcquireState()

	st := &StorageMinerActorState_I{
		Sectors_:    SectorsAMT_Empty(),
		PoStState_:  MinerPoStState_New_OK(rt.CurrEpoch()),
		ProvingSet_: SectorNumberSetHAMT_Empty(),
		Info_:       MinerInfo_New(ownerAddr, workerAddr, sectorSize, peerId),
	}

	UpdateRelease(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMinerActorCode_I) _rtCheckTemporaryFaultEvents(rt Runtime, sectorNumber sector.SectorNumber) {
	h, st := a.State(rt)
	checkSector, found := st.Sectors()[sectorNumber]
	Release(rt, h, st)

	if !found {
		return
	}

	storageWeightDesc := a._rtGetStorageWeightDescForSector(rt, sectorNumber)

	if checkSector.State() == SectorState_Active && rt.CurrEpoch() == checkSector.EffectiveFaultBeginEpoch() {
		checkSector.Impl().State_ = SectorState_TemporaryFault

		rt.Send(
			builtin.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveBegin,
			serde.MustSerializeParams(
				storageWeightDesc,
			),
			abi.TokenAmount(0),
		)

		delete(st.ProvingSet(), sectorNumber)
	}

	if checkSector.Is_TemporaryFault() && rt.CurrEpoch() == checkSector.EffectiveFaultEndEpoch() {
		checkSector.Impl().State_ = SectorState_Active
		checkSector.Impl().DeclaredFaultEpoch_ = epochUndefined
		checkSector.Impl().DeclaredFaultDuration_ = epochUndefined

		rt.Send(
			builtin.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd,
			serde.MustSerializeParams(
				storageWeightDesc,
			),
			abi.TokenAmount(0),
		)

		st.ProvingSet()[sectorNumber] = true
	}

	h, st = a.State(rt)
	st.Sectors()[sectorNumber] = checkSector
	UpdateRelease(rt, h, st)
}

func (a *StorageMinerActorCode_I) _rtCheckSectorExpiry(rt Runtime, sectorNumber sector.SectorNumber) {
	h, st := a.State(rt)
	checkSector, found := st.Sectors()[sectorNumber]
	Release(rt, h, st)

	if !found {
		return
	}

	if checkSector.State() == SectorState_PreCommit {
		if rt.CurrEpoch()-checkSector.PreCommitEpoch() > node_base.MAX_PROVE_COMMIT_SECTOR_EPOCH {
			a._rtDeleteSectorEntry(rt, sectorNumber)
			rt.SendFunds(builtin.BurntFundsActorAddr, checkSector.PreCommitDeposit())
		}
		return
	}

	// Note: the following test may be false, if sector expiration has been extended by the worker
	// in the interim after the Cron request was enrolled.
	if rt.CurrEpoch() >= checkSector.Info().Expiration() {
		a._rtTerminateSector(rt, sectorNumber, autil.SectorTerminationType_NormalExpiration)
	}
}

func (a *StorageMinerActorCode_I) _rtTerminateSector(rt Runtime, sectorNumber sector.SectorNumber, terminationType SectorTerminationType) {
	h, st := a.State(rt)
	checkSector, found := st.Sectors()[sectorNumber]
	Assert(found)
	Release(rt, h, st)

	storageWeightDesc := a._rtGetStorageWeightDescForSector(rt, sectorNumber)

	if checkSector.State() == SectorState_TemporaryFault {
		// To avoid boundary-case errors in power accounting, make sure we explicitly end
		// the temporary fault state first, before terminating the sector.
		rt.Send(
			builtin.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd,
			serde.MustSerializeParams(
				storageWeightDesc,
			),
			abi.TokenAmount(0),
		)
	}

	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnSectorTerminate,
		serde.MustSerializeParams(
			storageWeightDesc,
			terminationType,
		),
		abi.TokenAmount(0),
	)

	a._rtDeleteSectorEntry(rt, sectorNumber)
	delete(st.ProvingSet(), sectorNumber)
}

func (a *StorageMinerActorCode_I) _rtCheckSurprisePoStExpiry(rt Runtime) {
	rt.ValidateImmediateCallerIs(builtin.StoragePowerActorAddr)

	h, st := a.State(rt)

	if !st.PoStState().Is_Challenged() {
		// Already exited challenged state successfully prior to expiry.
		Release(rt, h, st)
		return
	}

	provingPeriod := indices.StorageMining_SurprisePoStProvingPeriod()
	if rt.CurrEpoch() < st.PoStState().As_Challenged().SurpriseChallengeEpoch()+provingPeriod {
		// Challenge not yet expired.
		Release(rt, h, st)
		return
	}

	numConsecutiveFailures := st.PoStState().As_Challenged().NumConsecutiveFailures() + 1

	Release(rt, h, st)

	if numConsecutiveFailures > indices.StoragePower_SurprisePoStMaxConsecutiveFailures() {
		// Terminate all sectors, notify power and market actors to terminate
		// associated storage deals, and reset miner's PoSt state to OK.
		terminatedSectors := []sector.SectorNumber{}
		for sectorNumber := range st.Sectors() {
			terminatedSectors = append(terminatedSectors, sectorNumber)
		}
		a._rtNotifyMarketForTerminatedSectors(rt, terminatedSectors)
	} else {
		// Increment count of consecutive failures, and continue.
		h, st = a.State(rt)
		st.Impl().PoStState_ = MinerPoStState_New_DetectedFault(numConsecutiveFailures)
		UpdateRelease(rt, h, st)
	}

	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerSurprisePoStFailure,
		serde.MustSerializeParams(
			numConsecutiveFailures,
		),
		abi.TokenAmount(0))
}

func (a *StorageMinerActorCode_I) _rtEnrollCronEvent(
	rt Runtime, eventEpoch abi.ChainEpoch, sectorNumbers []sector.SectorNumber) {

	rt.Send(
		builtin.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerEnrollCronEvent,
		serde.MustSerializeParams(
			eventEpoch,
			sectorNumbers,
		),
		abi.TokenAmount(0),
	)
}

func (a *StorageMinerActorCode_I) _rtDeleteSectorEntry(rt Runtime, sectorNumber sector.SectorNumber) {
	h, st := a.State(rt)
	delete(st.Sectors(), sectorNumber)
	UpdateRelease(rt, h, st)
}

func (a *StorageMinerActorCode_I) _rtUpdatePoStState(rt Runtime, state MinerPoStState) {
	h, st := a.State(rt)
	st.Impl().PoStState_ = state
	UpdateRelease(rt, h, st)
}

func (a *StorageMinerActorCode_I) _rtGetStorageWeightDescForSector(
	rt Runtime, sectorNumber sector.SectorNumber) autil.SectorStorageWeightDesc {

	h, st := a.State(rt)
	ret := st._getStorageWeightDescForSector(sectorNumber)
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _rtGetStorageWeightDescsForSectors(
	rt Runtime, sectorNumbers []sector.SectorNumber) []autil.SectorStorageWeightDesc {

	h, st := a.State(rt)
	ret := st._getStorageWeightDescsForSectors(sectorNumbers)
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _rtNotifyMarketForTerminatedSectors(rt Runtime, sectorNumbers []sector.SectorNumber) {
	h, st := a.State(rt)
	dealIDItems := []deal.DealID{}
	for _, sectorNo := range sectorNumbers {
		dealIDItems = append(dealIDItems, st._getSectorDealIDsAssert(sectorNo).Items()...)
	}
	dealIDs := &deal.DealIDs_I{Items_: dealIDItems}

	Release(rt, h, st)

	rt.Send(
		builtin.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_OnMinerSectorsTerminate,
		serde.MustSerializeParams(
			dealIDs,
		),
		abi.TokenAmount(0),
	)
}

func (a *StorageMinerActorCode_I) _rtVerifySurprisePoStOrAbort(rt Runtime, onChainInfo sector.OnChainSurprisePoStVerifyInfo) {
	h, st := a.State(rt)
	info := st.Info()

	Assert(st.PoStState().Is_Challenged())
	challengeEpoch := st.PoStState().As_Challenged().SurpriseChallengeEpoch()
	challengedSectors := st.PoStState().As_Challenged().ChallengedSectors()

	// verify no duplicate tickets
	challengeIndices := make(map[uint64]bool)
	for _, tix := range onChainInfo.Candidates() {
		if _, ok := challengeIndices[tix.ChallengeIndex()]; ok {
			rt.AbortStateMsg("Invalid Surprise PoSt. Duplicate ticket included.")
		}
		challengeIndices[tix.ChallengeIndex()] = true
	}

	TODO(challengedSectors)
	// TODO: Determine what should be the acceptance criterion for sector numbers
	// proven in SurprisePoSt proofs.
	//
	// Previous note:
	// Verify the partialTicket values
	// if !a._rtVerifySurprisePoStMeetsTargetReq(rt) {
	// 	rt.AbortStateMsg("Invalid Surprise PoSt. Tickets do not meet target.")
	// }

	randomnessK := rt.GetRandomness(challengeEpoch - node_base.SPC_LOOKBACK_POST)
	// regenerate randomness used. The PoSt Verification below will fail if
	// the same was not used to generate the proof
	postRandomness := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_SurprisePoStChallengeSeed, randomnessK, rt.CurrReceiver())

	UpdateRelease(rt, h, st)

	// Get public inputs

	sectorSize := info.SectorSize()
	postCfg := filproofs.SurprisePoStCfg(sectorSize)

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		Randomness_: onChainInfo.Randomness(),
		// EligibleSectors_: FIXME: verification needs these.
	}

	pv := filproofs.MakeSurprisePoStVerifier(postCfg)

	// Verify the PoSt Proof
	isVerified := pv.VerifySurprisePoSt(&pvInfo)

	if !isVerified {
		rt.AbortStateMsg("Surprise PoSt failed to verify")
	}
}

func (a *StorageMinerActorCode_I) _rtVerifySealOrAbort(rt Runtime, onChainInfo sector.OnChainSealVerifyInfo) {
	h, st := a.State(rt)
	info := st.Info()
	sectorSize := info.SectorSize()
	Release(rt, h, st)

	pieceInfos, err := sector.Deserialize_PieceInfos(rt.SendQuery(
		builtin.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_GetPieceInfosForDealIDs,
		serde.MustSerializeParams(
			sectorSize,
			onChainInfo.DealIDs(),
		),
	))
	Assert(err == nil)

	// Unless we enforce a minimum padding amount, this totalPieceSize calculation can be removed.
	// Leaving for now until that decision is entirely finalized.
	var totalPieceSize uint64
	for _, pieceInfo := range pieceInfos.Items() {
		pieceSize := pieceInfo.Size()
		totalPieceSize += pieceSize
	}

	unsealedCID, _ := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos.Items())

	sealCfg := sector.SealCfg_I{
		SectorSize_:  sectorSize,
		InstanceCfg_: info.InstanceCfg(),
		Partitions_:  info.SealPartitions(),
	}

	minerActorID, err := addr.IDFromAddress(rt.CurrReceiver())
	if err != nil {
		rt.AbortStateMsg("receiver must be ID address")
	}

	IMPL_TODO() // Use randomness APIs
	var svInfoRandomness abi.Randomness
	var svInfoInteractiveRandomness abi.Randomness

	svInfo := sector.SealVerifyInfo_I{
		SectorID_: &sector.SectorID_I{
			Miner_:  abi.ActorID(minerActorID),
			Number_: onChainInfo.SectorNumber(),
		},
		OnChain_: onChainInfo,

		// TODO: Make SealCfg sector.SealCfg from miner configuration (where is that?)
		SealCfg_: &sealCfg,

		Randomness_:            sector.SealRandomness(svInfoRandomness),
		InteractiveRandomness_: sector.InteractiveSealRandomness(svInfoInteractiveRandomness),
		UnsealedCID_:           unsealedCID,
	}

	sdr := filproofs.WinSDRParams(&sealCfg)

	isVerified := sdr.VerifySeal(&svInfo)

	if !isVerified {
		rt.AbortStateMsg("Sector seal failed to verify")
	}
}

func getSectorNums(m map[sector.SectorNumber]SectorOnChainInfo) []sector.SectorNumber {
	var l []sector.SectorNumber
	for i, _ := range m {
		l = append(l, i)
	}
	return l
}

func _surprisePoStSampleChallengedSectors(
	sampleRandomness abi.Randomness, provingSet []sector.SectorNumber) []sector.SectorNumber {

	IMPL_TODO()
	panic("")
}
