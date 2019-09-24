func (sis *SectorIndexerSubsystem) AddDealToSector(deal StorageDeal) AddDealToSectorResponse {
	addPieceResponse := sis.SectorBuilder.AddPiece(deal.PiecePath)
	pip := sis.StorageProofs.getPieceInclusionProof(deal)

	return AddDealToSectorResponse{
		sectorID: addPieceResponse.SectorInfo.ID
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
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) lookupSectorByExpiry(currentEpoch Epoch) []SectorID {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) purgeSectorWithNoLiveDeals(sectorIDs []SectorID) {
	panic("TODO")
}
