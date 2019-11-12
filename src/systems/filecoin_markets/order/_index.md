---
menuTitle: Orders
statusIcon: ⚠️
title: Market Orders - Asks
---

_Asks_ contain the terms on which a miner is willing to provide its services. They are propogated via gossipsub.

A `StorageAsk` contains basic storage deal terms of price, collateral, and minimum piece size (size of the smallest piece it is willing to store under these terms). It also contains a `Timestamp` for its creation in `ChainEpoch`, a `MaxDuration` for the max duration in `ChainEpoch` that a miner is willing to store under these terms, and a `MinDuration`. If a miner wishes to override an ask, it can issue a new ask with a higher sequence number (`SeqNo`). Clients look at all the `StorageAsks` in a gossip network and decide which miner to contact to enter into a deal. The deal negotiation process happens off chain and the client submits a `StorageDealProposal` to the miner, as detailed in Storage Deals, after an agreement is reached. 


TODO:

- Retrieval asks

{{< readfile file="order.id" code="true" lang="go" >}}

# Verifiability

TODO:

- write what parts of market orders are verifiable, and how
  - eg: miner storage ask could carry the amount of storage available (which should be at most (pledge - sectors sealed))
  - eg: client storage bid price could be checked against available money in the StorageMarket
