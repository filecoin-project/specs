---
menuTitle: ChainSync
title: ChainSync - synchronizing the Blockchain
entries:
- blocksync
---

# What is blockchain synchronization?

Blockchain synchronization ("sync") is a key part of a blockchain system.
It handles retrieval and propagation of blocks and transactions (messages), and
thus in charge of distributed state replication.
**This process is security critical -- problems here can be catastrophic to the
operation of a blockchain.**

# What is ChainSync?

`ChainSync` is the protocol Filecoin uses to synchronize its blockchain. It is
specific to Filecoin's choices in state representation and consensus rules,
but is general enough that it can serve other blockchains. `ChainSync` is a
group of smaller protocols, which handle different parts of the sync process.


# ChainSync Summary

At a high level, `ChainSync` does the following:

- **Part 1: Verify internal state (`INIT` state below)**
  - SHOULD verify data structures and validate local chain
  - Resource expensive verification MAY be skipped at nodes' own risk
- **Part 2: Bootstrap to the network (`BOOTSTRAP`)**
  - Step 1. Bootstrap to the network, and acquire a "secure enough" set of peers (more details below)
  - Step 2. Bootstrap to the `BlockPubsub` channels
  - Step 3. Listen and serve on Graphsync
- **Part 3: Synchronize trusted checkpoint state (`SYNC_CHECKPOINT`)**
  - Step 1. Start with a `TrustedCheckpoint` (defaults to `GenesisCheckpoint`).
  - Step 2. Get the block it points to, and that block's parents
  - Step 3. Graphsync the `StateTree`
- **Part 4: Catch up to the chain  (`CHAIN_CATCHUP`)**
  - Step 1. Maintain a set of `TargetHeads`, and select the `BestTargetHead` from it
  - Step 2. Synchronize to the latest heads observed, validating blocks towards them
  - Step 3. As validation progresses, `TargetHeads` and `BestTargetHead` will likely change
  - Step 4. Finish when node has "caught up" with `BestTargetHead` (retrieved all the state, linked to local chain, validated all the blocks, etc).
- **Part 5: Stay in sync, and participate in block propagation (`CHAIN_FOLLOW`)**
  - Step 1. If security conditions change, go back to Part 4 (`CHAIN_CATCHUP`)
  - Step 2. Receive, validate, and propagate received `Blocks`
  - Step 3. Now with greater certainty of having the best chain, finalize Tipsets, and advance chain state.

# libp2p Network Protocols

As a networking-heavy protocol, `ChainSync` makes heavy use of `libp2p`. In particular, we use two sets of protocols:

- **`libp2p.PubSub` a family of publish/subscribe protocols to propagate recent `Blocks`.**
  The concrete protocol choice impacts `ChainSync`'s effectiveness, efficiency, and security dramatically.
  For Filecoin v1.0 we will use `libp2p.Gossipsub`, a recent `libp2p` protocol that combines features and learnings
  from many excellent PubSub systems. In the future, Filecoin may use other `PubSub` protocols. *Important Note:* is entirely
  possible for Filecoin Nodes to run multiple versions simultaneously. That said, this specification *requires* that filecoin
  nodes `MUST` connect and participate in the main channel, using `libp2p.Gossipsub`.
- **`libp2p.PeerDiscovery` a family of discovery protocols, to learn about peers in the network.**
  This is especially important for security because network "Bootstrap" is a difficult problem in peer-to-peer networks.
  The set of peers we initially connect to may completely dominate our awareness of other peers, and therefore all state.
  We use a union of `PeerDiscovery` protocols as each by itself is not secure or appropriate for users' threat models.
  The union of these provides a pragmatic and effective solution.

More concretely, we use these protocols:

- **`libp2p.PeerDiscovery`**
  - **(required)** `libp2p.BootstrapList` a protocol that uses a persistent and user-configurable list of semi-trusted
    bootstrap peers. The default list includes a set of peers semi-trusted by the Filecoin Community.
  - **(required)** `libp2p.Gossipsub` a pub/sub protocol that -- as a side-effect -- disseminates peer information
  - **(optional/TODO)** `libp2p.PersistentPeerstore` a connectivity component that keeps persistent information about peers
    observed in the network throughout the lifetime of the node. This is useful because we resume and continually
    improve Bootstrap security.
  - **(optional/TODO)** `libp2p.DNSDiscovery` to learn about peers via DNS lookups to semi-trusted peer aggregators
  - **(optional/TODO)** `libp2p.HTTPDiscovery` to learn about peers via HTTP lookups to semi-trusted peer aggregators
  - **(optional)** `libp2p.KademliaDHT` a dht protocol that enables random queries across the entire network
