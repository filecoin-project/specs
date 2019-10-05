---
title: VM System Actors
entries:
- init_actor
- cron_actor
- account_actor
---

- There are two system actors required for VM processing:
  - [CronActor](#CronActor) - runs critical functions at every epoch
  - [InitActor](#InitActor) - initializes new actors
- There is one more VM level actor:
  - [AccountActor](#AccountActor) - for user accounts.
