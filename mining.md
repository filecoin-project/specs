# The Filecoin Mining Process

### What is the Filecoin Mining Process

An active participant in the filecoin consensus process is a storage miner and expected consensus block proposer. They are responsible for storing data for the filecoin network and also for driving the filecoin consensus process. Miners should constantly be performing Proofs of SpaceTime, and also checking if they have a winning `ticket` to propose a block for each round. We currently set rounds to take around 30 seconds, in order to account for network propagation around the world. The details of both processes are defined here.

While we refer to both storage miners and participants in expected consensus as "miners," strictly speaking only the latter are actively mining blocks (by participating in Filecoin consensus). Fulfilling storage orders and generating new blocks for block rewards are two wholly distinct ways to earn Filecoin.
With that said, it stands to reason that any storage miner would participate in Filecoin consensus (it effectively subsidizes their storage costs) and conversely, any participant in Filecoin consensus must be a storage miner (in order for them to appear in the power table). We therefore refer to these actors as "miners."

## Becoming a miner

### Registration

To become a miner, you must first register a new miner on-chain, set your pledge, and deposit your collateral. This is done through the storage market actor's `CreateMiner` method. Invoke that method with your desired pledge size, accompanying collateral, and public key for signature validation. The call will then create a new miner instance and return its address for you.

### Announcement

The next step is to place one or more storage market asks on the market. This is done through the storage markets `AddAsk` method. You may create a single ask for your entire storage, or partition your storage up in some way with multiple asks (at potentially different prices). 

### Deal Making

Once you have asks on the network, you must now wait for deal proposals from storage clients. Clients will look at all the miners announcing their prices, and use that information to select miners they want to store data with. As deal proposals come in and you accept the deals (TODO: add a section on deal acceptance strategy), you should start filling sectors with that data. Miners should continue to make deals until they run out of storage space.


### Committing Storage

Once a miner has filled up a sector with enough deals and sealed it, the next step is for them to commit that storage. To do this, they take the merkleroot commitments from the PoRep sealing setup and the seal proof, and call `CommitSector`.  Internally, `CommitSector` validates the proof, and then augments the miners `Power` by the appropriate amount.

It then checks if the miner is moving from having zero committed data to a non-zero amount of committed data, and if so sets the miners `ProvingPeriodStart` field to the current block.

TODO: sectors need to be globally unique. We can either do this by having the seal proof prove the sector is unique to this miner in some way, or by having a giant global map on-chain that we check against on each submission. I think that when we go towards sector aggregation, the latter option will become pretty much impossible, so we need to think about how that proof statement could work.


#### Sector Removal

When a miner no longer needs to store the data in a particular sector, they should remove it from their proving set by submitting the done sectors via the 'doneSet' parameter of `SubmitPost`. This will move the sectors to being 'removal candidates'. The collateral for those sectors remains locked up until the end of the next proving period, in case the removed sectors were still for valid deals and the client needs to slash them.

## The Miner Actor

**TODO** Where should this actually go? It feels a bit out of place right here.

After successfully calling `CreateMiner`, a miner actor will be created for you on-chain, and registered in the storage market. This miner, like all other Filecoin State Machine actors, has a fixed set of methods that can be used to interact with or control it.


