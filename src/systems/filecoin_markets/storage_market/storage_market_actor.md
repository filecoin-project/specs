---
title: Storage Market Actor
---

`StorageMarketActor` is responsible for processing and managing on-chain deals. This is also the entry point of all storage deals and data into the system. It maintains a mapping of `StorageDealID` to `StorageDeal` and keeps track of locked balances of `StorageClient` and `StorageProvider`. When a deal is posted on chain through the `StorageMarketActor`, it will first check if both transacting parties have sufficient balances locked up and include the deal on chain. On every successful submission of `PoStProof`, `StorageMarketActor` will credit the `StorageProvider` a fraction of the storage fee based on how many blocks have passed since the last `PoStProof`. In the event that there are sectors included in the `FaultSet`, `StorageMarketActor` will fetch deal information from the chain and `SlashStorageFault` for faulting on those deals. Similarly, when a `PoStProof` is missed by the end of a `ProvingPeriod`, `SlashStorageFault` will also be called by the `CronActor` to penalize `StorageProvider` for dropping a `StorageDeal`.

(You can see the _old_ Storage Market Actor [here](docs/systems/filecoin_markets/storage_market/storage_market_actor_old) )

# `StorageMarketActor` interface

{{< readfile file="storage_market_actor.id" code="true" lang="go" >}}

# `StorageMarketActor` implementation

{{< readfile file="storage_market_actor.go" code="true" lang="go" >}}


{{<label storage_deal_collateral>}}
# Storage Deal Collateral

Storage Deals have an associated collateral amount. This `StorageDealCollateral` is held in the `StorageMarketActor`.
Its value is agreed upon by the storage provider and client off-chain, but must be greater than a protocol-defined minimum in any deal. Storage providers will choose to offer greater collateral to signal high-quality storage to clients.

On `SectorFailureTimeout` (see {{<sref faults>}}), the `StorageDealCollateral` will be burned. In the future, the Filecoin protocol may be amended to send up to half of the collateral to storage clients as damages in such cases.

Upon graceful deal expiration, storage providers must wait for finality number of epochs (as defined in {{<sref finality>}}) before being able to withdraw their `StorageDealCollateral` from the `StorageMarketActor`.
