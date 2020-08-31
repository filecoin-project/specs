---
title: Virtual Machine
description: VM - Virtual Machine
bookCollapseSection: true
weight: 3
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# VM - Virtual Machine

An Actor in the Filecoin Blockchain is the equivalent of the smart contract in the Ethereum Virtual Machine. Actors carry the logic needed in order to submit transactions, proofs and blocks, among other things, to the Filecoin blockchain. Every actor is identified by a unique address.

The Filecoin Virtual Machine (VM) is the system component that is in charge of execution of all actors code. Execution of actors on the Filecoin VM (i.e., on-chain executions) incur a gas cost.

Any operation applied (i.e., executed) on the Filecoin VM produces an output in the form of a _State Tree_ (discussed below). The latest _State Tree_ is the current source of truth in the Filecoin Blockchain. The _State Tree_ is identified by a CID, which is stored in the IPLD store.