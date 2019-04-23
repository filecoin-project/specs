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

TODO(proofs team): Fill in how exactly this works.