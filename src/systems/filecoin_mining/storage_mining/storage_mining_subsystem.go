package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	indices "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

type Serialization = util.Serialization

// Note that implementations may choose to provide default generation methods for miners created
// without miner/owner keypairs. We omit these details from the spec.
// Also note that the pledge amount should be available in the ownerAddr in order for this call
// to succeed.
func (sms *StorageMiningSubsystem_I) CreateMiner(
	state stateTree.StateTree,
	ownerAddr addr.Address,
	workerAddr addr.Address,
	sectorSize util.UInt,
	peerId libp2p.PeerID,
	pledgeAmt actor.TokenAmount,
) (addr.Address, error) {

	params := make([]util.Serialization, 4)
	params[0] = addr.Serialize_Address(ownerAddr)
	params[1] = addr.Serialize_Address(workerAddr)
	params[2] = libp2p.Serialize_PeerID(peerId)

	ownerActor, ok := state.GetActor(ownerAddr)
	Assert(ok)

	unsignedCreationMessage := &msg.UnsignedMessage_I{
		From_:       ownerAddr,
		To_:         addr.StoragePowerActorAddr,
		Method_:     ai.Method_StoragePowerActor_CreateMiner,
		Params_:     params,
		CallSeqNum_: ownerActor.CallSeqNum(),
		Value_:      pledgeAmt,
		GasPrice_:   0,
		GasLimit_:   msg.GasAmount_SentinelUnlimited(),
	}

	var workerKey filcrypto.SigKeyPair // sms._keyStore().WorkerKey()
	signedMessage, err := msg.Sign(unsignedCreationMessage, workerKey)
	if err != nil {
		return nil, err
	}

	err = sms.FilecoinNode().MessagePool().Syncer().SubmitMessage(signedMessage)
	if err != nil {
		return nil, err
	}

	// WAIT for block reception with appropriate response from SPA
	util.IMPL_TODO()

	// harvest address from that block
	var storageMinerAddr addr.Address
	// and set in key store appropriately
	return storageMinerAddr, nil
}

func (sms *StorageMiningSubsystem_I) HandleStorageDeal(deal deal.StorageDeal) {
	sms.SectorIndex().AddNewDeal(deal)
	// stagedDealResponse := sms.SectorIndex().AddNewDeal(deal)
	// TODO: way within a node to notify different components
	// market.StorageProvider().NotifyStorageDealStaged(&storage_provider.StorageDealStagedNotification_I{
	// 	Deal_:     deal,
	// 	SectorID_: stagedDealResponse.SectorID(),
	// })
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() deal.StorageDeal {
	panic("TODO")
}

// triggered by new block reception and tipset assembly
func (sms *StorageMiningSubsystem_I) OnNewBestChain() {
	sms._runMiningCycle()
}

// triggered by wall clock
func (sms *StorageMiningSubsystem_I) OnNewRound() {
	sms._runMiningCycle()
}

func (sms *StorageMiningSubsystem_I) _runMiningCycle() {
	chainHead := sms._blockchain().BestChain().HeadTipset()
	sma := sms._getStorageMinerActorState(chainHead.StateTree(), sms._keyStore().MinerAddress())

	if sma.PoStState().Is_OK() {
		ePoSt := sms._tryLeaderElection(chainHead.StateTree(), sma)
		if ePoSt != nil {
			// Randomness for ticket generation in block production
			randomness1 := sms._consensus().GetTicketProductionRand(sms._blockchain().BestChain(), sms._blockchain().LatestEpoch())
			newTicket := sms.PrepareNewTicket(randomness1, sms._keyStore().MinerAddress())

			sms._blockProducer().GenerateBlock(ePoSt, newTicket, chainHead, sms._keyStore().MinerAddress())
		}
	} else if sma.PoStState().Is_Challenged() {
		sPoSt := sms._trySurprisePoSt(chainHead.StateTree(), sma)

		var gasLimit msg.GasAmount
		var gasPrice = actor.TokenAmount(0)
		IMPL_TODO("read from consts (in this case user set param)")

		sms._submitSurprisePoStMessage(chainHead.StateTree(), sPoSt, gasPrice, gasLimit)
	}
}

func (sms *StorageMiningSubsystem_I) _tryLeaderElection(currState stateTree.StateTree, sma StorageMinerActorState) sector.OnChainPoStVerifyInfo {
	// Randomness for ElectionPoSt
	randomnessK := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), sms._blockchain().LatestEpoch())

	input := sms.PreparePoStChallengeSeed(randomnessK, sms._keyStore().MinerAddress())
	postRandomness := sms._keyStore().WorkerKey().Impl().Generate(input).Output()

	// TODO: add how sectors are actually stored in the SMS proving set
	util.TODO()
	provingSet := make([]sector.SectorID, 0)

	candidates := sms.StorageProving().Impl().GenerateElectionPoStCandidates(postRandomness, provingSet)

	if len(candidates) <= 0 {
		return nil // fail to generate post candidates
	}

	winningCandidates := make([]sector.PoStCandidate, 0)

	var numMinerSectors uint64
	TODO() // update
	// numMinerSectors := uint64(len(sma.SectorTable().Impl().ActiveSectors_.SectorsOn()))

	for _, candidate := range candidates {
		sectorNum := candidate.SectorID().Number()
		sectorWeightDesc, ok := sma._getStorageWeightDescForSectorMaybe(sectorNum)
		if !ok {
			return nil
		}
		sectorPower := indices.ConsensusPowerForStorageWeight(sectorWeightDesc)
		if sms._consensus().IsWinningPartialTicket(currState, candidate.PartialTicket(), sectorPower, numMinerSectors) {
			winningCandidates = append(winningCandidates, candidate)
		}
	}

	if len(winningCandidates) <= 0 {
		return nil
	}

	postProof := sms.StorageProving().Impl().CreateElectionPoStProof(postRandomness, winningCandidates)

	var ctc sector.ChallengeTicketsCommitment // TODO: proofs to fix when complete
	electionPoSt := &sector.OnChainPoStVerifyInfo_I{
		CommT_:      ctc,
		Candidates_: winningCandidates,
		Randomness_: postRandomness,
		Proof_:      postProof,
	}

	return electionPoSt
}

