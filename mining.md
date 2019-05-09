### What is the Filecoin Mining Process

An active participant in the filecoin consensus process is a storage miner and expected consensus block proposer. They are responsible for storing data for the filecoin network and also for driving the filecoin consensus process. Miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block for each round. Rounds are currently set to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.

Any block proposer must be a storage miner, but storage miners can avoid performing the block proposer tasks, however in this way, they will be losing out on block rewards and transaction fees.


## The Miner Actor

After successfully calling `CreateMiner`, a miner actor will be created on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.

For details on the methods on the miner actor, see its entry in the [actors spec](actors.md#storage-miner-actor).

### Owner Worker distinction

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

### Storage Mining Cycle

Storage miners must continually produce Proofs of SpaceTime over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

#### Step 0: Registration

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage market actor's [`CreateMiner`](actors.md#createstorageminer) method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done through the storage markets [`AddAsk`](actors.md#addask) method. A miner may create a single ask for their entire storage, or partition their storage up in some way with multiple asks (at potentially different prices). 

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the section on [deal](storage-market.md#deal).

When they have a full sector, they should seal it. This is done by invoking [`PoRep.Seal`](proofs.md#seal) on the sector.

#### Step 1: Commit

When the miner has completed their first seal, they should post it on-chain using [CommitSector](actors.md#commitsector). If the miner had zero committed sectors prior to this call, this begins their proving period.

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.

TODO: sectors need to be globally unique. This can be done either by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain is checked against on each submission. As the system moves towards sector aggregation, the latter option will become unworkable, so more thought needs to go into how that proof statement could work.

#### Step 2: Proving Storage (PoSt creation)

At the beginning of their proving period, miners collect the proving set (the set of all live sealed sectors on the chain at this point), and then call `ProveStorage`. This process will take the entire proving period to complete.

```go
func ProveStorage(sectorSize BytesAmount, sectors []commR, startTime BlockHeight) (PoSTProof, []FaultSet) {
	var proofs []Proofs
	var seeds []Seed
	var faults []FaultSet
	for t := 0; t < ProvingPeriod; t += ReseedPeriod {
		seeds = append(seeds, GetSeedFromBlock(startTime+t))
		proof, faultset := GenPost(sectors, seeds[t], vdfParams)
		proofs = append(proofs, proof)
		faults = append(faults, faultset)
	}
	return GenPostSnark(sectorSize, sectors, seeds, proofs), faults
}
```

Note: See ['Proof of Space Time'](proofs.md#proof-of-space-time) for more details.

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

Having registered as a miner, it's time to start making and checking tickets. At this point, the miner should already be running chain validation, which includes keeping track of the latest [TipSets](./expected-consensus.md#tipsets) seen on the network.

For additional details around how consensus works in Filecoin, see the [expected consensus spec](./expected-consensus.md). For the purposes of this section, there is a consensus protocol (Expected Consensus) that guarantees a fair process for determining what blocks have been generated in a round, whether a miner should mine a block themselves, and some rules pertaining to how "Tickets" should be validated during block validation.

### Receiving Blocks

When receiving blocks from the network (via [block propagation](data-propagation.md)), a miner must do the following:

1. Check their validity (see [below](#block-validation)).
2. Assemble a TipSet with all valid blocks with common parents and the same number of tickets in their `Tickets` array.

A miner may sometimes receive blocks belonging to different TipSets (i.e. whose parents are not the same). In that case, they must choose which TipSet to mine on.

Chain selection is a crucial component of how the Filecoin blockchain works. Every chain has an associated weight accounting for the number of blocks mined on it and so the power (storage) they track. It is always preferable to mine atop a heavier TipSet rather than a lighter one. While a miner may be foregoing block rewards earned in the past, this lighter chain is likely to be abandoned by other miners forfeiting any block reward earned as miners converge on a final chain. For more on this, see [chain selection](./expected-consensus.md#chain-selection) in the Expected Consensus spec.

### Block Validation

The block structure and serialization is detailed in [the datastructures spec - block](data-structures.md#block). Check there for details on fields and types.

In order to validate a block coming in from the network at round `N` was well mined a miner must do the following:

```go
func VerifyBlock(blk Block) {
    // 1. Verify Signature
    pubk := GetPublicKey(blk.Miner)
    if !ValidateSignature(blk.Signature, pubk, blk) {
        Fatal("invalid block signature")
    }
    
    // 2. Verify Timestamp
    // first check that it is not in the future
    if blk.GetTime() > time.Now() {
        Fatal("block was generated too far in the future")
    }
    // next check that it is appropriately delayed from its parents
    if blk.GetTime() <= blk.minParentTime() + BLOCK_DELAY {
        Fatal("block was generated too soon")
    }

    // 3. Verify ParentWeight
    if blk.ParentWeight != ComputeWeight(blk.Parents) {
        Fatal("invalid parent weight")
    }

    // 4. Verify Tickets
    if !VerifyTickets(blk) {
        Fatal("tickets were invalid")
    }
    
    // 5. Verify ElectionProof
    randomnessLookbackTipset := RandomnessLookback(blk)
    lookbackTicket := minTicket(randomnessLookbackTipset)
    challenge := sha256.Sum(lookbackTicket)
    
    if !ValidateSignature(blk.ElectionProof, pubk, challenge) {
        Fatal("election proof was not a valid signature of the last ticket")
    }
    
    powerLookbackTipset := PowerLookback(blk)
    minerPower := GetPower(powerLookbackTipset.state, blk.Miner)
    totalPower := GetTotalPower(powerLookbackTipset.state)
    if !IsProofAWinner(blk.ElectionProof, minerPower, totalPower) {
        Fatal("election proof was not a winner")
    }
        
    // 6. Verify StateRoot
    state := GetParentState(blk.Parents)
    for i, msg := range blk.Messages {
        receipt := ApplyMessage(state, msg)
        if receipt != blk.MessageReceipts[i] {
            Fatal("message receipt mismatch")
        }
    }
    if state.Cid() != blk.StateRoot {
        Fatal("state roots mismatch")
    }
}
```

Ticket validation is detailed as follows:

```go
func RandomnessLookback(blk Block) TipSet {
	return chain.GetAncestorTipset(blk, K)
}

func PowerLookback(blk Block) TipSet {
	return chain.GetAncestorTipset(blk, L)
}

func IsProofAWinner(p ElectionProof, minersPower, totalPower Integer) bool {
	return Integer.FromBytes(sha256.Sum(p))*totalPower < minersPower*2^32
}

func VerifyTickets(b Block) error {
	// Start with the `Tickets` array
	// get the smallest ticket from the blocks parent tipset
	parTicket := GetSmallestTicket(b.Parents)

	// Verify each ticket in the chain of tickets. There will be one ticket
	// plus one ticket for each failed election attempt.
	for _, ticket := range b.Tickets {
		challenge := sha256.Sum(parTicket.Signature)

		// Check VDF
		if !VerifyVDF(ticket.VDFProof, ticket.VDFResult, challenge) {
			return "VDF was not run properly"
		}

		// Check VRF
		pubk := getPublicKeyForMiner(b.Miner)
		if !VerifySignature(ticket.Signature, pubk, ticket.VDFResult) {
			return "Ticket was not a valid signature over the parent ticket"
		}
		// in case mining this block generated multiple tickets
		parTicket = ticket
	}

	return nil
}
```


If all of this lines up, the block is valid. The miner repeats this for all blocks in a TipSet, and for all TipSets formed from incoming blocks.

Once they've ensured all blocks in the heaviest TipSet received were properly mined, they can mine on top of it. If they weren't, the miner may need to ensure the next heaviest `Tipset` was properly mined. This might mean the same `Tipset` with invalid blocks removed, or an altogether different one.

If no valid blocks are received, a miner may mine atop the same `TipSet` running leader election again using the next ticket in the ticket chain, and also generating a [new ticket](./expected-consensus.md#losing-tickets) in the process (see the [expected consensus spec](./expected-consensus.md) for more). 

### Ticket Generation

For details of ticket generation, see the [expected consensus spec](./expected-consensus.md#ticket-generation).

Ticket generation is the twin process of leader election (i.e. generating `ElectionProof`s). Every ticket scratch (i.e. round of leader election) has the miner generate a new ticket to include in the `Tickets` array of their block.

New tickets are generated using the smallest ticket from the parent TipSet at height `H` and added to the `Tickets` array. Generating a new ticket will take some amount of time (as imposed by the VDF in Expected Consensus).

Because of this, on expectation, as it is produced, the miner will hear about other blocks being mined on the network. By the time they have generated their new ticket, they can check whether they themselves are eligible to mine a new block.

If the lookback ticket yields a valid `ElectionProof`, the miner publishes their block (see [block generation](#block-creation)) including the new ticket and earns a block reward. They then assemble a new TipSet using any valid blocks they heard about while generating the ticket  (likely of height `H+1`) and mine atop the smallest ticket in that new TipSet.

### Scratching a losing ticket

If a miner fails to generate a valid `ElectionProof` using their lookback ticket, they may not yet publish a new block at height `H+1`. The miner will likely hear about other blocks being mined on the network at height `H+1` and can thus assemble a new TipSet to mine off of with these blocks (as above).

If the miner hears of no new blocks, they must instead draw a new ticket to scratch in order to try leader election again (as other miners will). In order to do so, they must generate a new ticket once more.

Now, rather than generating this new ticket from the smallest ticket from the parent TipSet (as above), the miner will instead use their ticket from the last round, now in the `Tickets` array.

This process is repeated until either a winning ticket is found (and block published) or a new valid TipSet comes in from the network. If a new TipSet comes in from the network, and it is heavier chain than the miner's own, they should abandon their process to mine atop this new block. Due to the way chain selection works in filecoin, a chain with more blocks will be preferred (see the [Expected Consensus spec](./expected-consensus.md#chain-selection) for more details).

The `Tickets` array in the block to be published grows with each round (and a new ticket generated).

### Block Creation

Scratching a winning ticket, and armed with a valid `ElectionProof`, a miner can now publish a new block!

To create a block, the eligible miner must compute a few fields:

- `Tickets` - An array containing a new ticket, and, if applicable, any intermediary tickets generated to prove appropriate delay for any failed election attempts. See [ticket generation](./expected-consensus.md#ticket-generation).
- `ElectionProof` - A signature over the final ticket from the `Tickets` array proving. See [ticket generation](./expected-consensus.md#ticket-generation).
- `Timestamp` - A Unix Timestamp generated at block creation. We use an unsigned integer to represent a UTC timestamp. The Timestamp in the newly created block must satisfy the following conditions:
  - the timestamp on the block is not in the future
  - the timestamp on the block is at least BLOCK_DELAY higher than the latest of its parents, with BLOCK_DELAY taking on the same value as that needed to generate a valid VDF proof for a new Ticket (currently set to 30 seconds).
- `ParentWeight` - As described in [Chain Weighting](./expected-consensus.md#chain-weighting).
- `ParentState` - Note that it will not end up in the newly generated block, but is necessary to compute to generate other fields. To compute this:
  - Take the `ParentState` of one of the blocks in the chosen parent set (invariant: this is the same value for all blocks in a given parent set).
  - For each block in the parent set, ordered by their tickets:
    - Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state. 
    - TODO: define message conflicts in the state-machine doc, and link to it from here
- `MsgRoot` - To compute this:
  - Select a set of messages from the mempool to include in the block.
  - Insert them into a Merkle Tree and take its root.
- `StateRoot` - Apply each chosen message to the `ParentState` to get this.
- `ReceiptsRoot` - To compute this:
  - Apply the set of messages selected above to the parent state, collecting invocation receipts as this happens.
  - Insert them into a Merkle Tree and take its root.
- `BlockSig` - A signature with the miner's private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after it propagates to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

An eligible miner can start by filling out `Parents`, `Tickets` and `ElectionProof` with values from the ticket checking process.

Next, they compute the aggregate state of their selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of the blocks' parent TipSet, sorting the parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipsets 'virtual block'. For more details on message execution and state transitions, see the [Filecoin state machine](state-machine.md) document.

Once the miner has the aggregate `ParentState`, they must apply the mining reward. This is done by adding the correct amount to the miner owner's account balance in the state tree. See [block reward](#block-rewards) for details.

Now, a set of messages is selected to put into the block. For each message, the miner subtracts `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain).

They then apply the messages state transition, and generate a receipt for it containing the total gas actually used by the execution, the executions exit code, and the return value (see [receipt](data-structures.md#message-receipt) for more details). Then, they refund the sender in the amount of `(msg.GasLimit - GasUsed) * msg.GasPrice`. In the event of a message processing error, the remaining gas is refunded to the user, and all other state changes are reverted. (Note: this is a divergence from the way things are done in Ethereum)

Each message should be applied on the resultant state of the previous message execution, unless that message execution failed, in which case all state changes caused by that message are thrown out. The final state tree after this process will be the block's `StateRoot`.

The miner merklizes the set of messages selected, and put the root in `MsgRoot`. They gather the receipts from each execution into a set, merklize them, and put that root in `ReceiptsRoot`. Finally, they set the `StateRoot` field with the resultant state.

{% hint style='info' %}
Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.
{% endhint %}

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
HalvingPeriodBlocks = 6 * 365 * 24 * 60 * 2
```

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

### Open Questions

- How should receipts for tipsets 'virtual blocks' be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.


### Future Work
There are many ideas for improving upon the storage miner, here are ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, and will not be implemented in the initial release of Filecoin, but could provide interesting capabilities down the road.
