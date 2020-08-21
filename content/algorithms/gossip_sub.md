---
title: "GossipSub"
weight: 6
dashboardWeight: 2.0
dashboardState: stable
dashboardAudit: complete
dashboardTests: 0
auditURL: https://gateway.ipfs.io/ipfs/QmWR376YyuyLewZDzaTHXGZr7quL5LB13HRFnNdSJ3CyXu/Least%20Authority%20-%20Gossipsub%20v1.1%20Final%20Audit%20Report%20%28v2%29.pdf
---

# GossipSub
---

Transaction messages and block headers alongside the message references are propagated using the [libp2p GossipSub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). **Both miners and full nodes must implement and run this protocol.** All pubsub messages are authenticated and must be [syntactically validated](message#message-syntax-validation) before being propagated further. The GossipSub Specification ++ provides all the implementation details and parameter settings. Here we provide an overview of the main characteristics and design choices of GossipSub.

GossipSub is a gossip-based pubsub protocol that is utilising two types of links to propagate messages: i) _mesh links_ that carry full messages in an _eager-push_ manner and ii) _gossip-links_ that carry message identifiers only and realise a _lazy-pull_ propagation model.

GossipSub includes a number of security extensions and mitigation strategies that make the protocol robust against attacks. Please refer to the [protocol's specification](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md) for details on GossipSub's design and parameter settings.