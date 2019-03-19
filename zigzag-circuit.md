# ZigZag: Offline PoRep Circuit Spec

ZigZag is the Proof of Replication used in Filecoin. The prover encodes the original data into a replica and commits to it. An offline PoRep proves that the commitment to the replica is a valid commitment of the encoded original data.

ZigZag has been presented by [Ben Fisch at EUROCRYPT19](https://eprint.iacr.org/2018/702.pdf).

#### ZigZag Overview

ZigZag PoRep is based on layering DRG graphs `l` times. The data represented in each DRG layer is the data encoded in the previous layer. The final layer is the replica (which in Filecoin terms is the sealed sector).

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
- `TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes)).
- `PARENT_NODES : UInt`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`.

#### Public Inputs

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `ReplicaId : Fr`: A unique identifier for the replica.

- `CommD : Fr`: the Merkle tree root hash of the original data (input to the first layer).

- `CommRLast : Fr`: The Merkle tree root hash of the final replica (output of the last layer).

- `CommRStar : Fr`: A commitment to each `l` layer's Merkle tree root hash `CommR[l]` and `ReplicaId`.

- `InclusionPath : [LAYERS][]Fr`: Inclusion path for the challenged data and replica node.

  Note: Each layer `l` has `LAYER_CHALLENGES[l]` inclusion paths.

- `ParentInclusionPath : [LAYERS][][PARENT_NODES]Fr`:  Inclusion path for the parent nodes of the corresponding `InclusionPath[l][c]` nodes.

  Note: Each layer `l` has `LAYER_CHALLENGES[l]` inclusion paths.

##### Design notes

- `CommRLast` is a public input, since we will be using it during Proof-of-Spacetime

#### Private Inputs

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `CommR : [LAYERS-1]Fr`: Commitment of the the encoded data at each layer. 

  Note: Size is `LAYERS-1` since the commitment to the last layer is `CommRLast`

- `DataProof : [LAYERS][][TREE_DEPTH-2]Fr`: Merkle tree inclusion proof for the current layer unencoded challenged nodes.

  Note: Size of proof per layer is `TREE_DEPTH-2` because it excludes (1) the roothash, and (2) the leaf value: (1) the root hash of the proof will be the corresponding unencoded data commitment for the current layer `l` (either `CommD` or `CommR[l-1]`), (2) the leaf value is `DataValue[l][c]`

- `ReplicaProof : [LAYERS][][TREE_DEPTH-2]Fr`: Merkle tree inclusion proof for the current layer encoded challenged nodes.

  Note: Size of proof per layer is `TREE_DEPTH-2` because it excludes (1) the roothash, and (2) the leaf value: (1) the root hash of the proof will be the corresponding unencoded data commitment for the current layer  (either `CommR[l]` or `CommRLast`), (2) the leaf value is `ReplicaValue[l][c]`

- `ParentProof : [LAYERS][][PARENT_NODES][TREE_DEPTH-2]Fr`: Pedersen hashes of the Merkle inclusion proofs of the parent nodes for each challenged node at layer `l`.

  Note: Size of proof per layer is `TREE_DEPTH-2` because it excludes (1) the roothash, and (2) the leaf value: (1) the root hash of the proof will be the corresponding unencoded data commitment for the current layer  (either `CommR[l]` or `CommRLast`), (2) the leaf value is `ParentValue[l][c][p]`

- `DataValue : [LAYERS][]Fr`: Value of the unencoded challenged nodes at layer `l`.

- `ReplicaValue : [LAYERS][]Fr`: Value of the encoded nodes for each challenged node at layer `l`.

- `ParentValue : [LAYERS][][PARENT_NODES]Fr`: Value of the parent nodes for each challenged node at layer `l`.

##### Design notes

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
Assign ReplicaIdBits : [255]Fr = Fr_to_bits(ReplicaId)
Check Packed(replica_id_bits) == ReplicaId

For l = 0..LAYERS:
  For c = 0..LAYERS_CHALLENGES[l]:

    Assign DataRoot : Fr = l == 0 ? CommD : CommR[l-1]
    Assign ReplicaRoot : Fr = l == LAYERS-1 ? CommRLast : CommR[l]

    // 2: Inclusion Proofs Checks
    // 2.1: Check inclusion proofs for data nodes are correct
    Check MerkleTreeVerify(DataRoot, InclusionPath[l][c], DataProof[l][c], DataValue[l][c])
    // 2.2: Check inclusion proofs for replica nodes are correct
    Check MerkleTreeVerify(ReplicaRoot, InclusionPath[l][c], ReplicaProof[l][c], ReplicaValue[l][c])
    // 2.3: Check inclusion proofs for parent nodes are correct
    For p = 0..PARENT_NODES:
      Check MerkleTreeVerify(ReplicaRoot, ParentInclusionPath[l][c][p], ParentProof[l][c][p])

    // 3: Encoding checks - Check that replica nodes have been correctly encoded
    For p = 0..PARENT_NODES:
      // 3.1: Check that each ParentValue is equal to its bit representation
      Assign parent = ParentValue[l][c][p]
      Assign ParentBits[l][c][p] : [255]Fr = Fr_to_bits(parent)
      Check Packed(ParentBits[l][c][p]) == parent

		// 3.2: Check that each key has generated correctly
    Assign PreImage = ReplicaIdBits || ParentBits[l][c][0] || .. || ParentBits[l][c][PARENT_NODES-1]
    Assign key : Fr = PedersenHash(PreImage)
    Check PedersenHash(PreImage) == key
    // 3.3: Check that the data has been encoded to a replica with the right key
    Check ReplicationInclusionValue[l][c] == InclusionValue[l][c] + k

    // 4: CommRStar check - Check that the CommRStar constructed correctly
    Check CommRStar == PedersenHash(ReplicaId || CommR[0] || .. || CommR[LAYERS-2] || CommRLast)
		// TODO check if we need to do packing/unpacking
```



#### Verification of offline porep proof

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Parent checks: For each `node = InclusionPaths_{l}_{c}`:
  - **Check** that all `ParentsInclusionPaths_{l}_{c}_{0..PARENT_NODES}` are the correct parent nodes of `node` in the DRG graph.
  - **Check** that the parent nodes are in numerical order.

