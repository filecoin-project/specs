---
title: libp2p
weight: 5
bookCollapseSection: true
dashboardWeight: 1
dashboardState: stable
dashboardTests: 0
dashboardAudit: done
dashboardAuditDate: '2019-10-10'
dashboardAuditURL: https://github.com/protocol/libp2p-vulnerabilities/blob/master/DRAFT_NCC_Group_ProtocolLabs_1903ProtocolLabsLibp2p_Report_2019-10-10_v1.1.pdf 
---

# Libp2p

[Libp2p](https://libp2p.io) is a modular network protocol stack for peer-to-peer networks. It consists of a catalogue of modules from which p2p network developers can select and reuse just the protocols they need, while making it easy to upgrade and interoperate between applications. This includes several protocols and algorithms to enable efficient peer-to-peer communication like peer discovery, peer routing and NAT Traversal. While libp2p is used by both IPFS and Filecoin, it is a standalone stack that can be used independently of these systems as well.

There are several implementations of libp2p, which can be found at the [libp2p GitHub repositoriy](https://github.com/libp2p). The specification of libp2p can be found in its [specs repo](https://github.com/libp2p/specs) and its documentation at [https://docs.libp2p.io](https://docs.libp2p.io).

Below we discuss how some of libp2p's components are used in Filecoin.

## DHT

The Kademlia DHT implementation of libp2p is used by Filecoin for peer discovery and peer exchange. Libp2p's [PeerID](https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md) is used as the ID scheme for Filecoin storage miners and more generally Filecoin nodes. One way that clients find miner information, such as a miner's address, is by using the DHT to resolve the associated PeerID to the miner's _Multiaddress_.

The Kademlia DHT implementation of libp2p in go can be found in its [GitHub repository](https://github.com/libp2p/go-libp2p-kad-dht).

## GossipSub

GossipSub is libp2p's pubsub protocol. Filecoin uses GossipSub for message and block propagation among Filecoin nodes. The recent hardening extensions of GossipSub include a number of techniques to make it robust against a variety of attacks.

Please refer to [GossipSub's Spec section](gossip_sub), or the protocol’s own and more complete [specification](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md) for details on the protocol’s design, implementation and parameter settings. A [technical report](https://arxiv.org/abs/2007.02754) is also available, which discusses the design rationale of the protocol.
