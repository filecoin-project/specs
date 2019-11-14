package storage_mining

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	power "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	storage_market "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	poster "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving/poster"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	util "github.com/filecoin-project/specs/util"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
)

const (
	Method_StorageMinerActor_SubmitSurprisePoSt = actor.MethodPlaceholder
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type State = StorageMinerActorState
type Any = util.Any
type Bool = util.Bool
type Bytes = util.Bytes
type InvocOutput = msg.InvocOutput
type Runtime = vmr.Runtime

var TODO = util.TODO

func (a *StorageMinerActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
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
func (st *StorageMinerActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) State {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (cs *ChallengeStatus_I) OnNewChallenge(currEpoch block.ChainEpoch) ChallengeStatus {
	cs.LastChallengeEpoch_ = currEpoch
	return cs
}

// Call by _onSuccessfulPoSt or _onMissedSurprisePoSt
func (cs *ChallengeStatus_I) OnChallengeResponse(currEpoch block.ChainEpoch) ChallengeStatus {
	cs.LastChallengeEndEpoch_ = currEpoch
	return cs
}

func (cs *ChallengeStatus_I) IsChallenged() bool {
	// true (isChallenged) when LastChallengeEpoch is later than LastChallengeEndEpoch
	return cs.LastChallengeEpoch() > cs.LastChallengeEndEpoch()
}

func (a *StorageMinerActorCode_I) _isChallenged(rt Runtime) bool {
	h, st := a.State(rt)
	ret := st._isChallenged(rt)
	Release(rt, h, st)
	return ret
}

func (st *StorageMinerActorState_I) _isChallenged(rt Runtime) bool {
	return st.ChallengeStatus().IsChallenged()
}

func (cs *ChallengeStatus_I) ShouldChallenge(currEpoch block.ChainEpoch, minChallengePeriod block.ChainEpoch) bool {
	return currEpoch > (cs.LastChallengeEpoch()+minChallengePeriod) && !cs.IsChallenged()
}

func (st *StorageMinerActorState_I) ShouldChallenge(rt Runtime, minChallengePeriod block.ChainEpoch) bool {
	return st.ChallengeStatus().ShouldChallenge(rt.CurrEpoch(), minChallengePeriod)
}

// called by StoragePowerActor to notify StorageMiner of PoSt Challenge (triggered by Cron)
func (a *StorageMinerActorCode_I) NotifyOfSurprisePoStChallenge(rt Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.CronActorAddr)

	if a._isChallenged(rt) {
		return rt.SuccessReturn() // silent return, dont re-challenge
	}

	a._expirePreCommittedSectors(rt)

	h, st := a.State(rt)
	st.ChallengeStatus().Impl().LastChallengeEpoch_ = rt.CurrEpoch()
	UpdateRelease(rt, h, st)
	return rt.SuccessReturn()
}


func (st *StorageMinerActorState_I) _processStagedCommittedSectors(rt Runtime) {
	for sectorNo, stagedInfo := range st.StagedCommittedSectors() {
		st.Sectors()[sectorNo] = stagedInfo.Sector()
		st.Impl().ProvingSet_.Add(sectorNo)
		st.SectorTable().Impl().CommittedSectors_.Add(sectorNo)
		st.SectorUtilization()[sectorNo] = stagedInfo.Utilization()
	}

	// empty StagedCommittedSectors
	st.Impl().StagedCommittedSectors_ = make(map[sector.SectorNumber]StagedCommittedSectorInfo)
}

func (a *StorageMinerActorCode_I) _slashDealsFromFaultReport(rt Runtime, sectorNumbers []sector.SectorNumber, action deal.StorageDealSlashAction) {

	h, st := a.State(rt)

	// This is not right but doing this for specing`
	dealIDs := deal.CompactDealSet(make([]byte, 0))

	for _, sectorNo := range sectorNumbers {

		utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
		activeDealIDs := utilizationInfo.DealExpirationQueue().ActiveDealIDs()
		dealIDs = dealIDs.Extend(activeDealIDs)

	}

	slashInfo := &deal.BatchDealSlashInfo_I{
		DealIDs_: dealIDs.DealsOn(),
		Action_:  action,
	}

	Release(rt, h, st)

	// TODO: Send(StorageMarketActor, ProcessSectorDealSlash)
	panic(slashInfo)

}

// construct FaultReport
// reset NewTerminatedFaults
func (a *StorageMinerActorCode_I) _submitFaultReport(
	rt Runtime,
	newDeclaredFaults sector.CompactSectorSet,
	newDetectedFaults sector.CompactSectorSet,
	newTerminatedFaults sector.CompactSectorSet,
) {
	faultReport := &sector.FaultReport_I{
		NewDeclaredFaults_:   newDeclaredFaults,
		NewDetectedFaults_:   newDetectedFaults,
		NewTerminatedFaults_: newTerminatedFaults,
	}

	rt.Abort("TODO") // TODO: Send(SPA, ProcessFaultReport(faultReport))
	panic(faultReport)

	// only terminatedFault will be slashed
	if len(newTerminatedFaults) > 0 {
		a._slashDealsFromFaultReport(rt, newTerminatedFaults.SectorsOn(), deal.SlashTerminatedFaults)
	}

	h, st := a.State(rt)
	st.SectorTable().Impl().TerminatedFaults_ = sector.CompactSectorSet(make([]byte, 0))
	UpdateRelease(rt, h, st)
}

// construct PowerReport from SectorTable
func (a *StorageMinerActorCode_I) _submitPowerReport(rt Runtime) {
	h, st := a.State(rt)
	newExpiredDealIDs := st._updateSectorUtilization(rt)
	activePower := st._getActivePower(rt)
	inactivePower := st._getInactivePower(rt)

	powerReport := &power.PowerReport_I{
		ActivePower_:   activePower,
		InactivePower_: inactivePower,
	}

	Release(rt, h, st)

	rt.Abort("TODO") // TODO: Send(SPA, ProcessPowerReport(powerReport))
	panic(powerReport)

	if len(newExpiredDealIDs) > 0 {
		batchDealPaymentInfo := &deal.BatchDealPaymentInfo_I{
			DealIDs_:               newExpiredDealIDs,
			Action_:                deal.ExpireStorageDeals,
			LastChallengeEndEpoch_: st.ChallengeStatus().LastChallengeEndEpoch(),
		}

		rt.Abort("TODO")
		panic(batchDealPaymentInfo)
		// Send(StorageMarketActor, ProcessSectorDealPayment(batchDealPaymentInfo))
	}
}

func (a *StorageMinerActorCode_I) _onMissedSurprisePoSt(rt Runtime) {
	h, st := a.State(rt)

	failingSectorNumbers := getSectorNums(st.Sectors())
	for _, sectorNo := range failingSectorNumbers {
		st._updateFailSector(rt, sectorNo, true)
	}
	st._updateExpireSectors(rt)
	UpdateRelease(rt, h, st)

	h, st = a.State(rt)
	newDetectedFaults := st.SectorTable().FailingSectors()
	newTerminatedFaults := st.SectorTable().TerminatedFaults()
	Release(rt, h, st)

	// Note: NewDetectedFaults is now the sum of all
	// previously active, committed, and recovering sectors minus expired ones
	// and any previously Failing sectors that did not exceed MaxFaultCount
	// Note: previously declared faults is now treated as part of detected faults
	a._submitFaultReport(
		rt,
		sector.CompactSectorSet(make([]byte, 0)), // NewDeclaredFaults
		newDetectedFaults,
		newTerminatedFaults,
	)

	a._submitPowerReport(rt)

	// end of challenge
	h, st = a.State(rt)
	st.ChallengeStatus().Impl().OnChallengeResponse(rt.CurrEpoch())
	st._processStagedCommittedSectors(rt)
	UpdateRelease(rt, h, st)
}

// If a Post is missed (either due to faults being not declared on time or
// because the miner run out of time, every sector is reported as failing
// for the current proving period.
func (a *StorageMinerActorCode_I) CheckSurprisePoStSubmissionHappened(rt Runtime) InvocOutput {
	TODO() // TODO: validate caller

	// we can return if miner has not yet gotten the chance to submit a surprise post
	if !a._isChallenged(rt) {
		// Miner gets out of a challenge when submit a successful PoSt
		// or when detected by CronActor. Hence, not being in isChallenged means that we are good here
		return rt.SuccessReturn()
	}

	// garbage collection - need to be called by cron once in a while
	a._expirePreCommittedSectors(rt)

	// oh no -- we missed it. rekt
	a._onMissedSurprisePoSt(rt)

	return rt.SuccessReturn()
}

func (a *StorageMinerActorCode_I) _verifyPoStSubmission(rt Runtime, postSubmission poster.PoStSubmission) bool {
	// 1. A proof must be submitted after the postRandomness for this proving
	// period is on chain
	// if rt.ChainEpoch < sm.ProvingPeriodEnd - challengeTime {
	//   rt.Abort("too early")
	// }

	// 2. A proof must be a valid snark proof with the correct public inputs
	// 2.1 Get randomness from the chain at the right epoch
	// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
	// 2.2 Generate the set of challenges
	// challenges := GenerateChallengesForPoSt(r, keys(sm.Sectors))
	// 2.3 Verify the PoSt Proof
	// verifyPoSt(challenges, TODO)

	rt.Abort("TODO") // TODO: finish
	return false
}

func (a *StorageMinerActorCode_I) _expirePreCommittedSectors(rt Runtime) {

	h, st := a.State(rt)
	for _, preCommitSector := range st.PreCommittedSectors() {

		elapsedEpoch := rt.CurrEpoch() - preCommitSector.ReceivedEpoch()

		if elapsedEpoch > sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
			delete(st.PreCommittedSectors(), preCommitSector.Info().SectorNumber())
			// TODO: potentially some slashing if ProveCommitSector comes late
		}
	}
	UpdateRelease(rt, h, st)

}

func (st *StorageMinerActorState_I) _updateSectorUtilization(rt Runtime) []deal.DealID {
	// TODO: verify if we should update Sector utilization for failing sectors
	// depends on decision around collateral requirement for inactive power
	// and what happens when a failing sector expires

	// this may not work
	ret := deal.CompactDealSet(make([]byte, 0))

	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOn() {

		utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
		newUtilization := utilizationInfo.CurrUtilization()

		currEpoch := rt.CurrEpoch()
		totalDealCount := len(st.Sectors()[sectorNo].SealCommitment().DealIDs())
		newExpiredDealIDs := deal.CompactDealSet(make([]byte, totalDealCount))

		queue := utilizationInfo.DealExpirationQueue()
		for queue.Peek().Expiration() <= currEpoch {
			expiredDeal := queue.Pop()
			newUtilization -= expiredDeal.PayloadPower()

			newExpiredDealIDs.Add(expiredDeal.DealID())
		}

		st.SectorUtilization()[sectorNo].Impl().CurrUtilization_ = newUtilization
		ret = ret.Extend(newExpiredDealIDs)
	}

	return ret.DealsOn()

}

func (st *StorageMinerActorState_I) _getActivePower(rt Runtime) block.StoragePower {
	var activePower = block.StoragePower(0)

	for _, sectorNo := range st.SectorTable().Impl().ActiveSectors_.SectorsOn() {
		utilizationInfo, found := st.SectorUtilization()[sectorNo]
		if !found {
			rt.Abort("sm._getActivePower: sectorNo not found in SectorUtilization")
		}
		activePower += utilizationInfo.CurrUtilization()
	}

	return activePower
}

func (st *StorageMinerActorState_I) _getInactivePower(rt Runtime) block.StoragePower {
	var inactivePower = block.StoragePower(0)

	// iterate over sectorNo in CommittedSectors, RecoveringSectors, and FailingSectors
	inactiveProvingSet := st.SectorTable().Impl().CommittedSectors_.Extend(st.SectorTable().RecoveringSectors())
	inactiveSectorSet := inactiveProvingSet.Extend(st.SectorTable().FailingSectors())

	for _, sectorNo := range inactiveSectorSet.SectorsOn() {
		utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
		inactivePower += utilizationInfo.CurrUtilization()
	}

	return inactivePower
}

// move Sector from Active/Failing
// into Cleared State which means deleting the Sector from state
// remove SectorNumber from all states on chain
// update SectorTable
func (st *StorageMinerActorState_I) _updateClearSector(rt Runtime, sectorNo sector.SectorNumber) {
	sectorState := st.Sectors()[sectorNo].State()
	switch sectorState.StateNumber {
	case SectorActiveSN:
		// expiration case
		st.SectorTable().Impl().ActiveSectors_.Remove(sectorNo)
	case SectorFailingSN:
		// expiration and termination cases
		st.SectorTable().Impl().FailingSectors_.Remove(sectorNo)
	default:
		// Committed and Recovering should not go to Cleared directly
		rt.Abort("invalid state in clearSector")
	}

	delete(st.Sectors(), sectorNo)
	delete(st.SectorUtilization(), sectorNo)
	st.ProvingSet_.Remove(sectorNo)
	st.SectorExpirationQueue().Remove(sectorNo)

	// Send message to SMA
}

// move Sector from Committed/Recovering into Active State
// reset FaultCount to zero
// update SectorTable
func (st *StorageMinerActorState_I) _updateActivateSector(rt Runtime, sectorNo sector.SectorNumber) {
	sectorState := st.Sectors()[sectorNo].State()
	switch sectorState.StateNumber {
	case SectorCommittedSN:
		st.SectorTable().Impl().CommittedSectors_.Remove(sectorNo)
	case SectorRecoveringSN:
		st.SectorTable().Impl().RecoveringSectors_.Remove(sectorNo)
	default:
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("sm._updateActivateSector: invalid state in activateSector")
	}

	st.Sectors()[sectorNo].Impl().State_ = SectorActive()
	st.SectorTable().Impl().ActiveSectors_.Add(sectorNo)
}

// failSector moves Sector from Active/Committed/Recovering into Failing State
// and increments FaultCount if asked to do so (DeclareFaults does not increment faultCount)
// move Sector from Failing to Cleared State if increment results in faultCount exceeds MaxFaultCount
// update SectorTable
// remove from ProvingSet
func (st *StorageMinerActorState_I) _updateFailSector(rt Runtime, sectorNo sector.SectorNumber, increment bool) {
	newFaultCount := st.Sectors()[sectorNo].State().FaultCount

	if increment {
		newFaultCount += 1
	}

	state := st.Sectors()[sectorNo].State()
	switch state.StateNumber {
	case SectorActiveSN:
		// wont be terminated from Active
		st.SectorTable().Impl().ActiveSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorCommittedSN:
		st.SectorTable().Impl().CommittedSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorRecoveringSN:
		st.SectorTable().Impl().RecoveringSectors_.Remove(sectorNo)
		st.SectorTable().Impl().FailingSectors_.Add(sectorNo)
		st.ProvingSet_.Remove(sectorNo)
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	case SectorFailingSN:
		// no change to SectorTable but increase in FaultCount
		st.Sectors()[sectorNo].Impl().State_ = SectorFailing(newFaultCount)
	default:
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("Invalid sector state in CronAction")
	}

	if newFaultCount > MAX_CONSECUTIVE_FAULTS {
		// TODO: heavy penalization: slash pledge collateral and delete sector
		// TODO: SendMessage(SPA.SlashPledgeCollateral)

		st._updateClearSector(rt, sectorNo)
		st.SectorTable().Impl().TerminatedFaults_.Add(sectorNo)
	}
}

// Decision is to currently account for power based on sector
// with at least one active deals and deals cannot be updated
// an alternative proposal is to account for power based on active deals
// an improvement proposal is to allow storage deal update in a sector

// TODO: decide whether declared faults sectors should be
// penalized in the same way as undeclared sectors and how

// this method is called by both SubmitElectionPoSt and SubmitSurprisePoSt
// - Process ProvingSet.SectorsOn()
//   - State Transitions
//     - Committed -> Active and credit power
//     - Recovering -> Active and credit power
//   - Process Active Sectors (pay miners)
// - Process ProvingSet.SectorsOff()
//     - increment FaultCount
//     - clear Sector and slash pledge collateral if count > MAX_CONSECUTIVE_FAULTS
// - Process Expired Sectors (settle deals and return storage collateral to miners)
//     - State Transition
//       - Failing / Recovering / Active / Committed -> Cleared
//     - Remove SectorNumber from Sectors, ProvingSet
// - Update ChallengeEndEpoch
func (a *StorageMinerActorCode_I) _onSuccessfulPoSt(rt Runtime, postSubmission poster.PoStSubmission) InvocOutput {
	h, st := a.State(rt)

	// The proof is verified, process ProvingSet.SectorsOn():
	// ProvingSet.SectorsOn() contains SectorCommitted, SectorActive, SectorRecovering
	// ProvingSet itself does not store states, states are all stored in Sectors.State
	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOn() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			// TODO: determine proper error here and error-handling machinery
			rt.Abort("Sector state not found in map")
		}
		switch sectorState.State().StateNumber {
		case SectorCommittedSN, SectorRecoveringSN:
			st._updateActivateSector(rt, sectorNo)
		case SectorActiveSN:
			// do nothing
			// deal payment is lazily evaluated when deals are expired during _updateSectorUtilization
			// or when miner calls CreditSectorDealPayment
		default:
			// TODO: determine proper error here and error-handling machinery
			rt.Abort("Invalid sector state in ProvingSet.SectorsOn()")
		}
	}

	// commit state change so that committed and recovering are now active

	// Process ProvingSet.SectorsOff()
	// ProvingSet.SectorsOff() contains SectorFailing
	// SectorRecovering is Proving and hence will not be in GetZeros()
	// heavy penalty if Failing for more than or equal to MAX_CONSECUTIVE_FAULTS
	// otherwise increment FaultCount in Sectors().State
	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOff() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			continue
		}
		switch sectorState.State().StateNumber {
		case SectorFailingSN:
			st._updateFailSector(rt, sectorNo, true)
		default:
			// TODO: determine proper error here and error-handling machinery
			rt.Abort("Invalid sector state in ProvingSet.SectorsOff")
		}
	}

	// Process Expiration.
	st._updateExpireSectors(rt)

	UpdateRelease(rt, h, st)

	h, st = a.State(rt)
	newTerminatedFaults := st.SectorTable().TerminatedFaults()
	Release(rt, h, st)

	a._submitFaultReport(
		rt,
		sector.CompactSectorSet(make([]byte, 0)), // NewDeclaredFaults
		sector.CompactSectorSet(make([]byte, 0)), // NewDetectedFaults
		newTerminatedFaults,
	)

	a._submitPowerReport(rt)

	// TODO: check EnsurePledgeCollateralSatisfied
	// pledgeCollateralSatisfied

	h, st = a.State(rt)
	st.ChallengeStatus().Impl().OnChallengeResponse(rt.CurrEpoch())
	st._processStagedCommittedSectors(rt)
	UpdateRelease(rt, h, st)

	return rt.SuccessReturn()

}

