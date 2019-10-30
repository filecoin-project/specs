---
title: "Expected Consensus"
---

{{<label expected_consensus>}}
## Algorithm

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high level, it operates by running a leader election every round in which, on expectation, one participant may be eligible to submit a block. EC guarantees that this winner will be anonymous until they reveal themselves by submitting a proof of their election (we call this proof an `Election Proof`). All valid blocks submitted in a given round form a `Tipset`. Every block in a Tipset adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see {{<sref chain_selection>}}.

At a very high level, with every new block generated, a miner will craft a new ticket from the prior one in the chain appended with the current epoch number (i.e. parentTipset.epoch + 1 to start). While on expectation at least one block will be generated at every round, in cases where no one finds a block in a given round, a miner will increment the round number and attempt a new leader election (using the new input) thereby ensuring liveness in the protocol.

The {{<sref storage_power_consensus>}} subsystem uses access to EC to use the following facilities:
- Access to verifiable randomness for the protocol, derived from {{<sref tickets>}}.
- Running and verifying {{<sref leader_election "leader election">}} for block generation.
- Access to a weighting function enabling {{<sref chain_selection>}} by the chain manager.
- Access to the most recently {{<sref finality "finalized tipset">}} available to all protocol participants.

{{<label tickets>}}
## Tickets

For leader election in EC, participants win in proportion to the power they have within the network.

A ticket is drawn from the past at the beginning of each new round to perform leader election. EC also generates a new ticket in every round for future use. Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any other data in the block.
On expectation, in Filecoin, every block header contains one ticket, though it could contain more if that block was generated over multiple rounds.

Tickets are used across the protocol as sources of randomness:
- The {{<sref sector_sealer>}} uses tickets to bind sector commitments to a given subchain.
- The {{<sref post_generator>}} likewise uses tickets to prove sectors remain committed as of a given block.
- EC uses them to run leader election and generates new ones for use by the protocol, as detailed below.

You can find the Ticket data structure {{<sref data_structures "here">}}.

### Comparing Tickets in a Tipset

Whenever comparing tickets is evoked in Filecoin, for instance when discussing selecting the "min ticket" in a Tipset, the comparison is that of the little endian representation of the ticket's VFOutput bytes.

## Tickets in EC

Within EC, a miner generates a new ticket in their block for every ticket they use running leader election, thereby ensuring the ticket chain is always as long as the block chain.

Tickets are used to achieve the following:
- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.


In practice, EC defines two different fields within a block:

- A `Ticket` field — this stores the new ticket generated during this block generation attempt. It is from this ticket that miners will sample randomness to run leader election in `K` rounds.
- An `ElectionProof` — this stores a proof that a given miner has won a leader election using the appropriate ticket `K` rounds back appended with the current epoch number. It proves that the leader was validly elected in this epoch.

```
But why the randomness lookback?

The randomness lookback helps turn independent ticket generation from a block one round back
into a global ticket generation game instead. Rather than having a distinct chance of winning or losing
for each potential fork in a given round, a miner will either win on all or lose on all
forks descended from the block in which the ticket is sampled.

This is useful as it reduces opportunities for grinding, across forks or sybil identities.

However this introduces a tradeoff:
- The randomness lookback means that a miner can know K rounds in advance that they will win,
decreasing the cost of running a targeted attack (given they have local predictability).
- It means electionProofs are stored separately from new tickets on a block, taking up
more space on-chain.

How is K selected?
- On the one end, there is no advantage to picking K larger than finality.
- On the other, making K smaller reduces adversarial power to grind.
```

### Ticket generation

This section discusses how tickets are generated by EC for the `Ticket` field.

At round `N`, a new ticket is generated using tickets drawn from the Tipset at round `N-1` (for more on how tickets are drawn see {{<sref ticket_chain>}}).

The miner runs the prior ticket through a Verifiable Random Function (VRF) to get a new unique output.

The VRF's deterministic output adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance).

We use the ECVRF algorithm from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve
  - Note that the operation type in step 2.1 is necessary to prevent an adversary from guessing an election proof for a miner ahead of time.

### Ticket Validation

Each Ticket should be generated from the prior one in the ticket-chain.

