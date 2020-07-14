---
title: System Actors
weight: 6
bookCollapseSection: true
dashboardAudit: 0
dashboardState: wip
dashboardInterface: wip
---

# System Actors
---

- There are two system actors required for VM processing:
  - [InitActor](init_actor.md) - initializes new actors, records the network name
  - [CronActor](cron_actor.md) - runs critical functions at every epoch
- There are two more VM level actors:
  - [AccountActor](account_actor.md) - for user accounts (non-singleton).
  - [RewardActor](reward_actor.md) - for block reward and token vesting (singleton).
