---
title: Retrieval Miner Node
weight: 6
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Retrieval Miner Node
---

```go
type RetrievalMinerNode interface {
  FilecoinNode

  blockchain.Blockchain
  markets.RetrievalMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