func (a *StorageMinerActorCode_I) SubmitSurprisePoSt(rt Runtime, postSubmission poster.PoStSubmission) InvocOutput {
	TODO() // TODO: validate caller

	if !a._isChallenged(rt) {
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("cannot SubmitSurprisePoSt when not challenged")
	}

	// Verify correct PoSt Submission
	isPoStVerified := a._verifyPoStSubmission(rt, postSubmission)
	if !isPoStVerified {
		// no state transition, just error out and miner should submitSurprisePoSt again
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("TODO")
	}

	return a._onSuccessfulPoSt(rt, postSubmission)

}

// Called by StoragePowerConsensus subsystem after verifying the Election proof
// and verifying the PoSt proof within the block. (this happens outside the VM)
// Assume ElectionPoSt has already been successfully verified when the function gets called.
func (a *StorageMinerActorCode_I) SubmitElectionPoSt(rt Runtime, postSubmission poster.PoStSubmission) InvocOutput {

	// TODO: validate caller
	// the caller MUST be the miner who won the block (who won the block should be callable as a a VM runtime call)
	// we also need to enforce that this call happens only once per block, OR make it not callable by special privileged messages
	TODO()

	if !a._isChallenged(rt) {
		rt.Abort("cannot SubmitElectionPoSt when not in election period")
	}

	// we do not need to verify post submission here, as this should have already been done
	// outside of the VM, in StoragePowerConsensus Subsystem. Doing so again would waste
	// significant resources, as proofs are expensive to verify.
	//
	// notneeded := a._verifyPoStSubmission(rt, postSubmission)

	return a._onSuccessfulPoSt(rt, postSubmission)

}

