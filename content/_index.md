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

The Filecoin protocol uses cryptographic proofs to ensure the two following guarantees:

- **Storage Based Consensus**: Miners' power in the consensus is proportional to their amount of storage. Miner increase their power by proving that they are dedicating unique storage to the network.
- **Verifiable Storage Market**: Miners must be proving that they are dedicating unique physical space for each copy of the clients data through a period of time.

We briefly describe the three main Filecoin proofs:

- __*Proof-of-Replication (PoRep)*__ proves that a Storage Miner is dedicating unique dedicated storage for each ***sector***. Filecoin Storage Miners collect new clients' data in a sector, run a slow encoding process (called `Seal`) and generate a proof (`SealProof`) that the encoding was generated correctly.

  In Filecoin, PoRep provides two guarantees: (1) *space-hardness*: Storage Miners cannot lie about the amount of space they are dedicating to Filecoin in order to gain more power in the consensus; (2) *replication*: Storage Miners are dedicating unique storage for each copy of their clients data.

- __*Proof of Spacetime*__ proves that an arbitrary number of __*sealed sectors*__ existed over a specified period of time in their own dedicated storage — as opposed to being generated on-the-fly at proof generation time.

- __*Piece-Inclusion-Proof*__ proves that a given __*piece*__ is contained within a specified __*sealed sector*__.



## Filecoin VM

The majority of Filecoin's user facing functionality (payments, storage market, power table, etc) is managed through the Filecoin State Machine. The network generates a series of blocks, and agrees which 'chain' of blocks is the correct one. Each block contains a series of state transitions called `messages`, and a checkpoint of the current `global state` after the application of those `messages`.

The `global state` here consists of a set of `actors`, each with their own private `state`.

An `actor` is the Filecoin equivalent of Ethereum's smart contracts, it is essentially an 'object' in the filecoin network with state and a set of methods that can be used to interact with it. Every actor has a Filecoin balance attributed to it, a `state` pointer, a `code` CID which tells the system what type of actor it is, and a `nonce` which tracks the number of messages sent by this actor. (TODO: the nonce is really only needed for external user interface actors, AKA `account actors`. Maybe we should find a way to clean that up?)

There are two routes to calling a method on an `actor`. First, to call a method as an external participant of the system (aka, a normal user with Filecoin) you must send a signed `message` to the network, and pay a fee to the miner that includes your `message`.  The signature on the message must match the key associated with an account with sufficient Filecoin to pay for the messages execution. The fee here is equivalent to transaction fees in Bitcoin and Ethereum, where it is proportional to the work that is done to process the message (Bitcoin prices messages per byte, Ethereum uses the concept of 'gas'. We also use 'gas').

Second, an `actor` may call a method on another actor during the invocation of one of its methods.  However, the only time this may happen is as a result of some actor being invoked by an external users message (note: an actor called by a user may call another actor that then calls another actor, as many layers deep as the execution can afford to run for).


For full implementation details, see the [VM Interpreter spec]({{<ref "/docs/impl/vm_interpreter">}}).



## Protocol Entities

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

  - In contrast, we use the term _static routine_ for a routine that only has access to its own internal state,
    such as a local daemon maintaining a key-value store. The following are the static routines in the Filecoin protocol:

      - [VM Interpreter]({{<ref "docs/impl/vm_interpreter.md">}})

- **_APIs_** are messages that can be sent to routines. A client's view of a given sub-protocol, such as a request to a miner node's [Storage Provider]({{<ref "docs/impl/storage_provider.md">}}) to store files in the storage market, may require the execution of a series of APIs.

- **_Nodes_** are complete software and hardware systems that interact with the protocol.
  A node might be constantly running several of the above _routines_, and exposing one or more _APIs_ locally and/or over the network,
  depending on the node configuration.
  The term _full node_ refers to an system that runs all of the above routines, and supports all of the APIs detailed in the spec.

- **_Actors_** are virtual entities embodied in the state of the Filecoin VM.
  Protocol actors are analogous to participants in smart contracts;
  an actor carries a FIL balance and can interact with other actors
  via the operations of the VM, but does not necessarily correspond to any particular software component.
