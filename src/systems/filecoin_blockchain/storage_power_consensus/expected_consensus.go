package storage_power_consensus

import (
	"math/big"

	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	util "github.com/filecoin-project/specs/util"
)

func (self *ExpectedConsensus_I) ComputeChainWeight(tipset chain.Tipset) block.ChainWeight {
	util.IMPL_FINISH()
	return block.ChainWeight(0)
	// see expected_consensus.md for detail

	// wPowerFactor := self.log2b(spa.GetTotalPower())
	// wBlocksFactor_num := (wPowerFactor * len(tipset.Blocks) * self.wParams.wRatio_num)
	// wBlocksFactor_den := self.expectedLeadersPerEpoch * self.wParams.wRatio_den
	// return tipset.ParentTipset.ChainWeight
	// 	+wPowerFactor * self.wParams.wPrecision
	// 	+(wBlocksFactor_num * self.wParams.wPrecision / wBlocksFactor_den)
}

func (self *ExpectedConsensus_I) IsValidConsensusFault(faults ConsensusFaultType, blocks []block.Block) bool {
	util.IMPL_FINISH()
	return false
	// 1. double-fork mining fault
	// return block1.Miner == block2.Miner && block1.Epoch == block2.Epoch

	// 2. time-offset mining fault
	// return block1.Miner == block2.Miner
	// && block1.Parents == block2.Parents

	// 3. parent grinding fault
	// return block1.Miner == block2.Miner
	// && abs(block1.Epoch - block2.Epoch) == 1
}

func (self *ExpectedConsensus_I) IsWinningChallengeTicket(challengeTicket util.Bytes, sectorPower block.StoragePower, networkPower block.StoragePower, numSectorsSampled util.UVarint, numSectorsMiner util.UVarint) bool {
	// Conceptually we are mapping the pseudorandom, deterministic hash output of the challenge ticket onto [0,1]
	// by dividing by 2^HashLen and comparing that to the sector's target.
	// if the challenge ticket hash / max hash val < sectorPower / totalPower * ec.ExpectedLeaders * numSectorsMiner / numSectorsSampled
	// it is a winning challenge ticket.
	// note that the sectorPower may differ based on the challenged sector

	// lhs := challengeTicket * totalPower * numSectorsSampled
	// rhs := maxTicket * activeSectorPower * numSectorsMiner * self.ExpectedLeaders
	lhs := util.BigFromBytes(challengeTicket[:])
	lhs = lhs.Mul(lhs, util.BigFromUint64(uint64(networkPower)))
	lhs = lhs.Mul(lhs, util.BigFromUint64(uint64(numSectorsSampled)))

	// TODO: remove const here
	SHA256Len := 256
	// sectorPower * 2^len(H)
	rhs := new(big.Int).Lsh(util.BigFromUint64(uint64(sectorPower)), uint(SHA256Len))
	rhs = rhs.Mul(rhs, util.BigFromUint64(uint64(numSectorsMiner)))
	rhs = rhs.Mul(rhs, big.NewInt(int64(self.expectedLeadersPerEpoch())))

	// lhs < rhs?
	return lhs.Cmp(rhs) == -1
}
