---
title: Sector Storage
entries:
  - components
# suppressMenu: true
---

type Piece []byte

struct SectorStorageSubsystem {
       SectorSealer   SectorSealer,
       SectorStore    SectorStore,
       SectorBuilders { SectorConfig : SectorBuilder },

       ProverId       UInt

       AddPiece(
         piece Piece,
	 sectorConfig SectorConfig,
       ) Error | bool {
           sectorbuilder := SectorBuilders[sectorConfig];
       	   piecePath := SectorStore.WritePiece(piece);
	   response := sectorBuilder.addPiece(PiecePath);

	   MaybeSeal(response.StagedSector, response.BytesRemaining);
       }

       MaybeSeal(StagedSector StagedSector, BytesRmaining Uint);
       
       Seal(stagedSector StagedSector) {
       	    sealedPath := SectorStore.AllocateSealedSector(stagedSector.SectorSize);
	    return SectorSealer.seal(stagedSector, sealedPath, ProverId);
       }

       RetrievePiece(CommD Commitment, Start UInt, Length UInt) Piece {
       }             
}
