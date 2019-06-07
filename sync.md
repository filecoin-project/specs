# Chain Syncing

This spec describes the Filecoin sync protocol, for related systems, see:

- [Bootstrapping](./bootstrap.md) which describes how a node builds a peer set in the first place.
- [Network Protocols](./network-protocols.md) on how Filecoin nodes can communicate with each other, with for instance, an [initial handshake](./network-protocols.md#hello-handshake), or [block syncing](./network-protocols.md#blocksync).
- [Operation](./operation.md) on various operations a functional Filecoin node needs to run, like [DHT routing](./operation.md#dht-for-peer-routing.md).

## What is chain syncing in Filecoin?

Chain syncing is the process a filecoin node runs to sync its internal chain state with new blocks from the network and new blocks it itself has mined.  A node syncs in two distinct modes: `syncing` and `caught up`.  Chain syncing updates both local storage of chain data and the head of the current heaviest observed chain.

'Syncing' mode and 'Caught up' mode are two distinct processes. 'Syncing' mode, or 'the initial sync' is a process that is triggered when a node is far enough behind the rest of the network. This process terminates once the node's 'head' is sufficiently close to the current chain head. This is a configurable default, but as a default we may use 10*Block_Time.

Once 'syncing' is complete, the 'caught up' sync process begins. This process keeps the node up to date with the rest of the network, and terminates only when the node is shut down.

## Interface

```go
type Syncer struct {
	// The heaviest known TipSet in the network.
	head TipSet

	// The interface for accessing and putting TipSets into local storage
	store ChainStore

	// The known genesis TipSet
	genesis TipSet

	// the current mode the syncer is in
	syncMode SyncMode

	// TipSets known to be invalid
	bad BadTipSetCache

  // handle to the block sync service
  bsync BlockSync

  //peer set
  peerSet []PeerID

  // peer heads
  // Note: clear cache on disconnects
  peerHeads map[PeerID][]Cid
  
  // trustedPeers will be used to kick off peer set expansion
  // Filecoin nodes will ship with a default set
  trustedPeers []PeerID
  
  // trusted heads are the latest heads as reported by the trusted peers
  trustedHeads map[PeerId][]Cid  
}

// BoostrapPeerThreshold is the threshold needed to move from peerSetExpansion phase
// (see bootstrap.md) to chain syncing. This is a node parameter that can be changed.
const BootstrapPeerThreshold = 25
// Likewise dupThreshold is the number of duplicate nodes the node can be routed to 
// making "randomPeer" requests before it assumes the network is too small to reach the
// bootstrapPeerThreshold and falls back to trying to sync with the trustedPeers.
const DupThreshold = 10
```



## General Operation

Whenever a node hears about a new head, it will sync to it as warranted. `InformNewHead()` is called both during bootstrapping and as peers send across new blocks.

```go
// InformNewHead informs the syncer about a new potential TipSet
// This should be called when connecting to new peers, and additionally
// when receiving new blocks from the network
func (syncer *Syncer) InformNewHead(from PeerID, head TipSet) {
  
  if syncMode == Bootstrap:
    // InformNewHead will only get called during PeerSetExpansion during bootstrapping
		SyncBootstrap(expanding = true)
	else if syncMode == CaughtUp:
		syncer.SyncCaughtUp(blk)
	}
}

// SyncBootstrap is used to synchronise your chain when first joining
// the network, or when rejoining after significant downtime.
func (syncer *Syncer) SyncBootstrap(bool expanding) {
  
  0. if in expansion mode, ensure we have enough peers to start syncing up
  1. if not, use trustedPeers/Heads to sync up to
  2. round robin through the peers to expand peerset until large enought, then
  3. get heaviest head from the peers
  4. ensure it has correct genesis block
  5. if first connection: validate the blocks from genesis to head
  	 else validate blocks from own head to head
  6. validate state transitions across blocks
	7. switch own head to new head
  8. move to caught up mode

// SyncCaughtUp is used to stay in sync once caught up to
// the rest of the network.
func (syncer *Syncer) SyncCaughtUp(maybeHead TipSet) error {

  0. on every incoming new block,
  1. validate block
  2. compare block weight to own head weight
  3. if heavier, switch head
  	 else discard
}
```



## Syncing Mode

A filecoin node syncs in `syncing` mode when entering the network for the first time, or after being separated for a sufficiently long period of time.  The exact period of time comes from the consensus protocol (TODO specify more concretely, for example how does this relate to the consensus.Punctual method?).

During `syncing` mode a node learns about the newest head of the blockchain through the secure bootstrapping protocol. The syncing protocol then syncs the block headers for that entire chain, and validates their linking. It then fetches all the messages for the chain, and checks all the state transitions between the blocks and that the blocks were correctly created.  If validation passes the node's head is updated to the head TipSet received from bootstrapping.

Note that this series of node validations during `syncing` is a potential vector for DoSing nodes in the network. In summary, a fake chain of blocks with a large height could be passed to a node, causing it to expend a lot of resources fetching that chain before eventually discarding it. The malicious node passing the fake chain can spoof their genesis block thereby foiling that validation step. To that end, Filecoin nodes will ship with a default confirmation time parameter, enabling nodes to discard nodes within confirmation_time/block_time rounds if they have not yet converged with the known chain.

In the event a node thus receives a bad block, it should drop the connection to the peer who sent it, replacing it in its peer set with another random peer.

In this mode of operation a filecoin node should not mine or send messages as it will not be able to successfully generate heaviest blocks or reference the correct state of the chain to verify that messages will execute as expected.

(TODO: should include discussion of a `Load()` call to make use of existing chain data on a node during "re-awakening" case of `syncing` mode.)

## Caught Up Mode

A filecoin node syncs in `caught up` mode after completing `syncing` mode. A node stays in this mode until it is shut down. New block cids are gossiped from the network through the hello protocol or the network's [block pubsub protocol](data-propagation.md#block-propagation). A node also obtains new block cids coming from its own successfully mined blocks.  These cids are input to the `caught up` syncing protocol.  If these cids belong to a TipSet already in the store then they are already synced and the syncing protocol finishes.  If not the syncing protocol resolves the TipSet corresponding to the input cids.  It checks that this TipSet is not in its badTipSet cache, and that this TipSet is not too far back in the chain using the consensus `Punctual` method.  It then resolves the parent TipSet by reading off the parent cids in the header of any block of the TipSet.  The above procedure repeats until either an error is found or the store contains the next TipSet.  In the case of an error bad TipSets and their children not already in the bad TipSet cache are added to the cache before the call to `collectTipSetCaughtUp` returns.

After collecting a chain up to an ancestor TipSet that was previously synced to the store the syncing protocol checks each TipSet of the new chain for validity one by one.  When the filecoin network runs Expected Consensus, or any other multiple parents consensus protocol, the syncing protocol must consider not only the TipSets in the new chain but also possible new-heaviest TipSets that are the union of TipSets in the new chain and TipSets already in the store.  In the case of Expected Consensus there is at most one such TipSet: the TipSet made up of the union of the first new TipSet in the new chain being synced and the largest TipSet with the same parents kept in the store.

To sync new TipSets the `caught up` syncing protocol first runs a consensus validation check on the TipSet.  If any TipSet is invalid the syncing protocol finishes.  If a TipSet is valid the syncer adds the TipSet to the chain store.  The syncing protocol then checks whether the TipSet is heavier than the current head using the consensus weighting rules.  If it is heavier the chain updates the state of the node to account for the new heaviest TipSet.

## Maintaining a fresh Peer Set

Syncing depends on the validity of a node's peer set. In order to ensure that the peer set remains representative of the network's state after bootstrap, a node should replace peers as peers disconnect. (TODO in future: cycle through peers on regular basis).

```go
// triggered when a peer disconnects
func (syncer *Syncer) replacePeer(peer PeerID) {
 	1. when peer disconnects
  2. remove from peerSet
  3. fetch new random peer to add to peerSet instead
}

func (syncer *Syncer) getRandomPeer(from PeerID) newPeer PeerID {
  1. fetch a node from DHT querying a random node ID
	2. pick the closest peer that returns
  3. if not in peer set,
  	a. handshake and ensure it uses same genesis block
  	b. add to peer set
  4. else, repeat
}
```

# Dependencies

Things that affect the chain syncing protocol.

**Consensus protocol**
- The consensus protocol should define a punctual function: `func Punctual([]TipSet chain) bool`. `Punctual(chain) == true` when a provided chain does not fork from the node's view of the current best chain 'too far in the past', and false otherwise.
- The fork selection rule.  This includes the weighting function.  As part of this in the context of EC the syncer must consider TipSets that are the union of independently propagated TipSets.

**Chain storage**

- The current chain syncing protocol requires that the chain store never stores an invalid TipSet. If it does, then chain sync may sync to an invalid chain.

- The current chain syncing protocol requires that the chain store never stores an invalid TipSet.
