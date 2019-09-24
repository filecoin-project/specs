func (sms *StorageMiningSubsystem) CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount) StorageMinerActor {
	ownerAddr := sms.generateOwnerAddress(workerPubKey)
	return spa.RegisterMiner(ownerAddr, workerPubKey)
}

type StorageDealStagedNotification = struct {
	Deal StorageDeal
	PieceRef CID
	Pip PieceInclusionProof
	SectorID SectorID
}

func (sms *StorageMiningSubsystem) HandleStorageDeal(deal StorageDeal, pieceRef CID) {
	AddDealToSectorResponse := sms.sectorIndexer.AddDealToSector(deal)
	storageProvider.NotifyStorageDealStaged(StorageDealStagedNotification {
		Deal: deal,
		PieceRef: pieceRef,
		Pip: AddDealToSectorResponse.pip,
		SectorID: AddDealToSectorResponse.sectorID,
	})
}

func (sms *StorageMiningSubsystem) generateOwnerAddress(workerPubKey PubKey) Addr {
	panic("TODO")
}

func (sms *StorageMiningSubsystem) CommitSectorError() StorageDeal {

}

func (sms *StorageMiningSubsystem) OnNewTipset(chain Chain, epoch Epoch, tipset Tipset) struct {} {
	sms.CurrentChain = chain
	sms.CurrentEpoch = epoch
	sms.CurrentTipset = tipset
}

func (sms *StorageMiningSubsystem) OnNewRound(newTipset Tipset) ElectionArtifacts {
	TK := storagePowerConsensus.TicketAtEpoch(sms.CurrentChain, sms.CurrentEpoch - k)
	T1 := storagePowerConsensus.TicketAtEpoch(sms.CurrentChain, sms.CurrentEpoch - 1)
	EP := DrawElectionProof(TK, workerPrivateKey)
	if newTipset {
		T0 := GenerateNextTicket(T1, workerPrivateKey)
	} else {
		T1 := GenerateNextTicket(T0, workerPrivateKey)
	}
	if storagePowerConsensus.TryLeaderElection(EP) {
		BlockProducer.GenerateBlock(EP, T0, sms.CurrentTipset, workerKey)
	} else {
		// TODO when not elected
	}	
}

func (sms *StorageMiningSubsystem) DrawElectionProof(tk Ticket, workerKey PrivateKey) ElectionProof {
	return generateElectionProof(tk, workerKey)
}

func (sms *StorageMiningSubsystem) GenerateNextTicket(t1 Ticket, workerKey PrivateKey) Ticket {

}