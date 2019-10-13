package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import poster "github.com/filecoin-project/specs/systems/filecoin_mining/poster"

func (sm *StorageMinerActor_I) SubmitPoSt(postSubmission poster.PoStSubmission) {
	// TODO
	// postRandomness := rt.Randomness(postSubmission.Epoch, 0)
	// challenges := GenerateChallengesForPoSt(r, sm.sectorStateSets().ActiveSet())
	// verifyPoSt(challenges, TODO)


	// sectorStateSets = sm.GetExpired()
	UpdateSectorState(sectorStateSets)

	// TODO: decision is to currently account for power based on sector
	// an alternative proposal is to account for power based on active deals



	panic("TODO")
}

func (sm *StorageMinerActor_I) UpdateSectorState(sectorStateSets sector.SectorStateSets) {
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

	sealCommitment := &sector.SealCommitment_I{
		UnsealedCID_: onChainInfo.UnsealedCID(),
		SealedCID_:   onChainInfo.SealedCID(),
		DealIDs_:     onChainInfo.DealIDs(),
		Expiration_:  latestDealExpiry,
	}

	_, found := sm.Sectors()[onChainInfo.SectorNumber()]

	if found {
		// TODO: throw error
		panic("sector already in there")
	}

	sm.Sectors()[onChainInfo.SectorNumber()] = sealCommitment

	// TODO write state change
}
