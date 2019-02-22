# Expected Consensus

This spec describes how to implement the protocol in general, for Filecoin-specific processes, see:

- [Mining Blocks](mining.md#mining-blocks) on how consensus is used.
- [Faults](./faults.md) on slashing.
- [Storage Market](./storage-market.md#the-power-table) on how the power table is maintained.

Expected Consensus is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one 
participant wins. The important part to note is "on expectation". In each round, any number of 
participants, or none, may find a winning ticket, and any miner holding a winning ticket may 
submit a block. All valid blocks submitted in a given round for a `TipSet`. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain.


The basic algorithm is as follows:

- Collect all incoming blocks.
- Check their validity as shown in [the mining spec.](./mining.md#block-validation)
- For each block that is valid, and not known to be based on a bad chain, place it into a `ChainTipsManager.`
  - Each block should be indexed by height, and by its set of parents.
  - A set of blocks with the same height, same parent set, and same number of null tickets is called a `TipSet`
- To check if you are a winner for round N:
  - Sample the secure chain randomness at block `N-K` to receive a ticket (see [Tickets](#tickets) for details). `K` is called the randomness lookback parameter, currently set to 1.
  - If the hash of the ticket is less than the ratio of your power over the total power in the network at the end of round N-L, you have a winning ticket. L is called the committee lookback parameter, currently set to 1.
- If you are a winner, generate a new block:
  - Select the smallest ticket from the heaviest `TipSet` at round N-1.
  - Generate a new ticket from it for inclusion in your new block.
  - Generate a new block as shown in [the mining spec](./mining.md#block-creation).

TODO: get accurate estimates for K and L, potentially merge both to a single param.

Note: Validity of blocks beyond appropriate ticket generation (defined below) is defined by the implementation. For the filecoin definition of a valid block, see the [mining spec](mining.md).

Expected consensus relies on weighted chains in order to quickly converge on 'one true chain'. 

We can describe the basic algorithm by looking in turn at its two major components: 

- Leader Election
- Chain Selection

## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their weight. In the case of Filecoin, participants and weights are drawn from the storage [power table](./storage-market.md#the-power-table).

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness drawn from the chain.

In order to pressure the network to converge on a single chain, each miner may only submit one block per round (see: `Slashing`). In cases in which no block is mined for a round, a block with multiple tickets in each array may be generated in a subsequent round. We deal with this case in the `Null Blocks` section.

### Tickets

We think of leader election in EC as a verifiable lottery, in which participants win in proportion to the storage they provide to the network. 

Tickets are drawn at the beginning of a new epoch. In EC, there is one round of leader election per epoch (i.e. a new ticket is drawn for each leader election). Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any data in the block. 

At a high-level, tickets must do the following:

- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.
- Prove appropriate delay between drawings — thereby preventing leaders from "rushing" the protocol by releasing blocks early (at the expense of fairness for miners with worse connectivity).
- Ensure a single drawing per round — derived in part from the above, thereby preventing miners from grinding on tickets (e.g. by repeatedly mining null blocks) within a round.

In practice, miners make use of tickets in two main places within a block:

- A `Tickets` array — this stores new tickets generated using the ticket 1 block back and proves appropriate delay. It is from this array that miners will sample randomness to run leader election in `K` rounds. We discuss generation in `Ticket generation`.
- An `ElectionProof` — this stores the winning lottery ticket generated using the smallest ticket (from the `Tickets` array) in the parent `Tipset` `K` round back. It proves that the leader was elected in this round. We discuss generation in `Checking election results`.

On expectation, the `Tickets` array will contain a single ticket. We discuss the case in which it contains more than one below (see [Null Blocks](#null-blocks)).

```
Lookback parameter and why tickets are stored in two different places in a block?

We use tickets for two main purposes: as a proof that blocks were appropriately delayed (both to prevent grinding in a round and a winner rushing the protocol), and as a proof that a leader was correctly elected. 

The proof of delay requires a ticket to be sampled from one block back, this is how tickets are generated for future use.

The proof of election however requires a ticket to be sampled from K rounds back. By sampling randomness from K rounds back, we turn independent events (each drawing from a block one round back) into a global lottery instead. That is that rather than having a distinct chance of winning or losing for each potential fork in a given round, a miner will either win on all or lose on all forks descended from the block in which randomness is sampled. This is useful as it reduces opportunities for grinding, across forks or across sybil identities.

This is a tradeoff however for two reasons:
- The lookback means that a miner can know K rounds in advance that they will win, decreasing the cost of running a targeted attack (given they have local predictability).
- It means we must store electionProofs separately from "DelayProofs" (i.e. Tickets) in a block, taking up more space on-chain.
```

#### Ticket generation

This section discusses how Tickets are generated for the `Tickets` array.

At round `N`, new tickets are generated using tickets drawn from the [Tipset](#tipsets) at round `N-1`. This ensures the miner cannot start the process of checking for a winning ticket until the correct round.

Because a Tipset can contain multiple blocks (see [Chain Selection](#chain-selection) below), the smallest ticket in the Tipset must be drawn otherwise the block will be invalid.

The miner runs the prior ticket through a Verifiable Delay Function (VDF) to get a new unique output. This approximates clock synchrony for miners, thereby ensuring miners have waited for an appropriate delay ahead of drawing a new ticket. It also ensures that miners will wait some amount of time before producing a new block, thereby ensuring miners with lesser connectivity are not penalized (e.g. by colluding miners rushing the protocol).

This output is then used as input into a Verifiable Random Function (VRF) generating a new ticket, different from any other miners'. This limits a miner's ability to alter one block to influence the ticket of the next block (given a miner does not know who will win a given round in advance).

Simply stated, the process of crafting a new ticket in round N is as follows: 

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

- At round N, the miner will draw the smallest ticket from the Tipset at round N-K, with K the randomness lookback parameter.
- The miner will use this ticket as input to a VRF (Verifiable Random Function), thereby ensuring secrecy: no other participant can generate this output without the miner's private key.
- The miner will compare the output to their power fraction from round N-L, with L the committee lookback parameter. If it is smaller, they have won and can mine a block inserting this winning `ElectionProof` in the block header for the block produced at round N. Else they wait to hear of another block generated in this round.

We have at round N:

```
electionProof = Sig(H(sort(lookbackTickets)[0]))

Sig: Signature with the miner's keypair, used as a VRF.
H: Cryptographic compression function
sort: bytewise sort
lookbackTickets: Winning tickets from each of the `Tipset` blocks at round N-K
```

It is important to note that a miner generates two artifacts: one, a ticket from derived from last block's ticket to prove that they have waited the appropriate delay, and two, an election proof from K blocks back (where K could be 1) to run leader election.

Typically, either a miner will generate a winning ticket (see [Block Generation](#block-generation) or will hear about a new block (or multiple) by the end of a round (and start mining atop the smallest ticket of this new tipset). The round may also have no successful miners.

### Null Blocks

In the case that nobody wins a ticket in a given round, a `Null block` may be inserted. A 'null block' doesn't actually exist as its own separate block, it simply serves as a way to signal that a number of rounds were skipped.

To 'mine' a null block, you will need to generate a new ticket to prove delay (and the fact that winning has taken multiple rounds), and a new election proof.

**Generating new tickets to prove delay**

Let's start with delay: take the new ticket generated by your initial mining attempt (which came from the previous `Tipset`), and use it as the 'parent ticket' to a new ticket generation process. Each time it is discovered that nobody has won a given round, every miner should use their previously generated ticket to repeat the ticket generation process, appending said ticket to their would-be block's `Ticket` array. This continues until some miner finds a winning ticket (see below), with the VDF ensuring that miners cannot grind through repeated null block generation (see more [here](https://github.com/filecoin-project/research/issues/31)).

Delay also ensures that a certain miner cannot "rush" the protocol by outputting a block before others have had a chance to (say geographically disadvantaged miners) but rather creates some fairness in the protocol. We currently set this delay to 30 seconds, given estimated network propagation times.

Thus, our full ticket generation algorithm (reprised from `Ticket`) is roughly:

```
// Ticket is created as an array, with the initial ticket
// coming from the parent tipset.
var Tickets []Signature
oldTicket := sort(parentTickets)[0]
newTicket := VRF(VDF(H(oldTicket)))
electionProof := VRF(H(ticketFromHeight(curHeight-K)))

Tickets = append(Tickets, newTicket)

// If the current ticket isn't a winner and the block isn't found by another miner,
// derive a ticket from the last losing ticket
for !winning(electionProof) && !blockFound()) {
 newTicket = VRF(VDF(H(newTicket)))
 newElectionProof = Sig(H(ticketFromHeight(curHeight+len(Tickets)-K)))
 Tickets = append(Tickets, newTicket)
}

// if the process yields a winning ticket, mine and put out a block
// containing the ticket array
if winning(electionProof) {
    mineBlock(electionProof, Tickets)
}
```

This ticket can be verified to have been generated in the appropriate number of rounds by looking at the 'Tickets' array, and ensuring that each subsequent ticket (leading to the final ticket in that block) was generated using the previous one in the array. Note that this has implications on block size, and client memory requirements, though on expectation, blocks should rarely grow too much because of this.

**Checking election results (after mining a null block)**

Failing to generate an election proof, a miner should discard their failed proof from the original mining attempt (drawn from `K` rounds back), and recompute a leader election proof at the next block. That is, the miner will now use the ticket sampled `K-1` rounds back to generate an election proof. They can then compare that proof with their power in the table N-L-1 blocks back. Each time it is discovered that nobody has won a given round, every miner should repeat the leader election process using the next ticket in the chain to generate a new ElectionProof. Once a miner finds a winning ticket, they can publish a block (see `Block Generation`).

This new block (with multiple tickets) will have a few key properties:

- All tickets in are signed by the same miner -- to avoid grinding through out-of-band collusion between miners exchanging tickets.
- The election proof was correctly generated from K rounds back, counting null blocks.

```
On looking back

Here are a few examples to help illustrate how the randomness lookback parameter, 'K', functions.

At a given height 'N' (or for round N, which is functionally equivalent), sample your ticket from N-K rounds to generate your election proof.
- This could mean that you are looking back J <= K actual blocks back: since we count null blocks, we may walk back multiple blocks (null blocks + actual block) within a single header.
- This could mean you are sampling a null block ticket (i.e. some element in a 'Tickets' array) to generate your Election Proof.
```

### Block Generation

When you have a winning election proof and corresponding ticket, you may create a block. For more on this, see the [Mining spec](./mining.md#block-creation).

## Chain Selection

As we saw, just as we can have no miners win in a round (leading to null block mining), multiple miners can be elected in a given round . This in turn means multiple blocks can be created in a round. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round

### Tipsets

All valid blocks generated in a round form a 'Tipset' that participants will attempt to mine off of in the subsequent round (see above). Tipsets are valid so long as:

- All blocks in a Tipset must have the same parents

This in turn implies that all blocks in a Tipset were mined at the same height. This rule is key to helping ensure that EC converges over time. While multiple new blocks can be mined in a round, subsequent blocks all mine off of a tipset bringing these blocks together.

Due to network propagation delay, it is possible for a miner in round N+1 to omit valid blocks mined at round N from their Tipset. This does not make the newly generated block invalid, it does however reduce its chances of being part of the canonical chain in the protocol.

### Chain Weighting

As we saw, it is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided. 

The weight at each block is equal to its `ParentWeight`, plus that block's delta weight. Delta
weight is a constant `V`, plus `X` - a function of the total power in the network as reported in the Power Table.  The exact value for `V` and the magnitude of the power ratio value are
still to be determined, but for now we can use `V = 10` and `X = log(TotalPower)`.

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given `TipSet` should have 
the same `ParentWeight` value) plus the delta weight of each parent. To make the 
computation a bit easier, we store a block's parent weight in the block itself (otherwise 
potentially long chain scans would be required to compute a given block's weight).

When selecting between tipsets of equal weight, choose the one with the smallest ticket (by bytewise comparison).

### Slashing

See the [Faults spec](./faults.md) for implementation details.

Due to the existence of potential forks, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This happens in the following instances: 

- A miner submits two blocks for the same round.
- A miner mines atop a 'Tipset' that should have contained their own block (same height, same parents) but does not.
  - While we cannot prove that block omission from a 'Tipset' is malicious (ie not due to network latency), a miner may omit their own block if they find a winning Ticket in the parent 'Tipset' though they themselves have submitted a smaller (losing) Ticket. For obvious reasons, we can assert they are aware of the block they previously mined.

Any node that detects this occurring should take both block headers, and submit them to the
network slasher. The network will then take all of that node's collateral, give a portion of it to
the reporter, and keep the rest.

TODO: It is unclear that rewarding the reporter any more than gas fees is the right thing to do. Needs thought.

Note: You may wonder what prevents miners from simply breaking up their power into multiple un-linkable miner actors  (or sybils) that will be able to mine on multiple chains without being caught mining at the same height at the same time. We call this the "why r u slashing" attack.

```
An attacker with 30% of the power has a 30% chance of mining a block at any given moment. During a fork, an attacker would have a 30% chance of winning _on each fork_, meaning they could continue to mine on both forks, except for 30% of the time they win, they will win on both chains and have to forgo publishing one of the blocks to avoid being slashed. This means that the miner loses 30% of their expected rewards, but is still able to mine on both chains. 

Now, if the miner instead controls two miners each with 15% of the total power in the network, they still have a 30% chance of winning on each fork (using both sybils), but they drop the probability of mining with the same miner at the same time on both chains down to 15% of their winnings (meaning they have to forgo 15% of their ‘successfully’ mined blocks). This continues on down, 3 identities is 10%, 4 is 7.6%, 5 is 6% and so on. So at the end of the day, the miner can mine on both chains with only a minimal loss in potential proceeds.

The above assumes that every election is an independent random process (even across forks). However, using a lookback parameter for seed sampling, the independent lottery drawing becomes a global lottery for all forks originating after the lookback (where the randomness was drawn). That is to say, given a common random seed and public key, each sybil will either win on all forks, or lose on all forks. This greatly decreases the chances that this attack succeeds, erasing the economic advantage sybils created.
```



## Implementation Notes

- When selecting messages from the mempool to include in your block, be aware that other miners may also generate blocks during this round, and if you want to maximize your fee earnings it may be best to select some messages at random (second in a duplicate earns no fees).

## Open Questions

- Parameter K, Parameter L
- Checkpointing Strategy
- Block confirmation time
- When selecting between two forks of equal weight, one strategy might be to select the 'Tipset' with the lowest number of null tickets for a given block height and weight.
- Should there be a minimum power required to participate in the consensus process?
- How long should we keep 'valid' candidate blocks around? Essentially the question is: when is finality?
- How should we assign block rewards in the expected consensus setting?
