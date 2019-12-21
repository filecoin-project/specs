package storage_mining

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMinerActorCode_I) GetWorkerVRFKey(rt Runtime) filcrypto.VRFPublicKey {
	h, st := a.State(rt)
	ret := st.Info().WorkerVRFKey()
	Release(rt, h, st)
	return ret
}

//////////////////
// SurprisePoSt //
//////////////////

// Called by StoragePowerActor to notify StorageMiner of PoSt Challenge.
func (a *StorageMinerActorCode_I) OnSurprisePoStChallenge(rt Runtime) {
	rt.ValidateImmediateCallerIs(addr.StoragePowerActorAddr)

	h, st := a.State(rt)

	// If already challenged, do not challenge again.
	// Failed PoSt will automatically reset the state to not-challenged.
	if st.PoStState().Is_Challenged() {
		Release(rt, h, st)
		return
	}

	// Do not challenge if the last successful PoSt was recent enough.
	if st.PoStState().Is_OK() && st.PoStState().As_OK().LastSuccessfulPoSt() >= rt.CurrEpoch()-SURPRISE_NO_CHALLENGE_PERIOD {
		Release(rt, h, st)
		return
	}

	numConsecutiveFailures := 0
	if st.PoStState().Is_Failing() {
		numConsecutiveFailures = st.PoStState().As_Failing().NumConsecutiveFailures()
	}

	st.Impl().PoStState_ = MinerPoStState_New_Challenged(rt.CurrEpoch(), numConsecutiveFailures)

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for SurprisePoSt challenge expiry.
	a._rtEnrollCronEvent(rt, rt.CurrEpoch()+PROVING_PERIOD, []sector.SectorNumber{})
}

// Invoked by miner's worker address to submit a response to a pending Surprise PoSt challenge.
func (a *StorageMinerActorCode_I) SubmitSurprisePoStResponse(rt Runtime, onChainInfo sector.OnChainPoStVerifyInfo) {
	h, st := a.State(rt)
	rt.ValidateImmediateCallerIs(st.Info().Worker())

	if !st.PoStState().Is_Challenged() {
		rt.AbortStateMsg("Not currently challenged")
	}

	Release(rt, h, st)

	a._rtVerifySurprisePoStOrAbort(rt, onChainInfo)
	a._rtUpdatePoStState(rt, MinerPoStState_New_OK(rt.CurrEpoch()))

	rt.Send(
		addr.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerSurprisePoStSuccess,
		[]util.Serialization{},
		actor.TokenAmount(0),
	)
}

//////////////////
// ElectionPoSt //
//////////////////

// Called by the VM interpreter once an ElectionPoSt has been verified.
// Assumes the block reward has already been granted to the storage miner actor.
// This only handles sector management.
func (a *StorageMinerActorCode_I) OnVerifiedElectionPoSt(rt Runtime) {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)

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
	rt.Send(
		addr.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_OnMinerSectorPreCommit_VerifyDealsOrAbort,
		[]util.Serialization{
			deal.Serialize_DealIDs(info.DealIDs()),
			sector.Serialize_SectorPreCommitInfo(info),
		},
		actor.TokenAmount(0),
	)

	h, st = a.State(rt)

	newSectorInfo := &SectorOnChainInfo_I{
		State_:            SectorState_PreCommit,
		Info_:             info,
		PreCommitDeposit_: depositReq,
		PreCommitEpoch_:   rt.CurrEpoch(),
		ProveCommitEpoch_: block.ChainEpoch(-1),
		DealWeight_:       util.BigFromInt(-1),
	}
	st.Sectors()[info.SectorNumber()] = newSectorInfo

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for PreCommit expiry check.
	expiryBound := rt.CurrEpoch() + sector.MAX_PROVE_COMMIT_SECTOR_EPOCH + 1
	a._rtEnrollCronEvent(rt, expiryBound, []sector.SectorNumber{info.SectorNumber()})

	if info.Expiration() > rt.CurrEpoch() {
		// Note: sector expiration may be in the past, in the case of committed capacity.
		a._rtEnrollCronEvent(rt, info.Expiration(), []sector.SectorNumber{info.SectorNumber()})
	}
}

