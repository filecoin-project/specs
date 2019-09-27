---
menuTitle: Retrieval Market
title: "Retrieval Market in Filecoin"
entries:
- retrieval_client
- retrieval_provider
---

## Components

Version 0 of the `retrieval market` protocol is what we (tentatively) will launch the filecoin network with. It is version zero because it will only be good enough to fit the bill as a way to pay another node for a file.

The main components are as follows:

- A payment channel actor (See [payment channels](payment-channels.md) for details)
- 'retrieval-v0' `libp2p` services
- A chain-based content routing interface
- A set of commands to interact with the above

## Retrieval V0 `libp2p` Services

The v0 `retrieval market` will initially be implemented as two `libp2p` services. It will be request response based, where the client who is requesting a file sends a `retrieval deal proposal` to the miner. The miner chooses whether or not to accept it, sends their response which (if they accept the proposal) includes a `signed retrieval deal`, followed by the actual requested content, streamed as a series of bitswap block messages, using a pre-order traversal of the dag. Each block should use the [bitswap block message format](https://github.com/ipfs/go-bitswap/blob/c980d7ed36f278e93828acf920f3a911e8263265/message/message.go#L228). This way, the client should be able to verify the data incrementally as it receives it. Once the client has received all the data, it should then send a payment channel SpendVoucher of the proposed amount to the miner. This protocol may be easily extended to include payments from the client to the miner every N blocks, but for now we omit that feature.
