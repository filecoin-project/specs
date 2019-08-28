# Proofs of Spacetime

__NOTE:__ __*Proof of Spacetime*__ is in transition. Current implementations are mocked, and the final design has not been implemented. Consumers may refer to the below for reference, but nothing should be implemented until the spec is updated and synchronized with what will be the canonical construction.

### GeneratePost

`GeneratePoSt` generates a __*Proof of Spacetime*__ over `POST_SECTORS_COUNT` __*sealed sectors*__ — identified by their `commR` commitments. This is accomplished by performing a series of merkle inclusion proofs (__*Proofs of Retrievability*__). Each proof is of a challenged node in a challenged sector. The challenges are generated pseudo-randomly, based on the provided `challengeSeed`. At each time step, a number of __*Proofs of Retrievability*__ are performed. The result of each such set of __*Proofs of Retrievability*__ is used to seed challenge generation for another iteration. Repeated and necessarily sequential generation of these __*Proofs of Retrievability*__ proves that the claimed __*sealed sectors*__ existed during the time required to generate them.

Since many __*sealed sectors*__ may be proved at once, it may be the case that one or more __*sealed sectors*__ has been lost, damaged, or otherwise become impossible to validly prove. In this case, a fault is recorded and returned in an array of faults. This allows provers to selectively default on individual __*sealed sector*__ proofs while still providing a verifiable proof of their aggregate __*Proof of Spacetime*__ claims.

```
GeneratePoSt
 (
  // request represents a request to generate a proof-of-spacetime.
  commRs         [POST_SECTORS_COUNT][32]byte,  // the commR commitments corresponding to the sealed sectors to prove
  challengeSeed  [32]byte,    // a pseudo-random value to be used in challenge generation
) err Error | (
  // response contains PoST proof and any faults that may have occurred.
  faults        []uint64,    // faults encountered while proving (by index of associated commR in the input)
  proof         []byte
)
```

### VerifyPoSt

`VerifyPoSt` is the functional counterpart to `GeneratePoSt`. It takes all of `GeneratePoSt`'s output, along with those of `GeneratePost`'s inputs required to identify the claimed proof. All inputs are required because verification requires sufficient context to determine not only that a proof is valid but also that the proof indeed corresponds to what it purports to prove.

```
VerifyPoSt
 (
  // request represents a request to generate verify a proof-of-spacetime.
  commRs        [POST_SECTORS_COUNT][32]byte,        // the commRs provided to GeneratePoSt
  challengeSeed [32]byte,
  faults        []uint64
  proof         []byte,            // Multi-SNARK proof returned by GeneratePoSt 
 ) err Error | 
  isValid bool                     // true iff the provided Proof of Spacetime is valid
```

------

## Piece Inclusion Proof

### PieceInclusionProof

A `PieceInclusionProof` contains a potentially complex merkle inclusion proof that all leaves included in `commP` (the piece commitment) are also included in `commD` (the sector data commitment).

```
struct PieceInclusionProof {
    Position uint,
    ProofElements [32]byte
}
```

### GeneratePieceInclusionProofs

`GeneratePieceInclusionProofs` takes a merkle tree and a slice of piece start positions and lengths (in nodes), and returns
a vector of `PieceInclusionProofs` corresponding to the pieces. For this method to work, the piece data used to validate pieces will need to be padded as necessary,
and pieces will need to be aligned (to 128-byte chunks due to details of __*preprocessing*__) when written. This assumes that pieces have been packed and padded according to the assumptions of the algorithm. For this reason, practical implementations should also provide a function to assist in correct packing of pieces.

```
GeneratePieceInclusionProofs
 (
  Tree MerkleTree,
  PieceStarts []uint
  PieceLengths uint,
 ) []PieceInclusionProof
```

`GeneratePieceInclusionProof` takes a merkle tree and the index positions of the first and last nodes
of the piece whose inclusion should be proved. It returns a corresponding `PieceInclusionProof`.
For the resulting proof to be valid, first_node must be <= last_node.

```
GeneratePieceInclusionProof
 (
  tree          MerkleTree,
  firstNode     uint,
  pieceLength   uint,
 ) err Error |  proof PieceInclusionProof
```

`VerifyPieceInclusionProof` takes a sector data commitment (`commD`), piece commitment (`commP`), sector size, and piece size.
Iff it returns true, then `PieceInclusionProof` indeed proves that all of piece's bytes were included in the merkle tree corresponding
to `commD` of a sector of `sectorSize`. The size inputs are necessary to prevent malicious provers from claiming to store the entire
piece but actually storing only the piece commitment. 

```
VerifyPieceInclusionProof
 (
  proof PieceInclusionProof,
  commD  [32]byte,
  commP [32]byte,
  sectorSize uint,
  pieceSize uint,
 ) err Error | IsValid bool // true iff the provided PieceInclusionProof is valid.
```
