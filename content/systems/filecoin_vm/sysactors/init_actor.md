---
title: InitActor
weight: 1
dashboardWeight: 2
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# InitActor
---

The `InitActor` has the power to create new actors, e.g., those that enter the system. It maintains a table resolving a public key and temporary actor addresses to their canonical ID-addresses. Invalid CIDs should not get committed to the state tree.

The actor should be able to construct a canonical ID address for the actor that will persist even in case of a chain re-organisation.

{{<embed src="/modules/actors/builtin/init/init_actor.go" lang="go">}}