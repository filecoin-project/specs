package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	sminact "github.com/filecoin-project/specs/actors/builtin/storage_miner"
	spowact "github.com/filecoin-project/specs/actors/builtin/storage_power"
	indices "github.com/filecoin-project/specs/actors/runtime/indices"
	serde "github.com/filecoin-project/specs/actors/serde"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	stateTree "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

type Serialization = util.Serialization

var Assert = util.Assert
var TODO = util.TODO

// Note that implementations may choose to provide default generation methods for miners created
// without miner/owner keypairs. We omit these details from the spec.
// Also note that the pledge amount should be available in the ownerAddr in order for this call
// to succeed.
func (sms *StorageMiningSubsystem_I) CreateMiner(
	state stateTree.StateTree,
	ownerAddr addr.Address,
	workerAddr addr.Address,
	sectorSize util.UInt,
	peerId peer.ID,
	pledgeAmt abi.TokenAmount,
) (addr.Address, error) {

	ownerActor, ok := state.GetActor(ownerAddr)
	Assert(ok)

	unsignedCreationMessage := &msg.UnsignedMessage_I{
		From_:       ownerAddr,
		To_:         builtin.StoragePowerActorAddr,
		Method_:     ai.Method_StoragePowerActor_CreateMiner,
		Params_:     serde.MustSerializeParams(ownerAddr, workerAddr, peerId),
		CallSeqNum_: ownerActor.CallSeqNum(),
		Value_:      pledgeAmt,
		GasPrice_:   0,
		GasLimit_:   msg.GasAmount_SentinelUnlimited(),
	}

	var workerKey filcrypto.SigKeyPair // sms._keyStore().Worker()
	signedMessage, err := msg.Sign(unsignedCreationMessage, workerKey)
	if err != nil {
		return addr.Undef, err
	}

	err = sms.Node().MessagePool().Syncer().SubmitMessage(signedMessage)
	if err != nil {
		return addr.Undef, err
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
	sma := sms._getStorageMinerActorState(chainHead.StateTree(), sms.Node().Repository().KeyStore().MinerAddress())

	if sma.PoStState().Is_OK() {
		ePoSt := sms._tryLeaderElection(chainHead.StateTree(), sma)
		if ePoSt != nil {
			// Randomness for ticket generation in block production
			randomness1 := sms._blockchain().BestChain().GetTicketProductionRandSeed(sms._blockchain().LatestEpoch())
			newTicket := sms.PrepareNewTicket(randomness1, sms.Node().Repository().KeyStore().MinerAddress())

			sms._blockProducer().GenerateBlock(ePoSt, newTicket, chainHead, sms.Node().Repository().KeyStore().MinerAddress())
		}
	} else if sma.PoStState().Is_Challenged() {
		sPoSt := sms._trySurprisePoSt(chainHead.StateTree(), sma)

		var gasLimit msg.GasAmount
		var gasPrice = abi.TokenAmount(0)
		util.IMPL_FINISH("read from consts (in this case user set param)")
		sms._submitSurprisePoStMessage(chainHead.StateTree(), sPoSt, gasPrice, gasLimit)
	}
}

func (sms *StorageMiningSubsystem_I) _tryLeaderElection(currState stateTree.StateTree, sma sminact.StorageMinerActorState) sector.OnChainElectionPoStVerifyInfo {
	// Randomness for ElectionPoSt

	randomnessK := sms._blockchain().BestChain().GetPoStChallengeRandSeed(sms._blockchain().LatestEpoch())
	input := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_ElectionPoStChallengeSeed, randomnessK, sms.Node().Repository().KeyStore().MinerAddress())
	// Use VRF to generate secret randomness
	postRandomness := sms.Node().Repository().KeyStore().WorkerKey().Impl().Generate(input).Output()

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
		sectorWeightDesc, ok := sma.GetStorageWeightDescForSectorMaybe(sectorNum)
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

	postProofs := sms.StorageProving().Impl().CreateElectionPoStProof(postRandomness, winningCandidates)

	electionPoSt := &sector.OnChainElectionPoStVerifyInfo_I{
		Candidates_: winningCandidates,
		Randomness_: postRandomness,
		Proofs_:     postProofs,
	}

	return electionPoSt
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(randomness abi.RandomnessSeed, minerActorAddr addr.Address) block.Ticket {
	// run it through the VRF and get deterministic output

	// take the VRFResult of that ticket as input, specifying the personalization (see data structures)
	// append the miner actor address for the miner generifying this in order to prevent miners with the same
	// worker keys from generating the same randomness (given the VRF)
	input := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_TicketProduction, randomness, minerActorAddr)

	// run through VRF
	vrfRes := sms.Node().Repository().KeyStore().WorkerKey().Impl().Generate(input)

	newTicket := &block.Ticket_I{
		VRFResult_: vrfRes,
		Output_:    vrfRes.Output(),
	}

	return newTicket
}

