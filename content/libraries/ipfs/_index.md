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

Filecoin is built on the same underlying stack as IPFS - including connecting nodes peer-to-peer via [libp2p](https://libp2p.io) and addressing data using [IPLD](https://ipld.io/). Therefore, it borrows many concepts from the InterPlanetary File System (IPFS), such as content addressing, the CID (which, strictly speaking, is part of the Multiformats specification) and Merkle-DAGs (which is part of IPLD). It also makes direct use of `Bitswap` (the data transfer algorithm in IPFS) and `UnixFS` (the file format built on top of IPLD Merkle-Dags).

## Bitswap

[Bitswap](https://github.com/ipfs/go-bitswap) is a simple peer-to-peer data exchange protocol, used primarily in IPFS, which can also be used independently of the rest of the pieces that make up IPFS. In Filecoin, `Bitswap` is used to request and receive blocks when a node is synchonized ("caught up") but `GossipSub` has failed to deliver some blocks to a node. 

Please refer to the [Bitswap specification](https://github.com/ipfs/specs/blob/master/BITSWAP.md) for more information.


## UnixFS

[UnixFS](https://github.com/ipfs/go-unixfs) is a protocol buffers-based format for describing files, directories, and symlinks in IPFS. `UnixFS` is used in Filecoin as a file formatting guideline for files submitted to the Filecoin network.

Please refer to the [UnixFS specification](https://github.com/ipfs/specs/blob/master/UNIXFS.md) for more information.
