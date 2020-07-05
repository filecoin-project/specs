---
title: System Actors
statusIcon: 🔁
entries:
- init_actor
- cron_actor
- account_actor
- reward_actor
---

# System Actors
---

- There are two system actors required for VM processing:
  - [InitActor](#systems__filecoin_vm__sysactors__init_actor) - initializes new actors, records the network name
  - [CronActor](#systems__filecoin_vm__sysactors__cron_actor) - runs critical functions at every epoch
- There are two more VM level actors:
  - [AccountActor](#systems__filecoin_vm__sysactors__account_actor) - for user accounts (non-singleton).
  - [RewardActor](#systems__filecoin_vm__sysactors__reward_actor) - for block reward and token vesting (singleton).
