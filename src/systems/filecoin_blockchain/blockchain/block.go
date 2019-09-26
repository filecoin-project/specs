package blockchain

import (
	base "github.com/filecoin-project/specs/systems/filecoin_blockchain"
	clock "github.com/filecoin-project/specs/systems/filecoin_nodes/clock"
	util "github.com/filecoin-project/specs/util"
)

func SmallerBytes(a, b util.Bytes) util.Bytes {
	if util.CompareBytesStrict(a, b) > 0 {
		return b
	}
	return a
}

func ExtractElectionSeed(lookbackTipset *Tipset_I) base.ElectionSeed {

	panic("TODO")
	// var ret []byte

	// for _, currBlock := range lookbackTipset.Blocks() {
	// 	for _, currTicket := range currBlock.Tickets() {

	// 		currSeed := Hash(
	// 			HashRole_ElectionSeedFromVRFOutput,
	// 			currTicket.VRFResult().bytes(),
	// 		)
	// 		if ret == nil {
	// 			ret = currSeed
	// 		} else {
	//             ret = SmallerBytes(currSeed, ret)
	//         }
	// 	}
	// }

	// Assert(ret != nil)
	// return ElectionSeed.FromBytesInternal(nil, ret)
}

// func GenerateElectionTicket(k VRFKeyPair, seed ElectionSeed) Ticket {
// 	var vrfResult VRFResult = VRFEval(k, seed.ToBytesInternal())

// 	var vdfInput []byte = Hash(
// 		HashRole_TicketVDFInputFromVRFOutput,
// 		vrfResult.ToBytesInternal(),
// 	)
// 	var vdfResult VDFResult = VDFEval(vdfInput)

// 	return &TicketI{
// 		vrfResult,
// 		vdfResult,
// 	}
// }

func (self *Block_I) ValidateTickets() bool {
	panic("TODO")

	// for _, tix := range self.Tickets_ {
		// panic("TODO")
		// tix.Validate()
	// }

	return true
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