{{<label leader_election>}}
## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the storage [power table](storage-market.md#the-power-table), where power is equivalent to storage provided through time.

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain tracks an independent ticket chain. These tickets are used as randomness inputs for Leader Election. Every block generated references an `ElectionProof` derived from a past ticket. The ticket chain is extended by the miner who generates a new block for each successful leader election.

### Running a leader election

Now, a miner must also check whether they are eligible to mine a block in this round.

To do so, the miner will use tickets from K rounds back as randomness to uniformly draw a value from 0 to 1. Comparing this value to their power, they determine whether they are eligible to mine. A user's `power` is defined as the ratio of the amount of storage they proved as of their last PoSt submission to the total storage in the network as of the current block.

We use the ECVRF algorithm (must yield a pseudorandom, deterministic output) from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve

If the miner wins the election in this round, it can use newEP, along with a newTicket to generate and publish a new block. Otherwise, it waits to hear of another block generated in this round.

It is important to note that every block contains two artifacts: one, a ticket derived from last block's ticket to extend the ticket-chain, and two, an election proof derived from the ticket `K` rounds back used to run leader election.

Succinctly, the process of crafting a new ElectionProof in round N is as follows. We use:

    The ECVRF algorithm (must yield a pseudorandom, deterministic output) from Goldberg et al. Section 5, with:
        Sha256 for our hashing function
        Secp256k1 for our curve

Note: We draw the miner power from the prior round. This means that a miner can win a block in an epoch in which they were supposed to have proven storage again but did not. Put another way, `ElectionProof` validity is checked prior to executing transactions (including that which cuts miner power) at every epoch.

If successful, the miner can craft a block, passing it to the block producer. If unsuccessful, it will wait to hear of another block mined this round to try again. In the case no other block was found in this round the miner can increment the epoch number and try leader election again using the same past ticket and new epoch number.
While a miner could try to run through multiple epochs in parallel in order to quickly generate a block, this effort will be futile as the rational majority of miners will reject blocks crafted with ElectionProofs whose epochs prove too high (i.e. in the future -- see below).

### Election Validation

In order to determine that the mined block was generated by an eligible miner, one must check its `ElectionProof`'s validity and that its input was generated using the current epoch value.

{{<label chain_selection>}}
## Chain Selection

Just as there can be 0 miners win in a round, multiple miners can be elected in a given round. This in turn means multiple blocks can be created in a round. In order to avoid wasting valid work done by miners, EC makes use of all valid blocks generated in a round.

### Chain Weighting

It is possible for forks to emerge naturally in Expected Consensus. EC relies on weighted chains in order to quickly converge on 'one true chain', with every block adding to the chain's weight. This means the heaviest chain should reflect the most amount of work performed, or in Filecoin's case, the most storage provided.

