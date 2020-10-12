---
title: "ChainExchange"
weight: 5
dashboardWeight: 1
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# ChainExchange

{{< hint info >}}
**Name**: Chain Exchange
**Protocol ID**: `/fil/chain/xchg/0.0.1`
{{< /hint >}}

ChainExchange is a simple request/response protocol that allows Filecoin nodes to request ranges of Tipsets and/or Messages from each other.

The `Request` message requests a chain segment of a given length by the hash of the highest Tipset in the segment (not necessarily heaviest Tipset of the current chain). For example, if the current height is at 5000, but a node is missing Tipsets between 4500-4700, then the `Request.Head` requested is 4700 and `Request.Length` is 200.

The `Options` allow the requester to specify whether they want to receive block headers of the Tipsets only, the transaction messages included in every block, or both.

The `Response` contains the requested chain segment in reverse iteration order. Each item in the `Chain` array contains either the block headers for that Tipset if the `Blocks` option bit in the request is set, or the messages across all blocks in that Tipset, if the `Messages` bit is set, or both, if both option bits are set.

Each `CompactedMessages` structure contains the BLS and `secp256k1` messages for the corresponding Tipset in unified arrays that encode each message once regardless of how many times it appears in the Tipset's blocks (to reduce the `Response` size). The mapping between message and the blocks into which they are included is encoded in the `BlsIncludes` and `SecpkIncludes` arrays. These arrays are indexed by the block index in the Tipset and contain in that position an array of message indexes that belong to that block. The messages themselves are in the `Bls` and `Secpk` arrays.

If not all Tipsets requested could be fetched, but the `Head` of the chain segment requested was successfully fetched (and potentially more contiguous Tipsets), then this is not considered an error and a `Partial` response code is returned. The node can continue extending the chain from the partial returned segment onwards.

```go
type Request struct {
    ## Head of the requested segment (block CIDs comprising the entire Tipset)
	Head [Cid]
    ## Length of the requested segment
	Length UInt
    ## Query options
    Options UInt
}
```

```go
type Options enum {
    # Include block headers
    | Headers 1
    # Include block messages
    | Messages 2
}

type Response struct {
    ## Response Status
    Status status
    ## Optional error message
    ErrorMessage string
    ## Returned segment containing block messages and/or headers
    Chain []*BSTipSet
}

type Status enum {
    ## Success: the entire segment requested was fetched.
    | Ok 0
    ## We could not fetch all Tipsets requested but at least we returned
    ## the `Head` requested (and potentially more contiguous Tipsets).
    ## Not considered an error.
    | Partial 101
    ## `Request.Head` not found.
    | NotFound 201
    ## Requester is making too many requests.
    | GoAway 202
    ## Internal error occurred.
    | InternalError 203
    ## Request was badly formed.
    | BadRequest 204
}

type CompactedMessages struct {
  ## Array of BLS messages in this tipset.
  Bls [Message]
  ## Array of messages indexes present in each block:
  ## `BlsIncludes[BI] -> [MI]`
  ##  * BI: block index in the tipset.
  ##  * MI: message index in `Bls` list.
  BlsIncludes [[Uint]]

  ## Array of `secp256k1` messages.
  Secpk [SignedMessage]
  ## Inclusion array, see `BlsIncludes`.
  SecpkIncludes [[UInt]]
}
```

## Example

For the set of arrays in the following `BSTipSet`, the corresponding messages per block are as shown below:

**BSTipSet**
```js
Blocks: [b0, b1]
CompactedMessages
```

**CompactedMessages**
```js
Secpk: [mA, mB, mC, mD]
SecpkIncludes: [[0, 1, 3], [1, 2, 0]]
```

**Messages corresponding to each Block**
```js
Block 'b0': [mA, mB, mD]
Block 'b1': [mB, mC, mA]
```

In other words, the first element of the `SecpkIncludes` array: `[0, 1, 3]` points to the 1st, 2nd and 4th element of the `Secpk` array: `mA, mB, mD`, which correspond to the 1st element of the `Blocks` array `b0`. Hence, `Block 'b0': [mA, mB, mD]`.

Similarly, the second element of the `SecpkIncludes` array: `[1, 2, 0]` points to the 2nd, 3rd and 1st element of the `Secpk` array: `mB, mC, mA`, which correspond to the 2nd element of the `Blocks` array `b1`. Hence, `Block 'b1': [mB, mC, mA]`.

See [Lotus](https://github.com/filecoin-project/lotus/tree/master/chain/exchange) for an example implementation of the protocol.
