---
title: Storage Power Actor
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Storage Power Actor

## `StoragePowerActorState` implementation

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/power/power_state.go" lang="go" symbol="State">}}

## `StoragePowerActor` implementation

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/power/power_actor.go" lang="go" symbol="Exports">}}

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/power/power_actor.go" lang="go" symbol="Constructor">}}

## The Power Table

The portion of blocks a given miner generates through leader election in EC (and so the block rewards they earn) is proportional to their `Quality-Adjusted Power Fraction` over time. That is, a miner whose quality adjusted power represents 1% of total quality adjusted power on the network should mine 1% of blocks on expectation.

SPC provides a power table abstraction which tracks miner power (i.e. miner storage in relation to network storage) over time. The power table is updated for new sector commitments (incrementing miner power), for failed PoSts (decrementing miner power) or for other storage and consensus faults.

Sector ProveCommit is the first time power is proven to the network and hence power is first added upon successful sector ProveCommit. Power is also added when a sector is declared as recovered. Miners are expected to prove over all their sectors that contribute to their power.

Power is decremented when a sector expires, when a sector is declared or detected to be faulty, or when it is terminated through miner invocation. Miners can also extend the lifetime of a sector through `ExtendSectorExpiration`.

The Miner lifecycle in the power table should be roughly as follows:

- `MinerRegistration`: A new miner with an associated worker public key and address is registered on the power table by the storage mining subsystem, along with their associated sector size (there is only one per worker).
- `UpdatePower`: These power increments and decrements are called by various storage actors (and must thus be verified by every full node on the network). Specifically:
  - Power is incremented at ProveCommit, as a subcall of `miner.ProveCommitSector` or `miner.ProveCommitAggregate`
  - Power of a partition is decremented immediately after a missed WindowPoSt (`DetectedFault`).
  - A particular sector's power is decremented when it enters into a faulty state either through Declared Faults or Skipped Faults.
  - A particular sector's power is added back after recovery is declared and proven by PoSt.
  - A particular sector's power is removed when the sector is expired or terminated through miner invocation.

To summarize, only sectors in the Active state will command power. A Sector becomes Active when it is added upon `ProveCommit`. Power is immediately decremented when it enters into the faulty state. Power will be restored when its declared recovery is proven. A sector's power is removed when it is expired or terminated through miner invocation.

## Pledge Collateral

Pledge Collateral is slashed for any fault affecting storage-power consensus, these include:

- faults to expected consensus in particular (see [Consensus Faults](expected_consensus#consensus-faults)), which will be reported by a slasher to the `StoragePowerActor` in exchange for a reward.
- faults affecting consensus power more generally, specifically uncommitted power faults (i.e. [Storage Faults](faults)), which will be reported by the `CronActor` automatically or when a miner terminates a sector earlier than its promised duration.

For a more detailed discussion on Pledge Collateral, please see the [Miner Collaterals section](filecoin_mining#miner_collaterals).
