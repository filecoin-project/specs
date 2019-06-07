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
    // not enough peers to sync yet...
    if expanding && len(syncer.peerHeads) < BootstrapPeerThreshold {
        return
    }
  
  	// if not expanding, start with trusted peers
  	syncPeers := syncer.peerSet
  	syncHeads := syncer.peerHeads
  	if (!expanding) {
    		syncPeers = syncer.trustedPeers
      	syncHeads = syncer.trustedHeads
  	}
  	// Will now get heaviest head from all the heads from our Peerset
    selectedHead := selectHead(syncer.peerHeads)

    cur := selectedHead
    var blockSet BlockSet
    for head.Height() > 0 {
        // NB: GetBlocks validates that the blocks are in-fact the ones we
        // requested, and that they are correctly linked to each other. It does
        // not validate any state transitions
        blks := syncer.bsync.GetBlocks(head, RequestWidth)
        blockSet.Insert(blks)

        head = blks.Last().Parents()
    }

    // Fetch all the messages for all the blocks in this chain
    // There are many ways to make this more efficient. For now, do the dumb thing
    blockSet.ForEach(func(b Block) {
        // FetchMessages should use bitswap to fetch any messages we don't have locally
        FetchMessages(b)
    })

  	// Ensure that the selectedHead has the right genesis block
  	// Should be checked for trusted nodes in InitialConnect() and for
  	// all others in addPeerToSet(), but we leave details up to impl.
  	selectedGenesis := blockSet.GetByHeight(0)
	  assert(selectedGenesis == genesis, "State failure: trying to sync to wrong chain")

    // Now, to validate some state transitions
    base := syncer.genesis
    for i := 1; i < selectedHead.Height(); i++ {
        next := blockSet.GetByHeight(i)
        if !ValidateTransition(base, next) {
            // TODO: do something productive here...
            Error("invalid state transition")
            return
        }
    }

    blockSet.PersistTo(syncer.store)
    syncer.head = bset.Head()
    syncer.syncMode = CaughtUp
}

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

// SyncCaughtUp is used to stay in sync once caught up to
// the rest of the network.
func (syncer *Syncer) SyncCaughtUp(maybeHead TipSet) error {
	chain, err := syncer.collectChainCaughtUp(maybeHead)
	if err != nil {
		return err
	}

	// possibleTs enumerates possible tipsets that are the union
	// of tipsets from the chain and the store
	for _, ts := range possibleTs(chain[1:]) {
		if err := consensus.Validate(ts, store); err != nil {
			return err
		}
		syncer.store.PutTipSet(ts)
		if consenus.Weight(ts) > consensus.Weight(head) {
			syncer.head = ts
		}
	}
	return nil
}

func (syncer *Syncer) collectChainCaughtUp(maybeHead TipSet) (Chain, error) {
	// fetch TipSet and messages via bitswap
	ts := tipsetFromCidOverNet(newHead)

	var chain Chain
	for {
		if !consensus.Punctual(ts) {
			syncer.bad.InvalidateChain(chain)
			syncer.bad.InvalidateTipSet(ts)
			return nil, errors.New("TipSet forks too far back from head")
		}

		chain.InsertFront(ts)

		if syncer.store.Contains(ts) {
			// Store has record of this TipSet.
			return chain, nil
		}
		parent := ts.ParentCid()
		ts = tipsetFromCidOverNet(parent)
		}
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
    delete(syncer.PeerSet, peer)
    delete(syncer.PeerHeads, peer)

  	newPeer := syncer.getRandomPeer()
  	// addPeerToSet will validate this peer (i.e. check it has right genesis, etc.)
  	while !syncer.addPeerToSet(newPeer) {
      newPeer = syncer.getRandomPeer(self.ownID)
  	}
}

func (syncer *Syncer) getRandomPeer(from PeerID) newPeer PeerID {
  	wantedPeer := generateRandomNodeID()
    // libp2p's GetClosestPeers gets back k nearest peers, ordered
    peers := dht.GetClosestPeers(from, wantedPeer)
	  // We only use the nearest to get a good random peer
    newPeer := peers[0]
		return newPeer
}

func (syncer *Syncer) addPeerToSet(newPeer PeerID) bool {
  	// verify new peer to check whether to include in peerSet
    if !s.peerSet.contains(newPeer) {
      peerChain := sayHello(newPeer)

      if peerChain.GenesisHash == GENESIS || trustedPeer.isAncestorOf(newPeer) {
        s.peerSet = append(s.peerSet, newPeer)
        s.InformNewHead(newPeer, trustedChain.HeaviestTipSet)
        return true
      }
    }
  return false
}
```

# Dependencies

Things that affect the chain syncing protocol.

**Consensus protocol**
- The consensus protocol should define a punctual function: `func Punctual([]TipSet chain) bool`. `Punctual(chain) == true` when a provided chain does not fork from the node's view of the current best chain 'too far in the past', and false otherwise.
- The fork selection rule.  This includes the weighting function.  As part of this in the context of EC the syncer must consider TipSets that are the union of independently propagated TipSets.

**Chain storage**

<<<<<<< HEAD
- The current chain syncing protocol requires that the chain store never stores an invalid TipSet. If it does, then chain sync may sync to an invalid chain.
=======
- The current chain syncing protocol requires that the chain store never stores an invalid TipSet.

# Open Questions

- Secure bootstrapping in `syncing` mode
- How do we handle the lag between the initial head bootstrapped in `syncing` mode and the network head once the first `SyncBootstrap` call is complete?  Likely we'll need multiple `SyncBootstrap` calls.  Should they be parallelized?
- The properties of the chain store implementation have significant impact on the design of the syncing protocol and the syncing protocol's resistance to Denial Of Service (DOS) attacks.  For example if the chain store naively keeps all blocks in storage nodes are more vulnerable to running out of space.  As another example the syncer assumes that the store always contains a punctual ancestor of the heaviest chain. Should the spec grow to include properties of chain storage so that the syncing protocol can guarantee a level of DOS resistance?  Should chain storage be completely up to the implementation?  Should the chain storage spec be a part of the syncing protocol?
>>>>>>> 168d0226bae61e0942a65b84f42cd99711c0fe29
