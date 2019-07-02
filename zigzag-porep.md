# Proof-of-Replication in Filecoin


Filecoin Proof-of-Replication (PoRep) 

This spec describes the specific Proof-of-Replication used in Filecoin called *ZigZag*. 

ZigZag has been presented by [Ben Fisch at EUROCRYPT19](https://eprint.iacr.org/2018/702.pdf).

## Introduction

### Background on Proof-of-Replication
*Proof-of-Replication* enables a prover *P* to convince a verifier *V* that *P* is storing a replica *R*, a physically independent copy of some data *D*, unique to *P*. The scheme is defined by a tuple of polynomial time algorithms (_Setup_, Replication, _Prove_, _Verify_). The assumption is that generation of a replica after _Replicate_  must be difficult (if not impossible) to generate.

- *Setup*: On setup, the public parameters of the proving systems are set.
- *Replicate*: On replication, either a party or both (depending on the scheme, in our case the prover only!) generate a unique permutation of the original data _D_, which we call replica _R_.
- _Prove_: On receiving a challenge, the prover must generate a proof that it is in possession of the replica and that it was derived from data _D_. The prover must only be able to respond to the challenge successfully if it is in possession of the replica, since would be difficult (if not impossible) to generate a replica that can be used to generate the proof at this stage
- _Verify_: On receiving the proof, the verifier checks the validity of the proof and accepts or rejects the proof. 

{{% mermaid %}}
sequenceDiagram
    Note right of Prover: CommD
    Prover-->>Prover: R, CommR ← Replicate(D) 
    Prover->>Verifier: CommR
    Verifier-->>Verifier: Generate random challenge
    Verifier->>Prover: challenge
    Prover-->>Prover: proof ← Prove(D, R, challenge)
    Prover->>Verifier: proof
{{% /mermaid %}}

#### Time-bounded Proof-of-Replication

**Timing assumption**. *Time-bounded Proof-of-Replication* are constructions of PoRep with timing assumptions. The assumption is that generation of the replica (hence the _Replication_) takes some time _t_ that is substantially larger than the time it takes to produce a proof (hence _time(Prove)_) and the round-trip time (_RTT_) for sending a challenge and receiving a proof.

**Distinguishing Malicious provers**. A malicious prover that does not have _R_, must obtain it (or generate it), before the _Prove_ step. A verifier can distinguish an honest prover from a malicious prover, since the malicious one will take too long to answer the challenge. A verifier will reject if receiving the proof from the prover takes longer than a timeout (bounded between proving time and replication time).

### Background on ZigZag

*ZigZag* is a specific Proof-of-Replication construction that we use in Filecoin. ZigZag has been designed by [Ben Fisch at EUROCRYPT19](https://eprint.iacr.org/2018/702.pdf).  In high level, ZigZag ensures that the *Replicate* step is a slow non-parallelizable sequential process by using some special type of graphs called Depth Robust Graphs (we refer to as "DRG").

**Encoding using DRGs**. By positioning data blocks into nodes in the DRG, we sequentially encode each node in the graph using its encoded parents. DRG graphs ensure that this process is not likely to be parallelizable.

**Layering DRGs**. ZigZag repeates this encoding by layering DRG graphs `LAYERS` times. The data represented in each DRG layer is the data encoded in the previous layer. We connect different layers using Bipartite Expander Graphs and at each layer, we reverse the graph edges with a technique which we call zigzag. The combination of DRGs, expander graphs and zigzag guarantee the security property of PoRep. The final encoded layer is the final replica

**Generating ZigZag proofs**. Given the following public parameters:

- `ReplicaId` is a unique replica identifier (see the Filecoin Proofs spec for details)
- `CommD` is the Merkle Tree root hash of the input data to the first layer
- `CommRStar` is the hash of the concatenation of the `ReplicaId` and all the `CommR`s.
- `CommRLast` is the hash of the last encoded DRG layer.

A ZigZag proof proves that some data whose committment is `CommD` has been used to run a `Replicate` algorithm and generated some data whose commitment is `CommRLast`.

A ZigZag proof consists of a set of challenged DRG nodes (both encoded and unencoded) for each layer, a set of parent nodes for each challenged node and a Merkle tree inclusion proof for each node provided. The verifier can then verify the correct encoding of each node and that the nodes given were consistent with the provers' commitments.

**Making proofs succinct with SNARKs**: The proof size in the ZigZag is too large for blockchain usage (~100MB), mostly due to the large amount of Merkle tree inclusion proofs required to achieve security. We use SNARKs to generate a proof of knowledge of a correct ZigZag proof. In other words, we implement the ZigZag proof verification algorithm in an arithmetic circuit and use SNARKs to prove that it was evaluated correctly.

The SNARK circuit proves that given a Merkle root `CommD`, `CommRLast`, and `commRStar`, that the prover knew the correct replicated data at each layer.

### ZigZag in Filecoin

Proof-of-Replication proves that a Storage Miner is dedicating unique dedicated storage for each ***sector***. Filecoin Storage Miners collect new clients' data in a sector, run a slow encoding process (called `Seal`) and generate a proof (`SealProof`) that the encoding was generated correctly.

In Filecoin, PoRep provides two guarantees: (1) *space-hardness*: Storage Miners cannot lie about the amount of space they are dedicating to Filecoin in order to gain more power in the consensus; (2) *replication*: Storage Miners are dedicating unique storage for each copy of their clients data. 

Glossary:

- __*sector:*__ a fixed-size block of data of `SECTOR_SIZE` bytes which generally contains clients' data.
- __*raw sector:*__ a concrete representation (on disk or in memory) of a sector's
- __*unsealed sector:*__ a concrete representation (on disk or in memory) of a sector's that follows the "Storage Format" described in [Filecoin Client Data Processing](client-data.md#storage-format) (currently `paddedfr32v1` is the required default).
- __*sealed sector:*__  a concrete representation (on disk or in memory) of the unique replica generated by `Seal` from an __*unsealed sector*__. 
- __*piece:*__ a block of data of at most `SECTOR_SIZE` bytes which is generally is a client's file or part of.

## Public Parameters

*The following public parameters are used in the ZigZag Replication and Proof Generation algorithms:*

- `LAYERS : uint`: Number of Depth Robust Graphs stacked layers.
- `EXPANSION_DEGREE: uint`: Degree of each Bipartite Expander Graph to extend dependencies between layers.
- `BASE_DEGREE: uint`: In-Degree of each Depth Robust Graph.
- `TREE_DEPTH: uint`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes/32 bytes per leaf)).
- `PARENT_COUNT : uint`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`.
- `GRAPH_SIZE: uint`: Number of nodes in the DRG.
- `GRAPH_SEED: uint`: Seed used for random number generation in `baseParents`.
- `NODE_SIZE: uint`: Size of each node in bytes.

*The following additional public parameters are required:*

- `TAPER` : Float: Fraction of each layer's challenges by which to reduce next-lowest layer's challenge count.
- `TAPER_LAYERS`: uint: Number of layers
  `Data` is a byte array initialized to the content of __*unsealed sector*__ and will be mutated in-place by the replication process.

### Hash Functions

We have describe three hash functions:

  - `KDFHash`: a hash function with 32-byte digest size: default is Blake2s-256
  - `CommRHash`: a hash function with 32-byte digest size: default is Blake2s-256
  - `RepCompress`: a hash function with 32-byte digest size: default is Pedersen Hashing over the Jubjub curve.

#### RepHash

`RepHash` is the process to generate a Merkle tree root hash of sealed sectors, unsealed sectors and of the intermediate steps of Proof-of-Replication. It takes as input some data respecting a valid Storage Format and outputs a Merkle root hash. 

`RepHash` is constructed from a balanced binary Merkle tree. The leaves of the merkle tree is the output of `RepCompress` on two adjacent 32 bytes blocks.

```go
type node [32]uint8

