---
title: Storage Power Actor
---

# `StoragePowerActor` interface

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor.id" code="true" lang="go" >}}

# `StoragePowerActorState` implementation

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor_state.go" code="true" lang="go" >}}

# `StoragePowerActorCode` implementation

{{< readfile file="/docs/actors/builtin/storage_power/storage_power_actor_code.go" code="true" lang="go" >}}

{{<label storage_power>}}
## On Power

Claimed power is assigned to every sector as a static function of its `SectorStorageWeightDesc` which includes `SectorSize`, `Duration`, and `DealWeight`. DealWeight is a measure that maps size and duration of active deals in a sector during its lifetime to its impact on power and reward distribution. A CommittedCapacity Sector will have a DealWeight of zero but all sectors have an explicit Duration which is defined from the ChainEpoch that the sector comes online in a ProveCommit message to the Expiration ChainEpoch of the sector. In principle, power is the number of votes a miner has in leader election and it is a point in time concept of storage. However, the exact function that maps `SectorStorageWeightDesc` to claimed `StoragePower` and `BlockReward` will be announced soon.

Nominal power is what goes into consensus from claimed power. This includes checks on whether miner's claimed power has met minimum miner size requirement for consensus, whether miner is in DetectedFault state which will result in loss of all power, and whether miner is under collateralized.

{{<label min_miner_size>}}
## Minimum Miner Size

In order to secure Storage Power Consensus, the system defines a minimum miner size required to participate in consensus.

Specifically, miners must have either at least `MIN_MINER_SIZE_STOR` of power (i.e. storage power currently used in storage deals) in order to participate in leader election. If no miner has `MIN_MINER_SIZE_STOR` or more power, miners with at least as much power as the smallest miner in the top `MIN_MINER_SIZE_TARG` of miners (sorted by storage power) will be able to participate in leader election. In plain english, take `MIN_MINER_SIZE_TARG = 3` for instance, this means that miners with at least as much power as the 3rd largest miner will be eligible to participate in consensus.

Miners smaller than this cannot mine blocks and earn block rewards in the network. Their power will not be counted as part of total network power. However, **it is important to note that such miners can still have their power faulted and be penalized accordingly**.

Accordingly, to bootstrap the network, the genesis block must include miners taking part in valid storage deals along with appropriate committed storage.

The `MIN_MINER_SIZE_TARG` condition will not be used in a network in which any miner has more than `MIN_MINER_SIZE_STOR` power. It is nonetheless defined to ensure liveness in small networks (e.g. close to genesis or after large power drops).

{{% notice placeholder %}}
The below values are currently placeholders.
{{% /notice %}}

We currently set:

- `MIN_MINER_SIZE_STOR = 1 << 40 Bytes` (100 TiB)
- `MIN_MINER_SIZE_TARG = 3

## Network recovery after halting

Placeholder where we will define a means of rebooting network liveness after it halts catastrophically (i.e. empty power table).

{{<label power_table>}}
# The Power Table

The portion of blocks a given miner generates through leader election in EC (and so the block rewards they earn) is proportional to their `Power Fraction` over time. That is, a miner whose storage represents 1% of total storage on the network should mine 1% of blocks on expectation.

SPC provides a power table abstraction which tracks miner power (i.e. miner storage in relation to network storage) over time. The power table is updated for new sector commitments (incrementing miner power), for failed PoSts (decrementing miner power) or for other storage and consensus faults. `_updatePowerEntriesFromClaimedPower` is called to update a particular miner's entry in the power table when its claimed power has changed.

Sector ProveCommit is the first time power is proven to the network and hence power is first added through `_rtAddPowerForSector` at `OnSectorProveCommit`. Power is also added when a sector's TemporaryFault period has ended. Miners are expected to prove over all their sectors that contribute to their power. `_rtDeductClaimedPowerForSectorAssert` is called to decrement a miner's power. This is called when a sector expires or invoked by miner through `OnSectorTerminate` and when a sector enters TemporaryFault through `OnSectorTemporaryFaultEffectiveBegin`. Both `_rtAddPowerForSector` and `_rtDeductClaimedPowerForSectorAssert` are currently called at `OnSectorModifyWeightDesc` as power is determined by `SectorStorageWeightDesc` and `SectorStorageWeightDesc` is only modified when a miner calls `ExtendSectorExpiration` to extend a sector's duration. This may or may not have an impact on power but the machinery is in place to preserve the flexibility.

The Miner lifecycle in the power table should be roughly as follows:

- MinerRegistration: A new miner with an associated worker public key and address is registered on the power table by the storage mining subsystem, along with their associated sector size (there is only one per worker).
- UpdatePower: These power increments and decrements are called by various storage actor (and must thus be verified by every full node on the network). Specifically:
    - Power is incremented at SectorProveCommit
    - All Power of a particular miner is decremented immediately after a missed SurprisePoSt (DetectedFault).
    - A particular sector's power is decremented `OnSectorTemporaryFaultEffectiveBegin` when TemporaryFault is declared on the sector.
    - A particular sector's power is added back `OnSectorTemporaryFaultEffectiveEnd` and miner is expected to prove over this sector. 
    - A particular sector's power is removed when the sector is terminated through sector expiration or miner invocation (`_rtTerminateSector`).
To summarize, only sectors in the Active state will command power. A Sector becomes Active after their first PoSt from Committed and Recovering stages. Power is immediately decremented when an Active Sector enters the Failing state (through DeclareFaults or Cron) and when an Active Sector expires.

{{<label pledge_collateral>}}
# Pledge Collateral

Consensus in Filecoin is secured in part by economic incentives enforced by Pledge Collateral.

Pledge collateral amount is committed based on power pledged to the system (i.e. proportional to number of sectors committed and sector size for a miner). It is a system-wide parameter and is committed to the `StoragePowerActor`. Pledge collateral can be posted by the `StorageMinerActor` at any time by a miner and its requirement is dependent on miner's power. Details around pledge collateral will be announced soon.

Pledge Collateral will be slashed when {{<sref consensus_faults>}} are reported to the `StoragePowerActor`'s `ReportConsensusFault` method, when a miner fails a SurprisePoSt (DetectedFault), or when a miner terminates a sector earlier than its duration.

Pledge Collateral is slashed for any fault affecting storage-power consensus, these include:

- faults to expected consensus in particular (see {{<sref consensus_faults>}}) which will be reported by a slasher to the `StoragePowerActor` in exchange for a reward.
- faults affecting consensus power more generally, specifically uncommitted power faults (i.e. {{<sref storage_faults>}}) which will be reported by the `CronActor` automatically or when a miner terminates a sector earlier than its promised duration.
