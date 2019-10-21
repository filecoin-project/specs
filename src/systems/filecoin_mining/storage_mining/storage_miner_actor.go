package storage_mining

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving/poster"
import proving "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving"
import power "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"

// If a Post is missed (either due to faults being not declared on time or
// because the miner run out of time, every sector is reported as failing
// for the current proving period.
func (sm *StorageMinerActor_I) CheckPoStSubmissionHappened() {
	var challengeEpoch block.ChainEpoch // TODO
	if sm.LastChallengePoSt() == challengeEpoch {
		return // success
	}

	// oh no -- we missed it. rekt
	sm.failSectors(getSectorNums(sm.SectorStates()))
	sm.expireSectors()
}

func (sm *StorageMinerActor_I) verifyPoStSubmission(postSubmission poster.PoStSubmission) bool {
	// 1. A proof must be submitted after the postRandomness for this proving
	// period is on chain
	// if rt.ChainEpoch < sm.ProvingPeriodEnd - challengeTime {
	//   panic("too early")
	// }

	// 2. A proof must be a valid snark proof with the correct public inputs
	// 2.1 Get randomness from the chain at the right epoch
	// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
	// 2.2 Generate the set of challenges
	// challenges := GenerateChallengesForPoSt(r, keys(sm.Sectors))
	// 2.3 Verify the PoSt Proof
	// verifyPoSt(challenges, TODO)

	panic("TODO")
	return true
}

// move Sector from Active/Committed/Recovering/Failing
// into Cleared State which means deleting the Sector from state
// remove SectorNumber from all states on chain
func (sm *StorageMinerActor_I) clearSector(sectorNo sector.SectorNumber) {
	delete(sm.Sectors(), sectorNo)
	delete(sm.SectorStates(), sectorNo)
	sm.ProvingSet_.Remove(sectorNo)
	sm.SectorExpirationQueue().Remove(sectorNo)
}

// move Sector from Committed/Recovering into Active State
// reset FaultCount to zero
func (sm *StorageMinerActor_I) activateSector(sectorNo sector.SectorNumber) {
	sm.SectorStates()[sectorNo] = SectorActive()
}

// failSector moves Sectors from Active/Committed/Recovering into Failing State
// and increments FaultCount
func (sm *StorageMinerActor_I) failSectors(sectors []sector.SectorNumber) {
	// TODO: detemine how much pledge collateral to slash
	// TODO: SendMessage(spa.SlashPledgeCollateral)
	var powerUpdate uint64 = 0

	// this can happen from all states
	for _, sectorNo := range sectors {
		state := sm.SectorStates()[sectorNo]
		switch state.StateNumber {
		case SectorActiveSN:
			// SlashStorageDealCollateral
			// SendMessage(sma.slashStorageDealCollateral(sc.DealIDs()))
			sm.incrementFaultCount(sectorNo)
			powerUpdate = powerUpdate - sm.Info().SectorSize()
		case SectorCommittedSN, SectorRecoveringSN:
			// SlashStorageDealCollateral
			// SendMessage(sma.slashStorageDealCollateral(sc.DealIDs()))
			sm.incrementFaultCount(sectorNo)
		case SectorFailingSN:
			sm.incrementFaultCount(sectorNo)
		default:
			// TODO: proper failure
			panic("Invalid sector state in CronAction")
		}
	}

	// TODO
	// Reset Proving Period and report power updates
	// sm.ProvingPeriodEnd_ = PROVING_PERIOD_TIME
	// SendMessage(sma.UpdatePower(powerUpdate))
}

// increment FaultCount and check if greater than MaxFaultCount
// move from Failing to Cleared State if so
// return true if results in sector terminated
func (sm *StorageMinerActor_I) incrementFaultCount(sectorNo sector.SectorNumber) bool {
	newFaultCount := sm.SectorStates()[sectorNo].FaultCount + 1
	sm.SectorStates()[sectorNo] = SectorFailing(newFaultCount)
	if newFaultCount > MAX_CONSECUTIVE_FAULTS {
		// TODO: heavy penalization: slash pledge collateral and delete sector
		// TODO: SendMessage(SPA.SlashPledgeCollateral)
		sm.clearSector(sectorNo)
		return true
	}
	// increment FaultCount
	// TODO: SendMessage(sma.SlashStorageDealCollateral)
	sm.ProvingSet_.Remove(sectorNo)
	return false
}

