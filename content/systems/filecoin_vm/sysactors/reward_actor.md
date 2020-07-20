---
title: RewardActor
---

# RewardActor
---

RewardActor is where unminted Filecoin tokens are kept. RewardActor contains a `RewardMap` which is a mapping from owner addresses to `Reward` structs. 

`Reward` struct is created to preserve the flexibility of introducing block reward vesting into the protocol. `MintReward` creates a new `Reward` struct and adds it to the `RewardMap`. 

A `Reward` struct contains a `StartEpoch` that keeps track of when this `Reward` is created, `Value` that represents the total number of tokens rewarded, and `EndEpoch` which is when the reward will be fully vested. `VestingFunction` is currently an enum to represent the flexibility of different vesting functions. `AmountWithdrawn` records how many tokens have been withdrawn from a `Reward` struct so far. Owner addresses can call `WithdrawReward` which will withdraw all vested tokens that the investor address has from the RewardMap so far. When `AmountWithdrawn` equals `Value` in a `Reward` struct, the `Reward` struct will be removed from the `RewardMap`.

{{<embed src="/specs-actors/actors/builtin/reward/reward_actor.go"  lang="go">}}
