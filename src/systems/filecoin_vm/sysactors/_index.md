---
title: System Actors
statusIcon: 🔁
entries:
- init_actor
- cron_actor
- account_actor
- reward_actor
---

- There are two system actors required for VM processing:
  - [CronActor](#CronActor) - runs critical functions at every epoch
  - [InitActor](#InitActor) - initializes new actors
- There are two more VM level actors:
  - [AccountActor](#AccountActor) - for user accounts (non-singleton).
  - [RewardActor](#RewardActor) - for block reward and token vesting (singleton).
