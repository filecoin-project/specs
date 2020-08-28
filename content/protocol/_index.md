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

## Mining Cycle

### Register a miner
A user registers a new `Miner` actor via `Power.CreateMiner`.

The user can now add storage capacity to the miner and engage in the market.

### Add storage capacity
A miner adds storage capacity to the network by adding sectors to a `Miner` actor.

1. The miner collects storage deals and publish them via `Market.PublishStorageDeals`.
2. The miner groups deals in a sector, runs the sealing process and submits the sealed sector information via `miner.PreCommitSector`. The sector state is now *precommitted*.
3. The miner waits `PreCommitChallengeDelay` since the on-chain precommit epoch, gets the `InteractiveRandomness`, generates a Seal Proof and submits it via `miner.ProveCommitSector`. The sector state is now *committed*.

See section on collateral requirements to understand what deposits are required.

### Mantain storage capacity

A miner mantains sectors *active* by timely submitting Proofs-of-Spacetime (PoSt).
A PoSt proves sectors are persistently stored through time.

A miner must timely submit `miner.SubmitWindowedPoSt` for their sectors according to their assigned deadlines (see deadline section).

#### Proving Period
A *proving period* is a period of `WPoStProvingPeriod` epochs in which a `Miner` actor is scheduled to prove its storage.
A *proving period* is evenly divided in `WPoStPeriodDeadlines` *deadlines*.

Each miner has a different start of proving period `ProvingPeriodStart` that is assigned at `Power.CreateMiner`.

#### Deadlines
A *deadline* is a period of `WPoStChallengeWindow` epochs that divides a proving period.
Sectors are assigned to a deadline on `miner.ProveCommitSector` and will remain assigned to it throughout their lifetime.
Sectors are also assigned to a partition (see Partitions section).

A miner must submit a `miner.SubmitWindowedPoSt` for each deadline.

There are four relevant epochs associated to a deadline:

| Name          | Distance from `Open`      | Description                                                                                                                   |
|---------------|---------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| `Open`        | `0`                       | Epoch from which a PoSt Proof for this deadline can be submitted.                                                             |
| `Close`       | `WPoStChallengeWindow`    | Epoch after which a PoSt Proof for this deadline will be rejected.                                                            |
| `FaultCutoff` | `-FaultDeclarationCutoff` | Epoch after which a `miner.DeclareFault` and `miner.DeclareFaultRecovered` for sectors in the upcoming deadline are rejected. |
| `Challenge`   | `-WPoStChallengeLookback` | Epoch at which the randomness for the challenges is available.                                                                |


#### Partitions
A partition is a set of sectors that is not larger than the Seal Proof allowed number of sectors `sp.WindowPoStPartitionSectors`.

Partitions are a by-product of our current proof mechanism. There is a limit of sectors (`sp.WindowPoStPartitionSectors`) that can be proven in a single SNARK proof. If more than this amount is required to be proven, more than one SNARK proof is required. Each SNARK proof is a partition.

Sectors are assigned to a partition at `miner.ProveCommitSector` and they can be re-arranged via `CompactPartitions`.

#### Declaring and recovering faults

### Mining Blocks
Filecoin miners attempt to mine block at every epoch by running the leader election process (see Consensus).

Miners win proportionally to their *active* power as reported in the power table,n only create blocks if they are *eligible*.

See the section on the block mining process.

#### Power accounting
The active QA power of a `Miner` actor is stored in `Power.Claims` and it is the sum of the QA power of all the active sectors of a miner.

The active QA power of the network is stored in `Power.TotalQualityAdjPower` and does not include power of miner below minimum miner size.

The QA power values used in the election are taken from `WinningPoStSectorSetLookback` epochs back.

#### Minimum miner size
A `Miner` actor must have a minimum miner size of `ConsensusMinerMinPower` in order for their power to count in the power table.

#### Eligibility requirement
A `Miner` actor is eligible to mine blocks if `miner.MinerEligibleForElection()` returns true.

The eligibility requirements are:
1. Initial Pledge requirement is met at the election epoch.
2. No active consensus faults at the election epoch.
3. Initial Pledge is sufficient to cover for consensus faults at the election epoch.
4. Miner had minimum miner size at the `WinningPoStSectorSetLookback`

## Market cycle

## Faults and Penalties

### Collaterals

| Name             | Type       | Deposited at              | Penalizations                                    | Withdrawable after        |
|------------------|------------|---------------------------|--------------------------------------------------|---------------------------|
| PreCommitDeposit | per sector | `miner.PreCommitDeposit`  | PreCommit Faults                                 | `miner.ProveCommitSector` |
| InitialPledge    | per sector | `miner.ProveCommitSector` | Declared Faults, Skipped Faults, Detected Faults | On sector expiration      |
|                  |            |                           |                                                  |                           |

#### PreCommitDeposit

`PreCommitDeposit` is the collateral that secures that a *precommitted* sectors becomes *committed* before `MaxProveCommitDuration` since the precommit on-chain inclusion.
It is deposited during `miner.PreCommitSector`, it is released on `miner.ProveCommitSector` and it is penalized on PreCommit faults.

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

## Security

### Attacks

#### PreCommit Attack

Aim: Store a compressible sector.

Context:
- Attacker has a `Miner` actor.

Execution:
- Attacker crafts a sealed sector that correctly encodes less than `$(1 - spacegap)$` of its size and incorrectly encodes the remaining space with zeroes.
- Attacker successfully submits `miner.PreCommitSector`.
- Attacker waits `PreCommitChallengeDelay` for `InteractiveRandomness`
- Attacker submits `miner.ProveCommitSector` if the proof succeeds, otherwise, abort.

Outcome: Attack will succeed with probability `$<2^-{\lambda}$`

Rationality mitigation:
- Require a deposit on `PreCommitDeposit` 
- Require the `PreCommitDeposit` fee to be equivalent to block rewards earned by the QApower of the sector in 5 years.

#### Fork-and-PreCommit Attack

Aim: Store a compressible sector.

Context:
- Attacker has a `Miner` actor.

Execution:
- Attacker sees the `InteractiveRandomness` of current block.
- Attacker crafts a compressible sealed sector that correctly encodes less than `$(1 - spacegap)$` of its size and incorrectly encodes the remaning space with zeroes.
- Attacker tries to do a fork longer of size `PreCommitChallengeDelay`.
- If successful, miner includes `miner.PreCommitSector` for the compressible sector.

Outcome: Attack will succeed with the probability of a successful fork of size `PreCommitChallengeDelay`.

Mitigation:
- Require `PreCommitChallengeDelay` to be large such that the probability of a successful fork is small for the largest miner.

#### PreCommit-and-Fork Attack

Note: This attack is a weaker version of Fork-and-PreCommit Attack

Aim: Store a compressible sector.

Context:
- Attacker has a `Miner` actor.

Execution:
- Attacker crafts a compressible sealed sector that correctly encodes less than `$(1 - spacegap)$` of its size and incorrectly encodes the remaning space with zeroes.
- Attacker waits `PreCommitChallengeDelay` for `InteractiveRandomness`
- If the proof succeeds, attacker submits `miner.ProveCommitSector`
- Otherwise, attacker forks the chain for `PreCommitChallengeDelay`

Mitigation:
- Require `PreCommitChallengeDelay` to be large such that the probability of a successful fork is small for the largest miner.
