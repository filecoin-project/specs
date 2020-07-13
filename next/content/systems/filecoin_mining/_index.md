---
title: "Storage Mining"
description: Storage Mining System - proving storage for producing blocks
bookCollapseSection: true
weight: 6
dashboardAudit: 1
dashboardState: permanent
dashboardInterface: stable
---

# Storage Mining
---


## Status Overview
{{< dashboard-level name="Storage Mining" open="true">}}

The Storage Mining System is the part of the Filecoin Protocol that deals with storing Client's
data, producing proof artifacts that demonstrate correct storage behavior, and managing the work
involved.

Storing data and producing proofs is a complex, highly optimizable process, with lots of tunable
choices. Miners should explore the design space to arrive at something that (a) satisfies protocol
and network-wide constraints, (b) satisfies clients' requests and expectations (as expressed in
`Deals`), and \(c) gives them the most cost-effective operation. This part of the Filecoin Spec
primarily describes in detail what MUST and SHOULD happen here, and leaves ample room for
various optimizations for implementers, miners, and users to make. In some parts, we describe
algorithms that could be replaced by other, more optimized versions, but in those cases it is
important that the **protocol constraints** are satisfied. The **protocol constraints** are
spelled out in clear detail (an unclear, unmentioned constraint is a "spec error").  It is up
to implementers who deviate from the algorithms presented here to ensure their modifications
satisfy those constraints, especially those relating to protocol security.
