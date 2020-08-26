---
title: Storage Miner Node
weight: 5
dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Storage Miner Node

```go
type StorageMinerNode interface {
  FilecoinNode

  systems.Blockchain
  systems.Mining
  markets.StorageMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
