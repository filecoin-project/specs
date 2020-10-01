---
title: Chain Manager
weight: 3
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# Chain Manager

The _Chain Manager_ is a central component in the blockchain system. It tracks and updates competing subchains received by a given node in order to select the appropriate blockchain head: the latest block of the heaviest subchain it is aware of in the system.

In so doing, the _chain manager_ is the central subsystem that handles bookkeeping for numerous other systems in a Filecoin node and exposes convenience methods for use by those systems, enabling systems to sample randomness from the chain for instance, or to see which block has been finalized most recently.


## Chain Extension

### Incoming block reception

For every incoming block, even if the incoming block is not added to the current heaviest tipset, the chain manager should add it to the appropriate subchain it is tracking, or keep track of it independently until either:
- it is able to add to the current heaviest subchain, through the reception of another block in that subchain, or
- it is able to discard it, as the block was mined before finality.

It is important to note that ahead of finality, a given subchain may be abandoned for another, heavier one mined in a given round. In order to rapidly adapt to this, the chain manager must maintain and update all subchains being considered up to finality.

Chain selection is a crucial component of how the Filecoin blockchain works. In brief, every chain has an associated weight accounting for the number of blocks mined on it and so the power (storage) they track. The full details of how selection works are provided in the [Chain Selection](expected_consensus#chain-selection) section.

**Notes/Recommendations:**
1. In order to make certain validation checks simpler, blocks should be indexed by height and by parent set. That way sets of blocks with a given height and common parents may be quickly queried.
2. It may also be useful to compute and cache the resultant aggregate state of blocks in these sets, this saves extra state computation when checking which state root to start a block at when it has multiple parents.
3. It is recommended that blocks are kept in the local datastore regardless of whether they are understood as the best tip at this point - this is to avoid having to refetch the same blocks in the future.

### ChainTipsManager

The Chain Tips Manager is a subcomponent of Filecoin consensus that is responsible for tracking all live tips of the Filecoin blockchain, and tracking what the current 'best' tipset is.

```go
// Returns the ticket that is at round 'r' in the chain behind 'head'
func TicketFromRound(head Tipset, r Round) {}

// Returns the tipset that contains round r (Note: multiple rounds' worth of tickets may exist within a single block due to losing tickets being added to the eventually successfully generated block)
func TipsetFromRound(head Tipset, r Round) {}

// GetBestTipset returns the best known tipset. If the 'best' tipset hasn't changed, then this
// will return the previous best tipset.
func GetBestTipset()

// Adds the losing ticket to the chaintips manager so that blocks can be mined on top of it
func AddLosingTicket(parent Tipset, t Ticket)
```
