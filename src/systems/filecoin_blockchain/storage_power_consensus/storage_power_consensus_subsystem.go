package fileName

func (self *StoragePowerConsensusSubsystem) ValidateBlock(block Block) {

	// 1. Verify miner has not been slashed and is still valid miner
	if self.PowerTable.LookupMinerStorage(block.MinerAddress()) <= 0 {
		return self.StoragePowerConsensusError("block miner not valid")
	}

	// 2. Verify ParentWeight
	if block.ParentWeight != self.computeTipsetWeight(block.ParentTipset()) {
		return self.StoragePowerConsensusError("invalid parent weight")
	}

	// 3. Verify Tickets
	if !block.ValidateTickets() {
		return self.StoragePowerConsensusError("tickets were invalid")
	}

	// 4. Verify ElectionProof
	// randomnessLookbackTipset := RandomnessLookback(blk)
	// lookbackTicket := minTicket(randomnessLookbackTipset)
	// challenge := blake2b(lookbackTicket)

	// if !ValidateSignature(blk.ElectionProof, pubk, challenge) {
	// 	Fatal("election proof was not a valid signature of the last ticket")
	// }

	// powerLookbackTipset := PowerLookback(blk)

	// lbStorageMarket := LoadStorageMarket(powerLookbackTipset.state)
	// minerPower := lbStorageMarket.PowerLookup(blk.Miner)
	// totalPower := lbStorageMarket.GetTotalStorage()
	// if !IsProofAWinner(blk.ElectionProof, minerPower, totalPower) {
	// 	Fatal("election proof was not a winner")
	// }

	return nil
}

func (self *StoragePowerConsensusSubsystem) computeTipsetWeight(tipset *Tipset) ChainWeight {
	panic("TODO")
}

func (self *StoragePowerConsensusSubsystem) StoragePowerConsensusError(string errMsg) StoragePowerConsensusError {
	return Error(errMsg)
}
