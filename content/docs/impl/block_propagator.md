# Block Receiver

{{<js>}}

- {{% todo %}} finish {{% /todo %}}




## Filecoin Data Propagation

The filecoin network needs to broadcast blocks and messages to all peers in the network. This document details how that process works.

Messages and block headers along side the message references are propagated using the [gossipsub libp2p pubsub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). Every full node must implement and run that protocol. All pubsub messages are authenticated and must be [syntactically validated](./validation.md#syntactical-validation) before being propagated further.

Further more, every full node must implement and offer the bitswap protocol and provide all Cid Referenced objects, it knows of, through it. This allows any node to fetch missing pieces (e.g. `Message`) from any node it is connected to. However, the node should fan out these requests to multiple nodes and not bombard any single node with too many requests at a time. A node may implement throttling and DDoS protection to prevent such a bombardment.

## Bitswap

Run bitswap to fetch and serve data (such as blockdata and messages) to and from other filecoin nodes. This is used to fill in missing bits during block propagation, and also to fetch data during sync.

There is not yet an official spec for bitswap, but [the protobufs](https://github.com/ipfs/go-bitswap/blob/master/message/pb/message.proto) should help in the interim.


## Block Propagation

Blocks are propagated over the libp2p pubsub channel `/fil/blocks`. The following structure is filled out with the appropriate information, serialized (with IPLD), and sent over the wire:

```sh
type BlockMessage struct {
  header BlockHeader
  secpkMessages []&SignedMessage
  blsMessages []&Message
}
```

Every `BlockMessage` received must be validated [through the syntactical check](validation.md#syntactical-validation) before being propagated again. If validation fails, it must not be propagated.


## Message Propagation

Messages are propagated over the libp2p pubsub channel `/fil/messages`. On this channel, every [serialised `SignedMessage`](data-structures.md#messages) is announced.

Upon receiving the message, its validity must be checked: the signature must be valid, and the account in question must have enough funds to cover the actions specified. If the message is not valid it should be dropped and must not be forwarded.

{{% notice todo %}}
discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.
{{% /notice %}}







## BlockSync

- **Name**: Block Sync
- **Protocol ID**: `/fil/sync/blk/0.0.1`

The blocksync protocol is a small protocol that allows Filecoin nodes to request ranges of blocks from each other. It is a simple request/response protocol.

The request requests a chain of a given length by the hash of its highest block. The `Options` allow the requester to specify whether or not blocks and messages to be included.

The response contains the requested chain in reverse iteration order. Each item in the `Chain` array contains the blocks for that tipset if the `Blocks` option bit in the request was set, and if the `Messages` bit was set, the messages across all blocks in that tipset. The `MsgIncludes` array contains one array of integers for each block in the `Blocks` array. Each of the arrays in `MsgIncludes` contains a list of indexes of messages from the `Messages` array that are in each `Block` in the blocks array.

```sh
type BlockSyncRequest struct {
    ## The TipSet being synced from
	start [&Block]
    ## How many tipsets to sync
	requestLength UInt
    ## Query options
    options Options
}
```

```sh
type Options enum {
    # Include only blocks
    | Blocks 0
    # Include only messages
    | Messages 1
    # Include messages and blocks
    | BlocksAndMessages 2
}

type BlockSyncResponse struct {
	chain [TipSetBundle]
	status Status
}

type TipSetBundle struct {
  blocks [Blocks]
  secpMsgs [SignedMessage]
  secpMsgIncludes [[UInt]]

  blsMsgs [Message]
  blsMsgIncludes [[Uint]]
}

type Status enum {
    ## All is well.
    | Success 0
    ## Sent back fewer blocks than requested.
    | PartialResponse 101
    ## Request.Start not found.
    | BlockNotFound 201
    ## Requester is making too many requests.
    | GoAway 202
    ## Internal error occured.
    | InternalError 203
    ## Request was bad
    | BadRequest 204
}
```

### Example

The TipSetBundle

```
Blocks: [b0, b1]
secpMsgs: [mA, mB, mC, mD]
secpMsgIncludes: [[0, 1, 3], [1, 2, 0]]
```

corresponds to:

```
Block 'b0': [mA, mB, mD]
Block 'b1': [mB, mC, mA]
```
