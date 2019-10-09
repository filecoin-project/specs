package storage_power_consensus

import (
	base "github.com/filecoin-project/specs/systems/filecoin_blockchain"
	blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain/blockchain"
)

type Block blockchain.Block

const (
	SPC_LOOKBACK_RANDOMNESS = 300 // this is EC.K maybe move it there. TODO
	SPC_LOOKBACK_TICKET     = 1   // we chain blocks together one after the other
)

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block Block) error {

	// 1. Verify miner has not been slashed and is still valid miner
	if spc.powerTable().LookupMinerStorage(block.MinerAddress()) <= 0 {
		return spc.StoragePowerConsensusError("block miner not valid")
	}

	minerPK := spc.StorageMiningSubsystem.GetMinerKeyByAddress(block.MinerAddress())
	// 2. Verify ParentWeight
	if block.ParentWeight != spc.computeTipsetWeight(block.ParentTipset()) {
		return spc.StoragePowerConsensusError("invalid parent weight")
	}

	// 3. Verify Tickets
	if !block.ValidateTickets(minerPK) {
		return spc.StoragePowerConsensusError("tickets were invalid")
	}

	// 4. Verify ElectionProof construction
	seed := block.ParentTipset().ExtractElectionSeed()
	if !block.ElectionProof.Validate(seed, minerPK) {
		return spc.StoragePowerConsensusError("election proof was not a valid signature of the last ticket")
	}

	// and value
	minerPower := spc.PowerTable.LookupMinerPowerFraction(block.MinerAddress)
	if !block.ElectionProof.IsWinning(minerPower) {
		return spc.StoragePowerConsensusError("election proof was not a winner")
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

func (pt *PowerTable_I) LookupMinerStorage(addr base.Address) UVarint {
	panic("")
}
func (pt *PowerTable_I) LookupMinerPowerFraction(addr base.Address) Float {
	panic("")
}
func (pt *PowerTable_I) RemovePower(addr base.Address) {
	panic("")
}