// Decision is to currently account for power based on sector
// with at least one active deals and deals cannot be updated
// an alternative proposal is to account for power based on active deals
// an improvement proposal is to allow storage deal update in a sector

// TODO: decide whether declared faults sectors should be
// penalized in the same way as undeclared sectors and how

// SubmitPoSt Workflow:
// - Verify PoSt Submission
// - Process ProvingSet.SectorsOn()
//   - State Transitions
//     - Committed -> Active and credit power
//     - Recovering -> Active and credit power
//   - Process Active Sectors (pay miners)
// - Process ProvingSet.SectorsOff()
//	   - increment FaultCount
//     - clear Sector and slash pledge collateral if count > MAX_CONSECUTIVE_FAULTS
// - Process Expired Sectors (settle deals and return storage collateral to miners)
//     - State Transition
//       - Failing / Recovering / Active / Committed -> Cleared
//     - Remove SectorNumber from Sectors, SectorStates, ProvingSet
func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission) {
	// Verify correct PoSt Submission
	isPoStVerified := sm.verifyPoStSubmission(postSubmission)
	if !isPoStVerified {
		// no state transition, just error out and miner should submitPoSt again
		// TODO: proper failure
		panic("TODO")
	}

	// TODO: make sure sm.LastChallengePost gets updated properly
	var challengeEpoch block.ChainEpoch
	sm.LastChallengePoSt_ = challengeEpoch

	// The proof is verified, process ProvingSet.SectorsOn():
	// ProvingSet.SectorsOn() contains SectorCommitted, SectorActive, SectorRecovering
	// ProvingSet itself does not store states, states are all stored in SectorStates
	var powerUpdate uint64 = 0
	for _, sectorNo := range sm.ProvingSet_.SectorsOn() {
		state := sm.SectorStates()[sectorNo]
		switch state.StateNumber {
		case SectorCommittedSN:
			sm.activateSector(sectorNo)
			powerUpdate = powerUpdate + sm.Info().SectorSize()
		case SectorRecoveringSN:
			// Note: SectorState.FaultCount is also reset to zero here
			sm.activateSector(sectorNo)
			powerUpdate = powerUpdate + sm.Info().SectorSize()
		case SectorActiveSN:
			// Process payment in all active deals
			// Note: this must happen before marking sectors as expired.
			// TODO: Pay miner in a single batch message
			// SendMessage(sma.ProcessStorageDealsPayment(sm.Sectors()[sectorNumber].DealIDs()))
		default:
			// TODO: proper failure
			panic("Invalid sector state in ProvingSet.SectorsOn()")
		}
	}

	// all checks pass, everything in proving set is now active
	var powerReport power.PowerReport
	activeSectorCount := uint64(len(sm.ProvingSet_.SectorsOn()))

	// Process ProvingSet.SectorsOff()
	// ProvingSet.SectorsOff() contains SectorFailing
	// SectorRecovering is Proving and hence will not be in GetZeros()
	// heavy penalty if Failing for more than or equal to MAX_CONSECUTIVE_FAULTS
	// otherwise increment FaultCount in SectorStates()
	inactivePower := uint64(len(sm.ProvingSet_.SectorsOff())) * sm.Info().SectorSize()
	powerReport.Impl().InactivePower_ = block.StoragePower(inactivePower)

	var terminatedFaultCount util.UVarint = 0
	var declaredFaultCount util.UVarint = 0
	for _, sectorNo := range sm.ProvingSet_.SectorsOff() {
		state := sm.SectorStates()[sectorNo]
		switch state.StateNumber {
		case SectorFailingSN:
			isTerminated := sm.incrementFaultCount(sectorNo)
			if isTerminated {
				terminatedFaultCount = terminatedFaultCount + 1
			} else {
				declaredFaultCount = declaredFaultCount + 1
			}
		default:
			// TODO: proper failure
			panic("Invalid sector state in ProvingSet.SectorsOff")
		}
	}

	// Process Expiration.
	// note: this may do an additional power update
	expireRes := sm.expireSectors()
	powerReport.Impl().ActivePower_ = block.StoragePower(activeSectorCount*sm.Info().SectorSize()) - expireRes

	powerReport.Impl().SlashDeclaredFaults_ = declaredFaultCount
	powerReport.Impl().SlashTerminatedFaults_ = terminatedFaultCount
	powerReport.Impl().SlashDetectedFaults_ = util.UVarint(0)

	// TODO: Send(spa.ProcessPowerReport(powerReport))

	// Reset Proving Period and report power updates
	// sm.ProvingPeriodEnd_ = PROVING_PERIOD_TIME

}

