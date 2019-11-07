---
menuTitle: Orders
statusIcon: ⚠️
title: Market Orders - Asks and Proposals
---

There are two primary types of market orders:
  - _Asks_ contain the terms on which a miner is willing to provide its services, and are propogated via gossipsub
  - _Proposals_ contain the client's proposed deal details, and are sent directly to a selected miner

{{< readfile file="order.id" code="true" lang="go" >}}

# Verifiability

TODO:

- write what parts of market orders are verifiable, and how
  - eg: miner storage ask could carry the amount of storage available (which should be at most (pledge - sectors sealed))
  - eg: client storage bid price could be checked against available money in the StorageMarket
