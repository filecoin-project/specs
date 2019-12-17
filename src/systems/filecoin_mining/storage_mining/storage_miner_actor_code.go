package storage_mining

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	market "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	sys "github.com/filecoin-project/specs/systems/filecoin_vm/sysactors"
	util "github.com/filecoin-project/specs/util"
)

const (
	Method_StorageMinerActor_GetOwnerKey = actor.MethodPlaceholder + iota
	Method_StorageMinerActor_GetWorkerKey
	Method_StorageMinerActor_ProcessVerifiedSurprisePoSt
	Method_StorageMinerActor_ProcessVerifiedElectionPoSt
	Method_StorageMinerActor_NotifyOfSurprisePoStChallenge
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type State = StorageMinerActorState
type Any = util.Any
type Bool = util.Bool
type Bytes = util.Bytes
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

var TODO = util.TODO

func (a *StorageMinerActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
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
func (st *StorageMinerActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) State {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (a *StorageMinerActorCode_I) GetOwnerKey(rt Runtime) filcrypto.VRFPublicKey {
	h, st := a.State(rt)
	ret := st.Info().OwnerKey()
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) GetWorkerKey(rt Runtime) filcrypto.VRFPublicKey {
	h, st := a.State(rt)
	ret := st.Info().WorkerKey()
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _isChallenged(rt Runtime) bool {
	h, st := a.State(rt)
	ret := st._isChallenged()
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _challengeHasExpired(rt Runtime) bool {
	h, st := a.State(rt)
	ret := st._challengeHasExpired(rt.CurrEpoch())
	Release(rt, h, st)
	return ret
}

func (a *StorageMinerActorCode_I) _shouldChallenge(rt Runtime) bool {
	h, st := a.State(rt)
	ret := st._shouldChallenge(rt.CurrEpoch())
	Release(rt, h, st)
	return ret
}

////////////////////////////////////////////////////////////////////////////////
// Surprise PoSt
////////////////////////////////////////////////////////////////////////////////

// called by StoragePowerActor to notify StorageMiner of PoSt Challenge (triggered by Cron)
func (a *StorageMinerActorCode_I) NotifyOfSurprisePoStChallenge(rt Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.StoragePowerActorAddr)

	// check that this is a valid challenge
	if !a._shouldChallenge(rt) {
		return rt.SuccessReturn() // silent return, dont re-challenge
	}

	a._expirePreCommittedSectors(rt)

	h, st := a.State(rt)
	// update challenge start epoch
	st.ChallengeStatus().Impl().OnNewChallenge(rt.CurrEpoch())
	UpdateRelease(rt, h, st)
	return rt.SuccessReturn()
}

// Called by the cron actor at every tick.
func (a *StorageMinerActorCode_I) OnCronTick(rt Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.CronActorAddr)

	a._checkSurprisePoStSubmissionHappened(rt)

	return rt.SuccessReturn()

}

// If the miner fails to respond to a surprise PoSt,
// cron triggers reporting every sector as failing for the current proving period.
func (a *StorageMinerActorCode_I) _checkSurprisePoStSubmissionHappened(rt Runtime) InvocOutput {

	// we can return if miner has not yet been challenged
	if !a._isChallenged(rt) {
		// Miner gets out of a challenge when submit a successful PoSt
		// or when detected by CronActor. Hence, not being in isChallenged means that we are good here
		return rt.SuccessReturn()
	}

	if a._challengeHasExpired(rt) {
		// garbage collection - need to be called by cron once in a while
		a._expirePreCommittedSectors(rt)

		// oh no -- we missed it. rekt
		a._onMissedSurprisePoSt(rt)

	}

	return rt.SuccessReturn()
}

// called by CheckSurprisePoSt above for miner who missed their post
func (a *StorageMinerActorCode_I) _onMissedSurprisePoSt(rt Runtime) {
	h, st := a.State(rt)

	failingSectorNumbers := getSectorNums(st.Sectors())
	for _, sectorNo := range failingSectorNumbers {
		st._updateFailSector(rt, sectorNo, true)
	}

	UpdateRelease(rt, h, st)

	a._expireSectors(rt)

	h, st = a.State(rt)

	newDetectedFaults := st.SectorTable().FailingSectors()
	newTerminatedFaults := st.SectorTable().TerminatedFaults()

	Release(rt, h, st)

	a._submitPowerReport(rt)

	// Note: NewDetectedFaults is now the sum of all
	// previously active, committed, and recovering sectors minus expired ones
	// and any previously Failing sectors that did not exceed MAX_CONSECUTIVE_FAULTS
	// Note: previously declared faults is now treated as part of detected faults
	a._slashCollateralForStorageFaults(
		rt,
		sector.CompactSectorSet(make([]byte, 0)), // NewDeclaredFaults
		newDetectedFaults,
		newTerminatedFaults,
	)

	// end of challenge
	// now that new power and faults are tracked move pointer of last challenge response up
	h, st = a.State(rt)
	st.ChallengeStatus().Impl().OnPoStFailure(rt.CurrEpoch())
	st._processStagedCommittedSectors(rt)
	UpdateRelease(rt, h, st)
}

func (a *StorageMinerActorCode_I) _expireSectors(rt Runtime) {

	h, st := a.State(rt)
	expiredSectorNos := st._updateExpireSectors(rt)

	// @param dealIDs
	expireSectorParams := make([]util.Serialization, 1)
	for _, sectorNo := range expiredSectorNos {
		dealIDs, ok := st._getSectorDealIDs(sectorNo)
		if !ok {
			panic("")
		}
		for _, dealID := range dealIDs {
			expireSectorParams = append(expireSectorParams, deal.Serialize_DealID(dealID))
		}
	}

	UpdateRelease(rt, h, st)

	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_ProcessDealExpiration,
		Params_: expireSectorParams,
	})
}

// construct PowerReport from SectorTable
func (a *StorageMinerActorCode_I) _submitPowerReport(rt Runtime) {
	h, st := a.State(rt)
	activePower, err := st._getActivePower()
	if err != nil {
		rt.AbortStateMsg(err.Error())
	}
	inactivePower, err := st._getInactivePower()
	if err != nil {
		rt.AbortStateMsg(err.Error())
	}

	// power report in processPowerReportParam
	_ = &spc.PowerReport_I{
		ActivePower_:   activePower,
		InactivePower_: inactivePower,
	}

	// @param powerReport spc.PowerReport
	processPowerReportParam := make([]util.Serialization, 1)

	Release(rt, h, st)

	// this will go through even if miners do not have the right amount of pledge collateral
	// when _submitPowerReport is called in DeclareFaults and _onMissedSurprisePoSt for power slashing
	// however in SubmitSurprisePoSt EnsurePledgeCollateralSatsified will be called
	// to ensure that miners have the required pledge collateral
	// otherwise, post submission will fail
	// Note: there is no power update in RecoverFaults and hence no EnsurePledgeCollatera or _submitPowerReport
	// Note: ElectionPoSt will always go through and some block rewards will go to LockedBalance in StoragePowerActor
	// if the block winning miner is undercollateralized
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StoragePowerActorAddr,
		Method_: spc.Method_StoragePowerActor_ProcessPowerReport,
		Params_: processPowerReportParam,
	})
}

