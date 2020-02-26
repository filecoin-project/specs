---
title: "Randomness"
---

{{<label randomness>}}

Randomness is used throughout the protocol in order to generate values and extend the blockchain.
Random values are drawn from the {{<sref ticket_chain>}} and appropriately formatted for usage.
We describe this formatting below.

## Encoding On-chain data for randomness

Entropy from the ticket-chain can be combined with other values to generate necessary randomness that can be
specific to (eg) a given miner address or epoch. To be used as part of entropy, these values are combined in 
objects that can then be CBOR-serialized according to their algebraic datatypes.

## Domain Separation Tags

Further, we define Domain Separation Tags with which we prepend random inputs when creating entropy.

All randomness used in the protocol must be generated in conjunction with a unique DST, as well as 
certain {{<sref crypto_signatures>}} and {{<sref vrf>}} usage.

## Drawing Randomness from the chain

Tickets are used as a source of on-chain randomness, generated with each new block created (see {{<sref tickets>}}).

Randomness is derived from a ticket's digest as follows, for a given epoch `n`, and ticket sought at epoch `e`:
```text
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
newRandomness = H(ticket.Digest() || e)
return newRandomness
```

In plain language, this means:

- Choose the smallest ticket in the Tipset if it contains multiple blocks.
- When sampling a ticket from an epoch with no blocks, draw the min ticket from the prior epoch with blocks

The ticket digest is concatenated with the wanted epoch number in order to ensure both:
- liveness for leader election -- in the case of null rounds, the new epoch number will output new randomness for LE
- distinct values for randomness sought before genesis -- where the genesis ticket will be returned

For instance, if in epoch `curr`, a miner wants randomness from `lookback` epochs back where `curr - lookback <= genesis`, 
the ticket randomness drawn would be `H(genesisTicket.digest || curr-lookback)` where the `genesisTicket` is the randomness included
in the genesis block (to be determined ahead of time to enable genesis participants to SEAL data ahead of time).

While the inclusion of the round number is not necessary for other uses of entropy, we include it as part of ticket drawing for uniformity/simplicity in the overall protocol.

See the `RandomnessAtEpoch` method below:
{{< readfile file="../struct/chain/chain.go" code="true" lang="go" >}}

## Forming Randomness Seeds

Different uses of randomness require randomness seeds predicated on a variety of inputs. For instance, we have:

- `TicketProduction` -- uses ticket and miner actor addr
- `ElectionPoStChallengeSeed` -- uses ticket and miner actor addr
- `WindowedPoStChallengeSeed` -- uses ticket and epoch number
...

In all cases, a ticket is used as the base of randomness (see {{<sref tickets>}}). In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, using blake2b as a hash function to generate a 256-bit output as follows (also see {{<sref tickets>}}):

In round `n`, for a given randomness lookback `l`, and serialized entropy `s`:

```text
ticket = DrawRandomness(n-l)
randSeed = ticket.digest

buffer = Bytes{}
buffer.append(IntToBigEndianBytes(AppropriateDST))
buffer.append(randSeed)
buffer.append(s)

randomness = H(buffer)
```

{{< readfile file="/docs/actors/actors/crypto/randomness.go" code="true" lang="go" >}}
{{< readfile file="/docs/systems/filecoin_blockchain/struct/chain/chain.go" code="true" lang="go" >}}

## Entropy to be used with randomness

We currently distinguish the following entropy needs per use:

- TicketProduction: requires MinerIDAddress
- ElectionPoStChallengeSeed: requires current epoch and MinerIDAddress
- WindowedPoStChallengeSeed: requires MinerIDAddress
- SealRandomness: requires MinerIDAddress
- InteractiveSealChallengeSeed: requires MinerIDAddress