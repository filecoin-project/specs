---
menuTitle: IpldStore
title: "IpldStore - Local Storage for hash-linked data"
---

{{< readfile file="../../../../libraries/ipld.id" code="true" lang="go" >}}

Filecoin datastructures are stored in [IPLD](https://ipld.io) format, a data format akin to json built for storage, retrieval and traversal of hash-linked data DAGs.

The Filecoin network relies primarily on two distinct IPLD GraphStores:

- One `ChainStore` which stores the blockchain, including block headers, associated messages, etc.
- One `StateStore` which stores the payload state from a given blockchain, or the `stateTree` resulting from all block messages in a given chain being applied to the genesis state by the {{<sref sys_vm "Filecoin VM">}}.

The `ChainStore` is downloaded by nodes from their peers during the bootstrapping phase of {{<sref chain_sync>}} and stored by the node thereafter. It is updated on every new block reception, or if the node syncs to a new best chain.

The `StateStore` is computed through execution of all block messages in a given `ChainStore` and stored by the node thereafter. It is updated with every new incoming block's processing by the {{<sref vm_interpreter>}} and referenced accordingly by new blocks produced atop it in the block {{<sref block "block header">}}'s `ParentState` field.

TODO:

- What is IPLD
  - hash linked data
  - from IPFS
- Why is it relevant to filecoin
  - all network datastructures are definitively IPLD
  - all local datastructures can be IPLD
- What is an IpldStore
  - local storage of dags
- How to use IpldStores in filecoin
  - pass it around
- One ipldstore or many
  - temporary caches
  - intermediately computed state
- Garbage Collection

