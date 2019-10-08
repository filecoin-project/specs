---
title: "Expected Consensus"
---

{{<label expected_consensus>}}
## Algorithm

Expected Consensus (EC) is a probabilistic Byzantine fault-tolerant consensus protocol. At a high level, it operates by running a leader election every round in which, on expectation, one participant may be eligible to submit a block. EC guarantees that this winner will be anonymous until they reveal themselves by submitting a proof of their election (we call this proof an `Election Proof`). All valid blocks submitted in a given round form a `Tipset`. Every block in a Tipset adds weight to its chain. The 'best' chain is the one with the highest weight, which is to say that the fork choice rule is to choose the heaviest known chain. For more details on how to select the heaviest chain, see {{<sref chain_selection>}}.

At a very high level, with every new block generated, a miner will craft a new ticket from the prior one in the chain. While on expectation at least one block will be generated at every round, in cases where no one finds a block in a given round, a miner may increment a given nonce as part of the input with which they attempt to run leader election in order to ensure liveness in the protocol. These nonces help mark block height. Every block in a given Tipset will contain election proofs with the same nonce (i.e. they are mined at the same height).

The {{<sref storage_power_consensus>}} subsystem uses access to EC to use the following facilities:
- Access to verifiable randomness for the protocol, derived from {{<sref tickets>}}.
- Running and verifying {{<sref leader_election "leader election">}} for block generation.
- Access to a weighting function enabling {{<sref chain_selection>}} by the chain manager.
- Access to the most recently {{<sref finality "finalized tipset">}} available to all protocol participants.

{{<label tickets>}}
## Tickets

One may think of leader election in EC as a verifiable lottery, in which participants win in proportion to the power they have within the network.

A ticket is drawn from the past at the beginning of each new round to perform leader election. EC also generates a new ticket in every round for future use. Tickets are chained independently of the main blockchain. A ticket only depends on the ticket before it, and not any other data in the block.
On expectation, in Filecoin, every block header contains one ticket, though it could contain more if that block was generated over multiple rounds.

Tickets are used across the protocol as sources of randomness:
- The {{<sref sector_sealer>}} uses tickets to bind sector commitments to a given subchain.
- The {{<sref post_generator>}} likewise uses tickets to prove sectors remain committed as of a given block.
- EC uses them to run leader election and generates new ones for use by the protocol, as detailed below.

You can find the Ticket data structure {{<sref data_structures "here">}}.

### Comparing Tickets in a Tipset

Whenever comparing tickets is evoked in Filecoin, for instance when discussing selecting the "min ticket" in a Tipset, the comparison is that of the little endian representation of the ticket's VFOutput bytes.

### The Ticket chain

While each Filecoin block header contains a ticket field, it is useful to provide nodes with a ticket chain abstraction.

Namely, tickets are used throughout the Filecoin system as sources of on-chain randomness. For instance,
- They are drawn by Storage Miners as SealSeeds to commit new sectors
- They are drawn by Storage Miners as PoStChallenges to generate PoSts
- They are drawn by the Storage Power subsystem as randomness in leader election to determine their eligibility to mine a block
- They are drawn by the Storage Power subsystem in order to generate new tickets for future use.

Each of these ticket uses may require drawing tickets at different chain heights, according to the security requirements of the particular protocol making use of tickets. Due to the nature of Filecoin's Tipsets and the possibility of using losing tickets (that did not yield leaders in leader election) for randomness at a given height, tracking the canonical ticket of a subchain at a given height can be arduous to reason about in terms of blocks. To that end, it is helpful to create a ticket chain abstraction made up of only those tickets to be used for randomness at a given height.

This ticket chain will track one-to-one with a block at each height in a given subchain, but omits certain details including other blocks mined at that height.

Simply, it is composed inductively as follows. For a given chain:

- At height 0, take the genesis block, return its ticket
- At height n+1, take the heaviest tipset in our chain at height n.
    - select the block in that tipset with the smallest final ticket, return its ticket

## Tickets in EC

Within EC, a miner generates a new ticket in their block for every ticket they use (or "scratch") running leader election, thereby ensuring the ticket chain is always as long as the block chain.

Tickets are used to achieve the following:
- Ensure leader secrecy -- meaning a block producer will not be known until they release their block to the network.
- Prove leader election -- meaning a block producer can be verified by any participant in the network.


