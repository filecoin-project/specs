---
title: Block Producer
---

{{<label block_producer>}}

# Mining Blocks

Having registered as a miner, it's time to start making and checking tickets. At this point, the miner should already be running chain validation, which includes keeping track of the latest tipsets seen on the network.

For additional details around how consensus works in Filecoin, see {{<sref expected_consensus>}}. For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner is eligible to mine a block itself, and other rules pertaining to the production of some artifacts required of valid blocks (e.g. Tickets, ElectionProofs).

## Mining Cycle

At any height `H`, there are three possible situations:

- The miner is eligible to mine a block: they produce their block and propagate it. They then resume mining at the next height `H+1`.
- The miner is not eligible to mine a block but has received blocks: they form a Tipset with them and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block and has received no blocks: prompted by their clock they run leader election again, incrementing the epoch number.

This process is repeated until either a winning ticket is found (and block published) or a new valid Tipset comes in from the network.

Let's illustrate this with an example.

Miner M is mining at Height H.
Heaviest tipset at H-1 is {B0}

- New Round:
    - M produces a ticket at H, from B0's ticket (the min ticket at H-1)
    - M draws the ticket from height H-K to generate a set of ElectionPoSt partial Tickets and uses them to run leader election
    - If M has no winning tickets
    - M has not heard about other blocks on the network.
- New Round:
    - Epoch/Height is incremented to H + 1.
    - M generates a new ElectionProof with this new epoch number.
    - If M has winning tickets
    - M generates a block B1 using the new ElectionProof and the ticket drawn last round.
    - M has received blocks B2, B3 from the network with the same parents and same height.
    - M forms a tipset {B1, B2, B3}

Anytime a miner receives new blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

## Block Creation

Scratching (a) winning ElectionPoSt ticket(s), and armed with the requisite `ElectionPoStOutput`, a miner can now publish a new block!

To create a block, the eligible miner must compute a few fields:

- `Parents` - the CIDs of the parent blocks.
- `ParentWeight` - The parent chain's weight (see {{<sref chain_selection>}}).
- `ParentState` - This is the state of the chain after all of the message from each of the `Parents` have been applied to their own `ParentState`. For more info on how to compute this, see the {{<sref vm_interpreter>}}.
- `ParentMessageReceipts` - To compute this:
  - Apply the set of messages from `Parents` to the `ParentState` as described above, collecting invocation receipts as this happens.
  - Insert them into a sharray and take its root.
- `Epoch` - The block's epoch, derived from the `Parents` epoch and the number of epochs it took to generate this block.
- `Timestamp` - A Unix Timestamp generated at block creation. We use an unsigned integer to represent a UTC timestamp (in seconds). The Timestamp in the newly created block must satisfy the following conditions:
  - the timestamp on the block corresponds to the current epoch (it is neither in the past nor in the future) as defined by the clock subsystem.
- `Ticket` - new ticket generated from that in the prior epoch (see {{<sref ticket_generation>}}).
- `Miner` - The block producer's miner actor address.
- `ElectionPoStVerifyInfo` - The byproduct of running an ElectionPoSt yielding requisite on-chain information (see {{<sref election_post>}}), namely:
  - An array of `PoStCandidate` objects, all of which include a winning partial ticket used to run leader election.
  - `PoStRandomness` used to challenge the miner's sectors and generate the partial tickets.
  - A `PoStProof` snark output to prove that the partial tickets were correctly generated.
- `Messages` - To compute this:
  - Select a set of messages from the mempool to include in the block.
  - Separate the messages into BLS signed messages and secpk signed messages
  - For the BLS messages:
    - Strip the signatures off of the messages, and insert all the bare `Message`s for them into a sharray.
    - Aggregate all of the bls signatures into a single signature and use this to fill out the `BLSAggregate` field
  - For the secpk messages:
    - Insert each of the secpk `SignedMessage`s into a sharray
  - Create a `TxMeta` object and fill each of its fields as follows:
    - `BLSMessages`: the root cid of the bls messages sharray
    - `SECPMessages`: the root cid of the secp messages sharray
  - The cid of this `TxMeta` object should be used to fill the `Messages` field of the block header.
- `BLSAggregate` - The aggregated signatures of all messages in the block that used BLS signing.
- `Signature` - A signature with the miner's worker account private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after it propagates to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

An eligible miner can start by filling out `Parents`, `Tickets` and `ElectionPoStVerifyInfo`.

Next, they compute the aggregate state of their selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of the blocks' parent Tipset, sorting the parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipset's messages (one can think of a tipset as a 'virtual block' of sorts).

They gather the receipts from each above message execution into a set, merklize them, and put that root in `ReceiptsRoot`. 

Once the miner has the aggregate `ParentState`, they must apply the block reward. This is done by adding the correct block reward amount to the miner owner's account balance in the state tree. The reward will be spendable immediately in this block.

Now, a set of messages is selected to put into the block. For each message, the miner subtracts `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain).

Finally, the miner can generate a Unix Timestamp to add to their block, to show that the block generation was appropriately delayed.

Now the block is complete, all that's left is to sign it. The miner serializes the block now (without the signature field), takes the sha256 hash of it, and signs that hash. They place the resultant signature in the `Signature` field.

## Block Broadcast

An eligible miner broadcasts the completed block to the network and assuming everything was done correctly, the network will accept it and other miners will mine on top of it, earning the miner a block reward!

Miners should output their valid block as soon as it is produced, otherwise they risk other miners receiving the block after the EPOCH_CUTOFF and not including them.

# Block Rewards

TODO: Rework this.

Over the entire lifetime of the protocol, 1,400,000,000 FIL (`TotalIssuance`) will be given out to miners. The rate at which the funds are given out is set to halve every six years, smoothly (not a fixed jump like in Bitcoin). These funds are initially held by the network account actor, and are transferred to miners in blocks that they mine. Over time, the reward will eventually become close zero as the fractional amount given out at each step shrinks the network account's balance to 0.

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

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

TODO: Ensure that if a miner earns a block reward while undercollateralized, then `min(blockReward, requiredCollateral-availableBalance)` is garnished (transfered to the miner actor instead of the owner).

## Notes on Block Reward Application

As mentioned above, every round, a miner checks to see if they have been selected as the leader for that particular round. Thus, it is possible that multiple miners may be selected as winners in a given round, and thus, that there will be multiple blocks with the same parents that are produced at the same block height (forming a Tipset). Each of the winning miners will apply the block reward directly to their actor's state in their state tree.

Other nodes will receive these blocks and form a Tipset out of the eligible blocks (those that have the same parents and are at the same block height). These nodes will then validate the Tipset. To validate Tipset state, the validating node will, for each block in the Tipset, first apply the block reward value directly to the mining node's account and then apply the messages contained in the block.

Thus, each of the miners who produced a block in the Tipset will receive a block reward. There will be no lockup. These rewards can be spent immediately.

Messages in Filecoin also have an associated transaction fee (based on the gas costs of executing the message). In the case where multiple winning miners included the same message in their blocks, only the first miner will be paid this transaction fee. The first miner is the miner with the lowest ticket value (sorted lexicographically).

# Open Questions

- How should receipts for tipsets be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.
