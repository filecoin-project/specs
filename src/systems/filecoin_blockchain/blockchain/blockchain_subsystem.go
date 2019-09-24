package fileName

// The semantic stage requires access to the chain which the block extends.
func validateBlockSemantics(block Block) {

	// 1. Verify Signature
	 pubk := BlockchainSubsystem.StateTree().GetMinerKey(block.MinerAddress())
	 if block.BlockSig().Verify(pubk, block) {
		return ErrInvalidBlock("invalid block signature")
	 }
 
	 // 2. Verify Timestamp
	 if block.Timestamp() < block.ParentTipset().LatestTimestamp()
		|| block.Timestamp() > blockchainSubsystem.
	 {
		 Fatal("block was generated too far in the future")
	 }
	 // next check that it is appropriately delayed from its parents including
	 // null blocks.
	 if blk.GetTime() <= blk.minParentTime()+(BLOCK_DELAY*len(blk.Tickets)) {
		 Fatal("block was generated too soon")
	 }
 
	 // 3. Verify miner has not been slashed and is still valid miner
	 curStorageMarket := LoadStorageMarket(blk.State)
	 if !curStorageMarket.IsMiner(blk.Miner) {
		 Fatal("block miner not valid")
	 }
 
	 // 4. Verify ParentWeight
	 if blk.ParentWeight != ComputeWeight(blk.Parents) {
		 Fatal("invalid parent weight")
	 }
 
	 // 5. Verify Tickets
	 if !VerifyTickets(blk) {
		 Fatal("tickets were invalid")
	 }
 
	 // 6. Verify ElectionProof
	 randomnessLookbackTipset := RandomnessLookback(blk)
	 lookbackTicket := minTicket(randomnessLookbackTipset)
	 challenge := blake2b(lookbackTicket)
 
	 if !ValidateSignature(blk.ElectionProof, pubk, challenge) {
		 Fatal("election proof was not a valid signature of the last ticket")
	 }
 
	 powerLookbackTipset := PowerLookback(blk)
 
	 lbStorageMarket := LoadStorageMarket(powerLookbackTipset.state)
	 minerPower := lbStorageMarket.PowerLookup(blk.Miner)
	 totalPower := lbStorageMarket.GetTotalStorage()
	 if !IsProofAWinner(blk.ElectionProof, minerPower, totalPower) {
		 Fatal("election proof was not a winner")
	 }
 
	 // 7. Verify Message Signatures
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
 
	 // 8. Validate State Transitions
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

