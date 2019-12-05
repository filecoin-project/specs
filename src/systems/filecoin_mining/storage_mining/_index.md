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

When managing their storage sectors as part of Filecoin mining, storage providers will account for where in the mining cycle their sectors are. For instance, has a sector been committed? Does it need a new PoSt?

Through this cycle, sectors will fall into either a miner's:
- *PreCommitted Sectors*: New sectors to be committed using PoRep
- `Proving Set`: Committed sectors to be proven over time using PoSts
  - Sectors can be faulted by the miner itself (`Declared Faults`).
- `Fault Set`: Committed sectors that have failed to be proven.
  - Fault recovery allows sectors to return to the `Proving Set`.

Through the above workflows, Filecoin {{<sref sector "Sectors">}} (or Storage more generally) can be describede as beeing in the following states:

- *Not-on-Fil*: This includes *Precommitted sectors* that have yet to be proven (and do not account for storage in Filecoin).
- **Active**: This includes sectors in the `Proving Set`.
- **Inactive**: This includes sectors in the `Fault Set` (i.e. through `Declared Faults`) and faults detected by the system (`Detected Faults`) which render all of a miner's storage inactive.

All of these states end up affecting both a miner's ability to serve as a storage provider in the Filecoin Storage Market, and their ability to participate in Filecoin consensus through Storage Power Consensus (see {{<sref storage_power>}}). We illustrate this below.

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}

### Sector in StorageMiner State Machine (new one)

{{< diagram src="diagrams/sector_state_fsm.dot.svg" title="Sector State (new one)" >}}

{{< diagram src="diagrams/sector_state_legend.dot.svg" title="Sector State Legend (new one)" >}}

### Sector in StorageMiner State Machine (both)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}

