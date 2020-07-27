---
title: "BlockSync"
weight: 5
dashboardWeight: 1
dashboardState: wip
dashboardAudit: 0
dashboardTests: 0
---

# BlockSync
---

{{< hint info >}}
**Name**: Block Sync  
**Protocol ID**: `/fil/sync/blk/0.0.1`
{{< /hint >}}

The blocksync protocol is a small protocol that allows Filecoin nodes to request ranges of blocks from each other. It is a simple request/response protocol.

The request requests a chain of a given length by the hash of its highest block. The `Options` allow the requester to specify whether or not blocks and messages are to be included.

The response contains the requested chain in reverse iteration order. Each item in the `Chain` array contains the blocks for that tipset if the `Blocks` option bit in the request was set, and if the `Messages` bit was set, the messages across all blocks in that tipset. The `MsgIncludes` array contains one array of integers for each block in the `Blocks` array. Each of the arrays in `MsgIncludes` contains a list of indexes of messages from the `Messages` array that are in each `Block` in the blocks array.

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
