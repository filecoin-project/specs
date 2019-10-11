package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

func (provider *StorageProvider_I) HandleNewStorageDealProposal(proposal deal.StorageDealProposal) {
	if provider.verifyStorageClient(proposal.Client(), proposal.ProposerSignature(), proposal.StoragePrice()) {
		provider.DealStatus()[proposal.PieceRef()] = deal.StorageDealProposed{}
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

func (provider *StorageProvider_I) HandleStorageDealQuery(dealCID deal.DealCID) deal.StorageDealStatus {
	dealStatus, found := provider.DealStatus()[dealCID]

	if found {
		return dealStatus
	}

	return deal.StorageDealNotFound{}
}
