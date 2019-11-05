package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (self *ExpectedConsensus_I) IsValidConsensusFault(faults ConsensusFaultType, blocks []block.Block) bool {
	panic("TODO")
}

func (self *ExpectedConsensus_I) IsWinningElectionProof(electionProof block.ElectionProof, workerAddr addr.Address) bool {
	panic("")
	// 1. Determine miner power fraction
	// minerPower := spc.PowerTable.GetMinerPower(workerAddr)
	// totalPower := spc.PowerTable.GetTotalPower()

	// // Conceptually we are mapping the pseudorandom, deterministic VRFOutput onto [0,1]
	// // by dividing by 2^HashLen (64 Bytes using Sha256) and comparing that to the miner's
	// // power (portion of network storage).
	// return (minerPower*2^(len(electionProof.Output)*8) < electionProof.Output*totalPower)
	return true
}
func (self *ExpectedConsensus_I) GetBlockRewards(electionProof ElectionProof, workerAddr addr.Address) UVarint {
	panic("")
	// draw := electionProof.output / electionProof.VRFResult_.MaxValue()
	// req := self.expectedMinersPerRound * spc.PowerTable.GetMinerPower(workerAddr) / spc.PowerTable.GetTotalPower()
	// rewardCount := ceil(req - draw)
	// reward := rewardCount * self.expectedRewardPerRound / self.expectedMinersPerRound
	// return reward
}

func (tix *Ticket_I) ValidateSyntax() bool {
	return tix.VRFResult_.ValidateSyntax()
}

func (tix *Ticket_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return tix.VRFResult_.Verify(input, pk)
}

func (ep *ElectionProof_I) ValidateSyntax() bool {
	return ep.VRFResult_.ValidateSyntax()
}

func (ep *ElectionProof_I) Verify(input util.Bytes, pk filcrypto.VRFPublicKey) bool {
	return ep.VRFResult_.Verify(input, pk)
}
