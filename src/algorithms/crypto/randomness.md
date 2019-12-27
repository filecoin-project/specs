---
title: "Randomness"
---

{{<label randomness>}}

Randomness is used throughout the protocol in order to generate values and extend the blockchain.
Random values are drawn from the {{<sref ticket_chain>}} and appropriately formatted for usage.
We describe this formatting below.

## Encoding On-chain data for randomness

Any randomness derived from on-chain values uses the following encodings to represent these values as bytes
- **Bytes**: Bytes
- **Ints**: Little endian uint64 representation
- **Strings**: ASCII
- **Objects**: Their specified Serialization, currently CBOR-based serialization defined on algebraic datatypes

## Domain Separation Tags

For {{<sref crypto_signatures>}} as well as {{<sref vrf>}} usage in the protocol, we define Domain Separation Tags with which we prepend random inputs.

The source of truth is defined below, but the currently defined DSTs are:
- for drawing randomness from an on-chain ticket:
    - `TicketDrawingDST = 1`
- for generating a new random ticket:
    - `TicketProductionDST = 2`
- for generating randomness for running PoSt (ElectionPoSt or SurprisePoSt):
    - `PoStChallengeSeedDST = 3`
- for selection of which miners to surprise:
    - `SurprisePoStSelectMinersDST = 4`

## Forming Randomness Seeds

Different uses of randomness require randomness seeds predicated on a variety of inputs. For instance, we have:
- `TicketDrawing` -- uses ticket and epoch number
- `TicketProduction` -- uses ticket and miner actor addr
- `PoStChallengeSeed` -- uses ticket and miner actor addr
- `SurprisePoStSelectMiners` -- uses ticket and epoch number

In all cases, a ticket is used as the base of randomness (see {{<sref tickets>}}). In order to make randomness seed creation uniform, the protocol derives all such seeds in the same way, as follows (also see {{<sref tickets>}}):
```text
For a given ticket's randomness ticket_randomness:
buffer = Bytes{}
buffer.append(IntToLittleEndianBytes(AppropriateDST))
buffer.append(-1) // a flag to be used in cases where FIL 
might need longer randomness outputs. Currently unused
buffer.append(ticket_randomness)
buffer.append(other needed inputs)

randomness = SHA256(buffer)
```

{{< readfile file="randomness.go" code="true" lang="go" >}}
{{< readfile file="chain.go" code="true" lang="go" >}}