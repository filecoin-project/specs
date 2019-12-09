---
menuTitle: Deals
statusIcon: 🔁
title: Market Deals
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties. All deal negotiation happen off chain and a request-response style storage deal protocol is in place to submit agreed-upon storage deals onto the network with `PublishStorageDeal` and `CommitSector` to gain storage power on chain. Hence, there is a `StorageDealProposal` and a `RetrievalDealProposal` that are half-signed contracts submitted by clients to be counter-signed and posted on-chain by the miners.

Filecoin Storage Market Deal Flow

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
- 4. Once the miner finishes packing a `Sector`, it generates a Sealed Sector and calls `StorageMinerActor.CommitSector` to verify the seal, store sector expiration, and record the mapping from `SectorNumber` to `SealCommitment`. It will also place this newly added `Sector` in the list of `CommittedSectors` in `StorageMinerActor`. `StorageMiner` does not earn any power for this newly added sector until its first PoSt has been submitted. Note that `CommitSector` can be called any time. However, sectors will be added to a staging buffer `StagedCommittedSectors` when miners are in the `Challenged` status (see 5 below).

# Receive Challenge

- 5. Miners enter the `Challenged` status whenever `NotifyOfPoStChallenge` is called by the chain. Miners will then have X Epoch as the ProvingPeriod to submit a successful PoSt before `CheckSurprisePoStSubmissionHappened` is called by the chain. Miners can only get out the challenge with `SubmitSurprisePoSt` or `onMissedSurprisePoSt`.
- 6. Miners are not allowed to call `DeclareFaults` or `RecoverFaults` when they are in the `Challneged` state but `CommitSector` is allowed and sectors will be added to a `StagedCommittedSectors` buffer. When miners get out of the `Challenged` status, `StagedCommittedSectors` will be copied over to their `Sectors`, `ProvingSet` and `SectorTable` and emptied.

# Declare and Recover Faults

- 7. Declared faults are penalized to a smaller degree than detected faults by `CronActor`. Miners declare failing sectors by invoking `StorageMinerActor.DeclareFaults` and X of the `StorageDealCollateral` will be slashed and power corresponding to these sectors will be tempororily lost. However, miners can only declare faults when they are not in `Challenged` status.
- 8. Miners can then recover faults by invoking `StorageMinerActor.RecoverFaults` and have sufficient `StorageDealCollateral` in their available balances. FaultySectors are recommitted and power is only restored at the next PoSt submission. Miners will not be able to invoke `RecoverFaults` when they are in the `Challenged` status.
- 9. Sectors that are failing for `storagemining.MAX_CONSECUTIVE_FAULTS` consecutive ChainEpochs will be cleared and result in `StoragePowerActor.SlashPledgeCollateral`.
  - TODO: set `X` parameter

# Submit PoSt

(TODO: move into Storage Mining)

On every PoSt Submission, the following steps happen.

- 10. `StorageMinerActor` first verifies the PoSt Submission. If PoSt is done correctly, all `Committed` and `Recovering` sectors will be marked as `Active` and power is credited to these sectors. Payments will be processed for deals that are `Active` by invoking `StorageMarketActor.ProcessStorageDealsPayment`.
- 11. For all sectors that are off from the `ProvingSet`, these sectors are failing. Increment `FaultCount` on these sectors and if any of these sectors are failing for `MAX_CONSECUTIVE_FAULTS` consecutive `ChainEpoch`, these sectors are terminated and cleared from the network.
- 13. Process sector expiration. Sectors expire when all deals in that sector have expired. Expired sectors will be cleared and `StorageDealCollateral` for both miners and users returned depending on the state that the sectors are in.
- 14. Submit `FaultReport` and `PowerReport` to `StoragePowerActor` for slashing and power accounting.
- 15. Check and ensure that Pledge Collateral is statisfied. TODO: some details are missing here, also related to ProvingPeriod depending on PoSt construction.
- 16. Update challenge status and add `Committed` sectors received during the challenge to the `Sectors`, `ProvingSet`, and `SectorTable`.
- 17. All Sectors will be considered in `DetectedFaults` when a miner fail to `SubmitSurprisePoSt` or `SubmitElectionPoSt` (jointly `SubmitPoSt` methods) in a proving period and detected by `onMissedSurprisePoSt` in `CheckSurprisePoStSubmissionHappened` (move to State 18).

# Detect Faults

(TODO: move into Storage Mining)

