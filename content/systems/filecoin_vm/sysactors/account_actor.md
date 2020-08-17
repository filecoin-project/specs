---
title: AccountActor
weight: 3
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# AccountActor
---

The `AccountActor` is responsible for user accounts. Account actors are not created by the `InitActor`, but their constructor is called by the system. Account actors are created by sending a message to a public-key style address. The address must be `BLS` or `SECP`, or otherwise there should be an exit error. The account actor is updating the state tree with the new actor address.

{{<embed src="/modules/actors/builtin/account/account_actor.go" lang="go">}}