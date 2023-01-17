---
title: Deal Flow
weight: 2
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Storage Deal Flow

![Deal Flow Sequence Diagram](diagrams/deal-flow.mmd)

## Add Storage Deal and Power

1. `StorageClient` and `StorageProvider` call `StorageMarketActor.AddBalance` to deposit funds into Storage Market.
   - `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made.
2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
   - `StorageProvider` verifies the `StorageDeal` by checking: - the address and signature of the `StorageClient`, - the proposal's `StartEpoch` is after the current Epoch, - (tentative) the `StorageClient` did not call withdraw in the last _X_ epochs (`WithdrawBalance` should take at least _X_ epochs) - _X_ is currently set to 0, but the setting will be re-considered in the near future. - both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
3. `StorageProvider` signs the `StorageDealProposal` by constructing an on-chain message.
   - `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` to publish this on-chain message which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`. However, the deals are not active at this point.
     - As a backup, `StorageClient` may call `PublishStorageDeals` with the `StorageDeal`, to activate the deal if they can obtain the signed on-chain message from `StorageProvider`.
     - It is possible for either `StorageProvider` or `StorageClient` to try to enter into two deals simultaneously with funds available only for one. Only the first deal to commit to the chain will clear, the second will fail with error `errorcode.InsufficientFunds`.
   - `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.

## Sealing sectors

4. Once a miner finishes packing a `Sector`, it generates a `SectorPreCommitInfo` and calls `PreCommitSector` or `PreCommitSectorBatch` with a `PreCommitDeposit`. It must call `ProveCommitSector` or `ProveCommitAggregate` with `SectorProveCommitInfo` within some bound to recover the deposit. Initial pledge will then be required at time of `ProveCommit`. Initial Pledge is usually higher than `PreCommitDeposit`. Recovered `PreCommitDeposit` will count towards Initial Pledge and miners only need to top up additional funds at `ProveCommit`. Excess `PreCommitDeposit`, when it is greater than Initial Pledge, will be returned to the miner. An expired `PreCommit` message will result in `PreCommitDeposit` being burned. All Sectors have an explicit expiration epoch declared during `PreCommit`. For sectors with deals, all deals must expire before sector expiration. The Miner gains power for this particular sector upon successful `ProveCommit`. For more details on the Sectors and the different types of deals that can be included in a Sector refer to the [Sector section](filecoin_mining#sector).

## Prove Storage

5. Miners have to prove that they hold unique copies of Sectors by submitting proofs according to the [Proof of SpaceTime](post) algorithm. Miners have to prove all their Sectors in regular time intervals in order for the system to guarantee that they indeed store the data they committed to store in the deal phase.

## Declare and Recover Faults

6. Miners can call `DeclareFaults` to mark certain Sectors as faulty to avoid paying Sector Fault Detection Fee. Power associated with the sector will be removed at fault declaration.
7. Miners can call `DeclareFaultsRecovered` to mark previously faulty sector as recovered. Power will be restored when recovered sectors pass WindowPoSt checks successfully.
8. A sector pays a Sector Fault Fee for every proving period during which it is marked as faulty.

## Skipped Faults

9. After a WindowPoSt deadline opens, a miner can mark one of their sectors as faulty and exempted by WindowPoSt checks, hence Skipped Faults. This could avoid paying a Sector Fault Detection Fee on the whole partition.

## Detected Faults

10. If a partition misses a WindowPoSt submission deadline, all previously non-faulty sectors in the partition are detected as faulty and a Fault Detection Fee is charged.

## Sector Expiration

11. Sector expires when its expiration epoch is reached and sector expiration epoch must be greater than the expiration epoch of all its deals.

## Sector Termination

12. Termination of a sector can be triggered in two ways. One when sector remains faulty for 42 consecutive days and the other when a miner initiates a termination by calling `TerminateSectors`. In both cases, a `TerminationFee` is penalized, which is in principle equivalent to how much the sector has earned so far. Miners are also penalized for the `DealCollateral` that the sector contains and remaining `DealPayment` will be returned to clients.

## Deal Payment and slashing

13. Deal payment and slashing are evaluated lazily through `updatePendingDealState` called at `CronTick`.
