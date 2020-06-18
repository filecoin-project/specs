---
menuTitle: IpldStore
title: "IpldStore - Local Storage for hash-linked data"
---

{{< readfile file="../../../../libraries/ipld/ipld.id" code="true" lang="go" >}}

IPLD is a set of libraries which allow for the interoperability of content-addressed data structures present in most distributed systems. It provides a fundamental 'common language' to primitive cryptographic hashing, enabling data to be verifiably referenced and retrieved between two independent protocols. For example, a user can reference a git commit in a blockchain transaction to create an immutible copy and timestamp, or a data from a DHT can be refenced and linked to in a smart contract. 

For a more in-depth explanation of IPLD's role in Filecoin, see [this section] of the specification.

## IPLD in filecoin

On the Filecoin network, IPLD is used in two ways:

- All system datastructures are stored in [IPLD](https://ipld.io) format, a data format akin to JSON but designed for storage, retrieval and traversal of hash-linked data DAGs.
- Files and data stored on the Filecoin network may also be stored in IPLD format. While this is not required, it offers the benefit of supporting [selectors](https://github.com/ipld/specs/blob/master/selectors/selectors.md) to retrieve a smaller subset of the total stored data, as opposed to inefficiently downloading the data set entirely.

IPLD handles marshalling and unmarshalling via a suite of codecs. IPLD codecs may support the complete Data Model or part of the Data Model. Two codecs that support the complete Data Model are [DAG-CBOR](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-cbor.md) and [DAG-JSON](https://github.com/ipld/specs/blob/master/block-layer/codecs/dag-json.md). 

Filecoin uses the **DAG-CBOR** codec for the serialization and deserialization of it's data structures and interacts with that data using the IPLD Data Model, upon which various tools are built. IPLD [Paths](https://github.com/ipld/specs/blob/master/data-model-layer/paths.md) are also used to address specific nodes within a linked data structure.

<-- Inclusion of IPLD 'kinds' here, or irrelevant? -->

### IpldStores

<-- to be expanded on -->

The Filecoin network relies primarily on two distinct IPLD GraphStores:

- One `ChainStore` which stores the blockchain, including block headers, associated messages, etc.
- One `StateStore` which stores the payload state from a given blockchain, or the `stateTree` resulting from all block messages in a given chain being applied to the genesis state by the {{<sref sys_vm "Filecoin VM">}}.

The `ChainStore` is downloaded by a node from their peers during the bootstrapping phase of {{<sref chain_sync>}} and is stored by the node thereafter. It is updated on every new block reception, or if the node syncs to a new best chain.

The `StateStore` is computed through the execution of all block messages in a given `ChainStore` and is stored by the node thereafter. It is updated with every new incoming block's processing by the {{<sref vm_interpreter>}}, and referenced accordingly by new blocks produced atop it in the block {{<sref block "block header">}}'s `ParentState` field.

TODO:

- What is an IpldStore
  - local storage of dags
- How to use IpldStores in filecoin
  - pass it around
- One ipldstore or many
  - temporary caches
  - intermediately computed state
- Garbage Collection
