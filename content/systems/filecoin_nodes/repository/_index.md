---
title: Repository
weight: 2
bookCollapseSection: true
dashboardWeight: 1
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Node Repository

The Filecoin node repository is simply local storage for system and chain data. It is an abstraction of the data which any functional Filecoin node needs to store locally in order to run correctly.

The repository is accessible to the node's systems and subsystems and can be compartmentalized from the node's `FileStore`.

The repository stores the node's keys, the IPLD data structures of stateful objects as well as the node configuration settings.

The Lotus implementation of the FileStore Repository can be found [here](https://github.com/filecoin-project/lotus/blob/master/node/repo/fsrepo.go).