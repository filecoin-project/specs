---
title: Block Syncer
---
In order to ensure they are always on the correct latest state of the blockchain,
a Filecoin node must continuously monitor and propagate blocks on the network.

When a node receives blocks, it must also _validate_ them.
Validation is split into two stages, syntactic and semantic.
The syntactic stage may be validated without reference to additional data, the semantic stage requires access to the chain which the block extends.

For clarity, we separate these stages into separate components:
syntactic validation is performed by the {{<sref block_syncer>}},
and once collected and validated, the blocks are forwarded
to the {{<sref chain_manager>}} which performs semantic validation and adds the blocks to the node's current view of the blockchain state.

The block syncer interfaces and functions are included here, but we expand on important details below for clarity.

goFile block_syncer

# Block reception and Syntactic Validation

- Called by: libp2p {{<sref block_sync>}}
- Calls on: {{<sref chain_manager>}} 

On reception of a new block over the appropriate libp2p channel (see gossib_sub), the _Block Syncer_'s _OnNewBlock_ method is invoked. Thereafter, the syncer must perform syntactic validation on the block to discard invalid blocks and forward the others for further validation by the {{<sref chain_manager>}}. At a high level, a syntactically valid block:

- must include a well-formed miner address
- must include at least one well-formed ticket
- must include an election proof which is a valid signature by the miner address of the final ticket
- must include at least one parent CID
- must include a positive parent weight
- must include a positive height
- must include a well-formed state root
- must include well-formed messages, and corresponding receipts CIDs
- must include a valid timestamp

## Timestamp Syntactic validation

In order to ensure that block producers release blocks as soon as possible, filecoin nodes have a cutoff time within each leader election round after which they cease to accept new blocks for this round.

Specifically, the block syncer validation process will use the {{<sref clock>}} subsystem to associate a wall clock time to the given round number. Any block coming in after the cutoff time is discarded.

In practice, appropriate parameters will not impact nodes regardless of their network connectivity, but helps clearly demarkate rounds by creating a buffer between a round's end and another's beginning.

# Block Propagation

- Called by: {{<sref block_producer>}}
- Calls on: libp2p {{<sref block_sync>}}

Blocks are propagated over the libp2p pubsub channel `/fil/blocks`. The following structure is filled out with the appropriate information, serialized (with IPLD), and sent over the wire:

```sh
type BlockMessage struct {
  header BlockHeader
  secpkMessages []&SignedMessage
  blsMessages []&Message
}
```