func (a *StorageMinerActorCode_I) ProveCommitSector(rt Runtime, info sector.SectorProveCommitInfo) {
	h, st := a.State(rt)
	sectorSize := st.Info().SectorSize()
	workerAddr := st.Info().Worker()
	rt.ValidateImmediateCallerIs(workerAddr)

	preCommitSector, found := st.Sectors()[info.SectorNumber()]
	if !found || preCommitSector.State() != SectorState_PreCommit {
		rt.AbortArgMsg("Sector not valid or not in PreCommit state")
	}

	if rt.CurrEpoch() > preCommitSector.PreCommitEpoch()+sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
		rt.AbortStateMsg("Deadline exceeded for ProveCommitSector")
	}

	TODO() // How are SealEpoch, InteractiveEpoch determined?

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
		addr.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_OnMinerSectorProveCommit_VerifyDealsOrAbort,
		[]util.Serialization{
			deal.Serialize_DealIDs(preCommitSector.Info().DealIDs()),
			sector.Serialize_SectorProveCommitInfo(info),
		},
		actor.TokenAmount(0),
	)

	dealWeight, ok := util.Deserialize_BigInt(rt.SendQuery(
		addr.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_GetWeightForDealSet,
		[]util.Serialization{
			deal.Serialize_DealIDs(preCommitSector.Info().DealIDs()),
		},
	))
	Assert(ok)

	h, st = a.State(rt)

	st.Sectors()[info.SectorNumber()] = &SectorOnChainInfo_I{
		State_:            SectorState_Active,
		Info_:             preCommitSector.Info(),
		PreCommitEpoch_:   preCommitSector.PreCommitEpoch(),
		ProveCommitEpoch_: rt.CurrEpoch(),
		DealWeight_:       dealWeight,
	}

	UpdateRelease(rt, h, st)

	// Request deferred Cron check for sector expiry.
	a._rtEnrollCronEvent(rt, preCommitSector.Info().Expiration(), []sector.SectorNumber{})

	// Notify SPA to update power associated to newly activated sector.
	rt.Send(
		addr.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnSectorProveCommit,
		[]util.Serialization{
			sector.Serialize_SectorSize(sectorSize),
			util.Serialize_BigInt(dealWeight),
		},
		actor.TokenAmount(0),
	)

	// Return PreCommit deposit to worker upon successful ProveCommit.
	rt.SendFunds(workerAddr, preCommitSector.PreCommitDeposit())
}

////////////
// Faults //
////////////

func (a *StorageMinerActorCode_I) DeclareTemporaryFaults(rt Runtime, sectorNumbers []sector.SectorNumber, duration block.ChainEpoch) {
	if duration <= block.ChainEpoch(0) {
		rt.AbortArgMsg("Temporary fault duration must be positive")
	}

	storageWeightDescs := a._rtGetStorageWeightDescsForSectors(rt, sectorNumbers)
	requiredFee := rt.CurrIndices().StorageMining_TemporaryFaultFee(storageWeightDescs, duration)

	RT_ConfirmFundsReceiptOrAbort_RefundRemainder(rt, requiredFee)
	rt.SendFunds(addr.BurntFundsActorAddr, requiredFee)

	effectiveBeginEpoch := rt.CurrEpoch() + DECLARED_FAULT_EFFECTIVE_DELAY
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

func (a *StorageMinerActorCode_I) OnDeferredCronEvent(rt Runtime, sectorNumbers []sector.SectorNumber) {
	rt.ValidateImmediateCallerIs(addr.StoragePowerActorAddr)

	for _, sectorNumber := range sectorNumbers {
		a._rtCheckTemporaryFaultEvents(rt, sectorNumber)
		a._rtCheckSectorExpiry(rt, sectorNumber)
	}

	a._rtCheckSurprisePoStExpiry(rt)
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
			addr.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveBegin,
			[]util.Serialization{actor_util.Serialize_SectorStorageWeightDesc(storageWeightDesc)},
			actor.TokenAmount(0),
		)
	}

	if checkSector.Is_TemporaryFault() && rt.CurrEpoch() == checkSector.EffectiveFaultEndEpoch() {
		checkSector.Impl().State_ = SectorState_Active
		checkSector.Impl().DeclaredFaultEpoch_ = block.ChainEpoch_None
		checkSector.Impl().DeclaredFaultDuration_ = block.ChainEpoch_None

		rt.Send(
			addr.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd,
			[]util.Serialization{actor_util.Serialize_SectorStorageWeightDesc(storageWeightDesc)},
			actor.TokenAmount(0),
		)
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
		if rt.CurrEpoch()-checkSector.PreCommitEpoch() > sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
			a._rtDeleteSectorEntry(rt, sectorNumber)
		}
		rt.SendFunds(addr.BurntFundsActorAddr, checkSector.PreCommitDeposit())
		return
	}

	// Note: the following test may be false, if sector expiration has been extended by the worker
	// in the interim after the Cron request was enrolled.
	if rt.CurrEpoch() >= checkSector.Info().Expiration() {
		storageWeightDesc := a._rtGetStorageWeightDescForSector(rt, sectorNumber)

		if checkSector.State() == SectorState_TemporaryFault {
			// To avoid boundary-case errors in power accounting, make sure we explicitly end
			// the temporary fault state first, before terminating the sector.
			rt.Send(
				addr.StoragePowerActorAddr,
				ai.Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd,
				[]util.Serialization{actor_util.Serialize_SectorStorageWeightDesc(storageWeightDesc)},
				actor.TokenAmount(0),
			)
		}

		rt.Send(
			addr.StoragePowerActorAddr,
			ai.Method_StoragePowerActor_OnSectorTerminate,
			[]util.Serialization{actor_util.Serialize_SectorStorageWeightDesc(storageWeightDesc)},
			actor.TokenAmount(0),
		)

		a._rtDeleteSectorEntry(rt, sectorNumber)
	}
}

