package storage_proving

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	util "github.com/filecoin-project/specs/util"
)

func (sps *StorageProvingSubsystem_I) VerifySeal(sv sector.SealVerifyInfo) StorageProvingSubsystem_VerifySeal_FunRet {
	registeredProof := sv.RegisteredProof()
	proofInstance := sector.RegisteredProofInstance(registeredProof)

	var result bool

	// TODO: Presumably this can be done with interfaces or whatever method we intend for such things,
	// but for now this expresses intent simply enough.
	switch proofInstance.Algorithm() {
	case sector.ProofAlgorithm_WinStackedDRGSeal:
		result = filproofs.WinSDRParams(proofInstance.Cfg().As_SealCfg()).VerifySeal(sv)
	}

	return StorageProvingSubsystem_VerifySeal_FunRet_Make_ok(StorageProvingSubsystem_VerifySeal_FunRet_ok(result)) //,
}

func (sps *StorageProvingSubsystem_I) ComputeUnsealedSectorCID(sectorSize sector.SectorSize, pieceInfos []sector.PieceInfo) StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet {
	unsealedCID, err := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos)

	if err != nil {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_err(StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_err(err))
	} else {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_unsealedSectorCID(
			StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_unsealedSectorCID(unsealedCID))
	}
}

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateElectionPoStCandidates(challengeSeed sector.PoStRandomness, sectorIDs []sector.SectorID) []sector.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * node_base.EPOST_SAMPLE_RATE_NUM / node_base.EPOST_SAMPLE_RATE_DENOM)

	var poster = sps.PoStGenerator()

	return poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)
}

func (sps *StorageProvingSubsystem_I) CreateElectionPoStProof(challengeSeed sector.PoStRandomness, candidates []sector.PoStCandidate) []sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: candidates,
	}

	var poster = sps.PoStGenerator()
	return poster.CreateElectionPoStProof(challengeSeed, witness)
}

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateSurprisePoStCandidates(challengeSeed sector.PoStRandomness, sectorIDs []sector.SectorID) []sector.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * node_base.SPOST_SAMPLE_RATE_NUM / node_base.SPOST_SAMPLE_RATE_DENOM)

	var poster = sps.PoStGenerator()

	return poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)
}

func (sps *StorageProvingSubsystem_I) CreateSurprisePoStProof(challengeSeed sector.PoStRandomness, candidates []sector.PoStCandidate) []sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: candidates,
	}

	var poster = sps.PoStGenerator()
	return poster.CreateSurprisePoStProof(challengeSeed, witness)
}