type ExpireSectorsResult struct {
	PowerLoss block.StoragePower
}

func (sm *StorageMinerActor_I) expireSectors() ExpireSectorsResult {
	var currEpoch block.ChainEpoch // TODO: replace this with runtime.CurrentEpoch()
	var powerLoss uint64 = 0

	queue := sm.SectorExpirationQueue()
	for queue.Peek().Expiration() <= currEpoch {
		expiredSectorNo := queue.Pop().SectorNumber()

		state := sm.SectorStates()[expiredSectorNo]
		// sc := sm.Sectors()[expiredSectorNo]
		switch state.StateNumber {
		case SectorActiveSN:
			// Note: in order to verify if something was stored in the past, one must
			// scan the chain. SectorNumber can be re-used.

			// Settle deals
			// SendMessage(sma.SettleExpiredDeals(sc.DealIDs()))
			sm.clearSector(expiredSectorNo)
			powerLoss = powerLoss + sm.Info().SectorSize()
		case SectorFailingSN:
			// TODO: check if there is any fault that we should handle here
			// If a SectorFailing Expires, return remaining StorageDealCollateral and remove sector
			// SendMessage(sma.SettleExpiredDeals(sc.DealIDs()))
			sm.clearSector(expiredSectorNo)
		default:
			// Note: SectorCommittedSN, SectorRecoveringSN transition first to SectorFailingSN, then expire
			// TODO: proper failure
			panic("Invalid sector state in SectorExpirationQueue")
		}
	}

	expireRes := ExpireSectorsResult{
		PowerLoss: block.StoragePower(powerLoss),
	}

	return expireRes

	// Update Power
	// SendMessage(sma.UpdatePower(powerUpdate)) // TODO

	// Return PledgeCollateral for active expirations
	// SendMessage(spa.Depledge) // TODO
}

// RecoverFaults checks if miners have sufficent collateral
// and adds SectorFailing into SectorRecovering
// - State Transition
//   - Failing -> Recovering with the same FaultCount
// - Add SectorNumber to ProvingSet
// Note that power is not updated until it is active
func (sm *StorageMinerActor_I) RecoverFaults(recoveringSet sector.CompactSectorSet) {
	// for all SectorNumber marked as recovering by recoveringSet
	for _, sectorNo := range recoveringSet.SectorsOn() {
		state := sm.SectorStates()[sectorNo]
		switch state.StateNumber {
		case SectorFailingSN:
			// Check if miners have sufficient balances in sma
			// SendMessage(sma.PublishStorageDeals) or sma.ResumeStorageDeals?
			// throw if miner cannot cover StorageDealCollateral

			// copy over the same FaultCount
			sm.SectorStates()[sectorNo] = SectorRecovering(state.FaultCount)
			sm.ProvingSet_.Add(sectorNo)
		default:
			// TODO: proper failure
			// TODO: consider this a no-op (as opposed to a failure), because this is a user
			// call that may be delayed by the chain beyond some other state transition.
			panic("Invalid sector state in RecoverFaults")
		}
	}
}

