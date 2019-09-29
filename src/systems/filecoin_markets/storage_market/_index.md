---
menuTitle: Storage Market
title: "Storage Market in Filecoin"
entries:
- storage_market_actor
- storage_participant
---

{{<label storage_market_subsystem>}}
Storage Market subsystem is the data entry point into the network. Storage miners only earn power from data stored in a storage deal and all deals live on the Filecoin network. Specific deal negotiation process happens off chain, clients and miners enter a storage deal after an agreement has been reached and post storage deals on the Filecoin network to earn block rewards and get paid for storing the data in the storage deal. A deal is only valid when it is posted on chain with signatures from both parties and at the time of posting, there are sufficient balances for both parties in `StorageMarketActor` to honor the deal in terms of deal price and deal collateral. Both `StorageClient` and `StorageProvider` can submit deals on chain once they have signatures from the counter party. A partially signed `StorageDeal` is called a `StorageDealProposal`. 

A `StorageDealProposal` can be propsed by any `StorageMarketParticipant` and sent to another participant who can counter sign the proposal with their public key and result in a `StorageDeal`. A `StorageDeal` is only in effect when it is submitted to and accepted by the `StorageMarketActor` on chain before the `ProposalExpiryEpoch`. Counter signer of the `StorageDeal` can either submit the deal directly onto the chain or return the deal to the proposer. In the event that a `StorageDeal` is reached but the counter signer does not submit the sign deal on chain, proposer can exit the storage deal by reducing its balance on `StorageMarketActor` before the deal is published so that the deal will be invalid. Hence, counter signer of the storage deal should publish the deal as soon as possible. Note that there are some UX challenges here that can be simplified with updatable storage on the network, where storage deals within a sector can be updated cheaply.

Storage Deal payments are processed at every successful PoSt submission and `StorageMarketActor` will move locked funds from `StorageClient` to `StorageProvider`. `SlashStorageDealCollateral` is also triggered on PoSt submission when a Sector containing a particular `StorageDeal` is faulty or miners fail to submit PoSt related to a `StorageDeal`. Note that `StorageProvider` does not need to be the same entity as the `StorageMinerActor` as long as the deal is stored on the Filecoin in at least one sector throughout the life time of the storage deal.

{{< readfile file="storage_market_subsystem.id" code="true" lang="go" >}}

