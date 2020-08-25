---
title: Deal Flow
weight: 1
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Storage Deal Flow
---

![Deal Flow Sequence Diagram](diagrams/deal-flow.mmd)

## Add Storage Deal and Power

1. `StorageClient` and `StorageProvider` call `StorageMarketActor.AddBalance` to deposit funds into Storage Market.
   - `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made.
2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
   - `StorageProvider` verifies the `StorageDeal` by checking address and signature of `StorageClient`, checking the proposal's `StartEpoch` is after the current Epoch, checking `StorageClient` did not call withdraw in the last X Epoch (`WithdrawBalance` should take at least X Epoch), checking both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
3. `StorageProvider` signs the `StorageDealProposal`  by constructing an on-chain message.
   - `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` to publish this on-chain message which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`. However, the deals are not active at this point.
     - As a backup, `StorageClient` may call `PublishStorageDeals` with the `StorageDeal`, to activate the deal if they can obtain the signed on-chain message from `StorageProvider`.
     - It is possible for either `StorageProvider` or `StorageClient` to try to enter into two deals simultaneously with funds available only for one. Only the first deal to commit to the chain will clear, the second will fail with error `errorcode.InsufficientFunds`.
   - `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.

## Sealing sectors

4. Once the miner finishes packing a `Sector`, it generates a SectorPreCommitInfo and calls PreCommitSector with a PreCommitDeposit. It must call ProveCommitSector with SectorProveCommitInfo within some bound to recover the deposit. An expired PreCommit message will result in PreCommitDeposit being burned. There are two types of sectors, Regular Sector and Committed Capacity Sector but all sectors have an explicit expiration epoch declared during PreCommit. For a Regular Sector with storage deals in it, all deals must expire before sector expiration. Miner gains power for this particular sector upon successful ProveCommit.

## Receive Challenge

5. Miners enter the `Challenged` status when receiving a SurprisePoSt challenge from the chain. Miners will then have X Epoch as the ProvingPeriod to submit a successful PoSt before the chain checks for SurprisePoSt expiry. Miners can only get out the challenge with `SubmitSurprisePoStResponse`.
6. Miners are allowed to DeclareTemporaryFault when they are in the `Challenged` state but this will not change the list of sectors challenged as `Challenged` state specifies a list of sectors to be challenged which is a snapshot of all Active sectors at the time of challenge. Miners are also allowed to call ProveCommit which will add to their ClaimedPower but their Nominal and Consensus Power are still zero whe  they are in either Challenged or DetectedFault state.

## Declare and Recover Faults

7. Declared faults are penalized to a smaller degree than DetectedFault. Miners declare failing sectors by invoking `DeclareTemporaryFaults` with a specified fault duration and associated `TemporaryFaultFee`. Miner will lose power associated with the sector when the TemporaryFault period begins.
8. The loss of power associated with TemporaryFault will be restored when the TemporaryFault period has ended and the miner is now expected to prove over that sector. Failure to do so will result in unsuccessful ElectionPoSt or unsuccessful SurprisePoSt that leads to detected faults.


## Detect Faults

9. `CronActor` triggers `StorageMinerActor._rtCheckSurprisePoStExpiry` through `StoragePowerActor` and checks if SurprisePoSt challenge has expired for a particular miner.
   - If no PoSt is submitted by the end of the `ProvingPeriod`, miner enters `DetectedFault` state, some `PledgeCollateral` is slashed, and all power is lost.
   - Miners will now have to wait for the next SurprisePoSt challenge.
   - If the faults persist for `MAX_CONSECUTIVE_FAULTS` then sectors are terminated and provider deal collateral is slashed. 

## Sector Expiration

10. Sector expiration is done via a scheduled Cron event `_rtCheckSectorExpiry`. Sector expires when its Expiration epoch is reached and sector expiration epoch must be greater than the expiration epoch of all its deals.

## Deal Payment and slashing

11.  Deal payment and slashing are evaluated lazily through `_updatePendingDealState` at WithdrawBalance and PublishStorageDeals events. The method is also called at `OnEpochTickEnd` on StorageMarketActor as a clean up mechanism.