func (sms *StorageMiningSubsystem_I) PreparePoStChallengeSeed(randomness util.Randomness, minerAddr addr.Address) util.Randomness {

	randInput := Serialize_PoStChallengeSeedInput(&PoStChallengeSeedInput_I{
		ticket_:    randomness,
		minerAddr_: minerAddr,
	})
	input := filcrypto.DeriveRand(filcrypto.DomainSeparationTag_PreparePoStChallengeSeed, randInput)
	return input
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(randomness util.Randomness, minerActorAddr addr.Address) block.Ticket {
	// run it through the VRF and get deterministic output

	// take the VRFResult of that ticket as input, specifying the personalization (see data structures)
	// append the miner actor address for the miner generifying this in order to prevent miners with the same
	// worker keys from generating the same randomness (given the VRF)
	randInput := block.Serialize_TicketProductionSeedInput(&block.TicketProductionSeedInput_I{
		PastTicket_: randomness,
		MinerAddr_:  minerActorAddr,
	})
	input := filcrypto.DeriveRand(filcrypto.DomainSeparationTag_TicketProduction, randInput)

	// run through VRF
	vrfRes := sms._keyStore().WorkerKey().Impl().Generate(input)

	newTicket := &block.Ticket_I{
		VRFResult_: vrfRes,
		Output_:    vrfRes.Output(),
	}

	return newTicket
}

// TODO: fix linking here
var node node_base.FilecoinNode

func (sms *StorageMiningSubsystem_I) _getStorageMinerActorState(stateTree stateTree.StateTree, minerAddr addr.Address) StorageMinerActorState {
	actorState, ok := stateTree.GetActor(minerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, err := node.LocalGraph().Get(ipld.CID(substateCID))
	if err != nil {
		panic("TODO")
	}
	// TODO fix conversion to bytes
	panic(substate)
	var serializedSubstate Serialization
	st, err := Deserialize_StorageMinerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (sms *StorageMiningSubsystem_I) _getStoragePowerActorState(stateTree stateTree.StateTree) spc.StoragePowerActorState {
	powerAddr := addr.StoragePowerActorAddr
	actorState, ok := stateTree.GetActor(powerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, err := node.LocalGraph().Get(ipld.CID(substateCID))
	if err != nil {
		panic("TODO")
	}

	// TODO fix conversion to bytes
	panic(substate)
	var serializedSubstate util.Serialization
	st, err := spc.Deserialize_StoragePowerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (sms *StorageMiningSubsystem_I) VerifyElectionPoSt(inds indices.Indices, header block.BlockHeader, onChainInfo sector.OnChainPoStVerifyInfo) bool {
	sma := sms._getStorageMinerActorState(header.ParentState(), header.Miner())
	spa := sms._getStoragePowerActorState(header.ParentState())

	pow, found := spa.PowerTable()[header.Miner()]
	if !found {
		return false
	}

	// 1. Verify miner has enough power (includes implicit checks on min miner size
	// and challenge status via SPA's power table).
	if pow == block.StoragePower(0) {
		return false
	}

	// 2. Verify partialTicket values are appropriate
	if !sms._verifyElection(header, onChainInfo) {
		return false
	}

	// verify the partialTickets themselves
	// 3. Verify appropriate randomness
	// TODO: fix away from BestChain()... every block should track its own chain up to its own production.
	randomness := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), header.Epoch())
	postRandomnessInput := sector.PoStRandomness(sms.PreparePoStChallengeSeed(randomness, header.Miner()))

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	// get worker key from minerAddr
	workerKey, err := sma.Info().Worker().GetKey()
	if err != nil {
		return false
	}

	if !postRand.Verify(postRandomnessInput, filcrypto.VRFPublicKey(workerKey)) {
		return false
	}

	// A proof must be a valid snark proof with the correct public inputs
	// 4. Get public inputs
	info := sma.Info()
	sectorSize := info.SectorSize()

	postCfg := sector.PoStCfg_I{
		Type_:        sector.PoStType_ElectionPoSt,
		SectorSize_:  sectorSize,
		WindowCount_: info.WindowCount(),
		Partitions_:  info.ElectionPoStPartitions(),
	}

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		PoStCfg_:    &postCfg,
		Randomness_: onChainInfo.Randomness(),
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{ElectionPoStCfg_: &postCfg})

	// 5. Verify the PoSt Proof
	isPoStVerified := sdr.VerifyElectionPoSt(&pvInfo)
	return isPoStVerified
}

func (sms *StorageMiningSubsystem_I) _verifyElection(header block.BlockHeader, onChainInfo sector.OnChainPoStVerifyInfo) bool {
	st := sms._getStorageMinerActorState(header.ParentState(), header.Miner())

	var numMinerSectors uint64
	TODO()
	// TODO: Decide whether to sample sectors uniformly for EPoSt (the cleanest),
	// or to sample weighted by nominal power.

	for _, info := range onChainInfo.Candidates() {
		sectorNum := info.SectorID().Number()
		sectorWeightDesc, ok := st._getStorageWeightDescForSectorMaybe(sectorNum)
		if !ok {
			return false
		}
		sectorPower := indices.ConsensusPowerForStorageWeight(sectorWeightDesc)
		if !sms._consensus().IsWinningPartialTicket(header.ParentState(), info.PartialTicket(), sectorPower, numMinerSectors) {
			return false
		}
	}
	return true
}

func (sms *StorageMiningSubsystem_I) _trySurprisePoSt(currState stateTree.StateTree, sma StorageMinerActorState) sector.OnChainPoStVerifyInfo {
	if !sma.PoStState().Is_Challenged() {
		return nil
	}

	// get randomness for SurprisePoSt
	challEpoch := sma.PoStState().As_Challenged().SurpriseChallengeEpoch()
	randomnessK := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), challEpoch)
	input := sms.PreparePoStChallengeSeed(randomnessK, sms._keyStore().MinerAddress())
	postRandomness := sms._keyStore().WorkerKey().Impl().Generate(input).Output()

	// TODO: add how sectors are actually stored in the SMS proving set
	util.TODO()
	provingSet := make([]sector.SectorID, 0)

	candidates := sms.StorageProving().Impl().GenerateSurprisePoStCandidates(postRandomness, provingSet)

	if len(candidates) <= 0 {
		// Error. Will fail this surprise post and must then redeclare faults
		return nil // fail to generate post candidates
	}

	winningCandidates := make([]sector.PoStCandidate, 0)
	for _, candidate := range candidates {
		if sma._verifySurprisePoStMeetsTargetReq(candidate) {
			winningCandidates = append(winningCandidates, candidate)
		}
	}

	postProof := sms.StorageProving().Impl().CreateSurprisePoStProof(postRandomness, winningCandidates)

	var ctc sector.ChallengeTicketsCommitment // TODO: proofs to fix when complete
	surprisePoSt := &sector.OnChainPoStVerifyInfo_I{
		CommT_:      ctc,
		Candidates_: winningCandidates,
		Randomness_: postRandomness,
		Proof_:      postProof,
	}
	return surprisePoSt
}

