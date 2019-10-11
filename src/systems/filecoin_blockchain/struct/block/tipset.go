package block

import (
	clock "github.com/filecoin-project/specs/systems/filecoin_nodes/clock"
)

func (ts *Tipset_I) MinTicket() base.Ticket {
	var ret Ticket

	for _, currBlock := range ts.Blocks() {
		tix := currBlock.Ticket()
		if ret == nil {
			ret = tix
		} else {
			ret = SmallerBytes(tix, ret)
		}
	}

	Assert(ret != nil)
	return tix
}

func (ts *Tipset_I) ValidateSyntax() bool {

	if len(ts.Blocks_) <= 0 {
		return false
	}

	panic("TODO")

	// parents := ts.Parents_
	// grandparent := parents[0].Parents_
	// for i := 1; i < len(parents); i++ {
	// 	if grandparent != parents[i].Parents_ {
	// 		return false
	// 	}
	// }

	// numTickets := len(ts.Blocks_[0].Tickets_)
	// for i := 1; i < len(ts.Blocks_); i++ {
	// 	if numTickets != len(ts.Blocks_[i].Tickets_) {
	// 		return false
	// 	}
	// }

	return true
}

func (ts *Tipset_I) LatestTimestamp() clock.Time {
	var latest clock.Time
	panic("TODO")
	// for _, blk := range ts.Blocks_ {
	// 	if blk.Timestamp().After(latest) || latest.IsZero() {
	// 		latest = blk.Timestamp()
	// 	}
	// }
	return latest
}

// func (tipset *Tipset_I) StateTree() stateTree.StateTree {
// 	var currTree stateTree.StateTree = nil
// 	for _, block := range tipset.Blocks() {
// 		if currTree == nil {
// 			currTree = block.StateTree()
// 		} else {
// 			Assert(block.StateTree().CID().Equals(currTree.CID()))
// 		}
// 	}
// 	Assert(currTree != nil)
// 	for _, block := range tipset.Blocks() {
// 		currTree = UpdateStateTree(currTree, block)
// 	}
// 	return currTree
// }
<<<<<<< HEAD:src/systems/filecoin_blockchain/struct/block/tipset.go
=======

func (bl *Block_I) TipsetAtEpoch(epoch Epoch) Tipset_I {
	dist := bl.Epoch() - epoch - 1
	current := bl.ParentTipset_
	parents := current.Parents_
	for i := range dist {
		current = parents
		parent = current.Parents_
	}

	return current
}

func (bl *Block_I) TicketAtEpoch(epoch base.Epoch) base.Ticket_I {
	ts := bl.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}

func (chain *Chain_I) TipsetAtEpoch(epoch base.Epoch) base.Tipset_I {
	dist := chain.HeadEpoch() - epoch
	current := chain.HeadTipset()
	parents := current.Parents_
	for i := range dist {
		current = parents
		parent = current.Parents_
	}

	return current
}

func (chain *Chain_I) TicketAtEpoch(epoch base.Epoch) base.Ticket_I {
	ts := chain.TipsetAtEpoch(epoch)
	return ts.MinTicket()
}

func (chain *Chain_I) FinalizedEpoch() Epoch {
	ep := HeadEpoch()
	return ep - spc.GetFinality()
}
>>>>>>> ticket sampling coded.:src/systems/filecoin_blockchain/blockchain/block.go
