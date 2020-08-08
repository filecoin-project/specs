---
title: Message Pool
weight: 2
bookCollapseSection: true
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Message Pool
---
The Message Pool, or `mpool` or `mempool` is a Memory Pool of Messages in the Filecoin blockchain system. It acts as the interface between Filecoin nodes (of any type) and the peer-to-peer network of other nodes used for off-chain message propagation. The message pool is used by nodes to maintain a set of messages to transmit to the Filecoin VM (for "on-chain" execution).

In order for a block message to end up in the blockchain it first has to be in the message pool. In reality, at least in the Lotus implementation of Filecoin, there is no central pool of messages stored somewhere. Instead, the message pool is an abstraction and is realised as a list of messages kept by every node in the network. Therefore, when a node puts a new message in the message pool, this message is propagated to the rest of the network using libp2p's pubsub protocol, GossipSub. Nodes need to subscribe to the corresponding pubsub topic in order to receive messages.

Message propagation using GossipSub does not happen immediately and therefore, there is some lag before message pools at different nodes can be in sync. In practice, and given continuous streams of messages being added to the message pool and the delay to propagate messages, the message pool is never synchronised across all nodes in the network. This is not a deficiency of the system, as the message pool does not _need_ to be synchronized across the network. 


{{< embed src="message_pool_subsystem.id" lang="go" >}}

The node types that use the message pool subsystem, as well as the specific reason they are using it for are as follows:

- storage market provider and client nodes: for transmission of deals on chain,
- storage miner nodes: for transmission of PoSts, sector commitments, deals, and other operations tracked on chain,
- verifier nodes: for transmission of potential faults on chain,
- relayer nodes: for forwarding and discarding messages appropriately.

The message pool should have a maximum size defined to avoid DoS attacks, where nodes are spammed and run out of memory. The recommended size for the message pool is 5000 Transaction messages.

The message pool subsystem is made of two components:

- The [Message Syncer](message_syncer.md) -- which receives and propagates messages.
- [Message Storage](message_storage.md) -- which caches messages according to a given policy.



