---
title: Storage Market Actor
---

`StorageMarketActor` is responsible for processing and managing on-chain deals. This is also the entry point of all storage deals and data into the system. It maintains a mapping of `StorageDealID` to `StorageDeal` and keeps track of locked balances of `StorageClient` and `StorageProvider`. When a deal is posted on chain through the `StorageMarketActor`, it will first check if both transacting parties have sufficient balances locked up and include the deal on chain. 

# `StorageMarketActorState` implementation

{{< readfile file="/docs/actors/actors/builtin/storage_market/storage_market_actor_state.go" code="true" lang="go" >}}

# `StorageMarketActor` implementation

{{< readfile file="/docs/actors/actors/builtin/storage_market/storage_market_actor.go" code="true" lang="go" >}}

# `StorageMarketActor` implementation

{{< readfile file="/docs/actors/actors/builtin/storage_market/storage_market_actor.go" code="true" lang="go" >}}


{{<label storage_deal_collateral>}}
# Storage Deal Collateral

There are two types of Storage Deal Collateral, ProviderDealCollateral and ClientDealCollateral. Both types of `StorageDealCollateral` are held in the `StorageMarketActor`.
Their values are agreed upon by the storage provider and client off-chain, but must be greater than a protocol-defined minimum in any deal. Storage providers will choose to offer greater provider deal collateral to signal high-quality storage to clients. Provider deal collateral is only slashed when a sector is terminated other than normal expiration. If a miner enters Temporary Fault for a sector and later recovers from it, no deal collateral will be slashed.

Upon graceful deal expiration, storage providers must wait for finality number of epochs (as defined in {{<sref finality>}}) before being able to withdraw their `StorageDealCollateral` from the `StorageMarketActor`.
