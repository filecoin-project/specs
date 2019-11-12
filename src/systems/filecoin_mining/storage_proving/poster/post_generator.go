package poster

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sectorIndex "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"

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

	// Poroposed answer: An alternative, which avoids the downsides of both of the above, by adding a new filproofs API call:

	sdr := filproofs.SDRParams(nil, postCfg)

	return sdr.GeneratePoStCandidates(challengeSeed, faults, sectorStore)
}

func GeneratePoStProofs(postCfg sector.PoStCfg, witness sector.PoStWitness) []sector.PoStProof {
	sdr := filproofs.SDRParams(nil, postCfg)

	var proofs []sector.PoStProof

	for _, candidate := range witness.Candidates() {
		proof := sdr.GeneratePoStProof(candidate.PrivateProof())
		proofs = append(proofs, proof)
	}

	return proofs
}

// This likely belongs elsewhere, but I'm not exactly sure where and wanted to encapsulate the proofs-related logic here. So this can be thought of as example usage.
// ticketThreshold is lowest eligible winning ticket (endianness?) for this PoSt. (Maybe should be 1+ the lowest. Define.)
func GeneratePoSt(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectors []sector.SectorID, sectorStore sectorIndex.SectorStore, ticketThreshold sector.ElectionTicket) []sector.PoStProof {
	candidates := GeneratePoStCandidates(postCfg, challengeSeed, faults, sectors, sectorStore)
	var winners []sector.ElectionCandidate

	for _, candidate := range candidates {
		if candidate.Ticket().Meets(ticketThreshold) {
			winners = append(winners, candidate)
		}
	}

	witness := sector.PoStWitness_I{
		Candidates_: candidates,
	}

	return GeneratePoStProofs(postCfg, sector.PoStWitness(&witness))

	panic("TODO")

}