func (st *StorageMinerActorState_I) _updateExpireSectors(rt Runtime) {
	currEpoch := rt.CurrEpoch()

	queue := st.SectorExpirationQueue()
	for queue.Peek().Expiration() <= currEpoch {
		expiredSectorNo := queue.Pop().SectorNumber()

		state := st.Sectors()[expiredSectorNo].State()

		switch state.StateNumber {
		case SectorActiveSN:
			// Note: in order to verify if something was stored in the past, one must
			// scan the chain. SectorNumber can be re-used.

			// do nothing about deal payment
			// it will be evaluated after _updateSectorUtilization
			st._updateClearSector(rt, expiredSectorNo)
		case SectorFailingSN:
			// TODO: check if there is any fault that we should handle here
			// If a SectorFailing Expires, return remaining StorageDealCollateral and remove sector
			// SendMessage(sma.SettleExpiredDeals(sc.DealIDs()))

			// a failing sector expires, no change to FaultCount
			st._updateClearSector(rt, expiredSectorNo)
		default:
			// Note: SectorCommittedSN, SectorRecoveringSN transition first to SectorFailingSN, then expire
			rt.Abort("Invalid sector state in SectorExpirationQueue")
		}
	}

	// Return PledgeCollateral for active expirations
	// SendMessage(spa.Depledge) // TODO
	rt.Abort("TODO: refactor use of this method in order for caller to send this message")
}

