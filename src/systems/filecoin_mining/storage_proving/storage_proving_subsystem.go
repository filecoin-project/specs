package storage_proving

import (
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	util "github.com/filecoin-project/specs/util"
)

func (sps *StorageProvingSubsystem_I) VerifySeal(sv abi.SealVerifyInfo) StorageProvingSubsystem_VerifySeal_FunRet {
	registeredProof := sv.OnChain.RegisteredProof

	verifier := filproofs.MakeSealVerifier(registeredProof)
	result := verifier.VerifySeal(sv)

	return StorageProvingSubsystem_VerifySeal_FunRet_Make_ok(StorageProvingSubsystem_VerifySeal_FunRet_ok(result)) //,
}

func (sps *StorageProvingSubsystem_I) ComputeUnsealedSectorCID(sectorSize abi.SectorSize, pieceInfos []abi.PieceInfo) StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet {
	unsealedCID, err := filproofs.ComputeUnsealedSectorCIDFromPieceInfos(sectorSize, pieceInfos)

	if err != nil {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_err(StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_err(err))
	} else {
		return StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_Make_unsealedSectorCID(
			StorageProvingSubsystem_ComputeUnsealedSectorCID_FunRet_unsealedSectorCID(unsealedCID))
	}
}

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateElectionPoStCandidates(challengeSeed abi.PoStRandomness, sectorIDs []abi.SectorID) []abi.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * node_base.EPOST_SAMPLE_RATE_NUM / node_base.EPOST_SAMPLE_RATE_DENOM)

	var poster = sps.PoStGenerator()

	return poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)
}

func (sps *StorageProvingSubsystem_I) CreateElectionPoStProof(challengeSeed abi.PoStRandomness, candidates []abi.PoStCandidate) []abi.PoStProof {
	var poster = sps.PoStGenerator()
	return poster.CreateElectionPoStProof(challengeSeed, candidates)
}

// TODO also return error
func (sps *StorageProvingSubsystem_I) GenerateSurprisePoStCandidates(challengeSeed abi.PoStRandomness, sectorIDs []abi.SectorID) []abi.PoStCandidate {
	numChallengeTickets := util.UInt(len(sectorIDs) * node_base.SPOST_SAMPLE_RATE_NUM / node_base.SPOST_SAMPLE_RATE_DENOM)

	var poster = sps.PoStGenerator()

	return poster.GeneratePoStCandidates(challengeSeed, numChallengeTickets, sectorIDs)
}

func (sps *StorageProvingSubsystem_I) CreateSurprisePoStProof(challengeSeed abi.PoStRandomness, candidates []abi.PoStCandidate) []abi.PoStProof {
	var poster = sps.PoStGenerator()
	return poster.CreateSurprisePoStProof(challengeSeed, candidates)
}
