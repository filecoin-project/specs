---
title: System Decomposition
---

# System Decomposition
---

## Implementing Systems

In order to make it easier to decouple functionality into systems, the Filecoin Protocol assumes
a set of functionality available to all systems. This functionality can be achieved by implementations
in a variety of ways, and should take the guidance here as a recommendation (SHOULD).

### System Requirements

All Systems, as defined in this document, require the following:

- **Repository:**
  - **Local `IpldStore`.** Some amount of persistent local storage for data structures (small structured objects).
    Systems expect to be initialized with an IpldStore in which to store data structures they expect to persist across crashes.
  - **User Configuration Values.** A small amount of user-editable configuration values.
    These should be easy for end-users to access, view, and edit.
  - **Local, Secure `KeyStore`.** A facility to use to generate and use cryptographic keys, which MUST remain secret to the
    Filecoin Node. Systems SHOULD NOT access the keys directly, and should do so over an abstraction (ie the `KeyStore`) which
    provides the ability to Encrypt, Decrypt, Sign, SigVerify, and more.
- **Local `FileStore`.** Some amount of persistent local storage for files (large byte arrays).
  Systems expect to be initialized with a FileStore in which to store large files.
  Some systems (like Markets) may need to store and delete large volumes of smaller files (1MB - 10GB).
  Other systems (like Storage Mining) may need to store and delete large volumes of large files (1GB - 1TB).
- **Network.** Most systems need access to the network, to be able to connect to their counterparts in other Filecoin Nodes.
  Systems expect to be initialized with a `libp2p.Node` on which they can mount their own protocols.
- **Clock.** Some systems need access to current network time, some with low tolerance for drift.
  Systems expect to be initialized with a Clock from which to tell network time. Some systems (like Blockchain)
  require very little clock drift, and require _secure_ time.

For this purpose, we use the `FilecoinNode` data structure, which is passed into all systems at initialization:

{{< embed src="../../systems/filecoin_nodes/node_base/filecoin_node.id" lang="go" >}}

{{< embed src="../../systems/filecoin_nodes/repository/repository_subsystem.id" lang="go" >}}

### System Limitations

Further, Systems MUST abide by the following limitations:

- **Random crashes.** A Filecoin Node may crash at any moment. Systems must be secure and consistent through crashes.
  This is primarily achived by limiting the use of persistent state, persisting such state through Ipld data structures,
  and through the use of initialization routines that check state, and perhaps correct errors.
- **Isolation.** Systems must communicate over well-defined, isolated interfaces. They must not build their critical
  functionality over a shared memory space. (Note: for performance, shared memory abstractions can be used to power
  IpldStore, FileStore, and libp2p, but the systems themselves should not require it.) This is not just an operational
  concern; it also significantly simplifies the protocol and makes it easier to understand, analyze, debug, and change.
- **No direct access to host OS Filesystem or Disk.** Systems cannot access disks directly -- they do so over the FileStore
  and IpldStore abstractions. This is to provide a high degree of portability and flexibility for end-users, especially
  storage miners and clients of large amounts of data, which need to be able to easily replace how their Filecoin Nodes
  access local storage.
- **No direct access to host OS Network stack or TCP/IP.** Systems cannot access the network directly -- they do so over the
  libp2p library. There must not be any other kind of network access. This provides a high degree of portability across
  platforms and network protocols, enabling Filecoin Nodes (and all their critical systems) to run in a wide variety of
  settings, using all kinds of protocols (eg Bluetooth, LANs, etc).


## What are systems? How do they work?
---

Filecoin decouples and modularizes functionality into loosely-joined `systems`.
Each system adds significant functionality, usually to achieve a set of important and tightly related goals.

For example, the Blockchain System provides structures like Block, Tipset, and Chain, and provides functionality
like Block Sync, Block Propagation, Block Validation, Chain Selection, and Chain Access. This is
separated from the Files, Pieces, Piece Preparation, and Data Transfer. Both of these systems are separated from
the Markets, which provide Orders, Deals, Market Visibility, and Deal Settlement.

### Why is System decoupling useful?

This decoupling is useful for:

- **Implementation Boundaries:** it is possible to build implementations of Filecoin that only implement a
  subset of systems. This is especially useful for _Implementation Diversity_: we want many implementations
  of security critical systems (eg Blockchain), but do not need many implementations of Systems that can be
  decoupled.
- **Runtime Decoupling:** system decoupling makes it easier to build and run Filecoin Nodes that isolate
  Systems into separate programs, and even separate physical computers.
- **Security Isolation:** some systems require higher operational security than others. System decoupling allows
  implementations to meet their security and functionality needs. A good example of this is separating Blockchain
  processing from Data Transfer.
- **Scalability:** systems, and various use cases, may drive different performance requirements for different opertators.
  System decoupling makes it easier for operators to scale their deployments along system boundaries.


### Filecoin Nodes don't need all the systems

Filecoin Nodes vary significantly, and do not need all the systems.
Most systems are only needed for a subset of use cases.

For example, the Blockchain System is required for synchronizing the chain, participating in secure consensus,
storage mining, and chain validation.
Many Filecoin Nodes do not need the chain and can perform their work by just fetching content from the latest
StateTree, from a node they trust.

Note: Filecoin does not use the "full node" or "light client" terminology, in wide use in Bitcoin and other blockchain
networks. In filecoin, these terms are not well defined. It is best to define nodes in terms of their capabilities,
and therefore, in terms of the Systems they run. For example:

- **Chain Verifier Node:** Runs the Blockchain system. Can sync and validate the chain. Cannot mine or produce blocks.
- **Client Node:** Runs the Blockchain, Market, and Data Transfer systems. Can sync and validate the chain. Cannot mine or produce blocks.
- **Retrieval Miner Node:** Runs the Market and Data Transfer systems. Does not need the chain. Can make Retrieval Deals
  (Retrieval Provider side). Can send Clients data, and get paid for it.
- **Storage Miner Node:** Runs the Blockchain, Storage Market, Storage Mining systems. Can sync and validate the chain.
  Can make Storage Deals (Storage Provider side). Can seal stored data into sectors. Can acquire
  storage consensus power. Can mine and produce blocks.

### Separating Systems

> How do we determine what functionality belongs in one system vs another?

Drawing boundaries between systems is the art of separating tightly related functionality from unrelated parts.
In a sense, we seek to keep tightly integrated components in the same system, and away from other unrelated
components. This is sometimes straightforward, the boundaries naturally spring from the data structures or
functionality. For example, it is straightforward to observe that Clients and Miners negotiating a deal
with each other is very unrelated to VM Execution.

Sometimes this is harder, and it requires detangling, adding, or removing abstractions. For
example, the `StoragePowerActor` and the `StorageMarketActor` were a single `Actor` previously. This caused
a large coupling of functionality across `StorageDeal` making, the `StorageMarket`, markets in general, with
Storage Mining, Sector Sealing, PoSt Generation, and more. Detangling these two sets of related functionality
requried breaking apart the one actor into two.

### Decomposing within a System

Systems themselves decompose into smaller subunits. These are sometimes called "subsystems" to avoid confusion
with the much larger, first-class Systems. Subsystems themselves may break down further. The naming here is not
strictly enforced, as these subdivisions are more related to protocol and implementation engineering concerns
than to user capabilities.