In practice, EC defines two different fields within a block:

- A `Ticket` field — this stores the new ticket generated during this block generation attempt. It is from this ticket that miners will sample randomness to run leader election in `K` rounds.
- An `ElectionProof` — this stores a proof that a given miner scratched a winning lottery ticket using the appropriate ticket `K` rounds back along with a nonce showing how many rounds generating the EP took. It proves that the leader was elected in this round.

```
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

### Ticket generation

This section discusses how tickets are generated by EC for the `Ticket` field.

At round `N`, a new ticket is generated using tickets drawn from the Tipset at round `N-1`. Because a Tipset can contain multiple blocks (see [Chain Selection](#chain-selection) below), the smallest ticket in the Tipset must be drawn otherwise the block will be invalid.

```
   ┌──────────────────────┐                     
   │                      │                     
   │                      │                     
   │┌────┐                │                     
   ││ TA │              A │                     
   └┴────┴────────────────┘                     
                                                
   ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─                      
                          │                     
   │                                            
    ┌────┐                │       TA < TB < TC  
   ││ TB │              B                       
    ┴────┘─ ─ ─ ─ ─ ─ ─ ─ ┘                     
                                                
   ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─                      
                          │                     
   │                                            
    ┌────┐                │                     
   ││ TC │              C                       
    ┴────┘─ ─ ─ ─ ─ ─ ─ ─ ┘                     
```

In the above diagram, a miner will use block A's Ticket to generate a new ticket (or an election proof farther in the future) since it is the smallest in the Tipset.

The miner runs the prior ticket through a Verifiable Random Function (VRF) to get a new unique output.

The VRF's deterministic output adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance).

We use the ECVRF algorithm from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve
  - Note that the operation type in step 2.1 is necessary to prevent an adversary from guessing an election proof for a miner ahead of time.

They are crafted as follows:

Input: parentTickets at round N-1, miner's private key SK
Output: newTicket

```
0. Prepare new ticket
    i. newTicket <-- New()
1. Draw prior ticket
    i. priorTicket <-- chain.TicketAtEpoch(N-1)
2. Run it through VRF and get determinstic output
    i. # take the VRFResult of that ticket as input, specifying the personalization (see data-structures)
            input <-- VRFPersonalization.Ticket | priorTicket.Output
    ii. # run it through the VRF and store the output in the new ticket
            newTicket.vrfResult <-- self.VRFKeyPair.Generate(input)
4. Return the new ticket
```

### Ticket Validation

Each Ticket should be generated from the prior one in the ticket-chain.

Succinctly, the process of verifying a block's tickets is as follows.
```text
Input: received block, storage market actor S, miner's public key PK, a public VDF validation key vk
Output: 0, 1

0. Get the ticket
    i. ticket <-- block.ticket
1. Verify its VRF Proof
    i. # get the appropriate parent
        parent = chain.TicketAtEpoch(N-1)
    ii. # generate the VRFInput
        input <-- VRFPersonalization.Ticket | parent.Output
    iii. # verify the VRF
        VRFState <-- ticket.VrfResult.Verify(input, PK)
        if VRFState == false:
            return 0
2. Return results
    return 1
