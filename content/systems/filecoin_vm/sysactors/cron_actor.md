---
title: CronActor
weight: 2
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# CronActor
---

The `CronActor` is responsible for sending messages to other registered actors at the end of every epoch. It interfaces with the Storage Power Actor on `OnEpochTickEnd` to check:
- `CronEventPreCommitExpiry`: if the miner has submitted their `ProveCommit` proofs,
- `CronEventProvingPeriod`: the proving period identified,
- `CronEventWorkerKeyChange`: if a miner has submitted a worker key address change.


```go
// Invoked by the system after all other messages in the epoch have been processed.
func (a Actor) EpochTick(rt vmr.Runtime, _ *adt.EmptyValue) *adt.EmptyValue {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

	var st State
	rt.State().Readonly(&st)
	for _, entry := range st.Entries {
		_, _ = rt.Send(entry.Receiver, entry.MethodNum, nil, abi.NewTokenAmount(0))
		// Any error and return value are ignored.
	}

	return nil
}
```

