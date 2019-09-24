---
title: Network Interface
---

```
type NetworkInterface struct {
  libp2p.Node

  MountProtocol() // ...
}
```

Filecoin nodes use the libp2p protocol for peer discovery, peer routing, and message multicast, and so on. Libp2p is a set of modular protocols common to the peer-to-peer networking stack. Nodes open connections with one another and mount different protocols or streams over the same connection. In the initial handshake, nodes exchange the protocols that each of them supports and all Filecoin related protcols will be mounted under `/filecoin/...` protocol identifiers.

Here is the list of libp2p protocols used by Filecoin.

- Graphsync: TODO
- Bitswap:  TODO
- Gossipsub: block headers and messages are broadcasted through a Gossip PubSub protocol where nodes can subscribe to topics such as `NewBlock`, `BlockHeader`, `BlockMessage`, etc and receive messages in those topics. When receiving messages related to a topic, nodes processes the message and forwards it to its peers who also subscribed to the same topic.
- Kad-DHT: Kademlia DHT is a distributed hash table with a logarithmic bound on the maximum number of lookups for a particular node. Kad DHT is used primarily for peer routing as well as peer discovery in the Filecoin protocol. 
- Bootstrap: Bootstrap is a list of nodes that a new node attempts to connect upon joining the network. The list of bootstrap nodes and their addresses are defined by the users.