```go
type StorageMinerActor interface {
    // AddAsk adds an ask via this miner to the storage markets orderbook
    AddAsk(price TokenAmount, expiry uint64) AskID

    // CommitSector adds a sector commitment to the chain.
    CommitSector(commD, commR, commRStar []byte, proof SealProof) SectorID

    // SubmitPoSt is used to submit a coalesced PoSt to the chain to convince the chain
    // that you have been actually storing the files you claim to be. A proof of
    // spacetime convinces the verifier that the sectors specified by the given
    // sector set are being correctly stored by the miner.
    // Also passed into this call are a set of FailureSets called 'faults'.
    // Each of these FailureSets specifies a set of sectors, and the time that they
    // failed at. During verification of the proof, these sectors will be removed.
    // recovered is the set of sectors that failed during the PoSt and were
    // recoverable by the miner. When making this call, the miner must submit enough
    // funds to pay fees for any faults that were recovered from. Collateral is deducted
    // for any sectors that were permanently lost.
    // The final parameter is the 'doneSet' this is the set of sectors the miner is
    // willingly removing, that they were able to prove for the entire proving period.
    SubmitPoSt(p PoSt, faults []FailureSet, recovered SectorSet, doneSet SectorSet)
    
    // IncreasePledge allows a miner to pledge more storage to the network
    IncreasePledge(addspace Integer)
    
    // SlashStorageFault is used to penalized this miner for missing their proofs
    SlashStorageFault()
    
    // GetCurrentProvingSet returns the current set of sectors that this miner is proving.
    // The next PoSt they submit will use this as an input. This method can be used by
    // clients to check that their data is being proven. (TODO: returning a list of commR 
    // hashes is a LOT of data, maybe find a better way?)
    GetCurrentProvingSet() [][]byte
    
    // ArbitrateDeal is called by a client of this miner whose deal with the miner has
    // data contained in a recently removed sector, and is not yet past the expiry date
    // of the deal.
    ArbitrateDeal(d Deal)
    
    // DePledge allows a miner to retrieve their pledged collateral.
    DePledge(amt TokenAmount)

    // GetOwner returns the address of the account that owns the miner. The owner is the
    // account that is authorized to control the miner, and is also where mining rewards
    // go to.
    GetOwner() Address

    // GetWorkerAddr (name still a WIP) returns the address corresponding to the
    // miners block signing key. Proof submissions for this miner must come from
    // this address. This is separate from the miners 'owner' address to allow
    // multiple miners to pay into a single address, while having separate signing keys
    // as well as to allow for miners to have their collateral put up by a different party
    // (e.g. for collateral pools)
    GetWorkerAddr() Address

    // GetPeerID returns the libp2p peer ID that this miner can be reached at.
    GetPeerID() libp2p.PeerID

    // UpdatePeerID is used to update the peerID this miner is operating under.
    UpdatePeerID(pid libp2p.PeerID)
}
```

The miner actor also has underlying state that is persisted on-chain. That state looks like this:

```go
type StorageMinerState struct {
    // Owner is the address of the account that owns this miner
    Owner Address
    
    // Worker is the address of the worker account for this miner
    Worker Address
    
    // PeerID is the libp2p peer identity that should be used to connect
    // to this miner
    PeerID peer.ID
    
    // PublicKey is the public portion of the key that the miner will use to sign blocks
    PublicKey PublicKey
    
    // PledgeBytes is the amount of space being offered by this miner to the network
    PledgeBytes BytesAmount
    
    // Collateral is locked up filecoin the miner has available to commit to storage.
    // When miners commit new sectors, tokens are moved from here to 'ActiveCollateral'
    // The sum of collateral here and in activecollateral should equal the required amount
    // for the size of the miners pledge.
    Collateral TokenAmount
    
    // ActiveCollateral is the amount of collateral currently committed to live storage
    ActiveCollateral TokenAmount
    
    // DePledgedCollateral is collateral that is waiting to be withdrawn
    DePledgedCollateral TokenAmount
    
    // DePledgeTime is the time at which the depledged collateral may be withdrawn
    DePledgeTime BlockHeight
    
    // Sectors is the set of all sectors this miner has committed
    Sectors SectorSet
    
    // ProvingSet is the set of sectors this miner is currently mining. It is only updated
    // when a PoSt is submitted (not as each new sector commitment is added)
    ProvingSet SectorSet
    
    // NextDoneSet is a set of sectors reported during the last PoSt submission as
    // being 'done'. The collateral for them is still being held until the next PoSt
    // submission in case early sector removal penalization is needed.
    NextDoneSet SectorSet
    
    // TODO: maybe this number is redundant with power
    LockedStorage BytesAmount
    
    // Power is the amount of power this miner has
    Power BytesAmount
}
```

### Owner Worker distinction

The miner actor has two distinct 'controller' addresses. One is the worker, which is the address which will be responsible for doing all of the work, submitting proofs, committing new sectors, and all other day to day activities. The owner address is the address that created the miner, paid the collateral, and has block rewards paid out to it. The reason for the distinction is to allow different parties to fulfil the different roles. One example would be for the owner to be a multisig wallet, or a cold storage key, and the worker key to be a 'hot wallet' key. 

### Storage Mining Cycle

Storage miners must continually produce proofs of space time over their storage to convince the network that they are actually storing the sectors that they have committed to. Each PoSt covers a miner's entire storage.

#### Step 0: Pre-Commit

Before doing anything else, a miner must first pledge some collateral for their storage and put up an ask to indicate their desired price.

