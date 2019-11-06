package poster

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sectorIndex "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"

// See "Proof-of-Spacetime Parameters" Section
// TODO: Unify with orient model.
const POST_CHALLENGE_DEADLINE = uint(480)

func GeneratePoStWitness(postCfg sector.PoStCfg, challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectors []sector.SectorID, sectorStore sectorIndex.SectorStore) sector.PoStWitness {
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

	return sdr.GeneratePoStWitness(challengeSeed, faults, sectorStore)
}

func GeneratePoStProof(postCfg sector.PoStCfg, witness sector.PoStWitness) sector.PoStProof {
	sdr := filproofs.SDRParams(nil, postCfg)
	return sdr.GeneratePoStProof(witness)
}
