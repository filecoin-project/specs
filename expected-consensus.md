# Expected Consensus

This spec describes how to implement the protocol in general, for Filecoin-specific processes, see:

- [Mining Blocks](mining.md#mining-blocks) on how consensus is used.
- [Faults](./faults.md) on slashing.
- [Storage Market](./storage-market.md#the-power-table) on how the power table is maintained.
- [Block data structure](data-structures.md#block) for details on fields and encoding. 

Important Concepts
- TipSet

- Weight

- Round

- Height

- Epoch

- Ticket

- ElectionProof

  

>  **Definition:** A `ticket` is used as a source of randomness in EC leader election. Every block depends on an `ElectionProof` derived from a ticket. At least one new ticket is produced with every new block.

> **Definition:** An `ElectionProof` is derived from a past ticket and included in every block header. The `ElectionProof` proves that the miner was eligible to mine a block in that round.

> **Definition:** A `round` is the period in which a miner runs leader election to attempt to generate a new block. A new ticket is produced at every round, consequently the duration of a round is currently bounded by the duration of the Verifiable Delay Function run to generate a ticket. Tickets can be counted by round as one will be produced for each round.

> **Definition:** `Height` refers to the number of `TipSets` a given block is built upon since genesis (height 0). Multiple blocks in a `TipSet` will have the same height.

> **Defintion:** An `epoch` is the period in which a new block is generated. There may be multiple rounds in an epoch.

> **Definition:** A `TipSet` is a set of blocks with the same parent set, and same number of tickets (which implies they will have been mined at the same height).

Expected Consensus is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one
participant may be eligible to submit a block. All valid blocks submitted in a given round form a `TipSet`. Every block in a TipSet adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on this, see [Chain Selection](#chain-selection).

The basic algorithms are as follows:

For each block received over the network, `OnBlockReceived` is called. `VerifyBlock` is defined in the [mining spec](mining.md#block-validation).

```go
func OnBlockReceived(blk Block) {
	// The exact definition of IsValid depends on the protocol
	// For Filecoin, see mining.md
	if VerifyBlock(blk) {
		ChainTipsMgr.Add(blk)
	}

	// Received an invalid block!
}
```

Separately, another process is running `Mine` to attempt to generate blocks.
```go
func Mine(minerKey PrivateKey) {
	for r := range rounds { // for each round
		bestTipset, tickets := ChainTipsMgr.GetBestTipsetAtRound(r - 1)

		ticket := GenerateTicket(minerKey, bestTipset)
		tickets.Append(ticket)

		// Generate an election proof and check if we win
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

`IsProofAWinner` is taken from [the mining doc](mining.md#block-validation).

```go
const RandomnessLookback = 1 // Also referred to as "K" in many places
const PowerLookback = 1      // Also referred to as "L" in many places

func CheckIfWinnerAtRound(key PrivateKey, n Integer, parentTipset Tipset) (bool, ElectionProof) {
	lbt := ChainTipsMgr.TicketFromRound(parentTipset, n-RandomnessLookback)

	eproof := ComputeElectionProof(lbt, key)

	tipset := ChainTipsMgr.TipsetFromRound(n - PowerLookback)
	minerPower := GetPower(tipset.state, key.Public())
	totalPower := GetTotalPower(tipset.state)

	return IsProofAWinner(eproof, minerPower, totalPower), eproof
}
```


TODO: get accurate estimates for K and L, potentially merge both to a single param.

Note: Validity of blocks beyond appropriate ticket generation (defined below) is defined by the specific protocol using EC. For the Filecoin definition of a valid block, see the [mining spec](mining.md).

The basic algorithm can be broken down by looking in turn at its two major components: 

- Leader Election
- Chain Selection

## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the storage [power table](./storage-market.md#the-power-table), where power is equivalent to storage provided through time.

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's use of EC, the blockchain tracks an independent ticket chain. These tickets are used as randomness for Leader Election. Every block generated references an `ElectionProof` derived from a past ticket. The ticket chain is extended by the miner who generates a new ticket with each attempted election.

In cases in which no winning ticket is found by any miner in a given round (i.e. no block is mined on the network), miners move on to the next ticket in the ticket chain to attempt a new leader election. When this happens, miners should nonetheless generate a new ticket prior to the new leader election, thereby appropriately prolonging the ticket chain (the block chain can never be longer than the ticket chain). This situation is fleshed out in the [Losing Tickets](#losing-tickets) section.

In order to pressure the network to converge on a single chain, each miner may only submit one block per round (see: `Slashing`).

TODO: pictures of ticket chain and block chain


### Tickets

One may think of leader election in EC as a verifiable lottery, in which participants win in proportion to the power they have within the network. 

A ticket is drawn from the past at the beginning of each new round, and a new ticket is generated in every round. Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any other data in the block. Nonetheless, in Filecoin, every block header contains one or more new tickets, thereby extending the ticket chain. A miner generates a new ticket in their block for every ticket they scratch running leader election, thereby ensuring the ticket chain is at least as long as the block chain.

At a high-level, tickets must do the following:

- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.
- Prove appropriate delay between drawings — thereby preventing leaders from "rushing" the protocol by releasing blocks early (at the expense of fairness for miners with worse connectivity).
- Ensure a single drawing per round — derived in part from the above, thereby preventing miners from grinding on tickets (e.g. by repeatedly drawing new tickets in the hopes of winning) within a round.

For tickets in EC, one may use the following:
```go
type Ticket struct {
	// The VDF Result is derived from the prior ticket in the ticket chain
	VDFResult BigInteger
	// The VDF proves a delay between tickets generated
	VDFProof BigInteger
	// This signature is generated by the miner's keypair run on the VDFResult.
	// It is the value that will be used to generate future tickets or ElectionProofs.
	Signature Signature
}
```

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

But why the lookback?

The lookback helps turn independent lotteries (ticket drawings from a block one round back)
into a global lottery instead. Rather than having a distinct chance of winning or losing
for each potential fork in a given round, a miner will either win on all or lose on all
forks descended from the block in which the ticket is sampled.

This is useful as it reduces opportunities for grinding, across forks or sybil identities.

However this introduces a tradeoff:
- The lookback means that a miner can know K rounds in advance that they will win,
decreasing the cost of running a targeted attack (given they have local predictability).
- It means electionProofs are stored separately from new tickets on a block, taking up
more space on-chain.
```

#### Ticket generation

This section discusses how tickets are generated for the `Tickets` array.

At round `N`, new tickets are generated using tickets drawn from the [TipSet](#tipsets) at round `N-1`. This ensures the miner cannot publish a new block (corresponding to the `ElectionProof` generated by a winning ticket `K` rounds back) until the correct round.

Because a Tipset can contain multiple blocks (see [Chain Selection](#chain-selection) below), the smallest ticket in the Tipset must be drawn otherwise the block will be invalid.

TODO: pictures of TipSet ticket drawing

The miner runs the prior ticket through a Verifiable Delay Function (VDF) to get a new unique output. This approximates clock synchrony for miners, thereby ensuring miners have waited an appropriate delay ahead of drawing a new ticket. It also establishes a lower-bound delay between production of new blocks at different heights. This helps ensure fairness, in that miners with lesser connectivity are not penalized and have a chance to produce blocks.

This output is then used as input into a Verifiable Random Function (VRF) generating a new ticket, different from any other miners'. This adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance).

Succinctly, the process of crafting a new ticket in round `N` is as follows: 

```
old_ticket := sort(parentTickets)[0]
new_ticket := Sig(VDF(H(old_ticket)))

new_ticket: ticket to be used at round N in case of block generation
old_ticket: ticket drawn from round N-1
H: Cryptographic compression function
sort: bytewise sort
VDF: Verifiable Delay Function
Sig: Signature with the miner's keypair, used as a Verifiable Random Function.
```

#### Checking election results

Now, a miner must also check whether they have won in this round.

The process is as follows:

- At round N, the miner will draw the smallest ticket from the TipSet at round `N-K`, with `K` the randomness lookback parameter.
  - In the case where there are multiple tickets to choose from at round `N-K` (i.e. if the TipSet eventually created has multiple blocks), miners should attempt to generate their `ElectionProof` from the ticket generated by the block with the smallest final ticket (i.e. not necessarily the smallest ticket generated at that round).
- The miner will use this ticket as input to a VRF (Verifiable Random Function), thereby ensuring secrecy: no other participant can generate this output without the miner's private key. In that sense whether the miner has scratched a winning ticket remains a secret until they release their `ElectionProof` along with a new block.
- The miner will compare the output to their power fraction from round `N-L`, with `L` the committee lookback parameter. If it is smaller, they have won and can mine a block inserting this winning `ElectionProof` in the block header for the block produced at round `N`. Else they wait to hear of another block generated in this round.

At round N:

```
electionProof = Sig(H(SampleRandomness(CurrentRound - K)))

Sig: Signature with the miner's keypair, used as a VRF.
H: Cryptographic compression function
sort: bytewise sort
SampleRandomness: returns the appropriate ticket from round N-K
```

It is important to note that a miner generates two artifacts: one, a ticket derived from last block's ticket to prove that they have waited the appropriate delay, and two, an election proof derived from the ticket `K` rounds back used to run leader election.

Typically, either a miner will generate a winning ticket (see [Block Generation](#block-generation) or will hear about a new block (or multiple) by the end of a round (and start mining atop the smallest ticket of this new TipSet). The round may also have no successful miners.

### Losing Tickets

In the case that everybody draws a losing ticket in a given round (i.e. no miner is eligible to produce a block), every miner can run leader election again by "scratching" (attempting to generate a new `ElectionProof` from) the next ticket in the chain. That is, miners will now use the ticket sampled `K-1` rounds back to generate a new `ElectionProof`. They can then compare that proof with their power in the table `N-(L-1)` blocks back. This is repeated until a miner scratches a winning ticket and can publish a block (see [Block Generation](#block-generation)).

In addition to each attempted `ElectionProof` generation, the miner will need to extend the ticket chain by generating a new ticket. They use the ticket they generated in the prior round, rather than the prior block's (as is normally used). This proves appropriate delay (given that finding a winning Ticket has taken multiple rounds). Thus, each time it is discovered that nobody has won in a given round, every miner should use their previously generated ticket to repeat the ticket generation process, appending said ticket to their would-be block's `Ticket` array. This continues until some miner finds a winning ticket (see below), ensuring that the ticket chain remains at least as long as the block chain.

As was stated above, in the case where there are multiple tickets to choose from at round `N-K` (i.e. if the TipSet eventually created has multiple blocks), miners should attempt to generate their `ElectionProof` from the ticket generated by the block with the smallest final ticket (i.e. not necessarily the smallest ticket generated at that round). Put another way, the block in a TipSet with the smallest final ticket prolongs the valid ticket chain.

TODO: Add a diagram to illustrate this

The VDF ensures fairness by enforcing that miners cannot grind through repeated losing tickets (see more [here](https://github.com/filecoin-project/research/issues/31)) and that a miner cannot "rush" the protocol by outputting a block before others have had a chance to (e.g. geographically disadvantaged miners). The VDF delay is currently to 30 seconds, given estimated network propagation times.

Thus, our full ticket generation algorithm (reprised from [Ticket Generation](#ticket-generation)) is roughly (ticket handling is simplified in the pseudocode below for legibility):

```go
// Ticket is created as an array, with the initial ticket
// coming from the parent TipSet.
var Tickets []Ticket
oldTicket := sort(parentTickets)[0]
newTicket := VRF(VDF(H(oldTicket)))
electionProof := VRF(H(ticketFromRound(curRound - K)))

Tickets = append(Tickets, newTicket)

// If the current ticket isn't a winner and the block isn't found by another miner,
// derive a ticket from the last ticket
for !IsProofAWinner(electionProof) && !blockFound() {
	newTicket = VRF(VDF(H(newTicket)))
	newElectionProof = Sig(H(ticketFromRound(curRound - K)))
	Tickets = append(Tickets, newTicket)
	curRound += 1
}

// if the process yields a winning ticket, mine and put out a block
// containing the ticket array
if winning(electionProof) {
	mineBlock(electionProof, Tickets)
}
```

A ticket can be verified to have been generated in the appropriate number of rounds by looking at the `Tickets` array, and ensuring that each subsequent ticket (leading to the final ticket in that array) was generated using the previous one in the array. Note that this has implications on block size, and client memory requirements, though on expectation, the `Ticket` array should only contain one value.

In fact, the length of repeated losing tickets in the ticket chain (equivalent to the length of generated tickets referenced by a single block, or the length of the `Tickets` array) decreases exponentially in the number of repeated losing tickets (see more [here](https://github.com/filecoin-project/go-filecoin/pull/1516)). In the unlikely case the number of losing tickets drawn by miners grows larger than the randomness lookback `K` (i.e. if a miner runs out of existing tickets on the ticket chain for use as randomness), a miner should proceed as usual using new tickets generated in this epoch for randomness. This has no impact on the protocol safety/validity.

New blocks (with multiple tickets) will have a few key properties:

- All tickets in the `Tickets` array are signed by the same miner -- to avoid grinding through out-of-band collusion between miners exchanging tickets.
- The `ElectionProof` was correctly generated from the ticket `K-|Tickets|-1` (with `|Tickets|` the length of the `Tickets` array) rounds back.

This means that valid `ElectionProof`s can be generated from tickets in the middle of the `Tickets` array.

### Block Generation

Once a miner has a winning election proof generated over `i` rounds and `i` corresponding tickets, they may create a block. For more on this, see the [Mining spec](./mining.md#block-creation).

## Chain Selection

Just as there can have 0 miners win in a round, multiple miners can be elected in a given round . This in turn means multiple blocks can be created in a round. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round

### Tipsets

All valid blocks generated in a round form a `TipSet` that participants will attempt to mine off of in the subsequent round (see above). TipSets are valid so long as:

- All blocks in a TipSet have the same parent TipSet
- All blocks in a TipSet have the same number of tickets in their `Tickets` array

The first condition implies that all blocks in a TipSet were mined at the same height (remember that height refers to block height as opposed to ticket round). This rule is key to helping ensure that EC converges over time. While multiple new blocks can be mined in a round, subsequent blocks all mine off of a TipSet bringing these blocks together. The second rule means blocks in a TipSet are mined in a same round.

Due to network propagation delay, it is possible for a miner in round N+1 to omit valid blocks mined at round N from their TipSet. This does not make the newly generated block invalid, it does however reduce its weight and chances of being part of the canonical chain in the protocol.

### Chain Weighting

TODO: ensure 'power' is properly and clearly defined

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided.

The weight at each block is equal to its `ParentWeight`, plus that block's delta weight. Delta
weight is a constant `V`, plus `X` - a function of the total power in the network as reported in the Power Table.  The exact value for `V` and the magnitude of the power ratio value are
still to be determined, but for now `V = 10` and `X = log(TotalPower)`.

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given TipSet should have 
the same `ParentWeight` value) plus the delta weight of each parent. To make the 
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise 
potentially long chain scans would be required to compute a given block's weight).

When selecting between TipSets of equal weight, a miner chooses the one with the smallest min ticket (by bytewise comparison).

### Slashing

See the [Faults spec](./faults.md) for implementation details.

Due to the existence of potential forks, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This is detectable when a miner submits two blocks that satisfy either of the following "slashing conditions":

(1) one block contains at least one ticket in its ticket array generated at the same round as one of the tickets in the other block's ticket array.
(2) one block's parent is a TipSet that could have validly included the other block according to TipSet validity rules, however the parent of the first block does not include the other block.

  - While it cannot be proven that a miner omits known blocks from a TipSet in general (i.e. network latency could simply mean the miner did not receive a particular block) in this case it can be proven because a miner must be aware of a block they mined in a previous round.

Any node that detects this occurring should take both block headers, and call [`storagemarket.SlashConsensusFault`](actors.md#slashconsensusfault). The network will then take all of that node's collateral, give a portion of it to
the reporter, and keep the rest.

TODO: It is unclear that rewarding the reporter any more than gas fees is the right thing to do. Needs thought. Tracking issue: https://github.com/filecoin-project/specs/issues/159

Note: One may wonder what prevents miners from simply breaking up their power into multiple un-linkable miner actors  (or sybils) that will be able to mine on multiple chains without being caught mining at the same round at the same time. We call this the "whyru slashing" attack.

TODO: Discuss with @zenground0 to ensure the below should not be removed

```
An attacker with 30% of the power has a 30% chance of mining a block at any given moment.
During a fork, an attacker would have a 30% chance of winning _on each fork_, meaning
they could continue to mine on both forks, except for 30% of the time they win, they
will win on both chains and have to forgo publishing one of the blocks to avoid being
slashed. This means that the miner loses 30% of their expected rewards, but is still
able to mine on both chains. 

Now, if the miner instead controls two miners each with 15% of the total power in
the network, they still have a 30% chance of winning on each fork (using both sybils),
but they drop the probability of mining with the same miner at the same time on both
chains down to 15% of their winnings (meaning they have to forgo 15% of their
‘successfully’ mined blocks). This continues on down, 3 identities is 10%, 4 is 7.6%,
5 is 6% and so on. So at the end of the day, the miner can mine on both chains with
only a minimal loss in potential proceeds.

The above assumes that every election is an independent random process (even across forks).
However, using a lookback parameter for seed sampling, the independent lottery drawing
becomes a global lottery for all forks originating after the lookback (where the
randomness was drawn). That is to say, given a common random seed and public key, each
sybil will either win on all forks, or lose on all forks. This greatly decreases the
chances that this attack succeeds, erasing the economic advantage sybils created.
```

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

- Parameter K, Parameter L
- Checkpointing Strategy
- Block confirmation time
- When selecting between two forks of equal weight, one strategy might be to select the 'Tipset' with the lowest number of linked tickets for a given block height and weight.
- Should there be a minimum power required to participate in the consensus process?
- How long should 'valid' candidate blocks be kept around? Essentially the question is: when is finality?
- How should block rewards be assigned in the expected consensus setting?
- VDF difficulty adjustment
