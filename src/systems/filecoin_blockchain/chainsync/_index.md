---
menuTitle: ChainSync
title: ChainSync - synchronizing the Blockchain
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

# Terms and Concepts

- `LastCheckpoint` the last hard social-consensus oriented checkpoint that `ChainSync` is aware of.
  This consensus checkpoint defines the minimum finality, and a minimum of history to build on.
  `ChainSync` takes `LastCheckpoint` on faith, and builds on it, never switching away from its history.
- `TargetHeads` a list of `BlockCIDs` that represent blocks at the fringe of block production.
  These are the newest and best blocks `ChainSync` knows about. They are "target" heads because
  `ChainSync` will try to sync to them. This list is sorted by "likelihood of being the best chain".
- `BestTargetHead` the single best chain head `BlockCID` to try to sync to.
  This is the first element of `TargetHeads`

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
  - Step 1. Maintain a set of `TargetHeads` (`BlockCIDs`), and select the `BestTargetHead` from it
  - Step 2. Synchronize to the latest heads observed, validating blocks towards them (requesting intermediate points)
  - Step 3. As validation progresses, `TargetHeads` and `BestTargetHead` will likely change, as new blocks at the production fringe will arrive,
    and some target heads or paths to them may fail to validate.
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


## Previous BlockSync protocol

