# Filecoin Network

The Filecoin network is built using libp2p building blocks, as transports and protocols, as well as some additional Filecoin specific protocols as outlined in [Network Protocols](./network-protocols.md).

## Required Protocols

Every full node must support the following libp2p protocols:

 - [gossipsub](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) ([for data announcements](./data-propagation.md))
 - [bitswap] ([for data exchange](./data-propagation.md))
 - Filecoin specific Protocols:
   - Hello Handshake
   - StorageDeal
   - BlockSync

## Transports, Streams & Encryption

Connections between nodes are end-to-end-encrypted and -authenticated, thus every node has a `peerId` associated with it.

Filecoin protocols are run as multiple streams over a single connection. Only one connection should exist between two nodes at any time, multiple streams within that connections must be used to exchange different protocols. This is usually handled from within the libp2p stack if configured properly.

This document doesn't make any assumptions on which transports and multiplexing features a particular node offers, however it must at least be able to respond and communicate through mplex channels via a secio transport.

{{% notice info %}}
Implementers are encouraged to offer and experiment with more protocols and transports, as long as the above mentioned requirements of ee2e, authentication and muxing are still met. In particular the authors encourage to look into TLS1.3 and QUIC as libp2p protocol/transport layers.
{{% /notice %}}

### Establishing connections

When a Filecoin node connects to the network, it may connect to multiple nodes at once. When a new node connects, it must be greeted with the `HelloMessage` through the hello-handshake protocol first. Only after that message has been sent can other streams be opened. If a node receives a `HelloMessage` with a GenesisBlock it doesn't support, it must immediately close the entire connection to that peer.

If the node learns through that handshake about newer blocks it should use that new information to sync up their chain.

### Syncing

Whenever a node learns about a new `BlockHead` it should attempt to import the block - let that be through the `HelloMessage` or through the [block pubsub protocol](data-propagation.md#block-propagation). For that it may fetch the ancestors of that block through bitswap until it is fully caught up. Importing in this context refers to the node confirming that the block is valid, as described in [Validation](./validation.md) and storing it locally. This mode is known as "syncing". During "syncing" the same node may not mine blocks.



## Establishing a network

Every node should aim to establish a set of stable connections to the network in order to stay up to date with latest changes and block propagation - aka to sync - as well as to participate in the processes of the network - e.g. mining.

In order to do that a node must be "bootstrapped" when starting up. At its simplest, bootstrapping can be thought of as secure network joining. It must enable a new node to get the latest head/chain that the majority of the network power is mining off and must be both:

- Robust — able to do so in adversarial networks (so long as honest nodes make up the majority of the network).
- Consistent — able to reliably get the canonical chain (even given poor latency or state loss in our initial peer set).

The specific bootstrapping method and heuristic - or a mix of them - is left up to the node implementation, but we provide some example ideas of how this could be done and what aspects should be taken into consideration for that, such as relying on both new nodes and new connections to ensure a node is connected to both a reliable and representative slice of the network

### Trusted bootnodes

As with many other networks, trusted entities - e.g. the Filecoin foundation, large exchanges, good actors in the community, etc - might publish a set of bootnodes one can connect to in order to establish network connections.

Keep in mind that, because these could easily be DoS'ed, a) large set of different entities should be pre-shipped but only a few should be tried at a time b) they might drop connections rather quickly after establishing (and ban that peer for a short period of time) in order to serve more nodes. Thus you want to make sure to send a discovery request to learn about more peers over the bootnodes to have nodes to connect to.

A node should not attempt to connect to more than 10 bootnodes at a time. And should diversify its set of nodes as quickly as possible.

### Previously known nodes

It is generally preferred that an implementation keeps permanent track of reliable peers and attempts to connect to them upon restart. This can be a simple list of previously established connections or the blow listed ranked reputation list.

### Peer Discovery requests

At this point the specification does not decide upon any particular peer discovery mechanism but leaves this up for the implementations to decide which ones they want to provide. However, the expectation is that all boot nodes provide some form peer discovery mechanism and publish that information along side the connectivity information. Whether that be BRAHMS or through RandomQueries of KADEMLIA-DHT, any node bootstrapping itself, should attempt to establish streams for peer discovery and may not fail just because their peer declines the given protocol stream.

On the other side, the attempt to establish an unsupported peer-discovery protocol must not - in and of itself - lead to a disconnect, as multiple protocols might be tried. Any protocol may only be tried once on the same connection however.

## Maintaining a stable network

Many aspects of the network topology assume some degree of liveness for filecoin to work. Thus and for it to work most efficiently, we aim for an overall stable network, that still allows for new nodes to join the network without interrupting.

In order to achieve that every node should aim for "future usefulness of peers" as outline below. However, as a general rule, if a node acts in violation of a `must` rule, any peer aware of that is totally in their right to drop any connection to the node without further warning. It may also keep a record of these offenses ban the peer from connecting again, if it finds them to continue offending. Such ban may be imposed upon - in increasing order - the peerId, on the IP+Port or finally the IP itself. One must keep in mind that IP addresses still rotate and thus every implementation must add a maximum timeout should they impose a ban on them, while PeerIDs can also be banned indefinitely.

### Optimise for usefulness of peers

A more sophisticated system to manage a healthy peer set is through optimising for usefulness of that peer to the node. in this system a node tracks all incoming messages and their costs in relation to the value it provided to them and records them as value any particular node provides to it - called "reputation". These calculations may also take into account time that passed or whether the other node would have to know about the uselessness of a particular message. Sending a useless message, in this system, is not understood as a hard violation, but as impolite behaviour and can be recorded as such. Thus creating a ranking by usefulness among the peer set.

An example would be that, continuously receive a pubsub message from a node after we've forwarded it to them prior. This isn't a hard failure, as though they aren't supposed to do that, we can't know if that is because of a faulty implementation or because of network delays queuing their message. However, we'd still note this as impolite and for every time this happens deduct from their reputation as it isn't useful to us.

A node attempts to always hold a certain amount of connections ("slots") to the network -- we recommend 25-50 on an on-the-shelf system. Whenever a new node connects, it can check their previously stored reputation or assign a default value and if that reputation is higher than the lowest currently connected nodes, may replace that connection (and thus drop the lowest quality connection) or otherwise refuse to take that connection.

In this system a node might also only record a strong decrease in reputation but not drop a connection to a peer directly even upon a strong violation, because the node may still be more useful than others. It is better to stay connected to a crappy network than to no network. However, this doesn't not free the node from still adhering to the spec itself - it should not forward said violation or its connection might still be righteously dropped. 

We also recommend to regularly check the ranking and drop and clear up the lowest 10% of slots, leaving 5% open for incoming connections and fill up the other 5% by connecting to other nodes it can find through peer discovery.

All this to create a local view of for that node most useful connections to the network. The reputation may be stored permanently and be available between restarts - thus providing a neat bootstrap start list, too. 

