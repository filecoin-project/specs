---
title: Blockchain Components
entries:
  - block_syncer
  - chain_manager
  - block_producer
# suppressMenu: true
---

The Filecoin blockchain is the main interface linking various actors in the Filecoin system. It ensures that the system's state is verifiably updated over time and dictates how nodes are meant to extend the network through block reception and validation and extend it through block propagation.

Its components include the:
- {{ <sref block_syncer> }} -- which receives and propagates blocks, maintaining sets of candidate chains on which the miner may mine and running syntactic validation on incoming blocks.
- {{ <sref chain_manager> }} -- which maintains a given chain's state, providing facilities to other blockchain subsystems which will query state about the latest chain in order to run, and ensuring incoming blocks are semantically validated before inclusion into the chain.
- {{ <sref block_producer> }} -- which is called in the event of a successful leader election in order to produce a new block that will extend the current heaviest chain before forwarding it to the syncer for propagation.