---
title: "Block Reception"
---

func (g *BlockValidationGraph_I) ConsiderBlock(block Block) {
	panic("TODO")
	// g.UnconnectedBlocks.AddBlock(block)
	// g.tryConnectBlockToFringe(block)
}

func (g *BlockValidationGraph_I) tryConnectBlockToFringe(block Block) {
	panic("TODO")

	// try to connect the block, and then try connecting its descendents.
	//
	// this algorithm should be breadth-first because we need to process the fringe
	// in order. Depth-first search may consider blocks whose parents are still
	// yet to be added
	// blocks := Queue < Block >
	// 	blocks.Enqueue(block)
	// for block := range blocks.Dequeue() {
	// 	if !g.ValidationFringe.HasTipset(block.Parents()) {
	// 		continue // ignore this one. not all of its parents are in fringe
	// 	}

	// 	children := g.UnconnectedBlocks.Children[block]
	// 	g.UnconnectedBlocks.Parents.Remove(block)
	// 	g.UnconnectedBlocks.Children.Remove(block)
	// 	g.ValidationFringe.AddBlock(block)
	// 	blocks.EnqueueAll(children)
	// }
}