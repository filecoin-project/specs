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

## Distinguishing between storage miners and block miners

There are two ways to earn Filecoin tokens in the Filecoin network:
- By participating in the {{<sref storage_market>}} as a storage provider and being paid by clients for file storage deals.
- By mining new blocks on the network, helping modify system state and secure the Filecoin consensus mechanism.

We must distinguish between both types of "miners" (storage and block miners). {{<sref leader_election>}} in Filecoin is predicated on a miner's storage power. Thus, while all block miners will be storage miners, the reverse is not necessarily true.

However, given Filecoin's "useful Proof-of-Work" is achieved through file storage (PoRep and PoSt), there is little overhead cost for storage miners to participate in leader election. Such a {{<sref storage_miner_actor>}} need only register with the {{<sref storage_power_actor>}} in order to participate in Expected Consensus and mine blocks.

{{<label storage_power>}}
## On Power

Per the above, we also clearly distinguish putting storage on-chain from gaining power in consensus (sometimes called "Storage Power") as follows:

**Consensus power in Filecoin is determined by in-deal storage**. For instance, if a miner had a 32GB sector, 20 of which were used as part of storage deals, only those 20 would contribute to said miner's *Storage Power*.

Thereafter, power can be either *active* or *inactive* as defined in {{<sref storage_mining_subsystem>}}.

{{<label tickets>}}
## Tickets

Tickets are used across the Filecoin protocol as sources of randomness:
- The {{<sref sector_sealer>}} uses tickets as SealSeeds to bind sector commitments to a given subchain.
- The {{<sref post_generator>}} likewise uses tickets as PoStChallenges to prove sectors remain committed as of a given block.
- They are drawn by the Storage Power subsystem as randomness in {{<sref leader_election>}} to determine their eligibility to mine a block
- They are drawn by the Storage Power subsystem in order to generate new tickets for future use.

Each of these ticket uses may require drawing tickets at different chain epochs, according to the security requirements of the particular protocol making use of tickets. Specifically, the ticket output (which is a SHA256 output) is used for randomness.

In Filecoin, every block header contains a single ticket.

You can find the Ticket data structure {{<sref data_structures "here">}}.

### Comparing Tickets in a Tipset

Whenever comparing tickets is evoked in Filecoin, for instance when discussing selecting the "min ticket" in a Tipset, the comparison is that of the little endian representation of the ticket's VFOutput bytes.

{{<label ticket_chain>}}
## The Ticket chain and drawing randomness

While each Filecoin block header contains a ticket field (see {{<sref tickets>}}), it is useful to think of a ticket chain abstraction.
Due to the nature of Filecoin's Tipsets and the possibility of using tickets from epochs that did not yield leaders to produce randomness at a given epoch, tracking the canonical ticket of a subchain at a given height can be arduous to reason about in terms of blocks. To that end, it is helpful to create a ticket chain abstraction made up of only those tickets to be used for randomness generation at a given height.

To read more about specifically how tickets are processed for randomness, see {{<sref randomness>}}.

To sample a ticket for a given epoch n:
```text
Set referenceTipsetOffset = 0
While true:
    Set referenceTipsetHeight = n - referenceTipsetOffset
    If blocks were mined at referenceTipsetHeight:
        ReferenceTipset = TipsetAtHeight(referenceTipsetHeight)
        Select the block in ReferenceTipset with the smallest final ticket, return its value (pastTicket).
    If no blocks were mined at referenceTipsetHeight:
        Increment referenceTipsetOffset
        (Repeat)
newRandomness = H(TicketDrawDST || index || Serialization(epoch || pastTicketOutput))
```

In english, this means two things:
- Choose the smallest ticket in the Tipset if it contains multiple blocks.
- When sampling a ticket from an epoch with no blocks, draw the min ticket from the prior epoch with blocks and concatenate it with
    - the wanted epoch number
    - hash this concatenation for a usable ticket value

See the `RandomnessAtEpoch` method below:
{{< readfile file="../struct/chain/chain.go" code="true" lang="go" >}}

The above means that ticket randomness is reseeded with every new block, but can indeed be derived by any miner for an arbitrary epoch number using a past epoch.

{{<label ticket_generation>}}
### Randomness Ticket generation

This section discusses how tickets are generated by EC for the `Ticket` field in every block header.

At round `N`, a new ticket is generated using tickets drawn from the Tipset at round `N-1` (as shown below).

The miner runs the prior ticket through a Verifiable Random Function (VRF) to get a new unique ticket which can later be derived for randomness (as shown above). The prior ticket is prepended with the ticket domain separation tag and concatenated with the miner actor address (to ensure miners using the same worker keys get different randomness).

To generate a ticket for a given epoch n:
```text
LastTicket = MinTicketValueAtEpoch(n-1)
newRandomness = VRF_miner(H(TicketProdDST || index || Serialization(pastTicket, minerActorAddress)))
```

The VRF's deterministic output adds entropy to the ticket chain, limiting a miner's ability to alter one block to influence a future ticket (given a miner does not know who will win a given round in advance).

We use the VRF from {{<sref vrf>}} for ticket generation in EC (see the `PrepareNewTicket` method below).

{{< readfile file="../../filecoin_mining/storage_mining/storage_mining_subsystem.go" code="true" lang="go" >}}


### Ticket Validation

Each Ticket should be generated from the prior one in the ticket-chain and verified accordingly as shown in `validateTicket` below.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}
{{< readfile file="storage_power_consensus_subsystem.go" code="true" lang="go" >}}

## Repeated Leader Election attempts

In the case that no miner is eligible to produce a block in a given round of EC, the storage power consensus subsystem will be called by the block producer to attempt another leader election by incrementing the nonce appended to the ticket drawn from the past in order to attempt to find a new winning `PartialTicket` and trying again. 
Note that a miner may attempt to grind through tickets by incrementing the nonce repeatedly until they find a winning ticket. However, any block so generated in the future will be rejected by other miners (with synchronized clocks) until that epoch's appropriate time.

{{<label min_miner_size>}}
## Minimum Miner Size

In order to secure Storage Power Consensus, the system defines a minimum miner size required to participate in consensus.

Specifically, miners must have either at least `MIN_MINER_SIZE_STOR` of active power (i.e. storage power currently used in storage deals) or `MIN_MINER_SIZE_PERC` of the network's active storage power to participate in leader election.

Miners smaller than this cannot mine blocks and earn block rewards in the network. Their power will not be counted as part of total network power. However, **it is important to note that such miners can still have their power faulted and be penalized accordingly**.

Accordingly, to bootstrap the network, the genesis block must include miners taking part in valid storage deals along with appropriate committed storage.

The `MIN_MINER_SIZE_PERC` condition will not be used in a network with more than `MIN_MINER_SIZE_STOR/MIN_MINER_SIZE_PERC` of active power. It is nonetheless defined to ensure liveness in small networks (e.g. close to genesis or after large power drops). Simply, a single miner can maintain network liveness for networks with less than `MIN_MINER_SIZE_STOR/MIN_MINER_SIZE_PERC` of active storage.

{{% notice placeholder %}}
The below values are currently placeholders.
{{% /notice %}}

We currently set:
- `MIN_MINER_SIZE_STOR = 1 << 40 Bytes` (100 TiB)
- `MIN_MINER_SIZE_PERC = .33`