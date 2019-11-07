package filecoin_blockchain

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// The semantic stage requires access to the chain which the block extends.
func (self *BlockchainSubsystem_I) validateBlockSemantics(block block.Block) {
	panic("TODO")
	// // 1. Verify Signature
	// pubk := self.StateTree().GetMinerKey(block.MinerAddress())
	// msg := append([]byte("BLOCK"), block...)
	// if block.BlockSig().Verify(pubk, msg) {
	// 	return ErrInvalidBlock("invalid block signature")
	// }

	// // 2. Verify Timestamp
	// if block.Timestamp() < block.ParentTipset().LatestTimestamp() {
	// 	return ErrInvalidBlock("block was generated too far in the past")
	// }
	// // next check that it is mined within the allowed time in the current epoch
	// if block.Timestamp() > self.epochCutoffTime() {
	// 	return ErrInvalidBlock("block was generated too far in the future")
	// }
	// // 3. Verify SPC dependencies: miner, weights, tickets, EP
	// SPCErr := self.SPCSubsystem.ValidateBlock(block)
	// if SPCErr {
	// 	return SPCErr
	// }

	// // 4. Verify Message Signatures
	// messages := LoadMessages(blk.Messages)
	// state := GetParentState(blk.Parents)

	// var blsMessages []Message
	// var blsPubKeys []PublicKey
	// for i, msg := range messages {
	// 	if IsBlsMessage(msg) {
	// 		blsMessages.append(msg)
	// 		blsPubKeys.append(state.LookupPublicKey(msg.From))
	// 	} else {
	// 		if !ValidateSignature(msg) {
	// 			Fatal("invalid message signature in block")
	// 		}
	// 	}
	// }

	// ValidateBLSSignature(blk.BLSAggregate, blsMessages, blsPubKeys)

	// // 5. Validate State Transitions
	// receipts := LoadReceipts(blk.MessageReceipts)
	// for i, msg := range messages {
	// 	receipt := ApplyMessage(state, msg)
	// 	if receipt != receipts[i] {
	// 		Fatal("message receipt mismatch")
	// 	}
	// }
	// if state.Cid() != blk.StateRoot {
	// 	Fatal("state roots mismatch")
	// }
}

func (self *BlockchainSubsystem_I) epochCutoffTime() {
	panic("TODO")
}

func (self *BlockchainSubsystem_R) HandleBlock(block block.Block) bool {
	panic("TODO")
}

func HandleBlock(block block.Block) bool {
	panic("TODO")
}

// LatestEpoch() Epoch
// BestChain() Chain
// CandidateChains() []Chain
// NewTipsets() BlockchainSubsystem_NewTipsets_FunRet
// NewBestTipset() BlockchainSubsystem_NewTipsets_FunRet
// SyncState() BlockchainSubsystem_SyncState_FunRet
// VerifySectorExists(sectorId SectorID) bool
// ValidateBlock(block Block) bool
// TryGenerateStateTree(block Block) stateTree.StateTree
// AssembleTipsets() []Tipset
// ChooseTipset(tipsets []Tipset) Tipset
// GetPostRandomness() Randomness
