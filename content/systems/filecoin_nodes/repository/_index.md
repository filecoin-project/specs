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
---

The Filecoin node repository is simply an abstraction denoting that data which any functional Filecoin node needs to store locally in order to run correctly.

The repo is accessible to the node's systems and subsystems and acts as local storage compartementalized from the node's `FileStore` (for instance).

It stores the node's keys, the IPLD datastructures of stateful objects and node configs.

{{<embed src="repository_subsystem.id" lang="go" >}}
