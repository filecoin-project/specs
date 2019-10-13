package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import (
	filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	address "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

func (sms *StorageMiningSubsystem_I) CreateMiner(ownerPubKey filcrypto.PubKey, workerPubKey filcrypto.PubKey, sectorSize util.UInt, peerId libp2p.PeerID) address.Address {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	// var pledgeAmt actor.TokenAmount TODO: unclear how to pass the amount/pay
	// TODO compute PledgeCollateral for 0 bytes
	return sms.StoragePowerActor().CreateStorageMiner(ownerAddr, workerPubKey, sectorSize, peerId)
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

func (sms *StorageMiningSubsystem_I) generateOwnerAddress(workerPubKey filcrypto.PubKey) address.Address {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() deal.StorageDeal {
	panic("TODO")
}

// triggered by new block reception and tipset assembly
func (sms *StorageMiningSubsystem_I) OnNewBestChain() {
	panic("")
	// new election, reset nonce
	// sms.electionNonce = 0
	// sms.tryLeaderElection()
}

// triggered by wall clock
func (sms *StorageMiningSubsystem_I) OnNewRound() {
	panic("")
	// repeat on prior tipset, increment nonce
	// sms.electionNonce += 1
	// sms.tryLeaderElection()
}

func (sms *StorageMiningSubsystem_I) tryLeaderElection() {
	panic("")
	// T1 := sms.Consensus.GetTicketProductionSeed(sms.CurrentChain, sms.Blockchain.LatestEpoch())
	// TK := sms.Consensus.GetElectionProofSeed(sms.CurrentChain, sms.Blockchain.LatestEpoch())

	// for _, worker := range sms.workers {
	// 	newTicket := PrepareNewTicket(worker.VRFKeyPair, T1)
	// 	newEP := DrawElectionProof(TK, sms.electionNonce, worker.VRFKeyPair)

	// 	if sms.Consensus.IsWinningLeaderElection(newEP, worker.address) {
	// 		BlockProducer.GenerateBlockHeader(newEP, newTicket, sms.CurrentTipset, worker.workerAddress)
	// 	}
	// }
}

func (sms *StorageMiningSubsystem_I) PrepareNewTicket(priorTicket block.Ticket, vrfKP filcrypto.VRFKeyPair) block.Ticket {
	panic("")
	// // 0. prepare new ticket
	// var newTicket Ticket

	// // 1. run it through the VRF and get deterministic output
	// // 1.i. take the VRFResult of that ticket as input, specifying the personalization (see data structures)
	// input := VRFPersonalization.Ticket
	// input.append(priorTicket.Output)
	// // 2.ii. run through VRF
	// newTicket.VRFResult := vrfKP.Generate(input)

	// return newTicket
}

func (sms *StorageMiningSubsystem_I) DrawElectionProof(lookbackTicket block.Ticket, nonce block.ElectionNonce, vrfKP filcrypto.VRFKeyPair) block.ElectionProof {
	panic("")
	// // 0. Prepare new election proof
	// var newEP ElectionProof

	// // 1. Run it through VRF and get determinstic output
	// // 1.i. # take the VRFOutput of that ticket as input, specified for the appropriate operation type
	// input := VRFPersonalization.ElectionProof
	// input.append(lookbackTicket.Output)
	// input.append(nonce)
	// // ii. # run it through the VRF and store the VRFProof in the new ticket
	// newEP.VRFResult := vrfKP.Generate
	// newEP.ElectionNonce := nonce
	// return newEP
}
