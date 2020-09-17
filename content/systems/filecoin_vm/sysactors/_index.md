---
title: System Actors
weight: 6
bookCollapseSection: true
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# System Actors

There are eleven (11) builtin System Actors in total, but not all of them interact with the VM. Each actor is identified by a _Code ID_ (or CID).

There are two system actors required for VM processing:
  - the [InitActor](sysactors#initactor), which initializes new actors and records the network name, and
  - the [CronActor](sysactors#cronactor), a scheduler actor that runs critical functions at every epoch.
There are another two actors that interact with the VM:
  - the [AccountActor](sysactors#accountactor) responsible for user accounts (non-singleton), and
  - the [RewardActor](sysactors#rewardactor) for block reward and token vesting (singleton).


The remaining seven (7) builtin System Actors that do not interact directly with the VM are the following:

- `StorageMarketActor`: responsible for managing storage and retrieval deals [[Market Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/market/market_actor.go)]
- `StorageMinerActor`: actor responsible to deal with storage mining operations and collect proofs [[Storage Miner Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/miner/miner_actor.go)]
- `MultisigActor` (or Multi-Signature Wallet Actor): responsible for dealing with operations involving the Filecoin wallet [[Multisig Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/multisig/multisig_actor.go)]
- `PaymentChannelActor`: responsible for setting up and settling funds related to payment channels [[Paych Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/paych/paych_actor.go)]
-  `StoragePowerActor`: responsible for keeping track of the storage power allocated at each storage miner [[Storage Power Actor](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/power/power_actor.go)]
- `VerifiedRegistryActor`: responsible for managing verified clients [[Verifreg Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/verifreg/verified_registry_actor.go)]
- `SystemActor`: general system actor [[System Actor Repo](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/system/system_actor.go)]

## CronActor

Built in to the genesis state, the `CronActor`'s dispatch table invokes the `StoragePowerActor` and `StorageMarketActor` for them to maintain internal state and process deferred events. It could in principle invoke other actors after a network upgrade.

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/cron/cron_actor.go"  lang="go">}}

## InitActor

The `InitActor` has the power to create new actors, e.g., those that enter the system. It maintains a table resolving a public key and temporary actor addresses to their canonical ID-addresses. Invalid CIDs should not get committed to the state tree.

Note that the canonical ID address does not persist in case of chain re-organization. The actor address or public key survives chain re-organization.

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/init/init_actor.go" lang="go">}}

## RewardActor

The `RewardActor` is where unminted Filecoin tokens are kept. The actor distributes rewards directly to miner actors, where they are locked for vesting. The reward value used for the current epoch is updated at the end of an epoch through a cron tick.

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/reward/reward_actor.go"  lang="go">}}

## AccountActor

The `AccountActor` is responsible for user accounts. Account actors are not created by the `InitActor`, but their constructor is called by the system. Account actors are created by sending a message to a public-key style address. The address must be `BLS` or `SECP`, or otherwise there should be an exit error. The account actor is updating the state tree with the new actor address.

{{<embed src="github:filecoin-project/specs-actors/actors/builtin/account/account_actor.go" lang="go" >}}