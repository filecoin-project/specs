# Block Miner

{{<js>}}

Block miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block at each height/in each round. Rounds are currently set to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.


## The Miner Actor

After successfully calling `CreateStorageMiner`, a miner actor will be created on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.

For details on the methods on the miner actor, see its entry in the [actors spec](actors.md#storage-miner-actor).

### Owner Worker distinction

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

### Storage Mining Cycle

Storage miners must continually produce Proofs of SpaceTime over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

#### Step 0: Registration

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage market actor's [`CreateStorageMiner`](actors.md#createstorageminer) method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done through the storage markets [`AddAsk`](actors.md#addask) method. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices).

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the section on [deal](storage-market.md#deal).

When they have a full sector, they should seal it. This is done by invoking [`PoRep.Seal`](proofs.md#seal) on the sector.

#### Step 1: Commit

When the miner has completed their first seal, they should post it on-chain using [CommitSector](actors.md#commitsector). If the miner had zero committed sectors prior to this call, this begins their proving period.

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.
For example, if a miner currently PoSts for 10 sectors, and commits to 20 more sectors. The next PoSt they submit (i.e. the one they're currently proving) will be for 10 sectors again, the subsequent one will be for 30.

TODO: sectors need to be globally unique. This can be done either by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain is checked against on each submission. As the system moves towards sector aggregation, the latter option will become unworkable, so more thought needs to go into how that proof statement could work.

#### Step 2: Proving Storage (PoSt creation)

```go
func ProveStorage(sectorSize BytesAmount, sectors []commR) PoStProof {
    challengeBlockHeight := miner.ProvingPeriodEnd - POST_CHALLENGE_TIME

    // Faults to be used are the currentFaultSet for the miner.
    faults := miner.currentFaultSet
    seed := GetRandFromBlock(challengeBlockHeight)
    return GeneratePoSt(sectorSize, sectors, seed, faults)
}
```

Note: See ['Proof of Space Time'](proof-of-spacetime.md) for more details.

The proving set remains consistent during the proving period. Any sectors added in the meantime will be included in the next proving set, at the beginning of the next proving period.

#### Step 3: PoSt Submission

When the miner has completed their PoSt, they must submit it to the network by calling [SubmitPoSt](actors.md#submitpost). There are two different times that this *could* be done.

1. **Standard Submission**: A standard submission is one that makes it on-chain before the end of the proving period. The length of time it takes to compute the PoSts is set such that there is a grace period between then and the actual end of the proving period, so that the effects of network congestion on typical miner actions is minimized.
2. **Penalized Submission**: A penalized submission is one that makes it on-chain after the end of the proving period, but before the generation attack threshold. These submissions count as valid PoSt submissions, but the miner must pay a penalty for their late submission. (See '[Faults](faults.md)' for more information)
   - Note: In this case, the next PoSt should still be started at the beginning of the proving period, even if the current one is not yet complete. Miners must submit one PoSt per proving period.

Along with the PoSt submission, miners may also submit a set of sectors that they wish to remove from their proving set. This is done by selecting the sectors in the 'done' bitfield passed to `SubmitPoSt`.

### Stop Mining

In order to stop mining, a miner must complete all of its storage contracts, and remove them from their proving set during a PoSt submission. A miner may then call [`DePledge()`](actors.md#depledge) to retrieve their collateral. `DePledge` must be called twice, once to start the cooldown, and once again after the cooldown to reclaim the funds. The cooldown period is to allow clients whose files have been dropped by a miner to slash them before they get their money back and get away with it.

### Faults

Faults are described in the [faults document](faults.md).

### On Being Slashed (WIP, needs discussion)

If a miner is slashed for failing to submit their PoSt on time, they currently lose all their pledge collateral. They do not necessarily lose their storage collateral. Storage collateral is lost when a miner's clients slash them for no longer having the data. Missing a PoSt does not necessarily imply that a miner no longer has the data. There should be an additional timeout here where the miner can submit a PoSt, along with 'refilling' their pledge collateral. If a miner does this, they can continue mining, their mining power will be reinstated, and clients can be assured that their data is still there.

TODO: disambiguate the two collaterals across the entire spec

Review Discussion Note: Taking all of a miners collateral for going over the deadline for PoSt submission is really really painful, and is likely to dissuade people from even mining filecoin in the first place (If my internet going out could cause me to lose a very large amount of money, that leads to some pretty hard decisions around profitability). One potential strategy could be to only penalize miners for the amount of sectors they could have generated in that timeframe.

## Mining Blocks

Having registered as a miner, it's time to start making and checking tickets. At this point, the miner should already be running chain validation, which includes keeping track of the latest [TipSets](expected-consensus.md#tipsets) seen on the network.

For additional details around how consensus works in Filecoin, see the [expected consensus spec](expected-consensus.md). For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner should mine a block themselves, and some rules pertaining to how "Tickets" should be validated during block validation.


### Ticket Generation

For details of ticket generation, see the [expected consensus spec](expected-consensus.md#ticket-generation).

New tickets are generated using the last ticket in the ticket-chain. Generating a new ticket will take some amount of time (as imposed by the VDF in Expected Consensus).

Because of this, on expectation, as it is produced, the miner will hear about other blocks being mined on the network. By the time they have generated their new ticket, they can check whether they themselves are eligible to mine a new block (see [block creation](#block-creation)).

At any height `H`, there are three possible situations:
- The miner is eligible to mine a block: they produce their block and form a TipSet with it and other blocks received in this round (if there are any), and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block but has received blocks: they form a TipSet with them and resume mining at the next height `H+1`.
- The miner is not eligible to mine a block and has received no blocks: they run leader election again, using:
    - their losing ticket from the last leader election to produce a new ticket (the `Tickets` array in the block to be published grows with each new ticket generated).
    - the ticket `H + 1 - K` blocks back to attempt to generate an `ElectionProof`.

This process is repeated until either a winning ticket is found (and block published) or a new valid TipSet comes in from the network.

Let's illustrate this with an example.

Miner M is mining at Height H.
Heaviest tipset at H-1 is {B0}
- New Round:
    - M produces a ticket at H, from B0's ticket (the min ticket at H-1)
    - M draws the ticket from height H-K to generate an ElectionProof
    - That ElectionProof is invalid
    - M has not heard about other blocks on the network.
- New Round:
    - M produces a ticket at H + 1 using the ticket produced at H last round.
    - M draws a ticket from height H+1-K to generate an ElectionProof
    - That ElectionProof is valid
    - M generates a block B1
    - M has received blocks B2, B3 from the network with the same parents and same height.
    - M forms a tipset {B1, B2, B3}
- Finding the new min ticket/extending the ticket chain:
    - M compares the final tickets in {B1,B2,B3} (each has two tickets in their `Tickets` array). B2 has the smallest final ticket. B2 should be used to extend the ticket chain, conceptually.
- New Round:
    - M produces a new ticket at H + 2 using B2's final ticket (the min final ticket in {B1, B2, B3})
    - M draws a ticket from H+2-K to generate an ElectionProof
    - That ElectionProof is invalid
    - M has received B4 from the network, mined atop {B1,B2,B3}
- New Round with M mining atop B4

Anytime a miner receives new blocks, it should evaluate which is the heaviest TipSet it knows about and mine atop it.

### Block Creation

Scratching a winning ticket, and armed with a valid `ElectionProof`, a miner can now publish a new block!

To create a block, the eligible miner must compute a few fields:

- `Tickets` - An array containing a new ticket, and, if applicable, any intermediary tickets generated to prove appropriate delay for any failed election attempts. See [ticket generation](expected-consensus.md#ticket-generation).
- `ElectionProof` - A signature over the final ticket from the `Tickets` array proving. See [checking election results](expected-consensus.md#checking-election-results).
- `ParentWeight` - As described in [Chain Weighting](expected-consensus.md#chain-weighting).
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
  - the timestamp on the block is not in the future (with ALLOWABLE_CLOCK_DRIFT grace to account for relative asynchrony)
  - the timestamp on the block is at least BLOCK_DELAY * len(block.Tickets) higher than the latest of its parents, with BLOCK_DELAY taking on the same value as that needed to generate a valid VDF proof for a new Ticket (currently set to 30 seconds).
  - We also recommend the use of a networkTime() function to be booted on node launch and run every so frequently to call on a networked time service (e.g. ntp) and ensure relative synchrony with the rest of the network.
- `BlockSig` - A signature with the miner's private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after it propagates to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

An eligible miner can start by filling out `Parents`, `Tickets` and `ElectionProof` with values from the ticket checking process.

Next, they compute the aggregate state of their selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of the blocks' parent TipSet, sorting the parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipsets 'virtual block'. For more details on message execution and state transitions, see the [Filecoin state machine](state-machine.md) document.

Once the miner has the aggregate `ParentState`, they must apply the block reward. This is done by adding the correct block reward amount to the miner owner's account balance in the state tree. The reward will be spendable immediately in this block. See [block reward](#block-rewards) for details on how the block reward is structured. See [Notes on Block Reward Application](#notes-on-block-reward-application) for some of the nuances in applying block rewards.

Now, a set of messages is selected to put into the block. For each message, the miner subtracts `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain).

They then apply the messages state transition, and generate a receipt for it containing the total gas actually used by the execution, the executions exit code, and the return value (see [receipt](data-structures.md#message-receipt) for more details). Then, they refund the sender in the amount of `(msg.GasLimit - GasUsed) * msg.GasPrice`. In the event of a message processing error, the remaining gas is refunded to the user, and all other state changes are reverted. (Note: this is a divergence from the way things are done in Ethereum)

Each message should be applied on the resultant state of the previous message execution, unless that message execution failed, in which case all state changes caused by that message are thrown out. The final state tree after this process will be the block's `StateRoot`.

The miner merklizes the set of messages selected, and put the root in `MsgRoot`. They gather the receipts from each execution into a set, merklize them, and put that root in `ReceiptsRoot`. Finally, they set the `StateRoot` field with the resultant state.

{{% notice info %}}
Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.
{{% /notice %}}

Finally, the miner can generate a Unix Timestamp to add to their block, to show that the block generation was appropriately delayed.

The miner will wait until BLOCK_DELAY has passed since the latest block in the parent set was generated to timestamp and send out their block. We recommend using NTP or another clock synchronization protocol to ensure that the timestamp is correctly generated (lest the block be rejected). While this timestamp does not provide a hard proof that the block was delayed (we rely on the VDF in the ticket-chain to do so), it provides some softer form of block delay by ensuring that honest miners will reject undelayed blocks.

Now the block is complete, all that's left is to sign it. The miner serializes the block now (without the signature field), takes the sha256 hash of it, and signs that hash. They place the resultant signature in the `BlockSig` field.

#### Block Broadcast

An eligible miner broadcasts the completed block to the network (via [block propagation](data-propagation.md)), and assuming everything was done correctly, the network will accept it and other miners will mine on top of it, earning the miner a block reward!

### Block Rewards

Over the entire lifetime of the protocol, 1,400,000,000 FIL (`TotalIssuance`) will be given out to miners. The rate at which the funds are given out is set to halve every six years, smoothly (not a fixed jump like in Bitcoin). These funds are initially held by the network account actor, and are transferred to miners in blocks that they mine. The reward amount remains fixed for a period of 1 week (given our 30 second block time, this  is 20,160 blocks, the `AdjustmentPeriod`) and is then adjusted. Over time, the reward will eventually become zero as the fractional amount given out at each step shrinks the network account's balance to 0.

The equation for the current block reward is of the form:

```
Reward = IV * (Decay ^ (BlockHeight / 20160))
```

`IV` is the initial value, and is computed by taking:

```
IV = TotalIssuance * (1 - Decay)
```

`Decay` is computed by:

```
Decay = e^(ln(0.5) / (HalvingPeriodBlocks / AdjustmentPeriod))
```

```
// Given one block every 30 seconds, this is how many blocks are in six years
HalvingPeriodBlocks = 6 * 365 * 24 * 60 * 2 = 6,307,200 blocks
```

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

TODO: Ensure that if a miner earns a block reward while undercollateralized, then `min(blockReward, requiredCollateral-availableBalance)` is garnished (transfered to the miner actor instead of the owner).

### Notes on Block Reward Application

As mentioned above, every round, a miner checks to see if they have been selected as the leader for that particular round (see [Secret Leader Election](expected-consensus.md#secret-leader-election) in the Expected Consensus spec for more detail). Thus, it is possible that multiple miners may be selected as winners in a given round, and thus, that there will be multiple blocks with the same parents that are produced at the same block height (forming a TipSet). Each of the winning miners will apply the block reward directly to their actor's state in their state tree.

Other nodes will receive these blocks and form a TipSet out of the eligible blocks (those that have the same parents and are at the same block height). These nodes will then validate the TipSet. The full procedure for how to verify a TipSet can be found above in [Block Validation](#block-validation). To validate TipSet state, the validating node will, for each block in the TipSet, first apply the block reward value directly to the mining node's account and then apply the messages contained in the block.

Thus, each of the miners who produced a block in the TipSet will receive a block reward. There will be no lockup. These rewards can be spent immediately.

Messages in Filecoin also have an associated transaction fee (based on the gas costs of executing the message). In the case where multiple winning miners included the same message in their blocks, only the first miner will be paid this transaction fee. The first miner is the miner with the lowest ticket value (sorted lexicographically). More details on message execution can be found in the [State Machine spec](state-machine.md#execution-calling-a-method-on-an-actor).

### Open Questions

- How should receipts for tipsets 'virtual blocks' be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.


### Future Work

There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
