package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

const FINALITY = 500

const (
	SPC_LOOKBACK_RANDOMNESS = 300      // this is EC.K maybe move it there. TODO
	SPC_LOOKBACK_TICKET     = 1        // we chain blocks together one after the other
	SPC_LOOKBACK_POST       = 1        // cheap to generate, should be set as close to current TS as possible
	SPC_LOOKBACK_SEAL       = FINALITY // should be set to finality
)

// Storage Power Consensus Subsystem

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block block.Block_I) error {
	panic("")
	// minerPK := PowerTable.GetMinerPublicKey(block.MinerAddress())
	// minerPower := PowerTable.GetMinerPower(block.MinerAddress())

	// // 1. Verify miner has not been slashed and is still valid miner
	// if minerPower <= 0 {
	// 	return StoragePowerConsensusError("block miner not valid")
	// }

	// // 2. Verify ParentWeight
	// if block.Weight() != computeTipsetWeight(block.Parents()) {
	// 	return errors.New("invalid parent weight")
	// }

	// // 3. Verify Tickets
	// if !validateTicket(block.Ticket, minerPK) {
	// 	return StoragePowerConsensusError("ticket was invalid")
	// }

	// // 4. Verify ElectionProof construction
	// if !ValidateElectionProof(block.Height, block.ElectionProof, block.MinerAddress) {
	// 	return StoragePowerConsensusError("election proof was not a valid signature of the last ticket")
	// }

	// // 5. and value
	// if !IsWinningElectionProof(block.ElectionProof, spa.GetMinerPower(), spa.GetTotalPower()) {
	// 	return StoragePowerConsensusError("election proof was not a winner")
	// }

	// return nil
}

func (spc *StoragePowerConsensusSubsystem_I) validateTicket(ticket block.Ticket, pk filcrypto.VRFPublicKey, minerAddr addr.Address) bool {
	randomness1 := spc.GetTicketProductionSeed(spc.blockchain().BestChain(), spc.blockchain().LatestEpoch())

	var input []byte
	input = append(input, byte(filcrypto.DomainSeparationTag_Case_Ticket))
	input = append(input, randomness1...)
	input = append(input, byte(filcrypto.InputDelimeter_Case_Bytes))
	input = append(input, minerAddr.ToBytes()...)
	return ticket.Verify(input, pk, minerAddr)
}

func (spc *StoragePowerConsensusSubsystem_I) ComputeChainWeight(tipset block.Tipset) block.ChainWeight {
	return spc.ec().ComputeChainWeight(tipset)
}

func (spc *StoragePowerConsensusSubsystem_I) StoragePowerConsensusError(errMsg string) StoragePowerConsensusError {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) IsWinningPartialTicket(partialTicket sector.PartialTicket) bool {

	// TODO: Send message to SPA to get ActivePowerInSector associated with sector
	activePowerInSector := uint64(block.StoragePower(0))

	// TODO: Send message to SPA
	totalPower := uint64(block.StoragePower(0))

	challengeTicket := block.SHA256(partialTicket)

	return spc.ec().IsWinningChallengeTicket(challengeTicket, activePowerInSector, totalPower)
}

func (spc *StoragePowerConsensusSubsystem_I) GetTicketProductionSeed(chain block.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_TICKET)
}

func (spc *StoragePowerConsensusSubsystem_I) GetElectionProofSeed(chain block.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_RANDOMNESS)
}

func (spc *StoragePowerConsensusSubsystem_I) GetSealSeed(chain block.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_SEAL)
}

func (spc *StoragePowerConsensusSubsystem_I) GetPoStChallenge(chain block.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_POST)
}

func (spc *StoragePowerConsensusSubsystem_I) ValidateElectionProof(height block.ChainEpoch, electionProof block.ElectionProof, minerAddr addr.Address) bool {
	panic("")
	// // 1. Check that ElectionProof was validated in appropriate time
	// if height > clock.roundTime {
	// 	return false
	// }

	// // 2. Determine that ticket was validly scratched
	// minerPK := spc.PowerTable.GetMinerPublicKey(workerAddr)
	// input := VRFPersonalizationElectionProof
	// TK := storagePowerConsensus.GetElectionProofSeed(sms.CurrentChain, sms.block.LatestEpoch())
	// input.append(TK.Output)
	// input.append(height)

	// return electionProof.Verify(input, minerPK)
}

func (spc *StoragePowerConsensusSubsystem_I) GetFinality() block.ChainEpoch {
	panic("")
	// return FINALITY
}

func (spc *StoragePowerConsensusSubsystem_I) FinalizedEpoch() block.ChainEpoch {
	panic("")
	// currentEpoch := rt.HeadEpoch()
	// return currentEpoch - spc.GetFinality()
}
