---
title: Storage Power Consensus
entries:
- expected_consensus
- storage_power_actor
---

The Storage Power Consensus subsystem tracks individual storage miners' storage-based power through given chains in an associated _Power Table_. It exposes methods to the rest of the blockchain system enabling it to update miner power based on their behavior (such as proving or committing sectors or misbehaving), as well as enabling it to run successive trials of leader election through Expected Consensus.

Much of the Storage Power Consensus' subsystem functionality is detailed in the code below but we touch upon some of its behaviors in more detail.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}

### Tipsets

All valid blocks generated in a round form a `Tipset` that participants will attempt to mine off of in the subsequent round (see above). Tipsets are valid so long as:

- All blocks in a Tipset have the same parent Tipset
- All blocks in a Tipset have the same number of tickets in their `Tickets` array

These conditions imply that all blocks in a Tipset were mined at the same height. This rule is key to helping ensure that EC converges over time. While multiple new blocks can be mined in a round, subsequent blocks all mine off of a Tipset bringing these blocks together. The second rule means blocks in a Tipset are mined in a same round.

The blocks in a tipset have no defined order in representation. During state computation, blocks in a tipset are processed in order of block ticket, breaking ties with the block CID bytes.

Due to network propagation delay, it is possible for a miner in round N+1 to omit valid blocks mined at round N from their Tipset. This does not make the newly generated block invalid, it does however reduce its weight and chances of being part of the canonical chain in the protocol as defined by EC's {{<sref chain_weighting>}} function.

## Losing Tickets and repeated leader election attempts