func (a *StorageMinerActorCode_I) ClaimDealPaymentsForSector(rt Runtime, sectorNo sector.SectorNumber) {
	h, st := a.State(rt)

	dealIDs, ok := st._getSectorDealIDs(sectorNo)
	if !ok {
		rt.AbortArgMsg("sm.ClaimDealPaymentsForSector: invalid sector number.")
	}

	lastPoStEpoch := st.ChallengeStatus().Impl().LastPoStSuccessEpoch()

	// @params dealIDs []deal.DealID activeDealIDs
	// @params lastPoSt block.ChainEpoch
	processDealPaymentParam := make([]util.Serialization, 1)
	for _, dealID := range dealIDs {
		processDealPaymentParam = append(processDealPaymentParam, deal.Serialize_DealID(dealID))
	}
	processDealPaymentParam = append(processDealPaymentParam, block.Serialize_ChainEpoch(lastPoStEpoch))
	Release(rt, h, st)

	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_ProcessDealPayment,
		Params_: processDealPaymentParam,
	})
}

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
func (a *StorageMinerActorCode_I) _onSuccessfulPoSt(rt Runtime) InvocOutput {
	h, st := a.State(rt)

	// The proof is verified, process ProvingSet.SectorsOn():
	// ProvingSet.SectorsOn() contains SectorCommitted, SectorActive, SectorRecovering
	// ProvingSet itself does not store states, states are all stored in Sectors.State
	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOn() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			rt.AbortStateMsg("Sector state not found in map")
		}
		switch sectorState.State().StateNumber {
		case SectorCommittedSN, SectorRecoveringSN:
			st._updateActivateSector(rt, sectorNo)
		case SectorActiveSN:
			// do nothing, deal payment is made lazily
		default:
			rt.AbortStateMsg("Invalid sector state in ProvingSet.SectorsOn()")
		}
	}

	// committed and recovering sectors are now active

	// Process ProvingSet.SectorsOff()
	// ProvingSet.SectorsOff() contains SectorFailing
	// SectorRecovering is Proving and hence will not be in SectorsOff()
	for _, sectorNo := range st.Impl().ProvingSet_.SectorsOff() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			continue
		}
		switch sectorState.State().StateNumber {
		case SectorFailingSN:
			// heavy penalty if Failing for more than or equal to MAX_CONSECUTIVE_FAULTS
			// otherwise increment FaultCount in Sectors().State
			st._updateFailSector(rt, sectorNo, true)
		default:
			rt.AbortStateMsg("Invalid sector state in ProvingSet.SectorsOff")
		}
	}

	UpdateRelease(rt, h, st)

	h, st = a.State(rt)
	newTerminatedFaults := st.SectorTable().TerminatedFaults()
	Release(rt, h, st)

	// Process Expiration.
	a._expireSectors(rt)
	a._submitPowerReport(rt)

	a._slashCollateralForStorageFaults(
		rt,
		sector.CompactSectorSet(make([]byte, 0)), // NewDeclaredFaults
		sector.CompactSectorSet(make([]byte, 0)), // NewDetectedFaults
		newTerminatedFaults,
	)

	h, st = a.State(rt)
	st.ChallengeStatus().Impl().OnPoStSuccess(rt.CurrEpoch())
	st._processStagedCommittedSectors(rt)
	UpdateRelease(rt, h, st)

	return rt.SuccessReturn()
}

