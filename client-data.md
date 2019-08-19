# Client Data Processing

This document describes two formats:

- **Transfer format**
- **Storage format**

## Transfer formats

The transfer format is the format to transfer a file over the network. This format SHALL be used for the initial transfer (from clients to storage miners) and for later retrievals (from storage miners to the clients).

The default transfer format is `unixfsv1`. Cliens MAY agree to use other formats of their preference.

### `unixfsv1`

The default transfer format is Unixfsv1 with the following parameters:

- Chunking: Fixed, 1MB
- Leaf Format: Raw
- Max Branch Width: 1024

For details on how UnixfsV1 works, see its spec [here](https://github.com/ipfs/specs/tree/master/unixfs).

## Storage Formats

The Storage Format MUST be use for generating Filecoin proofs and hashing sectors data. 

The current required storage format is `paddedfr32v1`.

### `paddedfr32v1`

A correctly formatted `paddedfr32v1` data must have:

- **Fr32 Padding**: Every 32 bytes blocks MUST contain two zeroes in the most significant bits (every 254 bits must be followed by 2 zeroes if interpreted as little-endian number). That is, for each block, `0x11000000 & block[31] == 0`.
- **Piece Padding**: In order to generate minimal `PieceInclusionProofs`, blocks of 32 zero bytes MUST be added so that the total number of blocks (including *piece padding*) is a power of two. **Piece Padding** can be omitted if the prover wishes to generate unaligned proofs. [NOTE: not yet fully specified.]

**Why do we need a special Storage Encoding Format?** In the Filecoin proofs we do operations in an arithmetic field of size `p`, where `p` is a prime of size `2^255`, hence the size of the data blocks must be smaller than `p`. We cautiously decide to have data blocks of size 254 to avoid possible overflows (data blocks numerical representation is bigger than `p`). 



## Piece Commitment

### commP

A piece commitment (`commP`) is the root hash of a piece that a client wants to store in Filecoin. It is generated using `RepHash` (as described in [Proof-of-Replication](zigzag-porep.md)) on some raw data which has been zero-padded to a multiple of 127 bytes, then preprocessed yielding `Fr32 padded` data which is a multiple of 128 bytes. 

