package storage_mining

func (sms *StorageMiningSubsystem_I) CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount) StorageMinerActor {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	return spa.RegisterMiner(ownerAddr, workerPubKey)
}

type StorageDealStagedNotification = struct {
	Deal     StorageDeal
	PieceRef CID
	Pip      PieceInclusionProof
	SectorID SectorID
}

func (sms *StorageMiningSubsystem_I) HandleStorageDeal(deal StorageDeal, pieceRef CID) {
	AddDealToSectorResponse := sms.sectorIndexer.AddDealToSector(deal)
	storageProvider.NotifyStorageDealStaged(StorageDealStagedNotification{
		Deal:     deal,
		PieceRef: pieceRef,
		Pip:      AddDealToSectorResponse.pip,
		SectorID: AddDealToSectorResponse.sectorID,
	})
}

func (sms *StorageMiningSubsystem_I) generateOwnerAddress(workerPubKey PubKey) Addr {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) CommitSectorError() StorageDeal {
	panic("TODO")
}

func (sms *StorageMiningSubsystem_I) OnNewTipset(chain Chain, epoch Epoch, tipset Tipset) struct{} {
	sms.CurrentChain = chain
	sms.CurrentEpoch = epoch
	sms.CurrentTipset = tipset
}

func (sms *StorageMiningSubsystem_I) OnNewRound(newTipset Tipset) ElectionArtifacts {
	ea := storagePowerConsensus.ElectionArtifacts(sms.CurrentChain, sms.CurrentEpoch)
	EP := DrawElectionProof(ea.TK, sms.workerPrivateKey)

	panic("TODO: fix this below")
	// if newTipset {
	// 	T0 := GenerateNextTicket(ea.T1, workerPrivateKey)
	// } else {
	// 	T1 := GenerateNextTicket(T0, workerPrivateKey)
	// }

	if storagePowerConsensus.TryLeaderElection(EP) {
		// TODO: move this into SPC or Blockchain
		// SMS should probably not have ability to call BlockProducer directly.
		BlockProducer.GenerateBlock(EP, T0, sms.CurrentTipset, workerKey)
	} else {
		// TODO when not elected
	}
}

func (sms *StorageMiningSubsystem_I) DrawElectionProof(tk Ticket, workerKey PrivateKey) ElectionProof {
	return generateElectionProof(tk, workerKey)
}

func (sms *StorageMiningSubsystem_I) GenerateNextTicket(t1 Ticket, workerKey PrivateKey) Ticket {
	panic("TODO")
}
