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

For {{<sref crypto_signatures>}} as well as {{<vrf>}} usage in the protocol, we define the following
Domain Separation Tags, prepending random inputs with bytes corresponding to the little endian uint64
representations of the following numbers:
- for generating a new random ticket:           TicketDST   = `1`
- for generating randomness for ElectionPoSt:   PoStDST     = `2`

{{< readfile file="domain.go" code="true" lang="go" >}}
{{< readfile file="randomness.go" code="true" lang="go" >}}

{{% notice todo %}}
**TODO**: sync with lotus about the following and accordingly, remove ot keep. It is currently not part of the 
Filecoin spec.

## Input Delimeters

As with Filecoin's use of hash function we also define an input delimeter to separate
concatenated elements for input to a signature (or hash), as a byte representation of the following:
- inputDelimeter = `0`
{{% /notice %}}