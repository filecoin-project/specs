package poster

import (
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

	sector_index "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"

	util "github.com/filecoin-project/specs/util"
)

type Serialization = util.Serialization

// See "Proof-of-Spacetime Parameters" Section
// TODO: Unify with orient model.
const POST_CHALLENGE_DEADLINE = uint(480)

func (pg *PoStGenerator_I) GeneratePoStCandidates(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, candidateCount int, sectors []sector.SectorID, sectorStore sector_index.SectorStore) []sector.ChallengeTicket {
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

	sdr := makeStackedDRGForPoSt(postCfg)
	var sectorNumbers []sector.SectorNumber
	for _, s := range sectors {
		sectorNumbers = append(sectorNumbers, s.Number())
	}

	return sdr.GeneratePoStCandidates(challengeSeed, sectorNumbers, candidateCount, sectorStore)
}

func (pg *PoStGenerator_I) GeneratePoStProof(postCfg sector.PoStCfg, witness sector.PoStWitness) sector.PoStProof {
	sdr := makeStackedDRGForPoSt(postCfg)
	var privateProofs []sector.PrivatePoStProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return sdr.GeneratePoStProof(privateProofs)
}

// This likely belongs elsewhere, but I'm not exactly sure where and wanted to encapsulate the proofs-related logic here. So this can be thought of as example usage.
// ticketThreshold is lowest non-winning ticket (endianness?) for this PoSt.
func GeneratePoSt(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectors []sector.SectorID, sectorStore sector_index.SectorStore, ticketThreshold sector.ElectionTicket) sector.PoStProof {
	candidates := GeneratePoStCandidates(postCfg, challengeSeed, faults, sectors, sectorStore)
	var winners []sector.ElectionCandidate

	for _, candidate := range candidates {
		if candidate.Ticket().IsBelow(ticketThreshold) {
			winners = append(winners, candidate)
		}
	}

	witness := sector.PoStWitness_I{
		Candidates_: winners,
	}

	return GeneratePoStProof(postCfg, sector.PoStWitness(&witness))
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
