---
title: RewardActor
weight: 4
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# RewardActor
---

The `RewardActor` is where unminted Filecoin tokens are kept. The `RewardActor` contains a `RewardMap` which is a mapping from owner addresses to `Reward` structs. 

A `Reward struct` is created to preserve the flexibility of introducing block reward vesting into the protocol. `MintReward` creates a new `Reward struct` and adds it to the `RewardMap`.

The award value used for the current epoch is updated at the end of an epoch through a cron tick. In the case previous epochs were null blocks this is the reward value as calculated at the last non-null epoch.


{{<embed src="/modules/actors/builtin/reward/reward_actor.go" lang="go">}}