---
title: Network Interface
weight: 3
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Interface
---

{{<embed src="network.id" lang="go" >}}


Filecoin nodes use the libp2p protocol for peer discovery, peer routing, and message multicast, and so on. Libp2p is a set of modular protocols common to the peer-to-peer networking stack. Nodes open connections with one another and mount different protocols or streams over the same connection. In the initial handshake, nodes exchange the protocols that each of them supports and all Filecoin related protcols will be mounted under `/fil/...` protocol identifiers.

Here is the list of libp2p protocols used by Filecoin.

- Graphsync: 
	- Graphsync is used to transfer blockchain and user data
	- [Draft spec](https://github.com/ipld/specs/blob/master/block-layer/graphsync/graphsync.md)
	- No filecoin specific modifications to the protocol id
- Gossipsub: 
	- block headers and messages are broadcasted through a Gossip PubSub protocol where nodes can subscribe to topics for blockchain data and receive messages in those topics. When receiving messages related to a topic, nodes process the message and forward it to peers who also subscribed to the same topic.
	- Spec is [here](https://github.com/libp2p/specs/tree/master/pubsub/gossipsub)
	- No filecoin specific modifications to the protocol id.  However the topic identifiers MUST be of the form `fil/blocks/<network-name>` and `fil/msgs/<network-name>`
- KademliaDHT: 
	- Kademlia DHT is a distributed hash table with a logarithmic bound on the maximum number of lookups for a particular node. Kad DHT is used primarily for peer routing as well as peer discovery in the Filecoin protocol.
	- Spec TODO [reference implementation](https://github.com/libp2p/go-libp2p-kad-dht)
	- The protocol id must be of the form `fil/<network-name>/kad/1.0.0`
- Bootstrap List: 
	- Bootstrap is a list of nodes that a new node attempts to connect to upon joining the network. The list of bootstrap nodes and their addresses are defined by the users.
- Peer Exchange: 
	- Peer Exchange is a discovery protocol enabling peers to create and issue queries for desired peers against their existing peers
	- spec [TODO](https://github.com/libp2p/specs/issues/222)
	- No Filecoin specific modifications to the protocol id.
- DNSDiscovery: Design and spec needed before implementing
- HTTPDiscovery: Design and spec needed before implementing
- Hello:
	- Hello protocol handles new connections to filecoin nodes to facilitate discovery
	- the protocol string is `fil/hello/1.0.0`. 

## Hello Spec 

### Protocol Flow

`fil/hello` is a filecoin specific protocol built on the libp2p stack.  It consists of two conceptual
procedures: `hello_connect` and `hello_listen`.   

`hello_listen`: `on new stream` -> `read peer hello msg from stream` -> `write latency message to stream` -> `close stream`

`hello_connect`: `on connected` -> `open stream` -> `write own hello msg to stream` -> `read peer latency msg from stream`  -> `close stream`

where stream and connection operations are all standard libp2p operations.  Nodes running the Hello Protocol should consume the incoming Hello Message and use it to help manage peers and sync the chain.

### Messages
{{<embed src="hello.id" lang="go" >}}


When writing the `HelloMessage` to the stream the peer must inspect its current head to provide accurate information.  When writing the `LatencyMessage` to the stream the peer should set `TArrival` immediately upon receipt and `TSent` immediately before writing the message to the stream.
