package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
import filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import address "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/address"
import base_blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain"
import blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain/blockchain"
import base_markets "github.com/filecoin-project/specs/systems/filecoin_markets"

// import storage_proving "github.com/filecoin-project/specs/systems/filecoin_mining/storage_proving"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

func (sms *StorageMiningSubsystem_I) CreateMiner(ownerPubKey filcrypto.PubKey, workerPubKey filcrypto.PubKey, pledgeAmt actor.TokenAmount) address.Address {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	return sms.StoragePowerActor().RegisterMiner(ownerAddr, workerPubKey)
}

func (sms *StorageMiningSubsystem_I) HandleStorageDeal(deal base_markets.StorageDeal, pieceRef ipld.CID) {
	stagedDealResponse := sms.SectorIndex().AddNewDeal(deal)
	sms.StorageProvider().NotifyStorageDealStaged(&StorageDealStagedNotification_I{
		Deal_:     deal,
		PieceRef_: pieceRef,
		SectorID_: stagedDealResponse.SectorID(),
	})
}

func (sms *StorageMiningSubsystem_I) generateOwnerAddress(workerPubKey filcrypto.PubKey) address.Address {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() base_markets.StorageDeal {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) OnNewTipset(chain blockchain.Chain, epoch blockchain.Epoch, tipset blockchain.Tipset) {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) OnNewRound(newTipset blockchain.Tipset) base_blockchain.ElectionArtifacts {
	panic("TODO: fix this below")

	// TODO this below has been commented due to incomplete implementation
	// ea := sms.Consensus().GetElectionArtifacts(sms.CurrentChain, sms.CurrentEpoch)
	// EP := sms.DrawElectionProof(ea.TK(), sms.workerPrivateKey)
	// if newTipset {
	// 	T0 := GenerateNextTicket(ea.T1, workerPrivateKey)
	// } else {
	// 	T1 := GenerateNextTicket(T0, workerPrivateKey)
	// }

	// if sms.Consensus().TryLeaderElection(EP) {
	// 	// TODO: move this into SPC or Blockchain
	// 	// SMS should probably not have ability to call BlockProducer directly.
	// 	sms.BlockProducer().GenerateBlock(EP, ea.T1(), sms.CurrentTipset(), workerKey)
	// } else {
	// 	// TODO when not elected
	// }

	// return ea
}

func (sms *StorageMiningSubsystem_I) DrawElectionProof(tk base_blockchain.Ticket, workerKey filcrypto.PrivKey) base_blockchain.ElectionProof {
	// return generateElectionProof(tk, workerKey)
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) GenerateNextTicket(t1 base_blockchain.Ticket, workerKey filcrypto.PrivKey) base_blockchain.Ticket {
	panic("TODO")
}

// TODO this should be moved into storage market
func (sp *StorageProvider_I) NotifyStorageDealStaged(storageDealNotification StorageDealStagedNotification) {
	panic("TODO")
}
