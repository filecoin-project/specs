# Proof-of-Replication

{{<js>}}

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

**Encoding using DRGs**. By positioning data blocks into nodes in the DRG, we sequentially encode each node in the graph using its encoded parents. The depth robustness property of these graphs ensure that this process is not likely to be parallelizable.

**Layering DRGs**. ZigZag repeates this encoding by layering DRG graphs `LAYERS` times. The data represented in each DRG layer is the data encoded in the previous layer. We connect different layers using Bipartite Expander Graphs and at each layer, we reverse the graph edges with a technique which we call zigzag. The combination of DRGs, expander graphs and zigzag guarantee the security property of PoRep. The final encoded layer is the final replica.

**Generating ZigZag proofs**. Given the following public parameters:

- `ReplicaId` is a unique replica identifier (see the Filecoin Proofs spec for details)
- `CommD` is the Merkle Tree root hash of the input data to the first layer
- `CommRStar` is the hash of the concatenation of the `ReplicaId` and all the `CommR`s.
- `CommRLast` is the hash of the last encoded DRG layer.

A ZigZag proof proves that some data whose committment is `CommD` has been used to run a `Replicate` algorithm and generated some data whose commitment is `CommRLast`.

A ZigZag proof consists of a set of challenged DRG nodes (both encoded and unencoded) for each layer, a set of parent nodes for each challenged node and a Merkle tree inclusion proof for each node provided. The verifier can then verify the correct encoding of each node and that the nodes given were consistent with the provers' commitments.

**Making proofs succinct with SNARKs**: The proof size in the ZigZag is too large for blockchain usage (~100MB), mostly due to the large amount of Merkle tree inclusion proofs required to achieve security. We use SNARKs to generate a proof of knowledge of a correct ZigZag proof. In other words, we implement the ZigZag proof verification algorithm in an arithmetic circuit and use SNARKs to prove that it was evaluated correctly.

The SNARK circuit proves that given a Merkle root `CommD`, `CommRLast`, and `commRStar`, the prover knew the correct replicated data at each layer.

### PoRep in Filecoin

Proof-of-Replication proves that a Storage Miner is dedicating unique storage for each ***sector***. Filecoin Storage Miners collect new clients' data in a sector, run a slow encoding process (called `Seal`) and generate a proof (`SealProof`) that the encoding was generated correctly.

In Filecoin, PoRep provides two guarantees: (1) *space-hardness*: Storage Miners cannot lie about the amount of space they are dedicating to Filecoin in order to gain more power in the consensus; (2) *replication*: Storage Miners are dedicating unique storage for each copy of their clients data. 

Glossary:

- __*sector:*__ a fixed-size block of data of `SECTOR_SIZE` bytes which generally contains clients' data.
- __*unsealed sector:*__ a concrete representation (on disk or in memory) of a sector's that follows the "Storage Format" described in [Client Data Processing](client-data.md#storage-format) (currently `paddedfr32v1` is the required default).
- __*sealed sector:*__  a concrete representation (on disk or in memory) of the unique replica generated by `Seal` from an __*unsealed sector*__. A sector contains one or more ***pieces***.
- __*piece:*__ a block of data of at most `SECTOR_SIZE` bytes which is generally a client's file or part of.

## ZigZag Construction

### Public Parameters

The following public parameters are used in the ZigZag Replication and Proof Generation algorithms:

TODO: the Appendix should explain why we picked those values


| name | type | description | value |
| --- | --- | --- | ---: |
| `SECTOR_SIZE` | `uint` | Number of nodes in the DRG in bytes | `68,719,476,736` |
| `LAYERS` | `uint` | Number of Depth Robust Graphs stacked layers. | `10` |
| `BASE_DEGREE` | `uint` | In-Degree of each Depth Robust Graph. | `5` |
| `EXPANSION_DEGREE` | `uint` | Degree of each Bipartite Expander Graph to extend dependencies between layers. | `8` |
| `GRAPH_SEED` |  `uint` | Seed used for random number generation in `baseParents`. | `TODO` |
| `NODE_SIZE` | `uint` | Size of each node in bytes.| `32B` |


