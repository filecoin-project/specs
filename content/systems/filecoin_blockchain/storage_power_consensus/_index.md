---
title: Storage Power Consensus
weight: 4
bookCollapseSection: true
dashboardWeight: 2
dashboardState: wip
dashboardAudit: wip
dashboardTests: 0
---

# Storage Power Consensus

TODO: remove all stale .id, .go files referenced

The Storage Power Consensus subsystem is the main interface which enables Filecoin nodes to agree on the state of the system. SPC accounts for individual storage miners' effective power over consensus in given chains in its [Power Table](storage_power_actor#the-power-table). It also runs [Expected Consensus](expected_consensus) (the underlying consensus algorithm in use by Filecoin), enabling storage miners to run leader election and generate new blocks updating the state of the Filecoin system.

Succinctly, the SPC subsystem offers the following services:

- Access to the [Power Table](storage_power_actor#the-power-table) for every subchain, accounting for individual storage miner power and total power on-chain.
- Access to [Expected Consensus](expected_consensus) for individual storage miners, enabling:

    - Access to verifiable randomness [Tickets](storage_power_consensus#tickets) as provided by [drand](drand) for the rest of the protocol.
    - Running [Leader Election](expected_consensus#secret-leader-election) to produce new blocks.
    - Running [Chain Selection](expected_consensus#chain-selection) across subchains using EC's weighting function.
    - Identification of [the most recently finalized tipset](expected_consensus#finality-in-ec), for use by all protocol participants.

Much of the Storage Power Consensus' subsystem functionality is detailed in the code below but we touch upon some of its behaviors in more detail.

{{<embed src="storage_power_consensus_subsystem.id" lang="go">}}

## Distinguishing between storage miners and block miners

There are two ways to earn Filecoin tokens in the Filecoin network:

- By participating in the [Storage Market](storage_market) as a storage provider and being paid by clients for file storage deals.
- By mining new blocks, extending the blockchain, securing the Filecoin consensus mechanism, and running smart contracts to perform state updates.

There are two types of "miners" (storage and block miners) to be distinguished. [Leader Election](expected_consensus#secret-leader-election) in Filecoin is predicated on a miner's storage power. Thus, while all block miners will be storage miners, the reverse is not necessarily true.

However, given Filecoin's "useful Proof-of-Work" is achieved through file storage ([PoRep](porep) and [PoSt](post)), there is little overhead cost for storage miners to participate in leader election. Such a [Storage Miner Actor](storage_miner_actor) need only register with the [Storage Power Actor](storage_power_actor) in order to participate in Expected Consensus and mine blocks.

## On Power
Quality-adjusted power is assigned to every sector as a static function of its _Sector Quality_ which includes `SectorSize`, `Duration`, and `DealWeight`. DealWeight is a measure that maps size and duration of active deals in a sector during its lifetime to its impact on power and reward distribution. Concretely, deal weight is defined as spacetime occupied by a deal type in a sector. A CommittedCapacity Sector (see Sector Types in [Storage Mining Subsystem](storage_mining)) will have a DealWeight of zero but all sectors have an explicit Duration which is defined from the ChainEpoch that the sector comes online in a ProveCommit message to the Expiration ChainEpoch of the sector. 

Quality-adjusted power is the number of votes a miner has in leader election and has been defined to increase linearly with the useful storage that a miner has committed to the network. 

The weight or quality of a sector depends on the deal made over the data inside the sector. There are generally three types of deals: the Committed Capacity (CC), where there is effectively no deal and the miner is storing arbitrary data inside the sector, the Regular Deals, where a miner and a client agree on a price in the market and the Verified Client deals, which give more power to the sector. We refer the reader to the [Sector](sector) and [Sector Quality](sector#sector_quality) section for details on Sector Types and Sector Quality, the [Verified Clients](verified_clients) section for more details on what a verified client is, and the [CryptoEconomics](cryptoecon) section for specific parameter values on the Deal Weights and Quality Multipliers. Sector quality multiplier of a sector is an average of deal quality multipliers weighted by the amount of spacetime each type of deal occupies in the sector.

More precisely,

- Raw-byte power: size of a sector in bytes.
- Quality-adjusted power: consensus power of stored data on the network, equal to Raw-byte power multiplied by Sector Quality Multiplier.

## Beacon Entries

The Filecoin protocol uses randomness produced by a [drand](drand) beacon to seed unbiasable randomness seeds for use in the chain (see [randomness](randomness)).

In turn these random seeds are used by:

- The [sector_sealer](sealing) as SealSeeds to bind sector commitments to a given subchain.
- The [post_generator](poster) as PoStChallenges to prove sectors remain committed as of a given block.
- The Storage Power subsystem as randomness in [leader_election](election_post) to determine their eligibility to mine a block.

This randomness may be drawn from various Filecoin chain epochs by the respective protocols that use them according to their security requirements.

It is important to note that a given Filecoin network and a given drand network
need not have the same round time, i.e. blocks may be generated faster or slower
by Filecoin than randomness is generated by drand.  For instance, if the drand
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
requires randomness for the Filecoin VM, the following method shows how to
extract the drand entry from the chain.
Note that the round may span multiple filecoin epochs if drand is slower; the
lowest epoch number block will contain the requested beacon entry. Similarly, if
there has been null rounds where the beacon should have been inserted, we need
to iterate on the chain to find where the entry is inserted.

```go
func GetRandomnessFromBeacon(e ChainEpoch, head ChainEpoch) (DrandEntry,error) {
  // get the drand round associated with the timestamp of this epoch.
  drandRound := MaxBeaconRoundForEpoch(e)
  // get the minimum drand timestamp associated with the drand round
  drandTs := drandGenesisTime + (drandPeriod-1) * drandRound 
  // get the minimum filecoin epoch associated with this timestamp
  minEpoch := (drandTs - filGenesisTime) / filEpochDuration 
  for minEpoch < head {
     // if this is not a null block, then it must have the entry we want
    if !chain.IsNullBlock(minEpoch) 
         // the requested drand entry must be in the list of drand entries
         // included in this block. If it is not the case, 
         // it means the block is invalid - but this condition is caught by the 
         // block validation logic.
         returns getDrandEntryFromBlockHeader(chain.Block(minEpoch))
    // otherwise, we need to continue progressing on the chain, i.e. maybe no
    // miner were elected or filecoin / drand outage
    minEpoch++
  }
}

func getDrandEntryFromBlockHeader(block,round) (DrandEntry,error) {
    for _,dr := range block.DrandEntries {
        if dr.Round == round {
            return dr
        }
    }
    return errors.New("drand entry not found in block")
}
```

### Fetch randomness from drand network

When mining, a miner can fetch entries from the drand network to include them in
the new block by calling the method `GetBeaconEntriesForEpoch`.

```go
GetBeaconEntriesForEpoch(epoch) []BeaconEntry {

    // special case genesis: the genesis block is pre-generated and so cannot include a beacon entry 
    // (since it will not have been generated). Hence, we only start checking beacon entries at the first block after genesis.
    // If that block includes a wrong beacon entry, we simply assume that a majority of honest miners at network birth will
    // simply fork.
    entries := []
    if epoch == 0 {
        return entries
    }

    maxDrandRound := MaxBeaconRoundForEpoch(epoch)

    // if checking the first post-genesis block, simply fetch the latest entry.
    if epoch == 1 {
        rand := drand.Public(maxDrandRound)
        return append(entries, rand)
    }

    // for the rest, fetch all drand entries generated between this epoch and last
    prevMaxDrandRound := MaxBeaconRoundForEpoch(epoch - 1)
    if (maxDrandRound == prevMaxDrandRound) {
        // no new beacon randomness
        return entries
    }

    entries := []
    curr := maxDrandRound
    for curr > prevMaxDrandRound {
        rand := drand.Public(curr)
        entries = append(entries, rand)
        curr -= 1
    }
    // return entries in increasing order
    reverse(entries)
    return entries
}
```

### Validating Beacon Entries on block reception

Per the above, a Filecoin chain will contain the entirety of the beacon's output from the Filecoin genesis to the current block.

Given their role in leader election and other critical protocols in Filecoin, a block's beacon entries must be validated for every block. See [drand](drand) for details. This can be done by ensuring every beacon entry is a valid signature over the prior one in the chain, using drand's [`Verify`](https://github.com/drand/drand/blob/763e9a252cf59060c675ced0562e8eba506971c1/chain/beacon.go#L76) endpoint as follows:

```go
// This need not be done for the genesis block
// We assume that blockHeader and priorBlockHeader are two valid subsequent headers where block was mined atop priorBlock
ValidateBeaconEntries(blockHeader, priorBlockHeader) error {
    currEntries := blockHeader.BeaconEntries
    prevEntries := priorBlockHeader.BeaconEntries

    // special case for genesis block (it has no beacon entry and so the first
    verifiable value comes at height 2, 
    // as with GetBeaconEntriesForEpoch()
    if priorBlockHeader.Epoch == 0 {
        return nil
    }

    maxRoundForEntry := MaxBeaconRoundForEpoch(blockHeader.Epoch)
    // ensure entries are not repeated in blocks
    lastBlocksLastEntry := prevEntries[len(prevEntries)-1]
    if lastBlocksLastEntry == maxRound && len(currEntries) != 0 {
        return errors.New("Did not expect a new entry in this round.")
    }

    // preparing to check that entries properly follow one another
    var entries []BeaconEntry
    // at currIdx == 0, must fetch last Fil block's last BeaconEntry
    entries := append(entries, lastBlocksLastEntry)
    entries := append(entries, currEntries...)

    currIdx := len(entries) - 1
    // ensure that the last entry in the header is not in the future (i.e. that this is not a Filecoin
    // block being mined with a future known drand entry).
    if entries[currIdx].Round != maxRoundForEntry {
        return fmt.Errorf("expected final beacon entry in block to be at round %d, got %d", maxRound, last.Round)
    }

    for currIdx >= 0 {
        // walking back the entries to ensure they follow one another
        currEntry := entries[currIdx]
        prevEntry := entries[currIdx - 1]
        err := drand.Verify(node.drandPubKey, prevEntry.Data, currEntry.Data, currEntry.Round)
        if err != nil {
            return err
        }
        currIdx -= 1
    }

    return nil
}
```

## Tickets

Filecoin block headers also contain a single "ticket" generated from its epoch's beacon entry. Tickets are used to break ties in the Fork Choice Rule, for forks of equal weight.

Whenever comparing tickets in Filecoin, the comparison is that of the ticket's VRFDigest's bytes.

### Randomness Ticket generation

At a Filecoin epoch n, a new ticket is generated using the appropriate beacon entry for epoch n.

The miner runs the beacon entry through a Verifiable Random Function (VRF) to get a new unique ticket. The beacon entry is prepended with the ticket domain separation tag and concatenated with the miner actor address (to ensure miners using the same worker keys get different tickets).

To generate a ticket for a given epoch n:
```text
randSeed = GetRandomnessFromBeacon(n)
newTicketRandomness = VRF_miner(H(TicketProdDST || index || Serialization(randSeed, minerActorAddress)))
```

We use the VRF from [Verifiable Random Functions](vrf) for ticket generation (see the `PrepareNewTicket` method below).

{{< embed src="../../filecoin_mining/storage_mining/storage_mining_subsystem.go" lang="go" >}}


### Ticket Validation

Each Ticket should be generated from the prior one in the VRF-chain and verified accordingly as shown in `validateTicket` below.

{{< embed src="storage_power_consensus_subsystem.id" lang="go" >}}
{{< embed src="storage_power_consensus_subsystem.go" lang="go" >}}

## Minimum Miner Size

In order to secure Storage Power Consensus, the system defines a minimum miner size required to participate in consensus.

Specifically, miners must have either at least `MIN_MINER_SIZE_STOR` of power (i.e. storage power currently used in storage deals) in order to participate in leader election. If no miner has `MIN_MINER_SIZE_STOR` or more power, miners with at least as much power as the smallest miner in the top `MIN_MINER_SIZE_TARG` of miners (sorted by storage power) will be able to participate in leader election. In plain english, take `MIN_MINER_SIZE_TARG = 3` for instance, this means that miners with at least as much power as the 3rd largest miner will be eligible to participate in consensus.

Miners smaller than this cannot mine blocks and earn block rewards in the network. Their power will still be counted in the total network (raw or claimed) storage power, even though their power will not be counted as votes for leader election. However, **it is important to note that such miners can still have their power faulted and be penalized accordingly**.

Accordingly, to bootstrap the network, the genesis block must include miners, potentially just CommittedCapacity sectors, to initiate the network.

The `MIN_MINER_SIZE_TARG` condition will not be used in a network in which any miner has more than `MIN_MINER_SIZE_STOR` power. It is nonetheless defined to ensure liveness in small networks (e.g. close to genesis or after large power drops).

## Network recovery after halting

Placeholder where we will define a means of rebooting network liveness after it halts catastrophically (i.e. empty power table).
