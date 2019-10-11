package storage_market

import base_blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

func (sma *StorageMarketActor_I) WithdrawBalance(balance deal.TokenAmount) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) AddBalance(balance deal.TokenAmount) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) CheckLockedBalance(participantAddr base_blockchain.Address) deal.TokenAmount {
	panic("TODO")
}

func (sma *StorageMarketActor_I) PublishStorageDeal(newStorageDeals []deal.StorageDeal) []PublishStorageDealResponse {
	panic("TODO")
}

func (sma *StorageMarketActor_I) verifyStorageDeal(newStorageDeal deal.StorageDeal) bool {
	panic("TODO")
}

func (sma *StorageMarketActor_I) generateStorageDealID(storageDeal deal.StorageDeal) deal.DealID {
	panic("TODO")
}

func (sma *StorageMarketActor_I) HandleCronAction() {
	panic("TODO")
}

func (sma *StorageMarketActor_I) closeExpiredStorageDeal(storageDealIDs []deal.DealID) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) handleStorageDealPayment(storageDealIDs []deal.DealID, currEpoch base_blockchain.Epoch) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) slashStorageDealCollateral(storageDealIDs []deal.DealID) {
	panic("TODO")
}
