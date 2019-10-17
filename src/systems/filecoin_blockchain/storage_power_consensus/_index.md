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

## Repeated leader election attempts

In the case that no miner is eligible to produce a block in a given round of EC, the storage power consensus subsystem will be called by the block producer to attempt another leader election by incrementing the nonce appended to the ticket drawn from the past in order to attempt to craft a new valid `ElectionProof` and trying again.

{{<label ticket_chain>}}
### The Ticket chain and randomness on-chain

While each Filecoin block header contains a ticket field (see {{<sref tickets>}}), it is useful to provide nodes with a ticket chain abstraction.

Namely, tickets are used throughout the Filecoin system as sources of on-chain randomness. For instance,
- The {{<sref sector_sealer>}} uses tickets as SealSeeds to bind sector commitments to a given subchain.
- The {{<sref post_generator>}} likewise uses tickets as PoStChallenges to prove sectors remain committed as of a given block.
- They are drawn by the Storage Power subsystem as randomness in {{<leader_election>}} to determine their eligibility to mine a block
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

## Drawing randomness for sector commitments

Tickets are used as input to the SEAL above in order to tie Proofs-of-Replication to a given chain, thereby preventing long-range attacks (from another miner in the future trying to reuse SEALs).

The ticket has to be drawn from a finalized block in order to prevent the miner from potential losing storage (in case of a chain reorg) even though their storage is intact.

Verification should ensure that the ticket was drawn no farther back than necessary by the miner. We note that tickets can uniquely be associated to a given round in the protocol (lest a hash collision be found), but that the round number is explicited by the miner in `commitSector`.

We present precisely how ticket selection and verification should work. In the below, we use the following notation:

- `F`-- Finality (number of rounds)
- `X`-- round in which SEALing starts
- `Z`-- round in which the SEAL appears (in a block)
- `Y`-- round announced in the SEAL `commitSector` (should be X, but a miner could use any Y <= X), denoted by the ticket selection
 - `T`-- estimated time for SEAL, dependent on sector size
 - `G = T + variance`-- necessary flexibility to account for network delay and SEAL-time variance.

We expect Filecoin will be able to produce estimates for sector commitment time based on sector sizes, e.g.:
`(estimate, variance) <--- SEALTime(sectors)`
G and T will be selected using these.

#### Picking a Ticket to Seal

When starting to prepare a SEAL in round X, the miner should draw a ticket from X-F with which to compute the SEAL.

#### Verifying a Seal's ticket

When verifying a SEAL in round Z, a verifier should ensure that the ticket used to generate the SEAL is found in the range of rounds [Z-T-F-G, Z-T-F+G].

#### In Detail

```
                               Prover
           ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
          │

          ▼
         X-F ◀───────F────────▶ X ◀──────────T─────────▶ Z
     -G   .  +G                 .                        .
  ───(┌───────┐)───────────────( )──────────────────────( )────────▶
      └───────┘                 '                        '        time
 [Z-T-F-G, Z-T-F+G]
          ▲

          └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─
                              Verifier
```

Note that the prover here is submitting a message on chain (i.e. the SEAL). Using an older ticket than necessary to generate the SEAL is something the miner may do to gain more confidence about finality (since we are in a probabilistically final system). However it has a cost in terms of securing the chain in the face of long-range attacks (specifically, by mixing in chain randomness here, we ensure that an attacker going back a month in time to try and create their own chain would have to completely regenerate any and all sectors drawing randomness since to use for their fork's power).

We break this down as follows:

- The miner should draw from `X-F`.
- The verifier wants to find what `X-F` should have been (to ensure the miner is not drawing from farther back) even though Y (i.e. the round of the ticket actually used) is an unverifiable value.
- Thus, the verifier will need to make an inference about what `X-F` is likely to have been based on:
  - (known) round in which the message is received (Z)
  - (known) finality value (F)
  - (approximate) SEAL time (T)
- Because T is an approximate value, and to account for network delay and variance in SEAL time across miners, the verifier allows for G offset from the assumed value of `X-F`: `Z-T-F`, hence verifying that the ticket is drawn from the range `[Z-T-F-G, Z-T-F+G]`.