// Create and return a binary merkle tree and its root commitment.
// len(leaves) must be a power of 2.
func RepHash(leaves []node) ([][]node, node) {
	rows = [][]node
	
	currentRow := leaves
	for height := 0; len(currentRow) > 1; height += 1 {
		rows.push(currentRow)
		var nextRow []node

		for i := 0; i < len(row)/2; i += 2 {
			left := row[i]
			right := row[i+1]

			// NOTE: Depending on choice of RepCompress, heightPart may be trimmed to fewer than 8 bits.
			heightPart := []uint8{height}
			
			input1 := append(heightPart, left...)
			input := append(input1, right...)
			hashed = RepCompress(input)
			nextRow = append(nextRow, hashed)
		}

		currentRow = nextRow
	}
    // The tree returned here is just a vector of rows for later use. Its representation is not part of the spec.
	return rows, currentRow[0]
}

```

### Depth Robust Graph Generation

```go
func baseParents(node uint) (parents [BASE_DEGREE]uint) {
    switch node {
        // Special case for the first node, it self references.
        // Special case for the second node, it references only the first one.
        case 0:
        case 1:
            for i := 0; i < BASE_DEGREE; i++ {
                parents[i] = 0
            }
        default:
            rng := ChaChaRng.from_seed(GRAPH_SEED)

            for k := 0; k < BASE_DEGREE; k++ {
                // iterate over m meta nodes of the ith real node
                // simulate the edges that we would add from previous graph nodes
                // if any edge is added from a meta node of jth real node then add edge (j,i)
                logi := floor(log2(node * BASE_DEGREE))
                j := rng.gen() % logi
                jj := min(node * BASE_DEGREE + k, 1 << (j + 1))
                backDist := rng.gen_range(max(jj >> 1, 2), jj + 1)
                out := (node * BASE_DEGREE + k - backDist) / BASE_DEGREE

                // remove self references and replace with reference to previous node
                if out == node {
                    parents[i] = node - 1
                } else {
                    parents[i] = out;
                }
            }

            sort(parents)
    }
}