func (sms *StorageMiningSubsystem_I) _submitSurprisePoStMessage(state stateTree.StateTree, sPoSt sector.OnChainPoStVerifyInfo, gasPrice actor.TokenAmount, gasLimit msg.GasAmount) error {
	params := make([]util.Serialization, 1)
	params[0] = sector.Serialize_OnChainPoStVerifyInfo(sPoSt)

	workerAddr, err := addr.Address_Make_Key(node_base.NETWORK, addr.KeyHash(sms._keyStore().WorkerKey().VRFPublicKey()))
	if err != nil {
		return err
	}
	worker, ok := state.GetActor(workerAddr)
	Assert(ok)
	unsignedCreationMessage := &msg.UnsignedMessage_I{
		From_:       workerAddr,
		To_:         sms._keyStore().MinerAddress(),
		Method_:     ai.Method_StorageMinerActor_SubmitSurprisePoStResponse,
		Params_:     params,
		CallSeqNum_: worker.CallSeqNum(),
		Value_:      actor.TokenAmount(0),
		GasPrice_:   gasPrice,
		GasLimit_:   gasLimit,
	}

	var workerKey filcrypto.SigKeyPair // sms._keyStore().WorkerKey()
	signedMessage, err := msg.Sign(unsignedCreationMessage, workerKey)
	if err != nil {
		return err
	}

	err = sms.FilecoinNode().MessagePool().Syncer().SubmitMessage(signedMessage)
	if err != nil {
		return err
	}

	return nil
}
