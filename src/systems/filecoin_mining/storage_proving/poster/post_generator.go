package poster

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sectorIndex "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
import util "github.com/filecoin-project/specs/util"

type Serialization = util.Serialization

// See "Proof-of-Spacetime Parameters" Section
// TODO: Unify with orient model.
const POST_CHALLENGE_DEADLINE = uint(480)

func GeneratePoStCandidates(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectors []sector.SectorID, sectorStore sectorIndex.SectorStore) []sector.ElectionCandidate {
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

	return sdr.GeneratePoStCandidates(challengeSeed, faults, sectorStore)
}

func GeneratePoStProof(postCfg sector.PoStCfg, witness sector.PoStWitness) sector.PoStProof {
	sdr := makeStackedDRGForPoSt(postCfg)
	var privateProofs []sector.PrivatePoStProof

	for _, candidate := range witness.Candidates() {
		privateProofs = append(privateProofs, candidate.PrivateProof())
	}

	return sdr.GeneratePoStProof(privateProofs)
}

// This likely belongs elsewhere, but I'm not exactly sure where and wanted to encapsulate the proofs-related logic here. So this can be thought of as example usage.
// ticketThreshold is lowest non-winning ticket (endianness?) for this PoSt.
func GeneratePoSt(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectors []sector.SectorID, sectorStore sectorIndex.SectorStore, ticketThreshold sector.ElectionTicket) sector.PoStProof {
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

func makeStackedDRGForPoSt(postCfg sector.PoStCfg) (sdr *filproofs.StackedDRG_I) {
	var cfg filproofs.SDRCfg_I

	switch postCfg.Type().(type) {
	case sector.PoStType_ElectionPoSt:
		cfg = filproofs.SDRCfg_I{
			ElectionPoStCfg_: postCfg,
		}
	case sector.PoStType_SurprisePoSt:
		cfg = filproofs.SDRCfg_I{
			SurprisePoStCfg_: postCfg,
		}
	}

	return filproofs.SDRParams(&cfg)
}

func Serialize_PoStSubmission(x PoStSubmission) Serialization {
	panic("TODO")
}

func Deserialize_PoStSubmission(x Serialization) (PoStSubmission, error) {
	panic("TODO")
}
