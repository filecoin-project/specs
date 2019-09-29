---
title: Storage Power Consensus
entries:
- expected_consensus
- storage_power_actor
---

{{<label storage_power_consensus>}}
The Storage Power Consensus subsystem is the main interface which enables Filecoin nodes to agree on the state of the system. SPC accounts for individual storage miners' effective power over consensus in given chains in its _Power Table_. It also runs _Expected Consensus_ (the underlying consensus algorithm in use by Filecoin), enabling storage miners to run leader election and generate new blocks updating the state of the Filecoin system.

Succinctly, the SPC subsystem offers the following services:
- Access to the _Power Table_ for every subchain, accounting for individual storage miner power and total power on-chain.
- Access to {{<sref expected_consensus>}} for individual storage miners, enabling:
    - Access to verifiable randomness {{<sref tickets>}}as needed in the rest of the protocol. 
    - Running  {{<sref leader_election>}} to produce new blocks.
    - Running {{<sref chain_selection>}} across subchains using EC's weighting function. 
    - Identification of {{<sref finality "the most recently finalized tipset">}}, for use by all protocol participants.

Much of the Storage Power Consensus' subsystem functionality is detailed in the code below but we touch upon some of its behaviors in more detail.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}

## Losing Tickets and repeated leader election attempts

In the case that everybody draws a losing ticket in a given round of EC (i.e. no miner is eligible to produce a block), the storage power consensus subsystem will allow a given miner to run leader election again by "scratching" (attempting to generate a new `ElectionProof` from) the next ticket in the chain. That is, miners will now use the ticket sampled `K-1` rounds back to generate a new `ElectionProof`. They can then compare that proof with their current power fraction. This is repeated until a miner scratches a winning ticket and can publish a block (see [Block Generation](#block-generation)).

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

## A note on miners' 'power fraction'

The portion of blocks a given miner generates (and so the block rewards they earn) is proportional to their `Power Fraction` over time.

This means miners should not be able to mine using power they have not yet proven. Conversly, it is acceptable for miners to mine with a slight delay between their proving storage and that proven storage being reflected in leader election. This is reflected in the height at which the Power Table's `GetTotalStorage` and `PowerLookup` methods are called.

To illustrate this, an example:

Miner M1 has a provingPeriod of 30. M1 submits a PoST at height 39. Their next `provingPeriodEnd` will be 69, but M1 can submit a new PoST at any height X, for X in (39, 69]. Let's assume X is 67.

At height Y in (39, 67], M1 will attempt to generate an `ElectionProof` using the storage market actor from height 39 for their own power (and an actor from Y for total network power); at height 68, M1 will use the storage market actor from height 67 for their own power, and the storage market actor from height 68 for total power and so on.