The following constants are computed from the public parameters:


| name | type | description | computation | value |
| --- | --- | --- | ---: | ---- |
| `PARENTS_COUNT` | `uint` | Total number of parent nodes |`EXPANSION_DEGREE + BASE_DEGREE` | `13` |
| `GRAPH_SIZE` | `uint` | Number of nodes in the graph  | `SECTOR_SIZE / NODE_SIZE` | `2,147,483,648` |
| `TREE_DEPTH` | `uint` | Height of the Merkle Tree of a sector | `LOG_2(GRAPH_SIZE)` | `31` |


The following additional public parameters are required:

- `TAPER` : `Float`: Fraction of each layer's challenges by which to reduce next-lowest layer's challenge count.
- `TAPER_LAYERS`: `uint`: Number of layers
  `Data` is a byte array initialized to the content of __*unsealed sector*__ and will be mutated in-place by the replication process.

### Hash Functions

We have describe three hash functions:

| name          | description                                                  | size of input | size of output | construction          |
| ------------- | ------------------------------------------------------------ | ------------- | -------------- | --------------------- |
| `KDFHash`     | Hash function used as a KDF to derive the key to encode a single node. | TODO          | `32B`          | `Blake2s-256`         |
| `CommRHash`   | Hash function used to hash all the commitments at every layer (`CommR`s) to generate `CommRStar` | TODO          | `32B`          | `Blake2s-256`         |
| `RepCompress` | Collision Resistant Hash function used for the Merkle tree.   | 2 x `32B` + integer height          | `32B`          | `JubjubPedersen`      |
| `RepHash`     | Merkle-tree based hash function used to generate commitments to sealed sectors, unsealed sectors, piece commitments and intermediate stepds of the Proof-of-Replication | TODO          | `32B`          | It uses `RepCompress` |

#### RepHash

`RepHash` is Merkle-tree based hash function used to generate commitments to sealed sectors, unsealed sectors, piece commitments and intermediate stepds of the Proof-of-Replication. The tree is a binary balanced Merkle-tree. The leaves of the Merkle tree are pairs of adjacent nodes.

`RepHash` inputs MUST respect a valid Storage Format.

```go
type node [32]uint8

// Create and return a balanced binary Merkle tree and its root commitment.
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
			hashed = RepCompress(input, height)
			nextRow = append(nextRow, hashed)
		}

		currentRow = nextRow
	}
    // The tree returned here is just a vector of rows for later use. Its representation is not part of the spec.
	return rows, currentRow[0]
}

```

### ZigZag Graph

The slow sequential encoding required is enforced by the depth robusness property of the ZigZag graph. 

**Encoding with ZigZag**: The data from a sector (of size `SECTOR_SIZE`) is divided in `NODE_SIZE` nodes (for a total of `GRAPH_SIZE` nodes) and arranged in a directed acyclic graph. The structure of the graph is used to encode the data sequentially: in order to encode a node, its parents must be encoded (see the "Layer Replication" section below). We repeat this process for `LAYERS` layers, where the input to a next layer is the output of the previous one.

**Generating the ZigZag graph**: The ZigZag graph is divided in `LAYERS` layers. Each layer is a directed acyclic graph and it combines a Depth Robust Graph (DRG) and a Bipartite Expander graph. 

We provide an algorithm (`ZigZag`) which computes the parents of a node. In high level, the parents of a node are computed by combining two algorithms: some parents (`BASE_DEGREE` of them) are computed via the `BucketSample` algorithm, others (`EXPANSION_DEGREE` of them) are computer via the `Chung` algorithm.In addition, every odd layer performs the "ZigZag" technique which reverts some edges and inverts some nodes.

#### `ZigZag`: ZigZag Graph algorithm 