func expandedParents(node uint) (parents []uint) {
    parents := make([]uint, EXPANSION_DEGREE)

    feistelKeys := []uint{1, 2, 3, 4}
    for i := 0, p := 0; i < EXPANSION_DEGREE; i++ {
        a := node * EXPANSION_DEGREE + i
        if graphIsReversed() {
            transformed := feistelInvertPermute(GRAPH_SIZE * EXPANSION_DEGREE, a, feistelKeys)
        } else {
            transformed := feistelPermute(GRAPH_SIZE * EXPANSION_DEGREE, a, feistelKeys)
        }
        other := transformed / EXPANSION_DEGREE
        // Collapse the output in the matrix search space to the row of the corresponding
        // node (losing the column information, that will be regenerated later when calling
        // back this function in the `reversed` direction).

        if graphIsReversed() {
            if other > node {
                parents[p] = other
                p += 1
            }
        } else if other < node {
            parents[p] = other
            p += 1
        }
}

func parents(rawNode uint) (parents [PARENT_COUNT]uint) {
    baseParents := baseParents(realIndex(rawNode))
    for i := 0; i < BASE_DEGREE; i++ {
        parents[i] = realIndex(baseParents[i])
    }

    // expandedParents takes raw_node
    expandedParents := expandedParents(rawNode)
    for i := 0; i < len(expandedParents); i++ {
        parents[BASE_DEGREE + i] = expandedParents[i]
    }

    // Pad so all nodes have correct degree.
    currentLength := BASE_DEGREE + len(expandedParents)
    for i := 0; i < PARENT_COUNT - currentLength; i++ {
        if graphIsReversed() {
            parents[currentLength + i] = GRAPH_SIZE - 1
        } else {
            parents[currentLength + i] = 0
        }
    }
    assert(len(parents) == PARENT_COUNT)
    sort(parents)

    for _, parent := range parents {
        if graphIsReversed() {
            assert(parent >= rawNode)
        } else {
            assert(parent <= rawNode)
        }
    }
}

func realIndex(rawNode uint) uint {
    if graphIsReversed() {
        return GRAPH_SIZE - 1 - rawNode
    } else {
        return rawNode
    }
}
```

### Layer Challenge Counts

TODO: define `Challenge` (or find existing definition)

```go
func ChallengesForLayer(challenge Challenge, layer uint) uint {

    switch challenge.Type {
      // TODO: remove ambiguity, there should not be a "fixed" case
        case Fixed:
            return challenge.Count
        case Tapered:
      // TODO: current calculation is incorrect and does not match claim6 from Fisch2019, it should look more like: https://observablehq.com/d/bbabac1947b79011#gen_zigzag_taper
            assert(layer < LAYERS)
            l := (LAYERS - 1) - layer
            r := 1.0 - TAPER;
            t := min(l, TAPER_LAYERS)

            totalTaper := pow(r, t)

            calculated := ceil(totalTaper * challenge.count)

            // Although implied by the call to `ceil()` above, be explicit
            // that a layer cannot contain 0 challenges.
            max(1, calculated)
    }
}
```

## Replication

> The Replication phase turns an *unsealed sector* into a *sealed sector*

Before running the `Replicate` algorithm, the prover must ensure that the sector is correctly formatted with a valid with the "Storage Format" described in [Filecoin Client Data Processing](client-data.md#storage-format) (currently `paddedfr32v1` is the required default).

The Replication Algorithm  proceeds as follows:

- Calculate `ReplicaID` using `Hash` (Blake2s):
```
ReplicaID := Hash(ProverID || SectorID || ticket)
```
- Perform `RepHash` on `Data` to yield `CommD` and `TreeD`:
```
CommD, TreeD = RepHash(data)
```

For each of `LAYERS` layers, `l`, perform one __*Layer Replication*__, yielding a replica, tree, and commitment (`CommR_<l>`) per layer:

```go
let layer_replicas = [LAYERS][nodes]uint8
let layer_trees = [LAYERS]MerkleTree
let CommRs = []commitment

let layer = data
for l in 0..layers {
	let layer_replica = ReplicateLayer(layer)
	layer_replicas[l] = layer_replica
	CommRs[l], layers_trees[l] = RepTree(layer_replica)
	layer = layer_replica
}
```

The replicated data is the output of the final __*Layer Replication*__,`layer_replicas[layers-1]`.
Set `CommRLast` to be  `CommR_<Layers>`.
Set `CommRStar` to be `CommRHash(ReplicaID || CommR_0 || CommR_<i> || ... || CommRLast)`.

```go
Replica := layer_replicas[layers - 1]
CommRLast :- CommRs[layers-1]
CommRStar := CommRHash(replicaID, ...CommRs)
```

### Layer Replication

TODO: Define `Graph`. We need to decide if this is an object we'll explicitly define or if its properties (e.g., `GRAPH_SIZE`) are just part of the replication parameters and all the functions just refer to the _same_ graphs being manipulated across the entire replication process. (At the moment I've avoided defining a `Graph` structure as in other specs I didn't see any object methods, just standalone functions.)

```go
func transformAndReplicateLayers(slothIter uint, replicaId Domain,
    data []byte) ([]Domain, []Tree) {

    var taus [LAYERS]Domain
    var auxs [LAYERS]Tree
    var sortedTrees [LAYERS]Tree

    for layer := 0; layer <= LAYERS; layer++ {
        treeData := merkleTree(data)
        sortedTrees.append(treeData)
        if layer < LAYERS {
            vdeEncode(slothIter, replicaId, data)
        }
        zigzag()
    }

    previousCommr := nil
    for _, replicaTree := range sortedTrees {
        commR := replicaTree.root()
        if previousCommr != nil {
            commD := previousCommr
            tau := (commR, commD)
            taus.append(tau)
        }
        auxs.append(replicaTree);
        previousCommr = commR
    }

    return (taus, auxs)
}
```

Note: The function `zigzag` just inverts an internal `bool` that tracks whether the `graphIsReversed()` or not.

```go
func vdeEncode(slothIter uint, replicaId uint, data []byte) {
    var parents [PARENT_COUNT]uint
    for n := 0; n < GRAPH_SIZE; n++ {
        var node uint
        if !graphIsReversed() {
            node = n
        } else {
            // If the graph is reversed, traverse in reverse order.
            node = GRAPH_SIZE - n - 1
        }

        parents = parents(node)

        key := kdf(replicaId, node, parents, data)

        start := node * NODE_SIZE
        end := start + NODE_SIZE;
        nodeData := data[start:end])
        encoded := slothEncode(key, nodeData, slothIter)

        data[start:end] = encoded
    }
}
```

### ZigZag operation: Feistel

TODO: Add `FEISTEL_ROUNDS` and `FEISTEL_BYTES` (or find its definitions)

```go
func permute(numElements uint, index uint, keys [FEISTEL_ROUNDS]uint) uint {
    u := feistelEncode(index, keys)

    while u >= numElements {
        u = feistelEncode(u, keys)
    }
    // Since we are representing `numElements` using an even number of bits,
    // that can encode many values above it, so keep repeating the operation
    // until we land in the permitted range.

    return u
}

