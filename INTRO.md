# Filecoin Spec

This is the official Filecoin protocol specification. It is a work in progress. While reading, if you notice any discrepancies, or issues, please open an issue on the [specs repo](https://github.com/filecoin-project/specs).

Any content that is written with `code ticks` has a specific definition to Filecoin and is defined in [the glossary](definitions.md).

Many sections of the spec use go type notation to describe the functionality of certain components. This is entirely a style preference by the authors and does not imply in any way that one must use go to implement Filecoin. 


The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119.


## Architecture

Filecoin is a distributed storage network operating over a blockchain. It uses
[libp2p](https://libp2p.io) for all network communications, and uses
[IPLD](https://ipld.io) for dealing with data. 

The network maintains consensus over the current state of a replicated [state machine](state-machine.md) through [expected consensus](expected-consensus.md). This replicated state machine is used to run the [filecoin storage market](storage-market.md). This market provides a place to buy and sell storage within the distributed network of filecoin miners. The market also provides the needed mechanisms to ensure that the data stored by the network is actually being stored as promised, without requiring client interaction.

Clients interact with the system by [sending messages](data-propogation.md#message-propogation) to the network. These messages are gathered up and included by miners in blocks. Each of these messages defines a state transition in the state machine. The simplest ones say just "move Filecoin from *this* account under my control to *this* other account".

