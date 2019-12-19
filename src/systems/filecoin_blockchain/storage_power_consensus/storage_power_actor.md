---
title: Storage Power Actor
---

# `StoragePowerActor` interface

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor.id" code="true" lang="go" >}}

# `StoragePowerActorState` implementation

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor_state.go" code="true" lang="go" >}}

# `StoragePowerActorCode` implementation

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor_code.go" code="true" lang="go" >}}


{{<label power_table>}}
# The Power Table

The portion of blocks a given miner generates through leader election in EC (and so the block rewards they earn) is proportional to their `Power Fraction` over time. That is, a miner whose storage represents 1% of total storage on the network should mine 1% of blocks on expectation.

SPC provides a power table abstraction which tracks miner power (i.e. miner storage in relation to network storage) over time. The power table is updated for new sector commitments (incrementing miner power), when PoSts fail to be put on-chain (decrementing miner power) or for other storage and consensus faults.

An invariant of the storage power consensus subsystem is that all storage in the power table must be verified. That is, miners can only derive power from storage they have already proven to the network.

In order to achieve this, Filecoin delays updating power for new sector commitments until the first valid PoSt in the next proving period corresponding to that sector. (TODO: potential delay this further in order to ensure that any power cut goes undetected at most as long as the shortest power delay on new sector commitments).

For instance, say a miner X does the following:
- In epoch 100: commits 10 TB
- In epoch 110: publishes a PoSt for their storage
- In epoch 120: commits another 10TB
- In epoch 135: publishes a new PoSt for their storage

Querying the power table for this miner at different rounds should yield (using the following shorthand as an illustration only):
- `Power(X, 90) == 0`
- `Power(X, 100) == 0`
- `Power(X, 110) == 0`
- `Power(X, 111) == 10`
- `Power(X, 120) == 10`
- `Power(X, 135) == 10`
- `Power(x, 136) == 20`

Conversely, storage faults only lead to power loss once they are detected (up to one proving period after the fault) so miners will mine with no more power than they have used to store data over time.

Put another way, power accounting in the SPC is delayed between storage being proven or faulted, and power being updated in the power table (and so for leader election). This ensures fairness over time.

The Miner lifecycle in the power table should be roughly as follows:

- MinerRegistration: A new miner with an associated worker public key and address is registered on the power table by the storage mining subsystem, along with their associated sector size (there is only one per worker).
- UpdatePower: These power increments and decrements are called by various storage actor (and must thus be verified by every full node on the network). Specifically:
    - Power is incremented to account for a new SectorCommitment at the first PoSt past the first ProvingPeriod.
    - All Power is decremented immediately after a missed PoSt.
    - Power is decremented immediately after faults are declared, proportional to the faulty sector size.
    - Power is incremented after a PoSt recovering from a fault.
    - Power is definitively removed from the Power Table past the sector failure timeout (see {{<sref storage_faults>}})
To summarize, only sectors in the Active state will command power. A Sector becomes Active after their first PoSt from Committed and Recovering stages. Power is immediately decremented when an Active Sector enters the Failing state (through DeclareFaults or Cron) and when an Active Sector expires.

{{<label pledge_collateral>}}
# Pledge Collateral

Consensus in Filecoin is secured in part by economic incentives enforced by Pledge Collateral.

Pledge collateral amount is committed based on power pledged to the system (i.e. proportional to number of sectors committed and sector size for a miner). It is a system-wide parameter and is committed to the `StoragePowerActor`. TODO: define parameter value. Pledge Collateral submission methods take on storage deals to determine the appropriate amount of collateral to be pledged. Pledge collateral can be posted by the `StorageMinerActor` at any time by a miner up to sector commitments. A sector commitment without the requisite posted pledge collateral will be deemed invalid.

Pledge Collateral will be slashed when {{<sref consensus_faults>}} are reported to the `StoragePowerActor`'s `ReportConsensusFault` method or when the `CronActor` calls the `StoragePowerActor`'s `ReportUncommittedPowerFault` method.

Pledge Collateral is slashed for any fault affecting storage-power consensus, these include:

- faults to expected consensus in particular (see {{<sref consensus_faults>}}) which will be reported by a slasher to the `StoragePowerActor` in exchange for a reward.
- faults affecting consensus power more generally, specifically uncommitted power faults (i.e. {{<sref storage_faults>}}) which will be reported by the `CronActor` automatically.
