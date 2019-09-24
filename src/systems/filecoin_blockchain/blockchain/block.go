package fileName

func (tix *Ticket) ValidateSyntax() bool {

	return VRFResult.validateSyntax()
		&& VDFResult.validateSyntax()
}

func (ep *ElectionProof) ValidateSyntax() bool {
	return VRFResult.validateSyntax()
}

func (ts *Tipset) ValidateSyntax() bool {

    if !len(ts.Blocks) > 0 {
        return false
    }

	parents := ts.Parents()
	grandparent := parents[0].Parents()
    for i := 1; i < len(parents); i++ {
        if grandparent != parents[i].Parents() {
            return false
		}
	}

	numTickets := len(ts.Blocks[0].Tickets())
    for i := 1; i < len(ts.Blocks); i++ {
        if numTickets != len(ts.Blocks[i].Tickets()) {
            return false
		}
	}

	return true
}

func (ts *Tipset) LatestTimestamp() {
	var latest Time
	for _, blk := ts.Blocks {
		if blk.Timestamp().After(latest) || latest.IsZero() {
			latest = blk.Timestamp()
		}
	}
	return latest
}

func (tipset *TipsetI) StateTree() StateTree {
	var currTree StateTree = nil
	for _, block := range tipset.Blocks() {
		if currTree == nil {
			currTree = block.StateTree()
		} else {
			Assert(block.StateTree().CID().Equals(currTree.CID()))
		}
	}
	Assert(currTree != nil)
	for _, block := range tipset.Blocks() {
		currTree = UpdateStateTree(currTree, block)
	}
	return currTree
}
