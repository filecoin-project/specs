---
title: "Randomness"
---

TODO: clean up stale .id/.go files

{{<label randomness>}}

Randomness is used throughout the protocol in order to generate values and extend the blockchain.
Random values are drawn from a {{<sref drand>}} beacon and appropriately formatted for usage.
We describe this formatting below.

## Encoding On-chain data for randomness

Entropy from the drand beacon can be combined with other values to generate necessary randomness that can be
specific to (eg) a given miner address or epoch. To be used as part of entropy, these values are combined in 
objects that can then be CBOR-serialized according to their algebraic datatypes.

## Domain Separation Tags

Further, we define Domain Separation Tags with which we prepend random inputs when creating entropy.

All randomness used in the protocol must be generated in conjunction with a unique DST, as well as 
certain {{<sref crypto_signatures>}} and {{<sref vrf>}} usage.

## Forming Randomness Seeds

Drand randomness entries are used as a source of on-chain randomness (see {{<sref random_seed "random seeds">}}).

The random seed is combined with a few elements for use as part of the protocol as follows:

- a DST (domain separation tag)
    - Different uses of randomness are distinguished by this type of personalization which ensures that randomness used for different purposes will not conflict with randomness used elsewhere in the protocol
- the epoch number, ensuring
    - liveness for leader election -- in the case no one is elected in a round and no new drand entry has appeared, the new epoch number will output new randomness for LE
    - other entropy, ensuring that randomness is modified as needed by other context-dependent entropy (e.g. a miner address if we want the randomness to be different for each miner).

While all elements are not needed for every use of entropy (e.g. the inclusion of the round number is not necessary prior to genesis or outside of leader election, other entropy is only used sometimes, etc), we draw randomness as follows for the sake of uniformity/simplicity in the overall protocol.

In all cases, a drand entry is used as the base of randomness (see {{<sref random_seed>}}). In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, using blake2b as a hash function to generate a 256-bit output as follows:

In round `n`, for a given randomness lookback `l`, and serialized entropy `s`:

```text
GetRandomness(dst, l, s):
    ticketDigest = beacon.GetRandomnessForEpoch(n-l)

    buffer = Bytes{}
    buffer.append(IntToBigEndianBytes(dst))
    buffer.append(randSeed)
    buffer.append(n-l) // the sought epoch
    buffer.append(s)

    return H(buffer)
```

{{< readfile file="/docs/actors/actors/crypto/randomness.go" code="true" lang="go" >}}
{{< readfile file="/docs/systems/filecoin_blockchain/struct/chain/chain.go" code="true" lang="go" >}}

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