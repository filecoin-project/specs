---
title: "Protocols"
weight: 2
dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
math-mode: true
---

# Protocol
---

## Key Concepts

### Sectors
A *sector* is a storage container in which miners store deals from the market.

A sector can be in one of the following states.

| State          | Description                                                                                                                                           |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Precommitted` | miner seals sector and submits `miner.PreCommitSector`                                                                                                |
| `Committed`    | miner generates a Seal proof and submits `miner.ProveCommitSector`                                                                                    |
| `Active`       | miner generate valid PoSt proofs and timely submits `miner.SubmitWindowedPoSt`                                                                        |
| `Faulty`       | miner fails to generate a proof (see Fault section)                                                                                                   |
| `Recovering`   | miner declared a faulty sector via `miner.DeclareFaultRecovered`                                                                                      |
| `Terminated`   | either sector is expired, or early terminated by a miner via `miner.TerminateSectors`, or was failed to be proven for 14 consecutive proving periods. |

See `Miner` actor for a more detailed explanation.

### Power
A miner wins blocks proportionally to their *Quality Adjusted Power*.

*Sector Quality Adjusted Power* is a weighted average of the quality of its space and it is based on the size, duration and quality of its deals.

Formula for calculating Sector Quality Adjusted Power (or QAp, often referred to as power):
- Calculate the `dealSpaceTime`: sum of the `duration*size` of each deal
- Calculate the `verifiedSpaceTime`: sum of the `duration*size` of each verified deal
- Calculate the spacetime without deals `baseSpaceTime`: `sectorSize*sectorDuration - dealSpaceTime - verifiedSpaceTime`
- Calculate the weighted average quality `avgQuality`: `(baseSpaceTime * QualityBaseMultiplier +  dealSpaceTime * DealWeightMultiplier + verifiedSpaceTime * VerifiedDealWeightMultiplier)/ (sectorSize * sectorDuration * QualityBaseMultiplier)`
- Calculate `sectorQuality`: `avgQuality * size`

During `miner.PreCommitSector`, the sector quality is calculated and stored in the sector information.

### Deals

#### Ask
#### Bids
#### Verified Deals

## Mining cycle

### Register a miner
A user registers a new `Miner` actor via `Power.CreateMiner`.

### Add storage capacity
A miner adds storage capacity to the network by adding sectors to a `Miner` actor.

- The miner collects storage deals and publish them via `Market.PublishStorageDeals`.
- The miner groups deals in a sector, runs the sealing process and submits the sealed sector information via `miner.PreCommitSector`. The sector state is now *precommitted*.
- The miner waits `PreCommitChallengeDelay` since the on-chain precommit epoch, gets the `InteractiveRandomness`, generates a Seal Proof and submits it via `miner.ProveCommitSector`. The sector state is now *committed*.

### Mantain storage capacity

A miner mantains sectors *active* by timely submitting Proofs-of-Spacetime (PoSt).
A PoSt proves sectors are persistently stored through time.

A miner must timely submit `miner.SubmitWindowedPoSt` for their sectors according to their assigned deadlines (see deadline section).

#### Proving Period
A *proving period* is a period of `WPoStProvingPeriod` epochs in which a `Miner` actor is scheduled to prove its storage.
A *proving period* is evenly divided in `WPoStPeriodDeadlines` *deadlines*.

Each miner has a different start of proving period `ProvingPeriodStart` that is assigned at `Power.CreateMiner`.

#### Deadline
A *deadline* is a period of `WPoStChallengeWindow` epochs that divides a proving period.
Sectors are assigned to a deadline on `miner.ProveCommitSector` and will remain assigned to it throughout their lifetime.
Sectors are also assigned to partitions.

There are four relevant epochs associated to a deadline.
- Open Epoch: Epoch from which a PoSt Proof for this deadline can be submitted.
- Close Epoch: Epoch after which a PoSt Proof for this deadline will be rejected.
- Fault Cutoff Epoch: Epoch after which a `DeclareFault` for sectors in the upcoming deadline are rejected.
- Challenge Epoch: Epoch at which the randomness for the challenges is available.

A miner must submit a `miner.SubmitWindowedPoSt` for each deadline.

#### Partitions
A partition is a set of sectors that is not larger than the Seal Proof `sp.WindowPoStPartitionSectors`.

Partitions are a by-product of our current proof mechanism. There is a limit of sectors (`sp.WindowPoStPartitionSectors`) that can be proven in a single SNARK proof. If more than this amount is required to be proven, more than one SNARK proof is required. Each SNARK proof is a partition.

Sectors are assigned to a partition at `miner.ProveCommitSector` and they can be re-arranged via `CompactPartitions`.

#### Deadline assignment

## Market cycle

## Faults and Penalties

### Collaterals
#### PreCommitDeposit

`PreCommitDeposit` is the collateral deposited on `miner.PreCommitSector` and it is returned on `miner.ProveCommit`. If the `miner.ProveCommit` is not executed before `MaxProveCommitDuration` epochs since precommit.

The `PreCommitDeposit` is locked separately from Initial Pledge and cannot be used as collateral for other faults.

#### InitialPledge
#### LockedFunds

### Fees
#### Fault Fee (FF)
Fault Fee (or "FF") = 2.15 BRapprox

#### Sector Penalization (SP)
Sector Penalization (or "SP") = 5 BRapprox

#### Termination Fee (TF)
Termination Fee (or "TF") = TODO

#### PreCommit Deposit Fee (PCD)
Sector Penalization (or "SP") = 20 BRapprox

#### Consensus Fault Fee (CFF)
Consensus Fault Fee (or "CFF") = 5 br

### Faults

A sector is considered faulty if a miner is not able to generate a proof.
There could be different reasons for which a sector can be faulty: power outage, hardware failure, internet connection loss.

If a sector state is not *active*, it will not count towards power.

#### PreCommit Faults
A sector fault is a "precommit fault" if it is *precommitted* after `MaxProveCommitDuration` epochs since precommit.

**Penalization**:
- Fee: PreCommitDeposit fee is charged on its first deadline. The fee is charged from the PreCommitDeposit collateral.
- Sector state: Sector is immediatedly marked as *deleted*.

#### Declared Faults
A sector fault is a "declared fault" if it is *faulty* on the deadline's fault cutoff epoch.

There are three cases of Declared Faults:
- **New Declared Fault**: A sector is *active*, and its Miner actor processes `miner.DeclareFault` before the deadline's fault cutoff epoch.
- **Declared Failed Recovery Fault**:  A sector is *recovering* and its Miner actor processes `miner.DeclareFault` before the deadline's fault cutoff epoch.
- **Continued Fault**: A sector was faulty in the previous proving period and continues beind so.

**Penalization:** 
- Fee: Fault fee is charged on its next deadline. The fee is charged from Initial Pledge.
- Sector state: Sector is immediatedly marked as *faulty*.

#### Skipped Faults
A sector fault is a "skipped fault" if the sector is not *faulty* before deadline's fault cutoff epoch, the miner submits a `miner.SubmittedWindowPoSt` marking the sector as skipped.

There two cases of Skipped Faults:
- **Active Skipped Fault**: A sector is *active* and the miner marked it as skipped at `miner.SubmittedWindowPoSt`.
- **Recovered Skipped Fault**: A sector is *recovering* and the miner marked it as skipped at `miner.SubmittedWindowPoSt`.

**Penalization**:
- Fee: Sector Penalization is immediatedly charged (on `miner.SubmittedWindowPoSt`). The fee is charged from Initial Pledge.
- Sector state: Sector is immediatedly marked as *faulty* (on `miner.SubmittedWindowPoSt`).

#### Detected Faults
A sector fault is a "detected fault" if the sector is not already *faulty* before deadline's fault cutoff epoch and no `miner.SubmittedWindowPoSt` is executed before the deadline's close epoch.

There two three of Skipped Faults:
- **Active Skipped Fault**: A sector is *active* and the miner missed `miner.SubmittedWindowPoSt`.
- **Recovered Skipped Fault**: A sector is *recovering* and the miner marked it as skipped at `miner.SubmittedWindowPoSt**.

**Penalization**:
- Fee: Sector Penalization is immediatedly charged (on `ProvingDeadline`). The fee is charged from Initial Pledge.
- Sector state: Sector is immediatedly marked as *faulty* (on `ProvingDeadline`).


## Rewards

### Block Rewards

### Storage Rewards