func feistelEncode(index uint, keys [FEISTEL_ROUNDS]uint) uint {
    left, right, rightMask, halfBits := commonSetup(index)

    for _, key := range keys {
        left, right = right, left ^ feistel(right, key, rightMask)
    }

    return  (left << halfBits) | right
}

func commonSetup(index uint) (uint, uint, uint, uint) {
    numElements := GRAPH_SIZE * EXPANSION_DEGREE
    nextPow4 := 4;
    halfBits := 1
    while nextPow4 < numElements {
        nextPow4 *= 4
        halfBits += 1
    }

    rightMask = (1 << halfBits) - 1
    leftMask = rightMask << halfBits

    right := index & rightMask
    left := (index & leftMask) >> halfBits

    return  (left, right, rightMask, halfBits)
}

// Round function of the Feistel network: `F(Ri, Ki)`. Joins the `right`
// piece and the `key`, hashes it and returns the lower `uint32` part of
// the hash filtered trough the `rightMask`.
func feistel(right uint, key uint, rightMask uint) uint {
    var data [FEISTEL_BYTES]uint

    var r uint
    if FEISTEL_BYTES <= 8 {
        data[0] = uint8(right >> 24)
        data[1] = uint8(right >> 16)
        data[2] = uint8(right >> 8)
        data[3] = uint8(right)

        data[4] = uint8(key >> 24)
        data[5] = uint8(key >> 16)
        data[6] = uint8(key >> 8)
        data[7] = uint8(key)

        hash := blake2b(data)

        r =   hash[0]) << 24
            | hash[1]) << 16
            | hash[2]) << 8
            | hash[3])
    } else {
        data[0]  = uint8(right >> 56)
        data[1]  = uint8(right >> 48)
        data[2]  = uint8(right >> 40)
        data[3]  = uint8(right >> 32)
        data[4]  = uint8(right >> 24)
        data[5]  = uint8(right >> 16)
        data[6]  = uint8(right >> 8)
        data[7]  = uint8(right)

        data[8]  = uint8(key >> 56)
        data[9]  = uint8(key >> 48)
        data[10] = uint8(key >> 40)
        data[11] = uint8(key >> 32)
        data[12] = uint8(key >> 24)
        data[13] = uint8(key >> 16)
        data[14] = uint8(key >> 8)
        data[15] = uint8(key)

        hash := blake2b(data)

        r =   hash[0] << 56
            | hash[1] << 48
            | hash[2] << 40
            | hash[3] << 32
            | hash[4] << 24
            | hash[5] << 16
            | hash[6] << 8
            | hash[7]
    }

    return r & rightMask
}

