---
title: Filproofs
entries:
 - components
# suppressMenu: true
---
'type SectorId    UInt64
type ProverId    UInt64
type Commitment  [32]byte
type DataElement [32]byte
type Rephash     func(left DataElement, right DataElement) hash Commitment
type Challenge   UInt64

type FaultSet   []Fault // More complicated than this.

// Constants can be zero-arg methods of structs (for now).

// Add data structures to minimize individual parameters.
// PublicInputs
// Configs.

type SectorConfig struct {
       SectorSize     UInt,
       SubsectorCount UInt,
       Partitions     Uint,
}

// This is metadata required to generate a PoSt proof for a sector.
// These should be stored and indexed somewhere by CommR.
type SectorPersistentAux struct {
       CommRLast      Commitment,
       CommC          Commitment,
}

type FilproofsSubsystem struct {
    Seal(
        SectorConfig   SectorConfig,
	unsealedPath   String,
	sealedPath     String,
	proverId       ProverId,
	ticket         Ticket, // Assuming this is defined  elsewhere.
	sectorId       SectorId) ProofError | SealResponse

    VerifySeal(
	CommD    Commitment,
	CommR    Commitment,
	Proof    SNARKProof,
	ProverId ProverId,
	Ticket   Ticket,
	SectorId SectorId,
	) ProofError | bool

    Unseal(
         // TODO
    )

   GeneratePoSt(
     Challenge Challenge,
     CommRs    [commmitments],
     Trees     MerkleTree<Rephash>,
     Faults    FaultSet
   )

   VerifyPoSt(
       // TODO
   )

    GenerateCommP()
    GenerateCommD()

    GeneratePieceInclusionProofs(
	Tree          MerkleTree<Rephash>,
	PieceStarts   []Uint,
	PieceLength   []Uint,
    ) Error | []PieceInclusionProof

    GeneratePieceInclusionProof(
	Tree          MerkleTree<Rephash>,
	PieceStart    Uint,
	PieceLength   Uint,
    ) Error | []PieceInclusionProof

   VerifyPieceInclusionProof(
   )

   MaxUnsealedBytesPerSector(
      SectorSize UInt
   ) UInt
}

type ProofError struct {
}

type SealResponse struct {
     CommD                 Commitment,
     CommR                 Commitment,
     Proof                 SealProof,
     PersistentAux         SectorPersistentAux,
     MerkleTreePath        Path, // TODO: This may be a partially-cached tree.
}

type SealPublicInputs {

}

type SealProofConfig {
    partitionCount UInt,
    subsectorsCount UInt,
}

type SNARKProof<curve, system> {
     config SealProofConfig,
     proofBytes []bytes
}

type FilecoinSNAKRProof<bls12-381, Groth16>

type SealProof struct {
     snarkProof      SNARKProof,
     susbsectorCount Uvarint,
}

type PieceInclusionProof {
}
