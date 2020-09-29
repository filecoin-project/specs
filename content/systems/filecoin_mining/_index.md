---
title: Storage Mining
bookCollapseSection: true
weight: 6
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Storage Mining

The Storage Mining System is the part of the Filecoin Protocol that deals with storing Client's
data and producing proof artifacts that demonstrate correct storage behavior.

Storage Mining is one of the most central parts of the Filecoin protocol overall, as it provides all the required consensus algorithms based on proven _storage power_ in the network. Miners are selected to mine blocks and extend the blockchain based on the storage power that they have committed to the network. Storage is added in unit of sectors and sectors are promises to the network that some storage will remain for a promised duration. In order to participate in Storage Mining, the storage miners have to: i) Add storage to the system, and ii) Prove that they maintain a copy of the data they have agreed to throughout the sector's lifetime.

Storing data and producing proofs is a complex, highly optimizable process, with lots of tunable
choices. Miners should explore the design space to arrive at something that (a) satisfies protocol
and network-wide constraints, (b) satisfies clients' requests and expectations (as expressed in
`Deals`), and (c) gives them the most cost-effective operation. This part of the Filecoin Spec
primarily describes in detail what MUST and SHOULD happen here, and leaves ample room for
various optimizations for implementers, miners, and users to make. In some parts, we describe
algorithms that could be replaced by other, more optimized versions, but in those cases it is
important that the **protocol constraints** are satisfied. The **protocol constraints** are
spelled out in clear detail.  It is up
to implementers who deviate from the algorithms presented here to ensure their modifications
satisfy those constraints, especially those relating to protocol security.
