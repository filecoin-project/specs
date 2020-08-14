---
title: RewardActor
weight: 4
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# RewardActor
---

The `RewardActor` is where unminted Filecoin tokens are kept. The `RewardActor` contains a `RewardMap` which is a mapping from owner addresses to `Reward` structs. 

A `Reward` struct is created to preserve the flexibility of introducing block reward vesting into the protocol. `MintReward` creates a new `Reward` struct and adds it to the `RewardMap`.

```go
type AwardBlockRewardParams struct {
	Miner     address.Address
	Penalty   abi.TokenAmount // penalty for including bad messages in a block, >= 0
	GasReward abi.TokenAmount // gas reward from all gas fees in a block, >= 0
	WinCount  int64           // number of reward units won, > 0
}
```

A `Reward` struct contains a `StartEpoch` that keeps track of when this `Reward` is created, a `Value` that represents the total number of tokens rewarded, and an `EndEpoch` which is when the reward will be fully vested. `VestingFunction` is currently an enum to represent the flexibility of different vesting functions. The `AmountWithdrawn` records how many tokens have been withdrawn from a `Reward` struct so far. Owner addresses can call `WithdrawReward` which will withdraw all vested tokens that the investor address has from the RewardMap so far. When `AmountWithdrawn` equals `Value` in a `Reward` struct, the `Reward` struct will be removed from the `RewardMap`.


The award value used for the current epoch is updated at the end of an epoch through a cron tick.  In the case previous epochs were null blocks this is the reward value as calculated at the last non-null epoch.

```go
type ThisEpochRewardReturn struct {
	ThisEpochReward         abi.TokenAmount
	ThisEpochRewardSmoothed *smoothing.FilterEstimate
	ThisEpochBaselinePower  abi.StoragePower
}

func (a Actor) ThisEpochReward(rt vmr.Runtime, _ *adt.EmptyValue) *ThisEpochRewardReturn {
	rt.ValidateImmediateCallerAcceptAny()

	var st State
	rt.State().Readonly(&st)
	return &ThisEpochRewardReturn{
		ThisEpochReward:         st.ThisEpochReward,
		ThisEpochBaselinePower:  st.ThisEpochBaselinePower,
		ThisEpochRewardSmoothed: st.ThisEpochRewardSmoothed,
	}
}
```