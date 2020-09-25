---
title: Deal Flow
weight: 2
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Storage Deal Flow

![Deal Flow Sequence Diagram](diagrams/deal-flow.mmd)

## Add Storage Deal and Power

1. `StorageClient` and `StorageProvider` call `StorageMarketActor.AddBalance` to deposit funds into Storage Market.
   - `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made.
2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
   - `StorageProvider` verifies the `StorageDeal` by checking:
   		- the address and signature of the `StorageClient`,
   		- the proposal's `StartEpoch` is after the current Epoch,
   		- the `StorageClient` did not call withdraw in the last **X** Epochs (`WithdrawBalance` should take at least **X** Epochs), 
   		- both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
3. `StorageProvider` signs the `StorageDealProposal`  by constructing an on-chain message.
   - `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` to publish this on-chain message which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`. However, the deals are not active at this point.
     - As a backup, `StorageClient` may call `PublishStorageDeals` with the `StorageDeal`, to activate the deal if they can obtain the signed on-chain message from `StorageProvider`.
     - It is possible for either `StorageProvider` or `StorageClient` to try to enter into two deals simultaneously with funds available only for one. Only the first deal to commit to the chain will clear, the second will fail with error `errorcode.InsufficientFunds`.
   - `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.

## Sealing sectors

4. Once the miner finishes packing a `Sector`, it generates a `SectorPreCommitInfo` and calls `PreCommitSector` with a `PreCommitDeposit`. It must call `ProveCommitSector` with `SectorProveCommitInfo` within some bound to recover the deposit. An expired `PreCommit` message will result in `PreCommitDeposit` being burned. All sectors have an explicit expiration epoch declared during `PreCommit`. For Sectors with Regular Deals, all deals must expire before sector expiration. The Miner gains power for this particular sector upon successful `ProveCommit`. For more details on the Sectors and the different types of deals that can be included in a Sector refer to the [Sector section](filecoin_mining#sector).

## Prove Storage

5. Miners have to prove that they hold unique copies of Sectors by submitting proofs according to the [Proof of SpaceTime](post) algorithm. Miners have to prove all their Sectors in regular time intervals in order for the system to guarantee that they indeed store the data they committed to store in the deal phase.
6. Miners are allowed to `DeclareTemporaryFault`. Miners are also allowed to call `ProveCommit` which will add to their ClaimedPower but their Nominal and Consensus Power are still zero for the Sectors for which they have in `DeclareTemporaryFault` state.

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
