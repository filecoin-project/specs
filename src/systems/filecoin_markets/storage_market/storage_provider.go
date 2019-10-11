package storage_market

import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

func (provider *StorageProvider_I) HandleNewStorageDealProposal(proposal deal.StorageDealProposal) {
	panic("TODO")
}

func (provider *StorageProvider_I) SignStorageDealProposal(proposal deal.StorageDealProposal) deal.StorageDeal {
	panic("TODO")
}

func (provider *StorageProvider_I) NotifyStorageDealStaged(storageDealNotification deal.StorageDealStagedNotification) {
	panic("TODO")
}

func (provider *StorageProvider_I) HandleStorageDealQuery(dealCID deal.DealCID) deal.StorageDealStatus {
	panic("TODO")
}

