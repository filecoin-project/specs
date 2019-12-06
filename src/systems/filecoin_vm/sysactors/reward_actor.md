---
title: RewardActor
---

RewardActor is where unminted and unvested Filecoin tokens are kept. At genesis, RewardActor is initiailized with investor accounts, tokens, and vesting schedule in a `RewardMap` which is a mapping from owner addressws to `Reward` structs. A `Reward` struct contains a `StartEpoch` that keeps track of when this `Reward` is created, `Value` that represents the total number of tokens rewarded, and `ReleaseRate` which is the linear rate of release in the unit of FIL per Epoch. `WithdrawAmount` records how many tokens have been withdrawn from a `Reward` struct so far. Owner addresses can call `WithdrawReward` which will withdraw all vested tokens that the investor address has from the RewardMap so far. When `WithdrawAmount` equals `Value` in a `Reward` struct, the `Reward` struct will be removed from the `RewardMap`.

`RewardMap` is also used in block reward minting to preserve the flexibility of introducing block reward vesting to the protocol. `MintReward` creates a new `Reward` struct and adds it to the `RewardMap`.

{{< readfile file="reward_actor.id" code="true" lang="go" >}}

{{< readfile file="reward_actor.go" code="true" lang="go" >}}
