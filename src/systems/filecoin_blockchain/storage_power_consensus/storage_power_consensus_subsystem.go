package storage_power_consensus

import (
	base "github.com/filecoin-project/specs/systems/filecoin_blockchain"
	blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain/blockchain"
)

const (
	SPC_LOOKBACK_RANDOMNESS = 300 // this is EC.K maybe move it there. TODO
	SPC_LOOKBACK_TICKET     = 1   // we chain blocks together one after the other
)

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block Block) {

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

func (spc *StoragePowerConsensusSubsystem_I) computeTipsetWeight(tipset *Tipset) ChainWeight {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) StoragePowerConsensusError(string errMsg) StoragePowerConsensusError {
	return Error(errMsg)
}

func (spc *StoragePowerConsensusSubsystem_I) GetElectionArtifacts(chain blockchain.Chain, epoch base.Epoch) base.ElectionArtifacts {
	return base.ElectionArtifacts{
		TK: spc.TicketAtEpoch(chain, epoch-SPC_LOOKBACK_RANDOMNESS),
		T1: spc.TicketAtEpoch(chain, epoch-SPC_LOOKBACK_TICKET),
	}
}
