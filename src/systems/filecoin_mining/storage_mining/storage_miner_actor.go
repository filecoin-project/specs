package storage_mining

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving/poster"

var CONSECUTIVE_FAULT_COUNT_LIMIT uint8

// If a Post is missed (either due to faults being not declared on time or
// because the miner run out of time, every sector is reported as faulty
// for the current proving period.
func (sm *StorageMinerActor_I) OnFaultBeingSpotted() {
	// var allFaultBitField FaultSet
	// TODO: ideally, both DeclareFault and OnFaultBeingSpotted call a method "ApplyFaultConsequences"
	// DeclareFaults(allFaultBitField)
	// slash pledge collateral
	panic("TODO")
}

func (sm *StorageMinerActor_I) verifyPoStSubmission(postSubmission poster.PoStSubmission) bool {
	panic("TODO")
	// 1. A proof must be submitted after the postRandomness for this proving
	// period is on chain
	{
		// if rt.ChainEpoch < sm.ProvingPeriodEnd - challengeTime {
		//   panic("too early")
		// }
	}

	// 2. A proof must be a valid snark proof with the correct public inputs
	{
		// 2.1 Get randomness from the chain at the right epoch
		// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
		// 2.2 Generate the set of challenges
		// challenges := GenerateChallengesForPoSt(r, keys(sm.Sectors))
		// 2.3 Verify the PoSt Proof
		// verifyPoSt(challenges, TODO)
	}
	return true
}

func extendSectors(toSectors map[sector.SectorNumber]sector.SealCommitment, fromSectors map[sector.SectorNumber]sector.SealCommitment) map[sector.SectorNumber]sector.SealCommitment {

	for sectorNo, sc := range toSectors {
		fromSectors[sectorNo] = sc
	}

	return fromSectors
}

func (sm *StorageMinerActor_I) activateUnprovenSectors(unprovenSectors map[sector.SectorNumber]sector.SealCommitment) {
	sm.ActiveSectors_ = extendSectors(sm.ActiveSectors(), sm.UnprovenSectors())
	sm.UnprovenSectors_ = map[sector.SectorNumber]sector.SealCommitment{}
	// TDOD check what committing state change here looks like
}

// decision is to currently account for power based on sector
// with at least one active deals and deals cannot be updated
// an alternative proposal is to account for power based on active deals
// an improvement proposal is to allow storage deal update in a sector

// TODO: decide whether declared faults sectors should be
// penalized in the same way as undeclared sectors
// Workflow:
// - Verify PoSt Submission
// - Process ProvingSet
//   - State Transitions
//     - Committed -> Active and credit power
//     - Recovering -> Active and credit power
//   - Remove Sector from ConsecutiveFaultCounts
// - Process Sectors in ConsecutiveFaultCounts
//	   - increment count
//     - clear Sector and slash pledge collateral if count >= CONSECUTIVE_FAULT_COUNT_LIMIT
// - Process Active Sectors (pay miners)
// - Process Expired Sectors (settle deals and return storage collateral to miners)
//     - State Transition
//       - Faulted / Recovering / Active / Committed -> Cleared
//     - Remove Sector from Sectors, ProvingSet
func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission) {
	// Verify correct PoSt Submission:
	isPoStVerified := sm.verifyPoStSubmission(postSubmission)
	if !isPoStVerified {
		// no state transition
		// just error out and miner should submitPoSt again
		// TODO: proper failture
		panic("TODO")
	}

	// The proof is verified, proceed to sector state transitions:
	var powerUpdate uint64

	// State change: Unproven -> Active
	// Note: this must be the first state transition check.
	// check if there are any sectors in UnprovenSectors_
	// if so, their PoSt has been verified, credit power for these sectors
	numUnprovenSector := uint64(len(sm.UnprovenSectors()))
	if numUnprovenSector > 0 {
		// activate unproven sectors will also empty sm.UnprovenSectors_
		sm.activateUnprovenSectors(sm.UnprovenSectors())
		powerUpdate = powerUpdate + numUnprovenSector*sm.Info().SectorSize()
	}

	// Process FaultySectors
	// Increment ConsecutiveFaultCounts
	for sectorNo, _ := range sm.FaultySectors() {
		prevFaultCount := sm.ConsecutiveFaultCounts()[sectorNo]
		newFaultCount := prevFaultCount + 1
		if newFaultCount >= CONSECUTIVE_FAULT_COUNT_LIMIT {
			// TODO heavy penalization
			// slash pledge collateral and delete sector?
			panic("TODO")
		} else {
			sm.ConsecutiveFaultCounts()[sectorNo] = newFaultCount
		}
	}

	// Pay all the Active sectors
	// Note: this must happen before marking sectors as expired.
	// TODO: Pay miner in a single batch message
	// for _, sc := range sm.ActiveSectors() {
	// SendMessage(sma.ProcessStorageDealsPayment(sc.DealIDs()))
	// }

	// State change: Active -> Deleted (because they are expired)
	// Note: this must happen as last state transition check to ensure that
	// payments and faults have been accounted for correctly
	// Handle expired sectors
	var numExpiredSectors = uint64(0)
	var currEpoch block.ChainEpoch // TODO: replace this with rt.State().Epoch()
	// go through sm.SectorExpirationQueue() and get the expiredSectorNumbers

	// TODO verify this
	expirationPeek := sm.SectorExpirationQueue().Peek().Expiration()
	for expirationPeek <= currEpoch {
		expiredSectorNo := sm.SectorExpirationQueue().Pop().SectorNumber()

		_, isUnproven := sm.UnprovenSectors()[expiredSectorNo]
		_, isFaulty := sm.FaultySectors()[expiredSectorNo]
		_, isActive := sm.ActiveSectors()[expiredSectorNo]

		// Note: this should never happen
		if isUnproven {
			panic("Expired sector is unproven")
		}

		if isFaulty {
			// TODO: check if there is any fault that we should handle here
			// Slash everything?
			panic("TODO")
			// delete(sm.Faults(), expiredSectorNumber)
		}

		if isActive {
			// Note: in order to verify if something was stored in the past, one must
			// scan the chain. SectorNumbers can be re-used.

			// Settle deals
			// expiredDealIDs := sm.ActiveSectors()[expiredSectorNo].DealIDs()
			delete(sm.ActiveSectors_, expiredSectorNo)
			numExpiredSectors = numExpiredSectors + 1
			// TODO: SendMessage(sma.SettleExpiredDeals(expiredDealIDs))

		}

		expirationPeek = sm.SectorExpirationQueue().Peek().Expiration()

	}

	// Reset Proving Period and report power updates
	// sm.ProvingPeriodEnd_ = PROVING_PERIOD_TIME
	// powerUpdate = powerUpdate - numExpiredSectors * sm.Info().SectorSize()
	// SendMessage(sma.UpdatePower(powerUpdate))

	// Return PledgeCollateral
	// TODO: SendMessage(spa.Depledge)
}

