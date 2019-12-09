---
title: Network Interface
statusIcon: ⚠️
---

{{< readfile file="network.id" code="true" lang="go" >}}


Filecoin nodes use the libp2p protocol for peer discovery, peer routing, and message multicast, and so on. Libp2p is a set of modular protocols common to the peer-to-peer networking stack. Nodes open connections with one another and mount different protocols or streams over the same connection. In the initial handshake, nodes exchange the protocols that each of them supports and all Filecoin related protcols will be mounted under `/filecoin/...` protocol identifiers.

Here is the list of libp2p protocols used by Filecoin.

- Graphsync: 
	- Graphsync is used to transfer blockchain and user data
	- [Draft spec](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md)
	- No filecoin specific modifications to the protocol id
- Gossipsub: 
	- block headers and messages are broadcasted through a Gossip PubSub protocol where nodes can subscribe to topics for blockchain data and receive messages in those topics. When receiving messages related to a topic, nodes processes the message and forwards it to its peers who also subscribed to the same topic.  
	- Spec is [here](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub)
	- No filecoin specific modifications to the protocol id.  However the topic identifiers MUST be of the form `fil/blocks/<network-name>` and `fil/msgs/<network-name>`
- KademliaDHT: 
	- Kademlia DHT is a distributed hash table with a logarithmic bound on the maximum number of lookups for a particular node. Kad DHT is used primarily for peer routing as well as peer discovery in the Filecoin protocol.
	- Spec TODO [reference implementation](https://github.com/libp2p/go-libp2p-kad-dht)
	- The protocol id must be of the form `fil/kad/<network-name>`
- Bootstrap List: 
	- Bootstrap is a list of nodes that a new node attempts to connect upon joining the network. The list of bootstrap nodes and their addresses are defined by the users.
- Peer Exchange: 
	- Peer Exchange is a discovery protocol enabling peers to create and issue queries for desired peers against their existing peers
	- spec [TODO](https://github.com/libp2p/specs/issues/222)
	- No Filecoin specific modifications to the protocol id.
- DNSDiscovery: Design and spec needed before implementing
- HTTPDiscovery: Design and spec needed before implementing
- Hello:
	- Hello protocol handles new connections to filecoin nodes.  It is an important part of the discovery process for ambient protocols (like KademliaDHT)
	- Spec TODO.
	- No Filecoin specific modifications to the protocol id.