Prior versions of this spec recommended a `BlockSync` protocol. This protocol definition is
[available here](https://github.com/filecoin-project/specs/blob/prevspec/network-protocols.md#blocksync).
Filecoin nodes are libp2p nodes, and therefore may run a variety
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


## ChainSync FSM: `INIT`

- beginning state. no network connections, not synchronizing.
- local state is loaded: internal data structures (eg chain, cache) are loaded
- `LastTrustedCheckpoint` is set the latest network-wide accepted `TrustedCheckpoint`
- `FinalityTipset` is set to finality achieved in a prior protocol run.
  - Default: If no later `FinalityTipset` has been achieved, set `FinalityTipset` to `LastTrustedCheckpoint`
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
  - The subscription is for both peer discovery and to start selecting best heads.
    Listing on pubsub from the start keeps the node informed about potential head changes.
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
  - In this state, the **chain MUST NOT advance** beyond whatever the node already has:
    - `FinalityTipset` does not change.
    - No new blocks are reported to consumers/users of `ChainSync` yet.
    - The chain state provided is the available `Blocks` and `StateTree` for all available epochs,
      specially the `FinalityTipset`.
    - finality must not move forward here because there are serious attack vectors where a node can be forced to end up on the wrong fork if finality advances before validation is complete up to the block production fringe.
  - Validation must advance, all the way to the block production fringe:
    - Validate the whole chain, from `FinalityTipset` to `BestTargetHead`
    - The node can reach `BestTargetHead` only to find out it was invalid, then has to update `BestTargetHead` with next best one, and sync to it
      (without having advanced `FinalityTipset` yet, as otherwise we may end up on the wrong fork)
- **security conditions to transition out:**
  - Gaps between `ChainSync.FinalityTipset ... ChainSync.BestTargetHead` have been closed:
    - All `Blocks` and their content MUST be fetched, stored, linked, and validated locally.
      This includes `BlockHeaders`, `Messages`, etc.
    - Bad heads have been expunged from `ChainSync.TargetHeads`. Bad heads include heads that initially
      seemed good but turned out invalid, or heads that `ChainSync` has failed to connect (ie. cannot
      fetch ancestors connecting back to `ChainSync.FinalityTipset` within a reasonable amount of time).
    - All blocks between `ChainSync.FinalityTipset ... ChainSync.TargetHeads` have been validated
      This means all blocks _before_ the best heads.
  - Not under a temporary network partition
- **transitions out:**
  - once gaps between `ChainSync.FinalityTipset ... ChainSync.TargetHeads` are closed: move to `CHAIN_FOLLOW`
  - (Perhaps moving to `CHAIN_FOLLOW` when 1-2 blocks back in validation may be ok.
    - we dont know we have the right head until we validate it, so if other heads of similar height are right/better, we wont know till then.)

## ChainSync FSM: `CHAIN_FOLLOW`

- While in theis state:
  - `ChainSync` is well-bootstrapped, and has an initial **trusted** `StateTree` to start from.
  - `ChainSync` fetches and validates blocks (see _Block Fetching and Validation_).
  - `ChainSync` is receiving and validating latest `Blocks` from `BlockPubsub`
  - `ChainSync` DOES NOT have unvalidated blocks between `ChainSync.FinalityTipset` and `ChainSync.TargetHeads`
  - `ChainSync` MUST drop back to another state if security conditions change.
  - Keep a set of gap measures:
    - `BlockGap` is the number of remaining blocks to validate between the Validated blocks and `BestTargetHead`.
      - (ie how many epochs do we need to validate to have validated `BestTargetHead`. does not include null blocks)
    - `EpochGap` is the number of epochs between the latest validated block, and `BestTargetHead` (includes null blocks).
    - `MaxBlockGap = 2`, which means how many blocks may `ChainSync` fall behind on before switching back to `CHAIN_CATCHUP` (does not include null blocks)
    - `MaxEpochGap = 10`, which means how many epochs may `ChainSync` fall behind on before switching back to `CHAIN_CATCHUP` (includes null blocks)
- **Chain State and Finality**:
  - In this state, the **chain MUST advance** as all the blocks up to `BestTargetHead` are validated.
  - New blocks are finalized as they cross the finality threshold (`ValidG.Heads[0].ChainEpoch - FinalityLookback`)
  - New finalized blocks are reported to consumers.
  - The chain state provided includes the `Blocks` and `StateTree` for the `Finality` epoch, as well as
    candidate `Blocks` and `StateTrees` for unfinalized epochs.
- **security conditions to transition out:**
  - Temporary network partitions (see _Detecting Network Partitions_).
  - Encounter gaps of `>MaxBlockGap` or `>MaxEpochGap` between Validated set and a new `ChainSync.BestTargetHead`
- **transitions out:**
  - if a temporary network partition is detected: move to `CHAIN_CATCHUP`
  - if `BlockGap > MaxBlockGap`: move to `CHAIN_CATCHUP`
  - if `EpochGap > MaxEpochGap`: move to `CHAIN_CATCHUP`
  - if node is shut down: move to `INIT`

# Block Fetching, Validation, and Propagation

## Notes on changing `TargetHeads` while syncing

- `TargetHeads` is changing, as `ChainSync` must be aware of the best heads at any time. reorgs happen, and our first set of peers could've been bad, we keep discovering others.
  - Hello protocol is good, but it's polling. unless node is constantly polllng, wont see all the heads.
  - `BlockPubsub` gives us the realtime view into what's actually going on.
  - weight can also be close between 2+ possible chains (long-forked), and `ChainSync` must select the right one (which, we may not be able to distinguish until validating all the way)
- fetching + validation are strictly faster per round on average than blocks produced/block time (if they're not, will always fall behind), so we definitely catch up eventually (and even quickly). the last couple rounds can be close ("almost got it, almost got it, there").

## General notes on fetching Blocks

- `ChainSync` selects and maintains a set of the most likely heads to be correct from among those received
  via `BlockPubsub`. As more blocks are received, the set of `TargetHeads` is reevaluated.
- `ChainSync` fetches `Blocks`, `Messages`, and `StateTree` through the `Graphsync` protocol.
- `ChainSync` maintains sets of `Blocks/Tipsets` in `Graphs` (see `ChainSync.id`)
- `ChainSync` gathers a list of `TargetHeads` from `BlockPubsub`, sorted by likelihood of being the best chain (see below).
- `ChainSync` makes requests for chains of `BlockHeaders` to close gaps between  `TargetHeads`
- `ChainSync` forms partial unvalidated chains of `BlockHeaders`, from those received via `BlockPubsub`, and those requested via `Graphsync`.
- `ChainSync` attempts to form fully connected chains of `BlockHeaders`, parting from `StateTree`, toward observed `Heads`
- `ChainSync` minimizes resource expenditures to fetch and validate blocks, to protect against DOS attack vectors.
  `ChainSync` employs **Progressive Block Validation**, validating different facets at different stages of syncing.
- `ChainSync` delays syncing `Messages` until they are needed. Much of the structure of the partial chains can
  be checked and used to make syncing decisions without fetching the `Messages`.

## Progressive Block Validation

- Blocks can be validated in progressive stages, in order to minimize resource expenditure.
- Validation computation is considerable, and a serious DOS attack vector.
- Secure implementations must carefully schedule validation and minimize the work done by pruning blocks without validating them fully.
- **Progressive Stages of Block Validation**
  - _(TODO: move this to blockchain/Block section)_
  - **BV0 - Syntactic Validation**: Validate data structure packing and ensure correct typing.
  - **BV1 - Light Consensus State Checks**: Validate `b.ChainWeight`, `b.ChainEpoch`, `b.MinerAddress`, are plausible (some ranges of bad values can be detected easily, especially if we have the state of the chain at `b.ChainEpoch - consensus.LookbackParameter`. Eg Weight and Epoch have well defined valid ranges, and `b.MinerAddress`
  must exist in the lookback state). This requires some chain state, enough to establish plausibility levels of each of these values. A node should be able to estimate valid ranges for `b.ChainEpoch` based on the `LastTrustedCheckpoint`. `b.ChainWeight` is easy if some of the relatively recent chain is available, otherwise hard.
  - **BV2 - Signature Validation**: Verify `b.BlockSig` is correct.
  - **BV3 - Verify ElectionProof**: Verify `b.ElectionProof` is correct. This requires having state for relevant lookback parameters.
  - **BV4 - Verify Ancestry links to chain**: Verify ancestry links back to trusted blocks. If the ancestry forks off before finality, or does not connect at all, it is a bad block.
  - **BV4 - Verify MessageSigs**: Verify the signatures on messages
  - **BV5 - Verify StateTree**: Verify the application of `b.Parents.Messages()` correctly produces `b.StateTree` and `b.MessageReceipts`
- These stages can be used partially across many blocks in a candidate chain, in order to prune out clearly bad blocks long before actually doing the expensive validation work.

Notes:
- in `CHAIN_CATCHUP`, if a node is receiving/fetching hundreds/thousands of `BlockHeaders`, validating signatures can be very expensive, and can be deferred in favor of other validation. (ie lots of BlockHeaders coming in through network pipe, dont want to bound on sig verification, other checks can help dump blocks on the floor faster (BV0, BV1)
- in `CHAIN_FOLLOW`, we're not receiving thousands, we're receiving maybe a dozen or 2 dozen packets in a few seconds. We receive cid w/ Sig and addr first (ideally fits in 1 packet), and can afford to (a) check if we already have the cid (if so done, cheap), or (b) if not, check if sig is correct before fetching header (expensive computation, but checking 1 sig is way faster than checking a ton). In practice likely that which one to do is dependent on miner tradeoffs. we'll recommend something but let miners decide, because one strat or the other may be much more effective depending on their hardware, on their bandwidth limitations, or their propensity to getting DOSed

## Progressive Block Propagation (or BlockSend)

- In order to make Block propagation more efficient, we trade off network round trips for bandwidth usage.
- **Motivating observations:**
  - Block propagation is one of the most security critical points of the whole protocol.
  - Bandwidth usage during Block propagation is the biggest rate limiter for network scalability.
  - The time it takes for a Block to propagate to the whole network is a critical factor in determining a secure `BlockTime`
  - Blocks propagating through the network should take as few _sequential_ roundtrips as possible, as these roundtrips impose serious block time delays. However, interleaved roundtrips may be fine. Meaning that `block.CIDs` may be propagated on their own, without the header, then the header without the messages, then the messages.
  - `Blocks` will propagate over a `libp2p.PubSub`. `libp2p.PubSub.Messages` will most likely arrive multiple times at a node. Therefore, using only the `block.CID` here could make this very cheap in bandwidth (more expensive in round trips)
  - `Blocks` in a single epoch may include the same `Messages`, and duplicate transfers can be avoided
  - `Messages` propagate through their own `MessagePubsub`, and nodes have a significant probability of already having a large fraction of the messages in a block. Since messages are the _bulk_ of the size of a `Block`, this can present great bandwidth savings.
- **Progressive Steps of Block Propagation**
  - **IMPORTANT NOTES**:
      - these can be effectively pipelined. The `receiver` is in control of what to pull, and when. It is up them to decide when to trade-off RTTs for Bandwidth.
      - If the `sender` is propagating the block at all to `receiver`, it is in their interest to provide the full content to `receiver` when asked. Otherwise the block may not get included at all.
      - Lots of security assumptions here -- this needs to be hyper verified, in both spec and code.
      - `sender` is a filecoin node running `ChainSync`, propagating a block via Gossipsub
        (as the originator, as another peer in the network, or just a Gossipsub router).
      - `receiver` is the local filecoin node running `ChainSync`, trying to get the blocks.
      - for `receiver` to `Pull` things from `sender`, `receiver`must conntect to `sender`. Usually `sender` is sending to `receiver` because of the Gossipsub propagation rules. `receiver` could choose to `Pull` from any other node they are connected to, but it is most likely `sender` will have the needed information. They usually may be more well-connected in the network.
  - **Step 1. (sender) `Push BlockHeader`**:
      - `sender` sends `block.BlockHeader` to `receiver` via Gossipsub:
          - `bh := Gossipsub.Send(h block.BlockHeader)`
          - This is a light-ish object (<4KB).
      - `receiver` receives `bh`.
          - This has many fields that can be validated before pulling the messages. (See **Progressive Block Validation**).
          - **BV0**, **BV1**, and **BV2** validation takes place before propagating `bh` to other nodes.
          - `receiver` MAY receive many advertisements for each winning block in an epoch in quick succession. this is because (a) many want propagation as fast as possible, (b) many want to make those network advertisements as light as reasonable, (c) we want to enable `receiver` to choose who to ask it from (usually the first party to advertise it, and that's what spec will recommend), and (d) want to be able to fall back to asking others if that fails (fail = dont get it in 1s or so)
  - **Step 2. (receiver) `Pull MessageCids`**:
      - upon receiving `bh`, `receiver` checks whether it already has the full block for `bh.BlockCID`. if not:
          - `receiver` requests `bh.MessageCids` from `sender`:
              - `bm := Graphsync.Pull(sender, SelectAMTCIDs(b.Messages))`
  - **Step 3. (receiver) `Pull Messages`**:
      - if `receiver` **DOES NOT** already have the all messages for `b.BlockCID`, then:
          - if `receiver` has _some_ of the messages:
              - `receiver` requests missing `Messages` from `sender`:
                  - `Graphsync.Pull(sender, SelectAll(bm[3], bm[10], bm[50], ...))` or
                  - ```
                    for m in bm {
                      Graphsync.Pull(sender, SelectAll(m))
                    }
                    ```
          - if `receiver` does not have any of the messages (default safe but expensive thing to do):
              - `receiver` requests all `Messages` from `sender`:
                  - `Graphsync.Pull(sender, SelectAll(bh.Messages))`
          - (This is the largest amount of stuff)
  - **Step 4. (receiver) `Validate Block`**:
      - the only remaining thing to do is to complete Block Validation.


<!--
  - **Step 1. (sender) `Send SignedBlock`**: (TODO: rename to `BlockStamp`)
      - Block propagation begins with the `sender` propagating the `block.SignedBlock`
          - `Gossipsub.Send(b block.SignedBlock)`
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

-->

# Calculations


## Security Parameters

- `Peers` >= 32 -- direct connections
  - ideally `Peers` >= {64, 128}
-

## Pubsub Bandwidth

These bandwidth calculations are used to motivate choices in `ChainSync`.

If you imagine that you will receive the header once per gossipsub peer (or if lucky, half of them), and that there is EC.E_LEADERS=10 blocks per round, then we're talking the difference between:

```
16 peers, 1 pkt  -- 1 * 16 * 10 = 160 dup pkts (256KB) in <5s
16 peers, 4 pkts -- 4 * 16 * 10 = 640 dup pkts (1MB)   in <5s

32 peers, 1 pkt  -- 1 * 32 * 10 =   320 dup pkts (512KB) in <5s
32 peers, 4 pkts -- 4 * 32 * 10 = 1,280 dup pkts (2MB)   in <5s

64 peers, 1 pkt  -- 1 * 32 * 10 =   320 dup pkts (1MB) in <5s
64 peers, 4 pkts -- 4 * 32 * 10 = 1,280 dup pkts (4MB)   in <5s
```

2MB in <5s may not be worth saving-- and maybe gossipsub can be much better about supressing dups.


# Notes (TODO: move elsewhere)

## Checkpoints

- A checkpoint is the CID of a block (not a tipset list of CIDs, or StateTree)
- The reason a block is OK is that it uniquely identifies a tipset.
- using tipsets directly would make Checkpoints harder to communicate. we want to make checkpoints a single hash, as short as we can have it. They will be shared in tweets, URLs, emails, printed into newspapers, etc. Compactness, ease of copy-paste, etc matters.
- we'll make human readable lists of checkpoints, and making "lists of lists" is more annoying.
- When we have `EC.E_PARENTS > 5` or `= 10`, tipsets will get annoyingly large.
- the big quirk/weirdness with blocks it that it also must be in the chain. (if you relaxed that constraint you could end up in a weird case where a checkpoint isnt in the chain and that's weird/violates assumptions).


![](https://user-images.githubusercontent.com/138401/67015561-8c929000-f0ab-11e9-847a-ec42f23b14da.png)

## Bootstrap chain stub

- the mainnet filecoin chain will need to start with a small chain stub of blocks.
- we must include some data in different blocks.
- we do need a genesis block -- we derive randomness from the ticket there. Rather than special casing, it is easier/less complex to ensure a well-formed chain always, including at the beginning
- A lot of code expects lookbacks, especially actor code. Rather than introducing a bunch of special case logic for what happens ostensibly once in network history (special case logic which adds complexity and likelihood of problems), it is easiest to assume the chain is always at least X blocks long, and the system lookback parameters are all fine and dont need to be scaled in the beginning of network's history.

## PartialGraph

The `PartialGraph` of blocks.

> Is a graph necessarily connected, or is this just a bag of blocks, with each disconnected subgraph being reported in heads/tails?

The latter.  the partial graph is a DAG fragment-- including disconnected components.
here's a visual example, 4 example PartialGraphs, with Heads and Tails. (note they aren't tipsets)

![](https://user-images.githubusercontent.com/138401/67014349-90bdae00-f0a9-11e9-9f29-bdca6c673c4b.png)
