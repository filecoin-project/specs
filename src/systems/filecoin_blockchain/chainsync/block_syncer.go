package chainsync

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

func (self *BlockSyncer_I) OnNewBlock(block block.Block) error {
	panic("TODO")
	// err := self.validateBlockSyntax(block)
	// if err {
	// 	return err
	// }

	// self.blockchainSubsystem_.HandleBlock(block)
	return nil
}

// The syntactic stage may be validated without reference to additional data (see block)
func (bs *BlockSyncer_I) validateBlockSyntax(block block.Block) error {
	panic("TODO")
	// if !block.MinerAddress().VerifySyntax(StorageMinerActor.Address.Protocol()) {
	// 	return ErrInvalidBlockSyntax("bad miner address syntax")
	// }

	// if !(len(block.Tickets) > 0) {
	// 	return ErrInvalidBlockSyntax("no tickets")
	// }

	// for _, tix := range block.Tickets {
	// 	if !tix.ValidateSyntax() {
	// 		return ErrInvalidBlockSyntax("bad ticket syntax")
	// 	}
	// }

	// if !block.ElectionProof.ValidateSyntax() {
	// 	return ErrInvalidBlockSyntax("bad election proof syntax")
	// }

	// if !block.ParentTipset().ValidateSyntax() {
	// 	return ErrInvalidBlockSyntax("invalid parent tipset")
	// }

	// if !block.ParentWeight() > 0 {
	// 	return ErrInvalidBlockSyntax("parent weight < 0")
	// }

	// if !block.Height() > 0 {
	// 	return ErrInvalidBlockSyntax("height < 0")
	// }

	// // if !block.StateTree().ValidateSyntax() {
	// //     return false
	// // }

	// for _, msg := range block.Messages() {
	// 	if !msg.ValidateSyntax() {
	// 		return ErrInvalidBlockSyntax("msg syntax invalid")
	// 	}
	// }

	// // TODO msg receipts

	// if block.Timestamp_ > bs.blockchainSubsystem.Clock.Now() {
	// 	return ErrInvalidBlockSyntax("bad timestamp")
	// }

	return nil

}

func (g *AncestryGraph_I) AddBlock(block block.Block) {
	panic("TODO")
}

func (g *BlockValidationGraph_I) ConsiderBlock(block block.Block) {
	panic("TODO")
	// g.UnconnectedBlocks.AddBlock(block)
	// g.tryConnectBlockToFringe(block)
}

func (g *BlockValidationGraph_I) tryConnectBlockToFringe(block block.Block) {
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

func (b *BlockPropagator_I) PropagateBlock(block block.SignedBlock) {
	panic("TODO")
	// NetworkBroadcast(BlockPubSub, block)
}
