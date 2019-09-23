package fileName

func InitStorageMiningSubsystem() *StorageMiningSubsystem {
	storageMinerActors := []StorageMinerActor{}
	sectorIndexerSubsystem := InitStorageIndexerSubsystem()
	return &StorageMiningSubsystem{
		storageMinerActors: storageMinerActors,
		sectorIndexerSubsystem: sectorIndexerSubsystem,
	}
}

func (sms *StorageMiningSubsystem) CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount) StorageMinerActor {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	return spa.RegisterMiner(ownerAddr, workerPubKey)
}

func (sms *StorageMiningSubsystem) HandleStorageDeal(deal StorageDeal, pieceRef CID) {
	AddPieceToSectorReturn := sms.sectorIndexerSubsystem.AddPieceToSector(deal, pieceRef)
	storageProvider.NotifyStorageDealStaged(struct {
		Deal: deal,
		PieceRef: pieceRef,
		Pip: AddPieceToSectorReturn.pip,
		SectorId: AddPieceToSectorReturn.sectorID
	})
}

func (sms *StorageMiningSubsystem) CommitmentSectorError() StorageDeal {

}

func (sms *StorageMiningSubsystem) OnNewTipset(chain Chain, epoch Epoch) StorageDeal {

}

func (sms *StorageMiningSubsystem) OnNewRound() StorageDeal {

}

func (sms *StorageMiningSubsystem) DrawElectionProof(tk Randomness, workerKey PrivateKey) ElectionProof {

}

func (sms *StorageMiningSubsystem) GenerateNextTicket(seed Randomness, workerKey PrivateKey) Ticket {

}