---
title: Storage Mining Cycle
dashboardWeight: 2
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Storage Mining Cycle
---

Block miners should constantly be performing Proofs of SpaceTime using [Election PoSt](election_post), and checking the outputted partial tickets to run [Leader Election](expected_consensus#secret-leader-election) and determine whether they can propose a block at each epoch. Epochs are currently set to take around X seconds, in order to account for election PoSt and network propagation around the world. The details of the mining cycle are defined here.

## Active Miner Mining Cycle

In order to mine blocks on the Filecoin blockchain a miner must be running [Block Validation](block) at all times, keeping track of recent blocks received and the heaviest current chain (based on [Expected Consensus](expected_consensus)).

With every new tipset, the miner can use their committed power to attempt to craft a new block.

For additional details around how consensus works in Filecoin, see [Expected Consensus](expected_consensus). For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner is eligible to mine a block itself, and other rules pertaining to the production of some artifacts required of valid blocks (e.g. Tickets, ElectionPoSt).

### Mining Cycle

After the chain has caught up to the current head using [ChainSync](chainsync). At a high-level, the mining process is as follows, (we go into more detail on epoch timing below):

- The node receives and transmits messages using the [Message Syncer](message_syncer)
- At the same time it [receives blocks](block_sync)
    - Each block has an associated timestamp and epoch (quantized time window in which it was crafted)
    - Blocks are validated as they come in [block validation](block)
- After an epoch's "cutoff", the miner should take all the valid blocks received for this epoch and assemble them into tipsets according to [Tipset validation rules](tipset) 
- The miner then attempts to mine atop the heaviest tipset (as calculated with [EC's weight function](expected_consensus#chain-selection)) using its smallest ticket to run leader election
    - The miner runs [Leader Election](expected_consensus#secret-leader-election) using the most recent [random](storage_power_consensus#beacon-entries)  output by a [drand](drand) beacon.
        - if this yields a valid `ElectionProof`, the miner generates a new [ticket](storage_power_consensus#tickets) and winning PoSt for inclusion in the block.
        - the miner then assembles a new block (see "block creation" below) and waits until this epoch's quantized timestamp to broadcast it 

This process is repeated until either the [Leader Election](expected_consensus#secret-leader-election) process yields a winning ticket (in EC) and the miner publishes a block or a new valid block comes in from the network.

At any height `H`, there are three possible situations:

- The miner is eligible to mine a block: they produce their block and propagate it. They then resume mining at the next height `H+1`.
- The miner is not eligible to mine a block but has received blocks: they form a Tipset with them and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block and has received no blocks: prompted by their clock they run leader election again, incrementing the epoch number.

Anytime a miner receives new valid blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

### Epoch Timing

{{<figure src="timing.png" title="Mining Cycle Timing">}}

The timing diagram above describes the sequence of block creation "mining", propagation and reception.

This sequence of events applies only when the node is in the `CHAIN_FOLLOW` syncing mode.  Nodes in other syncing modes do not mine blocks.

The upper row represents the conceptual consumption channel consisting of successive receiving periods `Rx` during which nodes validate incoming blocks.
The lower row is the conceptual production channel made up of a period of mining `M` followed by a period of transmission `Tx` (which lasts long enough for blocks to propagate throughout the network).  The lengths of the periods are not to scale.

The above diagram represents the important events within an epoch:

- **Epoch boundary**: change of current epoch. New blocks mined are mined in new epoch, and timestamped accordingly.
- **Epoch cutoff**: blocks from the prior epoch propagated on the network are no longer accepted. Miners can form a new tipset to mine on.

In an epoch, blocks are received and validated during `Rx` up to the prior epoch's cutoff. At the cutoff, the miner computes the heaviest tipset from the blocks received during `Rx`, and uses it as the head to build on during the next mining period `M`.  If mining is successful, the miner sets the block's timestamp to the epoch boundary and waits until the boundary to release the block. While some blocks could be submitted a bit later, blocks are all transmitted during `Tx`, the transmission period.

The timing validation rules are as follows:

- Blocks whose timestamps are not exactly on the epoch boundary are rejected.
- Blocks received with a timestamp in the future are rejected.
- Blocks received after the cutoff are rejected.
    - Note that those blocks are not invalid, just not considered for the miner's own tipset building. Tipsets received with such a block as a parent should be accepted.

In a fully synchronized network most of period `Rx` does not see any network traffic, only its beginning should. While there may be variance in operator mining time, most miners are expected to finish mining by the epoch boundary.

Let's look at an example, both use a block-time of 30s, and a cutoff at 15s.

- T = 0: start of epoch n
- T in [0, 15]: miner A receives, validates and propagates incoming blocks. Valid blocks should have timestamp 0.
- T = 15: epoch cutoff for n-1, A assembles the heaviest tipset and starts mining atop it.
- T = 25: A successfully generates a block, sets its timestamp to 30, and waits until the epoch boundary (at 30) to release it.
- T = 30: start of epoch n + 1, A releases its block for epoch n.
- T in [30, 45]: A receives and validates incoming blocks, their timestamp is 30.
- T = 45: epoch cutoff for n, A forms tipsets and starts mining atop the heaviest.
- T = 60: start of epoch n + 2.
- T in [60, 75]: A receives and validates incoming blocks
- T = 67: A successfully generates a block, sets it timestamp to 60 and releases it.
- T = 75: epoch cutoff for n+1...

Above, in epoch n, A mines fast, in epoch n+1 A mines slow. So long as the miner's block is between the epoch boundary and the cutoff, it will be accepted by other miners.

In practice miners should not be releasing blocks close to the epoch cutoff. Implementations may choose to locally randomize the exact time of the cutoff in order to prevent such behavior (while this means it may accept/reject blocks others do not, in practice this will not affect the miners submitting blocks on time).

## Full Miner Lifecycle

### Step 0: Registration and Market participation

To initially become a miner, a miner first registers a new miner actor on-chain. This is done through the storage power actor's `CreateStorageMiner` method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done off-chain as part of storage market functions. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices).

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the [Storage Market](storage_market). The miner will need to put up storage deal collateral for the deals they have entered into.

When they have a full sector, they should seal it. This is done by invoking the [Sector Sealer](sealer).

#### Owner/Worker distinction

The miner actor has two distinct 'controller' [addresses](address). One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

#### Changing Worker Addresses

Note that any change to worker keys after registration must be appropriately delayed in relation to randomness lookback for SEALing data (see [this issue](https://github.com/filecoin-project/specs/issues/415)).

### Step 1: Committing Sectors

When the miner has completed their first seal, they should post it on-chain using the [Storage Miner Actor's](storage_miner_actor) `ProveCommitSector` function. The miner will need to put up pledge collateral in proportion to the amount of storage they commit on chain. Miner will now gain power for this particular sector upon successful `ProveCommitSector`.

You can read more about sectors [here](sector) and how sector relates to power [here](storage_power_consensus#on-power).

### Step 2: Running Elections

Once the miner has power on the network, they can begin to submit `ElectionPoSts`. To do so, the miner must run a PoSt on a subset of their sectors in every round, using the outputted partial tickets to run leader election.

If the miner finds winning tickets, they are eligible to generate a new block and earn block rewards using the [Block Producer](block_producer).

Every successful PoSt submission will delay the next SurprisePoSt challenge the miner will receive.

In this period, the miner can still:

- commit new sectors
- be challenged with a SurprisePoSt
- declare faults

### Faults

If a miner detects [Storage Faults](faults#storage-faults) among their sectors (any sort of storage failure that would prevent them from crafting a PoSt), they should declare these faults with the `DeclareTemporaryFaults()` method of the Storage Miner Actor. 

The miner will be unable to craft valid PoSts over faulty sectors, thereby reducing their chances of winning Election and SurprisePoSts. By declaring a fault, the miner will no longer be challenged on that sector, and will lose power accordingly. The miner can specify how long the duration of their TemporaryFault and pay a TemporaryFaultFee.

A miner will no longer be able to declare faults after being challenged for a SurprisePoSt.

### Step 3: Deal/Sector Expiration

In order to stop mining, a miner must complete all of its storage deals. Once all deals in a sector have expired, the sector itself will expire thereby enabling the miner to remove the associated collateral from their account.

### Future Work

There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