```

{{<label leader_election>}}
## Secret Leader Election

Expected Consensus is a consensus protocol that works by electing a miner from a weighted set in proportion to their power. In the case of Filecoin, participants and powers are drawn from the storage [power table](storage-market.md#the-power-table), where power is equivalent to storage provided through time.

Leader Election in Expected Consensus must be Secret, Fair and Verifiable. This is achieved through the use of randomness used to run the election. In the case of Filecoin's EC, the blockchain tracks an independent ticket chain. These tickets are used as randomness inputs for Leader Election. Every block generated references an `ElectionProof` derived from a past ticket. The ticket chain is extended by the miner who generates a new ticket with each attempted election.

### Running a leader election

Now, a miner must also check whether they are eligible to mine a block in this round. For how Election Proofs are validated, see [election validation](mining.md#election-validation).

To do so, the miner will use tickets from K rounds back as randomness to uniformly draw a value from 0 to 1. Comparing this value to their power, they determine whether they are eligible to mine. A user's `power` is defined as the ratio of the amount of storage they proved as of their last PoSt submission to the total storage in the network as of the current block.

We use the ECVRF algorithm (must yield a pseudorandom, deterministic output) from [Goldberg et al. Section 5](https://tools.ietf.org/html/draft-irtf-cfrg-vrf-04#page-10), with:
  - Sha256 for our hashing function
  - Secp256k1 for our curve

If the miner scratches a winning ticket in this round, it can use newEP, along with a newTicket to generate and publish a new block (see [Block Generation](#block-generation)). Otherwise, it waits to hear of another block generated in this round.

It is important to note that every block contains two artifacts: one, a ticket derived from last block's ticket to prove that they have waited the appropriate delay, and two, an election proof derived from the ticket `K` rounds back used to run leader election.

Succinctly, the process of crafting a new ElectionProof in round N is as follows. We use:

    The ECVRF algorithm (must yield a pseudorandom, deterministic output) from Goldberg et al. Section 5, with:
        Sha256 for our hashing function
        Secp256k1 for our curve

Note: We draw the miner power from the prior round. This means that if a miner wins a block on their ProvingPeriodEnd even if they have not yet resubmitted a PoSt, they retain their power (until the next round).

At round N:

```
Input: electionSeed from N-K, miner's public key PK, miner's secret key SK, the Storage Power actor S, a Nonce nonce starting with value 0
Output: 1 or 0

0. Prepare new election proof
    i. newEP <-- New()
1. Draw electionSeed from K blocks back
    i. electionSeed <-- chain.TicketAtEpoch(N-K)
2. Run it through VRF and get determinstic output
    i. # take the VRFOutput of that ticket as input, specified for the appropriate operation type
        input <-- VRFPersonalization.ElectionProof | electionSeed.Output | nonce
    ii. # run it through the VRF and store the VRFProof in the new ticket
        newEP.VrfResult <-- self.VRFKeyPair.Generate(input)
        newEP.Nonce <-- nonce
3. Determine the miner's power fraction
    i. # Determine total storage this round
        S <-- storageMarket(N)
        p_n <-- S.GetTotalPower()
    ii. # Determine own power at prior tipSet
        p_m <-- LookupMinerPower(self.addr)
    iii. # Get power fraction
		p_f <-- p_m/p_n
4. Determine if the miner drew a winning lottery ticket
    # Conceptually we are mapping the pseudorandom, deterministic VRFOutput onto [0,1] by dividing by 2^HashLen (64 Bytes using Sha256) and comparing that to the miner's power (portion of network storage).
    # In practice, we actually multiply the power fraction by 2^HashLen for comparison with the ticket value in order to avoid dealing with floats.
    i.  # Map the miner's power onto [0, 2^HashLen] (64 Bytes using sha256)
        normalized_power <-- p_f * 2^HashLen
    ii. # Compare the miner's power fraction to the value of the ticket drawn: more power means win more often.
        # winning ticket
        if readLittleEndian(VRFOutput) <= normalized_power
            return newEP
        # otherwise parentTicket is not a winning lottery ticket
        else
            Return 0
Return 1
```

If successful, the miner can craft a block, passing it to the block producer. If unsuccessful, it will wait to hear of another block mined this round to try again. In the case no other block was found in this round the miner can increment nonce and try again.
While a miner could try to run through multiple nonces in parallel in order to quickly generate a block, this effort will be futile as the honest majority of miners will reject blocks crafted with ElectionProofs whose nonces prove too high (see below).

If successful, the miner can craft a block, passing it to the block producer. If unsuccessful, it will wait to hear of another block mined this round to try again. In the case no other block was found in this round the miner can increment nonce and try again.
While a miner could try to run through multiple nonces in parallel in order to quickly generate a block, this effort will be futile as the honest majority of miners will reject blocks crafted with ElectionProofs whose nonces prove too high (see below).

### Election Validation

In order to determine that the mined block was generated by an eligible miner, one must check its `ElectionProof`.

Succinctly, the process of verifying a block's election proof at round N, is as follows.

```text
Input: received block, storage power actor S, miner's public key PK, a public parameter K
Output: 0, 1

0. Get the election proof
        i. electionProof <-- block.electionProof
