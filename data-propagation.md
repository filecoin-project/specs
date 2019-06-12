# Filecoin Data Propagation

The filecoin network needs to broadcast blocks and messages to all peers in the network. This document details how that process works.

Messages and block headers along side the message references are propagated using the [gossipsub libp2p pubsub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). Every full node must implement and run that protocol. All pubsub messages are authenticated and must be [syntactically validated](./validation.md#syntactical-validation) before being propagated further.

Further more, every full node must implement and offer the bitswap protocol and provide all Cid Referenced objects, it knows of, through it. This allows any node to fetch missing pieces (e.g. `Message`) from any node it is connected to. However, the node should fan out these requests to multiple nodes and not bombard any single node with too many requests at a time. A node may implement throttling and DDoS protection to prevent such a bombardment.

## Bitswap

Run bitswap to fetch and serve data (such as blockdata and messages) to and from other filecoin nodes. This is used to fill in missing bits during block propagation, and also to fetch data during sync.

There is not yet an official spec for bitswap, but [the protobufs](https://github.com/ipfs/go-bitswap/blob/master/message/pb/message.proto) should help in the interim.


## Block Propagation

Blocks are propagated over the libp2p pubsub channel `/fil/blocks`. The following structure is filled out with the appropriate information, serialized (with CBOR-RPC), and sent over the wire:

```go
type BlockMsg struct {
  Header Block
  Messages []Cid
}
```

The array of message cids must match the `Messages` field in the block when used to construct a [sharray](sharray.md).

Every `BlockMsg` received must be validated [through the syntactical check](./validation.md#syntactical-validation) before being propagated again. If validation fails, it must not be propagated.


## Message Propagation

Messages are propagated over the libp2p pubsub channel `/fil/messages`. On this channel, every [serialised `Message`](data-structures.md#messages) is announced.

Upon receiving the message, its validity must be checked: the signature must be valid, and the account in question must have enough funds to cover the actions specified. If the message is not valid it should be dropped and must not be forwarded.

{{% notice todo %}}
discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.
{{% /notice %}}
