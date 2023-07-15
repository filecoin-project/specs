---
title: What are Systems?
weight: 1
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# What are Systems? How do they work?

Filecoin decouples and modularizes functionality into loosely-joined `systems`.
Each system adds significant functionality, usually to achieve a set of important and tightly related goals.

For example, the Blockchain System provides structures like Block, Tipset, and Chain, and provides functionality
like Block Sync, Block Propagation, Block Validation, Chain Selection, and Chain Access. This is
separated from the Files, Pieces, Piece Preparation, and Data Transfer. Both of these systems are separated from
the Markets, which provide Orders, Deals, Market Visibility, and Deal Settlement.

## Why is System decoupling useful?

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
- **Scalability:** systems, and various use cases, may drive different performance requirements for different operators.
  System decoupling makes it easier for operators to scale their deployments along system boundaries.

## Filecoin Nodes don't need all the systems

Filecoin Nodes vary significantly and do not need all the systems.
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

## Separating Systems

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
required breaking apart the one actor into two.

## Decomposing within a System

Systems themselves decompose into smaller subunits. These are sometimes called "subsystems" to avoid confusion
with the much larger, first-class Systems. Subsystems themselves may break down further. The naming here is not
strictly enforced, as these subdivisions are more related to protocol and implementation engineering concerns
than to user capabilities.
