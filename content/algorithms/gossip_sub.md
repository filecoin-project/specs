---
title: "GossipSub"
weight: 6
dashboardWeight: 2.0
dashboardState: stable
dashboardAudit: 1
dashboardTests: 0
---

# GossipSub
---

Transaction messages and block headers alongside the message references are propagated using the [libp2p GossipSub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). **Both miners and full nodes must implement and run this protocol.** All pubsub messages are authenticated and must be [syntactically validated](message#message-syntax-validation) before being propagated further. The GossipSub Specification ++ provides all the implementation details and parameter settings. Here we provide an overview of the main characteristics and design choices of GossipSub.

GossipSub is a gossip-based pubsub protocol that is utilising two types of links to propagate messages: i) _mesh links_ that carry full messages in an _eager-push_ manner and ii) _gossip-links_ that carry message identifiers only and realise a _lazy-pull_ propagation model.

1. _Mesh links_ build the GossipSub mesh. Each node connects to the mesh through a number of connections, `D`, which indicates the degree of the network. The degree, `D`, is accompanied by two thresholds, `D_low` and `D_high` that act as boundaries. When the number of connections increases above `D_high`, the node _prunes_ some of the connections, while when it decreases below `D_low` the node _grafts_ new connections. Both of these happen in order to keep the degree in the area `<D_low, D_high>`.
2. _Gossip links_ augment the message propagation performance of the protocol. Gossiping allows the network to operate on a low degree and therefore, keep traffic under certain levels.


## Control Messages

The protocol defines five control messages and a heartbeat:
- GRAFT a mesh link: this notifies the peer that it has been added to the local mesh of the grafting node.
- PRUNE a mesh link: this notifies the peer that it has been removed from the local mesh of the pruning peer.
- PRUNE-Peer Exchange (PX): the pruning peer sends a list of peer IDs to the pruned peer to help it connect to alternative peers and expand its mesh. There is a backoff timer associated with pruning to avoid re-grafting to recent peers for a period of time (set to 1min).
- IHAVE: gossip; this notifies the peer that the following messages were recently seen and are available on request.
- IWANT: gossip; request transmission of messages announced
in an IHAVE message.

The protocol also implements two simple messages to join and leave a topic, JOIN[topic] and LEAVE[topic], which are implemented through a GRAFT[topic] and LEAVE[topic] message, respectively.

## Mesh Maintenance

The protocol is running a mesh maintenance round every 1 sec. Two separate processes take place during mesh maintenance: i) every peer is adjusting its mesh (if needed), that is, peers are grafting or pruning other peers to bring their mesh to the desirable degree window `<D_low, D_high>`, and ii) peers emit gossip. Gossip is realised in three rounds, one in every “mesh maintenance round”. The rationale behind setting the number of rounds to three was to reach a certain level of network coverage. In our case, we wanted to ensure that gossip messages reach ~50% of nodes in the network.

## Peer Discovery

As part of its peer discovery tools, GossipSub supports: i) _Peer Exchange_, which
allows applications to bootstrap from a known set of bootstrap peers without an external peer discovery mechanism, and ii) _Explicit Peering Agreements_, where the application
can specify a list of peers to which nodes should connect when joining.

- **Peer Exchange.** This process is supported through either bootstrap nodes or other normal peers. Bootstrap nodes are maintained by system operators. They have to be stable and operate independently of the mesh construction, that is, bootstrap nodes do not maintain connections to the mesh. Bootstrap nodes maintain scores for peers they interact with (scoring is discussed later on in this section) and refuse to serve misbehaving peers and/or advertise misbehaving peers’ addresses to others. They also participate in propagating gossip. Their role is to facilitate formation and maintenance of the network. In terms of peer exchange between normal peers, whenever a node is excluded from the mesh, e.g., due to oversubscription, the pruning peer provides it with a list of alternative nodes, which it can use to reconnect, re-build or extend its mesh.
- **Explicit Peering Agreements.** With explicit peering, the application can specify a list of peers which nodes should connect to when joining. For every explicit peer, the router must establish and maintain a bidirectional (reciprocal) connection. Explicit peering connections exist outside the mesh, which in practice means that every new valid incoming message is forwarded to all of a peer’s explicit peers.


## Security Extensions

GossipSub incorporates an extensive set of security extensions that make it resilient against malicious behaviour. The security extensions are realised with: i) a peer scoring function and ii) a set of mitigation strategies.

### Peer Scoring

Every peer in a GossipSub-based network monitors the performance and behaviour of peers it knows of, i.e., both those that it is directly connected to in the mesh and those that it is interacting with through gossip. The score, captured in the `D_score` parameter, is not shared with other peers, that is, it is not a reputation system, but instead it is used by the node locally to identify whether a particular peer is behaving as expected or not. Based on that, nodes make grafting and pruning decisions driven by some of the mitigation strategies discussed next.
 
The score function takes six parameters into account and calculates the score as a weighted mix of parameters, 4 of them per topic and 3 of them globally applicable.

```
Score(p) = TopicCap(Σtᵢ*(w₁(tᵢ)*P₁(tᵢ) + w₂(tᵢ)*P₂(tᵢ) + w₃(tᵢ)*P₃(tᵢ) + w₃b(tᵢ)*P₃b(tᵢ) + w₄(tᵢ)*P₄(tᵢ))) + w₅*P₅ + w₆*P₆ + w₇*P₇
```
where `tᵢ` is the topic weight for each topic where per topic parameters apply.