Overview: on even layers, compute the DRG and the Bipartite Expander parents using respectively `BucketSample` and `ChungExpander`, on odd layers, compute the inverted DRG parents (using `BucketSample` on the inverted node: `GRAPH_SIZE - node - 1` and inverting the resulting parents `GRAPH_SIZE - parent -1` for each parent) and the reverted Biparate Expander parents (by reverting the edges)

##### Inputs

| name    | description                                       | Type   |
| ------- | ------------------------------------------------- | ------ |
| `node`  | The node for which the parents are being computed | `uint` |
| `layer` | The layer of the ZigZag graph                     | `uint` |

##### Outputs

| name      | description                                 | Type                  |
| --------- | ------------------------------------------- | --------------------- |
| `parents` | The parents of node `node` on layer `layer` | `[PARENTS_COUNT]uint` |

##### Algorithm

- If the `layer` is even:
  - Compute `drgParents = BucketSample(node)`
  - Compute `expanderParents = ChungExpander(node)`
  - Set `parents` to be the concatenation of `drgParents` and `expanderParents`

- If the `layer` is odd:
  - Invert the `node`: `inverted_node = GRAPH_SIZE - inverted_node - 1`
  - Compute the inverted DRG parents `invertedDRGParents`:
    - Compute `drgParents = BucketSample(n)`
    - For each `parent` in `drgParents`:
      - `invertedDRGParents.push(GRAPH_SIZE - parent - 1)`
  - Compute the reversed Expander parents: `reversedExpanderParents`:
    - Compute `reversedExpanderParents = InverseChungExpander(node)`
  - Set `parents` to be the concatenation of `invertedDrgParents` and `reversedExpanderParents`

##### Pseudocode

We provide below a more succinct representation of the algorithm:

```
func ZigZag(node uint, layer uint) {

  if layer % 2 == 0 {
  	// On even layers
  	let drgParents = BucketSample(node)
  	let expanderParents = ChungExpander(node)
  	return concat(drgParents, expanderParents)
  } else {
    // On odd layers
  	
  	// Inverting DRG parents
    let inverted_node = GRAPH_SIZE - node - 1
  	let drgParents = BucketSample(inverted_node)
  	let invertedDrgParents = []
  	for i in 0..drgParents.len() {
      invertedDrgParents.push(GRAPH_SIZE - drgParents[i] - 1)
  	}
  	// Reverting ChungExpander
		let reversedExpanderParents = ReversedChungExpander(node)
  	return concat(invertedDrgParents, reversedExpanderParents)
  }
}
```

##### Tests

- Each `parent` in `parents` MUST not be greater than `GRAPH_SIZE-1` and lower than `0`
- If `layer` is even:
  - Each `parent` in `parents` MUST be greater than `node` 
    - EXCEPT: if `node` is `0`, then all parents MUST be `0`
- if `layer` is odd:
  - Each `parent` in `parents` MUST be less than `node` 
    - EXCEPT: if `node` is `GRAPH_SIZE-1`, then all parents MUST be `GRAPH_SIZE-1`

##### Time-space tradeoff

Computing the parents using both `BucketSample` and `ChungExpander` (and `Reverse`) for every layer can be an expensive operation, however, this can be avoided by caching the parents. It is important to note that all the odd layers and all the even layers have the same structure.

#### `BucketSample`: Depth Robust Graphs algorithm

This section describes how to compute the "base parents" of the ZigZag graph, which is the equivalent of computing the parents of a Depth Robust Graph.

