---
title: Piece
weight: 2
bookCollapseSection: true
dashboardWeight: 1.5
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# The Filecoin Piece

The _Filecoin Piece_ is the main _unit of negotiation_ for data that users store on the Filecoin network. The Filecoin Piece is _not a unit of storage_, it is not of a specific size, but is upper-bounded by the size of the _Sector_. A Filecoin Piece can be of any size, but if a Piece is larger than the size of a Sector that the miner supports it has to be split into more Pieces so that each Piece fits into a Sector.

A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Storage Clients` and `Storage Miners` in `Deals`. `Storage Clients` hire `Storage Miners` to store `Pieces`.

The Piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a Piece and its proving tree, including both full and bandwidth-optimized
Piece data structures.

![Pieces, Proving Trees, and Piece Data Structures](pieces.png)

## Data Representation

It is important to highlight that data submitted to the Filecoin network go through several transformations before they come to the format at which the `StorageProvider` stores it.

Below is the process followed from the point a user starts preparing a file to store in Filecoin to the point that the provider produces all the identifiers of Pieces stored in a Sector.

The first three steps take place on the client side.

1. When a client wants to store a file in the Filecoin network, they start by producing the IPLD DAG of the file. The hash that represents the root node of the DAG is an IPFS-style CID, called _Payload CID_.

2. In order to make a _Filecoin Piece_, the IPLD DAG is serialised into a ["Content-Addressable aRchive" (.car)](https://github.com/ipld/specs/blob/master/block-layer/content-addressable-archives.md#summary) file, which is in raw bytes format. A CAR file is an opaque blob of data that packs together and transfers IPLD nodes. The _Payload CID_ is common between the CAR'ed and un-CAR'ed constructions. This helps later during data retrieval, when data is transferred between the storage client and the storage provider as we discuss later.

3. The resulting .car file is _padded_ with extra zero bits in order for the file to make a binary Merkle tree. To achieve a clean binary Merkle Tree the .car file size has to be in some power of two (^2) size. A padding process, called `Fr32 padding`, which adds two (2) zero bits to every 254 out of every 256 bits is applied to the input file. At the next step, the padding process takes the output of the `Fr32 padding` process and finds the size above it that makes for a power of two size. This gap between the result of the `Fr32 padding` and the next power of two size is padded with zeros.

In order to justify the reasoning behind these steps, it is important to understand the overall negotiation process between the `StorageClient` and a `StorageProvider`. The piece CID or CommP is what is included in the deal that the client negotiates and agrees with the storage provider. When the deal is agreed, the client sends the file to the provider (using GraphSync). The provider has to construct the CAR file out of the file received and derive the Piece CID on their side. In order to avoid the client sending a different file to the one agreed, the Piece CID that the provider generates has to be the same as the one included in the deal negotiated earlier.

The following steps take place on the `StorageProvider` side (apart from step 4 that can also take place at the client side).

4. Once the `StorageProvider` receives the file from the client, they calculate the Merkle root out of the hashes of the Piece (padded .car file). The resulting root of the clean binary Merkle tree is the **Piece CID**. This is also referred to as _CommP_ or _Piece Commitment_ and as mentioned earlier, has to be the same with the one included in the deal.

5. The Piece is included in a Sector together with data from other deals. The `StorageProvider` then calculates Merkle root for all the Pieces inside the Sector. The root of this tree is _CommD_ (aka _Commitment of Data_ or `UnsealedSectorCID`).

6. The `StorageProvider` is then sealing the sector and the root of the resulting Merkle root is the _CommRLast_.

7. Proof of Replication (PoRep), SDR in particular, generates another Merkle root hash called _CommC_, as an attestation that replication of the data whose commitment is _CommD_ has been performed correctly.

8. Finally, _CommR_ (or _Commitment of Replication_) is the hash of CommC || CommRLast.

**IMPORTANT NOTES:**

- `Fr32` is a 32-bit representation of a field element (which, in our case, is the arithmetic field of BLS12-381). To be well-formed, a value of type `Fr32` must _actually_ fit within that field, but this is not enforced by the type system. It is an invariant which must be perserved by correct usage.
  In the case of so-called `Fr32 padding`, two zero bits are inserted 'after' a number requiring at most 254 bits to represent. This guarantees that the result will be `Fr32`, regardless of the value of the initial 254 bits. This is a 'conservative' technique, since for some initial values, only one bit of zero-padding would actually be required.
- Steps 2 and 3 above are specific to the Lotus implementation. The same outcome can be achieved in different ways, e.g., without using `Fr32` bit-padding. However, any implementation has to make sure that the initial IPLD DAG is serialised and padded so that it gives a clean binary tree, and therefore, calculating the Merkle root out of the resulting blob of data gives the same **Piece CID**. As long as this is the case, implementations can deviate from the first three steps above.
- Finally, it is important to add a note related to the _Payload CID_ (discussed in the first two steps above) and the data retrieval process. The retrieval deal is negotiated on the basis of the _Payload CID_. When the retrieval deal is agreed, the retrieval miner starts sending the unsealed and "un-CAR'ed" file to the client. The transfer starts from the root node of the IPLD Merkle Tree and in this way the client can validate the _Payload CID_ from the beginning of the transfer and verify that the file they are receiving is the file they negotiated in the deal and not random bits.

## PieceStore

The `PieceStore` module allows for storage and retrieval of Pieces from some local storage. The piecestore's main goal is to help the [storage](https://github.com/filecoin-project/go-fil-markets/blob/master/storagemarket) and [retrieval market](https://github.com/filecoin-project/go-fil-markets/blob/master/retrievalmarket) modules to find where sealed data lives inside of sectors. The storage market writes the data, and retrieval market reads it in order to send out to retrieval clients.

The implementation of the PieceStore module can be found [here](https://github.com/filecoin-project/go-fil-markets/tree/master/piecestore).
