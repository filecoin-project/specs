# Expected Consensus

This spec describes how to implement the protocol in general, for Filecoin-specific processes, see:

- [Mining Blocks](mining.md#mining-blocks)
- [Faults](./faults.md)

Expected Consensus is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one 
participant wins. The important part to note is "on expectation". In each round, any number of 
participants, or none, may find a winning ticket, and any miner holding a winning ticket may 
submit a block. All valid blocks submitted in a given round for a TipSet. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain.


The basic algorithm is as follows:

- Collect all valid incoming blocks
- For each block that is valid, and not known to be based on a bad chain, place it into a `ChainTipsManager`
  - Each block should be indexed by height, and by its set of parents.
  - A set of blocks with the same height, same parent set, and same number of null tickets is called a `TipSet`
- To check if you are a winner for round N:
  - Select the heaviest `TipSet` for block N-1
  - Sample the secure chain randomness at block N-K to receive a ticket (see "Tickets" for details). K is called the randomness lookback parameter, currently set to 1.
  - If the hash of the ticket is less than the ratio of your power over the total power in the network at the end of block N-L, you have a winning ticket. L is called the committee lookback parameter, currently set to 1.

Note: Validity of blocks is defined by the implementation. For the filecoin definition of a valid block, see the [mining document](mining.md).

We can describe the basic algorithm by looking in turn at its two major components: 

- Leader Election
- Chain Selection

## Secret Leader Election

Expected Consensus is a Nakamoto-style consensus protocol that works by electing a miner from a weighted set in proportion to their weight. In the case of Filecoin, participants and weights are drawn from the storage power table. 

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of tickets drawn from the chain.

### Tickets

Tickets are how miners know if they are the leader of a given round. They act as a common and verifiable source of randomness miners can use to run leader election at every round and prove to other participants that they have indeed won.

Tickets are drawn at the beginning of a new epoch. In EC, there is one round of leader election per epoch (i.e. a new ticket is drawn for each leader election). Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any data in the block. 

Tickets must do the following:

- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.
- Prove appropriate delay between drawings -- meaning a participant can only draw a single ticket per round.

Each block stores its ticket in two places: an array `Tickets` which stores new tickets generated at every round for use as randomness `K` blocks in the future, and in an array `ElectionProofs` to prove that a miner has indeed been elected miner using a `Ticket` from `K` blocks back . On expectation, the arrays in a block will each contain a single ticket. In order to pressure the network to converge on a single chain, each miner may only submit one block per round (see: `Slashing`). In cases in which no block is mined for a round, a block with multiple tickets in each array may be generated in a subsequent round. We deal with this case in the `Null Blocks` section.

#### Ticket generation and delay

To generate a new ticket, take a ticket from the prior `Tipset`. Because a `Tipset` can contain multiple blocks (see `Chain Selection` below), the smallest ticket in the Tipset must be drawn (otherwise the block will be invalid). This ensures the miner cannot start the process of checking for a winning ticket until the correct round.

The miner then runs the prior ticket through a Verifiable Delay Function (VDF) to get a new unique output. This allows us to approximate clock synchrony for miners, thereby ensuring miners have waited for an appropriate delay ahead of drawing a new ticket.

This output is then used as input into a Verifiable Random Function (VRF) generating a new ticket, different from any other miners'. This limits a miner's ability to alter one block to influence the ticket of the next block (given a miner does not know who will win in a given round in advance).

Simply stated, the process of crafting a new ticket in round N is as follows: 

```
old_ticket := sort(parentTickets)[0]
new_ticket := VRF(VDF(H(old_ticket)))

new_ticket: ticket to be used at round N in case of block generation
old_ticket: ticket drawn from round N-1
H: Cryptographic compression function
sort: bytewise sort
VDF: Verifiable Delay Function
VRF: Verifiable Random Function
```

This ensures blocks are created after an appropriate delay, or put another way that miners only draw tickets once per round.

#### Checking election results

Now, a miner must also check whether they have won in this round.

The process is as follows:

- At block N, the miner will draw the smalelst ticket from the `Tipset` at block N-K, with K the randomness lookback parameter.
- The miner will use this ticket as input to a Verifiable Random Function, thereby ensuring secrecy: no other participant can generate this output without the miner's private key.
- The miner will compare the output to their power from block N-L, with L the committee lookback parameter. If it is smaller, they have won and can mine a block inserting this winning ticket in the `ElectionProofs` array in the block. Else they wait to hear of another block generated in this round.

We have at round N:

```
T = Sig(H(sort(lookbackTickets)[0]))

Sig: Signature with the miner's keypair
H: Cryptographic compression function
sort: bytewise sort
lookbackTickets: Winning tickets from each of the `Tipset` blocks at block N-K
```

It is important to note that a miner uses two tickets: one from last block to prove that they have waited the appropriate delay, and one from K blocks back (where K could be 1) to run leader election.

Typically, either a miner will generate a winning ticket (see `Block Generation`) or will hear about a new block (or many) by the end of a round (and start mining atop the smallest ticket of this new tipset). But the round may also have no successful miners.

### Null Blocks

In the case that nobody wins a ticket in a given round, a 'null block' may be inserted. A 'null block' doesn't technically need to exist as its own separate block, it simply serves as a way to signal that a number of rounds were skipped.

To 'mine' a null block, you will need to generate two new tickets, one proving appropriate delay (and the fact that winning has taken multiple rounds) to be used for leader election in a future round and another proving your election.

**Generating new tickets to prove delay**

Let's start with delay: take the new ticket generated by your initial mining attempt (which came from the previous `Tipset`), and use it as the 'parent ticket' to a new ticket generation process. Each time it is discovered that nobody has won a given round, every miner should use their previously generated ticket to repeat the ticket generation process, appending said ticket to their would-be block's `Ticket` array. This continues until some miner finds a winning ticket (see below), with the VDF ensuring that miners cannot grind through repeated null block generation (see more [here](https://github.com/filecoin-project/research/issues/31)).

Delay also ensures that a certain miner cannot "rush" the protocol by outputting a block before others have had a chance to (say geographically disadvantaged miners) but rather creates some fairness in the protocol. We currently set this delay to 30 seconds, given estimated network propagation times.

Thus, our full ticket generation algorithm (reprised from `Ticket`) is roughly:

```
// Ticket is created as an array, with the initial ticket
// coming from the parent tipset.
var Tickets []Signature
old_ticket := sort(parentTickets)[0]
new_ticket := VRF(VDF(H(old_ticket)))

Tickets = append(Tickets, new_ticket)

// If the current ticket isn't a winner and the block isn't found by another miner,
// derive a ticket from the last losing ticket
for !winning(new_ticket) && !blockFound()) {
 new_ticket = VRF(VDF(H(new_ticket)))
 Tickets = append(Tickets, new_ticket)
}

// if the process yields a winning ticket, mine and put out a block
// containing the ticket array
if winning(T) {
    mineBlock(ticket)
}
```

This ticket can be verified to have been generated in the appropriate number of rounds by looking at the tickets array, and ensuring that each subsequent ticket (leading to the winning ticket) was generated using the previous one in the array. Note that this has implications on block size, and client memory requirements, though on expectation, blocks should rarely grow too much because of this.

**Checking election results again**

Likewise, a miner should take their losing ticket from the original mining attempt (drawn from `K` blocks back), add it to the `ElectionProofs` array in the block and run a VRF on it once more, generating a new ticket to compare with their power in the table N-L blocks back. Each time it is discovered that nobody has won a given round, every miner should use their failed ticket to repeat the leader election process, appending said ticket to their would-be block's `ElectionProofs array. Once a miner finds a winning ticket, they can publish a block (see `Block Generation`).

This new block (with multiple tickets in each array) will have a few key properties:

- All tickets in each array are signed by the same miner -- to avoid grinding through out-of-band collusion between miners exchanging tickets.
- There is the same number of tickets in both arrays -- to ensure a miner has waited appropriately (in the `Ticket` array) before using their winning ticket (in the `ElectionProofs` array).

Thus, our ticket validation algorithm checks that the last ticket in the `ElectionProofs` array is a winning ticket and was adequately generated either from the parent set or from previous failed tickets.

### Block Generation

When you have found a winning ticket, you may create a block. For more on this, see the [Mining spec](./mining.md).

## Chain Selection

As we saw, just as we can have no miners win in a round (leading to null block mining), multiple miners can be elected in a given round . This in turn means multiple blocks can be created in a round. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round

### Tipsets

All valid blocks generated in a round form a `Tipset` that participants will attempt to mine off of in the subsequent round (see above). Tipsets are valid so long as:

- All blocks in a Tipset must have the same parents

This in turn implies that all blocks in a Tipset were mined at the same height. This rule is key to helping ensure that EC converges over time. While multiple new blocks can be mined in a round, subsequent blocks all mine off of a tipset bringing these blocks together.

Due to network propagation delay, it is possible for a miner in round r+1 to omit valid blocks mined at round r from their Tipset. This does not make the newly generated block invalid, it does however reduce its chances of being part of the canonical chain in the protocol.

### Chain Weighting

As we saw, it is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided. 

The weight at each block is equal to its `ParentWeight`, plus that block's delta weight. Delta
weight is a constant `V`, plus the ratio of the total power in the network controlled by the
miner of the block.  The exact value for `V` and the magnitude of the power ratio value are
still to be determined, but for now we can use `V = 10` and `(100 * MinerPower) / TotalPower`.

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

- When implementing parent selection, for now, always choose the single heaviest parent set for building on top of. More advanced strategies can be implemented later.
- When selecting messages from the mempool to include in your block, be aware that other miners may also generate blocks during this round, and if you want to maximize your fee earnings it may be best to select some messages at random (second in a duplicate earns no fees).

## Open Questions

- Parameter K, Parameter L
- Checkpointing Strategy
- Block confirmation time
- When selecting between two forks of equal weight, one strategy might be to select the 'Tipset' with the lowest number of null tickets for a given block height and weight.