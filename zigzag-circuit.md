# ZigZag: Offline PoRep Circuit Spec

ZigZag is the Proof of Replication used in Filecoin. The prover encodes the original data into a replica and commits to it. An offline PoRep proves that the commitment to the replica is a valid commitment of the encoded original data.

ZigZag has been presented by [Ben Fisch at EUROCRYPT19](https://eprint.iacr.org/2018/702.pdf).



#### ZigZag Overview

ZigZag PoRep is based on layering DRG graphs `l` times. The data represented in each DRG layer is the data encoded in the previous layer. The final layer is the replica (which in Filecoin terms is the sealed sector).

- `replicaId` is a unique replica identifier (see the Filecoin Proofs spec for details)
- `CommD` is the Merkle Tree root hash of the input data to the first layer
- `CommR_{l}` is the Merkle Tree hash of the output of the DRG encoding at each layer `l` 
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

## Offline PoRep circuit

**Public Parameters**: *Parameters that are embeded in the circuits or used to generate the circuit*

- `LAYERS : UInt`: Number of layers
- `LAYER_CHALLENGES[0..LAYERS] : UInt`: Number of challenges per layer
- `EXPANSION_DEGREE: UInt`: Degree of the bipartite expander graph to extend dependencies between layers
- `BASE_DEGREE: UInt`: Degree of each Depth Robust Graph
- `TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes))
- `PARENT_NODES : UInt`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`

**Public Inputs**: *Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `ReplicaId : Fr`: A unique identifier for the replica Id
- `CommD : Fr`: the Merkle tree root hash of the original data (input to the first layer)
- `CommRLast : Fr`: the Merkle tree root hash of the output of the last layer 
- `CommRStar : Fr`: the aggregate of each layer's Merkle tree root hash
- Inclusion paths: Binary representation of the Merkle tree path that must be proven packed into a single `Fr` element. We have the following inclusion paths:
  - `InclusionPaths_{l=0..LAYERS}_{0..LAYER_CHALLENGES[l]} : Fr`: At each layer `l` we have `LAYER_CHALLENGES[l]` inclusion paths.
  - `ParentsInclusionPaths_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[l]}_{0..PARENT_NODES} : Fr`: At each layer `l` we have `LAYER_CHALLENGES[l]` an inclusion path for each parent node of the corresponding `InclusionPaths_{l}_{c}`.

**Private Inputs**: *Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `CommR_{i=0..LAYERS-1}`: Commitment of the the encoded data at layer `i`. 
- Inclusion Proof: For each inclusion path in the public inputs, we provide a Merkle Tree path
  - `InclusionHash_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}_{0..TREE_DEPTH-1} : Fr`: Pedersen hashes of the Merkle inclusion proofs of the unencoded challenged nodes at layer `l`
  - `ReplicaInclusionHash_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]}_{0..TREE_DEPTH-1} : Fr`: Pedersen hashes of the Merkle inclusion proofs of the encoded challenged nodes at layer `l`
  - `ParentInclusionHash_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[i]}_{0..PARENT_NODES}_{0..TREE_DEPTH-1} : Fr`: Pedersen hashes of the Merkle inclusion proofs of the parent nodes for each challenged node at layer `l`
  - `InclusionLeaf_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]} : Fr`: Value of the unencoded challenged nodes at layer `l`
  - `ReplicaInclusionLeaf_{i=0..LAYERS}_{0..LAYER_CHALLENGES[i]} : Fr`: Value of the encoded nodes for each challenged node at layer `l`
  - `ParentInclusionLeaf_{l=0..LAYERS}_{c=0..LAYER_CHALLENGES[i]}_{0..PARENT_NODES} : Fr`: Value of the parent nodes for each challenged node at layer `l`

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
        Check MerkleTreeVerify(InclusionHash_{l}_{c}_{0..TREE_DEPTH-1})
        Check MerkleTreeVerify(ReplicaInclusionHash_{l}_{c}_{0..TREE_DEPTH-1})
        
        For p = 0..PARENT_NODES:
        	Check MerkleTreeVerify(ParentInclusionHash_{l}_{c}_{p}_{0..TREE_DEPTH-1})
        ```

      - Correct layer: **Check** that (1) inclusion proofs have their root hash equal to the  commitment of the previous layer and (2) replication and parent inclusion proofs have their root hash equal the commitment of the current layer. 

        *Note on inclusion proofs*: at the first layer, the inclusion proofs root hash must be equal to `CommD`
        *Note on replication proofs*: at the last layer, the replica and parent inclusion proofs must have their root hash equal to `CommRLast`

        *Note on proofs*: the hash at index `0` of an proof is the root of the Merkle tree

        ```
        Assign InclusionRoot = l=0 ? CommD : CommR_{l-1}
        Assign ReplicationInclusionRoot = l=LAYERS-1 ? CommRLast : CommR_{l}
        
        Check InclusionRoot === InclusionHash_{l}_{c}_{0}
        Check ReplicationInclusionRoot === ReplicationInclusionHash_{l}_{c}_{0}
        
        For parent = 0..PARENT_NODES:
        	Check ReplicationInclusionRoot === ParentInclusionHash_{l}_{c}_{parent}_{0}
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
        Assign pre_image = replica_id_bits || ParentBits_{l}_{c}_{0} || .. || ParentBits_{l}_{c}_{PARENT_NODES}
        
        Assign key : Fr = PedersenHash(pre_image)
        Check PedersenHash(pre_image) == key
        ```

      - **Check** correct encoding: the decoded leaf equals the input data leaf

        ```
        Check ReplicationInclusionValue_{l}_{c} == InclusionValue_{l}_{c} + k
        ```

- CommRStar check: **Check** that CommRStar is computed by concatenating the replica identifier and the commR at each layer (specified in the aux inputs)

  ```
  Check CommRStar == PedersenHash(ReplicaId || CommR_{0} || .. || CommR_{LAYERS-2} || CommRLast)
  // TODO check if we need to do packing/unpacking
  ```

**Verification of offline porep proof:**

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Parent checks: For each `node = InclusionPaths_{l}_{c}`:
  - **Check** that all `ParentsInclusionPaths_{l}_{c}_{0..PARENT_NODES}` are the correct parent nodes of `node` in the DRG graph.
  - **Check** that the parent nodes are in numerical order.
