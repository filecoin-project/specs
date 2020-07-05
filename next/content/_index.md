---
title: Home
---

# Introduction
---

Filecoin is a distributed storage network based on a blockchain mechanism.
Filecoin *miners* can elect to provide storage capacity for the network, and thereby
earn units of the Filecoin cryptocurrency (FIL) by periodically producing
cryptographic proofs that certify that they are providing the capacity specified.
In addition, Filecoin enables parties to exchange FIL currency
through transactions recorded in a shared ledger on the Filecoin blockchain.
Rather than using Nakamoto-style proof of work to maintain consensus on the chain, however,
Filecoin uses proof of storage itself: a miner's power in the consensus protocol
is proportional to the amount of storage it provides.

The Filecoin blockchain not only maintains the ledger for FIL transactions and
accounts, but also implements the Filecoin VM, a replicated state machine which executes
a variety of cryptographic contracts and market mechanisms among participants
on the network.
These contracts include *storage deals*, in which clients pay FIL currency to miners
in exchange for storing the specific file data that the clients request.
Via the distributed implementation of the Filecoin VM, storage deals
and other contract mechanisms recorded on the chain continue to be processed
over time, without requiring further interaction from the original parties
(such as the clients who requested the data storage).

## Status Overview

<i class="gg-incorrect"></i> Incorrect
<i class="gg-wip"></i> WIP
<i class="gg-incomplete"></i> Incomplete
<i class="gg-stable"></i> Stable
<i class="gg-permanent"></i> Permanent

{{< dashboard name="Introduction" >}}

## Overview Diagram
{{< details title="TODO" >}}
- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token
{{< /details >}}


{{< svg src="/intro/overview.dot.svg" title="Protocol Overview Diagram" />}}

## Protocol Flow Diagram

{{< svg src="/intro/full-deals-on-chain.mmd.svg" title="Deals on Chain" />}}

## Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{< svg src="/intro/filecoin.dot.svg" title="Protocol Overview Diagram" />}}


## Key Concepts

For clarity, we refer the following types of entities to describe implementations of the Filecoin protocol:

- **_Data structures_** are collections of semantically-tagged data members (e.g., structs, interfaces, or enums).

- **_Functions_** are computational procedures that do not depend on external state (i.e., mathematical functions,
  or programming language functions that do not refer to global variables).

- **_Components_** are sets of functionality that are intended to be represented as single software units
  in the implementation structure.
  Depending on the choice of language and the particular component, this might
  correspond to a single software module,
  a thread or process running some main loop, a disk-backed database, or a variety of other design choices.
  For example, the [ChainSync](./../algorithms/block_sync.md#example) is a component: it could be implemented
  as a process or thread running a single specified main loop, which waits for network messages
  and responds accordingly by recording and/or forwarding block data.

- **_APIs_** are messages that can be sent to components.
  A client's view of a given sub-protocol, such as a request to a miner node's
  [Storage Provider](/missing-link) component to store files in the storage market,
  may require the execution of a series of APIs.

- **_Nodes_** are complete software and hardware systems that interact with the protocol.
  A node might be constantly running several of the above _components_, participating in several _subsystems_,
  and exposing _APIs_ locally and/or over the network,
  depending on the node configuration.
  The term _full node_ refers to a system that runs all of the above components, and supports all of the APIs detailed in the spec.

- **_Subsystems_** are conceptual divisions of the entire Filecoin protocol, either in terms of complete protocols
  (such as the [Storage Market](/storage_market) or [Retrieval Market](/retrieval_market)), or in terms of functionality
  (such as the [VM - Virtual Machine](/sys_vm)). They do not necessarily correspond to any particular node or software component.

- **_Actors_** are virtual entities embodied in the state of the Filecoin VM.
  Protocol actors are analogous to participants in smart contracts;
  an actor carries a FIL currency balance and can interact with other actors
  via the operations of the VM, but does not necessarily correspond to any particular node or software component.