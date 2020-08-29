---
title: ChainSync
weight: 3
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# ChainSync

Blockchain synchronization ("sync") is a key part of a blockchain system.
It handles retrieval and propagation of blocks and transactions (messages), and
thus is in charge of distributed state replication.
As such, this process is security critical -- problems with state replication can have severe consequences to the
operation of a blockchain.

## ChainSync Overview

`ChainSync` is the protocol Filecoin uses to synchronize its blockchain. It is
specific to Filecoin's choices in state representation and consensus rules,
but is general enough that it can serve other blockchains. `ChainSync` is a
group of smaller protocols, which handle different parts of the sync process.

Chain synchronisation generally applies to three main processes:
1. Chain synchronization when a node first joins the network and needs to get to the current state before validating or extending the chain.
2. Chain synchronisation when a node has fell out of sync, e.g., due to a brief disconnection.
3. Continuous chain synchronisation to keep up with the latest messages and blocks.

There are three main protocols used to achieve synchronisation for these three processes.
- `GossipSub` is the libp2p pubsub protocol used to propagate messages and blocks. It is mainly used in the third process above when a node needs to stay in sync with new blocks being produced and propagated.
- `BlockSync` is used to synchronise specific parts of the chain, that is from and to a specific height.
- `Bitswap` is used to request and receive blocks, when a node is synchonized but GossipSub has failed to deliver some blocks to a node.
- `GraphSync` can be used to fetch parts of the blockchain as a more efficient version of `Bitswap`.

Filecoin nodes are libp2p nodes, and therefore may run a variety of other protocols. As with anything else in Filecoin, nodes MAY opt to use additional protocols to achieve the results. That said, nodes MUST implement the version of `ChainSync` as described in this spec in order to be considered implementations of Filecoin. 

## Terms and Concepts

- `LastCheckpoint` the last hard social-consensus oriented checkpoint that `ChainSync` is aware of.
  This consensus checkpoint defines the minimum finality, and a minimum of history to build on.
  `ChainSync` takes `LastCheckpoint` on faith, and builds on it, never switching away from its history.
- `TargetHeads` a list of `BlockCIDs` that represent blocks at the fringe of block production.
  These are the newest and best blocks `ChainSync` knows about. They are "target" heads because
  `ChainSync` will try to sync to them. This list is sorted by "likelihood of being the best chain". At this point this is simply realized through `ChainWeight`.
- `BestTargetHead` the single best chain head `BlockCID` to try to sync to.
  This is the first element of `TargetHeads`

## ChainSync State Machine

At a high level, `ChainSync` does the following:

- **Part 1: Verify internal state (`INIT` state below)**
  - SHOULD verify data structures and validate local chain
  - Resource expensive verification MAY be skipped at nodes' own risk
- **Part 2: Bootstrap to the network (`BOOTSTRAP`)**
  - Step 1. Bootstrap to the network, and acquire a "secure enough" set of peers (more details below)
  - Step 2. Bootstrap to the `GossipSub` channels
- **Part 3: Synchronize trusted checkpoint state (`SYNC_CHECKPOINT`)**
  - Step 1. Start with a `TrustedCheckpoint` (defaults to `GenesisCheckpoint`). The `TrustedCheckpoint` SHOULD NOT be verified in software, it SHOULD be verified by operators.
  - Step 2. Get the block it points to, and that block's parents
  - Step 3. Fetch the `StateTree`
- **Part 4: Catch up to the chain  (`CHAIN_CATCHUP`)**
  - Step 1. Maintain a set of `TargetHeads` (`BlockCIDs`), and select the `BestTargetHead` from it
  - Step 2. Synchronize to the latest heads observed, validating blocks towards them (requesting intermediate points)
  - Step 3. As validation progresses, `TargetHeads` and `BestTargetHead` will likely change, as new blocks at the production fringe will arrive,
    and some target heads or paths to them may fail to validate.
  - Step 4. Finish when node has "caught up" with `BestTargetHead` (retrieved all the state, linked to local chain, validated all the blocks, etc).
- **Part 5: Stay in sync, and participate in block propagation (`CHAIN_FOLLOW`)**
  - Step 1. If security conditions change, go back to Part 4 (`CHAIN_CATCHUP`)
  - Step 2. Receive, validate, and propagate received `Blocks`
  - Step 3. Now with greater certainty of having the best chain, finalize Tipsets, and advance chain state.