func (sms *StorageMiningSubsystem_I) _getStorageMinerActorState(stateTree stateTree.StateTree, minerAddr addr.Address) sminact.StorageMinerActorState {
	actorState, ok := stateTree.GetActor(minerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, ok := sms.Node().Repository().StateStore().Get(cid.Cid(substateCID))
	if !ok {
		panic("Couldn't find sma state")
	}
	// fix conversion to bytes
	util.IMPL_TODO(substate)
	var serializedSubstate Serialization
	st, err := sminact.Deserialize_StorageMinerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (sms *StorageMiningSubsystem_I) _getStoragePowerActorState(stateTree stateTree.StateTree) spowact.StoragePowerActorState {
	powerAddr := builtin.StoragePowerActorAddr
	actorState, ok := stateTree.GetActor(powerAddr)
	util.Assert(ok)
	substateCID := actorState.State()

	substate, ok := sms.Node().Repository().StateStore().Get(cid.Cid(substateCID))
	if !ok {
		panic("Couldn't find spa state")
	}

	// fix conversion to bytes
	util.IMPL_TODO(substate)
	var serializedSubstate util.Serialization
	st, err := spowact.Deserialize_StoragePowerActorState(serializedSubstate)

	if err == nil {
		panic("Deserialization error")
	}
	return st
}

func (sms *StorageMiningSubsystem_I) VerifyElectionPoSt(inds indices.Indices, header block.BlockHeader, onChainInfo sector.OnChainElectionPoStVerifyInfo) bool {
	sma := sms._getStorageMinerActorState(header.ParentState(), header.Miner())
	spa := sms._getStoragePowerActorState(header.ParentState())

	pow, found := spa.PowerTable()[header.Miner()]
	if !found {
		return false
	}

	// 1. Verify miner has enough power (includes implicit checks on min miner size
	// and challenge status via SPA's power table).
	if pow == abi.StoragePower(0) {
		return false
	}

	// 2. verify no duplicate tickets included
	challengeIndices := make(map[util.UInt]bool)
	for _, tix := range onChainInfo.Candidates() {
		if _, ok := challengeIndices[tix.ChallengeIndex()]; ok {
			return false
		}
		challengeIndices[tix.ChallengeIndex()] = true
	}

	// 3. Verify partialTicket values are appropriate
	if !sms._verifyElection(header, onChainInfo) {
		return false
	}

	// verify the partialTickets themselves
	// 4. Verify appropriate randomness
	// TODO: fix away from BestChain()... every block should track its own chain up to its own production.
	randomness := sms._blockchain().BestChain().GetPoStChallengeRandSeed(header.Epoch())
	input := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_ElectionPoStChallengeSeed, randomness, header.Miner())

	postRand := &filcrypto.VRFResult_I{
		Output_: onChainInfo.Randomness(),
	}

	// TODO if the workerAddress is secp then payload will be the blake2b hash of its public key
	// and we will need to recover the entire public key from worker before handing off to Verify
	// example of recover code: https://github.com/ipsn/go-secp256k1/blob/master/secp256.go#L93
	workerKey := sma.Info().Worker().Payload()
	// Verify VRF output from appropriate input corresponds to randomness used
	if !postRand.Verify(input, filcrypto.VRFPublicKey(workerKey)) {
		return false
	}

	// A proof must be a valid snark proof with the correct public inputs
	// 5. Get public inputs
	info := sma.Info()
	sectorSize := info.SectorSize()

	postCfg := filproofs.ElectionPoStCfg(sectorSize)

	pvInfo := sector.PoStVerifyInfo_I{
		OnChain_:    onChainInfo,
		Randomness_: onChainInfo.Randomness(),
	}

	pv := filproofs.MakeElectionPoStVerifier(postCfg)

	// 5. Verify the PoSt Proof
	isPoStVerified := pv.VerifyElectionPoSt(&pvInfo)

	return isPoStVerified
}

func (sms *StorageMiningSubsystem_I) _verifyElection(header block.BlockHeader, onChainInfo sector.OnChainElectionPoStVerifyInfo) bool {
	st := sms._getStorageMinerActorState(header.ParentState(), header.Miner())

	var numMinerSectors uint64
	TODO()
	// TODO: Decide whether to sample sectors uniformly for EPoSt (the cleanest),
	// or to sample weighted by nominal power.

	for _, info := range onChainInfo.Candidates() {
		sectorNum := info.SectorID().Number()
		sectorWeightDesc, ok := st.GetStorageWeightDescForSectorMaybe(sectorNum)
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

func (sms *StorageMiningSubsystem_I) _trySurprisePoSt(currState stateTree.StateTree, sma sminact.StorageMinerActorState) sector.OnChainSurprisePoStVerifyInfo {
	if !sma.PoStState().Is_Challenged() {
		return nil
	}

	// get randomness for SurprisePoSt
	challEpoch := sma.PoStState().As_Challenged().SurpriseChallengeEpoch()
	randomnessK := sms._blockchain().BestChain().GetPoStChallengeRandSeed(challEpoch)
	// unlike with ElectionPoSt no need to use a VRF
	postRandomness := filcrypto.DeriveRandWithMinerAddr(filcrypto.DomainSeparationTag_SurprisePoStChallengeSeed, randomnessK, sms.Node().Repository().KeyStore().MinerAddress())

	// TODO: add how sectors are actually stored in the SMS proving set
	util.TODO()
	provingSet := make([]sector.SectorID, 0)

	candidates := sms.StorageProving().Impl().GenerateSurprisePoStCandidates(sector.PoStRandomness(postRandomness), provingSet)

	if len(candidates) <= 0 {
		// Error. Will fail this surprise post and must then redeclare faults
		return nil // fail to generate post candidates
	}

	winningCandidates := make([]sector.PoStCandidate, 0)
	for _, candidate := range candidates {
		if sma.VerifySurprisePoStMeetsTargetReq(candidate) {
			winningCandidates = append(winningCandidates, candidate)
		}
	}

	postProofs := sms.StorageProving().Impl().CreateSurprisePoStProof(sector.PoStRandomness(postRandomness), winningCandidates)

	// var ctc sector.ChallengeTicketsCommitment // TODO: proofs to fix when complete
	surprisePoSt := &sector.OnChainSurprisePoStVerifyInfo_I{
		// CommT_:      ctc,
		Candidates_: winningCandidates,
		Proofs_:     postProofs,
	}
	return surprisePoSt
}

func (sms *StorageMiningSubsystem_I) _submitSurprisePoStMessage(state stateTree.StateTree, sPoSt sector.OnChainSurprisePoStVerifyInfo, gasPrice abi.TokenAmount, gasLimit msg.GasAmount) error {
	// TODO if workerAddr is not a secp key (e.g. BLS) then this will need to be handled differently
	workerAddr, err := addr.NewSecp256k1Address(sms.Node().Repository().KeyStore().WorkerKey().VRFPublicKey())
	if err != nil {
		return err
	}
	worker, ok := state.GetActor(workerAddr)
	Assert(ok)
	unsignedCreationMessage := &msg.UnsignedMessage_I{
		From_:       sms.Node().Repository().KeyStore().MinerAddress(),
		To_:         sms.Node().Repository().KeyStore().MinerAddress(),
		Method_:     ai.Method_StorageMinerActor_SubmitSurprisePoStResponse,
		Params_:     serde.MustSerializeParams(sPoSt),
		CallSeqNum_: worker.CallSeqNum(),
		Value_:      abi.TokenAmount(0),
		GasPrice_:   gasPrice,
		GasLimit_:   gasLimit,
	}

	var workerKey filcrypto.SigKeyPair // sms.Node().Repository().KeyStore().Worker()
	signedMessage, err := msg.Sign(unsignedCreationMessage, workerKey)
	if err != nil {
		return err
	}

	err = sms.Node().MessagePool().Syncer().SubmitMessage(signedMessage)
	if err != nil {
		return err
	}

	return nil
}