The properties of DRG graphs guarantee that a sector has been encoded with a slow, non-parallelizable process. We use the `BucketSample` algorithm that is based on DRSample ([ABH17](https://acmccs.github.io/papers/p1001-alwenA.pdf)) and described in [FBGB18](https://web.stanford.edu/~bfisch/porep_short.pdf) and generates a directed acyclic graph of in-degree `BASE_DEGREE`.

`BucketSample` DRG graphs are random graphs that can be deterministically generated from a seed; different seed lead with high probability to different graphs. In ZigZag, we use the same seed `GRAPH_SEED` for each layer of the ZigZag graph such that they are all based on the same DRG graph.

The parents of any node can be locally computed without computing the entire graph. We call the parents of a node calculated in this way *base parents*.

##### Pseudocode

```go
func BucketSample(node uint) (parents [BASE_DEGREE]uint) {
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
```

#### `ChungExpander`: Bipartite Expander Graphs

TODO: explain why we link nodes in the current layer

Every node in a layer has `EXPANSION_DEGREE` parents that are generated via the following algorithm. 

```go
func ChungExpander(node uint) (parents []uint) {
	parents := make([]uint, EXPANSION_DEGREE)

	feistelKeys := []uint{1, 2, 3, 4} // TODO
  
	for i := 0, p := 0; i < EXPANSION_DEGREE; i++ {
		a := node * EXPANSION_DEGREE + i
    transformed := feistelPermute(GRAPH_SIZE * EXPANSION_DEGREE, a, feistelKeys)
    other := transformed / EXPANSION_DEGREE
    if other < node {
      parents[p] = other
      p += 1
    }
  }
}

func ReverseChungExpander(node uint) (parents []uint) {
	parents := make([]uint, EXPANSION_DEGREE)

	feistelKeys := []uint{1, 2, 3, 4} // TODO
  
	for i := 0, p := 0; i < EXPANSION_DEGREE; i++ {
		a := node * EXPANSION_DEGREE + i
    transformed := invertFeistelPermute(GRAPH_SIZE * EXPANSION_DEGREE, a, feistelKeys)
    other := transformed / EXPANSION_DEGREE
    if other > node {
      parents[p] = other
      p += 1
    }
  }
}
```

##### Time-Space tradeoff

Computing these parents can be expensive (especially due to the hashing required by the Feistel algorithm). A miner can trade this computation by storing the expansion parents and the reversed expansion parents.

##### Feistel construction

We use three rounds of Feistel as a permutation to generate the parents of the Bipartite Expander graph.

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

// Inverts the `permute` result to its starting value for the same `key`.
func invertPermute(numElements uint, index uint, keys [FEISTEL_ROUNDS]uint) uint {
    u := feistelDecode(index, keys)

    while u >= numElements {
        u = feistelDecode(u, keys);
    }
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

func feistelDecode(index uint, keys [FEISTEL_ROUNDS]uint) uint {
    left, right, rightMask, halfBits := commonSetup(index)

    for _, key := range reversed(keys) {
        left, right = right ^ feistel(left, keys, rightMask), left
    }

    return (left << halfBits) | right
}
```

## Replication

> The Replication phase turns an *unsealed sector* into a *sealed sector*

Before running the `Replicate` algorithm, the prover must ensure that the sector is correctly formatted with a valid with the "Storage Format" described in [Filecoin Client Data Processing](client-data.md#storage-format) (currently `paddedfr32v1` is the required default).

TODO: inputs are missing

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

        key := KDFHash(replicaId, node, parents, data)

        start := node * NODE_SIZE
        end := start + NODE_SIZE;
        nodeData := data[start:end])
        encoded := slothEncode(key, nodeData, slothIter)

        data[start:end] = encoded
    }
}
```



## Proof Generation 

Overview:

- Challenge Derivation
- Proof Generation
- Circuit Proof Generation

TODO: write a single algorithm which includes the spec below

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

### Witness Generation

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

### SNARK Proof Generation

See [# ZigZag: Offline PoRep Circuit Spec]({{<relref "detail/zigzag_circuit_spec">}}) for details of Circuit Proof Generation.

## Appendex

### Layer Challenge Counts

TODO: define `Challenge` (or find existing definition)

TODO: we should just list current parameters and show this as a calculation for correctness, this should not mandatory to implement.

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