`ChainSync` uses the following _conceptual_ state machine. Since this is a _conceptual_ state machine,
implementations MAY deviate from implementing precisely these states, or dividing them strictly.
Implementations MAY blur the lines between the states. If so, implementations MUST ensure security
of the altered protocol.

![ChainSync State Machine](chainsync_fsm.dot)


## Peer Discovery

Peer discovery is a critical part of the overall architecture. Taking this wrong can have severe consequences for the operation of the protocol. This is especially important for security because network "Bootstrap" is a difficult problem in peer-to-peer networks. The set of peers a new node initially connects to may completely dominate the node's awareness of other peers, and therefore the view of the state of the network that the node has.

Peer discovery can be driven by arbitrary external means and is pushed outside the core functionality of the protocols involved in ChainSync (i.e., GossipSub, Bitswap, BlockSync). This allows for orthogonal, application-driven development and
no external dependencies for the protocol implementation. Nonetheless, the GossipSub protocol supports: i) Peer Exchange, which allows applications to bootstrap from a known set of bootstrap peers without an external peer discovery mechanism, and ii) Explicit Peering Agreements, where the application can specify a list of peers to which nodes should connect when joining.

### Peer Exchange

This process is supported through either bootstrap nodes or other normal peers. Bootstrap nodes must be maintained by system operators. They have to be stable and operate independently of protocol constructions, such as the GossipSub mesh construction, that is, bootstrap nodes do not maintain connections to the mesh.

For more details on Peer Exchange please refer to the [GossipSub specification](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub).

### Explicit Peering Agreements

With explicit peering agreements, the operators must specify a list of peers which nodes should connect to when joining. The protocol must have options available for these to be specified. For every explicit peer, the router must establish and maintain a bidirectional (reciprocal) connection.

## Progressive Block Validation

- [Blocks](block) may be validated in progressive stages, in order to minimize resource expenditure.
- Validation computation is considerable, and a serious DOS attack vector.
- Secure implementations must carefully schedule validation and minimize the work done by pruning blocks without validating them fully.
- `ChainSync` SHOULD keep a cache of unvalidated blocks (ideally sorted by likelihood of belonging to the chain), and delete unvalidated blocks when they are passed by `FinalityTipset`, or when `ChainSync` is under significant resource load.
- These stages can be used partially across many blocks in a candidate chain, in order to prune out clearly bad blocks long before actually doing the expensive validation work.

- **Progressive Stages of Block Validation**
  - **BV0 - Syntax**: Serialization, typing, value ranges.
  - **BV1 - Plausible Consensus**: Plausible miner, weight, and epoch values (e.g from chain state at `b.ChainEpoch - consensus.LookbackParameter`).
  - **BV2 - Block Signature**
  - **BV3 - Beacon entries**: Valid random beacon entries have been inserted in the block (see [beacon entry validation](storage_power_consensus#validating-beacon-entries-on-block-reception)).
  - **BV4 - ElectionProof**: A valid election proof was generated.
  - **BV5 - WinningPoSt**: Correct PoSt generated.
  - **BV6 - Chain ancestry and finality**: Verify block links back to trusted chain, not prior to finality.
  - **BV7 - Message Signatures**:
  - **BV8 - State tree**: Parent tipset message execution produces the claimed state tree root and receipts.

## Summary

When a node first joins the network it discovers peers (through the peer discovery discussed above) and joins the `/fil/blocks` and `/fil/msgs` GossipSub topics. It listens to new blocks being propagated by other nodes. It picks one block as the `BestTargetHead` and starts syncing the blockchain up to this height from the  `TrustedCheckpoint`, which by default is the `GenesisBlock` or `GenesisCheckpoint`. In order to pick the `BestTargetHead` the peer is comparing a combination of height and weight - the higher these values the higher the chances of the block being on the main chain. If there are two blocks on the same height, the peer should choose the one with the higher weight. Once the peer chooses the `BestTargetHead` it uses the BlockSync protocol to fetch the blocks and get to the current height. From that point on it is in `CHAIN_FOLLOW` mode, where it uses GossipSub to receive new blocks, or Bitswap if it hears about a block that it has not received through GossipSub.
