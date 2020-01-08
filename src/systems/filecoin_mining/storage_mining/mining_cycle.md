---
title: Storage Mining Cycle
---
{{<label mining_cycle>}}

Block miners should constantly be performing Proofs of SpaceTime using {{<sref election_post>}}, and checking the outputted partial tickets to run {{<sref leader_election>}} and determine whether they can propose a block at each epoch. Epochs are currently set to take around X seconds, in order to account for election PoSt and network propagation around the world. The details of the mining cycle are defined here.

## Active Miner Mining Cycle

In order to mine blocks on the Filecoin blockchain a miner must be running {{<sref block_validation>}} at all times, keeping track of recent blocks received and the heaviest current chain (based on {{<sref expected_consensus>}}).

With every new tipset, the miner can use their committed power to attempt to craft a new block.

For additional details around how consensus works in Filecoin, see {{<sref expected_consensus>}}. For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner is eligible to mine a block itself, and other rules pertaining to the production of some artifacts required of valid blocks (e.g. Tickets, ElectionPoSt).

### Continuous Mining Cycle

After the chain has caught up to the current head using {{<sref chain_sync>}}, the mining process is as follows:

- The node continuously receives and transmits messages using the {{<sref message_syncer>}}
- At the same time it continuously {{<sref block_sync "receives blocks">}}
    - Each block has an associated timestamp and epoch (quantized time window in which it was crafted)
    - Blocks are validated as they come in during an epoch (provided it is their epoch, see {{<sref block "validation">}})
- At the end of a given epoch, the miner should take all the valid blocks received for this epoch and assemble them into tipsets according to {{<sref tipset "tipset validation rules">}}
- The miner then attempts to mine atop the heaviest tipset (as calculated with {{<sref chain_selection "EC's weight function">}}) using its smallest ticket to run leader election
    - The miner runs an {{<sref election_post>}} on their sectors in order to generate partial tickets
    - The miner uses these tickets in order to run {{<sref leader_election>}}
        - if successful, the miner generates a new {{<sref tickets "randomness ticket">}} for inclusion in the block
        - the miner then assembles a new block (see "block creation" below) and broadcasts it

This process is repeated until either the {{<sref election_post>}} process yields a winning ticket (in EC) and a block published or a new valid {{<sref tipset>}} comes in from the network.

At any height `H`, there are three possible situations:

- The miner is eligible to mine a block: they produce their block and propagate it. They then resume mining at the next height `H+1`.
- The miner is not eligible to mine a block but has received blocks: they form a Tipset with them and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block and has received no blocks: prompted by their clock they run leader election again, incrementing the epoch number.

Anytime a miner receives new valid blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

### Timing

{{< diagram src="./diagrams/timing/timing.png" title="Mining Cycle Timing" >}}

The mining cycle relies on receiving and producing blocks concurrently.  The sequence of these events in time is given by the timing diagram above.  The upper row represents the conceptual consumption channel consisting of successive receiving periods `Rx` during which nodes validate and select blocks as chain heads.  The lower row is the conceptual production channel made up of a period of mining `M` followed by a period of transmission `Tx`.  The lengths of the periods are not to scale.

Blocks are received and validated during `Rx` up to the end of the epoch.  At the beginning of the next epoch, the heaviest tipset is computed from the blocks received during `Rx`, used as the head to build on during `M`.  If mining is successful a block is transmitted during `Tx`.  The epoch boundaries are as shown.

In a fully synchronized network most of period `Rx` does not see any network traffic, only the period lined up with `Tx`.  In practice we expect blocks from previous epochs to propagate during the remainder of `Rx`.  We also expect differences in operator mining time to cause additional variance.

This sequence of events applies only when the node is in the `CHAIN_FOLLOW` syncing mode.  Nodes in other syncing modes do not mine blocks.

## Full Miner Lifecycle

### Step 0: Registration and Market participation

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage power actor's `CreateStorageMiner` method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done off-chain as part of storage market functions. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices).

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the {{<sref storage_market>}}. The miner will need to put up storage deal collateral for the deals they have entered into.

When they have a full sector, they should seal it. This is done by invoking the {{<sref sector_sealer>}}.

#### Owner/Worker distinction

The miner actor has two distinct 'controller' {{<sref app_address "addresses">}}. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

#### Changing Worker Addresses

Note that any change to worker keys after registration (TODO: spec how this works) must be appropriately delayed in relation to randomness lookback for SEALing data (see [this issue](https://github.com/filecoin-project/specs/issues/415)).

### Step 1: Committing Sectors

When the miner has completed their first seal, they should post it on-chain using the {{<sref storage_miner_actor>}}'s `ProveCommitSector` function. The miner will need to put up pledge collateral in proportion to the amount of storage they commit on chain. If the miner had zero committed sectors prior to this call, this begins their proving period.

You can read more about sectors {{<sref sector "here">}}.

### Step 2: Getting Power

As part of the ${{<sref election_post>}} process. A miner will be challenged by the blockchain with a surprise PoSt once per proving period on expectation.

The miner's committed sectors will turn into {{<sref storage_power>}} enabling them to run {{<sref leader_election>}} after the first successful PoSt submission on these committed sectors. If the miner is committing sectors for the first time, a SurprisePoSt will move the sectors from the miner's `ProvingSet` into the active state, otherwise an ElectionPoSt could also do this.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next successful PoSt they submit to the chain.

### Step 3: Running Elections

Once the miner has power on the network, they can begin to submit `ElectionPoSts`. To do so, the miner must run a PoSt on a subset of their sectors in every round, using the outputted partial tickets to run leader election.

If the miner finds winning tickets, they are eligible to generate a new block and earn block rewards using the {{<sref block_producer>}}.

Every successful PoSt submission will delay the next SurprisePoSt challenge the miner will receive.

In this period, the miner can still:

- commit new sectors
- be challenged with a SurprisePoSt
- declare faults

### Faults

If a miner detects {{<sref storage_faults>}} among their sectors (any sort of storage failure that would prevent them from crafting a PoSt), they should declare these faults with the `DeclareFaults()` method of the Storage Miner Actor. 

The miner will be unable to craft valid PoSts over faulty sectors, thereby reducing their chances of winning Election and SurprisePoSts. By declaring a fault, the miner will no longer be challenged on that sector, and will lose power accordingly. The miner can reverse the fault by calling `RecoverFaults()` to add the sector back into the `ProvingSet` and regain power.

A miner will no longer be able to declare faults after being challenged for a SurprisePoSt.

### Step 4: Deal/Sector Expiration

In order to stop mining, a miner must complete all of its storage deals. Once all deals in a sector have expired, the sector itself will expire thereby enabling the miner to remove the associated collateral from their account.

### Future Work

There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
