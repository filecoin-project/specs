package chain

import (
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	builtin "github.com/filecoin-project/specs-actors/actors/builtin"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
)

// Returns the tipset at or immediately prior to `epoch`.
func (chain *Chain_I) TipsetAtEpoch(epoch abi.ChainEpoch) Tipset {
	current := chain.HeadTipset()
	for current.Epoch() > epoch {
		current = current.Parents()
	}

	return current
}

// Draws randomness from the tipset at or immediately prior to `epoch`.
func (chain *Chain_I) RandomnessAtEpoch(epoch abi.ChainEpoch) abi.RandomnessSeed {
	if epoch < genesis {
		genesisTS := chain.TipsetAtEpoch(genesis)
		genesisTix := genesisTS.MinTicket().DrawRandomness(genesis)
		buffer := []byte{}
		buffer = append(buffer, genesisTix...)
		buffer = append(buffer, BigEndianBytesFromInt(int64(epoch))...)
		return blake2b.Sum256(buffer)
	} else {
		ts := chain.TipsetAtEpoch(epoch)
		return ts.MinTicket().DrawRandomness()
	}
}

func (chain *Chain_I) GetTicketProductionRandSeed(epoch abi.ChainEpoch) abi.RandomnessSeed {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_TICKET)
}

func (chain *Chain_I) GetSealRandSeed(epoch abi.ChainEpoch) abi.RandomnessSeed {
	return chain.RandomnessAtEpoch(epoch - builtin.SPC_LOOKBACK_SEAL)
}

func (chain *Chain_I) GetPoStChallengeRandSeed(epoch abi.ChainEpoch) abi.RandomnessSeed {
	return chain.RandomnessAtEpoch(epoch - builtin.SPC_LOOKBACK_POST)
}
