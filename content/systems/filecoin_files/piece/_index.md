---
title: Piece
weight: 2
bookCollapseSection: true
dashboardWeight: 1.5
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# The Filecoin Piece & Data Representation in Filecoin
---

The _Filecoin Piece_ is the main _unit of negotiation_ for data that users store on the Filecoin network. The Filecoin Piece is _not a unit of storage_, it is not of a specific size, but is upper-bounded by the size of the _Sector_. A Filecoin Piece can be of any size, but if a Piece is larger than the size of a Sector that the miner supports it has to be split into more Pieces so that each Piece fits into a Sector.

A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Storage Clients` and `Storage Miners` in `Deals`. `Storage Clients` hire `Storage  Miners` to store `Pieces`. 

The Piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a Piece and its proving tree, including both full and bandwidth-optimized
Piece data structures.

![Pieces, Proving Trees, and Piece Data Structures](pieces.png)

## Data Representation in the Filecoin Network

It is important to highlight that data submitted to the Filecoin network go through several transformations before they come to the format at which the `StorageProvider` stores it.

Below is the process followed from the point a user starts preparing a file to store in Filecoin to the point that the provider produces all the identifiers of Pieces stored in a Sector.

The first three steps take place on the client side.

1. When a client wants to store a file in the Filecoin network, they start by producing the IPLD DAG of the file. The hash that represents root node of the DAG is an IPFS-style CID, called _Payload CID_.

2. In order to make a _Filecoin Piece_, the IPLD DAG is serialised into a ["Content-Addressable aRchive" (.car)](https://github.com/ipld/specs/blob/master/block-layer/content-addressable-archives.md#summary) file, which is in raw bytes format. A CAR file is an opaque blob of data that packs together and transfers IPLD nodes. The _Payload CID_ is common between the CAR'ed and un-CAR'ed constructions. This helps later during data retrieval, when data is transferred between the storage  client and the storage provider as we discuss later.

3. The resulting .car file is _padded_ with extra zero bits in order for the file to make a binary Merkle tree. To achieve a clean binary Merkle Tree the the .car file size has to be in some power of two (^2) size. Therefore, the padding process takes the input file and finds the size above the input size that makes for a power of two size. The gap between the input size and the next power of two size is padded with zeros. Finally, an extra padding process is applied, called `fr32 padding` which adds two (2) zero bits to every 254 bits (to make 256 bits in total).

The term `fr32` is derived from the name of a struct that Filecoin uses to represent the elements of the arithmetic field of a pairing-friendly curve, specifically Bls12-381—which justifies use of 32 bytes. `F` stands for "Field", while `r` is simply a mathematic letter-as-variable substitution used to denote the modulus of this particular field.

Given that padding involves only adding zeroes, the padding process described in step 3 above can be described or be implemented differently. For instance, `fr32` padding can be applied first to add 2 zero bits per 254 bits throughout the input client file. Once this is done, the process should find the next ^2 size above the size resulting from the `fr32` padding and add zeroes to reach this size.

In order to justify the reasoning behind these steps, it is important to understand the overall negotiation process between the `StorageClient` and a `StorageProvider`. The piece CID or CommP is what is included in the deal that the client negotiates and agrees with the storage provider. When the deal is agreed, the client sends the file to the provider (using GraphSync). The provider has to construct the CAR file out of the file received and derive the Piece CID on their side. In order to avoid the client sending a different file to the one agreed, the Piece CID that the provider generates has to be the same as the one included in the deal negotiated earlier.

The following steps take place on the `StorageProvider` side (apart from step 4 that can also take place at the client side).

4. Once the `StorageProvider` receives the file from the client, they calculate the Merkle root out of the hashes of the Piece (padded .car file). The resulting root of the clean binary Merkle tree is the **Piece CID**. This is also referred to as _CommP_ or _Piece Commitment_ and as mentioned earlier, has to be the same with the one included in the deal.

5. The Piece is included in a Sector together with data from other deals. The `StorageProvider` then calculates Merkle root for all the Pieces inside the Sector. The root of this tree is _CommD_ (aka _Commitment of Data_ or `UnsealedSectorCID`).

6. The `StorageProvider` is then sealing the sector and the root of the resulting Merkle root is the _CommRLast_.

7. Proof of Replication (PoRep), SDR in particular, generates another Merkle root hash called _CommC_, as an attestation that replication of the data whose commitment is _CommD_ has been performed correctly.

8. Finally, _CommR_ (or _Commitment of Replication_) is the hash of CommC || CommRLast.

Finally, it is important to add a note related to the _Payload CID_ (discussed in the first two steps above) and the data retrieval process. The retrieval deal is negotiated on the basis of the _Payload CID_. When the retrieval deal is agreed, the retrieval miner starts sending the unsealed and "un-CAR'ed" file to the client. The transfer starts from the root node of the IPLD Merkle Tree and in this way the client can validate the _Payload CID_ from the beginning of the transfer and verify that the file they are receiving is the file they negotiated in the deal and not random bits.