- 18. `CronActor` calls `StoragePowerActor.EpochTick` at every block. This calls `StorageMinerActor.CheckSurprisePoStSubmissionHappened` on all the miners whose `ProvingPeriod` is up.
  - If no PoSt is submitted by the end of the `ProvingPeriod`, `onMissedSurprisePoSt` detects the missing PoSt, and sets all sectors to `Failing`.
    - TODO: reword in terms of a conditional in the mining cycle
  - When there are sector faults are detected, some of `StorageDealCollateral` and `PledgeCollateral` are slashed, and power is lost.
    - If the faults persist for `storagemining.MAX_CONSECUTIVE_FAULTS` then sectors are removed/cleared from `StorageMinerActor`.

# Deal Code

{{< readfile file="deal.id" code="true" lang="go" >}}

# Deal Flow

{{< diagram src="diagrams/deal-flow.mmd.svg" title="Deal Flow Sequence Diagram" >}}

{{<label deal_states>}}
# Deal States

All on-chain economic activities in Filecoin start with the deal. This section aims to explain different states of a deal and their relationship with other concepts in the protocol such as Power, Payment, and Collaterals.

A deal has the following states:

- `Unpublished`: the deal has yet to be posted on chain.
- `Published`: the deal has been published and accepted by the chain but is not yet active as the sector containing the deal has not been proven.
- `Active`: the deal has been proven and not yet expired.
- `Deleted`: the deal has expired or the sector containing the deal has been terminated because of faults.

Note that `Unpublished` and `Deleted` states not tracked on chain. To reduce on-chain footprint, an `OnChainDeal` struct is created when a deal is published and it keeps track of a `LastPaymentEpoch` which defaults to -1 when a deal is in the `Published` state. A deal transitions into the `Active` state when `LastPaymentEpoch` is positive.

The following describes how a deal transitions between its different states.

- `Unpublished -> Published`: this is triggered by `StorageMarketActor.PublishStorageDeals` which validates new storage deals, locks necessary funds, generates deal IDs, and registers the storage deals in `StorageMarketActor`.
- `Published -> Deleted`: this is triggered by `StorageMinerActor.ProveCommitSector` during InteractivePoRep when the elapsed number of epochs between PreCommit and ProveCommit messages exceeds `MAX_PROVE_COMMIT_SECTOR_EPOCH`. ProveCommitSector will also trigger garbage collection on the list of published storage deals.
- `Published -> Active`: this is triggered by `ActivateStorageDeals` after successful `StorageMinerActor.ProveCommitSector`. It is okay for the StorageDeal to have already started (i.e. for `StartEpoch` to have passed) at this point but it must not have expired.
- `Active -> Deleted`: this can happen under the following conditions:
  - The deal itself has expired. This is triggered by `StorageMinerActorCode._submitPowerReport` which is called whenever a PoSt is submitted. Power associated with the deal will be lost, collaterals returned, and all remaining storage fees unlocked (allowing miners to call `WithdrawBalance` successfully).
  - The sector containing the deal has expired. This is triggered by `StorageMinerActorCode._submitPowerReport` which is called whenver a PoSt is submitted. Power associated with the deals in the sector will be lost, collaterals returned, and all remaining storage fees unlocked.
  - The sector containing the active deal has been terminated. This is triggered by `StorageMinerActor._submitFaultReport` for `TerminatedFaults`. No storage deal collateral will be slashed on fault declaration or detection, only on termination. A terminated fault is triggered when a sector is in the `Failing` state for `MAX_CONSECUTIVE_FAULTS` consecutive proving periods.

Given deal states and their transitions, the following are the relationships between deal states and other economic states and activities in the protocol.

- `Power`: only payload data in an Active storage deal counts towards power.
- `Deal Payment`: happens on `_onSuccessfulPoSt` and at deal/sector expiration through `_submitPowerReport`, paying out `StoragePricePerEpoch` for each epoch since the last PoSt.
- `Deal Collateral`: no storage deal collateral will be slashed for `NewDeclaredFaults` and `NewDetectedFaults` but instead some pledge collateral will be slashed given these faults' impact on consensus power. In the event of `NewTerminatedFaults`, all storage deal collateral and some pledge collateral will be slashed. Provider and client storage deal collaterals will be returned when a deal or a sector has expired. If a sector recovers from `Failing` within the `MAX_CONSECUTIVE_FAULTS` threshold, deals in that sector are still considered active. However, miners may need to top up pledge collateral when they try to `RecoverFaults` given the earlier slashing.

{{< diagram src="diagrams/deal-payment.mmd.svg" title="Deal States Sequence Diagram" >}}
