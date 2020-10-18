---
title: "GossipSub"
weight: 6
dashboardWeight: 2.0
dashboardState: stable
dashboardTests: 0
dashboardAudit: done
dashboardAuditURL: /#section-appendix.audit_reports.2020-06-03-gossipsub-design-and-implementation
dashboardAuditDate: '2020-06-03'
---

# GossipSub

Transaction messages and block headers alongside the message references are propagated using the [libp2p GossipSub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). **In order to guarantee interoperability between different implementations, all filecoin full nodes must implement and use this protocol.** All pubsub messages are authenticated and must be [syntactically validated](message#message-syntax-validation) before being propagated further.

GossipSub is a gossip-based pubsub protocol that is utilising two types of links to propagate messages: i) _mesh links_ that carry full messages in an _eager-push_ (i.e., proactive send) manner and ii) _gossip-links_ that carry message identifiers only and realise a _lazy-pull_ (i.e., reactive request) propagation model. Mesh links form a global mesh-connected structure, where, once messages are received they are forwarded in full to mesh-connected nodes, realizing an "eager-push" model. Instead, gossip-links are utilized periodically to complement the mesh structure. During gossip propagation, only message headers are sent to a selected group of nodes in order to inform them of messages that they might not have received before. In this case, nodes ask for the full message, hence, realizing a reactive request, or "lazy pull" model.

GossipSub includes a number of security extensions and mitigation strategies that make the protocol robust against attacks. Please refer to the [protocol's specification](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md) for details on GossipSub's design, implementation and parameter settings, or to the [technical report](https://arxiv.org/abs/2007.02754) for the design rationale and a more detailed evaluation of the protocol.