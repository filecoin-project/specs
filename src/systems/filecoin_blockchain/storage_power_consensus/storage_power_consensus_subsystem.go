package storage_power_consensus

func (self *StoragePowerConsensusSubsystem) ValidateBlock(block Block) {

	// 1. Verify miner has not been slashed and is still valid miner
	if self.PowerTable.LookupMinerStorage(block.MinerAddress()) <= 0 {
		return self.StoragePowerConsensusError("block miner not valid")
	}

	minerPK := self.StorageMiningSubsystem.GetMinerKeyByAddress(block.MinerAddress())
	// 2. Verify ParentWeight
	if block.ParentWeight != self.computeTipsetWeight(block.ParentTipset()) {
		return self.StoragePowerConsensusError("invalid parent weight")
	}

	// 3. Verify Tickets
	if !block.ValidateTickets(minerPK) {
		return self.StoragePowerConsensusError("tickets were invalid")
	}

	// 4. Verify ElectionProof construction
	seed := block.ParentTipset().ExtractElectionSeed()
	if !block.ElectionProof.Validate(seed, minerPK) {
		return self.StoragePowerConsensusError("election proof was not a valid signature of the last ticket")
	}

	// and value
	minerPower := self.PowerTable.LookupMinerPowerFraction(block.MinerAddress)
	if !block.ElectionProof.IsWinning(minerPower) {
		return self.StoragePowerConsensusError("election proof was not a winner")
	}

	return nil
}

func (self *StoragePowerConsensusSubsystem) computeTipsetWeight(tipset *Tipset) ChainWeight {
	panic("TODO")
}

func (self *StoragePowerConsensusSubsystem) StoragePowerConsensusError(string errMsg) StoragePowerConsensusError {
	return Error(errMsg)
}
