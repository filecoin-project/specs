package storage_power_consensus

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	util "github.com/filecoin-project/specs/util"
)

func (self *ExpectedConsensus_I) ComputeChainWeight(tipset block.Tipset) block.ChainWeight {
	panic("")
	// see expected_consensus.md for detail

	// wPowerFactor := self.log2b(spa.GetTotalPower())
	// wBlocksFactor_num := (wPowerFactor * len(tipset.Blocks) * self.wParams.wRatio_num)
	// wBlocksFactor_den := self.expectedBlocksPerEpoch * self.wParams.wRatio_den
	// return tipset.ParentTipset.ChainWeight
	// 	+wPowerFactor * self.wParams.wPrecision
	// 	+(wBlocksFactor_num * self.wParams.wPrecision / wBlocksFactor_den)
}

func (self *ExpectedConsensus_I) IsValidConsensusFault(faults ConsensusFaultType, blocks []block.Block) bool {
	panic("TODO")

	// 1. double-fork mining fault
	// return block1.Miner == block2.Miner && block1.Epoch == block2.Epoch

	// 2. time-offset mining fault
	// return block1.Miner == block2.Miner
	// && block1.ParentTipset == block2.ParentTipset

	// 3. parent grinding fault
	// return block1.Miner == block2.Miner
	// && abs(block1.Epoch - block2.Epoch) == 1
}

func (self *ExpectedConsensus_I) IsWinningElectionProof(electionProof block.ElectionProof, minerPower block.StoragePower, totalPower block.StoragePower) bool {
	panic("")
	// Conceptually we are mapping the pseudorandom, deterministic VRFOutput onto [0,1]
	// by dividing by 2^HashLen and comparing that to the miner's power (portion of network storage).
	// if the VRF Output is smaller than the miner's power fraction * expected number of blocks per round
	// it is a winning election proof.

	// return electionProof.Output()*totalPower < self.expectedBlocksPerEpoch*minerPower*electionProof.VRFResult_.MaxValue()
}

func (self *ExpectedConsensus_I) GetBlockRewards(electionProof block.ElectionProof, minerPower block.StoragePower, totalPower block.StoragePower) util.UVarint {
	panic("")
	// draw := electionProof.output / electionProof.VRFResult_.MaxValue()
	// req := self.expectedBlocksPerEpoch * minerPower / totalPower
	// rewardCount := ceil(req - draw)
	// reward := rewardCount * self.expectedRewardPerEpoch / self.expectedBlocksPerEpoch
	// return reward
}
