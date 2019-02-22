# Expected Consensus

This is a technical design document describing how to implement the protocol in general, for the Filecoin specific mining process that implements Expected Consensus see [Mining Blocks](mining.md#mining-blocks).

Expected Consensus is a probabilistic Byzantine fault-tolerant consensus protocol. At a high
level, it operates by running a leader election every round in which, on expectation, one 
participant wins. The important part to note is "on expectation". In each round, any number of 
participants, or none, may find a winning ticket, and any miner holding a winning ticket may 
submit a block. The 'best' chain is the one with the highest weight, which is to say that the fork 
choice rule is to choose the heaviest known chain.


The basic algorithm is as follows:

- Collect all valid incoming blocks
- For each block that is valid, and not known to be based on a bad chain, place it into a ChainTipsManager
  - Each block should be indexed by height, and by its set of parents.
  - A set of blocks with the same height, same parent set, and same number of null tickets is called a `TipSet`
- To check if you are a winner for round N:
  - Select the heaviest `TipSet` for round N-1
  - Sample the secure chain randomness to receive a ticket (see "Tickets" for details)
  - If the hash of the ticket is less than the ratio of your power over the total power in the network at the end of round N-K (where K is a lookback parameter, ideally behind finality, currently set to 1), you have a winning ticket.

Note: Validity of blocks is defined by the implementation. For the filecoin definition of a valid block, see the [mining document](mining.md).

## Chain Weighing

Expected consensus relies on weighted chains in order to quickly converge on 'one true chain'. 
The weight at each block is equal to its `ParentWeight`, plus that block's delta weight. Delta
weight is a constant `V`, plus `X` - a function of the total power in the network as reported in the Power Table.  The exact value for `V` and the magnitude of the power ratio value are
still to be determined, but for now we can use `V = 10` and `X = log(TotalPower)`.

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given `TipSet` should have 
the same `ParentWeight` value) plus the delta weight of each parent. To make the 
computation a bit easier, we store a block's parent weight in the block itself (otherwise 
potentially long chain scans would be required to compute a given block's weight).

When selecting between tipsets of equal weight, choose the one with the smallest ticket (by bytewise comparison).

## Tickets

Tickets are how miners know if they are the leader of a given round. Each new ticket sampling marks the beginning of a new epoch. In this case, there is one round per epoch. Tickets are
chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any data in the block.

Each block stores its ticket in an array `Tickets`. On expectation, the array in a block will contain a single ticket. In cases in which no block is mined for a round, a block with multiple tickets in the array may be generated in a subsequent round. We deal with this case in the `Null Blocks` section. To compute a ticket, use the following:

```
T = Sig(H(sort(parentTickets)[0]))

Sig: Signature with the miners keypair
H: Cryptographic compression function
sort: bytewise sort
```

