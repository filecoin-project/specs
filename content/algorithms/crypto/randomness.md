---
title: "Randomness"
weight: 3
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Randomness

TODO: clean up stale .id/.go files

Randomness is used throughout the protocol in order to generate values and extend the blockchain.
Random values are drawn from a [drand](drand) beacon and appropriately formatted for usage.
We describe this formatting below.

## Encoding Random Beacon randomness for on-chain use

Entropy from the drand beacon can be harvested into a more general data structure: a `BeaconEntry`, defined as follows:

```go
type BeaconEntry struct {
    // Drand Round for the given randomness
    Round       uint64
    // Drand Signature for the given Randomness, named Data as a more general name for random beacon output
    Data   []byte
}
```

The BeaconEntry is then combined with other values to generate necessary randomness that can be
specific to (eg) a given miner address or epoch. To be used as part of entropy, these values are combined in 
objects that can then be CBOR-serialized according to their algebraic datatypes.

## Domain Separation Tags

Further, we define Domain Separation Tags with which we prepend random inputs when creating entropy.

All randomness used in the protocol must be generated in conjunction with a unique DST, as well as 
certain [Signatures](signatures) and [Verifiable Random Function](vrf) usage.

## Forming Randomness Seeds

The beacon entry is combined with a few elements for use as part of the protocol as follows:

- a DST (domain separation tag)
    - Different uses of randomness are distinguished by this type of personalization which ensures that randomness used for different purposes will not conflict with randomness used elsewhere in the protocol
- the epoch number, ensuring
    - liveness for leader election -- in the case no one is elected in a round and no new beacon entry has appeared (i.e. if the beacon frequency is slower than that of block production in Filecoin), the new epoch number will output new randomness for LE (note that Filecoin uses liveness during a beacon outage).
    - other entropy, ensuring that randomness is modified as needed by other context-dependent entropy (e.g. a miner address if we want the randomness to be different for each miner).

While all elements are not needed for every use of entropy (e.g. the inclusion of the round number is not necessary prior to genesis or outside of leader election, other entropy is only used sometimes, etc), we draw randomness as follows for the sake of uniformity/simplicity in the overall protocol.

In all cases, a [drand](drand) signature is used as the base of randomness: it is hashed using blake2b in order to obtain a usable randomness seed. In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, using blake2b as a hash function to generate a 256-bit output as follows:

In round `n`, for a given randomness lookback `l`, and serialized entropy `s`:

```text
GetRandomness(dst, l, s):
    ticketDigest = beacon.GetRandomnessFromBeacon(n-l)

    buffer = Bytes{}
    buffer.append(IntToBigEndianBytes(dst))
    buffer.append(randSeed)
    buffer.append(n-l) // the sought epoch
    buffer.append(s)

    return H(buffer)
```


{{<embed src="/systems/filecoin_blockchain/struct/chain/chain.go" lang="go">}}

## Drawing tickets from the VRF-chain for proof inclusion

In some places, the protocol needs randomness drawn from the Filecoin blockchain's VRF-chain (which generates [tickets](storage_power_consensus#tickets) with each new block) rather than from the random beacon, in order to tie certain proofs to a particular set of Filecoin blocks (i.e. a given chain or fork).
In particular, `SealRandomness` must be taken from the VRF chain, in order to ensure that no other fork can replay the Seal (see [sealing](sealing) for more).

A ticket is drawn from the chain for randomness as follows, for a given epoch `n`, and ticket sought at epoch `e`:
```text
GetRandomnessFromVRFChain(e):
    While ticket is not set:
        Set wantedTipsetHeight = e
        if wantedTipsetHeight <= genesis:
            Set ticket = genesis ticket
        else if blocks were mined at wantedTipsetHeight:
            ReferenceTipset = TipsetAtHeight(wantedTipsetHeight)
            Set ticket = minTicket in ReferenceTipset
        If no blocks were mined at wantedTipsetHeight:
            wantedTipsetHeight--
            (Repeat)
    return ticket.Digest()
```

In plain language, this means:

- Choose the smallest ticket in the Tipset if it contains multiple blocks.
- When sampling a ticket from an epoch with no blocks, draw the min ticket from the prior epoch with blocks

This ticket is then combined with a Domain Separation Tag, the round number sought and appropriate entropy to form randomness for various uses in the protocol.

See the `GetRandomnessFromVRFChain` method below:
{{<embed src="/systems/filecoin_blockchain/struct/chain/chain.go" lang="go">}}

## Entropy to be used with randomness

As stated above, different uses of randomness may require added entropy. The CBOR-serialization of the inputs to this entropy must be used.

For instance, if using entropy from an object of type foo, its CBOR-serialization would be appended to the randomness in `GetRandomness()`. If using both an object of type foo and one of type bar for entropy, you may define an object of type baz (as below) which includes all needed entropy, and include its CBOR-serialization in `GetRandomness()`.

```text
type baz struct {
    firstObject     foo
    secondObject    bar
}
```

Currently, we distinguish the following entropy needs in the Filecoin protocol (this list is not exhaustive):

- TicketProduction: requires MinerIDAddress
- ElectionProofProduction: requires current epoch and MinerIDAddress -- epoch is already mixed in from ticket drawing so in practice is the same as just adding MinerIDAddress as entropy
- WinningPoStChallengeSeed: requires MinerIDAddress
- WindowedPoStChallengeSeed: requires MinerIDAddress
- WindowedPoStDeadlineAssignment: TODO @jake
- SealRandomness: requires MinerIDAddress
- InteractiveSealChallengeSeed: requires MinerIDAddress

The above uses of the MinerIDAddress ensures that drawn randomness is distinct for every miner drawing this (regardless of whether they share worker keys or not, eg -- in the case of randomness that is then signed as part of leader election or ticket production).
