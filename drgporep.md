# DRGPorep Circuit Design Doc

## Offline Proof

A Prover wants to convince a verifier that they have run the DRGPoRep replication algorithm on some data with Merkle root `comm_d` and that the Merkle root of the replicated data is `comm_r`.

#### High level description

Circuit inputs:

- (Public input) Commitments to the original data and the replica data
- (Public input) An identifier for the replica
- (Public input) The binary representation of the position in the Merkle tree of:
  - the challenged nodes in the original data
  - the challenged nodes in the replica data
  - the parents of the challenged nodes in the replica data
- (Aux input) Inclusion proofs for the set of challenged leaves in the original data and the replica data
- (Aux input) Inclusion proofs for the replica parent leaves of each challenged replica data leaf
- (Aux input) An identifier for the replica

Outside of circuit checks:

- Parent checks: **Check** that for each challenged data node, the parent replica nodes have the correct position (expressed as binary representation in the public inputs)
- Unique challenges checks:  **Check** that the challenged nodes are all different

Inside of circuit checks:

- Inclusion checks: **Check** that all the inclusion proofs are correct
- Encoding checks: **Check** that we can derive a decoding of the replica data that equals the original data
  - **Compute** a key: Hash the concatenatation the replica identifier with the replica parent leaves in dependency order in the DRG graph
  - **Compute** a Sloth decoding of the replica data using the key 
  - **Check** that the decoded leaf equals the original data leaf

#### Low level description

- **Public Parameters**: *Parameters that are embeded in the circuits or used to generate the circuit*
  - `jubjub_params: TODO`: TODO find correct structure, default: TODO
  - `degree: UInt`: Degree of the Depth Robust Graph
  - `sloth_iterations: UInt`: Number of iteration for Sloth enocding/decoding
  - `challenge_count: UInt`: Number of challenges
  - `merkle_tree_depth: UInt`: Depth of the Merkle Trees in the inclusion proofs

- **Public Inputs**: *Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*
  - `replica_id: Fr`  Unique identity for the replica (in Filecoin `H(replica_id, prover_id)`)
  - `comm_d: Fr`: Merkle root of the original data 
  - `comm_r: Fr`: Merkle root of the replicated data
  - `challenge_{0..challenge_count}/inclusion_checks`
    - `replica_inclusion/path/input: Fr`: Packed boolean vector that represents the authentication path for the replica inclusion proof; bool says if path is right (1) or left (0).
    - `data_inclusion/path/input: Fr`: Same previous one for data inclusion proofs
    - `parents_inclusion_{0..degree}/path/input: Fr`: Same previous one for data inclusion proofs

- **Private Inputs**: *Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*
  - AUX (input)
    - `challenge_{0..challenge_count}`
      - `inclusion_checks`
        - `replica_inclusion`
          - `value/num : Fr`: Leaf of the Merkle tree
          - `merkle tree hash {0..merkle_tree_depth} : TODO`: intermediary hashes for the merkle tree inclusion proof
        - `parent_inclusion_{0..degree}`
          - `value/num : Fr`: Leaf of the Merkle tree
          - `merkle tree hash {0..merkle_tree_depth} : TODO`: intermediary hashes for the merkle tree inclusion proof
        - `data_inclusion`
          - `value/num : Fr`: Leaf of the Merkle tree
          - `merkle tree hash {0..merkle_tree_depth} : TODO`: intermediary hashes for the merkle tree inclusion proof
      - `encoding_checks/parents_{0..degree}_bits/bit {0..255}: Fr`: Bit representation of the parent hashes
  - AUX (computed)
    - `encoding_checks/kdf: TODO`: TODO
    - `replica_id_bits/bit {0..255}/boolean: Fr`: bit at position *i*

- **Circuit `drgporep`**:

  - `replica_id_bits`: *Check `replica_id` is equal to its bit representation*

    ```
    Assign replica_id_bits = Fr_to_bits(replica_id)
    Check Packed(replica_id_bits) == replica_id
    ```

  - `challenge_{chall = 0..challenge_count}`

    - `inclusion_checks`: *Check inclusion proofs*

      - `replica_inclusion`: *Check inclusion for the challenged replica node*

        ```
        Check PoR(jubjub_params, replica_inclusion_proofs[chall], comm_r)
        ```

      - `data_inclusion`: *Check inclusion proof for the challenged data node*

        ```
        Check PoR(jubjub_params, data_inclusion_proofs[chall], comm_d)
        ```

      - `parents_inclusion_{parent = 0..(degree+expansion_degree)}`:  *Check inclusion proof for each parent of the challenged data node*

        ```
        Check PoR(jubjub_params, parents_inclusion_proofs[chall][parent], comm_r)
        ```

    - `encoding_checks`: *Check a data node was encoded correctly in replica node*

      - `parent_bits = parents_{parent = 0..degree}_bits`: *Check that a correct bit representation of the parents is known*

        ```
        Let leaf = /challenge_{chall}/inclusion_checks/parents_inclusion_{parent}/value
        Assign parent_bits = Fr_to_bits(leaf)
        Check Packed(bits) == leaf
        ```

      - `kdf`: *Check that the KDF was run correctly*

        ```
        Let leaf_bits[i = 0..deg] = challenge_{chall}/encoding_checks/parents_{i}_bits
        
        Assign key : Fr = KDF(replica_id_bits, leaf_bits)
        Check KDF(replica_id_bits, leaf_bits)
        ```

      - `sloth_decode`: *Check that the Sloth encoding was run correctly*

        ```
        Let leaf = challenge_{chall}/inclusion_checks/replica_inclusion/value
        
        Assign decoded = SlothDecode(key, leaf, sloth_iterations)
        Check SlothDecode(key, leaf, sloth_iterations)
        ```

      - `equality`: *Check that the decoded piece is equivalent to the challenged node*

        ```
        Let leaf = /challenge_{chall}/inclusion_checks/data_inclusion/value
        
        Check leaf == decoded
        ```



## Glossary

- **Fr**: Field element of BLS12-381
- **Merkle root**: Root hash of a binary Merkle tree
- **UInt**: Unsigned integer
- **{0..x}**: From 0 (included) to x (not included) (e.g. [0,x)] )
- **Check**: 
  - If there is an equality, create a constraint
  - otherwise, execute the function