These values were chosen carefully to identify the most important behavioural characteristics that differentiate a malicious from an honest node. These are:

- `P₁`: **Time in Mesh** for a topic. This is the time a peer has been in the mesh, capped to a small value and mixed with a small positive weight. This is intended to boost peers already in the mesh so that they are not prematurely pruned because of oversubscription.
- `P₂`: **First Message Deliveries** for a topic. This is the number of messages first delivered by the peer in the topic, mixed with a positive weight. This is intended to reward peers who first forward a valid message.
- `P₃`: **Mesh Message Delivery Rate** for a topic. This parameter is a threshold for the expected message delivery rate within the mesh in the topic. If the number of deliveries is above the threshold, then the value is 0. If the number is below the threshold, then the value of the parameter is the square of the deficit. This is intended to penalize peers in the mesh who are not delivering the expected number of messages so that they can be removed from the mesh. The parameter is mixed with a negative weight.
- `P₃b`: **Mesh Message Delivery Failures** for a topic. This is a sticky parameter that counts the number of mesh message delivery failures. Whenever a peer is pruned with a negative score, the parameter is augmented by the rate deficit at the time of prune. This is intended to keep history of prunes so that a peer that was pruned because of underdelivery cannot quickly get re-grafted into the mesh. The parameter is mixed with negative weight.
- `P₄`: **Invalid Messages** for a topic. This is the number of invalid messages delivered in the topic. This is intended to penalize peers who transmit invalid messages, according to application specific validation rules. It is mixed with a negative weight.
- `P₅`: **Application Specific** score. This is the score component assigned to the peer by the application itself, using application specific rules. The weight is positive, but the parameter itself has an arbitrary real value, so that the application can signal misbehaviour with a negative score or gate peers before an application specific handshake is completed.
- `P₆`: **IP Colocation Factor**. This parameter is a threshold for the number of peers using the same IP address. If the number of peers in the same IP exceeds the threshold, then the value is the square of the surplus, otherwise it is 0. This is intended to make it difficult to carry out sybil attacks by using a small number of IPs. The parameter is mixed with a negative weight.
- `P₇`: **Behavioural Penalty**. This parameter captures penalties applied for misbehaviour. The parameter has an associated (decaying) counter, which is explicitly incremented by the router on specific events. The value of the parameter is the square of the counter and is mixed with a negative weight.

**Parameter Decay**

The topic parameters are implemented using counters maintained internally by the router whenever an event of interest occurs. The counters _decay_ periodically so that their values are not continuously increasing and ensure that a large positive or negative score isn't sticky for the lifetime of the peer.

The decay interval is configurable by the application, with shorter intervals resulting in faster decay.

Each decaying parameter can have its own decay _factor_, which is a configurable parameter that controls how much the parameter will decay during each decay period.

For more details on the parameter calculation, please refer to the [GossipSub Specification](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md)


### GossipSub Mitigation Strategies

Mitigation strategies are the measures that GossipSub uses to defend against attacks. Some of the mitigation strategies are making use of the score function discussed above as a driver to make decisions on whether or not to prune a peer from the mesh (see "Controlled Mesh Maintenance"), while others are focusing on bypassing Sybil-dominated connections (see "Flood Publishing" and "Adaptive Gossip Dissemination").

**Controlled Mesh Maintenance.** During the mesh maintenance, every peer is pruning all peers it is connected to (in the mesh) that have a negative score. The peer keeps the highest scoring peers out of `D_score` and selects the remaining ones at random. These randomly chosen peers guarantee that the mesh is responsive to new nodes joining the network. The selection is done under the constraint that `D_out` peers are outbound connections; if the scoring plus random selection does not result in enough outbound connections, then we replace the random and lower scoring peers in the selection with outboud connection peers.
**Opportunistic Grafting.** In cases where a peer is stuck with a mesh of poorly-performing peers (e.g., due to a successful attack or node churn). Depending on the severity of the attack and the Sybil to honest ratio, the protocol might be slow to respond (e.g., pruning negative-scoring peers and choosing again from a poisoned pool). As a mitigation to this situation, the Opportunistic Grafting mechanism is periodically checking the median score of peers in its mesh against a threshold (the `opportunisticGraftThreshold`). If the median score of its immediately connected peers are below this threshold, the peer opportunistically grafts (at least) two peers with score above the median of its existing mesh connections.
**Flood Publishing.** Flood publishing has been identified as an efficient way to bypass sybil-dominated mesh connections of a peer that is under attack. With flood-publishing even if `D_high - 1` connections are Sybil-controlled, the newly published message will still make it to the rest of the network through the one remaining connection.
**Adaptive Gossip Dissemination.** Instead of emitting gossip to a fixed number of peers, the protocol should adjust according to the network size, which may vary wildly in case of Sybil attacks (i.e., large number of Sybil nodes). In order to achieve that, GossipSub is emitting gossip to `gossipFactor` number of peers, currently set to 25% of the peers known to the gossiping peer.
**Backoff on PRUNE.** This measure is an extension to the "controlled mesh maintenance" one. According to the Backoff on PRUNE, when a peer is excluded from the mesh due to poor scoring, it is not allowed to re-connect (to the peer where it was pruned from) for some period, which is called "backoff" and is currently set to 1 min.