// Called by StoragePowerConsensus subsystem after verifying the Election proof
// and verifying the PoSt proof in the block header.
// Assume ElectionPoSt has already been successfully verified (both proof and partial ticket
// value) when the function gets called.
// Likewise assume that the rewards have already been granted to the storage miner actor. This only handles sector management.
func (a *StorageMinerActorCode_I) ProcessVerifiedElectionPoSt(rt Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	// The receiver must be the miner who produced the block for which this message is created.
	Assert(rt.ToplevelBlockWinner() == rt.CurrReceiver())

	// we do not need to verify post submission here, as this should have already been done
	// outside of the VM, in StoragePowerConsensus Subsystem. Doing so again would waste
	// significant resources, as proofs are expensive to verify.
	//
	// notneeded := a._verifyPoStSubmission(rt)

	// the following will update last challenge response time
	return a._onSuccessfulPoSt(rt)
}

// called by verifier to update miner state on successful surprise post
// after it has been verified in the storage_mining_subsystem
func (a *StorageMinerActorCode_I) ProcessSurprisePoSt(rt Runtime, onChainInfo sector.OnChainPoStVerifyInfo) InvocOutput {
	TODO() // TODO: validate caller

	if !a._verifySurprisePoSt(rt, onChainInfo) {
		rt.Abort(
			exitcode.UserDefinedError(exitcode.PoStVerificationFailed),
			"PoSt verification failed")
	}

	// Ensure pledge collateral satisfied
	// otherwise, abort ProcessVerifiedSurprisePoSt
	a._ensurePledgeCollateralSatisfied(rt)

	return a._onSuccessfulPoSt(rt)

}

