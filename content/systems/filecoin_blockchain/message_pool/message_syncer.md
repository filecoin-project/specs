---
title: Message Syncer
weight: 1
dashboardWeight: 2
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Message Syncer
---

{{< hint warning >}}
TODO:

- explain message syncer works
- include the message syncer code
{{< /hint >}}


## Message Propagation

Messages are propagated over the libp2p pubsub topic `/fil/msgs`, using the [GossipSub](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) protocol. Every serialised `SignedMessage` is announced (see [Message](message)) in the `/fil/msgs` topic by any node participating in the network.

Upon receiving the message, nodes must check its validity: the signature must be valid, and the account in question must have enough funds to cover the actions specified. If the message is not valid it should be dropped and must not be forwarded.

The updated, hardened version of the GossipSub protocol includes a number of attack mitigation strategies. For instance, when a node receives an invalid message it assigns a negative _score_ to the sender peer. Peer scores are not shared with other nodes, but are rather kept locally by every peer for all other peers it is interacting with. If a peer's score drops below a threshold it is excluded from the scoring peer's mesh. We discuss more details on these settings in the GossipSub section. The full details can be found in the [GossipSub Specification](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub).

{{<hint warning>}}
**TODO:** discuss checking signatures and account balances, some tricky bits that need consideration. Does the fund check cause improper dropping? E.g. I have a message sending funds then use the newly constructed account to send funds, as long as the previous wasn't executed the second will be considered "invalid" ... though it won't be at the time of execution.
{{</hint>}}