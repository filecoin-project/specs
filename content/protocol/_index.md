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

A *sector* with no deals is called *Committed Capacity*.

#### Sector states
A sector can be in one of the following states.

| State          | Description                                                                                                                                           |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Precommitted` | Miner seals sector and submits `miner.PreCommitSector`                                                                                                |
| `Committed`    | Miner generates a Seal proof and submits `miner.ProveCommitSector`                                                                                    |
| `Active`       | Miner generate valid PoSt proofs and timely submits `miner.SubmitWindowedPoSt`                                                                        |
| `Faulty`       | Miner fails to generate a proof (see Fault section)                                                                                                   |
| `Recovering`   | Miner declared a faulty sector via `miner.DeclareFaultRecovered`                                                                                      |
| `Terminated`   | Either sector is expired, or early terminated by a miner via `miner.TerminateSectors`, or was failed to be proven for 14 consecutive proving periods. |

See `Miner` actor for a more detailed explanation.

#### Sector content

*Unsealed sector* refers to the data that is stored in a sector

*Sealed sector* refers to the output of the sealing process.

#### Sector quality

*Sector Quality Adjusted Power* is a weighted average of the quality of its space and it is based on the size, duration and quality of its deals.

| Name                         | Description                                           |
|------------------------------|-------------------------------------------------------|
| BaseMultiplier               | Multiplier for power for storage without deals.       |
| DealWeightMultiplier         | Multiplier for power for storage with deals.          |
| VerifiedDealWeightMultiplier | Multiplier for power for storage with verified deals. |


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

The user can now add storage capacity to their `Miner` actor and engage in the market.

### Add storage capacity
A miner adds storage capacity to the network by adding sectors to a `Miner` actor.

1. The miner collects storage deals and publish them via `Market.PublishStorageDeals`.
2. The miner groups deals in a sector, runs the sealing process and submits the sealed sector information via `miner.PreCommitSector` and deposits `PreCommitDeposit`. The sector state is now *precommitted*.
3. The miner waits `PreCommitChallengeDelay` since the on-chain precommit epoch, gets the `InteractiveRandomness`, generates a Seal Proof, submits it via `miner.ProveCommitSector` and deposits `InitialPledge`. The sector state is now *committed*.

See section on collateral requirements to understand what deposits are required.

### Mantain storage capacity
A miner mantains its sectors *active* by generating Proofs-of-Spacetime (PoSt) and timely submit `miner.SubmitWindowedPoSt` for their sectors.
A PoSt proves sectors are persistently stored through time.

Each miner proves all of its sector once per over a period of time called *proving period*; each sector must be proven by a particular time called deadline.

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

### Declare and recover faults

#### Reporting faulty sectors
A sector is marked as *faulty* when the miner is not able to submit a `miner.SubmittedWindowPoSt**
A sector can become faulty for different reasons, for example: power outage, hardware failure, internet connection loss.

Faults can either be declared by the user or detected by the network

