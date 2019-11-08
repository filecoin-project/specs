---
menuTitle: Orders
statusIcon: ⚠️
title: Market Orders - Asks
---

_Asks_ contain the terms on which a miner is willing to provide its services. They are propogated via gossipsub.

A `StorageAsk` contains basic storage deal terms of price, collateral, and minimum piece size (size of the smallest piece it is willing to store under these terms). It also contains a `Timestamp` for its creation, and `Expiry` for when the miner will stop accepting new deals under these terms. If a miner wishes to override an ask, it can issue an new ask with a higher sequence number (`SeqNo`).

Clients choose a `StorageAsk` and respond directly to that miner via a `StorageDealProposal`, detailed in Storage Deals.

TODO:
- confirm/clarify `Expiry` is NOT the longest duration the miner is willing to store
- Retrieval asks

{{< readfile file="order.id" code="true" lang="go" >}}

# Verifiability

TODO:

- write what parts of market orders are verifiable, and how
  - eg: miner storage ask could carry the amount of storage available (which should be at most (pledge - sectors sealed))
  - eg: client storage bid price could be checked against available money in the StorageMarket
