package storage_power_consensus

import (
	"errors"

	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

const (
	SPC_LOOKBACK_RANDOMNESS = 300 // this is EC.K maybe move it there. TODO
	SPC_LOOKBACK_TICKET     = 1   // we chain blocks together one after the other
)

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block block.BlockHeader) error {

	// 1. Verify miner has not been slashed and is still valid miner
	if spc.powerTable().LookupMinerStorage(block.MinerAddress()) <= 0 {
		return errors.New("block miner not valid")
	}

	minerPK := filcrypto.PubKey{} // TODO get the key from the state tree spc.StorageMiningSubsystem.GetMinerKeyByAddress(block.MinerAddress())
	// 2. Verify ParentWeight
	if block.Weight() != spc.computeTipsetWeight(block.Parents()) {
		return errors.New("invalid parent weight")
	}

	// 3. Verify Tickets
	if !block.ValidateTickets(minerPK) {
		return errors.New("tickets were invalid")
	}

	// 4. Verify ElectionProof construction
	seed := block.ExtractElectionSeed()
	if !block.ElectionProof().Verify(seed, minerPK) {
		return errors.New("election proof was not a valid signature of the last ticket")
	}

	// and value
	minerPower := spc.powerTable().LookupMinerPower(block.MinerAddress())
	totalPower := spc.powerTable().GetTotalPower()
	if !block.ElectionProof().IsWinning(float(minerPower) / totalPower) {
		return spc.StoragePowerConsensusError("election proof was not a winner")
	}

	return nil
}

func (spc *StoragePowerConsensusSubsystem_I) computeTipsetWeight(tipset block.Tipset) block.ChainWeight {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) TicketAtEpoch(chain chain.Chain, epoch block.ChainEpoch) block.Ticket {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) GetElectionArtifacts(chain chain.Chain, epoch block.ChainEpoch) block.ElectionArtifacts {
	return &block.ElectionArtifacts_I{
		TK_: spc.TicketAtEpoch(chain, epoch-SPC_LOOKBACK_RANDOMNESS),
		T1_: spc.TicketAtEpoch(chain, epoch-SPC_LOOKBACK_TICKET),
	}
}

func (pt *PowerTable_I) LookupMinerStorage(addr addr.Address) util.UVarint {
	panic("")
}
func (pt *PowerTable_I) LookupMinerPowerFraction(addr addr.Address) block.PowerFraction {
	panic("")
}
func (pt *PowerTable_I) RemovePower(addr addr.Address) {
	panic("")
}
