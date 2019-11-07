---
title: Storage Power Consensus
statusIcon: ✅
entries:
- storage_power_actor
---

{{<label storage_power_consensus>}}
The Storage Power Consensus subsystem is the main interface which enables Filecoin nodes to agree on the state of the system. SPC accounts for individual storage miners' effective power over consensus in given chains in its {{<sref power_table>}}. It also runs {{<sref expected_consensus>}} (the underlying consensus algorithm in use by Filecoin), enabling storage miners to run leader election and generate new blocks updating the state of the Filecoin system.

Succinctly, the SPC subsystem offers the following services:
- Access to the {{<sref power_table>}} for every subchain, accounting for individual storage miner power and total power on-chain.
- Access to {{<sref expected_consensus>}} for individual storage miners, enabling:
    - Access to verifiable randomness {{<sref tickets>}} as needed in the rest of the protocol.
    - Running  {{<sref leader_election>}} to produce new blocks.
    - Running {{<sref chain_selection>}} across subchains using EC's weighting function.
    - Identification of {{<sref finality "the most recently finalized tipset">}}, for use by all protocol participants.

Much of the Storage Power Consensus' subsystem functionality is detailed in the code below but we touch upon some of its behaviors in more detail.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}

# Distinguishing between storage miners and block miners

There are two ways to earn Filecoin tokens in the Filecoin network:
- By participating in the {{<sref storage_market>}} as a storage provider and being paid by clients for file storage deals.
- By mining new blocks on the network, helping modify system state and secure the Filecoin consensus mechanism.

We must distinguish between both types of "miners" (storage and block miners). {{<sref leader_election>}} in Filecoin is predicated on a miner's storage power. Thus, while all block miners will be storage miners, the reverse is not necessarily true.

However, given Filecoin's "useful Proof-of-Work" is achieved through file storage (PoRep and PoSt), there is little overhead cost for storage miners to participate in leader election. Such a {{<sref storage_miner_actor>}} need only register with the {{<sref storage_power_actor>}} in order to participate in Expected Consensus and mine blocks.

# Repeated leader election attempts

In the case that no miner is eligible to produce a block in a given round of EC, the storage power consensus subsystem will be called by the block producer to attempt another leader election by incrementing the nonce appended to the ticket drawn from the past in order to attempt to craft a new valid `ElectionProof` and trying again.

{{<label ticket_chain>}}
## The Ticket chain and randomness on-chain

While each Filecoin block header contains a ticket field (see {{<sref tickets>}}), it is useful to provide nodes with a ticket chain abstraction.

Namely, tickets are used throughout the Filecoin system as sources of on-chain randomness. For instance,
- The {{<sref sector_sealer>}} uses tickets as SealSeeds to bind sector commitments to a given subchain.
- The {{<sref post_generator>}} likewise uses tickets as PoStChallenges to prove sectors remain committed as of a given block.
- They are drawn by the Storage Power subsystem as randomness in {{<sref leader_election>}} to determine their eligibility to mine a block
- They are drawn by the Storage Power subsystem in order to generate new tickets for future use.

Each of these ticket uses may require drawing tickets at different chain heights, according to the security requirements of the particular protocol making use of tickets. Due to the nature of Filecoin's Tipsets and the possibility of using losing tickets (that did not yield leaders in leader election) for randomness at a given height, tracking the canonical ticket of a subchain at a given height can be arduous to reason about in terms of blocks. To that end, it is helpful to create a ticket chain abstraction made up of only those tickets to be used for randomness at a given height.

This ticket chain will track one-to-one with a block at each height in a given subchain, but omits certain details including other blocks mined at that height.

It is composed inductively as follows. For a given chain:

- At height 0, take the genesis block, return its ticket
- At height n+1, take the heaviest tipset in our chain at height n.
    - select the block in that tipset with the smallest final ticket, return its ticket

Because a Tipset can contain multiple blocks, the smallest ticket in the Tipset must be drawn otherwise the block will be invalid.

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
