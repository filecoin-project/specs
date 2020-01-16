package sector_index

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	smarkact "github.com/filecoin-project/specs/actors/builtin/storage_market"
)

func (sis *SectorIndexerSubsystem_I) AddNewDeal(deal smarkact.StorageDeal) StageDealResponse {
	return sis.Builder().StageDeal(deal)
}

// func (sis *SectorIndexerSubsystem_I) OnNewTipset(chain blockchain.Chain, epoch blockchain.Epoch) {
// 	panic("TODO")
// }

func (sis *SectorIndexerSubsystem_I) SectorsExpiredAtEpoch(epoch abi.ChainEpoch) []abi.SectorID {
	panic("TODO")
}

func (sis *SectorIndexerSubsystem_I) removeSectors(sectorIDs []abi.SectorID) {
	panic("TODO")
}
