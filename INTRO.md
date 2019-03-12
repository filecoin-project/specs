# Intro

This is the official Filecoin protocol specification. It is a work in progress. While reading, if you notice any discrepancies, or issues, please open an issue on the [specs repo](https://github.com/filecoin-project/specs).

This collection of pages specify the protocols comprising the Filecoin network. The goal of these specs is to provide sufficient detail that another implementation, written using only this document as reference, can be fully compatible with `go-filecoin` peers, but this is still a work in progress.

## Style
Any content that is written with `code ticks` has a specific definition to Filecoin and is defined in the [[#glossary]].

Many sections of the spec use go type notation to describe the functionality of certain components. This is entirely a style preference by the authors and does not imply in any way that one must use go to implement Filecoin.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.


## Overview

Filecoin is a distributed storage network, with shared state persisted to a blockchain.

The network maintains consensus over the current state of a replicated [[#state-machine]] through [[#expected-consensus]]. This replicated state machine is used to run [[#storage-market]]. This market provides a place to buy and sell storage within the distributed network of filecoin miners. The market also provides the needed mechanisms to ensure that the data stored by the network is actually being stored as promised, without requiring client interaction.

Clients interact with the system by [sending messages](data-propagation.md#message-propagation) to the network. These messages are gathered up and included by miners in blocks. Each of these messages defines a state transition in the state machine. The simplest messages say something like "move Filecoin from *this* account under my control to *this* other account", but more complex ones describe storage sector commitments, storage deals struck, and proofs of storage.

The Filecoin protocol itself is really a suite of protocols, including:
- the chain protocol for propagating the data that constitutes the blockchain
- the [block mining](mining.md) protocol for producing new blocks
- the [consensus mechanism](expected-consensus.md) and [rules](validation.md) for agreeing on canonical blockchain state
- the [storage market](storage-market.md) protocol for `storage miners` to sell storage and clients to purchase it
- the [retrieval market](retrieval-market.md) protocol for retrieving files
- the [payment channel](payments.md) channel protocol for transferring FIL tokens between actors

## Message Transport

Filecoin uses [libp2p](https://libp2p.io) for all network communications. libp2p provides transport-agnostic services for peer discovery, naming, routing, pubsub channels and a distributed record store, and there are full or partial [implementations](https://libp2p.io/implementations/) in a number of languages. This spec assumes the use of libp2p and its services and does not specify transport-level details. For more details on the exact wire protocol of libp2p, refer to the [libp2p specs](https://github.com/libp2p/specs).

Filecoin uses <a href="https://ipld.io">IPLD</a> for the representation and serialization of the majority of the data in the system. IPLD provides a canonical model for content-addressed data structures, providing a representation of basic data objects and links between them.