- **`libp2p.PubSub`**
  - **(required)** `libp2p.Gossipsub` the concrete `libp2p.PubSub` protocol `ChainSync` uses.

# Subcomponents

Aside from `libp2p`, `ChainSync` uses or relies on the following components:

- Libraries:
  - `ipld` data structures, selectors, and protocols
    - `ipld.Store` local persistent storage for `chain` datastructures
    - `ipld.Selector` a way to express requests for chain data structures
    - `ipfs.GraphSync` a general-purpose `ipld` datastructure syncing protocol
- Data Structures:
  - Data structures in the `chain` package: `Block, Tipset, Chain, Checkpoint ...`
  - `chainsync.BlockCache` a temporary cache of blocks, to constrain resource expended
  - `chainsync.AncestryGraph` a datastructure to efficiently link `Blocks`, `Tipsets`, and `PartialChains`
  - `chainsync.ValidationGraph` a datastructure for efficient and secure validation of `Blocks` and `Tipsets`

## Graphsync in ChainSync

`ChainSync` is written in terms of `Graphsync`. `ChainSync` adds blockchain and filecoin-specific
synchronization functionality that is critical for Filecoin security.

### Rate Limiting Graphsync responses (SHOULD)

When running Graphsync, Filecoin nodes must respond to graphsync queries. Filecoin requires nodes
to provide critical data structures to others, otherwise the network will not function. During
ChainSync, it is in operators' interests to provide data structures critical to validating,
following, and participating in the blockchain they are on. However, this has limitations, and
some level of rate limiting is critical for maintaining security in the presence of attackers
who might issue large Graphsync requests to cause DOS.

We recommend the following:

- **Set and enforce batch size rate limits.**
  Force selectors to be shaped like: `LimitedBlockIpldSelector(blockCID, BatchSize)` for a single
  constant `BatchSize = 1000`.
  Nodes may push for this equilibrium by only providing `BatchSize` objects in responses,
  even for pulls much larger than `BatchSize`. This forces subsequent pulls to be run, re-rooted
  appropriately, and hints at other parties that they should be requesting with that `BatchSize`.
- **Force all Graphsync queries for blocks to be aligned along cacheable bounderies.**
  In conjunction with a `BatchSize`, implementations should aim to cache the results of Graphsync
  queries, so that they may propagate them to others very efficiently. Aligning on certain boundaries
  (eg specific `ChainEpoch` limits) increases the likelihood many parties in the network will request
  the same batches of content.
  Another good cacheable boundary is the entire contents of a `Block` (`BlockHeader`, `Messages`,
  `Signatures`, etc).
- **Maintain per-peer rate-limits.**
  Use bandwidth usage to decide whether to respond and how much on a per-peer basis. Libp2p already
  tracks bandwidth usage in each connection. This information can be used to impose rate limits in
  Graphsync and other Filecoin protocols.
- **Detect and react to  DOS: restrict operation.**
  The safest implementations will likely detect and react to DOS attacks. Reactions could include:
  - Smaller `Graphsync.BatchSize` limits
  - Fewer connections to other peers
  - Rate limit total Graphsync bandwidth
  - Assign Graphsync bandwidth based on a peer priority queue
  - Disconnect from and do not accept connections from unknown peers
  - Introspect Graphsync requests and filter/deny/rate limit suspicious ones


### NaiveChainsync by just Graphsyncing blocks from BlockPubsub (MAY)

It is highly probable that many nodes will find success by running a protocol
as simple as this:

