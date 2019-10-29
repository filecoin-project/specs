---
menuTitle: Deals
statusIcon: 🔁
title: Market Deals
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties. All deal negotiation happen off chain and a request-response style storage deal protocol is in place to submit agreed-upon storage deals onto the network with `CommitSector` to gain storage power on chain. Hence, there is a `StorageDealProposal` and a `RetrievalDealProposal` that are half-signed contracts submitted by clients to be counter-signed and posted on-chain by the miners.

Filecoin Storage Market Deal Flow

# Add Storage Deal and Power

- 1. `StorageClient` and `StorageProvider` call `StorageMarketActor.AddBalance` to deposit funds into Storage Market. There are two fund states in the Storage Market, `Locked` and `Available`.
    - `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made. (move to state X)
- 2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
    - `StorageProvider` verifies the `StorageDeal` by checking address and signature of `StorageClient`, checking the proposal has not expired, checking `StorageClient` did not call withdraw in the last X Epoch, checking both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
- 3. `StorageProvider` signs the `StorageDealProposal`  gets a `StorageDeal`.
    - a. `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`. However, the deals are not active at this point.
      - As a backup, `StorageClient` MAY call `PublishStorageDeals` with the `StorageDeal`, to activate the deal.
      - It is possible for either `StorageProvider` or `StorageClient` to try to enter into two deals simultaneously with funds available only for one. Only the first deal to commit in the chain would clear, the second would fail with error `errorcode.InsufficientFunds`.
    - b. `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.
- 4. Once the miner finishes packing a `Sector`, it generates a Sealed Sector and calls `StorageMinerActor.CommitSector` to verify the seal, store sector expiration, and record the mapping from `SectorNumber` to `SealCommitment`. It will also place this newly added `Sector` in the list of `CommittedSectors` in `StorageMinerActor`. `StorageMiner` does not earn any power for this newly added sector until its first PoSt has been submitted.

# Declare and Recover Faults

- 5. Declared faults are penalized to a smaller degree than spotted faults by `CronActor`. Miners declare faulty sectors by invoking `StorageMinerActor.DeclareFaults` and X of the `StorageDealCollateral` will be slashed and power corresponding to these sectors will be tempororily lost.
- 6. Miners can then recover faults by invoking `StorageMinerActor.RecoverFaults` and have sufficient `StorageDealCollateral` in their available balances. FaultySectors are recommitted and power is only restored at the next PoSt submission.
- 7. Sectors that are declared faulty for `storagemining.MaxFaults` consecutive ChainEpochs will result in `StoragePowerActor.SlashPledgeCollateral`.
  - TODO: set `X` parameter

# Submit PoSt

(TODO: move into Storage Mining)

On every PoSt Submission, the following steps happen.

- 8. `StorageMinerActor` first verifies the PoSt Submission. All Sectors will be considered in `SpottedFaults` if PoSt submission has failed (move to State 14).
- 9. If `CommittedSectors` are proven in `PoStSubmission.SectorSet`, Storage Miner gains power for these newly committed sectors.
- 10. If there are `DeclaredFaultySectors` , `Sector` in that set will not be challenged.
- 11. For all other sectors, payment will be processed by invoking `StorageMarketActor.ProcessStorageDealsPayment` and miner available balances will be updated.
- 12. Decide which Sectors have expired by looking at the `SectorExpirationQueue`. Sectors expire when all deals in that Sector have expired. `StorageDealCollateral` for both miners and users will only be returned when all deals in the Sector have expired. This is done by calling `StorageMarketActor.SettleExpiredDeals` and the Sector will be deleted from `StorageMinerActor.Sectors`.

# Spot Faults

(TODO: move into Storage Mining)

- 13. `CronActor` calls `StoragePowerActor.EpochTick` at every block. This calls `StorageMinerActor.CheckPoSt` on all the miners whose `ProvingPeriod` is up.
  - If no PoSt is submitted by the end of the `ProvingPeriod`, `StorageMinerActor` spots the missing PoSt, and sets all sectors to `Failing`.
    - TODO: reword in terms of a conditional in the mining cycle
  - When there are sector faults are spotted, both `StorageDealCollateral` and `PledgeCollateral` are slashed, and power is lost.
    - If the faults persist for `storagemining.MaxFaults` then sectors are removed/cleared from `StorageMinerActor`.

# Deal Code

{{< readfile file="deal.id" code="true" lang="go" >}}