func (a *StorageMinerActorCode_I) _verifySurprisePoSt(rt Runtime, onChainInfo sector.OnChainPoStVerifyInfo) bool {
	h, st := a.State(rt)
	info := st.Info()

	// 1. Check that the miner in question is currently being challenged
	if !st._isChallenged() {
		// TODO: determine proper error here and error-handling machinery
		// rt.Abort("cannot SubmitSurprisePoSt when not challenged")
		rt.AbortStateMsg("Invalid Surprise PoSt. Miner not challenged.")
	}

	// 2. Check that the challenge has not expired
	// Check that miner can still submit (i.e. that the challenge window has not passed)
	// This will prevent miner from submitting a Surprise PoSt past the challenge period
	if st._challengeHasExpired(rt.CurrEpoch()) {
		rt.AbortStateMsg("Invalid Surprise PoSt. Challenge has expired.")
	}

	// 3. Verify the partialTicket values
	if !a._verifySurprisePoStMeetsTargetReq(rt) {
		rt.AbortStateMsg("Invalid Surprise PoSt. Tickets do not meet target.")
	}

	// verify the partialTickets themselves
	// 4. Verify appropriate randomness

	// pull from consts
	SPC_LOOKBACK_POST := uint64(0)
	randomnessEpoch := st.ChallengeStatus().LastChallengeEpoch()
	randomness := rt.Randomness(randomnessEpoch, SPC_LOOKBACK_POST)
	panic(randomness)                       // ignore circular dependency
	var postRandomnessInput util.Randomness // sms.PreparePoStChallengeSeed(randomness, actorAddr)

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	// get worker key from minerAddr
	out := rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     rt.ImmediateCaller(),
		Method_: Method_StorageMinerActor_GetWorkerKey,
	})
	temp, err := filcrypto.Deserialize_VRFPublicKey(out.ReturnValue())
	// ignore below line, redeclaring because of a code generation issue :/
	workerKey := filcrypto.VRFPublicKey(temp)

	if err != nil {
		rt.AbortArgMsg("miner has invalid owner address")
	}

	if !postRand.Verify(postRandomnessInput, workerKey) {
		rt.AbortStateMsg("Invalid Surprise PoSt. Invalid randomness.")
	}

	UpdateRelease(rt, h, st)

	// 5. Get public inputs
	sectorSize := info.SectorSize()

	postCfg := sector.PoStCfg_I{
		Type_:        sector.PoStType_SurprisePoSt,
		SectorSize_:  sectorSize,
		WindowCount_: info.WindowCount(),
		Partitions_:  info.SurprisePoStPartitions(),
	}

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		PoStCfg_:    &postCfg,
		Randomness_: onChainInfo.Randomness(),
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{SurprisePoStCfg_: &postCfg})

	// 6. Verify the PoSt Proof
	return sdr.VerifySurprisePoSt(&pvInfo)
}

// todo: define target
func (a *StorageMinerActorCode_I) _verifySurprisePoStMeetsTargetReq(rt Runtime) bool {
	util.TODO()
	return false
}

////////////////////////////////////////////////////////////////////////////////
// Faults
////////////////////////////////////////////////////////////////////////////////

// RecoverFaults checks if miners have sufficent collateral
// and adds SectorFailing into SectorRecovering
// - State Transition
//   - Failing -> Recovering with the same FaultCount
// - Add SectorNumber to ProvingSet
// Note that power is not updated until it is active
func (a *StorageMinerActorCode_I) RecoverFaults(rt Runtime, recoveringSet sector.CompactSectorSet) InvocOutput {
	TODO() // TODO: validate caller

	// RecoverFaults is only called when miners are not challenged
	if a._isChallenged(rt) {
		rt.AbortStateMsg("cannot RecoverFaults when sm isChallenged")
	}

	h, st := a.State(rt)

	// for all SectorNumber marked as recovering by recoveringSet
	for _, sectorNo := range recoveringSet.SectorsOn() {
		sectorState, found := st.Sectors()[sectorNo]
		if !found {
			rt.AbortStateMsg("Sector state not found in map")
		}
		switch sectorState.State().StateNumber {
		case SectorFailingSN:
			// pledge collateral is ensured at PoSt submission
			// no need to top up deal collateral because none was slashed during declared/detected faults

			// copy over the same FaultCount
			st.Sectors()[sectorNo].Impl().State_ = SectorRecovering(sectorState.State().FaultCount)
			st.Impl().ProvingSet_.Add(sectorNo)

			st.SectorTable().Impl().FailingSectors_.Remove(sectorNo)
			st.SectorTable().Impl().RecoveringSectors_.Add(sectorNo)

		default:
			// TODO: determine proper error here and error-handling machinery
			// TODO: consider this a no-op (as opposed to a failure), because this is a user
			// call that may be delayed by the chain beyond some other state transition.
			rt.AbortStateMsg("Invalid sector state in RecoverFaults")
		}
	}

	UpdateRelease(rt, h, st)

	return rt.SuccessReturn()
}

