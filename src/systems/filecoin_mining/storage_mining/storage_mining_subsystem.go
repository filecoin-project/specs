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
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

type Serialization = util.Serialization

func (sms *StorageMiningSubsystem_I) CreateMiner(
	ownerAddr addr.Address,
	workerAddr addr.Address,
	sectorSize util.UInt,
	peerId libp2p.PeerID,
) addr.Address {
	// ownerAddr := sms.generateOwnerAddress(workerPubKey)
	// var pledgeAmt actor.TokenAmount TODO: unclear how to pass the amount/pay
	// TODO compute PledgeCollateral for 0 bytes
	// return sms.StoragePowerActor().CreateStorageMiner(ownerAddr, workerPubKey, sectorSize, peerId)
	// TODO: access this from runtime
	// return sms.StoragePowerActor().CreateStorageMiner(ownerAddr, workerAddr, peerId)
	var minerAddr addr.Address
	return minerAddr
}

func (sms *StorageMiningSubsystem_I) HandleStorageDeal(deal deal.StorageDeal) {
	sms.SectorIndex().AddNewDeal(deal)
	// stagedDealResponse := sms.SectorIndex().AddNewDeal(deal)
	// TODO: way within a node to notify different components
	// markeet.StorageProvider().NotifyStorageDealStaged(&storage_provider.StorageDealStagedNotification_I{
	// 	Deal_:     deal,
	// 	SectorID_: stagedDealResponse.SectorID(),
	// })
}

func (sms *StorageMiningSubsystem_I) _generateOwnerAddress(workerPubKey filcrypto.PublicKey) addr.Address {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() deal.StorageDeal {
	panic("TODO")
}

// triggered by new block reception and tipset assembly
func (sms *StorageMiningSubsystem_I) OnNewBestChain() {
	sms._tryLeaderElection()
}

// triggered by wall clock
func (sms *StorageMiningSubsystem_I) OnNewRound() {
	sms._tryLeaderElection()
}

func (sms *StorageMiningSubsystem_I) _tryLeaderElection() {

	// Randomness for ElectionPoSt
	randomnessK := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), sms._blockchain().LatestEpoch())

	input := sms._preparePoStChallengeSeed(randomnessK, sms._keyStore().MinerAddress())
	postRandomness := sms._keyStore().WorkerKey().Impl().Generate(input).Output()

	// TODO: add how sectors are actually stored in the SMS proving set
	util.TODO()
	provingSet := make([]sector.SectorID, 0)

	candidates := sms.StorageProving().Impl().GenerateElectionPoStCandidates(postRandomness, provingSet)

	if len(candidates) <= 0 {
		return // fail to generate post candidates
	}

	// TODO Fix
	util.TODO()
	var currState stateTree.StateTree
	winningCandidates := make([]sector.PoStCandidate, 0)
	st := sms._getStorageMinerActorState(currState, sms._keyStore().MinerAddress())

	for _, candidate := range candidates {
		sectorNum := candidate.SectorID().Number()
		utilInfo, err := st._getUtilizationInfo(sectorNum)
		if err != nil {
			// panic(err)
			return
		}
		sectorPower := utilInfo.CurrUtilization()
		if sms._consensus().IsWinningPartialTicket(currState, candidate.PartialTicket(), sectorPower) {
			winningCandidates = append(winningCandidates, candidate)
		}
	}

	if len(winningCandidates) <= 0 {
		return
	}

	// Randomness for ticket generation in block production
	randomness1 := sms._consensus().GetTicketProductionRand(sms._blockchain().BestChain(), sms._blockchain().LatestEpoch())
	newTicket := sms.PrepareNewTicket(randomness1, sms._keyStore().MinerAddress())

	postProof := sms.StorageProving().Impl().CreateElectionPoStProof(postRandomness, winningCandidates)
	chainHead := sms._blockchain().BestChain().HeadTipset()

	sms._blockProducer().GenerateBlock(postProof, winningCandidates, newTicket, chainHead, sms._keyStore().MinerAddress())
}

func (sms *StorageMiningSubsystem_I) _preparePoStChallengeSeed(randomness util.Randomness, minerAddr addr.Address) util.Randomness {

	randInput := Serialize_PoStChallengeSeedInput(&PoStChallengeSeedInput_I{
		ticket_:    randomness,
		minerAddr_: minerAddr,
	})
	input := filcrypto.DomainSeparationTag_PoSt.DeriveRand(randInput)
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
	input := filcrypto.DomainSeparationTag_TicketProduction.DeriveRand(randInput)

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
	actorState := stateTree.GetActorState(minerAddr)
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
	actorState := stateTree.GetActorState(powerAddr)
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

func (sms *StorageMiningSubsystem_I) GetWorkerKeyByMinerAddress(minerAddr addr.Address) filcrypto.VRFPublicKey {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) VerifyElectionPoSt(header block.BlockHeader, onChainInfo sector.OnChainPoStVerifyInfo) bool {

	sma := sms._getStorageMinerActorState(header.ParentState(), header.Miner())
	spa := sms._getStoragePowerActorState(header.ParentState())

	// 1. Check that the miner in question is currently allowed to run election
	// Note that this is two checks, namely:
	// On SMA --> can the miner be elected per electionPoSt rules?
	// On SPA --> Does the miner's power meet the consensus minimum requirement?
	// we could bundle into a single call here for convenience
	if !sma._canBeElected(header.Epoch()) {
		return false
	}

	pow, err := sma._getActivePower()
	if err != nil {
		// TODO: better error handling
		return false
	}

	if !spa.ActivePowerMeetsConsensusMinimum(pow) {
		return false
	}

	// 2. Verify appropriate randomness
	// TODO: fix away from BestChain()... every block should track its own chain up to its own production.
	randomness := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), header.Epoch())
	postRandomnessInput := sector.PoStRandomness(sms._preparePoStChallengeSeed(randomness, header.Miner()))

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	if !postRand.Verify(postRandomnessInput, sms.GetWorkerKeyByMinerAddress(header.Miner())) {
		return false
	}

	// A proof must be a valid snark proof with the correct public inputs
	// 3. Get public inputs
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

