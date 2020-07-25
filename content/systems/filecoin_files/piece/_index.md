---
title: Piece
weight: 2
bookCollapseSection: true
dashboardWeight: 1.5
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# The Filecoin Piece & Data Representation in Filecoin
---

The _Filecoin Piece_ is the main _unit of negotiation_ for data that users store on the Filecoin network. The Filecoin Piece is _not a unit of storage_, it is not of a specific size, but is upper-bounded by the size of the _Sector_. A Filecoin Piece can be of any size, but if a Piece is larger than the size of a Sector that the miner supports it has to be split into more Pieces so that each Piece fits into a Sector.

A `Piece` is an object that represents a whole or part of a `File`,
and is used by `Clients` and `Miners` in `Deals`. `Clients` hire `Miners`
to store `Pieces`. 

The Piece data structure is designed for proving storage of arbitrary
IPLD graphs and client data. This diagram shows the detailed composition
of a Piece and its proving tree, including both full and bandwidth-optimized
Piece data structures.


![Pieces, Proving Trees, and Piece Data Structures](pieces.png)

## Data Representation in the Filecoin Network

It is important to highlight that data submitted to the Filecoin network go through several transformations before they come to the format at which the `StorageProvider` stores it.

1. When a piece of data, or file is submitted to Filecoin (in some raw system format) it is transformed into a _UnixFS DAG style data representation_ (in case it is not in this format already, e.g., from IPFS-based applications). The hash that represents the root of the IPLD DAG of the UnixFS file is the _Payload CID_, which is used in the Retrieval Market. The Payload CID is identical to an IPFS CID.
2. In order to make a _Filecoin Piece_ the UnixFS IPLD DAG is serialised into a ["Content-Addressable aRchive" (.car)](https://github.com/ipld/specs/blob/master/block-layer/content-addressable-archives.md#summary) file, which is in raw bytes format.
3. The resulting .car file is _padded_ with extra bits.
4. The next step is to calculate the Merkle root out of the hashes of the Piece. The resulting root of the Merkle tree is the **Piece CID**. This is also referred to as _CommP_ or _Piece Commitment_.
5. At this point, the Piece is included in a Sector together with data from other deals. The `StorageProvider` then calculates Merkle root for all the Pieces inside the Sector. The root of this tree is _CommD_ (aka _Commitment of Data_ or `UnsealedSectorCID`).
6. The `StorageProvider` is then sealing the sector and the root of the resulting Merkle root is the _CommR_ (or _Commitment of Replication_).

{{< hint info >}}
The code below is out-of-date
{{< /hint >}}

{{< embed src="piece.id" lang="go" >}}
