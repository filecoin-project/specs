package fileName

// func NewSectorIndexerSubsystem() *SectorIndexerSubsystem {
// 	panic("TODO")
// 	return struct{}
// }

func (sis *SectorIndexerSubsystem) AddDealToSector(deal StorageDeal) struct {
	sectorID SectorID,
	pip PieceInclusionProof
} {
	sectorBuilder := sis.selectSectorBuilderByDeal(deal)
	sectorBuilder.AddPiece(pieceRef)
	sis.SectorStore.AddPieceByRef(pieceRef)
	pip := sis.getPieceInclusionProof(deal)

	return struct {
		sectorID: sectorBuilder.sectorID,
		pip: pip,
	}
}

func getPieceInclusionProof(deal StorageDeal) PieceInclusionProof {
	panic("TODO")
}

