---
title: Deal States
weight: 2
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Storage Deal States

All on-chain economic activities in Filecoin start with the deal. This section aims to explain different states of a deal and their relationship with other concepts in the protocol such as Power, Payment, and Collaterals.

A deal has the following states:

- `Unpublished`: the deal has yet to be posted on chain.
- `Published`: the deal has been published and accepted by the chain but is not yet active as the sector containing the deal has not been proven.
- `Active`: the deal has been proven and not yet expired.
- `Deleted`: the deal has expired or the sector containing the deal has been terminated because of faults.

Note that `Unpublished` and `Deleted` states are not tracked on chain. To reduce on-chain footprint, an `OnChainDeal` struct is created when a deal is published and it keeps track of a `LastPaymentEpoch` which defaults to -1 when a deal is in the `Published` state. A deal transitions into the `Active` state when `LastPaymentEpoch` is positive.

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

![Deal States Sequence Diagram](diagrams/deal-payment.mmd)