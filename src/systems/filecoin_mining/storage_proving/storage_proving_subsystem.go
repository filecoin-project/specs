package storage_proving

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	sector_index "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"

	//	poster "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving/poster"
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

func (sps *StorageProvingSubsystem_I) ComputeUnsealedSectorCID(sectorSize sector.SectorSize, pieceInfos []*sector.PieceInfo_I) StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet {
	unsealedCID, err := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos)

	if err != nil {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_err(StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_err(err))
	} else {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_unsealedSectorCID(
			StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_unsealedSectorCID(unsealedCID))
	}
}

// TODO: get from consts
const EPOST_SAMPLE_NUM = 1
const EPOST_SAMPLE_DENOM = 25

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateElectionPoStCandidates(challengeSeed sector.PoStRandomness, sectorIDs []sector.SectorID) []sector.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * EPOST_SAMPLE_NUM / EPOST_SAMPLE_DENOM)

	// TODO: Get these correctly.
	var cfg sector.PoStCfg
	var sectorStore sector_index.SectorStore

	var poster = sps.PoStGenerator()

	poster.GeneratePoStCandidates(cfg, challengeSeed, numChallengeTickets, sectorIDs, sectorStore)

	todo := make([]sector.PoStCandidate, 0)
	return todo
}

func (sps *StorageProvingSubsystem_I) GenerateElectionPoStProof(challengeSeed sector.PoStRandomness, challengeTickets []sector.PoStCandidate) sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: challengeTickets,
	}

	// TODO: Get this correctly. Maybe move into PoStGenerator struct.
	var cfg sector.PoStCfg

	var poster = sps.PoStGenerator()
	return poster.GeneratePoStProof(cfg, witness)
}

// TODO: get from consts
const SPOST_SAMPLE_NUM = 1
const SPOST_SAMPLE_DENOM = 50

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateSurprisePoStCandidates(challengeSeed sector.PoStRandomness, sectorIDs []sector.SectorID) []sector.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * SPOST_SAMPLE_NUM / SPOST_SAMPLE_DENOM)

	// TODO: Get these correctly.
	var cfg sector.PoStCfg
	var sectorStore sector_index.SectorStore

	var poster = sps.PoStGenerator()

	poster.GeneratePoStCandidates(cfg, challengeSeed, numChallengeTickets, sectorIDs, sectorStore)

	todo := make([]sector.PoStCandidate, 0)
	return todo
}

func (sps *StorageProvingSubsystem_I) GenerateSurprisePoStProof(challengeSeed sector.PoStRandomness, challengeTickets []sector.PoStCandidate) sector.PoStProof {
	witness := &sector.PoStWitness_I{
		Candidates_: challengeTickets,
	}

	// TODO: Get this correctly. Maybe move into PoStGenerator struct.
	var cfg sector.PoStCfg

	var poster = sps.PoStGenerator()
	return poster.GeneratePoStProof(cfg, witness)
}
