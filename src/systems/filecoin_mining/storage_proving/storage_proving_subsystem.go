package storage_proving

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	util "github.com/filecoin-project/specs/util"
)

func (sps *StorageProvingSubsystem_I) VerifySeal(sv sector.SealVerifyInfo) StorageProvingSubsystem_VerifySeal_FunRet {
	cfg := filproofs.SDRCfg_I{
		SealCfg_: sv.SealCfg(),
	}
	sdr := filproofs.WinSDRParams(&cfg)

	result := sdr.VerifySeal(sv)

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

	poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)

	todo := make([]sector.PoStCandidate, 0)
	return todo
}

func (sps *StorageProvingSubsystem_I) CreateElectionPoStProof(challengeSeed sector.PoStRandomness, challengeTickets []sector.PoStCandidate) sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: challengeTickets,
	}

	var poster = sps.PoStGenerator()
	return poster.CreateElectionPoStProof(witness)
}

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateSurprisePoStCandidates(challengeSeed sector.PoStRandomness, sectorIDs []sector.SectorID) []sector.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * node_base.SPOST_SAMPLE_RATE_NUM / node_base.SPOST_SAMPLE_RATE_DENOM)

	var poster = sps.PoStGenerator()

	poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)

	todo := make([]sector.PoStCandidate, 0)
	return todo
}

func (sps *StorageProvingSubsystem_I) CreateSurprisePoStProof(challengeSeed sector.PoStRandomness, challengeTickets []sector.PoStCandidate) sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: challengeTickets,
	}

	var poster = sps.PoStGenerator()
	return poster.CreateSurprisePoStProof(witness)
}