##### User-declare faults
A miner can declare a sector as faulty ahead of time via `miner.DeclareFaults` (before the sector's deadline fault cutoff epoch) or at proof submission via `miner.SubmitWindowedPoSt`. If the miner does not declare a sector as faulty, they will not be able to generate a valid proof and submit a `miner.SubmitWindowedPoSt`.

##### Network-detected faults
At deadline close, the `Miner` actor will mark the sectors that have not been proven as faulty.

##### Penalizations
Power is removed from the `Power` actor when a sector is marked as *faulty*.

Fees are charged 
When a sector is marked as faulty its power is removed, a fee is charged (see the faults section for more details).

Counter for termination: If the sector is faulty for 14 proving periods, the sector is terminated and a termination fee is charged.

#### Recover faulty sectors
A faulty sector is marked as *recovering* via `miner.DeclareFaultRecovered` and it is marked as *active* at its next `miner.SubmittedWindowPoSt`.

Failure to submit a proof will result in a fault (Skipped fault if marked as skipped, Detected fault if no `miner.SubmittedWindowPoSt` appeared on chain).


### Sector Management

#### Recover Faults

#### Extend a sector

#### Expiration

#### Upgrading a sector quality

A Committed Capacity (CC) sector can be upgraded if non-faulty.

Protocol for sector upgrade
1. The miner calls `miner.PreCommitSector` with the upgrade flag and specifying the sector to be replaced. The new sector will have a different sector number.
2. The miner calls `miner.ProveCommitSector` which will schedule an early termination for the CC sector at its next deadline.
3. The CC sector will remain *active* until its next deadline, after which it will be marked as *terminated* (it will count for power and is subject to WindowPoSt submission)

Failures during upgrade:
- If the CC sector is faulty before `miner.PreCommitSector` of the new sector, the PreCommit will fail.
- If the CC sector is faulty before `miner.ProveCommitSector` of the new sector, the upgrade is aborted (old sector is not early terminated), but the new sector will be added to `Miner` actor.
- If the CC sector is faulty before its last deadline it is follow the ordinary protocol for the type of fault occurred.

#### Terminate sectors early

### Mining Blocks
Filecoin miners attempt to mine block at every epoch by running the leader election process (see Consensus).

Miners win blocks with a probability proportionally to their *active Quality Adjusted power* as reported in the power table `WinningPoStSectorSetLookback` epochs before the election.

Miners can mine blocks if they meet the mining requirements.

See the section on the block mining process.

#### Power accounting

The active QA power of a `Miner` actor is stored in `Power.Claims` and it is the sum of the QA power of all the active sectors of a miner.

The active QA power of the network is stored in `Power.TotalQualityAdjPower` and does not include power of miner below minimum miner size.

The QA power values used in the election are taken from `WinningPoStSectorSetLookback` epochs back.

#### Mining requirements 

| Require            | Valid at                                        |
|--------------------|-------------------------------------------------|
| Miner Validity     | Election epoch                                  |
| Minimum Miner Size | Election epoch - `WinningPoStSectorSetLookback` |
| Miner Eligibility  | Election epoch                                  |

##### Miner validity requirement
`Miner` actor exists  and has non zero power at the election epoch

##### Minimum miner size requirement
A `Miner` actor must have had a minimum miner size of `ConsensusMinerMinPower` exactly `WinningPoStSectorSetLookback` epochs before the election. 

##### Miner Eligibility requirement
A `Miner` actor is eligible to mine blocks if `miner.MinerEligibleForElection()` returns true for the election epoch.

The eligibility requirements are:
1. `Miner` actor is not in debt.
2. `Miner` actor has not active consensus faults.
3. `Miner`'s initial pledge is sufficient to cover for consensus fault (TODO).

### Earning rewards
Miner earn two types of rewards:
1. Mining rewards (block rewards and gas fees)
2. Storage rewards (deal payments)

#### Block Rewards
When a `Miner` actor creates a block, they earn a block reward via ``

## Market cycle

## Faults and Penalties

### Balances

#### Miner balance

A `Miner` actor has a balance

#### PreCommitDeposit

#### UnlockedFunds

#### LockedFunds


| Name             | Type       | Deposited at              | Penalizations                                                      | Withdrawable after        |
|------------------|------------|---------------------------|--------------------------------------------------------------------|---------------------------|
| PreCommitDeposit | per sector | `miner.PreCommitDeposit`  | PreCommit Faults                                                   | `miner.ProveCommitSector` |
| InitialPledge    | per sector | `miner.ProveCommitSector` | Declared Faults, Skipped Faults, Detected Faults, Consensus Faults | on sector termination.    |

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

### Storage Faults

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

#### Consensus Faults

## Rewards

### Block Rewards
Block rewards are assigned to a `Miner` actor at each block via `Reward.AwardBlockReward`.

Block rewards are vested over a period 

#### Block Reward

##### Win count


#### Gas fees

#### Vesting schedule

### Storage Rewards
#### Vesting schedule


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