```go
func (c *NaiveChainsync) Sync() {
  // for each of the blocks we see in the gossip channel
  for block := range c.BlockPubsub.NewBlocks() {
    c.IpldStore.Add(block)

    go func() {
      // fetch all the content with graphsync
      selector := ipldselector.SelectAll(block.CID)
      c.Graphsync.Pull(c.Peers, sel, c.IpldStore)

      // validate block
      c.ValidateBlock(block)
    }()
  }
}
```

However, this should be considered an honest-party optimization, unlikely to be secure or rational.
We expect that many nodes will synchronize blocks this way, or at least try it first, but we note
that such implementations will not be synchronizing securely or compatibly, and that a complete
implementation must include the full `ChainSync` protocol as described here.

## Previous BlockSync protocol

Prior versions of this spec recommended a `BlockSync` protocol. This protocol definition is
[available here](./blocksync1). Filecoin nodes are libp2p nodes, and therefore may run a variety
of other protocols, including this `BlockSync` protocol. As with anything else in Filecoin, nodes
MAY opt to use additional protocols to achieve the results.
That said, Nodes MUST implement the version of `ChainSync` as described in this spec in order to
be considered implementations of Filecoin. Test suites will assume this protocol.

# ChainSync State Machine

`ChainSync` uses the following _conceptual_ state machine. Since this is a _conceptual_ state machine,
implementations MAY deviate from implementing precisely these states, or dividing them strictly.
Implementations MAY blur the lines between the states. If so, implementations MUST ensure security
of the altered protocol.

State Machine:

{{< diagram src="chainsync_fsm.dot.svg" title="ChainSync State Machine" >}}

## Details for all states:

### Block Fetching and Validation

- `ChainSync` selects and maintains a set of the most likely heads to be correct from among those received
  via `BlockPubsub`. As more blocks are received, the set of `BestHeads` is reevaluated.
- `ChainSync` fetches `Blocks`, `Messages`, and `StateTree` through the `Graphsync` protocol.
- `ChainSync` maintains sets of `Blocks/Tipsets` in `Graphs` (see `ChainSync.id`)
- `ChainSync` gathers a list of `LatestTipsets` from `BlockPubsub`, sorted by likelihood of being the best chain (see below).
- `ChainSync` makes requests for chains of `BlockHeaders` to close gaps between  `LatestTipsets`
- `ChainSync` forms partial unvalidated chains of `BlockHeaders`, from those received via `BlockPubsub`, and those requested via `Graphsync`.
- `ChainSync` attempts to form fully connected chains of `BlockHeaders`, parting from `StateTree`, toward observed `Heads`
- `ChainSync` minimizes resource expenditures to fetch and validate blocks, to protect against DOS attack vectors.
  `ChainSync` employs **Progressive Block Validation**, validating different facets at different stages of syncing.
- `ChainSync` delays syncing `Messages` until they are needed. Much of the structure of the partial chains can
  be checked and used to make syncing decisions without fetching the `Messages`.

### Progressive Block Validation

- Blocks can be validated in progressive stages, in order to minimize resource expenditure.
- Validation computation is considerable, and a serious DOS attack vector.
- Secure implementations must carefully schedule validation and minimize the work done by pruning blocks without validating them fully.
- **Progressive Stages of Block Validation**
  - _(TODO: move this to blockchain/Block section)_
  - **BV0 - Syntactic Validation**: Validate data structure packing and ensure correct typing.
  - **BV1 - Light Consensus State Checks**: Validate `b.ChainWeight`, `b.ChainEpoch`, `b.MinerAddress`, are plausible (some ranges of bad values can be detected easily, especially if we have the state of the chain at `b.ChainEpoch - consensus.LookbackParameter`. Eg Weight and Epoch have well defined valid ranges, and `b.MinerAddress`
  must exist in the lookback state).
  - **BV2 - Signature Validation**: Verify `b.BlockSig` is correct.
  - **BV3 - Verify ElectionProof**: Verify `b.ElectionProof` is correct.
  - **BV4 - Verify Ancestry links to chain**: Verify ancestry links back to trusted blocks. If the ancestry forks off before finality, or does not connect at all, it is a bad block.
  - **BV4 - Verify MessageSigs**: Verify the signatures on messages
  - **BV5 - Verify StateTree**: Verify the application of `b.Parents.Messages()` correctly produces `b.StateTree` and `b.MessageReceipts`