In short, the weight at each block is equal to its `ParentWeight` plus that block's delta weight. Details of Filecoin's chain weighting function [are included here](https://observablehq.com/d/3812cd65c054082d).

Delta weight is a term composed of a few elements:
- wForkFactor: which seeks to cut the weight derived from rounds in which produced Tipsets do not correspond to what an honest chain is likely to have yielded (pointing to selfish mining or other non-collaborative miner behavior).
- wPowerFactor: which adds weight to the chain proportional to the total power backing the chain, i.e. accounted for in the chain's power table.
- wBlocksFactor: which adds weight to the chain proportional to the number of blocks mined in a given round. Like wForkFactor, it rewards miner cooperation (which will yield more blocks per round on expectation).

The weight should be calculated using big integer arithmetic with order of operations defined above. We use brackets instead of parentheses below for legibility. We have:

`w[r+1] = w[r] + (wPowerFactor[r+1] + wBlocksFactor[r+1]) * 2^8`

For a given tipset `ts` in round `r+1`, we define:

- `wPowerFactor[r+1]  = wFunction(totalPowerAtTipset(ts))`
- wBlocksFactor[r+1] =  `wPowerFactor[r+1] * wRatio * b / e`
  - with `b = |blocksInTipset(ts)|`
  - `e = expected number of blocks per round in the protocol`
  - and `wRatio in ]0, 1[`
Thus, for stability of weight across implementations, we take:
- wBlocksFactor[r+1] =  `(wPowerFactor[r+1] * b * wRatio_num) / (e * wRatio_den)`

We get:
- `w[r+1] = w[r] + wFunction(totalPowerAtTipset(ts)) * 2^8 + (wFunction(totalPowerAtTipset(ts)) * len(ts.blocks) * wRatio_num * 2^8) / (e * wRatio_den)`
 Using the 2^8 here to prevent precision loss ahead of the division in the wBlocksFactor.

 The exact value for these parameters remain to be determined, but for testing purposes, you may use:
 - `e = 5`
 - `wRatio = .5, or wRatio_num = 1, wRatio_den = 2`
- `wFunction = log2b` with
  - `log2b(X) = floor(log2(x)) = (binary length of X) - 1` and `log2b(0) = 0`. Note that that special case should never be used (given it would mean an empty power table.

```sh
Note that if your implementation does not allow for rounding to the fourth decimal, miners should apply the [tie-breaker below](#selecting-between-tipsets-with-equal-weight). Weight changes will be on the order of single digit numbers on expectation, so this should not have an outsized impact on chain consensus across implementations.
```

`ParentWeight` is the aggregate chain weight of a given block's parent set. It is calculated as
the `ParentWeight` of any of its parent blocks (all blocks in a given Tipset should have
the same `ParentWeight` value) plus the delta weight of each parent. To make the
computation a bit easier, a block's `ParentWeight` is stored in the block itself (otherwise
potentially long chain scans would be required to compute a given block's weight).

### Selecting between Tipsets with equal weight

When selecting between Tipsets of equal weight, a miner chooses the one with the smallest final ticket.

In the case where two Tipsets of equal weight have the same min ticket, the miner will compare the next smallest ticket (and select the Tipset with the next smaller ticket). This continues until one Tipset is selected.

The above case may happen in situations under certain block propagation conditions. Assume three blocks B, C, and D have been mined (by miners 1, 2, and 3 respectively) off of block A, with minTicket(B) < minTicket(C) < minTicket (D).

Miner 1 outputs their block B and shuts down. Miners 2 and 3 both receive B but not each others' blocks. We have miner 2 mining a Tipset made of B and C and miner 3 mining a Tipset made of B and D. If both succesfully mine blocks now, other miners in the network will receive new blocks built off of Tipsets with equal weight and the same smallest ticket (that of block B). They should select the block mined atop [B, C] since minTicket(C) < minTicket(D).

The probability that two Tipsets with different blocks would have all the same tickets can be considered negligible: this would amount to finding a collision between two 256-bit (or more) collision-resistant hashes.

{{<label finality>}}
## Finality in EC
EC enforces a version of soft finality whereby all miners at round N will reject all blocks that fork off prior to round N-F. For illustrative purposes, we can take F to be 500. While strictly speaking EC is a probabilistically final protocol, choosing such an F simplifies miner implementations and enforces a macroeconomically-enforced finality at no cost to liveness in the chain.

{{<label consensus_faults>}}
## Consensus Faults

Due to the existence of potential forks in EC, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This is detectable when a given miner submits two blocks that satisfy any of the following "consensus faults":

- (1) `double-fork mining fault`: two blocks mined off of the same Tipset at different epochs (i.e. with different `ElectionProof`s generated from the same input ticket).
- (2) `parent grinding fault`: one block's parent is a Tipset that provably should have included a given block but does not. While it cannot be proven that a missing block was willfully omitted in general (i.e. network latency could simply mean the miner did not receive a particular block), it can when a miner has successfully mined a block two epochs in a row and omitted one. That is, this condition should be evoked when a miner omits their own prior block. When a miner's block at epoch e + 1 references a Tipset that does not include the block they mined at e both blocks can be submitted to prove this fault.

Any node that detects either of the above events should submit both block headers to the `StoragePowerActor`'s `ReportConsensusFault` method. The "slasher" will receive a portion (TODO: define how much) of the offending miner's {{<sref pledge_collateral>}} as a reward for notifying the network of the fault.
(TODO: FIP of submitting commitments to block headers to prevent miners censoring slashers in order to gain rewards).

It is important to note that there exists a third type of consensus fault directly reported by the `CronActor` on `StorageDeal` failures via the `ReportUncommittedPowerFault` method:
- (3) `uncommitted power fault` which occurs when a miner fails to submit their `PostProof` and is thus participating in leader election with undue power (see {{<sref storage_faults>}}).
