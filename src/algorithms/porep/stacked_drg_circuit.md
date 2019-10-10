# Stacked DRG: Offline PoRep Circuit Spec

Stacked DRG PoRep is the Proof of Replication used in Filecoin. The prover encodes the original data into a replica and commits to it. An offline PoRep proves that the commitment to the replica is a valid commitment of the encoded original data.

Stacked DRG PoRep has been presented by [Ben Fisch at EUROCRYPT19](https://eprint.iacr.org/2018/702.pdf).

#### Stacked DRG PoRep Overview

Stacked DRG PoRep is based on layering DRG graphs `LAYERS` times. The data represented in each DRG layer is the data encoded in the previous layer. The final layer is the replica (which in Filecoin terms is the sealed sector).

- `ReplicaId` is a unique replica identifier (see the Filecoin Proofs spec for details)
- `CommD` is the Merkle Tree root hash of the input data to the first layer
- `CommR[l]` is the Merkle Tree hash of the output of the DRG encoding at each layer `l` 
- `CommRStar` is the hash of the concatenation of the `ReplicaId` and all the `CommR`s.

The (offline) proof size in the ZigZag is too large for blockchain usage (~3MB). We use SNARKs to generate a proof of knowledge of a correct ZigZag proof. In other words, we implement the ZigZag proof verification algorithm in an arithmetic circuit and use SNARKs to prove that it was evaluated correctly.

This circuit proves that given a Merkle root `CommD`, `CommRLast`, and `commRStar`, that the prover knew the correct replicated data at each layer.

#### Spec notation

- **Fr**: Field element of BLS12-381
- **UInt**: Unsigned integer
- **{0..x}**: From `0` (included) to `x` (not included) (e.g. `[0,x)` )
- **Check**: 
  - If there is an equality, create a constraint
  - otherwise, execute the function
- **Inclusion path**: Binary representation of the Merkle tree path that must be proven packed into a single `Fr` element.

## Offline PoRep circuit

#### Public Parameters

*Parameters that are embeded in the circuits or used to generate the circuit*

- `LAYERS : UInt`: Number of DRG layers.
- `LAYER_CHALLENGES : [LAYERS]UInt`: Number of challenges per layer.
- `EXPANSION_DEGREE: UInt`: Degree of each bipartite expander graph to extend dependencies between layers.
- `BASE_DEGREE: UInt`: Degree of each Depth Robust Graph.
- `TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes/32 bytes per leaf)).
- `PARENT_COUNT : UInt`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`.

#### Public Inputs

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `ReplicaId : Fr`: A unique identifier for the replica.
- `CommD : Fr`: the Merkle tree root hash of the original data (input to the first layer).
- `CommRLast : Fr`: The Merkle tree root hash of the final replica (output of the last layer).
- `CommRStar : Fr`: A commitment to each `l` layer's Merkle tree root hash `CommR[l]` and `ReplicaId`.
- `InclusionPath : [LAYERS][]Fr`: Inclusion path for the challenged data and replica leaf.
- `ParentInclusionPath : [LAYERS][][PARENT_COUNT]Fr`:  Inclusion path for the parents of the corresponding `InclusionPath[l][c]`.

##### Design notes

- `CommRLast` is a public input, since we will be using it during Proof-of-Spacetime.
- `InclusionPath` and `ParentInclusionPath`: Each layer `l` has `LAYER_CHALLENGES[l]` inclusion paths.

#### Private Inputs

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `CommR : [LAYERS-1]Fr`: Commitment of the the encoded data at each layer. 

  Note: Size is `LAYERS-1` since the commitment to the last layer is `CommRLast`

- `DataProof : [LAYERS][][TREE_DEPTH]Fr`: Merkle tree inclusion proof for the current layer unencoded challenged leaf.

- `ReplicaProof : [LAYERS][][TREE_DEPTH]Fr`: Merkle tree inclusion proof for the current layer encoded challenged leaves.

- `ParentProof : [LAYERS][][PARENT_COUNT][TREE_DEPTH]Fr`: Pedersen hashes of the Merkle inclusion proofs of the parent leaves for each challenged leaf at layer `l`.

- `DataValue : [LAYERS][]Fr`: Value of the unencoded challenged leaves at layer `l`.

- `ReplicaValue : [LAYERS][]Fr`: Value of the encoded leaves for each challenged leaf at layer `l`.

- `ParentValue : [LAYERS][][PARENT_COUNT]Fr`: Value of the parent leaves for each challenged leaf at layer `l`.

#### Circuit

##### High Level

In high level, we do 4 checks:

1. **ReplicaId Check**: Check the binary representation of the ReplicaId
2. **Inclusion Proofs Checks**: Check the inclusion proofs
3. **Encoding Checks**: Check that the data has been correctly encoding into a replica
4. **CommRStar Check**: Check that CommRStar has been generated correctly

##### Details

```go
// 1: ReplicaId Check - Check ReplicaId is equal to its bit representation
let ReplicaIdBits : [255]Fr = Fr_to_bits(ReplicaId)
assert(Packed(replica_id_bits) == ReplicaId)

let DataRoot, ReplicaRoot Fr

for l in range LAYERS {
  
  if l == 0 {
    DataRoot = CommD
  } else {
    DataRoot = CommR[l-1]
  }

  if l == LAYERS-1 {
    ReplicaRoot = CommRLast
  } else {
    ReplicaRoot = CommR[l]
  }
  
  for c in range LAYERS_CHALLENGES[l] {
    // 2: Inclusion Proofs Checks
    // 2.1: Check inclusion proofs for data leaves are correct
    assert(MerkleTreeVerify(DataRoot, InclusionPath[l][c], DataProof[l][c], DataValue[l][c]))
    // 2.2: Check inclusion proofs for replica leaves are correct
    assert(MerkleTreeVerify(ReplicaRoot, InclusionPath[l][c], ReplicaProof[l][c], ReplicaValue[l][c]))
    // 2.3: Check inclusion proofs for parent leaves are correct
    for p in range PARENT_COUNT {
      assert(MerkleTreeVerify(ReplicaRoot, ParentInclusionPath[l][c][p], ParentProof[l][c][p]))
    }

    // 3: Encoding checks - Check that replica leaves have been correctly encoded
    let ParentBits [PARENT_COUNT][255]Fr
    for p in range PARENT_COUNT {
      // 3.1: Check that each ParentValue is equal to its bit representation
      let parent = ParentValue[l][c][p]
      ParentBits[p] = Fr_to_bits(parent)
      assert(Packed(ParentBits[p]) == parent)
    }

    // 3.2: KDF check - Check that each key has generated correctly
    // PreImage = ReplicaIdBits || ParentBits[1] .. ParentBits[PARENT_NODES]
    let PreImage = ReplicaIdBits
    for parentbits in ParentBits {
      PreImage.Append(parentbits)
    }
    let key Fr = Blake2s(PreImage)
    assert(Blake2s(PreImage) == key)

    // 3.3: Check that the data has been encoded to a replica with the right key
    assert(ReplicaValue[l][c] == DataValue[l][c] + key)
  }
  
  // 4: CommRStar check - Check that the CommRStar constructed correctly
  let hash = ReplicaId
  for l in range LAYERS-1 {
    hash.Append(CommR[l])
  }
  hash.Append(CommRLast)

  assert(CommRStar == PedersenHash(hash))
  // TODO check if we need to do packing/unpacking
}
```



#### Verification of offline porep proof

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Parent checks: For each `leaf = InclusionPath[l][c]`:
  - **Check** that all `ParentsInclusionPaths_[l][c][0..PARENT_COUNT}` are the correct parent leaves of `leaf` in the DRG graph, if a leaf has less than `PARENT_COUNT`, repeat the leaf with the highest label in the graph.
  - **Check** that the parent leaves are in ascending numerical order.

