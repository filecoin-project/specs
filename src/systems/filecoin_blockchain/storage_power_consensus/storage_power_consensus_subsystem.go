package storage_power_consensus

import (
	"errors"

	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
	base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
	clock "github.com/filecoin-project/specs/systems/filecoin_nodes/clock"
)

const FINALITY = 500

const (
	SPC_LOOKBACK_RANDOMNESS = 300 // this is EC.K maybe move it there. TODO
	SPC_LOOKBACK_TICKET     = 1   // we chain blocks together one after the other
	SPC_LOOKBACK_POST       = 1   // cheap to generate, should be set as close to current TS as possible
	SPC_LOOKBACK_SEAL       = FINALITY  // should be set to finality
)

const VRFPersonalization {
	ticket: 0,
	electionProof: 1
}

// Storage Power Consensus Subsystem

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block blockchain.Block_I) error {
	minerPK := spc.PowerTable.GetMinerPublicKey(block.MinerAddress())
	minerPower := spc.PowerTable.GetMinerPower(block.MinerAddress())

	// 1. Verify miner has not been slashed and is still valid miner
	if minerPower <= 0 {
		return spc.StoragePowerConsensusError("block miner not valid")
	}

	// 2. Verify ParentWeight
	if block.Weight() != spc.computeTipsetWeight(block.Parents()) {
		return errors.New("invalid parent weight")
	}

	// 3. Verify Tickets
	if !validateTicket(block.Ticket, minerPK) {
		return spc.StoragePowerConsensusError("ticket was invalid")
	}

	// 4. Verify ElectionProof construction
	if !spc.ValidateElectionProof(block.Height, block.ElectionProof, block.MinerAddress) {
		return spc.StoragePowerConsensusError("election proof was not a valid signature of the last ticket")
	}

	// 5. and value
	if !spc.IsWinningElectionProof(block.ElectionProof, block.MinerAddress)
		return spc.StoragePowerConsensusError("election proof was not a winner")
	}

	return nil
}

func (spc *StoragePowerConsensusSubsystem_I) validateTicket(ticket base.Ticket, pk filcrypto.PublicKey) bool {
	T1 := storagePowerConsensus.GetTicketProductionSeed(sms.CurrentChain, sms.Blockchain.LatestEpoch())
	input := VRFPersonalization.Ticket
	input.append(T1.Output)
	return ticket.Verify(input, pk)
}

func (spc *StoragePowerConsensusSubsystem_I) computeTipsetWeight(tipset blockchain.Tipset) base.ChainWeight {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) StoragePowerConsensusError(errMsg string) base.StoragePowerConsensusError {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) GetTicketProductionSeed(chain blockchain.Chain, epoch Epoch) base_mining.SealSeed {
	return &base_mining.SealSeed_I {
		chain.TicketAtEpoch(epoch-SPC_LOOKBACK_TICKET)
	}
}

func (spc *StoragePowerConsensusSubsystem_I) GetElectionProofSeed(chain blockchain.Chain, epoch Epoch) base_mining.SealSeed {
	return &base_mining.SealSeed_I {
		chain.TicketAtEpoch(epoch-SPC_LOOKBACK_RANDOMNESS),
	}
}

func (spc *StoragePowerConsensusSubsystem_I) GetSealSeed(chain blockchain.Chain, epoch Epoch) base_mining.SealSeed {
	return &base_mining.SealSeed_I {
		chain.TicketAtEpoch(epoch-SPC_LOOKBACK_SEAL)
	}
}

func (spc *StoragePowerConsensusSubsystem_I) GetPoStChallenge(chain blockchain.Chain, epoch Epoch) base_mining.PoStChallenge {
	return &base_mining.PoStChallenge_I {
		chain.TicketAtEpoch(epoch-SPC_LOOKBACK_POST)
	}
}

func (spc *StoragePowerConsensusSubsystem_I) ValidateElectionProof(height BlockHeight, electionProof base.ElectionProof, workerAddr base.Address) bool {
	// 1. Check that ElectionProof was validated in appropriate time
	if (height + electionProof.ElectionNonce > clock.roundTime + 1) {
		return false
	}

	// 2. Determine that ticket was validly scratched
	minerPK := spc.PowerTable.GetMinerPublicKey(workerAddr)
	input := VRFPersonalization.ElectionProof
	TK := storagePowerConsensus.GetElectionProofSeed(sms.CurrentChain, sms.Blockchain.LatestEpoch())
	input.append(TK.Output)
	input.appent(electionProof.ElectionNonce)

	return electionProof.Verify(input, minerPK)
}

func (spc *StoragePowerConsensusSubsystem_I) IsWinningElectionProof(electionProof base.ElectionProof, workerAddr base.Address) bool {
	// 1. Determine miner power fraction
	minerPower := spc.PowerTable.GetMinerPower(workerAddr)
	totalPower := spc.PowerTable.GetTotalPower()

	// Conceptually we are mapping the pseudorandom, deterministic VRFOutput onto [0,1]
	// by dividing by 2^HashLen (64 Bytes using Sha256) and comparing that to the miner's
	// power (portion of network storage).
	return (minerPower * 2^(len(electionProof.Output)*8) < electionProof.Output * totalPower)
}

func (spc *StoragePowerConsensusSubsystem_I) GetFinality() Epoch {
	return FINALITY
}