---
title: Relayer Node
---

```
type RelayerNode interface {
  FilecoinNode

  blockchain.MessagePool
  markets.MarketOrderBook
}
```
