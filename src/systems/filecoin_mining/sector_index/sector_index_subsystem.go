package sector_index

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
)

func (sis *SectorIndexerSubsystem_I) AddNewDeal(deal deal.StorageDeal) StageDealResponse {
	return sis.Builder().StageDeal(deal)
}

// func (sis *SectorIndexerSubsystem_I) OnNewTipset(chain blockchain.Chain, epoch blockchain.Epoch) {
// 	panic("TODO")
// }

func (sis *SectorIndexerSubsystem_I) SectorsExpiredAtEpoch(epoch abi.ChainEpoch) []sector.SectorID {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem_I) removeSectors(sectorIDs []sector.SectorID) {
	panic("TODO")
}
