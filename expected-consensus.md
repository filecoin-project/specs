# Expected Consensus

**This spec describes how the expected consensus (EC) protocol works in general. To read more about Filecoin-specific processes, see:**

- [Mining Blocks](mining.md#mining-blocks) on how consensus is used in block mining.
- [Faults](faults.md) on slashing.
- [Storage market](storage-market.md#the-power-table) on how the `power table` is created and maintained.
- [Block data structure](data-structures.md#block) for details on fields and encoding.

## Important concepts and definitions
Some important concepts relevant to expected consensus are:
- [Verifiable Delay Function (VDF)](definitions.md#vdf)
- [Verifiable Random Function (VRF)](defintions.md#vrf)
- [TipSet](definitions.md#tipset)
- [Height](definitions.md#height)
- [Weight](definitions.md#weight)
- [Round](definitions.md#round) -- In the realm of EC, it is worth noting that a new ticket is produced at every round, consequently the duration of a round is currently bounded by the duration of the Verifiable Delay Function run to generate a ticket.
- [Power Fraction](definitions.md#power-fraction)
- [ElectionProof](definitions.md#electionproof)

## Algorithm

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one
participant may be eligible to submit a block. EC guarantees that this winner will be anonymous until they reveal themselves by submitting a proof of their election (we call this proof an `Election Proof`). All valid blocks submitted in a given round form a `TipSet`. Every block in a TipSet adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see [Chain Selection](#chain-selection).

The methods below describe the basic algorithm for EC.

Every time a block is received, we first verify whether the block is valid (`VerifyBlock` is defined in the [mining spec](mining.md#block-validation).) If the received block is valid, we add it to our TipSet.

```go
// Called when a block is received by this node
func OnBlockReceived(blk Block) {
	// The exact definition of VerifyBlock depends on the protocol
	// For Filecoin, see mining.md
	if VerifyBlock(blk) {
		ChainTipsMgr.Add(blk)
	}
  // Received an invalid block!
}
```

Meanwhile, the mining node runs a mining process to attempt to generate blocks. In `Mine`, the node identifies the best TipSet in each round and generates a ticket from it. In parallel, it uses a past ticket to try and generate a valid `ElectionProof` thereby enabling it to mine a new block. If no valid `ElectionProof` is produced, the miner mines a new ticket atop their old one and tries again at a new height.

```go
func Mine(minerKey PrivateKey) {
	for r := range rounds { // for each round
		bestTipset, tickets := ChainTipsMgr.GetBestTipsetAtRound(r - 1)

		ticket := GenerateTicket(minerKey, bestTipset)
		tickets.Append(ticket)

		// Generate an ElectionProof and check if we win
		// Note: even if we don't win, we want to make a ticket
		// in case we need to mine a null round
		win, proof := CheckIfWinnerAtRound(minerKey, r, bestTipset)
		if win {
			GenerateAndBroadcastBlock(bestTipset, tickets, proof)
		} else {
			// Even if we don't win, add our ticket to the tracker in
			// case we need to mine on it later.
			ChainTipsMgr.AddFailedTicket(bestTipset, tickets)
		}
	}
}
```

To check if the miner won the round, she runs `CheckIfWinnerAtRound`. In this method, the miner takes their ticket from a prior round (which round to look at is specified by the randomness lookback parameter), computes an ElectionProof, and returns whether the proof indicates that the miner has won the round. In the pseudocode below, `IsProofAWinner` is taken from [the mining doc](mining.md#block-validation).

```go
const RandomnessLookback = 1 // Also referred to as "K" in many places

func CheckIfWinnerAtRound(key PrivateKey, n Integer, parentTipset Tipset) (bool, ElectionProof) {
  lbt := ChainTipsMgr.TicketFromRound(parentTipset, n-RandomnessLookback)

  eproof := ComputeElectionProof(lbt, key)

  minerPower := GetPower(miner)
  totalPower := state.GetTotalPower()

  if IsProofAWinner(eproof, minerPower, totalPower)
    return eproof
  else
    return 0
}
```

Note: Validity of blocks beyond appropriate ticket generation (defined below) is defined by the specific protocol using EC. For the Filecoin definition of a valid block, see the [mining spec](mining.md).

The EC algorithm can be better understood by looking at its two major components in more detail:
- [Leader Election](#secret-leader-election)
- [Chain Selection](#chain-selection)

## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the storage [power table](storage-market.md#the-power-table), where power is equivalent to storage provided through time.

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain tracks an independent ticket chain. These tickets are used as randomness inputs for Leader Election. Every block generated references an `ElectionProof` derived from a past ticket. The ticket chain is extended by the miner who generates a new ticket with each attempted election.

In cases in which no winning ticket is found by any miner in a given round (i.e. no block, or a `null block`, is mined on the network), miners move on to the next ticket in the ticket chain to attempt a new leader election. Note that a null block is not included in a TipSet, since null blocks are, by definition, not blocks. When this happens, miners should nonetheless generate a new ticket prior to the new leader election, thereby appropriately prolonging the ticket chain (the block chain can never be longer than the ticket chain). This situation is fleshed out in the [Losing Tickets](#losing-tickets) section.

In order to pressure the network to converge on a single chain, each miner may only submit one block per round (see: [`Slashing`](#slashing)).

TODO: pictures of ticket chain and block chain

### Tickets

One may think of leader election in EC as a verifiable lottery, in which participants win in proportion to the power they have within the network.

A ticket is drawn from the past at the beginning of each new round, and a new ticket is generated in every round. Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any other data in the block. Nonetheless, in Filecoin, every block header contains one or more new tickets, thereby extending the ticket chain. A miner generates a new ticket in their block for every ticket they scratch running leader election, thereby ensuring the ticket chain is at least as long as the block chain.

At a high-level, tickets must do the following:

- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.
- Prove appropriate delay between drawings — thereby preventing leaders from "rushing" the protocol by releasing blocks early (at the expense of fairness for miners with worse connectivity).
- Ensure a single drawing per round — derived in part from the above, thereby preventing miners from grinding on tickets (e.g. by repeatedly drawing new tickets in the hopes of winning) within a round.

```text
ticket = {TODO} where (proof, value) <-- VDF(SK, x) for some seed x
```
You can find the Ticket data structure [here](data-structures.md#tickets).

In practice, EC defines two different fields within a block:

- A `Tickets` array — this stores new tickets generated during this epoch, or block generation attempt. It proves appropriate delay. It is from this array that miners will sample randomness to run leader election in `K` rounds. See [Ticket generation](#ticket-generation).
- An `ElectionProof` — this stores a proof that a given miner scratched a winning lottery ticket using the appropriate ticket `K` rounds back. It proves that the leader was elected in this round. See [Checking election results](#checking-election-results).

On expectation, the `Tickets` array will contain a single ticket. For cases in which it contains more than one, see [Losing Tickets](#losing-tickets).

```
On the two uses of tickets.

Tickets serve two main purposes.

1) As a proof that blocks were appropriately delayed. This both prevents miners grinding through tickets in a given epoch and a winner rushing the protocol by publishing
blocks as fast as possible (thereby penalizing poorly connected miners).

This is a result of a new ticket being generated by the miner for every leader election
attempt (using a ticket from the past). This ticket is generated using the ticket from
the prior block, or "losing" tickets generated in this epoch.

2) As a source of randomness used to prove that a leader was correctly elected.

This is done by generating an 'ElectionProof' derived from a ticket sampled K rounds
back.

But why the randomness lookback?

The randomness lookback helps turn independent lotteries (ticket drawings from a block one round back)
into a global lottery instead. Rather than having a distinct chance of winning or losing
for each potential fork in a given round, a miner will either win on all or lose on all
forks descended from the block in which the ticket is sampled.

This is useful as it reduces opportunities for grinding, across forks or sybil identities.

However this introduces a tradeoff:
- The randomness lookback means that a miner can know K rounds in advance that they will win,
decreasing the cost of running a targeted attack (given they have local predictability).
- It means electionProofs are stored separately from new tickets on a block, taking up
more space on-chain.
```

#### Ticket generation

This section discusses how tickets are generated for the `Tickets` array. For how tickets are validated, see [ticket validation](mining.md#ticket-validation).

At round `N`, new tickets are generated using tickets drawn from the [TipSet](#tipsets) at round `N-1`. This ensures the miner cannot publish a new block (corresponding to the `ElectionProof` generated by a winning ticket `K` rounds back) until the correct round. Because a Tipset can contain multiple blocks (see [Chain Selection](#chain-selection) below), the smallest ticket in the Tipset must be drawn otherwise the block will be invalid.

TODO: pictures of TipSet ticket drawing

The miner runs the prior ticket through a Verifiable Random Function (VRF) to get a new unique output. This output is then used as input into a Verifiable Delay Function (VDF), with the VRFProof, VDFProof and VDFOutput generating a new ticket for future use. 

The VRF's deterministic output adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance). The VDF approximates clock synchrony for miners, thereby ensuring miners have waited an appropriate delay ahead of drawing a new ticket. It also establishes a lower-bound delay between production of new blocks at different heights. This helps ensure fairness, in that miners with lesser connectivity are not penalized and have a chance to produce blocks.

Succinctly, the process of crafting a new `Ticket` in round `N` is as follows. We use:

- The ECVRF algorithm from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve
  - Note that the operation type in step 2.1 is necessary to prevent an adversary from guessing an election proof for a miner ahead of time.
- The TODO VDF impl.

```text
Input: parentTickets at round N-1, miner's private key SK
Output: newTicket

0. Prepare new ticket
	i. newTicket <-- New()
1. Draw prior ticket
	i. 	 # take the last ticket for each parent 'Tickets' array
			lastTickets <-- map(parentTickets, fun(x): x.last)
	ii.  # sort these tickets 
  			sortedTickets <-- Sort(lastTickets)
    iii. # draw the smallest ticket
  			parentTicket <-- min(sortedTickets)
2. Run it through VRF and get determinstic output
	i.   # take the VDFOutput of that ticket as input, specifying the personalization (see data-structures)
			input <-- VRFPersonalization.Ticket | parentTicket.VDFOutput
	ii.	 # run it through the VRF and store the VRFProof in the new ticket
			newTicket.VRFProof <-- ECVRF_prove(SK, input)
	iii. # draw a deterministic output from this
			VRFOutput <-- ECVRF_proof_to_hash(newTicket.VRFProof)
3. Run that deterministic output through a VDF
    i.  # run eval with our VDF and its evaluation k on VRFOutput
  			y, pi <-- Eval(ek, VRFOutput)
    ii. # Store the output and proof in our ticket
  			newTicket.VDFOutput <-- y
  			newTicket.VDFProof 	<-- pi
4. Return the new ticket
```

### Checking election results

Now, a miner must also check whether they are eligible to mine a block in this round. For how Election Proofs are validated, see [election validation](mining.md#election-validation).

To do so, the miner will use tickets from K blocks back as randomness to uniformly draw a value from 0 to 1. Comparing this value to their power, they determine whether they are eligible to mine. A user's `power` is defined as the ratio of the amount of storage they proved as of their last PoSt submission to the total storage in the network as of the current block.

Succinctly, the process of crafting a new `ElectionProof` in round `N` is as follows. We use:

- The ECVRF algorithm (must yield a pseudorandom, deterministic output) from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve
  - Note that the operation type in step 3.1 is not strictly necessary, but is used to distinguish this use of the VRF from that which generates tickets.

At round N:

```text
Input: parentTickets from N-K, miner's public key PK, miner's secret key SK, the Storage Market actor S
Output: 1 or 0

0. Prepare new election proof
	i. newEP <-- New()
1. Determine the miner's power fraction
	i. 	# Determine total storage this round
  		S_n <-- storageMarket(N)
  		p_n <-- S_n.GetTotalStorage()
    ii. # Determine own power as of last submitted PoSt
        appropriateHeight <-- self.ProvingPeriodEnd - ProvingPeriodDuration(self.SectorSize)
        S_m <-- storageMarket(appropriateHeight)
  		p_m <-- S_m.PowerLookup(self)
    iii. # Get power fraction
  		p_f <-- p_m/p_n
2. Draw parentTicket from K blocks back (see ticket creation above for example) 
3. Run it through VRF and get determinstic output
	i.   # take the VDFOutput of that ticket as input, specified for the appropriate operation type
		input <-- VRFPersonalization.ElectionProof | parentTicket.VDFOutput
	ii.	 # run it through the VRF and store the VRFProof in the new ticket
		newEP.VRFProof <-- ECVRF_prove(SK, input)
	iii. # draw a deterministic, pseudorandom output from this
		VRFOutput <-- ECVRF_proof_to_hash(newEP.VRFProof)
3. Determine if the miner drew a winning lottery ticket
	i.  # Map the VRFOutput onto [0,1], with HashLen of 32 Bytes using sha264
  		scratchValue <-- VRFOutput / {1}^HashLen
    ii. # Compare the miner's scratchValue to the miner's power fraction
        # winning ticket
        if scratchValue <= p_f
            return newEP
        # otherwise parentTicket is not a winning lottery ticket
        else 
            Return 0
```

If the miner scratches a winning ticket in this round, it can use newEP, along with a newTicket to generate and publish a new block (see [Block Generation](#block-generation)). Otherwise, it waits to hear of another block generated in this round.

It is important to note that every block contains two artifacts: one, a ticket derived from last block's ticket to prove that they have waited the appropriate delay, and two, an election proof derived from the ticket `K` rounds back used to run leader election.

#### Losing Tickets

In the case that everybody draws a losing ticket in a given round (i.e. no miner is eligible to produce a block), every miner can run leader election again by "scratching" (attempting to generate a new `ElectionProof` from) the next ticket in the chain. That is, miners will now use the ticket sampled `K-1` rounds back to generate a new `ElectionProof`. They can then compare that proof with their power in the table `N-(L-1)` rounds back. This is repeated until a miner scratches a winning ticket and can publish a block (see [Block Generation](#block-generation)).

In addition to each attempted `ElectionProof` generation, the miner will need to extend the ticket chain by generating another new ticket. They use the ticket they generated in the prior round, rather than the prior block's (as is normally used). This proves appropriate delay (given that finding a winning Ticket has taken multiple rounds).

Thus, each time it is discovered that nobody has won in a given round, every miner should append a new ticket to their would-be block's `Ticket` array. This continues until some miner finds a winning ticket (see below), ensuring that the ticket chain remains at least as long as the block chain.

The length of repeated losing tickets in the ticket chain (equivalent to the length of generated tickets referenced by a single block, or the length of the `Tickets` array) decreases exponentially in the number of repeated losing tickets (see more [here](https://github.com/filecoin-project/go-filecoin/pull/1516)). In the unlikely case the number of losing tickets drawn by miners grows larger than the randomness lookback `K` (i.e. if a miner runs out of existing tickets on the ticket chain for use as randomness), a miner should proceed as usual using new tickets generated in this epoch for randomness. This has no impact on the protocol safety/validity.

New blocks (with multiple tickets) will have a few key properties:

- All tickets in the `Tickets` array are signed by the same miner -- to avoid grinding through out-of-band collusion between miners exchanging tickets.
- The `ElectionProof` was correctly generated from the ticket `K-|Tickets|-1` (with `|Tickets|` the length of the `Tickets` array) rounds back.

This means that valid `ElectionProof`s can be generated from tickets in the middle of the `Tickets` array.

#### A note on miners' 'power fraction'

The portion of blocks a given miner generates (and so the block rewards they earn) is proportional to their `Power Fraction` over time.

This means miners should not be able to mine using power they have not yet proven. Conversly, it is acceptable for miners to mine with a slight delay between their proving storage and that proven storage being reflected in leader election. This is reflected in the height at which the [storage market actor](actors.md#storage-market-actor)'s `GetTotalStorage` and `PowerLookup` methods are called, as outlined [above](#checking-election-results).

The miner retrieves the appropriate state of the `storage market actor` at the height at which they submitted their last valid PoSt: that way, they account for storage which was already proven valid, i.e. they are mining with the power they had in their last proving period. That is the miner will get their own power at height `miner.ProvingPeriodEnd - miner.ProvingPeriod`.

This power is then compared to the total network power at the current height, in order to account for recent power changes from other miners (and so all miners attempting to mine in this round share 100% of the power).

To illustrate this, an example:

Miner M1 has a provingPeriod of 30. M1 submits a PoST at height 39. Their next `provingPeriodEnd` will be 69, but M1 can submit a new PoST at any height X, for X in (39, 69]. Let's assume X is 67.

At height Y in (39, 67], M1 will attempt to generate an `ElectionProof` using the storage market actor from height 39 for their own power (and an actor from Y for total network power); at height 68, M1 will use the storage market actor from height 67 for their own power, and the storage market actor from height 68 for total power and so on.

### Block Generation

Once a miner has a winning election proof generated over `i` rounds and `i` corresponding tickets, they may create a block. For more on this, see the [Mining spec](mining.md#block-creation).

## Chain Selection

Just as there can be 0 miners win in a round, multiple miners can be elected in a given round. This in turn means multiple blocks can be created in a round. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round.

### Tipsets

All valid blocks generated in a round form a `TipSet` that participants will attempt to mine off of in the subsequent round (see above). TipSets are valid so long as:

- All blocks in a TipSet have the same parent TipSet
- All blocks in a TipSet have the same number of tickets in their `Tickets` array

The first condition implies that all blocks in a TipSet were mined at the same height (remember that height refers to block height as opposed to ticket round). This rule is key to helping ensure that EC converges over time. While multiple new blocks can be mined in a round, subsequent blocks all mine off of a TipSet bringing these blocks together. The second rule means blocks in a TipSet are mined in a same round.

The blocks in a tipset have no defined order in representation. During state computation, blocks in a tipset are processed in order of block ticket, breaking ties with the block CID bytes.

Due to network propagation delay, it is possible for a miner in round N+1 to omit valid blocks mined at round N from their TipSet. This does not make the newly generated block invalid, it does however reduce its weight and chances of being part of the canonical chain in the protocol.

### Chain Weighting

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided.

The weight at each block is equal to its `ParentWeight`, plus that block's delta weight. Delta
weight is a constant `V`, plus `X` - a function of the total power in the network as reported in the Power Table.  The exact value for `V` and the magnitude of the power ratio value are
still to be determined, but for now `V = 10` and `X = log(TotalPower)`.

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given TipSet should have
the same `ParentWeight` value) plus the delta weight of each parent. To make the
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise
potentially long chain scans would be required to compute a given block's weight).

### Selecting between TipSets with equal weight

When selecting between TipSets of equal weight, a miner chooses the one with the smallest min ticket (by bytewise comparison).

In the case where two TipSets of equal weight have the same min ticket, the miner will compare the next smallest ticket (and select the TipSet with the next smaller ticket). This continues until one TipSet is selected.

The above case may happen in situations under certain block propagation conditions. Assume three blocks B, C, and D have been mined (by miners 1, 2, and 3 respectively) off of block A, with minTicket(B) < minTicket(C) < minTicket (D).

Miner 1 outputs their block B and shuts down. Miners 2 and 3 both receive B but not each others' blocks. We have miner 2 mining a TipSet made of B and C and miner 3 mining a TipSet made of B and D. If both succesfully mine blocks now, other miners in the network will receive new blocks built off of TipSets with equal weight and the same smallest ticket (that of block B). They should select the block mined atop [B, C] since minTicket(C) < minTicket(D).

The probability that two TipSets with different blocks would have all the same tickets can be considered negligible: this would amount to finding a collision between two 256-bit (or more) collision-resistant hashes.

### Slashing

See the [Faults spec](faults.md) for implementation details.

Due to the existence of potential forks, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This is detectable when a miner submits two blocks that satisfy either of the following "slashing conditions":

(1) one block contains at least one ticket in its ticket array generated at the same round as one of the tickets in the other block's ticket array.
(2) one block's parent is a TipSet that could have validly included the other block according to TipSet validity rules, however the parent of the first block does not include the other block.

  - While it cannot be proven that a miner omits known blocks from a TipSet in general (i.e. network latency could simply mean the miner did not receive a particular block) in this case it can be proven because a miner must be aware of a block they mined in a previous round.

Any node that detects this occurring should take both block headers, and call [`storagemarket.SlashConsensusFault`](actors.md#slashconsensusfault). The network will then take all of that node's collateral, give a portion of it to
the reporter, and keep the rest.

TODO: It is unclear that rewarding the reporter any more than gas fees is the right thing to do. Needs thought. Tracking issue: https://github.com/filecoin-project/specs/issues/159

Note: One may wonder what prevents miners from simply breaking up their power into multiple unlinkable miner actors  (or sybils) that will be able to mine on multiple chains without being caught mining at the same round at the same time. Read more about this [here](https://github.com/filecoin-project/consensus/issues/32).

### ChainTipsManager

The Chain Tips Manager is a subcomponent of Filecoin consensus that is technically up to the implementer, but since the pseudocode in previous sections reference it, it is documented here for clarity.

The Chain Tips Manager is responsible for tracking all live tips of the Filecoin blockchain, and tracking what the current 'best' tipset is.

```go
// Returns the ticket that is at round 'r' in the chain behind 'head'
func TicketFromRound(head Tipset, r Round) {}

// Returns the tipset that contains round r (Note: multiple rounds may exist within a single // tipset due to null blocks)
func TipsetFromRound(head Tipset, r Round) {}

// GetBestTipset returns the best known tipset. If the 'best' tipset hasnt changed, then this
// will return the previous best tipset and the null block we mined on top of it.
func GetBestTipset()

// Adds the failed ticket to the chaintips manager so that null blocks can be mined on top of
func AddFailedTicket(parent Tipset, t Ticket)
```

## Implementation Notes

- When selecting messages from the mempool to include in the block, be aware that other miners may also generate blocks during this round, and to maximize fee earnings it may be best to select some messages at random (second in a duplicate earns no fees).

## Open Questions

- Parameter K
- VDF difficulty adjustment