// RecoverFaults checks if miners have sufficent collateral
// and adds SectorFailing into SectorRecovering
// - State Transition
//   - Failing -> Recovering with the same FaultCount
// - Add SectorNumber to ProvingSet
// Note that power is not updated until it is active
func (a *StorageMinerActorCode_I) RecoverFaults(rt Runtime, recoveringSet sector.CompactSectorSet) InvocOutput {
	TODO() // TODO: validate caller

	// but miner can RecoverFaults in recovery before challenge
	if a._isChallenged(rt) {
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("cannot RecoverFaults when sm isChallenged")
	}

	h, st := a.State(rt)

	// for all SectorNumber marked as recovering by recoveringSet
	for _, sectorNo := range recoveringSet.SectorsOn() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			// TODO: determine proper error here and error-handling machinery
			rt.Abort("Sector state not found in map")
		}
		switch sectorState.State().StateNumber {
		case SectorFailingSN:
			// Check if miners have sufficient balances in sma

			// SendMessage(sma.PublishStorageDeals) or sma.ResumeStorageDeals?
			// throw if miner cannot cover StorageDealCollateral

			// Check if miners have sufficient pledgeCollateral

			// copy over the same FaultCount
			st.Sectors()[sectorNo].Impl().State_ = SectorRecovering(sectorState.State().FaultCount)
			st.Impl().ProvingSet_.Add(sectorNo)

			st.SectorTable().Impl().FailingSectors_.Remove(sectorNo)
			st.SectorTable().Impl().RecoveringSectors_.Add(sectorNo)

		default:
			// TODO: determine proper error here and error-handling machinery
			// TODO: consider this a no-op (as opposed to a failure), because this is a user
			// call that may be delayed by the chain beyond some other state transition.
			rt.Abort("Invalid sector state in RecoverFaults")
		}
	}

	UpdateRelease(rt, h, st)

	return rt.SuccessReturn()
}

