package storage_mining

import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving/poster"

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
// postSubmission: is the post for this proving period
// faultSet: the miner announces all of their current faults
// TODO(enhancement): decide whether declared faults sectors should be
// penalized in the same way as undeclared sectors
// Workflow:
// - Verify PoSt Submission
// - Process Unproven Sectors (move Sectors from Unproven to Active)
// - Process Faulty Sectors (penalize faults, recover sectors, delete faulty sectors)
// - Process Active Sectors (pay miners)
// - Process Expired Sectors (settle deals..)
// TODO: if something is faulty, move it to committed instead
func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission) {
	// Verify correct PoSt Submission:
	isPoStVerified := sm.verifyPoStSubmission(postSubmission)
	if !isPoStVerified {
		// TODO: mark all sectors as faulty
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

	// Pay all the Active sectors
	// Note: this must happen before marking sectors as expired.
	// TODO: Pay miner in a single batch message
	// for _, sc := range sm.ActiveSectors() {
	// SendMessage(sma.ProcessStorageDealsPayment(sc.DealIDs()))
	// }

	// State change: Active -> Deleted (because they are expired)
	// Note: this must happen as last state transition check to ensure that
	// payments and faults have been accounted for correctly.
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
// move SectorNumber and SealCommitment
// from sm.FaultySectors to UnprovenSectors
func (sm *StorageMinerActor_I) RecoverFaults(newFaults sector.FaultSet) {
	// 2. State change: Faulty -> Active
	{
		// Handle Recovered faults:
		// If a sector is not in sm.NextFaultSet at this point, it means that it
		// was just proven in this proving period.
		// However, if this had a counter in sm.Faults, then it means that it was
		// faulty in the previous proving period and then recovered.
		// In that case, reset the counter and resume power.

		// resumedSectorsCount := 0
		// for previouslyFaulty := range keys(sm.Faults) {
		//   if (previouslyFaulty not in sm.NextFaultSet() and previouslyFaulty in sm.UnprovenSectors_) {
		//     delete(sm.Faults, previouslyFaulty)
		//     resumedSectorsCount = resumedSectorsCount + 1
		//   }
		// }
		// powerUpdate = powerUpdate + resumedSectorsCount * sm.info().sectorSize()
	}
}

// DeclareFaults penalizes miners (slashStorageDealCollateral and suspendPower)
// and moves SectorNumber and SealCommitment
// from sm.ActiveSectors to sm.FaultySectors
func (sm *StorageMinerActor_I) DeclareFaults(newFaults sector.FaultSet) {

	// 1. State change: Active -> Faulty
	// Note: this must happen after Unproven->Active, in order to account for
	// sectors that have been committed in this proving period which also happen
	// to be faulty.
	// Note: if a Sector is faulty when it's in Unproven state, it wont pass verifyPoSt
	// TODO: @nicola verify this
	// Note: if the miner has not recovered the faults, they are re-declared
	// automatically.
	{
		// Handle faulty sectors
		// sm.DeclareFault(sm.nextFaultSet)
	}

	// Handle Fault

	// TODO: the faults that are declared after post challenge,
	// are faults for the next proving period

	// TODO: below is a bit complicated, it should be simplified.

	// Update Fault Set

	// var lostPower uint
	// for sectorNumber := range newFaultSet {
	//   // Avoid penalizing the miner multiple times in the same proving period
	//   if !(sector is in sm.NextFaultSet()) {
	//     // it is a new fault:
	//     lostPower = lostPower + sm.info().sectorSize()
	//     if (sectorNumber not in sm.Faults) sm.Faults[sectorNumber] = 0
	//     sm.Faults[sectorNumber] = sm.Faults[sectorNumber] + 1
	//     if (sm.Faults[sectorNumber] > MAX_CONSECUTIVE_FAULTS_ALLOWED) {
	//       Sector is lost, delete it, slash storage deal collateral and all pledge
	//       SendMessage(sma.SlashAllStorageDealCollateral(dealIDs)
	//       SendMessage(spa.SlashAllPledgeCollateral(sectorNumber)
	//       delete(sm.Sectors(), expiredSectorNumber)
	//     } else {
	//       append(sm.UnprovenSectors_, sectorNumber)
	//     }
	//   }
	// }

	// Delete Power
	// SendMessage(sma.UpdatePower(- lostPower))

	// Store updated fault set
	// sm.nextFaultSet.applyDiff(newFaultSet)

	// TODO: check if we want to penalize some collateral for losing some files

	panic("TODO")
}

func (sm *StorageMinerActor_I) verifySeal(onChainInfo sector.OnChainSealVerifyInfo) bool {
	panic("TODO")
	return true
}

func (sm *StorageMinerActor_I) checkIfSectorExists(sectorNo sector.SectorNumber) bool {
	_, isUnproven := sm.UnprovenSectors()[sectorNo]
	_, isActive := sm.ActiveSectors()[sectorNo]
	_, isFaulty := sm.FaultySectors()[sectorNo]
	if isUnproven || isActive || isFaulty {
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

	// determine lastDealExpiration from sma
	// TODO: proper onchain transaction
	// TODO: check if this makes sense as an onchain transaction
	// lastDealExpiration := SendMessage(sma, GetLastDealExpirationFromDealIDs(onChainInfo.DealIDs()))
	var lastDealExpiration block.ChainEpoch

	// no need to store the proof and randomseed in the state tree
	// just verify and drop it, only SealCommitments{CommD, CommR, DealIDs} on chain
	// TODO: @porcuquine verifies
	sealCommitment := &sector.SealCommitment_I{
		UnsealedCID_: onChainInfo.UnsealedCID(),
		SealedCID_:   onChainInfo.SealedCID(),
		DealIDs_:     onChainInfo.DealIDs(),
		Expiration_:  lastDealExpiration,
	}

	// add sector expiration to SectorExpirationQueue
	sm.SectorExpirationQueue().Add(&SectorExpirationQueuItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   lastDealExpiration,
	})

	sectorExists := sm.checkIfSectorExists(onChainInfo.SectorNumber())
	if sectorExists {
		//TODO: proper failure
		panic("Sector already exists")
	}

	// add SectorNumber and SealCommitment to UnprovenSectors
	// Note that SectorNumber is not in ActiveSectors yet
	// it will become Active at the next proving period
	sm.UnprovenSectors()[onChainInfo.SectorNumber()] = sealCommitment

	// TODO: write state change
}
