---
title: Storage Market Actor
---

`StorageMarketActor` is responsible for processing and managing on-chain deals. This is also the entry point of all storage deals and data into the system. It maintains a mapping of `StorageDealID` to `StorageDeal` and keeps track of locked balances of `StorageClient` and `StorageProvider`. When a deal is posted on chain through the `StorageMarketActor`, it will first check if both transacting parties have sufficient balances locked up and include the deal on chain. On every successful submission of `PoStProof`, `StorageMarketActor` will credit the `StorageProvider` a fraction of the storage fee based on how many blocks have passed since the last `PoStProof`. In the event that there are sectors included in the `FaultSet`, `StorageMarketActor` will fetch deal information from the chain and `SlashStorageFault` for faulting on those deals. Similarly, when a `PoStProof` is missed by the end of a `ProvingPeriod`, `SlashStorageFault` will also be called by the `CronActor` to penalize `StorageProvider` for dropping a `StorageDeal`.

{{< readfile file="storage_market_actor.id" code="true" lang="go" >}}
