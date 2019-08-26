---
title: "Filecoin Spec"
type: docs
---

{{<js>}}

# Filecoin Spec

{{% notice warning %}}
**Notice:** This is the official Filecoin protocol specification.
The goal of this document is to provide sufficient detail for other implementations to be fully compatible with [`go-filecoin`](https://github.com/filecoin-project/go-filecoin) peers, but this is still a work in progress.
If you notice any errors or discrepancies, please open an issue in the [specs repository](https://github.com/filecoin-project/specs).
{{% /notice %}}

## Overview

Filecoin is a distributed storage network, with shared state persisted to a blockchain.

The network maintains consensus over the current state of a replicated [state machine](state-machine.md) through [expected consensus](expected-consensus.md). This replicated state machine is used to run the [filecoin storage market](storage-market.md). This market provides a place to buy and sell storage within the distributed network of filecoin miners. The market also provides the needed mechanisms to ensure that the data stored by the network is actually being stored as promised, without requiring client interaction.

Clients interact with the system by [sending messages](data-propagation.md#message-propagation) to the network. These messages are gathered up and included by miners in blocks. Each of these messages defines a state transition in the state machine. The simplest messages say something like "move Filecoin from *this* account under my control to *this* other account", but more complex ones describe storage sector commitments, storage deals struck, and proofs of storage.

The Filecoin protocol itself is really a suite of protocols, including:

- the chain protocol for propagating the data that constitutes the blockchain
- the [block mining](mining.md) protocol for producing new blocks
- the [consensus mechanism](expected-consensus.md) and [rules](validation.md) for agreeing on canonical blockchain state
- all interacting with the [state machine](state-machine.md) and the [actors](actors.md) running on it:
  - the [storage market](storage-market.md) protocol for `storage miners` to sell storage and clients to purchase it
  - the [retrieval market](retrieval-market.md) protocol for retrieving files
  - the [payment channel](payments.md) channel protocol for transferring FIL tokens between actors




## Implementation Entities

For clarity, we introduce the following types of entities to describe implementations of the Filecoin protocol:

- **_Data structures_** are collections of semantically-tagged data members (e.g., structs, interfaces, or enums).

- **_Functions_** are computational procedures that do not depend on external state (i.e., mathematical functions,
  or programming language functions that do not refer to global variables).

- **_Routines_** are processes or threads, constantly running some specified main loop, which can send and receive messages.
  (The term "routine" alludes to the concept of a "Goroutine" in Golang.)

  - Within this category, we use the term _dynamic routine_ for a routine that can exchange messages with the outside world,
    including other routines and potentially the Filecoin network.
    The following are the dynamic routines in the Filecoin protocol:

      - [Storage Provider]({{<ref "docs/impl/storage_provider.md">}})
      - [Retrieval Provider]({{<ref "docs/impl/retrieval_provider.md">}})
      - [Block Miner]({{<ref "docs/impl/block_miner.md">}})
      - [Block Propagator]({{<ref "docs/impl/block_propagator.md">}})
      - [Chain Manager]({{<ref "docs/impl/chain_manager.md">}})
      - [Chain Verifier]({{<ref "docs/impl/chain_verifier.md">}})

  - In contrast, we use the term _static routine_ for a routine that only has access to its own internal state,
    such as a local daemon maintaining a key-value store. The following are the static routines in the Filecoin protocol:

      - [VM Interpreter]({{<ref "docs/impl/vm_interpreter.md">}})
      - [Sector Manager]({{<ref "docs/impl/sector_manager.md">}})

- **_APIs_** are messages that can be sent to routines. A client's view of a given sub-protocol, such as a request to a miner node's [Storage Provider]({{<ref "docs/impl/storage_provider.md">}}) to store files in the storage market, may require the execution of a series of APIs.

- **_Nodes_** are complete software and hardware systems that interact with the protocol.
  A node might be constantly running several of the above _routines_, and exposing one or more _APIs_ locally and/or over the network,
  depending on the node configuration.
  The term _full node_ refers to an system that runs all of the above routines, and supports all of the APIs detailed in the spec.

- **_Actors_** are virtual entities embodied in the state of the Filecoin VM.
  Protocol actors are analogous to participants in smart contracts;
  an actor carries a FIL balance and can interact with other actors
  via the operations of the VM, but does not necessarily correspond to any particular software component.
