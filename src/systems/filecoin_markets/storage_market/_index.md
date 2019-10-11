---
menuTitle: Storage Market
title: "Storage Market in Filecoin"
entries:
- storage_market_actor
- storage_market_participants
- faults
---

{{<label storage_market_subsystem>}}
Storage Market subsystem is the data entry point into the network. Storage miners only earn power from data stored in a storage deal and all deals live on the Filecoin network. Specific deal negotiation process happens off chain, clients and miners enter a storage deal after an agreement has been reached and post storage deals on the Filecoin network to earn block rewards and get paid for storing the data in the storage deal. A deal is only valid when it is posted on chain with signatures from both parties and at the time of posting, there are sufficient balances for both parties in `StorageMarketActor` to honor the deal in terms of deal price and deal collateral. 

Both `StorageClient` and `StorageProvider` need to first deposit Filecoin token into `StorageMarketActor` before participating in the storage market. `StorageClient` can then send a `StorageDealProposal` to the `StorageProvider` along with the data. A partially signed `StorageDeal` is called a `StorageDealProposal`. `StorageProvider` can then put this storage deal in their `Sector`, countersign the `StorageDealProposal` and result in a `StorageDeal`. A `StorageDeal` is only in effect when it is submitted to and accepted by the `StorageMarketActor` on chain before the `ProposalExpiryEpoch`. `StorageDeal` does not include a `StartEpoch` as it will come into effect at the block when the deal gets accepted into the network. Hence, `StorageProvider` should publish the deal as soon as possible.

`StorageDeal` payments are processed at every successful PoSt submission and `StorageMarketActor` will move locked funds from `StorageClient` to `StorageProvider`. `SlashStorageDealCollateral` is also triggered on PoSt submission when a Sector containing a particular `StorageDeal` is faulty or miners fail to submit PoSt related to a `StorageDeal`. Note that `StorageProvider` does not need to be the same entity as the `StorageMinerActor` as long as the deal is stored in at least one `Sector` throughout the life time of the storage deal.
