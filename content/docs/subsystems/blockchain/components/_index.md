---
title: Blockchain Components
entries:
  - block_receiver
  - block_propagator
  - chain_manager
  - block_producer
suppressMenu: true
---

In order to ensure they are always on the correct latest state of the blockchain,
a Filecoin node must continuously monitor and propagate blocks on the network.

When a node receives blocks, it must also _validate_ them.
Validation is split into two stages, syntactic and semantic.
The syntactic stage may be validated without reference to additional data,
while the semantic stage requires access to the chain which the block extends.
For clarity, we separate these stages into separate components:
syntactic validation is performed by the {{<sref block_receiver>}},
and once collected and validated, the blocks are forwarded
to the {{<sref chain_manager>}},
which performs semantic validation and adds the blocks to the node's
current view of the blockchain state.