// DeclareFaults penalizes miners (slashStorageDealCollateral and remove power)
// - State Transition
//   - Active / Commited / Recovering -> Failing
// - Update State in Sectors()
// - Remove Active / Commited / Recovering from ProvingSet
func (a *StorageMinerActorCode_I) DeclareFaults(rt Runtime, faultSet sector.CompactSectorSet) InvocOutput {
	TODO() // TODO: validate caller

	if a._isChallenged(rt) {
		rt.AbortStateMsg("cannot DeclareFaults when challenged")
	}

	h, st := a.State(rt)

	// fail all SectorNumber marked as Failing by faultSet
	for _, sectorNo := range faultSet.SectorsOn() {
		st._updateFailSector(rt, sectorNo, false)
	}

	UpdateRelease(rt, h, st)

	a._submitPowerReport(rt)

	a._slashCollateralForStorageFaults(
		rt,
		faultSet,                                 // NewDeclaredFaults
		sector.CompactSectorSet(make([]byte, 0)), // NewDetectedFaults
		sector.CompactSectorSet(make([]byte, 0)), // NewTerminatedFault
	)

	return rt.SuccessReturn()
}

func (a *StorageMinerActorCode_I) _slashDealsForStorageFault(rt Runtime, sectorNumbers []sector.SectorNumber, faultType sector.StorageFaultType) {

	h, st := a.State(rt)

	dealIDs := make([]deal.DealID, 0)
	for _, sectorNo := range sectorNumbers {
		sectorDealIDs, ok := st._getSectorDealIDs(sectorNo)
		if !ok {
			panic("")
		}
		dealIDs = append(dealIDs, sectorDealIDs...)
	}

	// @param dealIDs []deal.DealID
	// @param faultType sector.StorageFaultType
	processDealSlashParam := make([]util.Serialization, 2)
	for _, dealID := range dealIDs {
		processDealSlashParam = append(processDealSlashParam, deal.Serialize_DealID(dealID))
	}
	processDealSlashParam = append(processDealSlashParam, sector.Serialize_StorageFaultType(faultType))
	Release(rt, h, st)

	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_ProcessDealSlash,
		Params_: processDealSlashParam,
	})
}

func (a *StorageMinerActorCode_I) _slashPledgeForStorageFault(rt Runtime, sectorNumbers []sector.SectorNumber, faultType sector.StorageFaultType) {
	h, st := a.State(rt)

	affectedPower := block.StoragePower(0)
	for _, sectorNo := range sectorNumbers {
		sectorPower, ok := st._getSectorPower(sectorNo)
		if !ok {
			panic("")
		}
		affectedPower += sectorPower
	}

	// @param affectedPower block.StoragePower
	// @param faultType sector.StorageFaultType
	slashPledgeParams := make([]util.Serialization, 2)
	slashPledgeParams = append(slashPledgeParams, block.Serialize_StoragePower(affectedPower))
	slashPledgeParams = append(slashPledgeParams, sector.Serialize_StorageFaultType(faultType))

	Release(rt, h, st)

	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StoragePowerActorAddr,
		Method_: spc.Method_StoragePowerActor_SlashPledgeForStorageFault,
		Params_: slashPledgeParams,
	})
}

