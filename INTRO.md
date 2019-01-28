# Filecoin Spec

This is the official Filecoin protocol specification.

This collection of pages specify the protocols comprising the Filecoin network: storage, mining, retrieval, payments, and blockchain. It also specifies consensus rules for blockchain state, about which nodes must be in agreement to participate meaningfully. The goal of these specs is to provide sufficient detail for another implementation to be fully compatible with `go-filecoin` peers, but this is still a work in progress.

The repository also includes some implementation notes from the `go-filecoin` reference implementation.


### Style
Any content that is written with `code ticks` has a specific definition to Filecoin and is defined in [the glossary](definitions.md).

Many sections of the spec use go type notation to describe the functionality of certain components. This is entirely a style preference by the authors and does not imply in any way that one must use go to implement Filecoin. 

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.


## Overview

Filecoin is a distributed storage network, with shared state persisted to a blockchain. 

The network maintains consensus over the current state of a replicated [state machine](state-machine.md) through [expected consensus](expected-consensus.md). This replicated state machine is used to run the [filecoin storage market](storage-market.md). This market provides a place to buy and sell storage within the distributed network of filecoin miners. The market also provides the needed mechanisms to ensure that the data stored by the network is actually being stored as promised, without requiring client interaction.

Clients interact with the system by [sending messages](data-propogation.md#message-propogation) to the network. These messages are gathered up and included by miners in blocks. Each of these messages defines a state transition in the state machine. The simplest messages say something like "move Filecoin from *this* account under my control to *this* other account", but more complex ones describe storage sector commitments, storage deals struck, and proofs of storage.

The Filecoin protocol is really a suite of protocols, including:
- the chain protocol for sharing the block structures that constitute the blockchain,
- the [block mining](mining.md) protocol for producing new blocks,
- the [consensus mechanism](expected-consensus.md) and [rules](validation.md) for agreeing on blockchain state,
- the [storage market](storage-market.md) protocol for storage "miners" to offer storage and clients to take it (persisted to the blockchain),
- the [retrieval market](retrieval-market.md) protocol for retrieving files (off-chain),
- the [payments](payments.md) channel protocol for transferring FIL tokens between actors (mostly off-chain). 

### Message Transport

Filecoin uses [libp2p](https://libp2p.io) for all network communications. libp2p provides transport-agnostic services for peer discovery, naming, routing, pubsub channels and a distributed record store, and there are full or partial [implementations](https://libp2p.io/implementations/) in a number of languages. These specifications assume the use of libp2p and its servicesand do not specify transport-level details. It is not strictly necessary to use an implementation of libp2p, only to provide wire-compatible implementation of the libp2p services used by Filecoin; refer to the [libp2p specs](https://github.com/libp2p/specs).

The Filecoin messages and state machine use [IPLD](https://ipld.io) for data representation and serialization. IPLD provides a canonical model for content-addressed data structures, providing a representation of basic data objects and links between them. IPLD objects have standard representations as serial data (via [CBOR](http://cbor.io/) as well as JSON and other representation schemes.

