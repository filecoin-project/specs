package storage_power_consensus

import (
	"math"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	builtin "github.com/filecoin-project/specs-actors/actors/builtin"
	spowact "github.com/filecoin-project/specs-actors/actors/builtin/storage_power"
	acrypto "github.com/filecoin-project/specs-actors/actors/crypto"
	inds "github.com/filecoin-project/specs-actors/actors/runtime/indices"
	serde "github.com/filecoin-project/specs-actors/actors/serde"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
)

// Storage Power Consensus Subsystem

func (spc *StoragePowerConsensusSubsystem_I) ValidateBlock(block block.Block_I) error {
	util.IMPL_FINISH()
	return nil
}

func (spc *StoragePowerConsensusSubsystem_I) validateTicket(ticket block.Ticket, pk filcrypto.VRFPublicKey, minerActorAddr addr.Address) bool {
	randomness1 := spc.blockchain().BestChain().GetTicketProductionRandSeed(spc.blockchain().LatestEpoch())

	return ticket.Verify(randomness1, pk, minerActorAddr)
}

func (spc *StoragePowerConsensusSubsystem_I) ComputeChainWeight(tipset chain.Tipset) block.ChainWeight {
	return spc.ec().ComputeChainWeight(tipset)
}

func (spc *StoragePowerConsensusSubsystem_I) IsWinningPartialTicket(stateTree stateTree.StateTree, inds inds.Indices, partialTicket abi.PartialTicket, sectorUtilization abi.StoragePower, numSectors util.UVarint) bool {

	// finalize the partial ticket
	challengeTicket := acrypto.SHA256(abi.Bytes(partialTicket))

	networkPower := inds.TotalNetworkEffectivePower()

	sectorsSampled := uint64(math.Ceil(float64(node_base.EPOST_SAMPLE_RATE_NUM/node_base.EPOST_SAMPLE_RATE_DENOM) * float64(numSectors)))

	return spc.ec().IsWinningChallengeTicket(challengeTicket, sectorUtilization, networkPower, sectorsSampled, numSectors)
}

func (spc *StoragePowerConsensusSubsystem_I) _getStoragePowerActorState(stateTree stateTree.StateTree) spowact.StoragePowerActorState {
	powerAddr := builtin.StoragePowerActorAddr
	actorState, ok := stateTree.GetActor(powerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, ok := spc.node().Repository().StateStore().Get(cid.Cid(substateCID))
	util.Assert(ok)

	// fix conversion to bytes
	util.IMPL_FINISH(substate)
	var serializedSubstate util.Serialization
	var st spowact.StoragePowerActorState
	serde.MustDeserialize(serializedSubstate, &st)
	return st
}

func (spc *StoragePowerConsensusSubsystem_I) GetFinalizedEpoch(currentEpoch abi.ChainEpoch) abi.ChainEpoch {
	return currentEpoch - node_base.FINALITY
}