// reset NewTerminatedFaults
func (a *StorageMinerActorCode_I) _slashCollateralForStorageFaults(
	rt Runtime,
	newDeclaredFaults sector.CompactSectorSet, // diff value
	newDetectedFaults sector.CompactSectorSet, // diff value
	newTerminatedFaults sector.CompactSectorSet, // diff value
) {

	// only terminatedFault will result in collateral deal slashing
	if len(newTerminatedFaults) > 0 {
		a._slashDealsForStorageFault(rt, newTerminatedFaults.SectorsOn(), sector.TerminatedFault)
		a._slashPledgeForStorageFault(rt, newTerminatedFaults.SectorsOn(), sector.TerminatedFault)
	}

	if len(newDetectedFaults) > 0 {
		a._slashPledgeForStorageFault(rt, newDetectedFaults.SectorsOn(), sector.DetectedFault)
	}

	if len(newDeclaredFaults) > 0 {
		a._slashPledgeForStorageFault(rt, newDeclaredFaults.SectorsOn(), sector.DeclaredFault)
	}

	// reset terminated faults
	h, st := a.State(rt)
	st.SectorTable().Impl().TerminatedFaults_ = sector.CompactSectorSet(make([]byte, 0))
	UpdateRelease(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Sector Commitment
////////////////////////////////////////////////////////////////////////////////

func (a *StorageMinerActorCode_I) _verifySeal(rt Runtime, onChainInfo sector.OnChainSealVerifyInfo) bool {
	h, st := a.State(rt)
	info := st.Info()
	sectorSize := info.SectorSize()

	// @param sectorSize util.UVarint
	// @param dealIDs []deal.DealID onChainInfo.DealIDs()
	getPieceInfoParams := make([]util.Serialization, 2)
	getPieceInfoParams = append(getPieceInfoParams, sector.Serialize_SectorSize(sectorSize))
	for _, dealID := range onChainInfo.DealIDs() {
		getPieceInfoParams = append(getPieceInfoParams, deal.Serialize_DealID(dealID))
	}

	Release(rt, h, st)

	getPieceInfosReceipt := rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_GetPieceInfosForDealIDs,
		Params_: getPieceInfoParams,
	})

	getPieceInfosRet := getPieceInfosReceipt.ReturnValue()
	pieceInfos := sector.PieceInfosFromBytes(getPieceInfosRet)

	// Unless we enforce a minimum padding amount, this totalPieceSize calculation can be removed.
	// Leaving for now until that decision is entirely finalized.
	var totalPieceSize util.UInt
	for _, pieceInfo := range pieceInfos {
		pieceSize := (*pieceInfo).Size()
		totalPieceSize += pieceSize
	}

	unsealedCID, _ := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos)

	sealCfg := sector.SealCfg_I{
		SectorSize_:  sectorSize,
		WindowCount_: info.WindowCount(),
		Partitions_:  info.SealPartitions(),
	}

	// @param address addr.Address info.Worker()
	getActorIDParams := make([]util.Serialization, 1)
	getActorIDParams = append(getActorIDParams, addr.Serialize_Address(info.Worker()))

	getActorIDReceipt := rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.InitActorAddr,
		Method_: sys.Method_InitActor_GetActorIDForAddress,
		Params_: getActorIDParams,
	})

	getActorIDRet := getActorIDReceipt.ReturnValue()
	minerID, err := addr.Deserialize_Address(getActorIDRet)
	if err != nil {
		rt.AbortStateMsg("sm.verifySeal: failed to get actor id")
	}

	svInfo := sector.SealVerifyInfo_I{
		SectorID_: &sector.SectorID_I{
			MinerID_: minerID,
			Number_:  onChainInfo.SectorNumber(),
		},
		OnChain_: onChainInfo,

		// TODO: Make SealCfg sector.SealCfg from miner configuration (where is that?)
		SealCfg_: &sealCfg,

		Randomness_:            sector.SealRandomness(rt.Randomness(onChainInfo.SealEpoch(), 0)),
		InteractiveRandomness_: sector.InteractiveSealRandomness(rt.Randomness(onChainInfo.InteractiveEpoch(), 0)),
		UnsealedCID_:           unsealedCID,
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{SealCfg_: &sealCfg})
	return sdr.VerifySeal(&svInfo)
}

