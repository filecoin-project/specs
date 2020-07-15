---
title: IPLD Store
weight: 3
---

# IPLD Store - Local Storage for hash-linked data
---

{{<embed src="../../../../libraries/ipld/ipld.id" lang="go" >}}

IPLD is a set of libraries which allow for the interoperability of content-addressed data structures across different distributed systems. It provides a fundamental 'common language' to primitive cryptographic hashing, enabling data to be verifiably referenced and retrieved between two independent protocols. For example, a user can reference a git commit in a blockchain transaction to create an immutable copy and timestamp, or a data from a DHT can be referenced and linked to in a smart contract. 

## The Data Model

At its core, IPLD defines a [Data Model](https://github.com/ipld/specs/blob/master/data-model-layer/data-model.md) for representing data. The Data Model is designed for practical implementation across a wide variety of programming languages, while maintaining usability for content-addressed data and a broad range of generalized tools that interact with that data. 

The Data Model includes a range of standard primitive types (or "kinds"), such as booleans, integers, strings, nulls and byte arrays, as well as two recursive types: lists and maps. Because IPLD is designed for content-addressed data, it also includes a "link" primitive in its Data Model. In practice, links use the [CID](https://github.com/multiformats/cid) specification. IPLD data is organized into "blocks", where a block is represented by the raw, encoded data and its content-address, or CID. Every content-addressable chunk of data can be represented as a block, and together, blocks can form a coherent graph, or [Merkle DAG](https://docs.ipfs.io/guides/concepts/merkle-dag/).

Applications interact with IPLD via the Data Model, and IPLD handles marshalling and unmarshalling via a suite of codecs. IPLD codecs may support the complete Data Model or part of the Data Model. Two codecs that support the complete Data Model are [DAG-CBOR](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-cbor.md) and [DAG-JSON](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-json.md). These codecs are respectively based on the CBOR and JSON serialization formats but include formalizations that allow them to encapsulate the IPLD Data Model (including its link type) and additional rules that create a strict mapping between any set of data and it's respective content address (or hash digest). These rules include the mandating of particular ordering of keys when encoding maps, or the sizing of integer types when stored.

## IPLD in Filecoin

On the Filecoin network, IPLD is used in two ways:

- All system datastructures are stored in [IPLD](https://ipld.io) format, a data format akin to JSON but designed for storage, retrieval and traversal of hash-linked data DAGs.
- Files and data stored on the Filecoin network may also be stored in IPLD format. While this is not required, it offers the benefit of supporting [selectors](https://github.com/ipld/specs/blob/master/selectors/selectors.md) to retrieve a smaller subset of the total stored data, as opposed to inefficiently downloading the data set entirely.

IPLD provides a consistent and coherent abstraction above data that allows Filecoin to build and interact with complex, multi-block data structures, such as HAMT and Sharray. Filecoin uses the **DAG-CBOR** codec for the serialization and deserialization of its data structures and interacts with that data using the IPLD Data Model, upon which various tools are built. IPLD [Paths](https://github.com/ipld/specs/blob/master/data-model-layer/paths.md) are also used to address specific nodes within a linked data structure.


### IpldStores

The Filecoin network relies primarily on two distinct IPLD GraphStores:

- One `ChainStore` which stores the blockchain, including block headers, associated messages, etc.
- One `StateStore` which stores the payload state from a given blockchain, or the `stateTree` resulting from all block messages in a given chain being applied to the genesis state by the [Filecoin VM](systems/filecoin_vm).

The `ChainStore` is downloaded by a node from their peers during the bootstrapping phase of [Chain Sync](chainsync) and is stored by the node thereafter. It is updated on every new block reception, or if the node syncs to a new best chain.

The `StateStore` is computed through the execution of all block messages in a given `ChainStore` and is stored by the node thereafter. It is updated with every new incoming block's processing by the [VM Interpreter](interpreter), and referenced accordingly by new blocks produced atop it in the [block header's](block) `ParentState` field.