func (sms *StorageMiningSubsystem_I) VerifySurprisePoSt(header block.BlockHeader, onChainInfo sector.OnChainPoStVerifyInfo, posterAddr addr.Address) bool {

	st := sms._getStorageMinerActorState(header.ParentState(), header.Miner())

	// 1. Check that the miner in question is currently being challenged
	if !st._isChallenged() {
		// TODO: determine proper error here and error-handling machinery
		// rt.Abort("cannot SubmitSurprisePoSt when not challenged")
		return false
	}

	// 2. Check that the challenge has not expired
	// Check that miner can still submit (i.e. that the challenge window has not passed)
	// This will prevent miner from submitting a Surprise PoSt past the challenge period
	if st._challengeHasExpired(header.Epoch()) {
		return false
	}

	// A proof must be a valid snark proof with the correct public inputs

	// 3. Verify appropriate randomness
	randomnessEpoch := st.ChallengeStatus().LastChallengeEpoch()
	// TODO: fix away from BestChain()... every block should track its own chain up to its own production.
	randomness := sms._consensus().GetPoStChallengeRand(sms._blockchain().BestChain(), randomnessEpoch)
	postRandomnessInput := sms._preparePoStChallengeSeed(randomness, posterAddr)

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	if !postRand.Verify(postRandomnessInput, sms.GetWorkerKeyByMinerAddress(posterAddr)) {
		return false
	}

	// 4. Get public inputs
	info := st.Info()
	sectorSize := info.SectorSize()

	postCfg := sector.PoStCfg_I{
		Type_:        sector.PoStType_SurprisePoSt,
		SectorSize_:  sectorSize,
		WindowCount_: info.WindowCount(),
		Partitions_:  info.SurprisePoStPartitions(),
	}

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		PoStCfg_:    &postCfg,
		Randomness_: onChainInfo.Randomness(),
	}

	sdr := filproofs.WinSDRParams(&filproofs.SDRCfg_I{SurprisePoStCfg_: &postCfg})

	// 5. Verify the PoSt Proof
	isPoStVerified := sdr.VerifySurprisePoSt(&pvInfo)
	return isPoStVerified
}

func (sms *StorageMiningSubsystem_I) VerifyElection(header block.BlockHeader, onChainInfo sector.OnChainPoStVerifyInfo) bool {
	st := sms._getStorageMinerActorState(header.ParentState(), header.Miner())

	for _, info := range onChainInfo.Candidates() {
		sectorNum := info.SectorID().Number()
		utilInfo, err := st._getUtilizationInfo(sectorNum)
		if err != nil {
			// panic(err)
			return false
		}
		sectorPower := utilInfo.CurrUtilization()
		if !sms._consensus().IsWinningPartialTicket(header.ParentState(), info.PartialTicket(), sectorPower) {
			return false
		}
	}
	return true
}

// func (sms *StorageMiningSubsystem_I) submitPoStMessage(postSubmission poster.PoStSubmission) error {
// 	var workerAddress addr.Address
// 	var workerKeyPair filcrypto.SigKeyPair
// 	panic("TODO") // TODO: get worker address and key pair

// 	// TODO: is this just workerAddress, or is there a separation here
// 	// (worker is AccountActor, workerMiner is StorageMinerActor)?
// 	var workerMinerActorAddress addr.Address
// 	panic("TODO")

// 	var gasPrice msg.GasPrice
// 	var gasLimit msg.GasAmount
// 	panic("TODO") // TODO: determine gas price and limit

// 	var callSeqNum actor.CallSeqNum
// 	panic("TODO") // TODO: retrieve CallSeqNum from worker

// 	messageParams := actor.MethodParams([]actor.MethodParam{
// 		actor.MethodParam(poster.Serialize_PoStSubmission(postSubmission)),
// 	})

// 	unsignedMessage := msg.UnsignedMessage_Make(
// 		workerAddress,
// 		workerMinerActorAddress,
// 		Method_StorageMinerActor_SubmitPoSt,
// 		messageParams,
// 		callSeqNum,
// 		actor.TokenAmount(0),
// 		gasPrice,
// 		gasLimit,
// 	)

// 	signedMessage, err := msg.Sign(unsignedMessage, workerKeyPair)
// 	if err != nil {
// 		return err
// 	}

// 	err = sms.FilecoinNode().SubmitMessage(signedMessage)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