// Deals must be posted on chain via sma.PublishStorageDeals before PreCommitSector
// Optimization: PreCommitSector could contain a list of deals that are not published yet.
func (a *StorageMinerActorCode_I) PreCommitSector(rt Runtime, info sector.SectorPreCommitInfo) InvocOutput {
	TODO() // TODO: validate caller
	// can be called regardless of Challenged status

	// TODO: might be a good place for Treasury

	h, st := a.State(rt)

	msgValue := rt.ValueReceived()
	depositReq := st._getPreCommitDepositReq(rt)

	if msgValue < depositReq {
		rt.AbortFundsMsg("sm.PreCommitSector: insufficient precommit deposit.")
	}

	_, found := st.PreCommittedSectors()[info.SectorNumber()]

	if found {
		// no burn funds since miners can't do repeated precommit
		rt.AbortStateMsg("Sector already pre committed.")
	}

	st._assertSectorDidNotExist(rt, info.SectorNumber())

	// @param dealIDs []deal.DealID info.DealIDs()
	params := make([]util.Serialization, 1)
	for _, dealID := range info.DealIDs() {
		params = append(params, deal.Serialize_DealID(dealID))
	}

	Release(rt, h, st)

	// verify every DealID has been published and will not expire
	// before the MAX_PROVE_COMMIT_SECTOR_EPOCH + CurrEpoch
	// abort otherwise
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_VerifyPublishedDealIDs,
		Params_: params,
	})

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

	msgSender := rt.ImmediateCaller()

	h, st := a.State(rt)

	preCommitSector, precommitFound := st.PreCommittedSectors()[info.SectorNumber()]

	if !precommitFound {
		rt.AbortStateMsg("Sector not pre committed.")
	}

	st._assertSectorDidNotExist(rt, info.SectorNumber())

	// check if ProveCommitSector comes too late after PreCommitSector
	elapsedEpoch := rt.CurrEpoch() - preCommitSector.ReceivedEpoch()

	// if more than MAX_PROVE_COMMIT_SECTOR_EPOCH has elapsed
	if elapsedEpoch > sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
		// PreCommittedSectors is cleaned up at _expirePreCommitSectors triggered by Cron
		rt.Abort(
			exitcode.UserDefinedError(exitcode.DeadlineExceeded),
			"more than MAX_PROVE_COMMIT_SECTOR_EPOCH has elapsed")
	}

	onChainInfo := &sector.OnChainSealVerifyInfo_I{
		SealedCID_:        preCommitSector.Info().SealedCID(),
		SealEpoch_:        preCommitSector.Info().SealEpoch(),
		InteractiveEpoch_: info.InteractiveEpoch(),
		Proof_:            info.Proof(),
		DealIDs_:          preCommitSector.Info().DealIDs(),
		SectorNumber_:     preCommitSector.Info().SectorNumber(),
	}

	isSealVerified := st._verifySeal(rt, onChainInfo)
	if !isSealVerified {
		rt.Abort(
			exitcode.UserDefinedError(exitcode.SealVerificationFailed),
			"Seal verification failed")
	}

	minerInfo := st.Info()
	sectorSize := minerInfo.SectorSize()
	sectorPower := block.StoragePower(sectorSize)

	// Activate storage deals, abort if activation failed
	// @param expiration block.ChainEpoch info.Expiration()
	// @param sectorSize block.StoragePower sectorSize
	// @param dealIDs []deal.DealID onChainInfo.DealIDs()
	activateDealsParams := make([]util.Serialization, 1)
	activateDealsParams = append(activateDealsParams, block.Serialize_ChainEpoch(info.Expiration()))
	activateDealsParams = append(activateDealsParams, block.Serialize_StoragePower(sectorPower))
	for _, dealID := range onChainInfo.DealIDs() {
		activateDealsParams = append(activateDealsParams, deal.Serialize_DealID(dealID))
	}

	UpdateRelease(rt, h, st)

	activateDealsReceipt := rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StorageMarketActorAddr,
		Method_: market.Method_StorageMarketActor_ActivateDeals,
		Params_: activateDealsParams,
	})

	h, st = a.State(rt)

	activateDealsRet := activateDealsReceipt.ReturnValue()
	sectorActivePower, err := block.Deserialize_StoragePower(activateDealsRet)
	if err != nil {
		rt.AbortStateMsg("Failed to deserialize sector power")
	}

	// add sector expiration to SectorExpirationQueue
	st.SectorExpirationQueue().Add(&SectorExpirationQueueItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   info.Expiration(),
	})

	// no need to store the proof and randomseed in the state tree
	// verify and drop, only SealCommitment{CommR, DealIDs} on chain
	sealCommitment := &sector.SealCommitment_I{
		SealedCID_: onChainInfo.SealedCID(),
		DealIDs_:   onChainInfo.DealIDs(),
	}

	// add SectorNumber and SealCommitment to Sectors
	// set Sectors.State to SectorCommitted
	// Note that SectorNumber will only become Active at the next successful PoSt
	sealOnChainInfo := &SectorOnChainInfo_I{
		SealCommitment_: sealCommitment,
		State_:          SectorCommitted(),
		Power_:          block.StoragePower(sectorActivePower),
		Activation_:     rt.CurrEpoch(),
		Expiration_:     info.Expiration(),
	}

	if a._isChallenged(rt) {
		// move PreCommittedSector to StagedCommittedSectors if in Challenged status
		st.StagedCommittedSectors()[onChainInfo.SectorNumber()] = sealOnChainInfo
	} else {
		// move PreCommittedSector to CommittedSectors if not in Challenged status
		st.Sectors()[onChainInfo.SectorNumber()] = sealOnChainInfo
		st.Impl().ProvingSet_.Add(onChainInfo.SectorNumber())
		st.SectorTable().Impl().CommittedSectors_.Add(onChainInfo.SectorNumber())
	}

	// now remove SectorNumber from PreCommittedSectors (processed)
	delete(st.PreCommittedSectors(), preCommitSector.Info().SectorNumber())
	depositReq := st._getPreCommitDepositReq(rt)
	UpdateRelease(rt, h, st)

	// return deposit requirement to sender
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:    msgSender,
		Value_: depositReq,
	})

	return rt.SuccessReturn()
}

