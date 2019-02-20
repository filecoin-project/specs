# Chain Syncing

## What is chain syncing in Filecoin?

Chain syncing is the process a filecoin node runs to sync its internal chain state with new blocks from the network and new blocks it itself has mined.  A node syncs in two distinct modes: `syncing` and `caught up`.  Chain syncing updates both local storage of chain data and the head of the current heaviest observed chain.

'Syncing' mode and 'Caught up' mode are two distinct processes. 'Syncing' mode, or 'the initial sync' is a process that is triggered when a node is far enough behind the rest of the network. This process terminates once the nodes 'head' is    sufficiently far ahead. Once 'syncing' is complete, the 'caught up' sync process begins. This process keeps the node up to date with the rest of the network, and terminates only when the node is shut down.

## Interface
```go
type Syncer struct {
	// The heaviest known tipset in the network.
	head TipSet

	// The interface for accessing and putting tipsets into local storage
	store ChainStore

	// The known genesis tipset
	genesis TipSet
    
    // the current mode the syncer is in
    syncMode SyncMode

	// TipSets known to be invalid
	bad BadTipSetCache
    
    // handle to the block sync service
    bsync BlockSync
    
    // peer heads
    // Note: clear cache on disconnects
    peerHeads map[PeerID]Cid
}

const BootstrapPeerThreshold = 5

// InformNewHead informs the syncer about a new potential tipset
// This should be called when connecting to new peers, and additionally
// when receiving new blocks from the network
func (syncer *Syncer) InformNewHead(from PeerID, head TipSet) {
    switch syncer.syncMode {
    case Bootstrap:
        go SyncBootstrap(from, head)
    case CaughtUp:
        go syncer.SyncCaughtUp(blk)
    }
}

// SyncBootstrap is used to synchronise your chain when first joining
// the network, or when rejoining after significant downtime.
func (syncer *Syncer) SyncBootstrap() {
    syncer.syncLock.Lock()
    defer syncer.syncLock.Unlock()
    syncer.peerHeads[from] = head
    if len(syncer.peerHeads) < BootstrapPeerThreshold {
        // not enough peers to sync yet...
        return
    }
    
    selectedHead := selectHead(syncer.peerHeads)
    
    cur := selectedHead
    var blockSet BlockSet
    for head.Height() > 0 {
        // NB: GetBlocks validates that the blocks are in-fact the ones we
        // requested, and that they are correctly linked to eachother. It does
        // not validate any state transitions
        blks := syncer.bsync.GetBlocks(head, RequestWidth)
        blockSet.Insert(blks)
            
        head = blks.Last().Parents()
    }
    
    genesis := blockSet.GetByHeight(0)
    if genesis != syncer.genesis {
        // TODO: handle this...
        Error("We synced to the wrong chain!")
        return
    }
    
    // Fetch all the messages for all the blocks in this chain
    // There are many ways to make this more efficient. For now, do the dumb thing
    blockSet.ForEach(func(b Block) {
        FetchMessages(b)
    })
    
    // Now, to validate some state transitions
    base := genesis
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
	ts := tipsetFromCidOverNet(newHead) // lookup over network

	var chain Chain
	for {
		if !consensus.Punctual(ts) {
			syncer.bad.InvalidateChain(chain)
			syncer.bad.InvalidateTipSet(ts)
			return nil, errors.New("tipset forks too far back from head")
		}

		chain.InsertFront(ts)

		if syncer.store.Contains(ts) { // Store has record of this tipset.
			return chain, nil
		}
		parent := ts.ParentCid()
		ts, err = tipsetFromCidOverNet(parent)
		if err != nil {
			return nil, err
		}
	}
}
```

## Syncing Mode
A filecoin node syncs in `syncing` mode when entering the network for the first time, or after being separated for a sufficiently long period of time.  The exact period of time comes from the consensus protocol (TODO specify more concretely, for example how does this relate to the consensus.Punctual method?).

During `syncing` mode a node learns about the newest head of the blockchain through the secure bootstrapping protocol. The syncing protocol then syncs the block headers for that entire chain, and validates their linking. It then fetches all the messages for the chain, and checks all the state transitions between the blocks and that the blocks were correctly created.  If validation passes the node's head is updated to the head tipset received from bootstrapping.

In this mode of operation a filecoin node should not mine or send messages as it will not be able to successfully generate heaviest blocks or reference the correct state of the chain to verify that messages will execute as expected.

(TODO: should include discussion of a `Load()` call to make use of existing chain data on a node during "re-awakening" case of `syncing` mode.)

## Caught Up Mode
A filecoin node syncs in `caught up` mode after completing `syncing` mode. A node stays in this mode until they are shutdown. New block cids are gossiped from the network through the hello protocol or the network's [block pubsub protocol](data-propogation.md#block-propogation). A node also obtains new block cids coming from its own successfully mined blocks.  These cids are input to the `caught up` syncing protocol.  If these cids belong to a tipset already in the store then they are already synced and the syncing protocol finishes.  If not the syncing protocol resolves the tipset corresponding to the input cids.  It checks that this tipset is not in its badTipSet cache, and that this tipset is not too far back in the chain using the consensus `Punctual` method.  It then resolves the parent tipset by reading off the parent cids in the header of any block of the tipset.  The above procedure repeats until either an error is found or the store contains the next tipset.  In the case of an error bad tipsets and their children not already in the bad tipset cache are added to the cache before the call to `collectTipSetCaughtUp` returns.

After collecting a chain up to an ancestor tipset that was previously synced to the store the syncing protocol checks each tipset of the new chain for validity one by one.  When the filecoin network runs Expected Consensus, or any other multiple parents consensus protocol, the syncing protocol must consider not only the tipsets in the new chain but also possible new-heaviest tipsets that are the union of tipsets in the new chain and tipsets already in the store.  In the case of Expected Consensus there is at most one such tipset: the tipset made up of the union of the first new tipset in the new chain being synced and the largest tipset with the same parents kept in the store.

To sync new tipsets the `caught up` syncing protocol first runs a consensus validation check on the tipset.  If any tipset is invalid the syncing protocol finishes.  If a tipset is valid the syncer adds the tipset to the chain store.  The syncing protocol then checks whether the tipset is heavier than the current head using the consensus weighting rules.  If it is heavier the chain updates the state of the node to account for the new heaviest tipset.

# Dependencies
Things that affect the chain syncing protocol.

**Consensus protocol**
- The consensus protocol should define a punctual function: `func Punctual([]TipSet chain) bool`. `Punctual(chain) == true` when a provided chain does not fork from the node's view of the current best chain 'too far in the past', and false otherwise.
- The fork selection rule.  This includes the weighting function.  As part of this in the context of EC the syncer must consider tipsets that are the union of independently propagated tipsets.

**Chain storage**

- The current chain syncing protocol requires that the chain store never stores an invalid tipset.

# Open Questions
- Secure bootstrapping in `syncing` mode
- How do we handle the lag between the initial head bootstrapped in `syncing` mode and the network head once the first `SyncBootstrap` call is complete?  Likely we'll need multiple `SyncBootstrap` calls.  Should they be parallelized?
- The properties of the chain store implementation have significant impact on the design of the syncing protocol and the syncing protocol's resistance to Denial Of Service (DOS) attacks.  For example if the chain store naively keeps all blocks in storage nodes are more vulnerable to running out of space.  As another example the syncer assumes that the store always contains a punctual ancestor of the heaviest chain. Should the spec grow to include properties of chain storage so that the syncing protocol can guarantee a level of DOS resistance?  Should chain storage be completely up to the implementation?  Should the chain storage spec be a part of the syncing protocol?  