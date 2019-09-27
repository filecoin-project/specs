package sector_index

func (sis *SectorIndexerSubsystem) AddDealToSector(deal StorageDeal) AddDealToSectorResponse {
	addPieceResponse := sis.SectorBuilder.AddPiece(deal.PiecePath)
	pip := sis.StorageProofs.GetPieceInclusionProof(deal.PieceRef)

	return AddDealToSectorResponse{
		sectorID: addPieceResponse.SectorInfo.ID,
		pip: pip,
	}

	// sectorbuilder := SectorBuilders[sectorConfig];
	// piecePath := SectorStore.WritePiece(piece);
  // response := sectorBuilder.addPiece(PiecePath);
  // MaybeSeal(response.StagedSector, response.BytesRemaining);
}

// func (sis *SectorIndexerSubsystem) Seal(stagedSector StagedSector) {
// 	sealedPath := SectorStore.AllocateSealedSector(stagedSector.SectorSize);
// 	response := SectorSealer.Seal(stagedSector, sealedPath, ProverId);
// 	SectorStore.RegisterMetadata(
// 	SectorMetadata {
// 	  response.CommR,
// 		response.PersistentAux,
// 		response.PartialMerkleTreePath,
// 	});
// }

func (sis *SectorIndexerSubsystem) selectSectorBuilderByDeal(deal StorageDeal) SectorBuilder {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) indexSectorByDealExpiration(sectorID SectorID, deal StorageDeal) {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem) getPieceInclusionProof(pieceRef CID) PieceInclusionProof {
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
