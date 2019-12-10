---
title: Block Producer
---

{{<label block_producer>}}

# Mining Blocks

A miner registered with the storage power actor may begin generating and checking election tickets. 

In order to do so, the miner must be running chain validation, and be keeping track of the most recent
blocks received. A miner's new block will be based on parents from a previous epoch.

For additional details around how consensus works in Filecoin, see {{<sref expected_consensus>}}. For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner is eligible to mine a block itself, and other rules pertaining to the production of some artifacts required of valid blocks (e.g. Tickets, ElectionPoSt).

## Mining Cycle

At any height `H`, there are three possible situations:

- The miner is eligible to mine a block: they produce their block and propagate it. They then resume mining at the next height `H+1`.
- The miner is not eligible to mine a block but has received blocks: they form a Tipset with them and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block and has received no blocks: prompted by their clock they run leader election again, incrementing the epoch number.

This process is repeated until either a winning ticket is found (and block published) or a new valid Tipset comes in from the network.

Let's illustrate this with an example.

Miner M is mining at epoch H.
Heaviest tipset at H-1 is {B0}

- New Epoch:
    - M produces a ticket at H, from B0's ticket (the min ticket at H-1)
    - M draws the ticket from height H-K to generate a set of ElectionPoSt partial Tickets and uses them to run leader election
    - If M has no winning tickets
    - M has not heard about other blocks on the network.
- New Epoch:
    - Height is incremented to H + 1.
    - M generates a new ElectionProof with this new epoch number.
    - If M has winning tickets
    - M generates a block B1 using the new ElectionProof and the ticket drawn last epoch.
    - M has received blocks B2, B3 from the network with the same parents and same height.
    - M forms a tipset {B1, B2, B3}

Anytime a miner receives new valid blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

### Timing

{{< diagram src="../../../../diagrams/timing/timing.png" title="Mining Cycle Timing" >}}

The mining cycle relies on receiving and producing blocks concurrently.  The sequence of these events in time is given by the timing diagram above.  The upper row represents the conceptual consumption channel consisting of successive receiving periods `Rx` during which nodes validate and select blocks as chain heads.  The lower row is the conceptual production channel made up of a period of mining `M` followed by a period of transmission `Tx`.  The lengths of the periods are not to scale.

Blocks are received and validated during `Rx` up to the end of the epoch.  At the beginning of the next epoch, the heaviest tipset is computed from the blocks received during `Rx`, used as the head to build on during `M`.  If mining is successful a block is transmitted during `Tx`.  The epoch boundaries are as shown.

In a fully synchronized network most of period `Rx` does not see any network traffic, only the period lined up with `Tx`.  In practice we expect blocks from previous epochs to propagate during the remainder of `Rx`.  We also expect differences in operator mining time to cause additional variance.

This sequence of events applies only when the node is in the `CHAIN_FOLLOW` syncing mode.  Nodes in other syncing modes do not mine blocks.

## Block Creation

Producing a block for epoch `H` requires computing a tipset for epoch `H-1` (or possibly a prior epoch,
if no blocks were received for that epoch). Using the state produced by this tipset, a miner can
scratch winning ElectionPoSt ticket(s), and armed with the requisite `ElectionPoStOutput`, produce a new block.

See {{<sref vm_interpreter>}} for details of parent tipset evaluation, and {{block}} for constraints 
on valid block header values. 

To create a block, the eligible miner must compute a few fields:

- `Parents` - the CIDs of the parent tipset's blocks.
- `ParentWeight` - the parent chain's weight (see {{<sref chain_selection>}}).
- `ParentState` - the CID of the state root from the parent tipset state evaluation (see the {{<sref vm_interpreter>}}).
- `ParentMessageReceipts` - the CID of the root of an AMT containing receipts produced while computing `ParentState`.
- `Epoch` - the block's epoch, derived from the `Parents` epoch and the number of epochs it took to generate this block.
- `Timestamp` - a Unix timestamp, in seconds, generated at block creation.
- `Ticket` - a new ticket generated from that in the prior epoch (see {{<sref ticket_generation>}}).
- `Miner` - the block producer's miner actor address.
- `ElectionPoStVerifyInfo` - The byproduct of running an ElectionPoSt yielding requisite on-chain information (see {{<sref election_post>}}), namely:
  - An array of `PoStCandidate` objects, all of which include a winning partial ticket used to run leader election.
  - `PoStRandomness` used to challenge the miner's sectors and generate the partial tickets.
  - A `PoStProof` snark output to prove that the partial tickets were correctly generated.
- `Messages` - The CID of a `TxMeta` object containing message proposed for inclusion in the new block:
  - Select a set of messages from the mempool to include in the block, satisfying block size and gas limits
  - Separate the messages into BLS signed messages and secpk signed messages
  - `TxMeta.BLSMessages`: The CID of the root of an AMT comprising the bare `UnsignedMessage`s
  - `TxMeta.SECPMessages`: the CID of the root of an AMT comprising the `SignedMessage`s
- `BLSAggregate` - The aggregated signature of all messages in the block that used BLS signing.
- `Signature` - A signature with the miner's worker account private key (must also match the ticket signature) over the the block header's serialized representation (with empty signature). 

Note that the messages to be included in a block need not be evaluated in order to produce a valid block.
A miner may wish to speculatively evaluate the messages anyway in order to optimize for including
messages which will succeed in execution and pay the most gas.

The block reward is not evaluated when producing a block. It is paid when the block is included in a tipset in the following epoch.

The block's signature ensures integrity of the block after propagation, since unlike many PoW blockchains, 
a winning ticket is found independently of block generation.

## Block Broadcast

An eligible miner broadcasts the completed block to the network and, assuming everything was done correctly, 
the network will accept it and other miners will mine on top of it, earning the miner a block reward!

Miners should output their valid block as soon as it is produced, otherwise they risk other miners receiving the block after the EPOCH_CUTOFF and not including them.

# Block Rewards

TODO: Rework this.

Over the entire lifetime of the protocol, 1,400,000,000 FIL (`TotalIssuance`) will be given out to miners. The rate at which the funds are given out is set to halve every six years, smoothly (not a fixed jump like in Bitcoin). These funds are initially held by the reward actor, and are transferred to miners in blocks that they mine. Over time, the reward will eventually become close zero as the fractional amount given out at each step shrinks the network account's balance to 0.

The equation for the current block reward is of the form:

```
Reward = (IV * RemainingInNetworkActor) / TotalIssuance
```

`IV` is the initial value, and is set to:

```
IV = 153856861913558700202 attoFIL // 153.85 FIL
```

IV was derived from:
```
// Given one block every 30 seconds, this is how many blocks are in six years
HalvingPeriodBlocks = 6 * 365 * 24 * 60 * 2 = 6,307,200 blocks
λ = ln(2) / HalvingPeriodBlocks
IV = TotalIssuance * (1-e^(-λ)) // Converted to attoFIL (10e18)
```

Each of the miners who produced a block in a tipset will receive a block reward. 

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

TODO: Ensure that if a miner earns a block reward while undercollateralized, then `min(blockReward, requiredCollateral-availableBalance)` is garnished (transfered to the miner actor instead of the owner).
