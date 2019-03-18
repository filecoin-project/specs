# ZigZag: Offline PoRep Circuit Spec

ZigZag is the Proof of Replication used in Filecoin. The prover encodes the original data into a replica and commits to it. An offline PoRep proves that the commitment to the replica is a valid commitment of the encoded original data.

#### ZigZag Overview

ZigZag is the Proof of Replication used in Filecoin. The prover encodes the original data into a replica and commits to it. An offline PoRep proves that the commitment to the replica is a valid commitment of the encoded original data.

ZigZag PoRep is based on layering DRG graphs `l` times. The data represented in each DRG layer is the data encoded in the previous layer. The final layer is the replica (which in Filecoin terms is the sealed sector).

- `replicaId` is a unique replica identifier (see the Filecoin Proofs spec for details)
- `CommD` is the Merkle Tree root hash of the input data to the first layer
- `CommR_{i}` is the Merkle Tree hash of the output of the DRG encoding at each layer `i` 
- `CommRStar` is the hash of the concatenation of the `replicaId` and all the `commR`s.

In the PoRep offline proof as presented in Fisch2019, the proof is large for blockchain usage. We use SNARKs to generate a proof of knowledge of an offline proof. In other words, we use encode the ZigZag proof verification algorithm in an arithmetic circuit and use SNARKs to prove that it was evaluated correctly.

This circuit proves that given a Merkle root `CommD`, `CommR_l`, and `commRStar`, that the prover knew the correct replicated data at each layer.

## Offline PoRep circuit

**Public Parameters**: *Parameters that are embeded in the circuits or used to generate the circuit*

- `LAYERS : UInt`: Number of layers
- `LAYER_CHALLENGES[0..LAYERS]`: Number of challenges per layer
- `EXPANSION_DEGREE: UInt`: Degree of the bipartite expander graph to extend dependencies between layers
- `TREE_DEPTH: UInt`: Height of the Merkle tree

**Public Inputs**: *Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `ReplicaId : Fr`: A unique identifier for the replica Id
- `CommD : Fr`: the Merkle tree root hash of the original data (input to the first layer)
- `CommRLast : Fr`: the Merkle tree root hash of the output of the last layer 
- `CommRStar : Fr`: the aggregate of each layer's Merkle tree root hash
- Inclusion paths: Binary representation of the Merkle tree path that must be proven packed into a single `Fr` element. We have the following inclusion paths:
  - `InclusionPaths_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}`: At each layer `i` we have `LAYER_CHALLENGES[i]` inclusion paths.
  - `ParentsInclusionPaths_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[i]}_{EXPANSION_DEGREE+ BASE_DEGREE}`: At each layer `l` we have `LAYER_CHALLENGES[i]` an inclusion path for each parent node of the corresponding `InclusionPaths_{l}_{c}`.

**Private Inputs**: *Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `CommR_{i=0..LAYERS-1}`: Commitment of the the encoded data at layer `i`. 
- Inclusion Proof: For each inclusion path in the public inputs, we provide a Merkle Tree path
  - `InclusionHash_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}_{0..TREE_DEPTH-1}`
  - `ReplicaInclusionHash_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}_{0..TREE_DEPTH-1}`
  - `ParentInclusionHash_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[i]}_{EXPANSION_DEGREE+ BASE_DEGREE}_{0..TREE_DEPTH-1}`
  - `InclusionLeaf_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}`
  - `ReplicaInclusionLeaf_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}`
  - `ParentInclusionLeaf_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[i]}_{EXPANSION_DEGREE+ BASE_DEGREE}`

**Circuit:**

- **Check** `ReplicaId` is equal to its bit representation

  ```
  Assign replica_id_bits = Fr_to_bits(ReplicaId)
  Check Packed(replica_id_bits) == ReplicaId
  ```

- For each `l = 0..LAYERS`:

  - For each `c = 0..LAYERS_CHALLENGES[l]`

    - Inclusion checks:

      - Correct inclusions proofs: **Check**  that all the inclusion proofs are correct

        ```
        Check MerkleTreeVerify(InclusionHash_{l}_{c}_{0..TREE_DEPTH})
        Check MerkleTreeVerify(ReplicaInclusionHash_{l}_{c}_{0..TREE_DEPTH})
        
        For p = 0..EXPANSION_DEGREE + BASE_DEGREE:
        	Check MerkleTreeVerify(ParentInclusionHash_{l}_{c}_{p}_{0..TREE_DEPTH})
        ```

      - Correct layer: **Check** that `CommR_{l}` is matching the Replica Inclusion proofs root hash and `CommR_{l-1}` is matching the Inclusion proof root hash. (If `l=0`, use `CommD` instead, if `l=LAYER-1` use `CommRLast` instead).

        ```
        Check CommR_{l-1} === InclusionHash_{l}_{c}_{0}
        Check CommR_{l} === ReplicationInclusionHash_{l}_{c}_{0}
        
        For parent = 0..EXPANSION_DEGREE + BASE_DEGREE:
        	Check CommR_{l} === ParentInclusionHash_{l}_{c}_{parent}_{0}
        ```

    - Encoding checks: **Check** that a challenged replica node decodes to the correct data node.

      - **Check** that each parent has a correct bit representation:

        ```
        For each parent = ParentInclusionLeaf_{l}_{c}_{p}`:
          Assign ParentBits_{l}_{c}_{p} = Fr_to_bits(parent)
          Check Packed(ParentBits_{l}_{c}_{p}) == parent
        ```

      - **Check** that the KDF was run correctly:

        ```
        Assign leaf_bits = ParentBits_{l}_{c}_{0..EXPANSION_DEGREE+BASE_DEGREE}
        
        Assign key : Fr = KDF(replica_id_bits, leaf_bits)
        Check KDF(replica_id_bits, leaf_bits)
        ```

      - **Check** correct encoding: the decoded leaf equals the input data leaf

        ```
        Check ReplicationInclusionValue_{l}_{c} == InclusionValue_{l}_{c} + k
        ```

- CommRStar check: **Check** that CommRStar is computed by concatenating the replica identifier and the commR at each layer (specified in the aux inputs)

  ```
  Check CommRStar == PedersenHash(ReplicaId || CommR_{0} || .. || CommR_{LAYERS-1} || CommRLast)
  // TODO check if we need to do packing/unpacking
  ```

**Verification of offline porep proof:**

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Parent checks: For each `node = InclusionPaths_{l}_{c}`:
  - **Check** that all `ParentsInclusionPaths_{l}_{c}_{0..EXPANSION_DEGREE+BASE_DEGREE}` are the correct parent nodes of `node` in the DRG graph.
  - **Check** that the parent nodes are in numerical order.
