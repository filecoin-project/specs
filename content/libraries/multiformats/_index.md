---
title: Multiformats
weight: 4
dashboardWeight: 1
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# Multiformats - self describing protocol values

[Multiformats](https://multiformats.io/) is a set of formating rules and formatting structures for self-describing values. These values are useful both to the data layer (IPLD) and to the network layer (libp2p). Multiformats includes specifications for the Content Identifier (CID) used by IPLD and IPFS, the multicodec, multibase and multiaddress (used by libp2p).

Please refer to the [Multiformats repository](https://github.com/multiformats) for more information.

## CIDs - Content IDentifiers

Filecoin references data using IPLD's Content Identifier (CID).
A CID is effectively a hash value, prefixed with its hash function (multihash) as well as extra labels, such as a `codec` and a `multibase` to inform applications about how to deserialize the given data.

For a more detailed specification, we refer the reader to the
[CID specification](https://github.com/multiformats/cid).

## Multihash - self describing hash values

Multihash is a protocol for differentiating outputs from various well-established cryptographic hash functions, addressing sizes and encoding considerations. Please refer to the [Multihash specification](https://github.com/multiformats/multihash) for more information.

## Multiaddr - self describing network addresses

Multiaddresses are composable and future-proof network addresses used by libp2p. See the [Multiaddr specification](https://github.com/multiformats/multiaddr) for more information.
