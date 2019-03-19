# Filecoin Bootstrapping Routine

This spec describes how to implement the Filecoin bootstrapping protocol, for related systems, see:

- [Sync](sync.md) on how Filecoin nodes sync with bootstrap nodes to catch up to the chain and stay synced after the initial bootstrapping has occured.
- [Network Protocols](./network-protocols.md) on how Filecoin nodes can communicate with each other, with for instance, an [initial handshake](./network-protocols.md#hello-handshake), or [block syncing](./network-protocols.md#blocksync).
- [Operation](./operation.md) on various operations a functional Filecoin node needs to run, like [DHT routing](./operation.md#dht-for-peer-routing.md).

## What is bootstrapping in Filecoin?

When a Filecoin node first comes online, it must find the latest head for the chain and then initiate the syncing routine in order to catch up to the chain state and maintain it thereafter as new blocks come in. This is covered in the [sync spec](./sync.md). Bootstrapping is the process through which a node finds the chain's latest head in the first place (thereby enabling it to sync), through other peers in the network.

At its simplest, bootstrapping can be thought of as secure network joining. It must enable a new node to get the latest head/chain that the majority of the network power is mining off and must be both:

- Robust — able to do so in adversarial networks (so long as honest nodes make up the majority of the network).

- Consistent — able to reliably get the canonical chain (even given poor latency or state loss in our initial peer set).

The specific bootstrapping heuristic could be left up to the node implementation, but we provide a simple one here as a candidate construction. This spec is divided as follows:

- Initial Connection — ensuring the node connects to an honest set of initial peers.
- Peer Set Expansion — ensuring that these peers can serve up the network's canonical head reliably in spite of their individual latencies or state loss.
- Chain Selection — choosing the most likely canonical chain thereafter.

We leave trustless initial peer selection as future work.

## Initial Connection

Initial peer selection is key to avoiding eclipse attacks. If the nodes you initially connect to are malicious, then they can simply serve you malicious other peers to all of your requests, and effectively prevent you from ever discovering the real network. Filecoin assumes 51% of rational or honest miners. Given no prior information, as is the case when bootstrapping, distinguishing a rational node from an adversarial one would require sampling a majority of the network.

This is obviously impractical, so we choose to start the bootstrap sequence with nodes by connecting with an inital set of trusted peers. These can be nodes run by well-known entities (the Filecoin foundation, large exchanges, good actors in the community, etc) or manually specified by the node operator in the config file. We leave trustless (or fully decentralized) bootstrapping for future work.

A strawman bootstrapping sequence would have nodes use these initial peers to sync up to the canonical chain. However, new connection requests might effecitvely DoS these peers, making the Filecoin network unavailable to new nodes. Therefore, Filecoin bootstrapping seeks to "load-balance" initial requests to ensure the network remains widely available.

We connect to each of these initial peers using the `hello handshake` (see in [network protocols](./network-protocols.md#hello-handshake)), getting back

```go
type HelloMessage struct {
HeaviestTipSet []Cid
HeaviestTipSetHeight uint64
GenesisHash Cid
}
```

followed by a `GetPeers` request, getting back

```go
type PeersMessage struct {
PeerSet []PeerID
}
```

We will use these new peers to sync with the chain, using the chain heads from our initial set of trusted nodes as a means of verifying that these peers are mining the canonical chain. These initial peers do not return their latest heads, but rather their latest confirmed heads, thereby ensuring other honest peers should be aware of these as well.

We assume no long-lived network partitions in Filecoin, meaning these confirmed heads will all be on the same chain (though potentially at different heights). A node should initially connect to 5 trusted nodes to begin bootstrapping.

```go
type Bootstrapper struct {
    trustedPeers []PeerID
    trustedHeads map[[]Cid]uint64
    genesisHash Cid
    
    bootstrapPeers []PeerID
    
}

func InitialConnect(config ConfigFile, b Bootstrapper, s Syncer) {
    b.trustedPeers = configFile.bootstrapPeers
    
    for trustedPeer := range b.trustedPeers {
        trustedChain := sayHello(trustedPeer)
        otherPeers := getPeers(trustedPeer).PeerSet
        
        b.trustedHeads[trustedChain.HeaviestTipSet] = trustedChain.HeaviestTipSetHeight
    	
        if b.genesisHash == nil {
	        b.genesisHash = trustedChain.GenesisHash           
        } else if b.genesisHash != trustedChain.GenesisHash {
            Fatal("One of your trusted peers is bad.")
        }
        
        b.bootstrapPeers = append(b.bootstrapPeers, otherPeers)
    }
    
    s.genesis = b.genesisHash
}
```



## Peer Set Expansion

Once connected to some portion of the network, a node will need to expand their peer set in order to get a good sample of the network and a live latest head. To do so, use the DHT to do random `FindPeer` requests. We do so until we have requested 5 new peers from each peer in the PeerSet.

For each of our 30 peers, we perform hello handshakes, and run the following verifications:

- Ensure that each of our bootstrapper's `trustedHeads.keys` are included as ancestors of our peers' latestHeads.
- If they are not, include that peer in our syncer's `BadTipsetCache`.
- If they are, add these peers and latestHeads to our syncer's peerSet/peerHeads (with offset initialConnections)

Repeat this process requesting new peers until your PeerHeads contains at least 25 valid peers, to be used to select a latest head to mine off of.

## Chain Selection

We now have 25 peers mining off of the same trusted ancestor heads. Disagreements in their chains have to do with the fact that these may not yet be confirmed blocks. We simply want to start mining off of the heaviest of their chains to sync up to and mine off of. 

```go
func selectHead(heads map[PeerID]TipSet) TipSet {
    headsArr := toArray(heads)
    sel := headsArr[0]
    for i := 1; i < len(headsArr); i++ {
        cur := headsArr[i]
        
        if cur.IsAncestorOf(sel) {
            continue
        }
        if sel.IsAncestorOf(cur) {
            sel = cur
            continue
        }
        
        nca := NearestCommonAncestor(cur, sel)
        if sel.Height() - nca.Height() > ForkLengthThreshold {
        	// TODO: handle this better than refusing to sync
        	Fatal("Conflict exists in heads set")
        }

        if cur.Weight() > sel.Weight() {
            sel = cur
        }
    }
    return sel
}
```

The selected chain should be passed to the syncer, and the ‘initial sync’ should begin.