// Inverts the `permute` result to its starting value for the same `key`.
func invertPermute(numElements uint, index uint, keys [FEISTEL_ROUNDS]uint) uint {
    u := feistelDecode(index, keys)

    while u >= numElements {
        u = feistelDecode(u, keys);
    }
    return u
}

func feistelDecode(index uint, keys [FEISTEL_ROUNDS]uint) uint {
    left, right, rightMask, halfBits := commonSetup(index)

    for _, key := range reversed(keys) {
        left, right = right ^ feistel(left, keys, rightMask), left
    }

    return (left << halfBits) | right
}
```



## Proof Generation 

Overview:

- Challenge Derivation
- Proof Generation
- Circuit Proof Generation

### Challenge Derivation

This is the Fiat-Shamir transform that turns the interactive Proof-of-Replication into non-interactive in the Random Oracle model.

TODO: define `Domain` (for practical purposes a `uint`) and `LayerChallenges` (or find existing definition).

```go
// TODO: we should replace the word commitment with the word seed, this will be more interactive porep friendly
func DeriveChallenges(challenges LayerChallenges, layer uint, leaves uint,
    replicaId Domain, commitment Domain, k uint) []uint {

    n := challenges.ChallengesForLayer(layer)
    var derivedChallenges [n]uint
    for i := 0; i < n; i++ {
        bytes := []byte(replicaId)
        bytes.append([]byte(commitment));
        bytes.append(layer);
        bytes.append(toLittleEndian(n * k + i))

        // For now, we cannot try to prove the first or last node, so make 
        // sure the challenge can never be 0 or leaves - 1.
        big_mod_challenge := blake2s(bytes) % (leaves - 2);
        derivedChallenges[i] = big_mod_challenge + 1
    }
}
```

### Challenge Generation

TODO: we may need to remove this section.

Calculate `LAYER_CHALLENGES : [LAYERS]uint`: Number of challenges per layer. (This will be passed to the ZigZag circuit proof.)

Derive challenges for each layer (call `DeriveChallenges()`).

### Proof Generation

```go
let layer_proofs = []

