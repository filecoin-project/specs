---
title: "GossipSub"
weight: 6
dashboardWeight: 1.5
dashboardState: wip
dashboardAudit: complete
dashboardTests: 0
dashboardAuditURL: https://gateway.ipfs.io/ipfs/QmWR376YyuyLewZDzaTHXGZr7quL5LB13HRFnNdSJ3CyXu/Least%20Authority%20-%20Gossipsub%20v1.1%20Final%20Audit%20Report%20%28v2%29.pdf
---

# GossipSub
---

Messages and block headers along side the message references are propagated using the [gossipsub libp2p pubsub router](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub). Every full node must implement and run that protocol. All pubsub messages are authenticated and must be [syntactically validated](message#message-syntax-validation) before being propagated further.

Further more, every full node must implement and offer the bitswap protocol and provide all Cid Referenced objects, it knows of, through it. This allows any node to fetch missing pieces (e.g. `Message`) from any node it is connected to. However, the node should fan out these requests to multiple nodes and not bombard any single node with too many requests at a time. A node may implement throttling and DDoS protection to prevent such a bombardment.

## Bitswap

Run bitswap to fetch and serve data (such as blockdata and messages) to and from other filecoin nodes. This is used to fill in missing bits during block propagation, and also to fetch data during sync.

There is not yet an official spec for bitswap, but [the protobufs](https://github.com/ipfs/go-bitswap/blob/master/message/pb/message.proto) should help in the interim.
