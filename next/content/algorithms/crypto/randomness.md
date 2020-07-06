---
title: "Randomness"
---

# Randomness
---

{{< hint danger >}}
Issues with labels and readfile throughout
{{< /hint >}}


{{< hint danger >}}
Issue with label
{{< /hint >}}


{{</* label randomness */>}}

Randomness is used throughout the protocol in order to generate values and extend the blockchain.
Random values are drawn from the [Ticket Chain](\missing-link) and appropriately formatted for usage.
We describe this formatting below.

## Encoding On-chain data for randomness

Entropy from the ticket-chain can be combined with other values to generate necessary randomness that can be
specific to (eg) a given miner address or epoch. To be used as part of entropy, these values are combined in 
objects that can then be CBOR-serialized according to their algebraic datatypes.

## Domain Separation Tags

Further, we define Domain Separation Tags with which we prepend random inputs when creating entropy.

All randomness used in the protocol must be generated in conjunction with a unique DST, as well as 
certain [Signatures](\missing-link) and [Verifiable Random Function](\missing-link) usage.

## Drawing tickets for randomness from the chain

Tickets are used as a source of on-chain randomness, generated with each new block created (see [Tickets](\missing-link)).

A ticket is drawn from the chain for randomness as follows, for a given epoch `n`, and ticket sought at epoch `e`:
```text
RandomnessSeedAtEpoch(e):
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

{{< hint danger >}}
Issue with readfile
{{< /hint >}}

See the `RandomnessSeedAtEpoch` method below:
{{</* readfile file="../struct/chain/chain.go" code="true" lang="go" */>}}

## Forming Randomness Seeds

The drawn ticket digest is combined with a few elements to make up randomness for use as part of the protocol.

- a DST (domain separation tag)
    - Different uses of randomness are distinguished by this type of personalization which ensures that randomness used for different purposes will not conflict with randomness used elsewhere in the protocol
- the epoch number, ensuring
    - liveness for leader election -- in the case of null rounds, the new epoch number will output new randomness for LE
    - distinct values for randomness sought before genesis -- where the genesis ticket will be returned
    - For instance, if in epoch `curr`, a miner wants randomness from `lookback` epochs back where `curr - lookback <= genesis`, the ticket randomness drawn would be based on `genesisTicket.digest` where the `genesisTicket` is the randomness included in the genesis block. Using the epoch as part of randomness composition ensures that randomness drawn at various epochs prior to genesis has different values.
- other entropy,
    - ensuring that randomness is modified as needed by other context-dependent entropy (e.g. a miner address if we want the randomness to be different for each miner).

While all elements are not needed for every use of entropy (e.g. the inclusion of the round number is not necessary prior to genesis or outside of leader election, other entropy is only used sometimes, etc), we draw randomness as follows for the sake of uniformity/simplicity in the overall protocol.

In all cases, a ticket is used as the base of randomness (see [Tickets](\missing-link)). In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, using blake2b as a hash function to generate a 256-bit output as follows (also see [Tickets](\missing-link)):

In round `n`, for a given randomness lookback `l`, and serialized entropy `s`:

```text
GetRandomness(dst, l, s):
    ticketDigest = RandomnessSeedAtEpoch(n-l)

    buffer = Bytes{}
    buffer.append(IntToBigEndianBytes(dst))
    buffer.append(randSeed)
    buffer.append(n-l) // the sought epoch
    buffer.append(s)

    return H(buffer)
```

{{< hint danger >}}
Issue with readfile
{{< /hint >}}

{{</* readfile file="/docs/actors/actors/crypto/randomness.go" code="true" lang="go" */>}}
{{</* readfile file="/docs/systems/filecoin_blockchain/struct/chain/chain.go" code="true" lang="go" */>}}

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
- ElectionPoStChallengeSeed: requires current epoch and MinerIDAddress -- epoch is already mixed in from ticket drawing so in practice is the same as just adding MinerIDAddress as entropy
- WindowedPoStChallengeSeed: requires MinerIDAddress
- SealRandomness: requires MinerIDAddress
- InteractiveSealChallengeSeed: requires MinerIDAddress

The above uses of the MinerIDAddress ensures that drawn randomness is distinct for every miner drawing this (regardless of whether they share worker keys or not, eg -- in the case of randomness that is then signed as part of leader election or ticket production).