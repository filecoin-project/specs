package poster

import (
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	util "github.com/filecoin-project/specs/util"
)

type Serialization = util.Serialization

// See "Proof-of-Spacetime Parameters" Section
// TODO: Unify with orient model.
const POST_CHALLENGE_DEADLINE = uint(480)

func (pg *PoStGenerator_I) GeneratePoStCandidates(challengeSeed abi.PoStRandomness, candidateCount int, sectors []abi.SectorID) []abi.PoStCandidate {
	// Question: Should we pass metadata into FilProofs so it can interact with SectorStore directly?
	// Like this:
	// PoStReponse := SectorStorageSubsystem.GeneratePoSt(sectorSize, challenge, faults, sectorsMetatada);

	// Question: Or should we resolve + manifest trees here and pass them in?
	// Like this:
	// trees := sectorsMetadata.map(func(md) { SectorStorage.GetMerkleTree(md.MerkleTreePath) });
	// Done this way, we redundantly pass the tree paths in the metadata. At first thought, the other way
	// seems cleaner.
	// PoStReponse := SectorStorageSubsystem.GeneratePoSt(sectorSize, challenge, faults, sectorsMetadata, trees);

	// For now, dodge this by passing the whole SectorStore. Once we decide how we want to represent this, we can narrow the call.

	return filproofs.GenerateElectionPoStCandidates(challengeSeed, sectors, candidateCount, pg.SectorStore())
}

func (pg *PoStGenerator_I) CreateElectionPoStProof(randomness abi.PoStRandomness, postCandidates []abi.PoStCandidate) []abi.PoStProof {
	var privateProofs []abi.PrivatePoStCandidateProof

	for _, candidate := range postCandidates {
		privateProofs = append(privateProofs, candidate.PrivateProof)
	}

	return filproofs.CreateElectionPoStProof(privateProofs, randomness)
}

func (pg *PoStGenerator_I) CreateSurprisePoStProof(randomness abi.PoStRandomness, postCandidates []abi.PoStCandidate) []abi.PoStProof {
	var privateProofs []abi.PrivatePoStCandidateProof

	for _, candidate := range postCandidates {
		privateProofs = append(privateProofs, candidate.PrivateProof)
	}

	return filproofs.CreateSurprisePoStProof(privateProofs, randomness)
}
