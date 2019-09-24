package fileName

// The semantic stage requires access to the chain which the block extends.
func (self *BlockchainSubsystem) validateBlockSemantics(block Block) {

	// 1. Verify Signature
	pubk := self.StateTree().GetMinerKey(block.MinerAddress())
	if block.BlockSig().Verify(pubk, block) {
		return ErrInvalidBlock("invalid block signature")
	}

	// 2. Verify Timestamp
	if block.Timestamp() < block.ParentTipset().LatestTimestamp() {
		return ErrInvalidBlock("block was generated too far in the past")
	}
	// next check that it is mined within the allowed time in the current epoch
	if block.Timestamp() > self.epochCutoffTime() {
		return ErrInvalidBlock("block was generated too far in the future")
	}
	// 3. Verify SPC dependencies: miner, weights, tickets, EP
	SPCErr := self.SPCSubsystem.ValidateBlock(block)
	if SPCErr {
		return SPCErr
	}

	// 4. Verify Message Signatures
	messages := LoadMessages(blk.Messages)
	state := GetParentState(blk.Parents)

	var blsMessages []Message
	var blsPubKeys []PublicKey
	for i, msg := range messages {
		if IsBlsMessage(msg) {
			blsMessages.append(msg)
			blsPubKeys.append(state.LookupPublicKey(msg.From))
		} else {
			if !ValidateSignature(msg) {
				Fatal("invalid message signature in block")
			}
		}
	}

	ValidateBLSSignature(blk.BLSAggregate, blsMessages, blsPubKeys)

	// 5. Validate State Transitions
	receipts := LoadReceipts(blk.MessageReceipts)
	for i, msg := range messages {
		receipt := ApplyMessage(state, msg)
		if receipt != receipts[i] {
			Fatal("message receipt mismatch")
		}
	}
	if state.Cid() != blk.StateRoot {
		Fatal("state roots mismatch")
	}
}

func (state StateTree) LookupPublicKey(a Address) PubKey {
	act := state.GetActor(a)
	if !act.Code == AccountActor {
		Fatal("only account actors have public keys")
	}

	ast := LoadAccountActorState(act)
	if act.Address.Type == BLS {
		return ExtractBLSPubKey(act.Address)
	}
	Fatal("can only look up public keys for BLS controlled accounts")
}

func (self *BlockchainSubsystem) epochCutoffTime() {
	panic("TODO")
}

func (self *Address) verifySyntax(addrType Address_Protocol) bool {
	panic("TODO")
	// switch aType := addrType; aType {
	// case Address_Protocol.Secp256k1():
	// 	// 80 Bytes
	// 	return len(self)
	// case Address_Protocol.ID():
	// 	// ?
	// case Address_Protocol.Actor():
	// 	// Blake2b - 64 Bytes
	// case Address_Protocol.BLS():
	// 	// BLS-12_381 - 48 Byte PK
	// }
}
