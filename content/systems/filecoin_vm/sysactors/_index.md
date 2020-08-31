---
title: System Actors
weight: 6
bookCollapseSection: true
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# System Actors

There are eleven (11) builtin System Actors in total, but not all of them interact with the VM. Each actor is identified by a _Code ID_ (or CID).

There are two system actors required for VM processing:
  - the [InitActor](init_actor.md), which initializes new actors and records the network name, and
  - the [CronActor](cron_actor.md), a scheduler actor that runs critical functions at every epoch.
There are another two actors that interact with the VM:
  - the [AccountActor](account_actor.md) responsible for user accounts (non-singleton), and
  - the [RewardActor](reward_actor.md) for block reward and token vesting (singleton).


The remaining seven (7) builtin System Actors that do not interact directly with the VM are the following:

- `StorageMarketActor`: responsible for managing storage and retrieval deals [[Market Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/market/market_actor.go)]
- `StorageMinerActor`: actor responsible to deal with storage mining operations and collect proofs [[Storage Miner Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/miner/miner_actor.go)]
- `MultisigActor` (or Multi-Signature Wallet Actor): responsible for dealing with operations involving the Filecoin wallet [[Multisig Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/multisig/multisig_actor.go)]
- `PaymentChannelActor`: responsible for setting up and settling funds related to payment channels [[Paych Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/paych/paych_actor.go)]
-  `StoragePowerActor`: responsible for keeping track of the storage power allocated at each storage miner [[Storage Power Actor](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/power/power_actor.go)]
- `VerifiedRegistryActor`: responsible for managing verified clients [[Verifreg Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/verifreg/verified_registry_actor.go)]
- `SystemActor`: general system actor [[System Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/system/system_actor.go)]
