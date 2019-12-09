---
menuTitle: Deal Flow
statusIcon: 🔁
title: Storage Deal Flow
---

{{< diagram src="diagrams/deal-flow.mmd.svg" title="Deal Flow Sequence Diagram" >}}

# Add Storage Deal and Power

- 1. `StorageClient` and `StorageProvider` call `StorageMarketActor.AddBalance` to deposit funds into Storage Market. There are two fund states in the Storage Market, `Locked` and `Available`.
    - `StorageClient` and `StorageProvider` can call `WithdrawBalance` before any deal is made. (move to state X)
- 2. `StorageClient` and `StorageProvider` negotiate a deal off chain. `StorageClient` sends a `StorageDealProposal` to a `StorageProvider`.
    - `StorageProvider` verifies the `StorageDeal` by checking address and signature of `StorageClient`, checking the proposal's `StartEpoch` is after the current Epoch, checking `StorageClient` did not call withdraw in the last X Epoch (`WithdrawBalance` should take at least X Epoch), checking both `StorageProvider` and `StorageClient` have sufficient available balances in `StorageMarketActor`.
- 3. `StorageProvider` signs the `StorageDealProposal`  by constructing an on-chain message.
    - `StorageProvider` calls `PublishStorageDeals` in `StorageMarketActor` to publish this on-chain message which will generate a `DealID` for each `StorageDeal` and store a mapping from `DealID` to `StorageDeal`. However, the deals are not active at this point.
      - As a backup, `StorageClient` MAY call `PublishStorageDeals` with the `StorageDeal`, to activate the deal if they can obtain the signed on-chain message from `StorageProvider`.
      - It is possible for either `StorageProvider` or `StorageClient` to try to enter into two deals simultaneously with funds available only for one. Only the first deal to commit to the chain will clear, the second will fail with error `errorcode.InsufficientFunds`.
    - `StorageProvider` calls `HandleStorageDeal` in `StorageMiningSubsystem` which will then add the `StorageDeal` into a `Sector`.

# Sealing sectors

- 4. Once the miner finishes packing a `Sector`, it generates a Sealed Sector and calls `StorageMinerActor.CommitSector` to verify the seal, store sector expiration, and record the mapping from `SectorNumber` to `SealCommitment`. It will also place this newly added `Sector` in the list of `CommittedSectors` in `StorageMinerActor`. `StorageMiner` does not earn any power for this newly added sector until its first PoSt has been submitted. Note that `CommitSector` can be called any time. However, sectors will be added to a staging buffer `StagedCommittedSectors` when miners are in the `Challenged` status (see 5 below).

# Receive Challenge

- 5. Miners enter the `Challenged` status whenever `NotifyOfPoStChallenge` is called by the chain. Miners will then have X Epoch as the ProvingPeriod to submit a successful PoSt before `CheckSurprisePoStSubmissionHappened` is called by the chain. Miners can only get out the challenge with `SubmitSurprisePoSt` or `onMissedSurprisePoSt`.
- 6. Miners are not allowed to call `DeclareFaults` or `RecoverFaults` when they are in the `Challneged` state but `CommitSector` is allowed and sectors will be added to a `StagedCommittedSectors` buffer. When miners get out of the `Challenged` status, `StagedCommittedSectors` will be copied over to their `Sectors`, `ProvingSet` and `SectorTable` and emptied.

# Declare and Recover Faults

- 7. Declared faults are penalized to a smaller degree than detected faults by `CronActor`. Miners declare failing sectors by invoking `StorageMinerActor.DeclareFaults` and X of the `StorageDealCollateral` will be slashed and power corresponding to these sectors will be tempororily lost. However, miners can only declare faults when they are not in `Challenged` status.
- 8. Miners can then recover faults by invoking `StorageMinerActor.RecoverFaults` and have sufficient `StorageDealCollateral` in their available balances. FaultySectors are recommitted and power is only restored at the next PoSt submission. Miners will not be able to invoke `RecoverFaults` when they are in the `Challenged` status.
- 9. Sectors that are failing for `storagemining.MaxFaults` consecutive ChainEpochs will be cleared and result in `StoragePowerActor.SlashPledgeCollateral`.
  - TODO: set `X` parameter

# Submit PoSt

On every PoSt Submission, the following steps happen.

- 10. `StorageMinerActor` first verifies the PoSt Submission. If PoSt is done correctly, all `Committed` and `Recovering` sectors will be marked as `Active` and power is credited to these sectors. Payments will be processed for deals that are `Active` by invoking `StorageMarketActor.ProcessStorageDealsPayment`.
- 11. For all sectors that are off from the `ProvingSet`, these sectors are failing. Increment `FaultCount` on these sectors and if any of these sectors are failing for `MaxFaultCount` consecutive `ChainEpoch`, these sectors are terminated and cleared from the network.
- 13. Process sector expiration. Sectors expire when all deals in that sector have expired. Expired sectors will be cleared and `StorageDealCollateral` for both miners and users returned depending on the state that the sectors are in.
- 14. Submit `FaultReport` and `PowerReport` to `StoragePowerActor` for slashing and power accounting.
- 15. Check and ensure that Pledge Collateral is statisfied. TODO: some details are missing here, also related to ProvingPeriod depending on PoSt construction.
- 16. Update challenge status and add `Committed` sectors received during the challenge to the `Sectors`, `ProvingSet`, and `SectorTable`.
- 17. All Sectors will be considered in `DetectedFaults` when a miner fail to `SubmitSurprisePoSt` or `SubmitElectionPoSt` (jointly `SubmitPoSt` methods) in a proving period and detected by `onMissedSurprisePoSt` in `CheckSurprisePoStSubmissionHappened` (move to State 18).

# Detect Faults

- 18. `CronActor` calls `StoragePowerActor.EpochTick` at every block. This calls `StorageMinerActor.CheckSurprisePoStSubmissionHappened` on all the miners whose `ProvingPeriod` is up.
  - If no PoSt is submitted by the end of the `ProvingPeriod`, `onMissedSurprisePoSt` detects the missing PoSt, and sets all sectors to `Failing`.
    - TODO: reword in terms of a conditional in the mining cycle
  - When there are sector faults are detected, some of `StorageDealCollateral` and `PledgeCollateral` are slashed, and power is lost.
    - If the faults persist for `storagemining.MaxFaultCount` then sectors are removed/cleared from `StorageMinerActor`.
