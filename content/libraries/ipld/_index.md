---
title: IPLD
bookCollapseSection: true
weight: 4
dashboardWeight: 1
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# IPLD

The InterPlanetary Linked Data or [IPLD](https://ipld.io/) is the data model of the content-addressable web. It provides standards and formats to build Merkle-DAG data-structures, like those that represent a filesystem. IPLD allows us to treat all hash-linked data structures as subsets of a unified information space, unifying all data models that link data via hashes as instances of IPLD. This means that data can be linked and referenced from totally different data structures in a global namespace. This is a very useful feature that is used extensively in Filecoin.

IPLD introduces several concepts and protocols, such as the concept of content addressing itself, codecs such as DAG-CBOR, file formats such as Content Addressable aRchives (CARs), and protocols such as GraphSync.

Please refer to the [IPLD specifications repository](https://github.com/ipld/specs) for more information.

## DAG-CBOR encoding

All Filecoin system data structures are stored using DAG-CBOR (which is an IPLD codec). DAG-CBOR is a more strict subset of CBOR with a predefined tagging scheme, designed for storage, retrieval and traversal of hash-linked data DAGs.

Files and data stored on the Filecoin network are also stored using various IPLD codecs (not necessarily DAG-CBOR). IPLD provides a consistent and coherent abstraction above data that allows Filecoin to build and interact with complex, multi-block data structures, such as HAMT and AMT. Filecoin uses the DAG-CBOR codec for the serialization and deserialization of its data structures and interacts with that data using the IPLD Data Model, upon which various tools are built. IPLD Selectors are also used to address specific nodes within a linked data structure (see GraphSync below).

Please refer to the [DAG-CBOR specification](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-cbor.md) for more information.

## Content Addressable aRchives (CARs)

The Content Addressable aRchives (CAR) format is used to store content addressable objects in the form of IPLD block data as a sequence of bytes; typically in a file with a `.car` filename extension.

The CAR format is used to produce a _Filecoin Piece_ (the main representation of files in Filecoin) by serialising its IPLD DAG. The `.car` file then goes through further transformations to produce the _Piece CID_.

Please refer to the [CAR specification](https://github.com/ipld/specs/blob/master/block-layer/content-addressable-archives.md) for further information.

## GraphSync

GraphSync is a request-response protocol that synchronizes _parts_ of a graph (an authenticated Directed Acyclic Graph - DAG) between different peers. It uses _selectors_ to identify the specific subset of the graph to be synchronized between different peers.

GraphSync is used by Filecoin in order to synchronize parts of the blockchain.

Please refer to the [GraphSync specification](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md) for more information.
