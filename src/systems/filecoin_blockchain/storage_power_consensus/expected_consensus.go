package storage_power_consensus

import (
	"math/big"

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

func (self *ExpectedConsensus_I) IsWinningChallengeTicket(challengeTicket util.Bytes, maxTicket util.Bytes, sectorPower block.StoragePower, totalPower block.StoragePower, sampleNum util.UVarint, sampleDenom util.UVarint) bool {
	// Conceptually we are mapping the pseudorandom, deterministic hash output of the challenge ticket onto [0,1]
	// by dividing by 2^HashLen and comparing that to the sector's target.
	// if the challenge ticket hash is smaller than active power in the sector size / network total power * sectorSampled * ec.ExpectedLeaders
	// it is a winning challenge ticket.

	// lhs := challengeTicket * totalPower * sampleDenom
	// rhs := maxTicket * minerPower * sampleNum * self.expectedBlocksPerEpoch
	lhs := util.BigFromBytes(challengeTicket[:])
	lhs = lhs.Mul(lhs, util.BigFromUint64(uint64(totalPower)))
	lhs = lhs.Mul(lhs, util.BigFromUint64(uint64(sampleDenom)))

	// TODO: remove const here
	SHA256Len := 256
	// sectorPower * 2^len(H)
	rhs := new(big.Int).Lsh(util.BigFromUint64(uint64(sectorPower)), uint(SHA256Len))
	rhs = rhs.Mul(rhs, util.BigFromUint64(uint64(sampleNum)))
	rhs = rhs.Mul(rhs, big.NewInt(int64(self.expectedBlocksPerEpoch())))

	// lhs < rhs?
	return lhs.Cmp(rhs) == -1
}

func (self *ExpectedConsensus_I) GetBlockRewards(electionProof block.ElectionProof, minerPower block.StoragePower, totalPower block.StoragePower) util.UVarint {
	panic("")
}
