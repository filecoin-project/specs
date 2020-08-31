---
title: "Actors"
weight: 3
dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
math-mode: true
---

# Actors
---

## System

## Init

### Exec

## Reward

## Cron

## Power

### CreateMiner

{{< mermaid >}}
sequenceDiagram
    autonumber
    Account->>Power: CreateMiner
    Power->>Init: Exec
{{</ mermaid >}}


## Market

## Account

## Miner
A *miner* is an user that participates in the consensus protocol and provides storage capacity for the market.

A `Miner` actor is an entity in the Filecoin VM that can make deal, seal sectors, mine blocks and earn block rewards. A miner can control one or more `Miner` actors.

> Why having one miner instead of many?

A `Miner` actor can only be associated with a Seal proof. If a miner wants to have 32GiB and 64GiB SDR sectors, they must have two different `Miner` actors.

### Keys
A `Miner` actor can be controlled by different key types: one *worker key*, one *owner key* and up to ten *control keys*.

The *owner key* is the key used at `Power.CreateMiner` and it cannot change, but it can be used to change the *worker key*.

The *worker key* is used for the leader election step in consensus and it can be changed via `miner.WorkerKeyChange` using the *owner key*.

All the keys can be used to interact with the `Miner` actor. We reccomend to avoid using the *owner key* and use the *control keys* for daily operation.

If the miner calls `miner.WorkerKeyChange`, the *worker key* cannot be used until the key is updated (this requires `WorkerKeyChangeDelay` epochs). This implies that the miner won't be able to participate in leader election with the previous key.*Control keys* are updated immediatedly.

### Miner Duties
- Timely follow up `miner.PreCommitSectors` with valid `miner.ProveCommitSector` (see PreCommit Faults)
- Timely submit `miner.SubmittedWindowPoSt` on healty sectors (see Faults)
- Timely submit `miner.DeclareFaults` on faulty sectors to reduce penalties (see Declared Faults)

### State Machine

{{< mermaid >}}
stateDiagram
    Null --> Precommitted: PreCommitSectors
    Precommitted --> Committed: CommitSectors
    Precommitted --> Deleted: CronPreCommitExpiry (PCD)
    Committed --> Active: SubmittedWindowPoSt
    Committed --> Faulty: DeclareFault\nSubmitWindowPoSt (SP)\nProvingDeadline (SP)
    Committed --> Terminated: TerminateSectors\n(TF)
    Faulty --> Active: SubmittedWindowPoSt (FF)
    Faulty --> Faulty: ProvingDeadline (FF)
    Faulty --> Recovering: DeclareFaultRecovered
    Faulty --> Terminated: EarlyExpiration (TF)\nTerminateSectors (TF)
    Recovering --> Active: SubmittedWindowPoSt (FF)
    Recovering --> Faulty: DeclareFault\nProvingDeadline (SP)
    Recovering --> Terminated: TerminateSectors (TF)
    Active --> Active: SubmittedWindowPoSt
    Active --> Faulty: DeclareFault\nSubmitWindowPoSt (SP)\nProvingDeadline (SP)
    Active --> Terminated: CronExpiration\nTerminateSectors (TF)
    Terminated --> Deleted: CompactSectors
{{</ mermaid >}}


| State        | Counts for power | Requires WindowPoSt |
|--------------|------------------|---------------------|
| Precommitted | -                | -                   |
| Committed    | -                | Yes                 |
| Faulty       | -                | -                   |
| Recovering   | -                | Yes                 |
| Active       | Yes              | Yes                 |
| Terminated   | -                | -                   |
| Deleted      | -                | -                   |

Note: see section on power accounting for a more detailed explanation of when a sector counts for power.

### State

### Methods

Instances of `Miner` are created only via `Power.CreateMiner`.


| Method ID | Method name |
|-----------|---------------|
|	      	`2` | `ControlAddresses` |
|	      	`3` | `ChangeWorkerAddress` |
|	      	`4` | `ChangePeerID` |
|		      `5` | `SubmitWindowedPoSt` |
|		      `6` | `PreCommitSector` |
|		      `7` | `ProveCommitSector` |
|		      `8` | `ExtendSectorExpiration` |
|		      `9` | `TerminateSectors` |
|		     `10` | `DeclareFaults` |
|		     `11` | `DeclareFaultsRecovered` |
|		     `12` | `OnDeferredCronEvent` |
|		     `13` | `CheckSectorProven` |
|		     `14` | `AddLockedFund` |
|	       `15` | `ReportConsensusFault` |
|	       `16` | `WithdrawBalance` |
|		     `17` | `ConfirmSectorProofsValid`|
|		     `18` | `ChangeMultiaddrs` |
|		     `19` | `CompactPartitions` |
|		     `20` | `CompactSectorNumbers` |

### ControlAddresses

### ChangeWorkerAddress

### ChangePeerID

### SubmitWindowedPoSt
{{< mermaid >}}
sequenceDiagram
    autonumber
    Account->>Miner: SubmitWindowedPoSt
    Miner ->>Power: UpdateClaimedPower
    Miner ->>BurntFunds: Send
    Miner ->>Power: UpdatePledgeTotal
{{</ mermaid >}}

### PreCommitSector

A sealed sector is specific to its data, to its `Miner` actor and to a point in the chain (a chain reference called `SealRandomness`).

#### Internal Messages
{{< mermaid >}}
sequenceDiagram
    autonumber
    Account->>Miner: PreCommitSector
    Miner ->> Reward: ThisEpochReward
    Miner ->> Power: CurrentTotalPower
    Miner ->> Market: VerifyDealsForActivation
{{</ mermaid >}}

### ProveCommitSector

#### Internal Messages
{{< mermaid >}}
sequenceDiagram
    Account ->>Miner: ProveCommitSector
    Miner ->>Power: SubmitPoRepForBulkVerify
    Power ->>Miner: ConfirmSectorProofsValid
    Miner ->>Market: ActivateDeals
    Miner ->>Power: UpdateClaimedPower
    Miner ->>Power: UpdatePledgeTotal
{{</ mermaid >}}

### ExtendSectorExpiration

### TerminateSectors

### DeclareFaults

### DeclareFaultsRecovered

### OnDeferredCronEvent

### CheckSectorProven

### AddLockedFund

### ReportConsensusFault

### WithdrawBalance

### ConfirmSectorProofsValid

### ChangeMultiaddrs

### CompactPartitions

### CompactSectorNumbers

# Parameters

Parameter table here
