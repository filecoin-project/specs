---
menuTitle: Storage Miner
statusIcon: 🔁
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

# Sector in StorageMiner State Machine (new one)

{{< diagram src="diagrams/sector_state_fsm.dot.svg" title="Sector State (new one)" >}}

{{< diagram src="diagrams/sector_state_legend.dot.svg" title="Sector State Legend (new one)" >}}

# Sector in StorageMiner State Machine (both)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}

# ElectionPoSt

To enable short PoSt response time, miners are required submit a PoSt when they are elected to mine a block, hence PoSt on election or ElectionPoSt. When miners win a block, they need to immediately generate a PoStProof and submit that along with the ElectionProof. Both the ElectionProof and PoStProof are checked at Block Validation by `StoragePowerConsenusSubsystem`. When a block is included, a special message is added that calls `SubmitElectionPoSt` which will process sector updates in the same way as successful `SubmitSurprisePoSt` do. Note that it will be a single PoStProof in a single `SubmitElectionPoSt` message instead of a Commit-then-Prove scheme.

# SurprisePoSt

The chain keeps track of the last time that a miner received and submitted a Challenge (Election or Surprise). It randomly samples X miners to surprise at every Epoch. For every miner sampled, a `NotifyOfPoStSurpriseChallenge` is issued and sent as an on-chain message to surprised `StorageMinerActor`. However, if the `StorageMinerActor` is already proving a SurprisePoStChallenge (`IsChallenged` is True) or the `StorageMinerActor` has received a challenge (by Election or Surprise) within the `MIN_CHALLENGE_PERIOD` (`ShouldChallenge` is False), then the PoSt surprise notification will be ignored and return success.
