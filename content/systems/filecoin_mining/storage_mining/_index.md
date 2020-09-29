---
title: Storage Miner
bookCollapseSection: true
weight: 2
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Storage Miner

## Storage Mining Subsystem

The Filecoin Storage Mining Subsystem ensures a storage miner can effectively commit storage to the Filecoin protocol in order to both:

- Participate in the Filecoin [Storage Market](storage_market) by taking on client data and participating in storage deals.
- Participate in Filecoin [Storage Power Consensus](storage_power_consensus) by verifying and generating blocks to grow the Filecoin blockchain and earning block rewards and fees for doing so.

The above involves a number of steps to putting on and maintaining online storage, such as:

- Committing new storage (see [Sector](sector), [Sector Sealing](sector#sealing) and [PoRep](sdr))
- Continuously proving storage (see [Election PoSt](election_post))
- Declaring [storage faults](sector#sector-faults) and recovering from them.


![Sector State Machine](diagrams/sector_state_machine.dot)

![Sector State Machine Legend](diagrams/sector_state_machine_legend.dot)

## Filecoin Proofs

### Proof of Replication

A [Proof of Replication (PoRep)](sdr) is a proof that a Miner has correctly generated a unique _replica_ of some underlying data.

In practice, the underlying data is the _raw data_ contained in an Unsealed Sector, and a PoRep is a SNARK proof that the _sealing process_ was performed correctly to produce a Sealed Sector (See [Sealing a Sector](#Sealing-a-Sector)).

It is important to note that the replica should not only be unique to the miner, but also to the time when a miner has actually created the replica, i.e., sealed the sector. This means that if the same miner produces a sealed sector out of the same raw data twice, then this would count as a different replica.

When Miners commit to storing data, they must first produce a valid Proof of Replication.

#### Proof of Spacetime

A [Proof of Spacetime (aka PoSt)](post) is a long-term assurance of a Miner's continuous storage of their Sectors' data. _This is not a single proof,_ but a collection of proofs the Miner has submitted over time. Periodically, a Miner must add to these proofs by submitting a **WindowPoSt**:
* Fundamentally, a WindowPoSt is a collection of merkle proofs over the underlying data in a Miner's Sectors.
* WindowPoSts bundle proofs of various leaves across groups of Sectors (called **Partitions**).
* These proofs are submitted as a single SNARK.

The historical and ongoing submission of WindowPoSts creates assurance that the Miner has been storing, and continues to store the Sectors they agreed to store in the storage deal.

Once a Miner successfully adds and ProveCommits a Sector, the Sector is assigned to a Deadline: a specific window of time during which PoSts must be submitted. The day is broken up into 48 individual Deadlines of 30 minutes each, and ProveCommitted Sectors are assigned to one of these 48 Deadlines.
* PoSts may only be submitted for the currently-active Deadline. Deadlines are open for 30 minutes, starting from the Deadline's "Open" epoch and ending at its "Close" epoch.
* PoSts must incorporate randomness pulled from a random beacon. This randomness becomes publicly available at the Deadline's "Challenge" epoch, which is 20 epochs prior to its "Open" epoch.
* Deadlines also have a `FaultCutoff` epoch, 70 epochs prior to its "Open" epoch. After this epoch, Faults can no longer be declared for the Deadline's Sectors.

## Miner Accounting

A Miner's financial gain or loss is affected by the following three actions:
1. Miners deposit tokens to act as collateral for their `PreCommitted` and `ProveCommitted` Sectors
2. Miners earn tokens from block rewards, when they are elected to mine a new block and extend the blockchain.
3. Miners lose tokens if they fail to prove storage of a sector and are given penalties as a result.

### Balance Requirements

A Miner's token balance MUST cover ALL of the following:
* **PreCommit Deposits**: When a Miner PreCommits a Sector, they must supply a "precommit deposit" for the Sector, which acts as collateral. If the Sector is not ProveCommitted on time, this deposit is removed and burned.
* **Initial Pledge**: When a Miner ProveCommits a Sector, they must supply an "initial pledge" for the Sector, which acts as collateral. If the Sector is terminated, this deposit is removed and burned along with rewards earned by this sector up to a limit.
* **Locked Funds**: When a Miner receives tokens from block rewards, the tokens are locked and added to the Miner's vesting table to be unlocked linearly over some future epochs.

### Faults, Penalties and Fee Debt

A Sector's PoSts must be submitted on time, or that Sector is marked "faulty." There are two types of faults:
* **Declared fault**: When the Miner explicitly declares a Sector "faulty" _before_ its Deadline's FaultCutoff.
* **Undeclared fault**: When the Miner does not explicitly declare a Sector "faulty," but their submitted PoSt does not contain a proof for that Sector.

A Miner may accrue penalties for many reasons:
* **PreCommit Expiry Penalty**: Occurs if a Miner fails to `ProveCommit` a PreCommitted Sector in time. This happens the first time that a miner declares that it proves a sector and falls into the PoRep consensus.
* **Undeclared Fault Penalty**: Occurs if a Miner fails to submit a PoSt for a Sector on time.
* **Declared Fault Penalty**: Occurs if a Miner fails to submit a PoSt for a Sector on time, but they declare the Sector faulty before the system finds out (in which case the fault falls in the "Undeclared Fault Penalty" above). **This penalty fee is lower than the undeclared fault penalty**, in order to incentivize Miners to declare faults early.
* **Ongoing Fault Penalty**: Occurs every Proving Period a Miner fails to submit a PoSt for a Sector.
* **Termination Penalty**: Occurs if a Sector is forcibly terminated before its expiration.
* **Consensus Fault Penalty**: Occurs if a Miner commits a consensus fault and is reported.

When a Miner accrues penalties, the amount penalized is tracked as "Fee Debt." If a Miner has Fee Debt, they are restricted from certain actions until the amount owed is paid off. Miners with Fee Debt may not:
* PreCommit new Sectors
* Declare faulty Sectors "recovered"
* Withdraw balance

Faults are implied to be "temporary" - that is, a Miner that temporarily loses internet connection may choose to declare some Sectors for their upcoming Deadline as faulty, because the Miner knows they will regain the ability to submit proofs for those Sectors eventually. This declaration allows the Miner to still submit a valid proof for their Deadline (minus the faulty Sectors). This is very important for Miners, as missing a Deadline's PoSt entirely incurs a high penalty.
