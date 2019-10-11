package sector_index

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	// piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	mkt "github.com/filecoin-project/specs/systems/filecoin_markets"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
)

func (sis *SectorIndexerSubsystem_I) AddNewDeal(deal mkt.StorageDeal) StageDealResponse {
	return sis.Builder().StageDeal(deal)
}

// func (sis *SectorIndexerSubsystem_I) OnNewTipset(chain blockchain.Chain, epoch blockchain.Epoch) {
// 	panic("TODO")
// }

func (sis *SectorIndexerSubsystem_I) SectorsExpiredAtEpoch(epoch block.ChainEpoch) []sector.SectorID {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem_I) removeSectors(sectorIDs []sector.SectorID) {
	panic("TODO")
}
