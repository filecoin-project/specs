# Filecoin Client Data Processing

Data that clients wish to store must be pre-processed before it can be sent to miners and have a storage deal made for it.

This processing is comprised of two top-level pieces:

1. **Transfer Encoding**
2. **Storage Encoding**

Transfer Encoding is a pre-processing for the file to allow it to be transferred over the network more easily. Both for the initial transfer to the storage miner, and for later retrievals. This mechanism can be done in any way that is supported by the client and the miner (alternative strategies may be implemented).

Storage Encoding is a pre-processing for the file that allows the Filecoin storage market and proofs to reference the data. Filecoin operates on a special type of merkle tree generated with a function we call `RepHash`. The root hash of this merkle tree is called the `Piece Commitment` or `CommP` and is included in the deal when it is proposed to a miner.

The next sections will go into some additional detail on how each of these encodings works.

## Transfer Encoding

Transfer Encoding may be computed in a variety of ways, as long as it is supported by both the client and miner involved. Here, we document the `default` transfer encoding.

### `default`

The default transfer encoding is Unixfsv1 with the following parameters:

- Chunking: Fixed, 1MB
- Leaf Format: Raw
- Max Branch Width: 1024

For details on how UnixfsV1 works, see its spec [here](https://github.com/ipfs/specs/tree/master/unixfs).

## Storage Encoding

Storage Encoding is computed by building a merkle tree out of the data using `RepHash`.

Inputs to `RepHash` must first be preprocessed then padded.

__*Preprocessing*__ adds two zero bits after every 254 bits of original data, yielding a sequence of 32-byte blocks, each of which contains two zeroes in the most-significant bits, when interpreted as a little-endian number. That is, for each block, `0x11000000 & block[31] == 0`.

Next, __*piece padding*__ is added: blocks of 32 zero bytes are added so that the total number of blocks (including __*piece padding*__) is a power of two.

`RepHash` constructs a binary merkle tree from the resulting blocks, designated as the *leaves* — by applying the __*RepHash Compression Function*__, `RepCompress`, to adjacent pairs of leaves. The final result is the merkle root of the constructed tree.

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
