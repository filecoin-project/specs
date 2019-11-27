package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
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

	// Draw randomness from chain for ElectionPoSt and Ticket Generation
	// Randomness for ticket generation in block production
	randomness1 := sms._consensus().GetTicketProductionSeed(sms._blockchain().BestChain(), sms._keyStore().OwnerAddress(), sms._blockchain().LatestEpoch())

	// Randomness for ElectionPoSt
	randomnessK := sms._consensus().GetPoStChallenge(sms._blockchain().BestChain(), sms._keyStore().OwnerAddress(), sms._blockchain().LatestEpoch())

	// TODO: @why @jz align on this
	for _, worker := range sms._keyStore().Workers() {

		var input []byte
		input = append(input, spc.VRFPersonalizationPoSt)
		input = append(input, randomnessK...)

		postRandomness := worker.VRFKeyPair().Impl().Generate(input).Output()
		// TODO: add how sectors are actually stored in the SMS proving set
		provingSet := make([]sector.SectorID, 0)

		candidates := sms.StorageProving().Impl().GenerateElectionPoStCandidates(postRandomness, provingSet)

		if len(candidates) <= 0 {
			return // fail to generate post candidates
		}

		winningCandidates := make([]sector.PoStCandidate, 0)

		for _, candidate := range candidates {
			// TODO align on worker address
			if sms._consensus().IsWinningPartialTicket(candidate.PartialTicket()) {
				winningCandidates = append(winningCandidates, candidate)
			}
		}

		if len(winningCandidates) <= 0 {
			return
		}

		newTicket := sms.PrepareNewTicket(randomness1, worker.VRFKeyPair())
		postProof := sms.StorageProving().Impl().GenerateElectionPoStProof(postRandomness, winningCandidates)
		chainHead := sms._blockchain().BestChain().HeadTipset()

		sms._blockProducer().GenerateBlock(postProof, winningCandidates, newTicket, chainHead, worker.Address())

	}
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(randomness util.Randomness, vrfKP filcrypto.VRFKeyPair) block.Ticket {
	// run it through the VRF and get deterministic output

	// take the VRFResult of that ticket as input, specifying the personalization (see data structures)
	var input []byte
	input = append(input, spc.VRFPersonalizationTicket)
	input = append(input, randomness...)

	// run through VRF
	vrfRes := vrfKP.Generate(input)

	newTicket := &block.Ticket_I{
		VRFResult_: vrfRes,
		Output_:    vrfRes.Output(),
	}

	return newTicket
}
