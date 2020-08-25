---
title: Node Types
weight: 1
bookCollapseSection: true
dashboardWeight: 1
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

Nodes in the Filecoin network are primarily identified in terms of the services they provide. The type of node, therefore, depends on which services a node provides. A basic set of services in the Filecoin network include:
- chain verification
- storage market client
- storage market provider
- retrieval market client
- retrieval market provider
- storage mining

Any node participating in the Filecoin network should provide the _chain verification_ service as a minimum. Depending on which extra services a node provides on top of chain verification, it gets the corresponding functionality and Node Type "label".

Nodes can be realized with a repository (directory) in the host in a one-to-one relationship - that is, one repo belongs to a single node. That said, one host can implement multiple Filecoin nodes by having the corresponding repositories.

A Filecoin implementation should support the following subsystems, or types of nodes:

- **Chain Verifier Node:** this is the minimum functionality that a node needs to have in order to participate in the Filecoin network. This type of node cannot play an active role in the network, unless it implements **Client Node** functionality, described below. A Chain Verifier Node must synchronise the chain (ChainSync) when it first joins the network to reach current consensus. From then on, the node must constantly be fetching any addition to the chain (i.e., receive the latest blocks) and validate them to reach consensus state.
- **Client Node:** this type of node builds on top of the **Chain Verifier Node** and must be implemented by any application that is building on the Filecoin network. This can be thought of as the main infrastructure node (at least as far as interaction with the blockchain is concerned) of applications such as exchanges or decentralised storage applications building on Filecoin. The node should implement the _storage market and retrieval market client_ services. The client node should interact with the Storage and Retrieval Markets, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.
- **Retrieval Miner Node:** this node type is extending the **Chain Verifier Node** to add _retrieval miner_ functionality, that is, participate in the retrieval market. As such, this node type needs to implement the _retrieval market provider_ service, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.
- **Storage Miner Node:** this type of node must implement all of the required functionality for validating, creating and adding blocks to extend the blockchain. It should implement the chain verification, storage mining and storage market provider services, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.

# Node Interface

{{<embed src="content/externals/lotus/node/repo/interface.go" lang="go" >}}

# Chain Verifier Node

```go
type ChainVerifierNode interface {
  FilecoinNode

  systems.Blockchain
}
```
{{<embed src="content/externals/lotus/node/impl/full.go" lang="go" >}}

# Client Node

```go
type ClientNode struct {
  FilecoinNode

  systems.Blockchain
  markets.StorageMarketClient
  markets.RetrievalMarketClient
  markets.MarketOrderBook
  markets.DataTransfers
}
```
{{<embed src="content/externals/lotus/node/impl/client/client.go" lang="go" >}}

# Storage Miner Node

```go
type StorageMinerNode interface {
  FilecoinNode

  systems.Blockchain
  systems.Mining
  markets.StorageMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```
{{<embed src="content/externals/lotus/node/impl/storminer.go" lang="go" >}}

# Retrieval Miner Node

```go
type RetrievalMinerNode interface {
  FilecoinNode

  blockchain.Blockchain
  markets.RetrievalMarketProvider
  markets.MarketOrderBook
  markets.DataTransfers
}
```

# Relayer Node

```go
type RelayerNode interface {
  FilecoinNode

  blockchain.MessagePool
  markets.MarketOrderBook
}
```

