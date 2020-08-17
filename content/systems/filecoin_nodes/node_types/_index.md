---
title: Node Types
weight: 1
bookCollapseSection: true
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---


Fundamentally, there are only two major node types in the Filecoin network and this is how the Lotus implementation is realising _node types_ in the strict sense of the word. Nodes are identified with a repository (directory) in the host in a one-to-one relationship - that is one repo belongs to a single node. That said one host might realise multiple Filecoin nodes by having the corresponding repositories. The two major node types are:

The **Full Node:** this is the chain validation or chain verifier node. A Full Node must synchronise the chain (ChainSync) when it first joins the network to reach current consensus. From then on, the node must constantly be fetching any addition to the chain (i.e., receive the latest blocks) and validate them to reach consensus state.
The **Miner Node:** Miner nodes are the backbone of the Filecoin blockchain network. They have all the functionality of the Full Node, but also extend this functionality to create blocks, submit blocks to the blockchain and ultimately extend the blockchain.

In order to realise the full potential of the Filecoin network, there is more functionality that needs to be added on top of the _base_ node types described above. This extra functionality can be thought of as a service or subsystem that run on top of the two node types above. We refer to a node that realises extra functionality as _node_ too, for clarity. However, note that different node types are not mutually exclusive to each other, meaning that one physical node might have to implement more than one type of node functionality.

A Filecoin implementation should support the following subsystems, or types of nodes:

- **Chain Verifier Node:** this is the **Full Node** described above. This type of node cannot play an active role in the network, unless it implements **Client Node** functionality, described below.
- **Client Node:** this type of node builds on top of the **Full or Chain Verifier Node** and must be implemented by any application that is building on the Filecoin network. This can be thought of as the main infrastructure node (at least as far as interaction with the blockchain is concerned) of applications such as exchanges or decentralised storage applications building on Filecoin. The node should implement and interact (as a client) with the Storage and Retrieval Markets, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.
- **Retrieval Miner Node:** this node type is extending the **Full or Chain Verifier Node** to add _retrieval miner_ functionality, that is, participate in the retrieval market. As such, this node type needs to implement the retrieval market provider subsystem, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.
- **Storage Miner Node:** this type of node is building on top of the Miner Node described above and must implement all of its functionality for validating, creating and adding blocks to the blockchain. It should implement the storage mining subsystem, the storage market provider subsystem, keep the Market Order Book and be able to do Data Transfers through the Data Transfer Module.

