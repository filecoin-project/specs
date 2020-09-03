---
title: IPFS
bookCollapseSection: true
weight: 2
dashboardWeight: 1
dashboardState: stable
dashboardAudit: wip
dashboardTests: 0
---

# IPFS

Although Filecoin borrows many concepts from the InterPlanetary File System (IPFS), such as content addressing, the CID (which, strictly speaking, is part of the Multiformats specification) and Merkle-DAGs (which is part of IPLD), it only makes use of one of its protocols, that is, `Bitswap` and a file format, that is `UnixFS`.

## Bitswap

[Bitswap](https://github.com/ipfs/go-bitswap) is a simple peer-to-peer data exchange protocol, used primarily in IPFS, but which can be used independently of the rest of the pieces that make up IPFS. In Filecoin, `Bitswap` is used to request and receive blocks, when a node is synchonized ("caught up"), but `GossipSub` has failed to deliver some blocks to a node. 

Please refer to the [Bitswap specification](https://github.com/ipfs/specs/blob/master/BITSWAP.md) for more information.


## UnixFS

[UnixFS](https://github.com/ipfs/go-unixfs) is a protocol buffers-based format for describing files, directories, and symlinks in IPFS. `UnixFS` is used in Filecoin as a file formatting guideline for files submitted to the Filecoin network.

Please refer to the [UnixFS specification](https://github.com/ipfs/specs/blob/master/UNIXFS.md) for more information.

