---
title: Message Syncer
weight: 1
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Message Syncer
---

Transaction messages are sorted in the `mpool` of miners as they arrive according to cryptoeconomic rules followed by the miner.

Block messages are ordered as they arrive at miner nodes in order to extend the blockchain.


## Message Propagation

Messages are propagated over libp2p pubsub topics, using the [GossipSub](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) protocol. Every serialised `SignedMessage` is announced (see [Message](message)) in the corresponding topic by any node participating in the network.

There are two pubsub topics related to the mempool: i) the `/fil/msgs/` topic that carries transaction messages and, ii) the `/fil/blocks/` topic that carries block messages. The process is as follows:
1. When a client wants to carry out a transaction in the Filecoin network, they publish a transaction message in the `/fil/msgs/` topic.
2. The message propagates to all other nodes in the network using GossipSub and eventually ends up in the `mpool` of all miners.
3. Depending on cryptoeconomic rules, some miner will eventually pick the transaction message from the `mpool` (together with other transaction messages) and include it in a block.
4. The miner publishes the newly-mined block in the `/fil/blocks/` pubsub topic and the block message propagates to all nodes in the network (including the nodes that published the transactions included in this block).

Nodes must check that incoming messages are valid, that is, that they have a valid signature. If the message is not valid it should be dropped and must not be forwarded.

The updated, hardened version of the GossipSub protocol includes a number of attack mitigation strategies. For instance, when a node receives an invalid message it assigns a negative _score_ to the sender peer. Peer scores are not shared with other nodes, but are rather kept locally by every peer for all other peers it is interacting with. If a peer's score drops below a threshold it is excluded from the scoring peer's mesh. We discuss more details on these settings in the GossipSub section. The full details can be found in the [GossipSub Specification](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub).

{{<hint warning>}}
**TODO:** discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.
{{</hint>}}