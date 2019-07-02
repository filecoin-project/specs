# Filecoin Client Data Processing

Client must pre-process their data into a valid formats before sending them to miners, creating storage deals and generating mining proofs.

There are two formats:

- **Transfer Format**
- **Storage Format**

## Transfer Format

The Transfer Format is the format to transfer a file over the network. This format is used for the initial transfer (from clients to storage miners) and for later retrievals (from storage miners to the clients).

The default Transfer Format is `unixfsv1`, however clients and miner might agree to use other formats of their preference.

### `unixfsv1`

The default transfer encoding is Unixfsv1 with the following parameters:

- Chunking: Fixed, 1MB
- Leaf Format: Raw
- Max Branch Width: 1024

For details on how UnixfsV1 works, see its spec [here](https://github.com/ipfs/specs/tree/master/unixfs).

## Storage Format

The Storage Format it's the required format for generating Filecoin proofs and hashing sectors data. 

The current Storage Format is `paddedfr32v1`.

### `paddedfr32v1`

A correctly formatted `paddedfr32v1` data must have:

- **Fr32 Padding**: Every 32 bytes blocks must contain two zeroes in the most significant bits (every 254 bits must be followed by 2 zeroes if interpreted as little-endian number). That is, for each block, `0x11000000 & block[31] == 0`.
- **Piece Padding**: blocks of 32 zero bytes are added so that the total number of blocks (including *piece padding*) is a power of two.

**Why do we need a special Storage Encoding Format?** In the Filecoin proofs we do operations in an arithmetic field of size `p`, where `p` is a prime of size `2^255`, hence the size of the data blocks must be smaller than `p`. We cautiously decide to have data blocks of size 254 to avoid possible overflows (data blocks numerical representation is bigger than `p`). 



## Piece Commitment

TODO: check if this section is not already repeated

A piece commitment (`CommP`) is a commitment to a file that a client wants to store in Filecoin. A client must store `CommP` in order to check the integrity of a stored file for future retreaval and include it in the deal when proposing it to a miner.

### CommP

`CommP` is the root hash of a piece that a clients wants to store in Filecoin and it is generated on using `RepHash` on some raw data that respects the Storage Formats. 

### RepHash

```
TODO: this should be moved to the proofs spec
```

`RepHash` is the process to generate a Merkle tree root hash of sealed sectors, unsealed sectors and of the intermediate steps of Proof-of-Replication. It takes as input some data respecting a valid Storage Format and outputs a Merkle root hash. 

`RepHash` is constructed from a balanced binary Merkle tree. The leaves of the merkle tree is the output of `RepCompress` on two adjacent 32 bytes blocks.

```go
type node [32]uint8

// Create and return the root of a binary merkle tree.
// len(leaves) must be a power of 2.
func RepHash(leaves node) node {
	currentRow := leaves
	for height := 0; len(currentRow) > 1; height += 1 {
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

	return currentRow[0]
}

```