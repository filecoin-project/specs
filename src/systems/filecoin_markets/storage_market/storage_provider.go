package storage_market

import base_blockchain "github.com/filecoin-project/specs/systems/filecoin_blockchain"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

func (provider *StorageProvider_I) HandleNewStorageDealProposal(proposal deal.StorageDealProposal) {
	panic("TODO")
}

func (provider *StorageProvider_I) signStorageDealProposal(proposal deal.StorageDealProposal) deal.StorageDeal {
	panic("TODO")
}

func (provider *StorageProvider_I) verifyStorageClient(address base_blockchain.Address, signature deal.Signature) bool {
	panic("TODO")
}

func (provider *StorageProvider_I) NotifyStorageDealStaged(storageDealNotification StorageDealStagedNotification) {
	panic("TODO")
}

func (provider *StorageProvider_I) HandleStorageDealQuery(dealCID deal.DealCID) deal.StorageDealStatus {
	panic("TODO")
}