- These stages can be used partially across many blocks in a candidate chain, in order to prune out clearly bad blocks long before actually doing the expensive validation work.

### Progressive Block Propagation

- In order to make Block propagation more efficient, we trade off network round trips for bandwidth usage.
- **Motivating observations:**
  - `Block` propagation is one of the most security critical points of the whole protocol.
  - Bandwidth usage during `Block` propagation is the biggest rate limiter for network scalability.
  - The time it takes for a `Block` to propagate to the whole network is a critical factor in determining a secure `BlockTime`
  - Blocks propagating through the network should take as few _sequential_ roundtrips as possible, as these roundtrips impose serious block time delays. However, interleaved roundtrips may be fine. Meaning that `block.CIDs` may be propagated on their own, without the header, then the header without the messages, then the messages.
  - `Blocks` will propagate over a `libp2p.PubSub`. `libp2p.PubSub.Messages` will most likely arrive multiple times at a node. Therefore, using only the `block.CID` here could make this very cheap in bandwidth (more expensive in round trips)
  - `Blocks` in a single epoch may include the same `Messages`, and duplicate transfers can be avoided
  - `Messages` propagate through their own `MessagePubsub`, and nodes have a significant probability of already having a large fraction of the messages in a block. Since messages are the _bulk_ of the size of a `Block`, this can present great bandwidth savings.
- **Progressive Steps of Block Propagation**
  - **IMPORTANT NOTES**:
      - these can be effectively pipelined. The `receiver` is in control of what to pull, and when. It is up them to decide when to trade-off RTTs for Bandwidth.
      - If the `sender` is propagating the block at all to `receiver`, it is in their interest to provide the full content to `receiver` when asked. Otherwise the block may not get included at all.
      - Lots of security assumptions here -- this needs to be hyper verified, in both spec and code.
  - **Step 1. (sender) `Send SignedBlock`**:
      - Block propagation begins with the `sender` propagating the `block.SignedBlock` - `Gossipsub.Send(b block.SignedBlock)`
      - This is a very light object (~200 bytes), which fits into one network packet.
      - This has the BlockCID, the MinerAddress and the Signature. This enables parties to (a) learn that a block from a particular miner is propagating, (b) validate the `MinerAddress` and `Signature`, and decide whether to invest more resources pulling the `BlockHeader`, or the `Messages`.
      - Note that this can propagate to ther rest of the network before the next steps complete.
  - **Step 2. (receiver) `Pull BlockHeader`**:
      - if `receiver` **DOES NOT** already have the `BlockHeader` for `b.BlockCID`, then:
        - `receiver` requests `BlockHeader` from `sender`:
            - `bh := Graphsync.Pull(sender, SelectCID(b.BlockCID))`
        - This is a light-ish object (<4KB).
        - This has many fields that can be validated before pulling the messages. (See **Progressive Block Validation**).
  - **Step 3. (receiver) `Pull MessageHints` (TODO: validate/decide on this)**:
      - if `receiver` **DOES NOT** already have the full block for `b.BlockCID`, then:
        - `receiver` requests `MessagePoolHints` from `sender`: ...
        - `MessagePoolHints` are TBD -- this is a compressed representation of messages that can expedite message propagation by leveraging prior `MessagePool` syncing.
        - (This is an extension and not required. can just do Step 4.)
  - **Step 4. (receiver) `Pull Messages`**:
      - if `receiver` **DOES NOT** already have the full block for `b.BlockCID`, then:
        - if `receiver` has _some_ of the messages:
          - `receiver` requests missing `Messages` and `MessageReceipts` from `sender`:
              - `Graphsync.Pull(sender, Select(m3, m10, m50, ...))`
        - if `receiver` does not have any of the messages (default safe but expensive thing to do):
          - `receiver` requests all `Messages` and `MessageReceipts` from `sender`:
              - `Graphsync.Pull(sender, SelectAll(bh.Messages), SelectAll(bh.MessageReceipts))`
        - (This is the largest amount of stuff)
  - **Step 5. (receiver) `Validate Block`**:
      - the only remaining thing to do is to complete Block Validation.


## ChainSync FSM: `INIT`

