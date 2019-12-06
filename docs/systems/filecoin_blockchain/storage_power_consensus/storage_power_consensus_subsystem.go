package storage_power_consensus

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
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
}

func (spc *StoragePowerConsensusSubsystem_I) validateTicket(ticket block.Ticket, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {
	randomness1 := spc.GetTicketProductionRand(spc.blockchain().BestChain(), spc.blockchain().LatestEpoch())

	return ticket.Verify(randomness1, pk, minerActorAddr)
}

func (spc *StoragePowerConsensusSubsystem_I) ComputeChainWeight(tipset chain.Tipset) block.ChainWeight {
	return spc.ec().ComputeChainWeight(tipset)
}

func (spc *StoragePowerConsensusSubsystem_I) StoragePowerConsensusError(errMsg string) StoragePowerConsensusError {
	panic("TODO")
}

func (spc *StoragePowerConsensusSubsystem_I) IsWinningPartialTicket(stateTree stateTree.StateTree, partialTicket sector.PartialTicket, sectorUtilization block.StoragePower) bool {

	// finalize the partial ticket
	challengeTicket := filcrypto.SHA256(partialTicket)

	st := spc._getStoragePowerActorState(stateTree)
	activePower := st._getActivePower()

	// TODO: pull from constants
	SAMPLE_NUM := util.UVarint(1)
	SAMPLE_DENOM := util.UVarint(25)

	return spc.ec().IsWinningChallengeTicket(challengeTicket, sectorUtilization, activePower, SAMPLE_NUM, SAMPLE_DENOM)
}

// TODO: fix linking here
var node node_base.FilecoinNode

func (spc *StoragePowerConsensusSubsystem_I) _getStoragePowerActorState(stateTree stateTree.StateTree) StoragePowerActorState {
	powerAddr := addr.StoragePowerActorAddr
	actorState := stateTree.GetActorState(powerAddr)
	substateCID := actorState.State()

	substate, err := node.LocalGraph().Get(ipld.CID(substateCID))
	if err != nil {
		panic("TODO")
	}

	// TODO fix conversion to bytes
	panic(substate)
	var serializedSubstate util.Serialization
	st, err := Deserialize_StoragePowerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (spc *StoragePowerConsensusSubsystem_I) GetTicketProductionRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_TICKET)
}

func (spc *StoragePowerConsensusSubsystem_I) GetSealRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_SEAL)
}

func (spc *StoragePowerConsensusSubsystem_I) GetPoStChallengeRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - SPC_LOOKBACK_POST)
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
