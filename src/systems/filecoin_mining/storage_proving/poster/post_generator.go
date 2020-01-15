package poster

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

	util "github.com/filecoin-project/specs/util"
)

type Serialization = util.Serialization

// See "Proof-of-Spacetime Parameters" Section
// TODO: Unify with orient model.
const POST_CHALLENGE_DEADLINE = uint(480)

func (pg *PoStGenerator_I) GeneratePoStCandidates(challengeSeed sector.PoStRandomness, candidateCount int, sectors []sector.SectorID) []sector.PoStCandidate {
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

	return filproofs.GenerateElectionPoStCandidates(pg.PoStCfg(), challengeSeed, sectors, candidateCount, pg.SectorStore())
}

func (pg *PoStGenerator_I) CreateElectionPoStProof(randomness sector.PoStRandomness, witness sector.PoStWitness) []sector.PoStProof {
	var privateProofs []sector.PrivatePoStCandidateProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return filproofs.CreateElectionPoStProof(pg.PoStCfg(), privateProofs, randomness)
}

func (pg *PoStGenerator_I) CreateSurprisePoStProof(postCfg sector.PoStInstanceCfg, randomness sector.PoStRandomness, witness sector.PoStWitness) []sector.PoStProof {
	var privateProofs []sector.PrivatePoStCandidateProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return filproofs.CreateSurprisePoStProof(pg.PoStCfg(), privateProofs, randomness)
}
