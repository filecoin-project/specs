---
title: Multiformats
weight: 3
dashboardWeight: 1
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# Multiformats

[Multiformats](https://multiformats.io/) is a set of self-describing protocol values. These values are useful both to the data layer (IPLD) and to the network layer (libp2p). Multiformats includes specifications for the Content Identifier (CID) used by IPLD and IPFS, the multicodec, multibase and multiaddress (used by libp2p).

Please refer to the [Multiformats repository](https://github.com/multiformats) for more information.

## CIDs

Filecoin references data using IPLD's Content Identifier (CID).

A CID is a hash digest prefixed with identifiers for its hash function and codec. This means you can validate and decode data by with only this identifier.

When CIDs are printed as strings they also use multibase to identify the base encoding being used.

For a more detailed specification, we refer the reader to the
[CID specification](https://github.com/multiformats/cid).

## Multihash

A Multihash is a set of self-describing hash values. Multihash is used for differentiating outputs from various well-established cryptographic hash functions, addressing sizes and encoding considerations.

Please refer to the [Multihash specification](https://github.com/multiformats/multihash) for more information.

## Multiaddr

A Multiadddress is a self-describing network address. Multiaddresses are composable and future-proof network addresses used by libp2p.

Please refer to the [Multiaddr specification](https://github.com/multiformats/multiaddr) for more information.
