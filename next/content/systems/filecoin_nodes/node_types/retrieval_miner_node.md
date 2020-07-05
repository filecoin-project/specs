---
title: Retrieval Miner Node
---

# Retrieval Miner Node
---

```
type RetrievalMinerNode interface {
  FilecoinNode

  blockchain.Blockchain
  markets.RetrievalMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
