---
menuTitle: Storage Miner
statusIcon: 🔁
title: Storage Miner
entries:
  - mining_cycle
  - storage_miner_actor
---

{{<label storage_mining_subsystem>}}
# Filecoin Storage Mining Subsystem

The Filecoin Storage Mining Subsystem ensures a storage miner can effectively commit storage to the Filecoin protocol in order to both:

- participate in the Filecoin {{<sref storage_market>}} by taking on client data and participating in storage deals.
- participate in Filecoin {{<sref storage_power_consensus>}}, verifying and generating blocks to grow the Filecoin blockchain and earning block rewards and fees for doing so.

The above involves a number of steps to putting on and maintaining online storage, such as:

- Committing new storage (see Sealing and PoRep)
- Continously proving storage (see {{<sref election_post>}})
- Declaring storage faults and recovering from them.

## Sector Types

There are two types of sectors, Regular Sectors with storage deals in them and Committed Capacity (CC) Sectors with no deals. All sectors require an expiration epoch that is declared upon PreCommit and sectors are assigned a StartEpoch at ProveCommit. Start and Expiration epoch collectively define the lifetime of a Sector. Length and size of active deals in a sector's lifetime determine the `DealWeight` of the sector. `SectorSize`, `Duration`, and `DealWeight` statically determine the power assigned to a sector that will remain constant throughout its lifetime. More details on cost and reward for different sector types will be announced soon.

## Sector States

When managing their storage {{<sref sector "sectors">}}  as part of Filecoin mining, storage providers will account for where in the {{<sref mining_cycle>}} their sectors are. For instance, has a sector been committed? Does it need a new PoSt? Most of these operations happen as part of cycles of chain epochs called `Proving Period`s each of which yield high confidence that every miner in the chain has proven their power (see {{<sref election_post>}}).

There are three states that an individual sector can be in:

- `PreCommit` when a sector has been added through a PreCommit message.
- `Active` when a sector has been proven through a ProveCommit message and when a sector's TemporaryFault period has ended.
- `TemporaryFault` when a miner declares fault on a particular sector.

Sectors enter `Active` from `PreCommit` through a ProveCommit message that serves as the first proof for the sector. PreCommit requires a PreCommit deposit which will be returned upon successful and timely ProveCommit. However, if there is no matching ProveCommit for a particular PreCommit message, the deposit will be burned at PreCommit expiration.

A particular sector enters `TemporaryFault` from `Active` through `DeclareTemporaryFault` with a specified period. Power associated with the sector will be lost immediately and miner needs to pay a `TemporaryFaultFee` determined by the power suspended and the duration of suspension. At the end of the declared duration, faulted sectors automatically regain power and enter `Active`. Miners are expected to prove over this recovered sector. Failure to do so may result in failing ElectionPoSt or `DetectedFault` from failing SurprisePoSt. 

{{< diagram src="diagrams/sector_state_machine.dot.svg" title="Sector State Machine" >}}

{{< diagram src="diagrams/sector_state_machine_legend.dot.svg" title="Sector State Machine Legend" >}}

### Miner PoSt State

`MinerPoStState` keeps track of a miner's state in responding to PoSt and there are three states in `MinerPoStState`:

- `OK` miner has passed either a ElectionPoSt or a SurprisePoSt sufficiently recently.
- `Challenged` miner has been selected to prove its storage via SurprisePoSt and is currently in the Challenged state
- `DetectedFault` miner has failed at least one SurprisePoSt, indicating that all claimed storage may not be proven. Miner has lost power on its sector and recovery can only proceed by a successful response to a subsequent SurprisePoSt challenge, up until the limit of number of consecutive failures.

`DetectedFault` is a miner-wide PoSt state when all sectors are considered inactive. All power is lost immediately and pledge collateral is slashed. If a miner remains in `DetectedFault` for more than MaxConsecutiveFailures, all sectors will be terminated, both power and market actors will be notified.

`ProvingSet` consists of sectors that miners are required to generate proofs against and is what counts towards miners' power. In other words, `ProvingSet` is a set of all `Active` sectors for a particular miner. `ProvingSet` is only relevant when the miner is in OK stage of its `MinerPoStState`. When a miner is in the `Challenged` state, `ChallengedSectors` specify the list of sectors to be challenged which is the `ProvingSet` before the challenge is issued thus allowing more sectors to be added while it is in the `Challenged` state.

{{< diagram src="diagrams/miner_post_state_machine.dot.svg" title="Miner PoSt State Machine" >}}

{{< diagram src="diagrams/miner_post_state_machine_legend.dot.svg" title="Miner PoSt State Machine Legend" >}}