In the case that everybody draws a losing ticket in a given round of EC (i.e. no miner is eligible to produce a block), the storage power consensus subsystem should prompt a given miner node to run leader election again by "scratching" (attempting to generate a new `ElectionProof` from) the next ticket in the chain. That is, miners will now use the ticket sampled `K-1` rounds back to generate a new `ElectionProof`. They can then compare that proof with their current power fraction. This is repeated until a miner scratches a winning ticket and can publish a block (see [Block Generation](#block-generation)).

In addition to each attempted `ElectionProof` generation, the miner will need to extend the ticket chain by generating another new ticket. They use the ticket they generated in the prior round, rather than the prior block's (as is normally used). This proves appropriate delay (given that finding a winning Ticket has taken multiple rounds).

Thus, each time it is discovered that nobody has won in a given round, every miner should append a new ticket to their would-be block's `Ticket` array. This continues until some miner finds a winning ticket (see below), ensuring that the ticket chain remains at least as long as the block chain.

The length of repeated losing tickets in the ticket chain (equivalent to the length of generated tickets referenced by a single block, or the length of the `Tickets` array) decreases exponentially in the number of repeated losing tickets. In the unlikely case the number of losing tickets drawn by miners grows larger than the randomness lookback `K` (i.e. if a miner runs out of existing tickets on the ticket chain for use as randomness), a miner should proceed as usual using new tickets generated in this epoch for randomness. This has no impact on the protocol safety/validity.

New blocks (with multiple tickets) will have a few key properties:

- All tickets in the `Tickets` array are signed by the same miner -- to avoid grinding through out-of-band collusion between miners exchanging tickets.
- The `ElectionProof` was correctly generated from the ticket `K-|Tickets|` (with `|Tickets|` the length of the `Tickets` array) rounds back.

This means that valid `ElectionProof`s can be generated from tickets in the middle of the `Tickets` array. In cases where there are multiple tickets to choose from (i.e. a `Tipset` made up of multiple blocks mined atop losing tickets), a miner must use tickets from the `Tickets` array in the block whose final ticket is the min-ticket (though that may not be the min ticket in that round).


                  ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
                                         │
                  │
                   ┌────┬────┬────┐      │
                  ││T1-A│T2-A│T3-A│    A
                   ┴────┴────┴────┘─ ─ ─ ┘   T1-A < T1-B
                                             T2-A < T2-B
                  ┌──────────────────────┐   T3-A > T3-B
                  │                      │
                  │                      │
                  │┌────┬────┬────┐      │
                  ││T1-B│T2-B│T3-B│    B │
                  └┴────┴────┴────┴──────┘
In the above case, because T3-B < T3-A, then a miner will use T*-B to generate election proofs, even if T1-A < T1-B.

## The Ticket chain

While each Filecoin block header contains a ticket array, it is useful to provide nodes with a ticket chain abstraction.

Namely, tickets are used throughout the Filecoin system as sources of on-chain randomness. For instance,
- They are drawn by Storage Miners as SealSeeds to commit new sectors
- They are drawn by Storage Miners as PoStChallenges to generate PoSts
- They are drawn by the Storage Power subsystem as randomness in leader election to determine their eligibility to mine a block
- They are drawn by the Storage Power subsystem in order to generate new tickets for future use.

Each of these ticket uses may require drawing tickets at different chain heights, according to the security requirements of the particular protocol making use of tickets. Due to the nature of Filecoin's Tipsets and the possibility of using losing tickets (that did not yield leaders in leader election) for randomness at a given height, tracking the canonical ticket of a subchain at a given height can be arduous to reason about in terms of blocks. To that end, it is helpful to create a ticket chain abstraction made up of only those tickets to be used for randomness at a given height.

This ticket chain will track one-to-one with a block at each height in a given subchain, but omits certain details including other blocks mined at that height.

Simply, it is composed inductively as follows:

- At height 0, take the genesis block, return its ticket
- At height n+1, take the block used at height n.
  - If the ticket it returned in the previous step is not the final ticket in its ticket array, return the next one.
  - If the ticket it returned in the previous step is the final ticket in its ticket array,
    - use EC to select the next heaviest tipset in the subchain.
    - select the block in that tipset with the smallest final ticket, return its first ticket

We illustrate this below:
```
                     ┌────────────────────────────────────────────────────────────────────────────────────┐
                      │                                                                                    │
                      ▼                                  ◀─────────────────────────────────────────────────┼────────────────────────────────────────────┐
          ┌─┐        ┌─┐        ┌─┐        ┌─┐          ┌─┐                                       ┌─┐      │    ┌─┐         ┌─┐         ┌─┐           ┌─┤
   ◀──────│T1◀───────│T2◀───────│T3◀───────│T4◀─────────│T5◀───────────────    ...    ◀───────────T1+k─────┼───T2+k◀────────T3+k────────T4+k──────────T5│k
          └─┘        └─┘        └─┘        └─┘          └─┘                                       └─┘      │    └─┘         └─┘         └─┘           └─┤
                                 ▲                       │                                                 │                 │           │             ││
           │          │          │          │                                                      │       │     │                                      │
                                 ├───────────────────────┼─────────────────────────────────────────────────┼─────────────────┼────┐      │             ││
           │          │                     │                                                      │       │     │                │                     │
                       ─ ─ ─ ─ ─ ┴ ─   ─ ─ ─             ┘                                            ─ ─ ─│─ ─ ─          ─ ┘    │      └ ─ ─ ─   ─ ─ ┘│
        ┌──┼──────┬─┬┐         ┌──┼─┼─┼──┬─┬┐        ┌──┼──────┬─┬┐                             ┌──┼─┼────┬─┬┐         ┌──┼──────┬─┬┐        ┌──┼─┼────┬─┬┐
        │         │E1│         │         │E2│        │         │E3│                             │         E1+l         │         E2+l        │         E3+l
        │  │      └─┘│         │  │ │ │  └─┘│        │  │      └─┘│                             │  │ │    └─┘│         │  │      └─┘│        │  │ │    └─┘│
        │            │         │            │        │            │                             │            │         │            │        │            │
 ◀──────│  ▼         │◀────────│  ▼ ▼ ▼     │◀───────│  ▼         │◀──────     ...     ◀────────│  ▼ ▼       │◀────────│  ▼         │◀───────│  ▼ ▼       │
        │ ┌─┐        │         │ ┌─┬─┬─┐    │        │ ┌─┐        │                             │ ┌─┬─┐      │         │ ┌─┐        │        │ ┌─┬─┐      │
        │ │T1     B1 │         │ │T2 │ │ B2 │        │ │T5     B3 │                             │ │ │ │  B1+l│         │ │ │    B2+l│        │ │ │ │  B3+l│
        └─┴─┴────────┘         └─┴─┴─┴─┴────┘        └─┴─┴────────┘                             └─┴─┴─┴──────┘         └─┴─┴────────┘        └─┴─┴─┴──────┘
```



The above represents an instance of a block-chain, and its associated ticket-chain. The details of how it works should be clear by the end of this section, but quickly, some observations:

- Block 1 contains a single ticket T1
- Block 2 contains 3 tickets T2, T3, T4 meaning it was likely generated after 2 failed leader election attempts in the network.
- Block B1+l has two tickets and an Election Proof E1+l that was generated using T2. This means B1+l's miner tried to generate an election proof using T1 and failed, succeeding on their second attempt with T2.

## A note on miners' 'power fraction'

The portion of blocks a given miner generates (and so the block rewards they earn) is proportional to their `Power Fraction` over time.

This means miners should not be able to mine using power they have not yet proven. Conversly, it is acceptable for miners to mine with a slight delay between their proving storage and that proven storage being reflected in leader election. This is reflected in the height at which the Power Table's `GetTotalStorage` and `PowerLookup` methods are called.

To illustrate this, an example:

Miner M1 has a provingPeriod of 30. M1 submits a PoST at height 39. Their next `provingPeriodEnd` will be 69, but M1 can submit a new PoST at any height X, for X in (39, 69]. Let's assume X is 67.

At height Y in (39, 67], M1 will attempt to generate an `ElectionProof` using the storage market actor from height 39 for their own power (and an actor from Y for total network power); at height 68, M1 will use the storage market actor from height 67 for their own power, and the storage market actor from height 68 for total power and so on.