for l in 0..LAYERS {
  let replica = layer_replicas[l]
  let replica_tree = layer_trees[l]
  
  for c in derive_challenges(LAYER_CHALLENGES[l])
    data_inclusion_proof = inclusion_proof(data[c], DataTree, CommR_<l>)
    replica_inclusion_proof = inclusion_proof(replica[c], replica_tree, CommR_<l+1>) || FAIL// Prove the replica. TODO explain replica[].
    
    // *** let kdf_preimage = [replica_id] ***
    let parent_replica_inclusion_proofs = []
    for p in parents(c) {
      // *** kdf_preimage.push(p)***
      parent_replica_inclusion_proofs.push(inclusion_proof(p, CommR_<l+1>))
    }
    // *** let key = kdf(kdf_preimage); ***
    
    // *** encode(key, data[c]) == replica[c]
    // *** We don't actually need to encode in the proof. ***
    // TODO: move this ***stuff*** to verification.

    layer_proof.push((data_inclusion_proof, replication_inclusion_proof, parent_replica_inclusion_proofs))    
  } 
}

return layer_proofs, CommRstar, CommRLast
```

TODO: reconcile outputs of non-circuit proof with inputs to circuit proof.

### Circuit Proof Generation

See [# ZigZag: Offline PoRep Circuit Spec](zigzag-circuit.md) for details of Circuit Proof Generation.