In English, that is the signature (using the miner's keypair) over the hash of the smallest ticket in the parent set. Note that this signature must be deterministic.

The advantages of this method of leader selection are that the miner cannot start the process of checking 
for a winning ticket until the correct round, and the miner cannot alter one block to influence the ticket of the next block. The downside is that it introduces some level of concurrency in the process.

In order to pressure the network to converge on a single chain, each
miner may only submit one block per round. Any miner caught submitting multiple winning tickets for a 
given round may be reported and slashed (see: Slashing).

## Block Generation

When you have found a winning ticket, you may then create a block. To create a block, first compute a few fields:

- `Tickets` - An array containing your winning ticket, and, if applicable, the failed intermediary tickets, or `NullTickets` for any null blocks you mined on
- `ParentWeight` - As described above in "Chain Weighing"
- `ParentState` - To compute this:
  -  Take the `ParentState` of one of the blocks in your chosen parent set (invariant: this is the same value for all blocks in a given parent set)
  -  For each block in the parent set, ordered by their tickets:
    -   Apply each message in the block to the parent state, in order. If a message was already applied in a previous block, skip it.
    - Transaction fees are given to the miner of the block that the first occurance of the message is included in. If there are two blocks in the parent set, and they both contain the exact same set of messages, the second one will receive no fees.
    - It is valid for messages in two different blocks of the parent set to conflict, that is, A conflicting message from the combined set of messages will always error.  Regardless of conflicts all messages are applied to the state.
- `Messages` - Select a set of messages from the mempool to include in your block
- `StateRoot` - Apply each of your chosen messages to the `ParentState` to get this
- `Signature` - A signature with your private key (must also match the ticket signature) over the entire block. This is to ensure that nobody tampers with the block after we propogate it to the network, since unlike normal PoW blockchains, a winning ticket is found independently of block generation.

## Null Blocks

Typically, either a miner will generate a winning ticket or will hear about a new block (or many) by the end of a round (and start mining atop the smallest ticket of this new tipset). In the case that nobody wins a ticket in a given round, a 'null block' may be inserted. A 'null block' doesn't technically need to exist as its own separate block, it simply serves as a way to signal that a number of rounds were skipped.

To 'mine' a null block, take the failed ticket from your initial mining attempt, and use it as the 'parent ticket' to a new ticket generation process. Each time it is discovered that nobody has won a given round, every miner should use their failed ticket to repeat the ticket checking process, appending said ticket to their would-be block's ticket array. This continues until some miner finds a winning ticket, ensuring that miners cannot grind through repeated null block generation (see more [here](https://github.com/filecoin-project/research/issues/31)).

Thus, our full ticket generation algorithm (reprised from `Ticket`) is roughly:

```
// Ticket is created as an array, with the initial ticket
// coming from the parent tipset.
var Tickets []Signature
T := Sig(H(sort(parentTickets)[0]))
Tickets = append(Tickets, T)

// If the current ticket isn't a winner and the block isn't found by another miner,
// derive a ticket from the last losing ticket
for !winning(T) && !blockFound()) {
 T = Sig(H(T))
 Tickets = append(Tickets, T)
}

// if the process yields a winning ticket, mine and put out a block
// containing the ticket array
if winning(T) {
    mineBlock(ticket)
}
```

This ticket can be verified to have been generated in the appropriate number of rounds by looking at the tickets array, and ensuring that each subsequent ticket (leading to the winning ticket) was generated using the previous one in the array. Note that this has implications on block size, and client memory requirements, though on expectation, blocks should rarely grow too much because of this.

Verification of a block should also ensure that all tickets in the array were signed by the same miner (to avoid grinding through out-of-band collusion between miners exchanging tickets).

Thus, our ticket validation algorithm checks that the last ticket is a winning ticket and was adequately generated either from the parent set or from previous failed tickets.

## Slashing

A miner should be slashed if they are provably deviating from the honest protocol.

This happens in the following instances: 

- A miner submits two blocks for the same round.
- A miner mines atop a 'Tipset' that should have contained their own block (same height, same parents) but does not.
  - While we cannot prove that block omission from a 'Tipset' is malicious (ie not due to network latency), a miner may omit their own block if they find a winning Ticket in the parent 'Tipset' though they themselves have submitted a smaller (losing) Ticket. For obvious reasons, we can assert they are aware of the block they previously mined.

Any node that detects this occurring should take both block headers, and submit them to the
network slasher. The network will then take all of that node's collateral, give a portion of it to
the reporter, and keep the rest.

TODO: It is unclear that rewarding the reporter any more than gas fees is the right thing to do. Needs thought.

TODO: elaborate

Note(why): This slashing is insufficient to protect against most attacks. Miners can cheaply break up their power into multiple un-linkable miner actors that will be able to mine on multiple chains without being caught mining at the same height at the same time.

## Implementation Notes

- When implementing parent selection, for now, always choose the single heaviest parent set for building on top of. More advanced strategies can be implemented later.
- When selecting messages from the mempool to include in your block, be aware that other miners may also generate blocks during this round, and if you want to maximize your fee earnings it may be best to select some messages at random (second in a duplicate earns no fees).

## Open Questions

- Whether to look back and how far for seed for leader election.
- When selecting between two forks of equal weight, one strategy might be to select the 'Tipset' with the lowest number of null tickets for a given block height and weight.