package storage_power_consensus

import (
	"math"

	spowact "github.com/filecoin-project/specs/actors/builtin/storage_power"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	inds "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

// Storage Power Consensus Subsystem

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block block.Block_I) error {
	util.IMPL_FINISH()
	return nil
}

func (spc *StoragePowerConsensusSubsystem_I) validateTicket(ticket block.Ticket, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {
	randomness1 := spc.GetTicketProductionRand(spc.blockchain().BestChain(), spc.blockchain().LatestEpoch())

	return ticket.Verify(randomness1, pk, minerActorAddr)
}

func (spc *StoragePowerConsensusSubsystem_I) ComputeChainWeight(tipset chain.Tipset) block.ChainWeight {
	return spc.ec().ComputeChainWeight(tipset)
}

func (spc *StoragePowerConsensusSubsystem_I) IsWinningPartialTicket(stateTree stateTree.StateTree, inds inds.Indices, partialTicket sector.PartialTicket, sectorUtilization block.StoragePower, numSectors util.UVarint) bool {

	// finalize the partial ticket
	challengeTicket := filcrypto.SHA256(partialTicket)

	networkPower := inds.TotalNetworkEffectivePower()

	sectorsSampled := uint64(math.Ceil(float64(node_base.EPOST_SAMPLE_RATE_NUM/node_base.EPOST_SAMPLE_RATE_DENOM) * float64(numSectors)))

	return spc.ec().IsWinningChallengeTicket(challengeTicket, sectorUtilization, networkPower, sectorsSampled, numSectors)
}

func (spc *StoragePowerConsensusSubsystem_I) _getStoragePowerActorState(stateTree stateTree.StateTree) spowact.StoragePowerActorState {
	powerAddr := addr.StoragePowerActorAddr
	actorState, ok := stateTree.GetActor(powerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, ok := spc.node().Repository().StateStore().Get(ipld.CID(substateCID))
	util.Assert(ok)

	// fix conversion to bytes
	util.IMPL_FINISH(substate)
	var serializedSubstate util.Serialization
	st, err := spowact.Deserialize_StoragePowerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (spc *StoragePowerConsensusSubsystem_I) GetTicketProductionRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_TICKET)
}

func (spc *StoragePowerConsensusSubsystem_I) GetSealRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_SEAL)
}

func (spc *StoragePowerConsensusSubsystem_I) GetPoStChallengeRand(chain chain.Chain, epoch block.ChainEpoch) util.Randomness {
	return chain.RandomnessAtEpoch(epoch - node_base.SPC_LOOKBACK_POST)
}

func (spc *StoragePowerConsensusSubsystem_I) GetFinalizedEpoch(currentEpoch block.ChainEpoch) block.ChainEpoch {
	return currentEpoch - node_base.FINALITY
}
