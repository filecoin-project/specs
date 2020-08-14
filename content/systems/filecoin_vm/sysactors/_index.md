---
title: System Actors
weight: 6
bookCollapseSection: true
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# System Actors
---

There are eleven (11) builtin System Actors in total, but not all of them interact with the VM.

```go
type BuiltinActor struct {
	actor abi.Invokee
	code  cid.Cid
}
```

Each actor is identified by a _Code ID_ (or CID), not to be confused with the IPFS-style _Content ID_, (also CID).

```go
// Code is the CodeID (cid) of the actor.
func (b BuiltinActor) Code() cid.Cid {
	return b.code
}
```

The eleven (11) builtin System Actors are the following:

```go

func BuiltinActors() []BuiltinActor {
	return []BuiltinActor{
		{
			// Account Actor: manages user accounts
			actor: account.Actor{},
			code:  builtin.AccountActorCodeID,
		},
		{
			// Cron Actor: a scheduler actor that runs critical functions at every epoch
			actor: cron.Actor{},
			code:  builtin.CronActorCodeID,
		},
		{
			// Init Actor: responsible to initialize new actors
			actor: init_.Actor{},
			code:  builtin.InitActorCodeID,
		},
		{
			// Market Actor: actor responsible to manage storage and retrieval deals
			actor: market.Actor{},
			code:  builtin.StorageMarketActorCodeID,
		},
		{
			// Storage Miner Actor: actor responsible to deal storage mining operations, collect proofs, etc.
			actor: miner.Actor{},
			code:  builtin.StorageMinerActorCodeID,
		},
		{
			// The Multi-Signature Wallet Actor
			actor: multisig.Actor{},
			code:  builtin.MultisigActorCodeID,
		},
		{
			// Payment Channel Actor: responsible for setting up and settling funds related to payment channels.
			actor: paych.Actor{},
			code:  builtin.PaymentChannelActorCodeID,
		},
		{
			// Storage Power Consensus Actor
			actor: power.Actor{},
			code:  builtin.StoragePowerActorCodeID,
		},
		{
			// Reward Actor: responsible for managing block rewards
			actor: reward.Actor{},
			code:  builtin.RewardActorCodeID,
		},
		{
			// System Actor
			actor: system.Actor{},
			code:  builtin.SystemActorCodeID,
		},
		{
			// Verified Registry Actor
			actor: verifreg.Actor{},
			code:  builtin.VerifiedRegistryActorCodeID,
		},
	}
}

```

There are two system actors required for VM processing:
  - the [InitActor](init_actor.md), which initializes new actors and records the network name, and
  - the [CronActor](cron_actor.md), a scheduler actor that runs critical functions at every epoch.
There are another two actors that interact with the VM:
  - the [AccountActor](account_actor.md) responsible for user accounts (non-singleton), and
  - the [RewardActor](reward_actor.md) for block reward and token vesting (singleton).
