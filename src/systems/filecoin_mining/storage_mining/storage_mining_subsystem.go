package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
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

func (sms *StorageMiningSubsystem_I) generateOwnerAddress(workerPubKey filcrypto.PublicKey) addr.Address {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() deal.StorageDeal {
	panic("TODO")
}

// triggered by new block reception and tipset assembly
func (sms *StorageMiningSubsystem_I) OnNewBestChain() {
	sms.tryLeaderElection()
}

// triggered by wall clock
func (sms *StorageMiningSubsystem_I) OnNewRound() {
	sms.tryLeaderElection()
}

func (sms *StorageMiningSubsystem_I) tryLeaderElection() {

	// Draw randomness from chain for ElectionPoSt and Ticket Generation
	// Randomness for ticket generation in block production
	randomness1 := sms.consensus().GetTicketProductionSeed(sms.blockchain().BestChain(), sms.keyStore().OwnerAddress(), sms.blockchain().LatestEpoch())

	// Randomness for ElectionPoSt
	randomnessK := sms.consensus().GetPoStChallenge(sms.blockchain().BestChain(), sms.keyStore().OwnerAddress(), sms.blockchain().LatestEpoch())

	// TODO: @why @jz align on this
	for _, worker := range sms.keyStore().Workers() {

		var input []byte
		input = append(input, spc.VRFPersonalizationPoSt)
		input = append(input, randomnessK...)

		postRandomness := worker.VRFKeyPair().Impl().Generate(input).Output()
		// TODO: add how sectors are actually stored in the SMS proving set
		provingSet := make([]sector.SectorID, 0)

		challengeTickets := sms.StorageProving().Impl().GeneratePoStCandidates(postRandomness, provingSet)

		if len(challengeTickets) <= 0 {
			return // fail to generate post candidates
		}

		winningCTs := make([]sector.ChallengeTicket, 0)

		for _, ct := range challengeTickets {
			// TODO align on worker address
			if sms.consensus().IsWinningChallengeTicket(ct) {
				winningCTs = append(winningCTs, ct)
			}
		}

		if len(winningCTs) <= 0 {
			return
		}

		newTicket := sms.PrepareNewTicket(randomness1, worker.VRFKeyPair())
		postProof := sms.StorageProving().Impl().GeneratePoStProof(postRandomness, winningCTs)
		chainHead := sms.blockchain().BestChain().HeadTipset()

		sms.blockProducer().GenerateBlock(postProof, winningCTs, newTicket, chainHead, worker.Address())

	}
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(randomness util.Randomness, vrfKP filcrypto.VRFKeyPair) block.Ticket {
	// run it through the VRF and get deterministic output

	// take the VRFResult of that ticket as input, specifying the personalization (see data structures)
	// append the miner actor address for the miner generifying this in order to prevent miners with the same
	// worker keys from generating the same randomness (given the VRF)
	var input []byte
	input = append(input, byte(filcrypto.DomainSeparationTag_Case_Ticket))
	input = append(input, randomness...)
	input = append(input, byte(filcrypto.InputDelimeter_Case_Bytes))
	input = append(input, minerAddr.AddrToBytes()...)

	// run through VRF
	vrfRes := vrfKP.Generate(input)

	newTicket := &block.Ticket_I{
		VRFResult_: vrfRes,
		Output_:    vrfRes.Output(),
	}

	return newTicket
}

func (sms *StorageMiningSubsystem_I) DrawElectionProof(randomness block.Randomness, height block.ChainEpoch, vrfKP filcrypto.VRFKeyPair) block.ElectionProof {
	panic("")
	// // 0. Prepare new election proof
	// var newEP ElectionProof

	// // 1. Run it through VRF and get determinstic output
	// // 1.i. # take the VRFOutput of that ticket as input, specified for the appropriate operation type
	// var input []byte
	// input = append (input, filcrypto.ElectionTag)
	// input = append(input, lookbackTicket.Output)
	// input = append(input, height)
	// // ii. # run it through the VRF and store the VRFProof in the new ticket
	// newEP.VRFResult := vrfKP.Generate(input)
	// return newEP
}

func (sms *StorageMiningSubsystem_I) submitPoStMessage(postSubmission poster.PoStSubmission) error {
	var workerAddress addr.Address
	var workerKeyPair filcrypto.SigKeyPair
	panic("TODO") // TODO: get worker address and key pair

	// TODO: is this just workerAddress, or is there a separation here
	// (worker is AccountActor, workerMiner is StorageMinerActor)?
	var workerMinerActorAddress addr.Address
	panic("TODO")

	var gasPrice msg.GasPrice
	var gasLimit msg.GasAmount
	panic("TODO") // TODO: determine gas price and limit

	var callSeqNum actor.CallSeqNum
	panic("TODO") // TODO: retrieve CallSeqNum from worker

	messageParams := actor.MethodParams([]actor.MethodParam{
		actor.MethodParam(poster.Serialize_PoStSubmission(postSubmission)),
	})

	unsignedMessage := msg.UnsignedMessage_Make(
		workerAddress,
		workerMinerActorAddress,
		Method_StorageMinerActor_SubmitPoSt,
		messageParams,
		callSeqNum,
		actor.TokenAmount(0),
		gasPrice,
		gasLimit,
	)

	signedMessage, err := msg.Sign(unsignedMessage, workerKeyPair)
	if err != nil {
		return err
	}

	err = sms.FilecoinNode().SubmitMessage(signedMessage)
	if err != nil {
		return err
	}

	return nil
}