// DeclareFaults penalizes miners (slashStorageDealCollateral and remove power)
// TODO: decide how much storage collateral to slash
// - State Transition
//   - Active / Commited / Recovering -> Failing
// - Update State in Sectors()
// - Remove Active / Commited / Recovering from ProvingSet
func (a *StorageMinerActorCode_I) DeclareFaults(rt Runtime, faultSet sector.CompactSectorSet) InvocOutput {
	TODO() // TODO: validate caller

	if a._isChallenged(rt) {
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("cannot DeclareFaults when challenged")
	}

	h, st := a.State(rt)

	// fail all SectorNumber marked as Failing by faultSet
	for _, sectorNo := range faultSet.SectorsOn() {
		st._updateFailSector(rt, sectorNo, false)
	}

	UpdateRelease(rt, h, st)

	a._submitFaultReport(
		rt,
		faultSet,                                 // DeclaredFaults
		sector.CompactSectorSet(make([]byte, 0)), // DetectedFaults
		sector.CompactSectorSet(make([]byte, 0)), // TerminatedFault
	)

	a._submitPowerReport(rt)

	return rt.SuccessReturn()
}

func (a *StorageMinerActorCode_I) TallySectorDealPayment(rt Runtime, sectorNo sector.SectorNumber) {
	TODO() // verify caller

	h, st := a.State(rt)

	utilizationInfo := st._getUtilizationInfo(rt, sectorNo)
	sectorIsActive := st.SectorTable().Impl().ActiveSectors_.Contain(sectorNo)

	if !sectorIsActive {
		rt.Abort("sm.GetSectorPaymentInfo: sector not active")
	}

	activeDealSet := utilizationInfo.DealExpirationQueue().ActiveDealIDs()
	batchDealPaymentInfo := &deal.BatchDealPaymentInfo_I{
		DealIDs_:               activeDealSet.DealsOn(),
		Action_:                deal.TallyStorageDeals,
		LastChallengeEndEpoch_: st.ChallengeStatus().LastChallengeEndEpoch(),
	}

	Release(rt, h, st)

	panic(batchDealPaymentInfo)
	// Send(StorageMarketActor, ProcessSectorDealPayment(batchDealPaymentInfo))

}

