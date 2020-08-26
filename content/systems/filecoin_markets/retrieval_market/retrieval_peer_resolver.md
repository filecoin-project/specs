---
title: "Retrieval Peer Resolver"
weight: 1
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Retrieval Peer Resolver

The `peer resolver` is a content routing interface to discover retrieval miners that have a given Piece.

It can be backed by both a local store of previous storage deals or by querying the chain.

```go
// PeerResolver is an interface for looking up providers that may have a piece
type PeerResolver interface {
	GetPeers(payloadCID cid.Cid) ([]RetrievalPeer, error) // TODO: channel
}
```