# Filecoin ZigZag Proof of Replication

ZigZagDrgPorep is a layered PoRep which replicates layer by layer.
Between layers, the graph is 'reversed' in such a way that the dependencies expand with each iteration.
This reversal is not a straightforward inversion -- so we coin the term 'zigzag' to describe the transformation.
Each graph can be divided into base and expansion components.
The 'base' component is an ordinary DRG. The expansion component attempts to add a target (expansion_degree) number of connections
between nodes in a reversible way. Expansion connections are therefore simply inverted at each layer.
Because of how DRG-sampled parents are calculated on demand, the base components are not. Instead, a same-degree
DRG with connections in the opposite direction (and using the same random seed) is used when calculating parents on demand.
For the algorithm to have the desired properties, it is important that the expansion components are directly inverted at each layer.
However, it is fortunately not necessary that the base DRG components also have this property.


__*Filecoin ZigZag Proof of Replication*__ is the process by which raw data is transformed into a replica and a proof of replication.

This comprises the following steps:
 - Preprocessing
 - Sector Padding
 - Replication
 - Proof Generation
 - Circuit Proof Generation
 
 Together, __*Replication*__, __*Proof Generation*__, and __*Circuit Proof Generation*__ constitute the `Seal` operation, as described in [Filecoin Proofs](proofs.md).
 
## Hash Functions

 *__Filecoin ZigZag Proof of Replication__ as described here is generic over the following hash functions:*

  - KDF hash: a hash function with 32-byte digest size: default is Blake2s
  - CommR Hash: a hash function with 32-byte digest size: default is Blake2s
  - RepCompress: a hash function with 32-byte digest size: default is pedersen hashing over jubjub curve.

## Preprocessing

Raw data to be replicated is first preprocessed to yield an __*unsealed sector*__. The size of the raw data is __*RawSize*__ bytes.

First apply preprocessing as described in [Filecoin Client Data Processing](client-data.md):

> Preprocessing adds two zero bits after every 254 bits of original data, yielding a sequence of 32-byte blocks, each of which contains two zeroes in the most-significant bits, when interpreted as a little-endian number. That is, for each block, 0x11000000 & block[31] == 0.

## Sector Padding
After preprocessing, the data is padded with zero bytes so that the total length of the data is the sector size, __*SectorSize*__. __*SectorSize*__ must be one of a set of explicitly supported sizes and must be a power-of-two multiple of 32 bytes. The preprocessed and padded data is now considered to be an __*unsealed sector*__.

## Replication

#### Public Parameters

