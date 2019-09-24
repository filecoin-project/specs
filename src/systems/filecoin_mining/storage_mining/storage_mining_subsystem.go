package fileName

// func NewStorageMiningSubsystem() *StorageMiningSubsystem {
// 	storageMinerActors := []StorageMinerActor{}
// 	sectorIndexerSubsystem := InitStorageIndexerSubsystem()
// 	return &StorageMiningSubsystem{
// 		storageMinerActors: storageMinerActors,
// 		sectorIndexerSubsystem: sectorIndexerSubsystem,
// 	}
// }

func (sms *StorageMiningSubsystem) CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount) StorageMinerActor {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	return spa.RegisterMiner(ownerAddr, workerPubKey)
}

func (sms *StorageMiningSubsystem) HandleStorageDeal(deal StorageDeal, pieceRef CID) {
	AddDealToSectorResponse := sms.sectorIndexer.AddDealToSector(deal)
	storageProvider.NotifyStorageDealStaged(struct {
		Deal: deal,
		PieceRef: pieceRef,
		Pip: AddDealToSectorResponse.pip,
		SectorID: AddDealToSectorResponse.sectorID
	})
}

func (sms *StorageMiningSubsystem) generateOwnerAddress(workerPubKey PubKey) Addr {
	panic("TODO")
}

func (sms *StorageMiningSubsystem) CommitSectorError() StorageDeal {

}

func (sms *StorageMiningSubsystem) OnNewTipset(chain Chain, epoch Epoch, tipset Tipset) struct {} {
	sms.currentChain := chain
	sms.currentEpoch := epoch
	sms.currentTipset := tipset
}

func (sms *StorageMiningSubsystem) OnNewRound() ElectionArtifacts {
	TK := storagePowerConsensus.TicketAtEpoch(sms.chain, sms.epoch - k)
	T1 := storagePowerConsensus.TicketAtEpoch(sms.chain, sms.epoch - 1)
	EP := DrawElectionProof(TK, workerPrivateKey)
	if NewTipset {
		T0 := GenerateNextTicket(T1, workerPrivateKey)
	} else {
		T1 := GenerateNextTicket(T0, workerPrivateKey)
	}

	if storagePowerConsensus.TryLeaderElection(EP) {
		BlockProducer.GenerateBlock(EP, T0, sms.currentTipset, workerKey)
	} else {

	}
}

func (sms *StorageMiningSubsystem) DrawElectionProof(tk Ticket, workerKey PrivateKey) ElectionProof {
	return generateElectionProof(tk, workerKey)
}

func (sms *StorageMiningSubsystem) GenerateNextTicket(t1 Ticket, workerKey PrivateKey) Ticket {

}