- beginning state. no network connections, not synchronizing.
- local state is loaded: internal data structures (eg chain, cache) are loaded
- `LastTrustedCheckpoint` is set the latest network-wide accepted `TrustedCheckpoint`
- **Chain State and Finality**:
  - In this state, the **chain MUST NOT advance** beyond whatever the node already has.
  - No new blocks are reported to consumers.
  - The chain state provided is whatever was loaded from prior executions (worst case is `LastTrustedCheckpoint`)
- **security conditions to transition out:**
  - local state and data structures SHOULD be verified to be correct
    - this means validating any parts of the chain or `StateTree` the node has, from `LastTrustedCheckpoint` on.
  - `LastTrustedCheckpoint` is well-known across the Filecoin Network to be a true `TrustedCheckpoint`
    - this SHOULD NOT be verified in software, it SHOULD be verified by operators
    - Note: we ALWAYS have at least one `TrustedCheckpoint`, the `GenesisCheckpoint`.
- **transitions out:**
  - once done verifying things: move to `BOOTSTRAP`

## ChainSync FSM: `BOOTSTRAP`

- `network.Bootstrap()`: establish connections to peers until we satisfy security requirement
  - for better security, use many different `libp2p.PeerDiscovery` protocols
- `BlockPubsub.Bootstrap()`: establish connections to `BlockPubsub` peers
- `Graphsync.Serve()`: set up a Graphsync service, that responds to others' queries
- **Chain State and Finality**:
  - In this state, the **chain MUST NOT advance** beyond whatever the node already has.
  - No new blocks are reported to consumers.
  - The chain state provided is whatever was loaded from prior executions (worst case is `LastTrustedCheckpoint`).
- **security conditions to transition out:**
  - `Network` connectivity MUST have reached the security level acceptable for `ChainSync`
  - `BlockPubsub` connectivity MUST have reached the security level acceptable for `ChainSync`
  - "on time" blocks MUST be arriving through `BlockPubsub`
- **transitions out:**
  - once bootstrap is deemed secure enough:
    - if node does not have the `Blocks` or `StateTree` corresponding to `LastTrustedCheckpoint`: move to `SYNC_CHECKPOINT`
    - otherwise: move to `CHAIN_CATCHUP`

## ChainSync FSM: `SYNC_CHECKPOINT`

- While in this state:
  - `ChainSync` is well-bootstrapped, but does not yet have the `Blocks` or `StateTree` for `LastTrustedCheckpoint`
  - `ChainSync` issues `Graphsync` requests to its peers randomly for the `Blocks` and `StateTree` for `LastTrustedCheckpoint`:
    - `ChainSync`'s counterparts in other peers MUST provide the state tree.
    - It is only semi-rational to do so, so `ChainSync` may have to try many peers.
    - Some of these requests MAY fail.
- **Chain State and Finality**:
  - In this state, the **chain MUST NOT advance** beyond whatever the node already has.
  - No new blocks are reported to consumers.
  - The chain state provided is the available `Blocks` and `StateTree` for `LastTrustedCheckpoint`.
- **Important Notes:**
  - `ChainSync` needs to fetch several blocks: the `Block` pointed at by `LastTrustedCheckpoint`, and its direct `Block.Parents`.
  - Nodes only need hashing to validate these `Blocks` and `StateTrees` -- no block validation or state machine computation is needed.
  - The initial value of `LastTrustedCheckpoint` is `GenesisCheckpoint`, but it MAY be a value later in Chain history.
  - `LastTrustedCheckpoint` enables efficient syncing by making the implicit economic consensus of chain history explicit.
  - By allowing fetching of the `StateTree` of `LastTrustedCheckpoint` via `Graphsync`, `ChainSync` can yield much more
    efficient syncing than comparable blockchain synchronization protocols, as syncing and validation can start there.
  - Nodes DO NOT need to validate the chain from `GenesisCheckpoint`. `LastTrustedCheckpoint` MAY be a value later in Chain history.
  - Nodes DO NOT need to but MAY sync earlier `StateTrees` than `LastTrustedCheckpoint` as well.
