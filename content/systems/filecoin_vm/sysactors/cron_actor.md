---
title: CronActor
weight: 2
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# CronActor
---

The `CronActor` is responsible for scheduling jobs for several operations that need to take place between actors. This is primarily realised in the form of sending messages to other registered actors at the end of every epoch. It interfaces with the Storage Power Actor on `OnEpochTickEnd` to check:
- `CronEventPreCommitExpiry`: if the miner has submitted their `ProveCommit` proofs (given enough time for the miner to complete the proofs of replication and finalize the commit phase). If the miner has completed their commit phase on time, they can claim the `PreCommiDeposit` back. If not, the `PreCommitDeposit` is lost and no power is added for this miner to the power table.
- `CronEventProvingPeriod`: the proving period identified,
- `CronEventWorkerKeyChange`: if a miner has submitted a worker key address change.


{{<embed src="/modules/actors/builtin/cron/cron_actor.go" lang="go">}}