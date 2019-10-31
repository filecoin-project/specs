package block

import (
	"bytes"

	clock "github.com/filecoin-project/specs/systems/filecoin_nodes/clock"
)

func (ts *Tipset_I) MinTicket() Ticket {
	var ret Ticket

	for _, currBlock := range ts.Blocks() {
		tix := currBlock.Ticket()
		if ret == nil {
			ret = tix
		} else {
			smaller := SmallerBytes(tix.Output(), ret.Output())
			if bytes.Equal(smaller, tix.Output()) {
				ret = tix
			}
		}
	}

	return ret
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