func (a *StorageMinerActorCode_I) _rtCheckSurprisePoStExpiry(rt Runtime) {
	TODO() // Fill in from constants
	PROVING_PERIOD := block.ChainEpoch(0)
	MAX_SURPRISE_CONSECUTIVE_FAILURES := 5

	rt.ValidateImmediateCallerIs(addr.StoragePowerActorAddr)

	h, st := a.State(rt)

	if !st.PoStState().Is_Challenged() {
		// Already exited challenged state successfully prior to expiry.
		Release(rt, h, st)
		return
	}

	if rt.CurrEpoch() < st.PoStState().As_Challenged().SurpriseChallengeEpoch()+PROVING_PERIOD {
		// Challenge not yet expired.
		Release(rt, h, st)
		return
	}

	numConsecutiveFailures := st.PoStState().As_Challenged().NumConsecutiveFailures() + 1

	Release(rt, h, st)

	rt.Send(
		addr.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerSurprisePoStFailure,
		[]util.Serialization{
			util.Serialize_Int(numConsecutiveFailures),
		},
		actor.TokenAmount(0))

	if numConsecutiveFailures > MAX_SURPRISE_CONSECUTIVE_FAILURES {
		// Terminate all sectors, notify power and market actors to terminate
		// associated storage deals, and reset miner's PoSt state to OK.
		terminatedSectors := []sector.SectorNumber{}
		for sectorNumber := range st.Sectors() {
			terminatedSectors = append(terminatedSectors, sectorNumber)
		}
		a._rtNotifyForTerminatedSectors(rt, terminatedSectors, actor_util.SectorTerminationType_SurprisePoStFailure)

		h, st := a.State(rt)
		st.Impl().Sectors_ = SectorsAMT_Empty()
		st.Impl().PoStState_ = MinerPoStState_New_OK(rt.CurrEpoch())
		UpdateRelease(rt, h, st)
	} else {
		// Increment count of consecutive failures, and continue.
		h, st = a.State(rt)
		st.Impl().PoStState_ = MinerPoStState_New_Failing(numConsecutiveFailures)
		UpdateRelease(rt, h, st)
	}
}

