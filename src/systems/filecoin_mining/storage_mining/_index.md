---
menuTitle: Storage Miner
title: Storage Miner
entries:
  - mining_cycle
  - storage_miner_actor
  - mining_scheduler
---

{{<label storage_mining_subsystem>}}

TODO:

- rename "Storage Mining Worker" ?

Filecoin Storage Mining Subsystem

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}

# Sector in StorageMiner State Machine (on chain)

{{< diagram src="diagrams/sector_chain_fsm.dot.svg" title="Sector State Machine (on chain)" >}}

# Sector in StorageMiner State Machine (off chain)

{{< diagram src="diagrams/sector_offchain_fsm.dot.svg" title="Sector State Machine (off chain)" >}}

# Sector in StorageMiner State Machine (both)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}