- Pseudocode 1: a basic version of `SYNC_CHECKPOINT`:
    ```go
    func (c *ChainSync) SyncCheckpoint() {
        while !c.HasCompleteStateTreeFor(c.LastTrustedCheckpoint) {
            selector := ipldselector.SelectAll(c.LastTrustedCheckpoint)
            c.Graphsync.Pull(c.Peers, sel, c.IpldStore)
            // Pull SHOULD NOT pull what c.IpldStore already has (check first)
            // Pull SHOULD pull from different peers simultaneously
            // Pull SHOULD be efficient (try different parts of the tree from many peers)
            // Graphsync implementations may not offer these features. These features
            // can be implemented on top of a graphsync that only pulls from a single
            // peer and does not check local store first.
        }
        c.ChainCatchup() // on to CHAIN_CATCHUP
    }
    ```
- **security conditions to transition out:**
  - `StateTree` for `LastTrustedCheckpoint` MUST be stored locally and verified (hashing is enough)
- **transitions out:**
  - once node receives and verifies complete `StateTree` for `LastTrustedCheckpoint`: move to `CHAIN_CATCHUP`

## ChainSync FSM: `CHAIN_CATCHUP`

- While in this state:
  - `ChainSync` is well-bootstrapped, and has an initial **trusted** `StateTree` to start from.
  - `ChainSync` is receiving latest `Blocks` from `BlockPubsub`
  - `ChainSync` starts fetching and validating blocks (see _Block Fetching and Validation_ above).
  - `ChainSync` has unvalidated blocks between `ChainSync.FinalityTipset` and `ChainSync.TargetHeads`
- **Chain State and Finality**:
  - In this state, the **chain MUST NOT advance** beyond whatever the node already has.
  - No new blocks are reported to consumers.
  - The chain state provided is the available `Blocks` and `StateTree` for all available epochs,
    specially the `FinalityTipset`.
- **security conditions to transition out:**
  - Gaps between `ChainSync.FinalityTipset ... ChainSync.BestTargetHead` have been closed:
    - All `Blocks` and their content MUST be fetched, stored, linked, and validated locally.
      This includes `BlockHeaders`, `Messages`, etc.
    - Bad heads have been expunged from `ChainSync.BestHeads`. Bad heads include heads that initially
      seemed good but turned out invalid, or heads that `ChainSync` has failed to connect (ie. cannot
      fetch ancestors connecting back to `ChainSync.FinalityTipset` within a reasonable amount of time).
    - All blocks between `ChainSync.FinalityTipset ... ChainSync.BestHeads` have been validated
      This means all blocks _before_ the best heads.
  - Not under a temporary network partition
- **transitions out:**
  - once gaps between `ChainSync.FinalityTipset ... ChainSync.BestHeads` are closed: move to `CHAIN_FOLLOW`

## ChainSync FSM: `CHAIN_FOLLOW`

- While in this state:
  - `ChainSync` is well-bootstrapped, and has an initial **trusted** `StateTree` to start from.
  - `ChainSync` fetches and validates blocks (see _Block Fetching and Validation_).
  - `ChainSync` is receiving and validating latest `Blocks` from `BlockPubsub`
  - `ChainSync` DOES NOT have unvalidated blocks between `ChainSync.FinalityTipset` and `ChainSync.TargetHeads`
  - `ChainSync` MUST drop back to another state if security conditions change.
- **Chain State and Finality**:
  - In this state, the **chain MUST advance** as all the blocks up to `BestTargetHead` are validated.
  - New blocks are finalized as they cross the finality threshold (`ValidG.Heads[0].ChainEpoch - FinalityLookback`)
  - New finalized blocks are reported to consumers.
  - The chain state provided includes the `Blocks` and `StateTree` for the `Finality` epoch, as well as
    candidate `Blocks` and `StateTrees` for unfinalized epochs.
- **security conditions to transition out:**
  - Temporary network partitions (see _Detecting Network Partitions_).
  - Encounter gaps of >2 epochs between Validated set and a new `ChainSync.BestTargetHead`
- **transitions out:**
  - if a temporary network partition is detected: move to `CHAIN_CATCHUP`
  - if gaps of >2 epochs form between the Validated set and `ChainSync.BestTargetHead`: move to `CHAIN_CATCHUP`
  - if node is shut down: move to `INIT`
