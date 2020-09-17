---
title: Storage Power Actor
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Storage Power Actor

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/power/power_state.go" lang="go" title="StoragePowerActorState implementation" symbol="State">}}


{{<embed src="github:filecoin-project/specs-actors/actors/builtin/power/power_actor.go" lang="go" title="StoragePowerActor implementation" >}}

## The Power Table

The portion of blocks a given miner generates through leader election in EC (and so the block rewards they earn) is proportional to their `Power Fraction` over time. That is, a miner whose storage represents 1% of total storage on the network should mine 1% of blocks on expectation.

SPC provides a power table abstraction which tracks miner power (i.e. miner storage in relation to network storage) over time. The power table is updated for new sector commitments (incrementing miner power), for failed PoSts (decrementing miner power) or for other storage and consensus faults.

Sector ProveCommit is the first time power is proven to the network and hence power is first added upon successful sector ProveCommit. Power is also added when a sector's TemporaryFault period has ended. Miners are expected to prove over all their sectors that contribute to their power. 

Power is decremented when a sector expires, when a sector enters TemporaryFault, or when it is invoked by miners through Sector Termination. Miners can also extend the lifetime of a sector through `ExtendSectorExpiration` and thus modifying `SectorStorageWeightDesc`. This may or may not have an impact on power but the machinery is in place to preserve the flexibility.

The Miner lifecycle in the power table should be roughly as follows:

- MinerRegistration: A new miner with an associated worker public key and address is registered on the power table by the storage mining subsystem, along with their associated sector size (there is only one per worker).
- UpdatePower: These power increments and decrements are called by various storage actor (and must thus be verified by every full node on the network). Specifically:
    - Power is incremented at SectorProveCommit
    - All Power of a particular miner is decremented immediately after a missed SurprisePoSt (DetectedFault).
    - A particular sector's power is decremented when its TemporaryFault begins.
    - A particular sector's power is added back when its TemporaryFault ends and miner is expected to prove over this sector. 
    - A particular sector's power is removed when the sector is terminated through sector expiration or miner invocation.

To summarize, only sectors in the Active state will command power. A Sector becomes Active when it is added upon ProveCommit. Power is immediately decremented upon when TemporaryFault begins on an Active sector or when the miner is in Challenged or DetectedFault state. Power will be restored when TemporaryFault has ended and when the miner successfully responds to a SurprisePoSt challenge. A sector's power is removed when it is terminated through either miner invocation or normal expiration. 

## Pledge Collateral

Consensus in Filecoin is secured in part by economic incentives enforced by Pledge Collateral.

Pledge collateral amount is committed based on power pledged to the system (i.e. proportional to number of sectors committed and sector size for a miner). It is a system-wide parameter and is committed to the `StoragePowerActor`. Pledge collateral can be posted by the `StorageMinerActor` at any time by a miner and its requirement is dependent on miner's power. Details around pledge collateral will be announced soon.

Pledge Collateral will be slashed when [Consensus Faults](expected_consensus#consensus-faults) are reported to the `StoragePowerActor`'s `ReportConsensusFault` method, when a miner fails a SurprisePoSt (DetectedFault), or when a miner terminates a sector earlier than its duration.

Pledge Collateral is slashed for any fault affecting storage-power consensus, these include:

- faults to expected consensus in particular (see [Consensus Faults](expected_consensus#consensus-faults))  which will be reported by a slasher to the `StoragePowerActor` in exchange for a reward.
- faults affecting consensus power more generally, specifically uncommitted power faults (i.e. [Storage Faults](faults#storage-faults)) which will be reported by the `CronActor` automatically or when a miner terminates a sector earlier than its promised duration.
