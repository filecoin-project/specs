package fileName


func InitSectorIndexerSubsystem() *SectorIndexerSubsystem {
	panic("TODO")
	return struct{}
}

func (sis *SectorIndexerSubsystem) AddPieceToSector(deal StorageDeal, pieceRef CID) struct {
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

