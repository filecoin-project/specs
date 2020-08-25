---
title: PieceStore

dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# PieceStore - Storing and indexing pieces

A `PieceStore` is an object that can store and retrieve pieces
from some local storage. The `PieceStore` additionally keeps
an index of pieces.

{{< embed src="piece_store.id" lang="go" >}}