*The following public parameters are shared with the ZigZag circuit proof:*

 - `LAYERS : UInt`: Number of DRG layers.
 - `EXPANSION_DEGREE: UInt`: Degree of each bipartite expander graph to extend dependencies between layers.
 - `BASE_DEGREE: UInt`: Degree of each Depth Robust Graph.
 - `TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is (log_2(Size of original data in bytes/32 bytes per leaf)).
 - `PARENT_COUNT : UInt`: Defined as `EXPANSION_DEGREE+BASE_DEGREE`.

*The following additional public parameters are required:*

 - `TAPER` : Float: Fraction of each layer's challenges by which to reduce next-lowest layer's challenge count.
 - `TAPER_LAYERS`: UInt: Number of layers 
 `Data` is a byte array initialized to the content of __*unsealed sector*__ and will be mutated in-place by the replication process.

Replication proceeds as follows:

- Calculate `ReplicaID` using `Hash` (Blake2s):
```
ReplicaID := Hash(ProverID || SectorID || ticket)
```
- Perform `RepHash` on `Data` to yield `CommD` and `TreeD`:
```
CommD, TreeD = RepHash(data)
```

`RepHash` constructs a binary merkle tree from the resulting blocks, designated as the __*leaves*__ — by applying the __*RepHash Compression Function*__, `RepCompress`, to adjacent pairs of leaves. The final result is the merkle root of the constructed tree.

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

## Proof Generation

### Challenge Generation

Calculate `LAYER_CHALLENGES : [LAYERS]UInt`: Number of challenges per layer. (This will be passed to the ZigZag circuit proof.)

Derive challenges for each layer (call `derive_challenges()`).

### Layer Challenge Counts

TODO: define `Challenge` (or find existing definition)

```go
func ChallengesForLayer(challenge Challenge, layer uint) uint {

    switch challenge.Type {
        case Fixed:
            return challenge.Count
        case Tapered:
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

### Challenge Derivation

TODO: define `Domain` and `LayerChallenges` (or find existing definition)

```go
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

### Layer Replication
TODO: define `TransformedLayers` (or find existing definition)

Note: representing the (simpler) sequential replication case.

```go
func transformAndReplicateLayers(graph Graph, slothIter uint,
    replicaId Domain, data []byte) TransformedLayers {

    var taus [LAYERS]Domain
    var auxs [LAYERS]Tree
    var sortedTrees [LAYERS]Tree

    currentGraph := graph
    for layer := 0; layer <= LAYERS; layer++ {
        treeData := currentGraph.merkleTree(data)
        sortedTrees.append(treeData)
        if layer < LAYERS {
            encode(currentGraph, slothIter, replicaId, data)
        }
        currentGraph = currentGraph.zigzag()
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

### Graph Structure
TODO: define `parents()` and `expansion_parents()` and use in `replicate_layer()`.
(`parents()` is actually used in `encode()`, should `encode()` be defined in the previous section then?)
TODO: define `Graph`.

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
            rng := ChaChaRng.from_seed(Graph.seed)
            // FIXME: Might be an implementation detail worth hiding.

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
        if Graph.reversed {
            transformed := feistelInvertPermute(Graph.size * EXPANSION_DEGREE, a, feistelKeys)
        } else {
            transformed := feistelPermute(Graph.size * EXPANSION_DEGREE, a, feistelKeys)
        }
        other := transformed / EXPANSION_DEGREE
        // Collapse the output in the matrix search space to the row of the corresponding
        // node (losing the column information, that will be regenerated later when calling
        // back this function in the `reversed` direction).

        if Graph.reversed {
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
        if Graph.reversed {
            parents[currentLength + i] = Graph.size - 1
        } else {
            parents[currentLength + i] = 0
        }
    }
    assert(len(parents) == PARENT_COUNT)
    sort(parents)

    for _, parent := range parents {
        if Graph.forward {
            assert(parent <= rawNode)
        } else {
            assert(parent >= rawNode)
        }
    }
}

func realIndex(rawNode uint) uint {
    if Graph.reversed {
        return Graph.size - 1 - rawNode
    } else {
        return rawNode
    }
}
```

### ZigZag operation: Feistel

TODO: Add `FEISTEL_ROUNDS` and `FEISTEL_BYTES` (or find its definitions)

```go
func permute(numElements uint, index uint, keys [FEISTEL_ROUNDS]uint) uint {
    u := encode(index, keys)

    while u >= numElements {
        u = encode(u, keys)
    }
    // Since we are representing `numElements` using an even number of bits,
    // that can encode many values above it, so keep repeating the operation
    // until we land in the permitted range.

    return u
}

func encode(index uint, keys [FEISTEL_ROUNDS]uint) uint {
    left, right, rightMask, halfBits := commonSetup(index)

    for _, key := range keys {
        left, right = right, left ^ feistel(right, key, rightMask)
    }

    return  (left << halfBits) | right
}

func commonSetup(index uint) (uint, uint, uint, uint) {
    numElements := Graph.size * EXPANSION_DEGREE
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
    u := decode(index, keys)

    while u >= numElements {
        u = decode(u, keys);
    }
    return u
}

func decode(index uint, keys [FEISTEL_ROUNDS]uint) uint {
    left, right, rightMask, halfBits := commonSetup(index)

    for _, key := range reversed(keys) {
        left, right = right ^ feistel(left, keys, rightMask), left
    }

    return (left << halfBits) | right
}
```

## Proof Generation

```
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

## Circuit Proof Generation
See [# ZigZag: Offline PoRep Circuit Spec](zigzag-circuit.md) for details of Circuit Proof Generation.
