---
menuTitle: IpldStore
title: "IpldStore - Local Storage for content-addressed data"
---

{{< readfile file="../../../../libraries/ipld/ipld.id" code="true" lang="go" >}}

Filecoin data structures are stored using [IPLD](https://ipld.io). IPLD is an abstraction layer and set of libraries for content-addressed data. IPLD serves as the data layer of [IPFS](https://ipfs.io) but is designed for use across many distributed data systems and provides interfaces in multiple languages for multiple data encoding formats (codecs), including JSON, CBOR, Git, Bitcoin.

<!-- TODO: insert final version of https://github.com/ipld/specs/issues/244#issuecomment-637269338 when ready, it'll land on the IPLD specs repo README -->

At its core, IPLD defines a [Data Model](https://github.com/ipld/specs/blob/master/data-model-layer/data-model.md) for representing data. The Data Model is designed for practical implementation across a wide variety of programming languages, while being as useful as possible for content-addressed data and a broad range of generalized tools to interact with that data. The Data Model includes a range of standard primitive types (or "kinds"), such as booleans, integers, strings, nulls and byte arrays, as well as two recursive types: lists and maps. Because IPLD is designed for content-addressed data, it also includes a "link" primitive in its Data Model. In practice, links use the [CID](https://github.com/multiformats/cid) specification. IPLD data is organized into "blocks", where a block is represented by the raw, encoded data and its content-address, or CID. Every content-addressable chunk of data can be represented as a block, and together, blocks can form a coherent graph, or [Merkle DAG](https://docs.ipfs.io/guides/concepts/merkle-dag/).

Applications interact with IPLD via the Data Model, and IPLD handles marshalling and unmarshalling via a suite of codecs. IPLD codecs may support the complete Data Model or part of the Data Model. Two codecs that support the complete Data Model are [DAG-CBOR](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-cbor.md) and [DAG-JSON](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-json.md). These codecs are respectively based on the CBOR and JSON serialization formats but include formalizations that allow them to encapsulate the IPLD Data Model (including its link type) and additional rules that create a strict mapping between any set of data and its respective content address (or hash digest). These rules include the mandating of particular ordering of keys when encoding maps, or the sizing of integer types when stored.

Filecoin uses the **DAG-CBOR** codec for the serialization and deserialization of its data structures and interacts with that data using the IPLD Data Model and various tools that are built for the IPLD Data Model. Filecoin makes use of IPLD [Selectors]([selectors](https://github.com/ipld/specs/blob/master/selectors/selectors.md)) to address and retrieve subsets of linked data graphs. IPLD [Paths](https://github.com/ipld/specs/blob/master/data-model-layer/paths.md) are also used to address specific nodes within a linked data structure.

IPLD provides a consistent and coherent abstraction above data that allows Filecoin to build and interact with complex, multi-block data structures, such as its HAMT and Sharray <!-- TODO: links, correct names, other multi-block structures, AMT ? -->.

### IpldStores

The Filecoin network relies primarily on two distinct IPLD GraphStores:

- One `ChainStore` which stores the blockchain, including block headers, associated messages, etc.
- One `StateStore` which stores the payload state from a given blockchain, or the `stateTree` resulting from all block messages in a given chain being applied to the genesis state by the {{<sref sys_vm "Filecoin VM">}}.

The `ChainStore` is downloaded by a node from its peers during the bootstrapping phase of {{<sref chain_sync>}} and stored by the node thereafter. It is updated on every new block reception, or if the node syncs to a new best chain.

The `StateStore` is computed through the execution of all block messages in a given `ChainStore` and stored by the node thereafter. It is updated with every new incoming block's processing by the {{<sref vm_interpreter>}} and referenced accordingly by new blocks produced atop it in the block {{<sref block "block header">}}'s `ParentState` field.

TODO:

- What is an IpldStore
  - local storage of dags
- How to use IpldStores in filecoin
  - pass it around
- One ipldstore or many
  - temporary caches
  - intermediately computed state
- Garbage Collection

