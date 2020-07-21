---
title: Formats and Serialization 
weight: 4
dashboardAudit: 1
dashboardState: wip
---

# Data Formats and Serialization
---

Filecoin seeks to make use of as few data formats as needed, with well-specced serialization rules to
better protocol security through simplicity and enable interoperability amongst implementations of the 
Filecoin protocol.

Read more on design considerations [here for CBOR-usage](https://github.com/filecoin-project/specs/issues/621) and [here for int types in Filecoin](https://github.com/filecoin-project/specs/issues/615).

## Data Formats

Filecoin in-memory data types are mostly straightforward.
Implementations should support two integer types: Int (meaning native 64-bit integer), and BigInt (meaning arbitrary length)
and avoid dealing with floating-point numbers to minimize interoperability issues across programming languages and implementations.

You can also read more on [data formats as part of randomness generation](randomness) in the Filecoin protocol.

## Serialization

Data `Serialization` in Filecoin ensures a consistent format for serializing in-memory data for transfer
in-flight and in-storage. Serialization is critical to protocol security and interoperability across
implementations of the Filecoin protocol, enabling consistent state updates across Filecoin nodes.

All data structures in Filecoin are [CBOR](https://tools.ietf.org/html/rfc7049)-tuple encoded.
That is, any data structures used in the Filecoin system (structs in this spec) should be serialized
as CBOR-arrays with items corresponding to the data structure fields in their order of declaration.

You can find the encoding structure for major data types in CBOR [here](https://tools.ietf.org/html/rfc7049#section-2.1).

For illustration, an in-memory map would be represented as a CBOR-array of the keys and values listed in some
pre-determined order. A near-term update to the serialization format will involve tagging fields appropriately
to ensure appropriate serialization/deserialization as the protocol evolves.