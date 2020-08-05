---
title: Repository
weight: 2
bookCollapseSection: true
dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Repository - Local Storage for Chain Data and Systems

The Filecoin node repository is simply an abstraction of the data which any functional Filecoin node needs to store locally in order to run correctly.

The repository is accessible to the node's systems and subsystems and acts as local storage, which, for example, is compartmentalized from the node's `FileStore`.

The repository stores the node's keys, the IPLD data structures of stateful objects as well as the node configuration settings.

<!--
{{< hint info >}}
**Code out of date**  
{{< /hint >}}

{{<embed src="repository_subsystem.id" lang="go" >}}
-->