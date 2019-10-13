---
menuTitle: Deals
title: Market Deals
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties. All deal negotiation happen off chain and a request-response style storage deal protocol is in place to submit agreed-upon storage deals onto the network with `CommitSector` to gain storage power on chain. Hence, there is a `StorageDealProposal` and a `RetrievalDealProposal` that are half sign contracts submitted by clients to be counter-signed and posted on-chain by the miners. 

Filecoin Storage Market Deal Flow

1.`StorageClient` and `StorageProvider` deposit funds to `StorageMarketActor`
a. `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made. (move to state X)
2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
a. `StorageProvider` verifies the `StorageDeal` by checking address and signature of `StorageClient`, checking the proposal has not expired, checking `StorageClient` did not call withdraw in the last X Epoch, checking both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
3. `StorageProvider` signs the `StorageDealProposal` and gets a `StorageDeal`.
a. `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.
b. `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`.
4. `StorageMiningSubsystem` calls `CommitSector`
5. Payment, Expiration, and Faults

{{< readfile file="deal.id" code="true" lang="go" >}}
