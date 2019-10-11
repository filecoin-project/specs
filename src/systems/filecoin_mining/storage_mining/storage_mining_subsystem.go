package storage_mining

// import sectoridx "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
// import spc "github.com/filecoin-project/specs/systems/filecoin_blockchain/storage_power_consensus"
// import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import filcrypto "github.com/filecoin-project/specs/libraries/filcrypto"
import libp2p "github.com/filecoin-project/specs/libraries/libp2p"
import address "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
import util "github.com/filecoin-project/specs/util"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

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
	// sms.StorageProvider().NotifyStorageDealStaged(&storage_provider.StorageDealStagedNotification_I{
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

func (sms *StorageMiningSubsystem_I) OnNewTipset(chain chain.Chain, epoch block.ChainEpoch, tipset block.Tipset) {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) OnNewRound(newTipset block.Tipset) block.ElectionArtifacts {
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

func (sms *StorageMiningSubsystem_I) DrawElectionProof(tk block.Ticket, workerKey filcrypto.PrivKey) block.ElectionProof {
	// return generateElectionProof(tk, workerKey)
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) GenerateNextTicket(t1 block.Ticket, workerKey filcrypto.PrivKey) block.Ticket {
	panic("TODO")
}
