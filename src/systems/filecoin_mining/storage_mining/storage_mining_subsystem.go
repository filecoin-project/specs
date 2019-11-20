package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

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
	// new election, incremented height
	Randomness1 := sms.consensus().GetTicketProductionSeed(sms.blockchain().BestChain(), sms.blockchain().LatestEpoch())
	RandomnessK := sms.consensus().GetElectionProofSeed(sms.blockchain().BestChain(), sms.blockchain().LatestEpoch())

	for _, worker := range sms.keyStore().Workers() {
		newTicket := sms.PrepareNewTicket(Randomness1, worker.VRFKeyPair())
		newEP := sms.DrawElectionProof(RandomnessK, sms.blockchain().LatestEpoch(), worker.VRFKeyPair())

		if sms.consensus().IsWinningElectionProof(newEP, worker.Address()) {
			sms.blockProducer().GenerateBlock(newEP, newTicket, sms.blockchain().BestChain().HeadTipset(), worker.Address())
		}
	}
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(randomness block.Randomness, vrfKP filcrypto.VRFKeyPair) block.Ticket {
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

func (sms *StorageMiningSubsystem_I) DrawElectionProof(randomness block.Randomness, height block.ChainEpoch, vrfKP filcrypto.VRFKeyPair) block.ElectionProof {
	panic("")
	// // 0. Prepare new election proof
	// var newEP ElectionProof

	// // 1. Run it through VRF and get determinstic output
	// // 1.i. # take the VRFOutput of that ticket as input, specified for the appropriate operation type
	// input := VRFPersonalization.ElectionProof
	// input.append(lookbackTicket.Output)
	// input.append(height)
	// // ii. # run it through the VRF and store the VRFProof in the new ticket
	// newEP.VRFResult := vrfKP.Generate(input)
	// return newEP
}
