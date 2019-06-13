---
title: Filecoin Bootstrapping Routine
type: "docs"
---
# Filecoin Bootstrapping Routine

This spec describes the Filecoin bootstrapping protocol, it must be read along with:

- [Sync](sync.md) on how Filecoin nodes sync with bootstrap nodes to catch up to the chain and stay synced after the initial bootstrapping has occured. The `Syncer` interface and many of the methods used here are defined in the sync spec.

For related systems, see:

- [Network Protocols](network-protocols.md) on how Filecoin nodes can communicate with each other, with for instance, an [initial handshake](network-protocols.md#hello-handshake), or [block syncing](network-protocols.md#blocksync).
- [Operation](operation.md) on various operations a functional Filecoin node needs to run, like [DHT routing](operation.md#dht-for-peer-routing.md).

## What is bootstrapping in Filecoin?

When a Filecoin node first comes online, it must find the latest head for the chain and then initiate the syncing routine in order to catch up to the chain state and maintain it thereafter as new blocks come in. This is covered in the [sync spec](sync.md). Bootstrapping is the process through which a node finds the chain's latest head in the first place (thereby enabling it to sync), through other peers in the network.

At its simplest, bootstrapping can be thought of as secure network joining. It must enable a new node to get the latest head/chain that the majority of the network power is mining off and must be both:

- Robust — able to do so in adversarial networks (so long as honest nodes make up the majority of the network).
- Consistent — able to reliably get the canonical chain (even given poor latency or state loss in our initial peer set).

The specific bootstrapping heuristic could be left up to the node implementation, but we provide a simple one here as a candidate construction. This spec is divided as follows:

- Initial Connection — ensuring the node connects to an honest set of initial peers.

- Peer Set Expansion — ensuring that these peers can serve up the network's canonical head reliably in spite of their individual latencies or non-byzantine state loss.


Chain Selection is addressed in the [sync spec](sync.md).

We leave trustless initial peer selection as future work.

## Initial Connection

Initial peer selection is key to avoiding eclipse attacks. If the nodes you initially connect to are malicious, then they can simply serve you malicious other peers to all of your requests, and effectively prevent you from ever discovering the real network. Filecoin assumes 51% of rational or honest miners. Given no prior information, as is the case when bootstrapping, distinguishing a rational node from an adversarial one would require sampling a majority of the network.

This is obviously impractical, so we choose to start the bootstrap sequence with nodes by connecting with an inital set of trusted peers. These can be nodes run by well-known entities (the Filecoin foundation, large exchanges, good actors in the community, etc) or manually specified by the node operator in the config file. Default Filecoin nodes will ship with a set of default nodes though these can be changed at will.

We leave trustless (or fully decentralized) bootstrapping for future work.

A strawman bootstrapping sequence would have nodes use these initial peers to sync up to the canonical chain. However, new connection requests might effecitvely DoS these peers, making the Filecoin network unavailable to new nodes. Therefore, Filecoin bootstrapping seeks to "load-balance" initial requests to ensure the network remains widely available.

We connect to each of these initial peers and then initiate the `hello handshake` (see in [network protocols](network-protocols.md#hello-handshake)), getting back

```go
type HelloMessage struct {
HeaviestTipSet []Cid
HeaviestTipSetHeight uint64
GenesisHash Cid
}
```

We will use these new peers to sync with the chain, using the chain heads from our initial set of trusted nodes as a means of verifying that these peers are mining the canonical chain. These initial peers do not return their latest heads, but rather their latest confirmed heads, thereby ensuring other honest peers should be aware of these as well.

We assume no long-lived network partitions in Filecoin, meaning these confirmed heads will all be on the same chain (though potentially at different heights). A node should initially connect to 5 trusted nodes to begin bootstrapping.

```go
func (s *Syncer) InitialConnect(config ConfigFile) {
    b.trustedPeers = configFile.bootstrapPeers
    
    for trustedPeer := range s.trustedPeers {
        trustedChain := sayHello(trustedPeer)
        
        s.trustedHeads[trustedChain.HeaviestTipSet] = trustedChain.HeaviestTipSetHeight
    		
        if trustedChain.GenesisHash != GENESIS {
            Fatal("One of your trusted peers is bad.")
        }        
    }
    
    s.genesis = GENESIS
 		s.ExpandPeerSet()
}
```



## Peer Set Expansion

Once connected to some portion of the network, a node will need to expand their peer set in order to get a good sample of the network and a live latest head. To do so, use the DHT to do random `FindPeer` requests. We do so until we have requested 5 new peers from each peer in the PeerSet.

We look for random peers so our peer set makes up a "representative" set of the network to avoid both hotspots and potential attacks (syncing from given clusters).

In that sense, during the initial bootstrapping, we treat the Filecoin nodes network and the underlying libp2p network as two distinct networks for bootstrapping purposes. We specifically distinguish this bootstrapping phase (meant to help us sync to the chain and which doesn't use the DHT structure to pick nodes) from the node bootstrapping in which Filecoin nodes use the DHT to find nearby nodes with whom to exchange messages.

We will successively query random IDs (25 times) round-robinning our trusted nodes as starting points. We'll keep the nodes with IDs closest to our random one in our peer set.

```go
func (s *Syncer) ExpandPeerSet() {
  dht := getDHT()
  
  // get a random node from each trusted peer until we have 25
  trustedPeerNum := 0
  // we'll look for duplicates found sampling the network as a heuristic for interrupting the search if the network is too small
  dupCounter := 0
  for len(s.peerSet) < BootstrapPeerThreshold && dupCounter < DupThreshold {
    trustedPeer := s.trustedPeers[trustedPeerNum]
   	newPeer := syncer.getRandomPeer(trustedPeer)
    // adding a peer to the peer set will trigger syncing if our peer set has grown over the threshold
    if !syncer.addPeerToSet(newPeer) {
      dupCounter += 1
    }
    
    // round robin
    trustedPeerNum += 1
    trustedPeerNum %= len(s.trustedPeers)
  }
  
  // if we've reached the threshold, abandon expansion and sync with trusted heads instead
  if dupCounter >= DupThreshold {
    s.SyncBootstrap(false)
  }
}
```

For each of our BootstrapPeerThreshold peers, we perform hello handshakes, and run the following verifications:

- Ensure that each of our bootstrapper's `trustedHeads.keys` are included as ancestors of our peers' latestHeads.
- If they are not, include that peer in our syncer's `BadTipsetCache`.
- If they are, add these peers and latestHeads to our syncer's peerSet/peerHeads

Repeat this process requesting new peers until your PeerHeads contains at least 25 valid peers, to be used to select a latest head to mine off of.

If the random search yields too many duplicates, the peer set expansion is abandoned, under the assumption that the network is too small to support the peer expansion (this would happen early in the network's life or if the number of peers dipped below the threshold at some point).

In that case, the protocol will instead try to sync with the trustedPeers directly. In the steady-state, we assume the trustedPeers would reject such syncing requests under normal network conditions.