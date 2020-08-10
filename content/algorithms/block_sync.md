---
title: "BlockSync"
weight: 5
dashboardWeight: 1
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# BlockSync
---

{{< hint info >}}
**Name**: Block Sync  
**Protocol ID**: `/fil/sync/blk/0.0.1`
{{< /hint >}}

BlockSync is a simple request/response protocol that allows Filecoin nodes to request ranges of blocks from each other, when they have run out of sync, e.g., during downtime. Given that the Filecoin blockchain is extended in Tipsets (i.e., groups of blocks), rather than in blocks, the BlockSync protocol should also operate in terms of Tipsets.

The request message requests a chain segment of a given length by the hash of its highest block. It is worth noting that this does not necessarily apply to the head (i.e., latest tipset) of the current chain, but it can also apply to earlier segments. For example, if the current height is at 5000, but a node is missing blocks between 4500-4700, then the `Head` requested is 4700 and the length is 200.

The `Options` allow the requester to specify whether they want to receive block headers only, the transaction messages included in every block or both.

The response contains the requested chain segment in reverse iteration order. Each item in the `Chain` array contains either the block headers for that tipset if the `Blocks` option bit in the request was set, or the messages across all blocks in that tipset, if the `Messages` bit was set, or both, if both option bits were set.

The `MsgIncludes` array contains one array of integers for each block in the `Blocks` array. Each of the `Blocks` arrays in `MsgIncludes` contains a list of message indexes from the `Messages` array that are in each `Block` in the blocks array.

If not all tipsets requested could be fetched, but the `Head` of the chain segment requested was fetched, then this is not considered an error, given that the node can continue extending the chain from the `Head` onwards.

```go
type BlockSyncRequest struct {
    ## The TipSet being synced from
	start [&Block]
    ## How many tipsets to sync
	requestLength UInt
    ## Query options
    options Options
}
```

```go
type Options enum {
    # Include only blocks
    | Blocks 1
    # Include only messages
    | Messages 2
    # Include messages and blocks
    | BlocksAndMessages 3
}

type BlockSyncResponse struct {
	chain [TipSetBundle]
	status Status
}

type TipSetBundle struct {
  blocks [Blocks]

  blsMsgs [Message]
  blsMsgIncludes [[Uint]]

  secpMsgs [SignedMessage]
  secpMsgIncludes [[UInt]]
}

type Status enum {
    ## All is well.
    | Success 0
    ## We could not fetch all blocks requested (but at least we returned
	## the `Head` requested). Not considered an error.
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

## Example

The TipSetBundle

```js
Blocks: [b0, b1]
secpMsgs: [mA, mB, mC, mD]
secpMsgIncludes: [[0, 1, 3], [1, 2, 0]]
```

corresponds to:

```js
Block 'b0': [mA, mB, mD]
Block 'b1': [mB, mC, mA]
```
