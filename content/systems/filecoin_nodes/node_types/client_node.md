---
title: Client Node
weight: 4
dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Client Node
---

```go
type ClientNode struct {
  FilecoinNode

  systems.Blockchain
  markets.StorageMarketClient
  markets.RetrievalMarketClient
  markets.MarketOrderBook
  markets.DataTransfers
}
```
