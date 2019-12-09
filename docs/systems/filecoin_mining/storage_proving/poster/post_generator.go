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

	sdr := makeStackedDRGForPoSt(pg.PoStCfg())
	return sdr.GenerateElectionPoStCandidates(challengeSeed, sectors, candidateCount, pg.SectorStore())
}

func (pg *PoStGenerator_I) CreateElectionPoStProof(postCfg sector.PoStCfg, randomness sector.PoStRandomness, witness sector.PoStWitness) sector.PoStProof {
	sdr := makeStackedDRGForPoSt(postCfg)
	var privateProofs []sector.PrivatePoStCandidateProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return sdr.CreateElectionPoStProof(privateProofs, randomness)
}

func (pg *PoStGenerator_I) CreateSurprisePoStProof(postCfg sector.PoStCfg, randomness sector.PoStRandomness, witness sector.PoStWitness) sector.PoStProof {
	sdr := makeStackedDRGForPoSt(postCfg)
	var privateProofs []sector.PrivatePoStCandidateProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return sdr.CreateSurprisePoStProof(privateProofs, randomness)
}

func makeStackedDRGForPoSt(postCfg sector.PoStCfg) (sdr *filproofs.WinStackedDRG_I) {
	var cfg filproofs.SDRCfg_I

	switch postCfg.Type() {
	case sector.PoStType_ElectionPoSt:
		cfg = filproofs.SDRCfg_I{
			ElectionPoStCfg_: postCfg,
		}
	case sector.PoStType_SurprisePoSt:
		cfg = filproofs.SDRCfg_I{
			SurprisePoStCfg_: postCfg,
		}
	}

	return filproofs.WinSDRParams(&cfg)
}
