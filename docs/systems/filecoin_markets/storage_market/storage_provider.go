package storage_market

import (
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"

	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
)

// import deal_status "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market"

func (provider *StorageProvider_I) HandleNewStorageDealProposal(proposal deal.StorageDealProposal) {
	if provider.verifyStorageClient(proposal.Client(), proposal.ProposerSignature(), proposal.StoragePrice()) {
		// status := &deal.StorageDealStatus_StorageDealProposed_I{}
		// s := deal.StorageDealStatus_Make_StorageDealProposed(status)
		provider.DealStatus()[proposal.PieceRef()] = StorageDealProposed
		// TODO notify StorageClient that a deal has been received
		// TODO notify StorageMiningSubsystem to add deals to sector
		provider.signStorageDealProposal(proposal)
		// DO THIS TODAY Call StorageMarketActor.publishStorageDeal()
	}
}

func (provider *StorageProvider_I) signStorageDealProposal(proposal deal.StorageDealProposal) deal.StorageDeal {
	// TODO add signature to the proposal
	// TODO notify StorageClient that a deal has been signed
	panic("TODO")
}

func (provider *StorageProvider_I) rejectStorageDealProposal(proposal deal.StorageDealProposal) {
	provider.DealStatus()[proposal.PieceRef()] = StorageDealRejected
	// TODO send notification to client
}

func (provider *StorageProvider_I) verifyStorageClient(address addr.Address, signature deal.Signature, price actor.TokenAmount) bool {
	// TODO make call to StorageMarketActor
	// balance, found := StorageMarketActor.Balances()[address]

	// if !found {
	// 	return false
	// }

	// if balance < price {
	// 	return false
	// }

	// TODO Check on Signature
	// return true
	panic("TODO")
}

// TODO: func (provider *StorageProvider_I) NotifyStorageDealStaged(storageDealNotification StorageDealStagedNotification) {
// 	panic("TODO")
// }

func (provider *StorageProvider_I) HandleStorageDealQuery(dealCID deal.DealCID) StorageDealStatus {
	dealStatus, found := provider.DealStatus()[dealCID]

	if found {
		return dealStatus
	}

	return StorageDealNotFound
}

// TODO this should be moved into storage market
func (sp *StorageProvider_I) NotifyStorageDealStaged(storageDealNotification StorageDealStagedNotification) {
	panic("TODO")
}
