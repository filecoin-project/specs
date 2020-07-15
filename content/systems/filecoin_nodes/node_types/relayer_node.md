---
title: Relayer Node
weight: 7
---

# Relayer Node
---

```go
type RelayerNode interface {
  FilecoinNode

  blockchain.MessagePool
  markets.MarketOrderBook
}
```