// DeclareFaults penalizes miners (slashStorageDealCollateral and suspendPower)
// TODO: decide how much storage collateral to slash
// - State Transition
//   - Active / Commited / Recovering -> Failing
// - Update SectorStates
// - Remove Active / Commited / Recovering from ProvingSet
func (sm *StorageMinerActor_I) DeclareFaults(faultSet sector.CompactSectorSet) {

	var powerUpdate uint64 = 0

	// get all SectorNumber marked as Failing by faultSet
	for _, sectorNo := range faultSet.SectorsOn() {
		state := sm.SectorStates()[sectorNo]

		switch state.StateNumber {
		case SectorActiveSN:
			// will be incremented at end of proving period (by post or cron)
			sm.SectorStates()[sectorNo] = SectorFailing(0)
			powerUpdate = powerUpdate - sm.Info().SectorSize()
		case SectorCommittedSN:
			sm.SectorStates()[sectorNo] = SectorFailing(0)
		case SectorRecoveringSN:
			sm.SectorStates()[sectorNo] = SectorFailing(0)
		default:
			// TODO: proper failure
			panic("Invalid sector state in DeclareFaults")
		}
	}

	// Suspend power
	// SendMessage(spa.UpdatePower(powerUpdate))
}

func (sm *StorageMinerActor_I) verifySeal(onChainInfo sector.OnChainSealVerifyInfo) bool {
	// TODO: verify seal @nicola
	// TODO: get var sealRandomness sector.SealRandomness from onChainInfo.Epoch
	// TODO: sm.verifySeal(sectorID SectorID, comm sector.OnChainSealVerifyInfo, proof SealProof)

	// verifySeal will also generate CommD on the fly from CommP and PieceSize

	var pieceInfos []sector.PieceInfo // = make([]sector.PieceInfo, 0)

	for dealId := range onChainInfo.DealIDs() {
		// FIXME: Actually get the deal info from the storage market actor and use it to create a sector.PieceInfo.
		_ = dealId

		pieceInfos = append(pieceInfos, nil)
	}

	new(proving.StorageProvingSubsystem_I).VerifySeal(&sector.SealVerifyInfo_I{
		SectorID_: &sector.SectorID_I{
			MinerID_: sm.Info().Worker(), // TODO: This is actually miner address. MinerID needs to be derived.
			Number_:  onChainInfo.SectorNumber(),
		},

		OnChain_:    onChainInfo,
		PieceInfos_: pieceInfos,
	})
	return true
}

func (sm *StorageMinerActor_I) checkIfSectorExists(sectorNo sector.SectorNumber) bool {
	_, sectorExists := sm.Sectors()[sectorNo]
	if sectorExists {
		return true
	}
	return false
}

// Currently deals must be posted on chain via sma.PublishStorageDeals before CommitSector
// TODO(optimization): CommitSector could contain a list of deals that are not published yet.
func (sm *StorageMinerActor_I) CommitSector(onChainInfo sector.OnChainSealVerifyInfo) {
	isSealVerified := sm.verifySeal(onChainInfo)
	if !isSealVerified {
		// TODO: proper failure
		panic("Seal is not verified")
	}

	sectorExists := sm.checkIfSectorExists(onChainInfo.SectorNumber())
	if sectorExists {
		//TODO: proper failure
		panic("Sector already exists")
	}

	// determine lastDealExpiration from sma
	// TODO: proper onchain transaction
	// lastDealExpiration := SendMessage(sma, GetLastDealExpirationFromDealIDs(onChainInfo.DealIDs()))
	var lastDealExpiration block.ChainEpoch

	// add sector expiration to SectorExpirationQueue
	sm.SectorExpirationQueue().Add(&SectorExpirationQueueItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   lastDealExpiration,
	})

	// no need to store the proof and randomseed in the state tree
	// verify and drop, only SealCommitments{CommR, DealIDs} on chain
	sealCommitment := &sector.SealCommitment_I{
		SealedCID_:  onChainInfo.SealedCID(),
		DealIDs_:    onChainInfo.DealIDs(),
		Expiration_: lastDealExpiration, // TODO decide if we need this too
	}

	// add SectorNumber and SealCommitment to Sectors
	// set SectorState to SectorCommitted
	// Note that SectorNumber will only become Active at the next successful PoSt
	sm.Sectors()[onChainInfo.SectorNumber()] = sealCommitment
	sm.SectorStates()[onChainInfo.SectorNumber()] = SectorCommitted()
	sm.ProvingSet_.Add(onChainInfo.SectorNumber())

	// TODO: write state change
}

func getSectorNums(m map[sector.SectorNumber]SectorState) []sector.SectorNumber {
	var l []sector.SectorNumber
	for i, _ := range m {
		l = append(l, i)
	}
	return l
}