func (a *StorageMinerActorCode_I) _isSealVerificationCorrect(rt Runtime, onChainInfo sector.OnChainSealVerifyInfo) bool {
	h, st := a.State(rt)
	info := st.Info()
	sectorSize := info.SectorSize()
	dealIDs := onChainInfo.DealIDs()
	params := make([]actor.MethodParam, 1+len(dealIDs))

	Release(rt, h, st) // if no modifications made; or

	// TODO: serialize method param as {sectorSize,  DealIDs...}.
	receipt := rt.SendCatchingErrors(&msg.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: storage_market.MethodGetUnsealedCIDForDealIDs,
		Params_: params,
	})

	if receipt.ExitCode() == exitcode.InvalidSectorPacking {
		return false
	}

	ret := receipt.ReturnValue()

	pieceInfos := sector.PieceInfosFromBytes(ret)

	// Unless we enforce a minimum padding amount, this totalPieceSize calculation can be removed.
	// Leaving for now until that decision is entirely finalized.
	var totalPieceSize util.UInt
	for _, pieceInfo := range pieceInfos {
		pieceSize := (*pieceInfo).Size()
		totalPieceSize += pieceSize
	}

	unsealedCID, _ := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos)

	sealCfg := sector.SealCfg_I{
		SectorSize_:     sectorSize,
		SubsectorCount_: info.SubsectorCount(),
		Partitions_:     info.Partitions(),
	}
	svInfo := sector.SealVerifyInfo_I{
		SectorID_: &sector.SectorID_I{
			MinerID_: info.Worker(), // TODO: This is actually miner address. MinerID needs to be derived.
			Number_:  onChainInfo.SectorNumber(),
		},
		OnChain_: onChainInfo,

		// TODO: Make SealCfg sector.SealCfg from miner configuration (where is that?)
		SealCfg_: &sealCfg,

		Randomness_:            sector.SealRandomness(rt.Randomness(onChainInfo.SealEpoch(), 0)),
		InteractiveRandomness_: sector.InteractiveSealRandomness(rt.Randomness(onChainInfo.InteractiveEpoch(), 0)),
		UnsealedCID_:           unsealedCID,
	}

	sdr := filproofs.SDRParams(&filproofs.SDRCfg_I{SealCfg_: &sealCfg})
	return sdr.VerifySeal(&svInfo)
}

