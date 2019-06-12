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
let layer_replicas = [LAYERS][nodes]Uint8
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

```rust
    pub fn challenges_for_layer(&self, layer: usize) -> usize {
        match self {
            LayerChallenges::Fixed { count, .. } => *count,
            LayerChallenges::Tapered {
                taper,
                taper_layers,
                count,
                layers,
            } => {
                assert!(layer < *layers);
                let l = (layers - 1) - layer;

                let r: f64 = 1.0 - *taper;
                let t = min(l, *taper_layers);
                let total_taper = r.powi(t as i32);

                let calculated = (total_taper * *count as f64).ceil() as usize;

                // Although implied by the call to `ceil()` above, be explicit that a layer cannot contain 0 challenges.
                max(1, calculated)
            }
        }
    }
```

### Challenge Derivation
```rust
pub fn derive_challenges<D: Domain>(
    challenges: &LayerChallenges,
    layer: u8,
    leaves: usize,
    replica_id: &D,
    commitment: &D,
    k: u8,
) -> Vec<usize> {
    let n = challenges.challenges_for_layer(layer as usize);
    (0..n)
        .map(|i| {
            let mut bytes = replica_id.into_bytes();
            let j = ((n * k as usize) + i) as u32;
            bytes.extend(commitment.into_bytes());
            bytes.push(layer);
            bytes.write_u32::<LittleEndian>(j).unwrap();

            let hash = blake2s(bytes.as_slice());
            let big_challenge = BigUint::from_bytes_le(hash.as_ref());

            // For now, we cannot try to prove the first or last node, so make sure the challenge can never be 0 or leaves - 1.
            let big_mod_challenge = big_challenge % (leaves - 2);
            big_mod_challenge.to_usize().unwrap() + 1
        })
        .collect()
}
```

### Layer Replication
TODO: define and use `replicate_layer()`

### Graph Structure
TODO: define `parents()` and `expansion_parents()` and use in `replicate_layer()`.
```rust
fn parents(&self, node: usize) -> Vec<usize> {
    let m = self.base_degree;

    match node {
        // Special case for the first node, it self references.
        0 => vec![0; m as usize],
        // Special case for the second node, it references only the first one.
        1 => vec![0; m as usize],
        _ => {
            // seed = self.seed | node
            let mut seed = [0u32; 8];
            seed[0..7].copy_from_slice(&self.seed);
            seed[7] = node as u32;
            let mut rng = ChaChaRng::from_seed(&seed);

            let mut parents = Vec::with_capacity(m);
            for k in 0..m {
                // iterate over m meta nodes of the ith real node
                // simulate the edges that we would add from previous graph nodes
                // if any edge is added from a meta node of jth real node then add edge (j,i)
                let logi = ((node * m) as f32).log2().floor() as usize;
                let j = rng.gen::<usize>() % logi;
                let jj = cmp::min(node * m + k, 1 << (j + 1));
                let back_dist = rng.gen_range(cmp::max(jj >> 1, 2), jj + 1);
                let out = (node * m + k - back_dist) / m;

                // remove self references and replace with reference to previous node
                if out == node {
                    parents.push(node - 1);
                } else {
                    assert!(out <= node);
                    parents.push(out);
                }
            }

            parents.sort_unstable();

            parents
        }
    }
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