After that, they need to make deals with clients and begin filling up sectors with data. For more information on making deals, see the section on [deal flow](storage-market.md#deal-flow)

When they have a full sector, they should seal it.

#### Step 1: Commit

When the miner has completed their first seal, they should post it on-chain using [CommitSector](actors.md#commit-sector). This starts their proving period.

The proving period is a fixed amount of time in which the miner must submit a Proof of Space Time to the network.

During this period, the miner may also commit to new sectors, but they will not be included in proofs of space time until the next proving period starts.


#### Step 2: Proving Storage (PoSt creation)

At the beginning of their proving period, miners collect the proving set (the set of all live sealed sectors on the chain at this point), and then call `ProveStorage`.

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


### Stop Mining

In order to stop mining, a miner must complete all of its storage contracts, and remove them from their proving set during a PoSt submission. A miner may then call `DePledge()` to retrieve their collateral (Note: depledging requires two calls to the chain, and a 'cooldown' period).

### Faults

Faults are described in the [faults document](../faults.md).

### On Being Slashed (WIP, needs discussion)

If a miner is slashed for failing to submit their PoSt on time, they currently lose all their pledge collateral. They do not necessarily lose their storage collateral. Storage collateral is lost when a miners clients slash them for no longer having the data. Missing a PoSt does not necessarily imply that a miner no longer has the data. There should be an additional timeout here where the miner can submit a PoSt, along with 'refilling' their pledge collateral. If a miner does this, they can continue mining, their mining power will be reinstated, and clients can be assured that their data is still there.

Review Discussion Note: Taking all of a miners collateral for going over the deadline for PoSt submission is really really painful, and is likely to dissuade people from even mining filecoin in the first place (If my internet going out could cause me to lose a very large amount of money, that leads to some pretty hard decisions around profitability). One potential strategy could be to only penalize miners for the amount of sectors they could have generated in that timeframe. 

## Mining Blocks

Now that you are a real life filecoin miner, it's time to start making and checking tickets. At this point, you should already be running chain validation, which includes keeping track of the latest [`Tipsets`](./expected-consensus.md#tipsets) that you've seen on the network.

For additional details around how consensus works in Filecoin, see the [expected consensus spec](./expected-consensus.md). For the purposes of this section, we have a consensus protocol (Expected Consensus) that guarantees us a fair process for determining what blocks have been generated in a round and whether a miner should mine a block themselves, and some rules pertaining to how "Tickets" should be validated during block validation.

### Receiving Blocks

Receiving blocks from the network, you must do the following:

1. Check their validity (see below).
2. Assemble a `Tipset` with all valid blocks with common parents.

You may sometimes receive blocks belonging to different `Tipsets` (i.e. whose parents are not the same). In that case, you must choose which tipset to mine on.

Chain selection is a crucial component of how the Filecoin blockchain works. In short, every chain has an associated weight accounting for the number of blocks mined on it and the power (storage) they track. It is always preferable to mine atop a heavier block rather than a lighter one. While you may be foregoing block rewards earned in the past, this lighter chain is likely to be abandoned by other miners forfeiting any block reward earned. For more on this, see [chain selection](./expected-consensus.md#chain-selection) in the Expected Consensus spec.

### Block structure

In order to go over block validation in detail, we quickly go over block structure (see more in [block creation](#block-creation)).

Note: This may not yet be the format the go-filecoin codebase is yet using, but it describes what the 'ideal' block structure should look like. Implementations may name things differently, or use different types. For a detailed layout of exactly how a block should look and be serialized, see [the datastructures spec](data-structures.md#block).

We have:

```go
type Block struct {
    // Parents are the cids of this blocks parents
    Parents []*cid.Cid

    // Tickets is the array of tickets used to prove delay in block generation.
    // Generated from tickets in its parent tipset.
    // On expectation there will be only one, but there could be multiple, null blocks
    // on which the winning ticket was found.
    Tickets []Signature

    // ElectionProofs is the array of proofs that the leader was indeed elected.
    // Generated from tickets in the lookback tipset.
    // On expectation there will be only one, but there could be multiple, null blocks
    // on which the winning ticket was found.
    ElectionProofs []Signature
    
    // Parent Weight is the aggregate chain weight of this block's parent set.
    ParentWeight float64
        
    // MsgRoot is the root of the merklized set of state transitions in this block
    MsgRoot *cid.Cid
    
    // ReceiptsRoot is the root of the merklized set of invocation receipts matching
    // the messages in this block.
    ReceiptsRoot *cid.Cid
    
    // StateRoot is the resultant state represented by this block after the in-order 
    // application of its set of messages.
    StateRoot *cid.Cid
    
    // BlockSig is a signature over all the other fields in the block with the miners
    // private key
    BlockSig Signature
}
```

To start, you must validate that the block was well mined.

### Block Validation

In order to validate a block coming in from the network at height `N` was well mined you must do the following:

1. Validate `BlockSig`

2. Ensure that the block content is appropriately referenced in the header. 

3. Validate `ParentWeight`
   - Each of the blocks in the block's `ParentTipset` must have the same `ParentWeight`.
   - The new block's ParentWeight must have been properly calculated using the [chain weighting function](./expected-consensus#chain-weighting).

4. Validate `StateRoot`
   - In order to do this, you must ensure that applying the block's messages to the `ParentState` (not included in block header) appropriately yields the `StateRoot` and receipts in the `ReceiptRoot`.

5. You must validate that the tickets in the `Tickets` array are valid, thereby proving appropriate delay for this block creation.
   - Ensure that the new ticket was generated from the smallest ticket in the block's parent tipset (at height `N-1`).
   - Recompute the new ticket, using the miner's public key, ensuring it was computed appropriately.

6. You must validate that the tickets in `ElectionProofs` array were correctly generated by the miner, and that they are eligible to mine.
   - Ensure that the proof was generated using the smallest ticket at height `N-K` (the lookback tipset).
   - Recompute the proof, using the miner's public key, and check that the ticket is a valid signature over that value.
   - Verify that the proof is indeed smaller than the miner's power as reported in the power table at height `N-L`.

7. In a case where the block contains multiple values in either the `Tickets` or `ElectionProofs` arrays
   - Ensure that all tickets in both arrays are signed by the same key.
   - For the `Tickets` array, ensure that each ticket was used to generate the next, starting from the smallest ticket in the Parent tipset (at `N-1`).
   - For the `ElectionProofs` array, ensure that each ticket was used to generate the next, starting from the smallest ticket in the Lookback tipset (at `N-K`).
   - Ensure that both arrays contain the same number of tickets.

We detail ticket validation as follows:

```go
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
    lookbackTicket := selectSmallestTicket(randomnessLookback(b))

    for _, ticket := range b.ElectionProofs {
            challenge := sha256.Sum(lookbackTicket)
	
        	// Check VRF
            pubk := getPublicKeyForMiner(b.Miner)
            if !pubk.VerifySignature(ticket, lookbackTicket) {
                return "Ticket was not a valid signature over the lookback ticket"
            }
            // in case this block was mined atop null blocks
        lookbackTicket = ticket
    }
    
    state := getStateTree(powerLookback(b))
    minersPower := state.getPowerForMiner(b.Miner)
    totalPower := state.getTotalPower()
    
    if !IsTicketAWinner(lookbackTicket, minersPower, totalPower) {
       return "Ticket was not a winning ticket"
    }
    
    // Winner!
    return nil
}
```

If all of this lines up, the block is valid. Repeat for all blocks in a tipset.

Once you've ensured all blocks in the `Tipset` received were properly mined, you can mine on top of it. If it wasn't, ensure the next heaviest `Tipset` was properly mined (this might mean the same `Tipset` with invalid blocks removed, or an altogether different one (whose blocks have a different parent set).

If none were, you may need to mine null blocks instead (see the [expected consensus spec](./expected-consensus.md#null-blocks) for more). 

### Ticket Generation

We detail ticket generation in the [expected consensus spec](./expected-consensus#ticket-generation).

### Mining a losing ticket with new blocks on the network

Generating a new ticket will take you some amount of time (as imposed by the VDF in Expected Consensus). If you find yourself with a losing ticket, on expectation you will hear about at least one other block being mined on the network. If so, you should verify the validity of these incoming blocks repeating the above process for a new `Tipset`.

#### Mining a losing ticket with no new blocks on the network

If no new blocks appear in the round, you may attempt to mine the same `Tipset` again. We call this mining a null block (i.e. mining atop the failed ticket you generated in your previous attempt). 

To start, you should insert your losing ticket into the `Tickets` array, then repeat the above process (from `Ticket Generation`) using your failed ticket from the previous round rather than the smallest ticket from the parent tipset (multiple null blocks in a row may be found).  This will generate a new ticket.

Repeat this process until you either find a winning ticket or hear about new blocks to mine atop of from the network. If a new block comes in from the network, and it is on a heavier chain than your own, you should abandon your null block mining to mine atop this new block. Due to the way chain selection works in filecoin, a chain with fewer null blocks will be preferred (see the [Expected Consensus spec](./expected-consensus.md#chain-selection) for more details).

#### Mining a winning ticket

If you mine a winning ticket, you may proceed to block creation thereby earning a block reward.

### Block Creation

When you have found a winning ticket, it's time to create your very own block!

To create a block, first compute a few fields:

- `Tickets` - An array containing a new ticket, and, if applicable, any intermediary tickets generated to prove appropriate delay for any null blocks you mined on. See [ticket generation](./expected-consensus.md#ticket-generation).
- `ElectionProofs` - An array containing your winning ticket proving election, and, if applicable, the failed intermediary tickets for any null blocks you mined on. See [ticket generation](./expected-consensus.md#ticket-generation).
- `ParentWeight` - As described in [Chain Weighting](./expected-consensus.md#chain-weighting).
- `ParentState` - Note that it will not end up in the newly generated block, but is necessary to compute to generate other fields. To compute this:
  - Take the `ParentState` of one of the blocks in your chosen parent set (invariant: this is the same value for all blocks in a given parent set).
  - For each block in the parent set, ordered by their tickets:
    - Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state.
- `MsgRoot` - To compute this:
  - Select a set of messages from the mempool to include in your block.
  - Insert them into a Merkle Tree and take its root.
- `ReceiptsRoot` - To compute this:
  - Apply the set of messages selected above to the parent state, collecting invocation receipts as you go.
  - Insert them into a Merkle Tree and take its root.
- `StateRoot` - Apply each of your chosen messages to the `ParentState` to get this.
- `BlockSig` - A signature with your private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after we propogate it to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

Start by filling out `Parents`, `Tickets` and `ElectionProofs` with values from the ticket checking process.

Next, compute the aggregate state of your selected parent blocks, the `ParentState`. This is done by taking the aggregate parent state of *their* parent tipset, sorting your parent blocks by their tickets, and applying each message in each block to that state. Any message whose nonce is already used (duplicate message) in an earlier block should be skipped (application of this message should fail anyway). Note that re-applied messages may result in different receipts than they produced in their original blocks, an open question is how to represent the receipt trie of this tipsets 'virtual block'.

Once you have the aggregate `ParentState`, select a set of messages to put into your block. The miner may include a block reward message in this set.  For this reward message to pass validation it must be the first message in the serialized list of messages and must not claim more than the protocol defined block reward constant.  The `From` field of the block reward message must equal the protocol defined network address constant.  The message should not include a signature.  Apply each message of the message set to your aggregate parent state in order to compute your block's resultant state. Gather the receipts from each execution into a set. Now, merklize the set of messages you selected, and put the root in `MsgRoot`. Merklize the receipts, and put that root in `ReceiptsRoot`. Finally, set the `StateRoot` field with your resultant state.

Note that the `ParentState` field from the expected consensus document is left out, this is to help minimize the size of the block header. The parent state for any given parent set should be computed by the client and cached locally.

Now you have a filled out block, all that's left is to sign it. Serialize the block now (without the signature field), take the sha256 hash of it, and sign that hash. Place the resultant signature in the `BlockSig` field, and then you're done.

#### Block Broadcast

Broadcast your completed block to the network, and assuming you've done everything correctly, the network will accept it, and other miners will mine on top of it, earning you a block reward!

### Open Questions

- How should receipts for tipsets 'virtual blocks' be referenced? It is common for applications to provide the merkleproof of a receipt to prove that a transaction was successfully executed.


### Future Work
There are many ideas for improving upon the storage miner, here we note some ideas that may be potentially implemented in the future.

- **Sector Resealing**: Miners should be able to 're-seal' sectors, to allow them to take a set of sectors with mostly expired pieces, and combine the not-yet-expired pieces into a single (or multiple) sectors.
- **Sector Transfer**: Miners should be able to re-delegate the responsibility of storing data to another miner. This is tricky for many reasons, so we won't implement it for the initial release of Filecoin, but this could provide some really interesting capabilities down the road.
