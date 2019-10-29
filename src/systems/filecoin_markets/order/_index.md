---
menuTitle: Orders
statusIcon: 🛑
title: Market Orders - Asks and Bids
---

TODO:

- Write asks
- Write bids
- Write how market orders propagate (gossipsub)

{{< readfile file="order.id" code="true" lang="go" >}}

# Verifiability

TODO:

- write what parts of market orders are verifiable, and how
  - eg: miner storage ask could carry the amount of storage available (which should be at mot (pledge - sectors sealed))
  - eg: client storage bid price could be checked against available money in the StorageMarket
