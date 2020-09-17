---
title: Storage Market Actor
weight: 1
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# Storage Market Actor

`StorageMarketActor` is responsible for processing and managing on-chain deals. This is also the entry point of all storage deals and data into the system. It maintains a mapping of `StorageDealID` to `StorageDeal` and keeps track of locked balances of `StorageClient` and `StorageProvider`. When a deal is posted on chain through the `StorageMarketActor`, it will first check if both transacting parties have sufficient balances locked up and include the deal on chain. 

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/market/market_state.go" lang="go" symbol="State" title="Storage Market Actor State">}}

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/market/market_actor.go" lang="go" title="Storage Market Actor" >}}

The Storage Market Actor Balance states and mutations can be found [here](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/market/market_balances.go).

## Storage Deal Collateral

Apart from [Initial Pledge Collateral and Block Reward Collateral](miner_collaterals) discussed earlier, the third form of collateral is provided by the storage provider to _collateralize deals_, is called _Storage Deal Collateral_ and is held in the `StorageMarketActor`.

There is a minimum amount of collateral required by the protocol to provide a minimum level of guarantee, which is agreed upon by the storage provider and client off-chain. However, miners can offer a higher deal collateral to imply a higher level of service and reliability to potential clients. Given the increased stakes, clients may associate additional provider deal collateral beyond the minimum with an increased likelihood that their data will be reliably stored.

Provider deal collateral is only slashed when a sector is terminated before the deal expires. If a miner enters Temporary Fault for a sector and later recovers from it, no deal collateral will be slashed.

This collateral is returned to the storage provider when all deals in the sector successfully conclude. Upon graceful deal expiration, storage providers must wait for finality number of epochs (as defined in [Finality](expected_consensus#finality-in-ec)) before being able to withdraw their `StorageDealCollateral` from the `StorageMarketActor`.


{{<katex>}}
MinimumProviderDealCollateral = \\[0.2cm] \ \ \ \ \ \ \ \ 5\% \times FILCirculatingSupply \times \frac{DealRawByte}{max(NetworkBaseline, NetworkRawBytePower)}
{{</katex>}}
