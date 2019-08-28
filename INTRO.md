---
title: Introduction
type: docs
---

# Filecoin Spec

{{% notice warning %}}
**Warning:** This is the official Filecoin protocol specification. It is a work in progress. While reading, if you notice any discrepancies, or issues, please open an issue on the [specs repo](https://github.com/filecoin-project/specs).
{{% /notice %}}

<!-- the below comment is replaced with the pdf link on site builds -->
You can also download the full spec in [PDF format](./).

This collection of pages specify the protocols comprising the Filecoin network. The goal of these specs is to provide sufficient detail that another implementation, written using only this document as reference, can be fully compatible with `go-filecoin` peers, but this is still a work in progress.

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

### Message Transport

Filecoin uses [libp2p](https://libp2p.io) for all network communications. libp2p provides transport-agnostic services for peer discovery, naming, routing, pubsub channels and a distributed record store, and there are full or partial [implementations](https://libp2p.io/implementations/) in a number of languages. This spec assumes the use of libp2p and its services and does not specify transport-level details. That said, to be at least minimally compatible with other Filecoin nodes, it must support at least the [mplex](https://github.com/libp2p/specs/tree/master/mplex) stream multiplexer, and the [secio](TODO, spec in PR) encrypted transport protocols. For more details on the exact wire protocol of libp2p, refer to the [libp2p specs](https://github.com/libp2p/specs).

Filecoin uses [IPLD](https://ipld.io) for the representation and serialization of the majority of the data in the system. IPLD provides a canonical model for content-addressed data structures, providing a representation of basic data objects and links between them.
