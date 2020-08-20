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

The `CronActor` is responsible for scheduling jobs for several operations that need to take place between actors. The `CronActor` can provide services to other actors related to task scheduling, such as `OnEpochTickEnd` call the `InitActor` to create a new actor, contact the `StoragePowerActor` to check if proofs have been submitted, or call the `RewardActor` to update parameters. This is primarily realised in the form of sending messages to other registered actors at the end of every epoch. The `CronActor` interfaces with the Storage Power Actor on `OnEpochTickEnd` to check:
- `CronEventPreCommitExpiry`: if the miner has submitted their `ProveCommit` proofs.
- `CronEventProvingPeriod`: the proving period identified,
- `CronEventWorkerKeyChange`: if a miner has submitted a worker key address change.

{{<embed src="/modules/actors/builtin/cron/cron_actor.go" lang="go">}}