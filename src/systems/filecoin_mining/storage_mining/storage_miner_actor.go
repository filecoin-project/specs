package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/poster"

// TODO(enhancement): decision is to currently account for power based on sector
// an alternative proposal is to account for power based on active deals
// postSubmission: is the post for this proving period
// faultSet: the miner announces all of their current faults
// Workflow:
// - New Committed Sectors (add power)
// - Faulty Sectors (penalize faults, recover sectors)
// - Active Sectors (pay)
// - Expired Sectors
func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission, nextFaultSet sector.FaultSet) {
	// Verify Proof
	// TODO
	// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
	// challenges := GenerateChallengesForPoSt(r, keys(sm.Sectors))
	// verifyPoSt(challenges, TODO)

	// TODO: Enter newly introduced sector

	// Handle faulty sectors
	// sm.NextFaultSet_ = all zeros
	// sm.DeclareFault(nextFaultSet)

	// Handle Recovered faults:
	// If a sector is not in sm.NextFaultSet at this point, it means that it
	// was just proven in this proving period.
	// However, if this had a counter in sm.Faults, then it means that it was
	// faulty in the previous proving period and then recovered.
	// In that case, reset the counter and resume power.

	// resumedSectorsCount := 0
	// for previouslyFaulty := range keys(sm.Faults) {
	//   if (previouslyFaulty not in sm.NextFaultSet()) {
	//     delete(sm.Faults, previouslyFaulty)
	//     resumedSectorsCount = resumedSectorsCount + 1
	//   }
	// }
	// SendMessage(SPA.UpdatePower(rt.SenderAddress, resumedSectorsCount * sm.info.sectorSize()))

	// Pay miner
	// TODO: batch into a single message
	// for _, sealCommitment := range sm.Sectors {
	//   if sector is not in sm.Faults {
	//     SendMessage(sma.ProcessStorageDealsPayment(sealCommitment.DealIDs))
	//   }
	// }

	// Handle expired sectors
	// var expiredSectorsNumber [SectorNumber]
	// go through sm.SectorExpirationQueue() and get the expiredSectorsNumber

	for expiredSectorNumber := range expiredSectorsNumber {
		// Settle deals
		// TODO: SendMessage(SMA.SettleExpiredDeals(sm.Sectors()[expiredSectorNumber].DealIDs()))
		// Note: in order to verify if something was stored in the past, one must
		// scan the chain. SectorNumbers can be re-used.

		// TODO: check if there is any fault that we should handle here

		// Clean up data structures
		// delete(sm.Sectors(), expiredSectorNumber)
		// TODO: maybe nextFaultSet[expiredSectorNumber] = 0
		// TODO: SPA must return the pledge collateral to the miner
	}
	// newPower := - len(expiredSectorsNumber) * sm.info.sectorSize()
	// SendMessage(SPA.UpdatePower(rt.SenderAddress, newPower))

	panic("TODO")
}

func (sm *StorageMinerActor_I) DeclareFault(newFaults sector.FaultSet) {
	// Handle Fault

	// TODO: the faults that are declared after post challenge,
	// are faults for the next proving period

	// Update Fault Set

	// TODO: batch into a single message
	// for sectorNumber := range newFaultSet {
	//   // Avoid penalizing the miner multiple times in the same proving period
	//   if !(sector is in sm.NextFaultSet()) {
	//     if (sectorNumber not in sm.Faults) sm.Faults[sectorNumber] = 0
	//     sm.Faults[sectorNumber] = sm.Faults[sectorNumber] + 1

	//     if (sm.Faults[sectorNumber] > MAX_CONSECUTIVE_FAULTS_ALLOWED) {
	//       Sector is lost, delete it, slash storage deal collateral and all pledge
	//       SendMessage(sma.SlashAllStorageDealCollateral(dealIDs)
	//       SendMessage(spa.SlashAllPledgeCollateral(sectorNumber)
	//       delete(sm.Sectors(), expiredSectorNumber)
	//     } else {
	//       dealIDs := sm.Sectors()[sectorNumber].DealIDs()
	//       SendMessage(sma.SlashStorageDealCollateral(dealIDs)
	//     }
	//   }
	// }

	// sm.nextFaultSet.applyDiff(newFaultSet)

	panic("TODO")
}

// TODO: currently deals must be posted on chain via sma.PublishStorageDeals
// as an optimization, in the future CommitSector could contain a list of
// deals that are not published yet.
func (sm *StorageMinerActor_I) CommitSector(onChainInfo sector.OnChainSealVerifyInfo) {
	// TODO verify
	// var sealRandomness sector.SealRandomness
	// TODO: get sealRandomness from onChainInfo.Epoch
	// TODO: sm.verifySeal(sectorID SectorID, comm sector.OnChainSealVerifyInfo, proof SealProof)

	// TODO: you don't need to store the proof and the randomseed in the state tree
	// just verify it and drop it, just SealCommitments{CommD, CommR, DealIDs}
	// TODO: @porcuquine verifies

	// TODO: this is the latest expiry date of all deals.
	var latestDealExpiry ChainEpoch
	// for deal := range onChainInfo.DealIDs() {
	//   if deal.Expiry > latestDealExpiry {
	//     latestDealExpiry = deal.Expiry
	//   }
	// }

	sealCommitment := &sector.SealCommitment_I{
		UnsealedCID_: onChainInfo.UnsealedCID(),
		SealedCID_:   onChainInfo.SealedCID(),
		DealIDs_:     onChainInfo.DealIDs(),
		Expiration_:  latestDealExpiry,
	}

	sm.SectorExpirationQueue.Add(&SectorExpirationQueuItem_I{
		SectorNumber_: onChainInfo.SectorNumber(),
		Expiration_:   latestDealExpiry,
	})

	_, found := sm.Sectors()[onChainInfo.SectorNumber()]

	if found {
		// TODO: throw error
		panic("sector already in there")
	}

	sm.Sectors()[onChainInfo.SectorNumber()] = sealCommitment

	// TODO write state change
}
