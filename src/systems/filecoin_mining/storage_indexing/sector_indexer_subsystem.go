package fileName

// func NewSectorIndexerSubsystem() *SectorIndexerSubsystem {
// 	panic("TODO")
// 	return struct{}
// }

func (sis *SectorIndexerSubsystem) AddDealToSector(deal StorageDeal) AddDealToSectorResponse {
	sectorBuilder := sis.selectSectorBuilderByDeal(deal)
	sectorBuilder.AddPiece(pieceRef)
	sis.SectorStore.AddPieceByRef(pieceRef)
	pip := sis.getPieceInclusionProof(deal)

	return AddDealToSectorResponse{
		sectorID: sectorBuilder.sectorID,
		pip: pip,
	}
}

func (sis *SectorIndexerSubsystem) selectSectorBuilderByDeal(deal StorageDeal) SectorBuilder {
	panic("TODO")
}


func (sis *SectorIndexerSubsystem) indexSectorByDealExpiration(sectorID SectorID, deal StorageDeal) {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) getPieceInclusionProof(deal StorageDeal) PieceInclusionProof {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) OnNewTipset(chain Chain, epoch Epoch) {
	sectorIDs := lookupSectorByExpiry(epoch)
	purgeSectorWithNoLiveDeals(sectorIDs)
}

func (sis *SectorIndexerSubsystem) lookupSectorByExpiry(currentEpoch Epoch) []SectorID {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) purgeSectorWithNoLiveDeals(sectorIDs []SectorID) {
	panic("TODO")
}