1. Verify it was generated in the appropriate round, specifically, for a block generated at round N, if the last block was generated at round N-L, the EP's nonce value should be L. Any EP with a nonce L' > L should be rejected.
        i. # get nonce
            Nonce <-- electionProof.Nonce
            cur_round <-- chain.Current_Round()
            if block.Parents.Height() + Nonce > cur_round:
                return 0
2. Get the total power, miner power
        i. # get total market power
            S <-- storageMarket(N)
            p_n <-- S.GetTotalPower()
        ii. # get miner power
            p_m <-- LookupMinerPower(miner.addr)
        iii. # Get power fraction
              p_f <-- p_m/p_n
3. Get the appropriate ticket from the ticket chain
        i. # Get the tipset K rounds back
            ticket <-- chain.TicketAtEpoch(N-K)
4. Verify Election Proof validity
        i. # generate the VRFInput from the scratched ticket
            input <-- VRFPersonalization.ElectionProof | scratchedTicket.VDFOutput
        ii. # Check that the election proof was correctly generated by the miner
            # using the appropriate ticket
            VRFState <-- electionProof.Verify(input, miner.PK)
            if VRFState == false:
                return 0
        iii. Ensure that the scratched ticket is a winner
            # internally -> map p_f onto [0, 2^HashLen] and compare the miner's scratchValue to the miner's normalized power fraction
            if electionProof.IsWinning(p_f) == false:
                return 0
5. Everything checks out, it's a valid election proof
        return 1
```

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

`w[r+1] = w[r] + (wForkFactor[r+1](wPowerFactor[r+1] + wBlocksFactor[r+1]))`

For a given tipset `ts` in round `r+1`, we define:

- wPowerFactor[r+1]  = log2b(totalPowerAtTipset(ts))
- wBlocksFactor[r+1] = v * |blocksInTipset(ts)| = v * k
  - with v = vt * log2b(totalPowerAtTipset(ts) * vs)/vs


Take X -> Bin(eNumberOfBlocksPerRound * numberOfMinersInPowerTable, 1/numberOfMinersInPowerTable), for simplicity, we take:
X -> Bin(eNumberOfBlocksPerRound * 10,000, 1/10,000). We have:
- wForkFactor[r+1]   = 1 if |blocksInTipset(ts)| > E[X] - 2*stdDev[X], CDF(X, |blocksInTipset(ts)) otherwise.

with:
    - log2b(X) = floor(log2(x)), the binary length of X.
    - stdDev[X] = sqrt(eNumberOfBlocksPerRound * (1 - 1/10,000))
    - CDF(X, |blocksInTipset(ts)) = Sum_i=0^k ((10,000* choose i) (1/10,000)^i (9,999/10,000)n- i


```sh
Note that if your implementation does not allow for rounding to the fourth decimal, miners should apply the [tie-breaker below](#selecting-between-tipsets-with-equal-weight). Weight changes will be on the order of single digit numbers on expectation, so this should not have an outsized impact on chain consensus across implementations.
```

 The exact value for these parameters remain to be determined, but for testing purposes, you may use:
 - `vt = 350`
 - `vs = E-5`

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
EC enforces a version of soft finality whereby all miners at round N will reject all blocks coming in prior to round N-F. For illustrative purposes, we can take F to be 500. While strictly speaking EC is a probabilistically final protocol, choosing such an F simplifies miner implementations and enforces a socially-bound finality at no cost to liveness in the chain.

## Slashing in EC

Due to the existence of potential forks, a miner can try to unduly influence protocol fairness. This means they may choose to disregard the protocol in order to gain an advantage over the power they should normally get from their storage on the network. A miner should be slashed if they are provably deviating from the honest protocol.

This is detectable when a miner submits two blocks that satisfy either of the following "slashing conditions":

(1) one block contains the same electionProof nonce as one of the other blocks'.
(2) one block's parent is a Tipset that could have validly included the other block according to Tipset validity rules, however the parent of the first block does not include the other block.

  - While it cannot be proven that a miner omits known blocks from a Tipset in general (i.e. network latency could simply mean the miner did not receive a particular block) in this case it can be proven because a miner must be aware of a block they mined in a previous round.

Any node that detects this occurring should take both block headers, and call [`storagemarket.SlashConsensusFault`](actors.md#slashconsensusfault). The network will then take all of that node's collateral, give a portion of it to
the reporter, and keep the rest.
