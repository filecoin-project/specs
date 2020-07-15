---
title: Retrieval Miner Node
weight: 6
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