func (st *StorageMinerActorState_I) _assertSectorDidNotExist(rt Runtime, sectorNo sector.SectorNumber) {
	_, found := st.Sectors()[sectorNo]
	if found {
		rt.Abort("sm._assertSectorDidNotExist: sector already exists.")
	}
}

func (st *StorageMinerActorState_I) _getUtilizationInfo(rt Runtime, sectorNo sector.SectorNumber) sector.SectorUtilizationInfo {
	utilizationInfo, found := st.SectorUtilization()[sectorNo]

	if !found {
		rt.Abort("sm._getUtilizationInfo: utilization info not found.")
	}

	return utilizationInfo
}

func (st *StorageMinerActorState_I) _initializeUtilizationInfo(rt Runtime, deals []deal.StorageDeal) sector.SectorUtilizationInfo {

	var dealExpirationQueue deal.DealExpirationQueue
	var maxUtilization block.StoragePower
	var lastExpiration block.ChainEpoch

	for _, d := range deals {

		dealExpiration := d.Proposal().EndEpoch()

		if dealExpiration > lastExpiration {
			lastExpiration = dealExpiration
		}

		// TODO: verify what counts towards power here
		// There is PayloadSize, OverheadSize, and Total, see piece.id
		dealPayloadPower := block.StoragePower(d.Proposal().PieceSize().PayloadSize())

		queueItem := &deal.DealExpirationQueueItem_I{
			DealID_:       d.ID(),
			PayloadPower_: dealPayloadPower,
			Expiration_:   dealExpiration,
		}
		dealExpirationQueue.Add(queueItem)
		maxUtilization += dealPayloadPower

	}

	initialUtilizationInfo := &sector.SectorUtilizationInfo_I{
		DealExpirationQueue_: dealExpirationQueue,
		MaxUtilization_:      maxUtilization,
		CurrUtilization_:     maxUtilization,
	}

	return initialUtilizationInfo

}

// Deals must be posted on chain via sma.PublishStorageDeals before PreCommitSector
// TODO(optimization): PreCommitSector could contain a list of deals that are not published yet.
func (a *StorageMinerActorCode_I) PreCommitSector(rt Runtime, info sector.SectorPreCommitInfo) InvocOutput {
	TODO() // TODO: validate caller

	// can be called regardless of Challenged status

	// TODO: might take collateral in case no ProveCommit follows within sometime
	// TODO: collateral also penalizes repeated precommit to get randomness that one likes
	// TODO: might be a good place for Treasury

	h, st := a.State(rt)

	_, found := st.PreCommittedSectors()[info.SectorNumber()]

	if found {
		// TODO: burn some funds?
		rt.Abort("Sector already pre committed.")
	}

	st._assertSectorDidNotExist(rt, info.SectorNumber())

	Release(rt, h, st)

	// verify every DealID has been published and will not expire
	// before the MAX_PROVE_COMMIT_SECTOR_EPOCH + CurrEpoch
	// abort otherwise
	// Send(SMA.VerifyPublishedDealIDs(info.DealIDs()))

	h, st = a.State(rt)

	precommittedSector := &PreCommittedSector_I{
		Info_:          info,
		ReceivedEpoch_: rt.CurrEpoch(),
	}
	st.PreCommittedSectors()[info.SectorNumber()] = precommittedSector

	UpdateRelease(rt, h, st)
	return rt.SuccessReturn()
}