func (a *StorageMinerActorCode_I) _ensurePledgeCollateralSatisfied(rt Runtime) {
	emptyParams := make([]util.Serialization, 0)
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     addr.StoragePowerActorAddr,
		Method_: spc.Method_StoragePowerActor_EnsurePledgeCollateralSatisfied,
		Params_: emptyParams,
	})
}

func (a *StorageMinerActorCode_I) _expirePreCommittedSectors(rt Runtime) {

	h, st := a.State(rt)

	expiredSectorNum := 0
	inactiveDealIDs := make([]deal.DealID, 0)

	for _, preCommitSector := range st.PreCommittedSectors() {

		elapsedEpoch := rt.CurrEpoch() - preCommitSector.ReceivedEpoch()

		if elapsedEpoch > sector.MAX_PROVE_COMMIT_SECTOR_EPOCH {
			delete(st.PreCommittedSectors(), preCommitSector.Info().SectorNumber())
			expiredSectorNum += 1
			inactiveDealIDs = append(inactiveDealIDs, preCommitSector.Info().DealIDs()...)
		}
	}

	depositToBurn := st._getPreCommitDepositReq(rt)
	UpdateRelease(rt, h, st)

	// send funds to BurntFundsActor
	if depositToBurn > 0 {
		rt.SendPropagatingErrors(&vmr.InvocInput_I{
			To_:    addr.BurntFundsActorAddr,
			Value_: depositToBurn,
		})
	}

	// clear inactive deals
	if len(inactiveDealIDs) > 0 {
		// @param []deal.DealID inactiveDealIDs
		dealsToClearParams := make([]util.Serialization, 1)
		for _, dealID := range inactiveDealIDs {
			dealsToClearParams = append(dealsToClearParams, deal.Serialize_DealID(dealID))
		}

		rt.SendPropagatingErrors(&vmr.InvocInput_I{
			To_:     addr.StorageMarketActorAddr,
			Method_: market.Method_StorageMarketActor_ClearInactiveDealIDs,
			Params_: dealsToClearParams,
		})
	}
}

func getSectorNums(m map[sector.SectorNumber]SectorOnChainInfo) []sector.SectorNumber {
	var l []sector.SectorNumber
	for i, _ := range m {
		l = append(l, i)
	}
	return l
}