// RecoverFaults checks if miners have sufficent collateral
// - State Transition
//   - Faulted -> Recovering with the same FaultCount
// - Add SectorNumbers to ProvingSet
// Note that power is not updated until it is active
func (sm *StorageMinerActor_I) RecoverFaults(recoveredSectorNo []sector.SectorNumber) {
	// State change: Faulty -> Unproven
	for _, sectorNo := range recoveredSectorNo {
		sc, isFaulty := sm.FaultySectors()[sectorNo]
		if !isFaulty {
			// TODO proper failure
			panic("Sector was not at fault")
		}

		// Check if miners have sufficient balances in sma
		// SendMessage(sma.PublishStorageDeals) or sma.ResumeStorageDeals?
		// TODO need to remove storage deals from sma at fault?
		// throw if miner cannot cover StorageDealCollateral

		// update Sectors
		delete(sm.FaultySectors(), sectorNo)
		// WRONG: FaultCounts should be reset when an unproven sector is proven
		delete(sm.ConsecutiveFaultCounts(), sectorNo)
		sm.UnprovenSectors()[sectorNo] = sc
	}
}

// DeclareFaults penalizes miners (slashStorageDealCollateral and suspendPower)
// - State Transition
//   - Active -> Faulted
// - Remove Faulted from ProvingSet
func (sm *StorageMinerActor_I) DeclareFaults(faultSet sector.CompactSectorSet) {

	// TODO: DeclareFaults can also move from Unproven to Faulty

	// State change: Active -> Faulty
	// Note: this must happen after Unproven->Active, in order to account for
	// sectors that have been committed in this proving period which also happen
	// to be faulty.
	// Note: if a Sector is faulty when it's in Unproven state, it wont pass verifyPoSt
	// TODO: @nicola verify this
	// Note: if the miner has not recovered the faults, they are re-declared
	// automatically.

	// slash storage collateral
	// TODO: decide how much storage collateral to slash
	// SendMessage(sma.slashStorageDealCollateral(sc.DealIDs()))

	// Update mapping


	// suspend power
	// powerDiff := uint64(len(faultySectorNo) * sm.Info().SectorSize())
	// SendMessage(spa.UpdatePower(-powerDiff))
}

func (sm *StorageMinerActor_I) verifySeal(onChainInfo sector.OnChainSealVerifyInfo) bool {
	panic("TODO")
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
// TODO: as an optimization, in the future CommitSector could contain a list of
// deals that are not published yet.
func (sm *StorageMinerActor_I) CommitSector(onChainInfo sector.OnChainSealVerifyInfo) {
	// TODO verify seal @nicola
	// var sealRandomness sector.SealRandomness
	// TODO: get sealRandomness from onChainInfo.Epoch
	// TODO: sm.verifySeal(sectorID SectorID, comm sector.OnChainSealVerifyInfo, proof SealProof)
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
	// TODO: check if this makes sense as an onchain transaction
	// lastDealExpiration := SendMessage(sma, GetLastDealExpirationFromDealIDs(onChainInfo.DealIDs()))
	var lastDealExpiration block.ChainEpoch

	// add sector expiration to SectorExpirationQueue
	sm.SectorExpirationQueue().Add(&SectorExpirationQueuItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   lastDealExpiration,
	})

	// no need to store the proof and randomseed in the state tree
	// just verify and drop it, only SealCommitments{CommD, CommR, DealIDs} on chain
	// TODO: @porcuquine verifies
	sealCommitment := &sector.SealCommitment_I{
		UnsealedCID_: onChainInfo.UnsealedCID(),
		SealedCID_:   onChainInfo.SealedCID(),
		DealIDs_:     onChainInfo.DealIDs(),
		Expiration_:  lastDealExpiration,
	}

	// add SectorNumber and SealCommitment to Sectors
	// set SectorState to SectorCommitted
	// Note that SectorNumber will only become Active at the next successful PoSt
	sm.Sectors()[onChainInfo.SectorNumber()] = sealCommitment
	sm.SectorStates()[onChainInfo.SectorNumber()] = SectorCommitted

	// TODO: write state change
}
