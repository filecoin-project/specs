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
    - M draws the ticket from height H-K to generate an ElectionProof
    - That ElectionProof is invalid
    - M has not heard about other blocks on the network.
- New Round:
    - Epoch/Height is incremented to H + 1.
    - M generates a new ElectionProof with this new epoch number.
    - That ElectionProof is valid
    - M generates a block B1 using the new ElectionProof and the ticket drawn last round.
    - M has received blocks B2, B3 from the network with the same parents and same height.
    - M forms a tipset {B1, B2, B3}

Anytime a miner receives new blocks, it should evaluate what is the heaviest Tipset it knows about and mine atop it.

## Block Creation

Scratching a winning ticket, and armed with a valid `ElectionProof`, a miner can now publish a new block!

To create a block, the eligible miner must compute a few fields:

- `Ticket` - new ticket generated from that in the prior epoch (see {{<sref ticket_generation>}).
- `ElectionProof` - A specific signature over the min_ticket from `randomness_lokkback` epochs back (see {{<sref leader_election>}}).
- `ParentWeight` - The parent chain's weight (see {{<sref chain_selection>}}).
- `Parents` - the CIDs of the parent blocks.
- `ParentState` - Note that it will not end up in the newly generated block, but is necessary to compute to generate other fields. To compute this:
  - Take the `ParentState` of one of the blocks in the chosen parent set (invariant: this is the same value for all blocks in a given parent set).
  - For each block in the parent set, ordered by their tickets:
    - Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state.
    - TODO: define message conflicts in the state-machine doc, and link to it from here
- `MsgRoot` - To compute this:
  - Select a set of messages from the mempool to include in the block.
  - Separate the messages into BLS signed messages and secpk signed messages
  - For the BLS messages:
    - Strip the signatures off of the messages, and insert all the bare `Message`s for them into a sharray.
    - Aggregate all of the bls signatures into a single signature and use this to fill out the `BLSAggregate` field
  - For the secpk messages:
    - Insert each of the secpk `SignedMessage`s into a sharray
  - Create a `TxMeta` object and fill each of its fields as follows:
    - `blsMessages`: the root cid of the bls messages sharray
    - `secpkMessages`: the root cid of the secp messages sharray
  - The cid of this `TxMeta` object should be used to fill the `MsgRoot` field of the block header.
- `BLSAggregate` - The aggregated signatures of all messages in the block that used BLS signing.
- `StateRoot` - Apply each chosen message to the `ParentState` to get this.
  - Note: first apply bls messages in the order that they appear in the blsMsgs sharray, then apply secpk messages in the order that they appear in the secpkMessages sharray.
- `ReceiptsRoot` - To compute this:
  - Apply the set of messages to the parent state as described above, collecting invocation receipts as this happens.
  - Insert them into a sharray and take its root.
- `Timestamp` - A Unix Timestamp generated at block creation. We use an unsigned integer to represent a UTC timestamp (in seconds). The Timestamp in the newly created block must satisfy the following conditions:
  - the timestamp on the block corresponds to the current epoch (it is neither in the past nor in the future) as defined by the clock subsystem.
- `BlockSig` - A signature with the miner's private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after it propagates to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

An eligible miner can start by filling out `Parents`, `Tickets` and `ElectionProof` with values from the ticket checking process.

Next, they compute the aggregate state of their selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of the blocks' parent Tipset, sorting the parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipset's messages (one can think of a tipset as a 'virtual block' of sorts).

Once the miner has the aggregate `ParentState`, they must apply the block reward. This is done by adding the correct block reward amount to the miner owner's account balance in the state tree. The reward will be spendable immediately in this block.

Now, a set of messages is selected to put into the block. For each message, the miner subtracts `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain).

They then apply the messages state transition, and generate a receipt for it containing the total gas actually used by the execution, the executions exit code, and the return value . Then, they refund the sender in the amount of `(msg.GasLimit - GasUsed) * msg.GasPrice`. In the event of a message processing error, the remaining gas is refunded to the user, and all other state changes are reverted. (Note: this is a divergence from the way things are done in Ethereum)

Each message should be applied on the resultant state of the previous message execution, unless that message execution failed, in which case all state changes caused by that message are thrown out. The final state tree after this process will be the block's `StateRoot`.

The miner merklizes the set of messages selected, and put the root in `MsgRoot`. They gather the receipts from each execution into a set, merklize them, and put that root in `ReceiptsRoot`. Finally, they set the `StateRoot` field with the resultant state.

{{% notice info %}}
Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.
{{% /notice %}}

Finally, the miner can generate a Unix Timestamp to add to their block, to show that the block generation was appropriately delayed.

The miner will wait until BLOCK_DELAY has passed since the latest block in the parent set was generated to timestamp and send out their block. We recommend using NTP or another clock synchronization protocol to ensure that the timestamp is correctly generated (lest the block be rejected). While this timestamp does not provide a hard proof that the block was delayed (we rely on the VDF in the ticket-chain to do so), it provides some softer form of block delay by ensuring that honest miners will reject undelayed blocks.

Now the block is complete, all that's left is to sign it. The miner serializes the block now (without the signature field), takes the sha256 hash of it, and signs that hash. They place the resultant signature in the `BlockSig` field.

## Block Broadcast

An eligible miner broadcasts the completed block to the network and assuming everything was done correctly, the network will accept it and other miners will mine on top of it, earning the miner a block reward!

# Block Rewards

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
