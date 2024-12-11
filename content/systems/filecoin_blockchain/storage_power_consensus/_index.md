---
title: Storage Power Consensus
weight: 4
bookCollapseSection: true
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Storage Power Consensus

The Storage Power Consensus (SPC) subsystem is the main interface which enables Filecoin nodes to agree on the state of the system. Storage Power Consensus accounts for individual storage miners' effective power over consensus in given chains in its [Power Table](storage_power_actor#the-power-table). It also runs [Expected Consensus](expected_consensus) (the underlying consensus algorithm in use by Filecoin), enabling storage miners to run leader election and generate new blocks updating the state of the Filecoin system.

Succinctly, the SPC subsystem offers the following services:

- Access to the [Power Table](storage_power_actor#the-power-table) for every subchain, accounting for individual storage miner power and total power on-chain.
- Access to [Expected Consensus](expected_consensus) for individual storage miners, enabling:

  - Access to verifiable randomness [Tickets](storage_power_consensus#tickets) as provided by [drand](drand) for the rest of the protocol.
  - Running [Leader Election](expected_consensus#secret-leader-election) to produce new blocks.
  - Running [Chain Selection](expected_consensus#chain-selection) across subchains using EC's weighting function.
  - Identification of [the most recently finalized tipset](expected_consensus#finality-in-ec), for use by all protocol participants.

## Distinguishing between storage miners and block miners

There are two ways to earn Filecoin tokens in the Filecoin network:

- By participating in the [Storage Market](storage_market) as a storage provider and being paid by clients for file storage deals.
- By mining new blocks, extending the blockchain, securing the Filecoin consensus mechanism, and running smart contracts to perform state updates as a [Storage Miner](filecoin_mining#storage_mining).

There are two types of "miners" (storage and block miners) to be distinguished. [Leader Election](expected_consensus#secret-leader-election) in Filecoin is predicated on a miner's storage power. Thus, while all block miners will be storage miners, the reverse is not necessarily true.

However, given Filecoin's "useful Proof-of-Work" is achieved through file storage ([PoRep](porep) and [PoSt](post)), there is little overhead cost for storage miners to participate in leader election. Such a [Storage Miner Actor](storage_miner_actor) need only register with the [Storage Power Actor](storage_power_actor) in order to participate in Expected Consensus and mine blocks.

## On Power

Quality-adjusted power is assigned to every sector as a static function of its **_Sector Quality_** which includes: i) the **Sector Spacetime**, which is the product of the sector size and the promised storage duration, ii) the **Deal Weight** that converts spacetime occupied by deals into consensus power, iii) the **Deal Quality Multiplier** that depends on the type of deal done over the sector (i.e., CC, Regular Deal or Verified Client Deal), and finally, iv) the **Sector Quality Multiplier**, which is an average of deal quality multipliers weighted by the amount of spacetime each type of deal occupies in the sector.

The **Sector Quality** is a measure that maps size, duration and the type of active deals in a sector during its lifetime to its impact on power and reward distribution.

The quality of a sector depends on the deals made over the data inside the sector. There are generally three types of deals: the _Committed Capacity (CC)_, where there is effectively no deal and the miner is storing arbitrary data inside the sector, the _Regular Deals_, where a miner and a client agree on a price in the market and the _Verified Client_ deals, which give more power to the sector. We refer the reader to the [Sector](sector) and [Sector Quality](sector#sector-quality) sections for details on Sector Types and Sector Quality, the [Verified Clients](verified_clients) section for more details on what a verified client is, and the [CryptoEconomics](cryptoecon) section for specific parameter values on the Deal Weights and Quality Multipliers.

**Quality-Adjusted Power** is the number of votes a miner has in the [Secret Leader Election](expected_consensus#secret-leader-election) and has been defined to increase linearly with the useful storage that a miner has committed to the network.

More precisely, we have the following definitions:

- _Raw-byte power_: the size of a sector in bytes.
- _Quality-adjusted power_: the consensus power of stored data on the network, equal to Raw-byte power multiplied by the Sector Quality Multiplier.

## Beacon Entries

The Filecoin protocol uses randomness produced by a [drand](drand) beacon to seed unbiasable randomness seeds for use in the chain (see [randomness](randomness)).

In turn these random seeds are used by:

- The [sector_sealer](sealing) as SealSeeds to bind sector commitments to a given subchain.
- The [post_generator](poster) as PoStChallenges to prove sectors remain committed as of a given block.
- The Storage Power subsystem as randomness in [Secret Leader Election](expected_consensus#secret-leader-election) to determine how often a miner is chosen to mine a new block.

This randomness may be drawn from various Filecoin chain epochs by the respective protocols that use them according to their security requirements.

It is important to note that a given Filecoin network and a given drand network
need not have the same round time, i.e. blocks may be generated faster or slower
by Filecoin than randomness is generated by drand. For instance, if the drand
beacon is producing randomness twice as fast as Filecoin produces blocks, we
might expect two random values to be produced in a Filecoin epoch, conversely if
the Filecoin network is twice as fast as drand, we might expect a random value
every other Filecoin epoch. Accordingly, depending on both networks'
configurations, certain Filecoin blocks could contain multiple or no drand
entries.
Furthermore, it must be that any call to the drand network for a new randomness
entry during an outage should be blocking, as noted with the `drand.Public()`
calls below.
In all cases, Filecoin blocks must include all drand beacon outputs generated
since the last epoch in the `BeaconEntries` field of the block header. Any use
of randomness from a given Filecoin epoch should use the last valid drand entry
included in a Filecoin block. This is shown below.

### Get drand randomness for VM

For operations such as PoRep creation, proof validations, or anything that
requires randomness for the Filecoin VM, there should be a method that
extracts the drand entry from the chain correctly.
Note that the round may span multiple filecoin epochs if drand is slower; the
lowest epoch number block will contain the requested beacon entry. Similarly, if
there has been null rounds where the beacon should have been inserted, we need
to iterate on the chain to find where the entry is inserted. Specifically, the next non-null block must contain the drand entry requested by definition.

### Fetch randomness from drand network

When mining, a miner can fetch entries from the drand network to include them in
the new block.

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/beacon/drand/drand.go" lang="go" symbol="DrandBeacon">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/beacon/beacon.go" lang="go" symbol="BeaconEntriesForBlock">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/beacon/drand/drand.go" lang="go" symbol="MaxBeaconRoundForEpoch">}}

### Validating Beacon Entries on block reception

A Filecoin chain will contain the entirety of the beacon's output from the Filecoin genesis to the current block.

Given their role in leader election and other critical protocols in Filecoin, a block's beacon entries must be validated for every block. See [drand](drand) for details. This can be done by ensuring every beacon entry is a valid signature over the prior one in the chain, using drand's [`Verify`](https://github.com/drand/drand/blob/763e9a252cf59060c675ced0562e8eba506971c1/chain/beacon.go#L76) endpoint as follows:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/beacon/beacon.go" lang="go" symbol="ValidateBlockValues">}}

## Tickets

Filecoin block headers also contain a single "ticket" generated from its epoch's beacon entry. Tickets are used to break ties in the Fork Choice Rule, for forks of equal weight.

Whenever comparing tickets in Filecoin, the comparison is that of the ticket's VRF Digest's bytes.

### Randomness Ticket generation

At a Filecoin epoch `n`, a new ticket is generated using the appropriate beacon entry for epoch `n`.

The miner runs the beacon entry through a Verifiable Random Function (VRF) to get a new unique ticket. The beacon entry is prepended with the ticket domain separation tag and concatenated with the miner actor address (to ensure miners using the same worker keys get different tickets).

To generate a ticket for a given epoch n:

```text
randSeed = GetRandomnessFromBeacon(n)
newTicketRandomness = VRF_miner(H(TicketProdDST || index || Serialization(randSeed, minerActorAddress)))
```

[Verifiable Random Functions](vrf) are used for ticket generation.

### Ticket Validation

Each Ticket should be generated from the prior one in the VRF-chain and verified accordingly.

## Minimum Miner Size

In order to secure Storage Power Consensus, the system defines a minimum miner size required to participate in consensus.

Specifically, miners must have either at least `MIN_MINER_SIZE_STOR` of power (i.e. storage power currently used in storage deals) in order to participate in leader election. If no miner has `MIN_MINER_SIZE_STOR` or more power, miners with at least as much power as the smallest miner in the top `MIN_MINER_SIZE_TARG` of miners (sorted by storage power) will be able to participate in leader election. In plain english, take `MIN_MINER_SIZE_TARG = 3` for instance, this means that miners with at least as much power as the 3rd largest miner will be eligible to participate in consensus.

Miners smaller than this cannot mine blocks and earn block rewards in the network. Their power will still be counted in the total network (raw or claimed) storage power, even though their power will not be counted as votes for leader election. However, **it is important to note that such miners can still have their power faulted and be penalized accordingly**.

Accordingly, to bootstrap the network, the genesis block must include miners, potentially just CommittedCapacity sectors, to initiate the network.

The `MIN_MINER_SIZE_TARG` condition will not be used in a network in which any miner has more than `MIN_MINER_SIZE_STOR` power. It is nonetheless defined to ensure liveness in small networks (e.g. close to genesis or after large power drops).
