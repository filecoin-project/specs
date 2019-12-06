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
- Continously proving storage ({see {<sref election_post>}})
- Declaring storage faults and recovering from them.

## Storage States

When managing their storage {{<sref sector "sectors">}}  as part of Filecoin mining, storage providers will account for where in the mining cycle their sectors are. For instance, has a sector been committed? Does it need a new PoSt? Most of these operations happen as part of cycles of chain epochs called `Proving Period`s each of which yield high confidence that every miner in the chain has proven their power (see {{<sref election_post>}}).

Through most of this cycle, sectors will be part of the miner's `Proving Set`. This is the set of sectors that the miners are supposed to generate proofs against. It includes:
- `Committed Sectors` which have an associated PoRep but have but yet to generate a PoSt, 
- `Active Sectors` which have successfully been proven and are used in {{<sref storage_power_consensus>}} (SPC), and
- `Recovering Sectors` which were faulted and are poised to recover through a new PoSt.

Sectors in the `Proving Set` can be faulted and marked as `Failing` by the miner itself declaring a fault (`Declared Faults`), or through automatic fault detection by the network using on-chain data (`Detected Faults`).

A sector that is in the `Failing` state for three consecutive `Proving Period`s will be terminated (`Terminated Faults`) meaning its data will be deemed unrecoverable by PoSt and the sector will have to be SEALed once more (as part of PoRep).

Conversely `RecoverFaults()` can be called any time by the miner on a failing sector to return it to the `ProvingSet` and attempt to prove data is being stored once more. For instance an Active sector might move into the failing state during a power outage (through a declared or detected fault). At the end of the outage, the miner may call `RecoverFaults` to transition the state to `Recovering` before proving it once more and returning it to `Active`. 

Through the above workflows, Filecoin (or Storage more generally) can more generally be described as being **Active** (`Active` sectors), or **Inactive** (`Committed`, `Recovering` and `Failing` sectors). Specifically, active storage is used in SPC (and in-deal active storage counts toward {{<sref storage_power>}}) whereas inactive storage is not. Both active and inactive storage is tracked by the system and affect a miner's ability to serve as a storage provider in the Filecoin Storage Market. 

We illustrate these states below.

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}

### Sector in StorageMiner State Machine (new one)

{{< diagram src="diagrams/sector_state_fsm.dot.svg" title="Sector State (new one)" >}}

{{< diagram src="diagrams/sector_state_legend.dot.svg" title="Sector State Legend (new one)" >}}

### Sector in StorageMiner State Machine (both)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}

