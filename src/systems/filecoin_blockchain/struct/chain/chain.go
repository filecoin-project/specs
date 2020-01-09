package chain

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	util "github.com/filecoin-project/specs/util"
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
func (chain *Chain_I) RandomnessAtEpoch(epoch abi.ChainEpoch) util.Bytes {
	ts := chain.TipsetAtEpoch(epoch)
	return ts.MinTicket().DrawRandomness(epoch)
}

func (chain *Chain_I) GetTicketProductionRand(epoch abi.ChainEpoch) abi.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_TICKET)
}

func (chain *Chain_I) GetSealRand(epoch abi.ChainEpoch) abi.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_SEAL)
}

func (chain *Chain_I) GetPoStChallengeRand(epoch abi.ChainEpoch) abi.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_POST)
}
