---
title: Filecoin Nodes
bookCollapseSection: true
weight: 1
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Filecoin Nodes
---

This section is providing the required details in order to implement any of the different Filecoin nodes. Although node types in the Lotus implementation of Filecoin are less strictly defined than in other blockchain networks, there are different properties and features that different types of nodes should implement.

In this section we also discuss issues related to storage of system files in Filecoin nodes. Note that by storage in this section we do not refer to the storage that a node is committing for mining in the network, but rather the local storage repositories that it needs to have available for keys and IPLD data among other things.

In this section we are also discussing the network interface and how nodes find and connect with each other, how they interact and propagate messages using libp2p, as well as how to set the node's clock.
