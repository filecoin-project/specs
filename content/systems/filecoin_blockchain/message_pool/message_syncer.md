---
title: Message Propagation
weight: 1
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Message Propagation

The message pool has to interface with the libp2p pubsub [GossipSub](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) protocol. This is because transaction messages are propagated over [GossipSub](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) the corresponding `/fil/msgs/` _topic_. Every [Message](message) is announced in the corresponding `/fil/msgs/` topic by any node participating in the network.

There are two main pubsub topics related to transactions and blocks: i) the `/fil/msgs/` topic that carries transactions and, ii) the `/fil/blocks/` topic that carries blocks. The `/fil/msgs/` topic is linked to the `mpool`. The process is as follows:

1. When a client wants to carry out a transaction in the Filecoin network, they publish a transaction message in the `/fil/msgs/` topic.
2. The message propagates to all other nodes in the network using GossipSub and eventually ends up in the `mpool` of all miners.
3. Depending on cryptoeconomic rules, some miner will eventually pick the transaction message from the `mpool` (together with other transaction messages) and include it in a block.
4. The miner publishes the newly-mined block in the `/fil/blocks/` pubsub topic and the block message propagates to all nodes in the network (including the nodes that published the transactions included in this block).

Nodes must check that incoming transaction messages are valid, that is, that they have a valid signature. If the message is not valid it should be dropped and must not be forwarded.

The updated, hardened version of the GossipSub protocol includes a number of attack mitigation strategies. For instance, when a node receives an invalid message it assigns a negative _score_ to the sender peer. Peer scores are not shared with other nodes, but are rather kept locally by every peer for all other peers it is interacting with. If a peer's score drops below a threshold it is excluded from the scoring peer's mesh. We discuss more details on these settings in the GossipSub section. The full details can be found in the [GossipSub Specification](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub).

NOTES:

- _Fund Checking:_ It is important to note that the `mpool` logic is not checking whether there are enough funds in the account of the transaction message issuer. This is checked by the miner before including a transaction message in a block.
- _Message Sorting:_ Transaction messages are sorted in the `mpool` of miners as they arrive according to cryptoeconomic rules followed by the miner and in order for the miner to compose the next block.
