package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/poster"


	// TODO(enhancement): decision is to currently account for power based on sector
	// an alternative proposal is to account for power based on active deals
func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission) {
	// Verify Proof
	// TODO
	// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
	// challenges := GenerateChallengesForPoSt(r, keys(sm.Sectors))
	// verifyPoSt(challenges, TODO)

	// Pay miner
	// for _, sealCommitment := range sm.Sectors {
	//   SendMessage(sma.ProcessStorageDealsPayment(sealCommitment.DealIDs))
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
		// FaultSet[expiredSectorNumber] = 0
	}

	// Update Power and Handle Faults

	// newPower := power - len(expiredSectorsNumber) * sm.info.sectorSize()
	// SendMessage(SPA.UpdatePower(newPower))
	// TODO: either SM or SPA must return the pledge collateral to the miner
	panic("TODO")
}

func (sm *StorageMinerActor_I) DeclareFault(faultSet sector.FaultSet) {
	// Handle Fault
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
		Expiration_: latestDealExpiry,
	})

	_, found := sm.Sectors()[onChainInfo.SectorNumber()]

	if found {
		// TODO: throw error
		panic("sector already in there")
	}

	sm.Sectors()[onChainInfo.SectorNumber()] = sealCommitment

	// TODO write state change
}
