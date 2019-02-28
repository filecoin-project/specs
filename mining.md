> TODO @nicola: there is a lot of redundancy in this doc, I would divide into three categories:
> - Storage Mining Cycle (have a single more compact organised set of sections)
> - Block Mining Cycle (mostly as is)

# The Filecoin Mining Process
### What is the Filecoin Mining Process

An active participant in the filecoin consensus process is a storage miner and expected consensus block proposer. They are responsible for storing data for the filecoin network and also for driving the filecoin consensus process. Miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block for each round. We currently set rounds to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.

Any block proposer must be a storage miner, but storage miners can avoid performing the block proposer tasks, however in this way, they will be losing out on block rewards and transaction fees.


## The Miner Actor

After successfully calling `CreateMiner`, a miner actor will be created for you on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.

For details on the methods on the miner actor, see its entry in the [actors spec](actors.md#storage-miner-actor).

### Owner Worker distinction

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key.

### Storage Mining Cycle

Storage miners must continually produce Proofs of SpaceTime over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

#### Step 0: Registration

To initially become a miner, a miner first register a new miner actor on-chain. This is done through the storage market actor's [`CreateMiner`](actors.md#createminer) method. The call will then create a new miner actor instance and return its address.

The next step is to place one or more storage market asks on the market. This is done through the storage markets [`AddAsk`](actors.md#addask) method. A miner may create a single ask for their entire storage, or partition your storage up in some way with multiple asks (at potentially different prices). 

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the section on [deal flow](storage-market.md#deal-flow).

When they have a full sector, they should seal it. This is done by invoking [`PoRep.Seal`](proofs.md#seal) on the sector.

#### Step 1: Commit

When the miner has completed their first seal, they should post it on-chain using [CommitSector](actors.md#commit-sector). If the miner had zero committed sectors prior to this call, this begins their proving period.

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.

TODO: sectors need to be globally unique. We can either do this by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain that we check against on each submission. I think that when we go towards sector aggregation, the latter option will become pretty much impossible, so we need to think about how that proof statement could work.

#### Step 2: Proving Storage (PoSt creation)

At the beginning of their proving period, miners collect the proving set (the set of all live sealed sectors on the chain at this point), and then call `ProveStorage`. This process will take the entire proving period to complete.

```go
func ProveStorage(sectors []commR, startTime BlockHeight) (PoSTProof, []Fault) {
    proofs []Proofs
    seeds []Seed
    faults []Fault
    for t := 0; t < ProvingPeriod; t += ReseedPeriod {
        seeds = append(seeds, GetSeedFromBlock(startTime + t))
        proof, fault := GenPost(sectors, seeds[t], vdfParams)
        proofs = append(proofs, proof)
        faults = append(faults, fault)
    }
    return GenPostSnark(sectors, seeds, proofs), faults
}
```

Note: See ['Proof of Space Time'](proofs.md#proof-of-space-time) for more details.

The proving set remains consistent during the proving period. Any sectors added in the meantime will be included in the next proving set, at the beginning of the next proving period.

#### Step 3: PoSt Submission

When the miner has completed their PoSt, they must submit it to the network by calling [SubmitPoSt](actors.md#submit-post). There are two different times that this *could* be done.

1. **Standard Submission**: A standard submission is one that makes it on-chain before the end of the proving period. The length of time it takes to compute the PoSts is set such that there is a grace period between then and the actual end of the proving period, so that the effects of network congestion on typical miner actions is minimized.
2. **Penalized Submission**: A penalized submission is one that makes it on-chain after the end of the proving period, but before the generation attack threshold. These submissions count as valid PoSt submissions, but the miner must pay a penalty for their late submission. (See '[Faults](../faults.md)' for more information)
   - Note: In this case, the next PoSt should still be started at the beginning of the proving period, even if the current one is not yet complete. Miners must submit one PoSt per proving period.

Along with the PoSt submission, miners may also subit a set of sectors that they wish to remove from their proving set. This is done by selecting the sectors in the 'done' bitfield passed to `SubmitPoSt`.

### Stop Mining

In order to stop mining, a miner must complete all of its storage contracts, and remove them from their proving set during a PoSt submission. A miner may then call [`DePledge()`](actors.md#depledge) to retrieve their collateral. `DePledge` must be called twice, once to start the cooldown, and once again after the cooldown to reclaim the funds. The cooldown period is to allow clients whose files have been dropped by a miner to slash them before they get their money back and get away with it.

### Faults

Faults are described in the [faults document](../faults.md).

### On Being Slashed (WIP, needs discussion)

If a miner is slashed for failing to submit their PoSt on time, they currently lose all their pledge collateral. They do not necessarily lose their storage collateral. Storage collateral is lost when a miners clients slash them for no longer having the data. Missing a PoSt does not necessarily imply that a miner no longer has the data. There should be an additional timeout here where the miner can submit a PoSt, along with 'refilling' their pledge collateral. If a miner does this, they can continue mining, their mining power will be reinstated, and clients can be assured that their data is still there.

TODO: disambiguate the two collaterals across the entire spec

Review Discussion Note: Taking all of a miners collateral for going over the deadline for PoSt submission is really really painful, and is likely to dissuade people from even mining filecoin in the first place (If my internet going out could cause me to lose a very large amount of money, that leads to some pretty hard decisions around profitability). One potential strategy could be to only penalize miners for the amount of sectors they could have generated in that timeframe. 

## Mining Blocks

Now that you are a real life filecoin miner, it's time to start making and checking tickets. At this point, you should already be running chain validation, which includes keeping track of the latest [TipSets](./expected-consensus.md#tipsets) that you've seen on the network.

For additional details around how consensus works in Filecoin, see the [expected consensus spec](./expected-consensus.md). For the purposes of this section, we have a consensus protocol (Expected Consensus) that guarantees us a fair process for determining what blocks have been generated in a round and whether a miner should mine a block themselves, and some rules pertaining to how "Tickets" should be validated during block validation.

### Receiving Blocks

When receiving blocks from the network (via [block propagation](data-propogation.md)), you must do the following:

1. Check their validity (see [below](#block-validation)).
2. Assemble a TipSet with all valid blocks with common parents.

You may sometimes receive blocks belonging to different TipSets (i.e. whose parents are not the same). In that case, you must choose which TipSet to mine on.

Chain selection is a crucial component of how the Filecoin blockchain works. Every chain has an associated weight accounting for the number of blocks mined on it and the power (storage) they track. It is always preferable to mine atop a heavier `TipSet` rather than a lighter one. While you may be foregoing block rewards earned in the past, this lighter chain is likely to be abandoned by other miners forfeiting any block reward earned. For more on this, see [chain selection](./expected-consensus.md#chain-selection) in the Expected Consensus spec.

### Block Validation

The block structure and serialization is detailed in  [the datastructures spec](data-structures.md#block). Check there for details on fields and types.

In order to validate a block coming in from the network at height `N` was well mined you must do the following:

TODO: 'ticket height' -> 'round number'

1. Validate that `BlockSig` validates with the miners public.
2. Validate `ParentWeight`
   - Each of the blocks in the block's `ParentTipset` must have the same `ParentWeight`.
   - The new block's ParentWeight must have been properly calculated using the [chain weighting function](./expected-consensus#chain-weighting).
3. Validate `StateRoot`
   - In order to do this, you must ensure that applying the block's messages to the `ParentState` (not included in block header) appropriately yields the `StateRoot` and receipts in the `ReceiptRoot`.
4. You must validate that the tickets in the `Tickets` array are valid, thereby proving appropriate delay for this block creation.
   - Ensure that the new ticket was generated from the smallest ticket in the block's parent tipset (at round `N-1`, inclusive of null tickets).
   - Recompute the new ticket, using the miner's public key, ensuring it was computed appropriately.
5. You must validate that the `ElectionProof` was correctly generated by the miner, and that they are eligible to mine.
   - Ensure that the proof was generated using the smallest ticket at round `N-K` (the lookback ticket).
   - Validate the proof using the miner's public key, to check that the `ElectionProof` is a valid signature over that lookback ticket.
   - Verify that the proof is indeed smaller than the miner's power fraction as reported in the power table at round `N-L`.
6. In a case where the block contains multiple values in the `Tickets` array
   - Ensure that all tickets in both arrays are signed by the same key.
   - Ensure that each ticket was used to generate the next, starting from the smallest ticket in the Parent tipset (at `N-1`).
   

We detail ticket validation as follows:

```go
func RandomnessLookback(blk Block) TipSet {
    return chain.GetAncestorTipset(blk, L)
}

func PowerLookback(blk Block) TipSet {
    return chain.GetAncestorTipset(blk, K)
}

func IsTicketAWinner(t Ticket, minersPower, totalPower Integer) bool {
    return ToFloat(sha256.Sum(ticket)) * totalPower < minersPower
}

func VerifyTicket(b Block) error {
    // 1. start with the `Tickets` array
    // get the smallest ticket from the blocks parent tipset
    parTicket := selectSmallestTicket(b.Parents)
    
    // Verify each ticket in the chain of tickets. There will be one ticket
    // plus one ticket for each null block. Only the final ticket must be a
    // 'winning' ticket.
    for _, ticket := range b.Tickets {
    	challenge := sha256.Sum(parTicket)
    	
        // Check VDF
        if !Verify(b.VDFProof, challenge) {
            return "VDF was not run properly"
        }
        
        // Check VRF   
    	pubk := getPublicKeyForMiner(b.Miner)
    	if !pubk.VerifySignature(ticket, parTicket) {
        	return "Ticket was not a valid signature over the parent ticket"
    	}
        // in case this block was mined atop null blocks
        parTicket = ticket
    }
    
    // 2. Check leader election
    // get the smallest ticket from the lookback tipset
    lookbackTicket := selectSmallestTicket(RandomnessLookback(b))
    challenge := sha256.Sum(lookbackTicket)
	
    // Check VRF
    pubk := getPublicKeyForMiner(b.Miner)
    if !pubk.VerifySignature(b.ElectionProof, lookbackTicket) {
		return "Ticket was not a valid signature over the lookback ticket"
	}
    
    state := getStateTree(PowerLookback(b))
    minersPower := state.getPowerForMiner(b.Miner)
    totalPower := state.getTotalPower()
    if !IsTicketAWinner(b.ElectionProof, minersPower, totalPower) {
       return "Ticket was not a winning ticket"
    }
    
    // Winner!
    return nil
}
```


If all of this lines up, the block is valid. Repeat for all blocks in a tipset.

Once you've ensured all blocks in the `Tipset` received were properly mined, you can mine on top of it. If it wasn't, ensure the next heaviest `Tipset` was properly mined (this might mean the same `Tipset` with invalid blocks removed, or an altogether different one (whose blocks have a different parent set).

If none were, you may need to mine null tickets instead (see the [expected consensus spec](./expected-consensus.md#null-blocks) for more). 

TODO: rename section in EC doc to be 'null tickets' and update link

### Ticket Generation

We detail ticket generation in the [expected consensus spec](./expected-consensus#ticket-generation).

### Mining a losing ticket with new blocks on the network

Generating a new ticket will take you some amount of time (as imposed by the VDF in Expected Consensus). If you find yourself with a losing ticket, on expectation you will hear about at least one other block being mined on the network. If so, you should verify the validity of these incoming blocks repeating the above process for a new `Tipset`.

TODO: find a better title here

#### Mining a losing ticket with no new blocks on the network

If no new blocks appear in the round, you may attempt to mine the same `Tipset` again. In order to do this, simply We call this mining a null block (i.e. mining atop the failed ticket you generated in your previous attempt). 

To start, you should insert your losing ticket into the `Tickets` array, then repeat the above process (from `Ticket Generation`) using your failed ticket from the previous round rather than the smallest ticket from the parent tipset (multiple null blocks in a row may be found).  This will generate a new ticket.

Repeat this process until you either find a winning ticket or hear about new blocks to mine atop of from the network. If a new block comes in from the network, and it is on a heavier chain than your own, you should abandon your null block mining to mine atop this new block. Due to the way chain selection works in filecoin, a chain with fewer null blocks will be preferred (see the [Expected Consensus spec](./expected-consensus.md#chain-selection) for more details).

#### Mining a winning ticket

If you mine a winning ticket, you may proceed to block creation thereby earning a block reward.

### Block Creation

When you have found a winning ticket, it's time to create your very own block!

To create a block, first compute a few fields:

- `Tickets` - An array containing a new ticket, and, if applicable, any intermediary tickets generated to prove appropriate delay for any null blocks you mined on. See [ticket generation](./expected-consensus.md#ticket-generation).
- `ElectionProof` - A signature over the final ticket from the `Tickets` array proving. See [ticket generation](./expected-consensus.md#ticket-generation).
- `ParentWeight` - As described in [Chain Weighting](./expected-consensus.md#chain-weighting).
- `ParentState` - Note that it will not end up in the newly generated block, but is necessary to compute to generate other fields. To compute this:
  - Take the `ParentState` of one of the blocks in your chosen parent set (invariant: this is the same value for all blocks in a given parent set).
  - For each block in the parent set, ordered by their tickets:
    - Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state.
    > NOTE (@porcuquine): What does it mean for messages to 'conflict'? Can we define it here or refer to the definition elsewhere?
    
    TODO: define this in the state-machine doc, and link to it from here
- `MsgRoot` - To compute this:
  - Select a set of messages from the mempool to include in your block.
  - Insert them into a Merkle Tree and take its root.
- `StateRoot` - Apply each of your chosen messages to the `ParentState` to get this.
- `ReceiptsRoot` - To compute this:
  - Apply the set of messages selected above to the parent state, collecting invocation receipts as you go.
  - Insert them into a Merkle Tree and take its root.
- `BlockSig` - A signature with your private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after we propogate it to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

Start by filling out `Parents`, `Tickets` and `ElectionProof` with values from the ticket checking process.

Next, compute the aggregate state of your selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of *their* parent tipset, sorting your parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipsets 'virtual block'. For more details on message execution and state transitions, see the [Filecoin state machine](state-machine.md) document.

Once you have the aggregate `ParentState`, you must apply the mining reward. This is done by adding the correct amount to the miner owner's account balance in the state tree. (TODO: link to block reward calculation. Currently, the block reward is a fixed 1000 filecoin).

Now, a set of messages is selected to put into the block. For each message, subtract `msg.GasPrice * msg.GasLimit` from the sender's account balance, returning a fatal processing error if the sender does not have enough funds (this message should not be included in the chain). Then apply the messages state transition, and generate a receipt for it containing the total gas actually used by the execution, the executions exit code, and the return value (see [receipt](data-structures#message-receipt) for more details). Then, refund the sender in the amount of `(msg.GasLimit - GasUsed) * msg.GasPrice`. In the event of a message processing error, the remaining gas is refunded to the user, and all other state changes are reverted. (Note: this is a divergence from the way things are done in Ethereum)

Each message should be applied on the resultant state of the previous message execution, unless that message execution failed, in which case all state changes caused by that message are thrown out. The final state tree after this process will be your blocks `StateRoot`.

Merklize the set of messages you selected, and put the root in `MsgRoot`. Gather the receipts from each execution into a set, merklize them, and put that root in `ReceiptsRoot`. Finally, set the `StateRoot` field with your resultant state.

Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.

Now the block is complete, all that's left is to sign it. Serialize the block now (without the signature field), take the sha256 hash of it, and sign that hash. Place the resultant signature in the `BlockSig` field.

#### Block Broadcast

Broadcast the completed block to the network (via [block propagation](data-propogation.md)), and assuming everything was done correctly, the network will accept it, and other miners will mine on top of it, earning you a block reward!

### Block Rewards


### Open Questions

- How should receipts for tipsets 'virtual blocks' be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.


### Future Work
There are many ideas for improving upon the storage miner, here we note some ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, so we won't implement it for the initial release of Filecoin, but this could provide some really interesting capabilities down the road.