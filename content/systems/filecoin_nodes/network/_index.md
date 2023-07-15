---
title: Network Interface
weight: 3
dashboardWeight: 1
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Network Interface

Filecoin nodes use several protocols of the libp2p networking stack for peer discovery, peer routing and block and message propagation. Libp2p is a modular networking stack for peer-to-peer networks. It includes several protocols and mechanisms to enable efficient, secure and resilient peer-to-peer communication. Libp2p nodes open connections with one another and mount different protocols or streams over the same connection. In the initial handshake, nodes exchange the protocols that each of them supports and all Filecoin related protocols will be mounted under `/fil/...` protocol identifiers.

The complete specification of libp2p can be found at [https://github.com/libp2p/specs](https://github.com/libp2p/specs).
Here is the list of libp2p protocols used by Filecoin.

- **Graphsync:** Graphsync is a protocol to synchronize graphs across peers. It is used to reference, address, request and transfer blockchain and user data between Filecoin nodes. The [draft specification of GraphSync](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md) provides more details on the concepts, the interfaces and the network messages used by GraphSync. There are no Filecoin-specific modifications to the protocol id.

- **Gossipsub:** Block headers and messages are propagating through the Filecoin network using a gossip-based pubsub protocol acronymed _GossipSub_. As is traditionally the case with pubsub protocols, nodes subscribe to topics and receive messages published on those topics. When nodes receive messages from a topic they are subscribed to, they run a validation process and i) pass the message to the application, ii) forward the message further to nodes they know of being subscribed to the same topic. Furthermore, v1.1 version of GossipSub, which is the one used in Filecoin is enhanced with security mechanisms that make the protocol resilient against security attacks. The [GossipSub Specification](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub) provides all the protocol details pertaining to its design and implementation, as well as specific settings for the protocols parameters. There have been no Filecoin-specific modifications to the protocol id. However the topic identifiers MUST be of the form `fil/blocks/<network-name>` and `fil/msgs/<network-name>`

- **Kademlia DHT:** The Kademlia DHT is a distributed hash table with a logarithmic bound on the maximum number of lookups for a particular node. In the Filecoin network, the Kademlia DHT is used primarily for peer discovery and peer routing. In particular, when a node wants to store data in the Filecoin network, they get a list of miners and their node information. This node information includes (among other things) the PeerID of the miner. In order to connect to the miner and exchange data, the node that wants to store data in the network has to find the Multiaddress of the miner, which they do by querying the DHT. The [libp2p Kad DHT Specification](https://github.com/libp2p/go-libp2p-kad-dht) provides implementation details of the DHT structure. For the Filecoin network, the protocol id must be of the form `fil/<network-name>/kad/1.0.0`.

- **Bootstrap List:** This is a list of nodes that a new node attempts to connect to upon joining the network. The list of bootstrap nodes and their addresses are defined by the users (i.e., applications).

- **Peer Exchange:** This protocol is the realisation of the peer discovery process discussed above at Kademlia DHT. It enables peers to find information and addresses of other peers in the network by interfacing with the DHT and create and issue queries for the peers they want to connect to.
