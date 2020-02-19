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

## Forming Randomness Seeds

Different uses of randomness require randomness seeds predicated on a variety of inputs. For instance, we have:

- `TicketProduction` -- uses ticket and miner actor addr
- `ElectionPoStChallengeSeed` -- uses ticket and miner actor addr
- `WindowedPoStChallengeSeed` -- uses ticket and epoch number
...

In all cases, a ticket is used as the base of randomness (see {{<sref tickets>}}). In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, as follows (also see {{<sref tickets>}}):

For a given randomness lookback `l`, and serialized entropy `s`:

```text
ticket = draw_ticket_from_chain(l)
ticket_randomness = ticket.digest

buffer = Bytes{}
buffer.append(IntToBigEndianBytes(AppropriateDST))
buffer.append(ticket_randomness)
buffer.append(s)

randomness = H(buffer)
```

{{< readfile file="/docs/actors/actors/crypto/randomness.go" code="true" lang="go" >}}
{{< readfile file="/docs/systems/filecoin_blockchain/struct/chain/chain.go" code="true" lang="go" >}}

## Drawing randomness prior to genesis

Any randomness tickets drawn from farther back than genesis will be drawn using the genesis ticket concatenated with the wanted epoch number.

For instance, if in epoch `curr`, a miner wants randomness from `lookback` epochs back where `curr - lookback < genesis`, 
the ticket randomness drawn would be `H(genesisTicket.digest || curr-lookback)` where the `genesisTicket` is the randomness included
in the genesis block (to be determined ahead of time to enable genesis participants to SEAL data ahead of time).