func (a *StorageMinerActorCode_I) _rtEnrollCronEvent(
	rt Runtime, eventEpoch block.ChainEpoch, sectorNumbers []sector.SectorNumber) {

	rt.Send(
		addr.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnMinerEnrollCronEvent,
		[]util.Serialization{
			block.Serialize_ChainEpoch(eventEpoch),
			sector.Serialize_SectorNumber_Array(sectorNumbers),
		},
		actor.TokenAmount(0),
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
	rt Runtime, sectorNumber sector.SectorNumber) actor_util.SectorStorageWeightDesc {

	h, st := a.State(rt)
	ret := st._getStorageWeightDescForSector(sectorNumber)
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _rtGetStorageWeightDescsForSectors(
	rt Runtime, sectorNumbers []sector.SectorNumber) []actor_util.SectorStorageWeightDesc {

	h, st := a.State(rt)
	ret := st._getStorageWeightDescsForSectors(sectorNumbers)
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _rtProcessTemporaryFaultEnd(rt Runtime, sectorNumber sector.SectorNumber) {
	TODO()
}

func (a *StorageMinerActorCode_I) _rtNotifyForTerminatedSectors(
	rt Runtime, sectorNumbers []sector.SectorNumber, terminationType SectorTerminationType) {

	// Notify StoragePowerActor to adjust power.
	storageWeightDescs := a._rtGetStorageWeightDescsForSectors(rt, sectorNumbers)

	rt.Send(
		addr.StoragePowerActorAddr,
		ai.Method_StoragePowerActor_OnSectorTerminate,
		[]util.Serialization{
			actor_util.Serialize_SectorStorageWeightDesc_Array(storageWeightDescs),
		},
		actor.TokenAmount(0),
	)

	// If termination is not via normal expiration, then must also notify StorageMarketActor
	// to terminate associated storage deals.
	if terminationType != actor_util.SectorTerminationType_NormalExpiration {
		h, st := a.State(rt)
		dealIDItems := []deal.DealID{}
		for _, sectorNo := range sectorNumbers {
			dealIDItems = append(dealIDItems, st._getSectorDealIDsAssert(sectorNo).Items()...)
		}
		dealIDs := &deal.DealIDs_I{Items_: dealIDItems}

		Release(rt, h, st)

		rt.Send(
			addr.StorageMarketActorAddr,
			ai.Method_StorageMarketActor_OnMinerSectorsTerminate,
			[]util.Serialization{
				deal.Serialize_DealIDs(dealIDs),
			},
			actor.TokenAmount(0),
		)
	}
}

func (a *StorageMinerActorCode_I) _rtVerifySurprisePoStOrAbort(rt Runtime, onChainInfo sector.OnChainPoStVerifyInfo) {
	h, st := a.State(rt)
	info := st.Info()

	// Verify the partialTicket values
	if !a._rtVerifySurprisePoStMeetsTargetReq(rt) {
		rt.AbortStateMsg("Invalid Surprise PoSt. Tickets do not meet target.")
	}

	// verify the partialTickets themselves
	// Verify appropriate randomness

	TODO() // pull from consts
	SPC_LOOKBACK_POST := block.ChainEpoch(0)

	Assert(st.PoStState().Is_Challenged())
	challengeEpoch := st.PoStState().As_Challenged().SurpriseChallengeEpoch()
	randomnessEpoch := challengeEpoch - SPC_LOOKBACK_POST

	TODO() // extract randomness
	util.PARAM_FINISH(randomnessEpoch)
	var postRandomnessInput util.Randomness // sms.PreparePoStChallengeSeed(randomness, actorAddr)

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	if !postRand.Verify(postRandomnessInput, info.WorkerVRFKey()) {
		rt.AbortStateMsg("Invalid Surprise PoSt. Invalid randomness.")
	}

	UpdateRelease(rt, h, st)

	// Get public inputs
	postCfg := sector.PoStCfg_I{
		Type_:        sector.PoStType_SurprisePoSt,
		SectorSize_:  info.SectorSize(),
		WindowCount_: info.WindowCount(),
		Partitions_:  info.SurprisePoStPartitions(),
	}

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		PoStCfg_:    &postCfg,
		Randomness_: onChainInfo.Randomness(),
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{SurprisePoStCfg_: &postCfg})

	// Verify the PoSt Proof
	isVerified := sdr.VerifySurprisePoSt(&pvInfo)

	if !isVerified {
		rt.AbortStateMsg("Surprise PoSt failed to verify")
	}
}

// todo: define target
func (a *StorageMinerActorCode_I) _rtVerifySurprisePoStMeetsTargetReq(rt Runtime) bool {
	util.TODO()
	return false
}

func (a *StorageMinerActorCode_I) _rtVerifySealOrAbort(rt Runtime, onChainInfo sector.OnChainSealVerifyInfo) {
	h, st := a.State(rt)
	info := st.Info()
	sectorSize := info.SectorSize()
	Release(rt, h, st)

	pieceInfos, err := sector.Deserialize_PieceInfos(rt.SendQuery(
		addr.StorageMarketActorAddr,
		ai.Method_StorageMarketActor_GetPieceInfosForDealIDs,
		[]util.Serialization{
			sector.Serialize_SectorSize(sectorSize),
			deal.Serialize_DealIDs(onChainInfo.DealIDs()),
		},
	))
	Assert(err == nil)

	// Unless we enforce a minimum padding amount, this totalPieceSize calculation can be removed.
	// Leaving for now until that decision is entirely finalized.
	var totalPieceSize util.UInt
	for _, pieceInfo := range pieceInfos.Items() {
		pieceSize := pieceInfo.Size()
		totalPieceSize += pieceSize
	}

	unsealedCID, _ := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos.Items())

	sealCfg := sector.SealCfg_I{
		SectorSize_:  sectorSize,
		WindowCount_: info.WindowCount(),
		Partitions_:  info.SealPartitions(),
	}

	minerActorID, err := rt.CurrReceiver().GetID()
	if err != nil {
		rt.AbortStateMsg("receiver must be ID address")
	}

	svInfo := sector.SealVerifyInfo_I{
		SectorID_: &sector.SectorID_I{
			Miner_:  minerActorID,
			Number_: onChainInfo.SectorNumber(),
		},
		OnChain_: onChainInfo,

		// TODO: Make SealCfg sector.SealCfg from miner configuration (where is that?)
		SealCfg_: &sealCfg,

		Randomness_:            sector.SealRandomness(rt.Randomness(onChainInfo.SealEpoch(), 0)),
		InteractiveRandomness_: sector.InteractiveSealRandomness(rt.Randomness(onChainInfo.InteractiveEpoch(), 0)),
		UnsealedCID_:           unsealedCID,
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{SealCfg_: &sealCfg})

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
