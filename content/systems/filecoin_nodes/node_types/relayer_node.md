---
title: Relayer Node
weight: 7
dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Relayer Node

```go
type RelayerNode interface {
  FilecoinNode

  blockchain.MessagePool
  markets.MarketOrderBook
}
```