func (a *StorageMinerActorCode_I) ProveCommitSector(rt Runtime, info sector.SectorProveCommitInfo) InvocOutput {
	TODO() // TODO: validate caller

	h, st := a.State(rt)

	preCommitSector, precommitFound := st.PreCommittedSectors()[info.SectorNumber()]

	if !precommitFound {
		rt.Abort("Sector not pre committed.")
	}

	st._assertSectorDidNotExist(rt, info.SectorNumber())

	// check if ProveCommitSector comes too late after PreCommitSector
	elapsedEpoch := rt.CurrEpoch() - preCommitSector.ReceivedEpoch()

	// if more than MAX_PROVE_COMMIT_SECTOR_EPOCH has elapsed
	if elapsedEpoch > sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
		// TODO: potentially some slashing if ProveCommitSector comes late

		// TODO: remove dealIDs from PublishedDeals

		// expired
		delete(st.PreCommittedSectors(), preCommitSector.Info().SectorNumber())
		UpdateRelease(rt, h, st)
		return rt.ErrorReturn(exitcode.UserDefinedError(0)) // TODO: user dfined error code?
	}

	onChainInfo := &sector.OnChainSealVerifyInfo_I{
		SealedCID_:        preCommitSector.Info().SealedCID(),
		SealEpoch_:        preCommitSector.Info().SealEpoch(),
		InteractiveEpoch_: info.InteractiveEpoch(),
		Proof_:            info.Proof(),
		DealIDs_:          preCommitSector.Info().DealIDs(),
		SectorNumber_:     preCommitSector.Info().SectorNumber(),
	}

	isSealVerificationCorrect := st._isSealVerificationCorrect(rt, onChainInfo)
	if !isSealVerificationCorrect {
		// TODO: determine proper error here and error-handling machinery
		rt.Abort("Seal verification failed")
	}

	// TODO: check EnsurePledgeCollateralSatisfied
	// pledgeCollateralSatisfied

	_, utilizationFound := st.SectorUtilization()[onChainInfo.SectorNumber()]
	if utilizationFound {
		rt.Abort("sm.ProveCommitSector: sector number found in SectorUtilization")
	}

	// ActivateStorageDeals
	// Ok if deal has started
	// Send(SMA.ActivateDeals(onChainInfo.DealIDs())
	// abort if activation failed

	// TODO: proper onchain transaction
	// deals := SendMessage(sma, GetDeals(onChainInfo.DealIDs()))
	var deals []deal.StorageDeal
	initialUtilization := st._initializeUtilizationInfo(rt, deals)
	lastDealExpiration := initialUtilization.DealExpirationQueue().LastDealExpiration()

	// add sector expiration to SectorExpirationQueue
	st.SectorExpirationQueue().Add(&SectorExpirationQueueItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   lastDealExpiration,
	})

	// no need to store the proof and randomseed in the state tree
	// verify and drop, only SealCommitment{CommR, DealIDs} on chain
	sealCommitment := &sector.SealCommitment_I{
		SealedCID_:  onChainInfo.SealedCID(),
		DealIDs_:    onChainInfo.DealIDs(),
		Expiration_: lastDealExpiration, // TODO decide if we need this too
	}

	// add SectorNumber and SealCommitment to Sectors
	// set Sectors.State to SectorCommitted
	// Note that SectorNumber will only become Active at the next successful PoSt
	sealOnChainInfo := &SectorOnChainInfo_I{
		SealCommitment_: sealCommitment,
		State_:          SectorCommitted(),
	}

	if st._isChallenged(rt) {
		// move PreCommittedSector to StagedCommittedSectors if in Challenged status
		stagedSectorInfo := &StagedCommittedSectorInfo_I{
			Sector_:      sealOnChainInfo,
			Utilization_: initialUtilization,
		}

		st.StagedCommittedSectors()[onChainInfo.SectorNumber()] = stagedSectorInfo
	} else {
		// move PreCommittedSector to CommittedSectors if not in Challenged status
		st.Sectors()[onChainInfo.SectorNumber()] = sealOnChainInfo
		st.Impl().ProvingSet_.Add(onChainInfo.SectorNumber())
		st.SectorTable().Impl().CommittedSectors_.Add(onChainInfo.SectorNumber())
		st.SectorUtilization()[onChainInfo.SectorNumber()] = initialUtilization
	}

	// now remove SectorNumber from PreCommittedSectors (processed)
	delete(st.PreCommittedSectors(), preCommitSector.Info().SectorNumber())
	UpdateRelease(rt, h, st)

	return rt.SuccessReturn()
}

func getSectorNums(m map[sector.SectorNumber]SectorOnChainInfo) []sector.SectorNumber {
	var l []sector.SectorNumber
	for i, _ := range m {
		l = append(l, i)
	}
	return l
}
