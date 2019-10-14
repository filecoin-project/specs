package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (sma *StorageMarketActor_I) WithdrawBalance(balance actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this from VM runtime

	if balance < 0 {
		panic("Error")
	}

	senderBalance, found := sma.Balances()[msgSender]
	if !found {
		panic("Error") // TODO replace with fatal
	}

	if senderBalance.Available() < balance {
		panic("Error") // TODO replace with fatal
	}

	senderBalance.Impl().Available_ = senderBalance.Available() - balance
	sma.Balances()[msgSender] = senderBalance

	// TODO send funds to msgSender with `send` in VM runtime
}

func (sma *StorageMarketActor_I) AddBalance(balance actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this

	// TODO subtract balance from msgSender
	// TODO add balance to StorageMarketActor
	if balance < 0 {
		panic("Error")
	}

	senderBalance, found := sma.Balances()[msgSender]
	if found {
		senderBalance.Impl().Available_ = senderBalance.Available() + balance
		sma.Balances()[msgSender] = senderBalance
	} else {
		sma.Balances()[msgSender] = &StorageParticipantBalance_I{
			Locked_:    0,
			Available_: balance,
		}
	}
}

func (sma *StorageMarketActor_I) CheckLockedBalance(participantAddr addr.Address) actor.TokenAmount {
	var msgSender addr.Address // TODO replace this

	return sma.Balances()[msgSender].Locked()
}

func (sma *StorageMarketActor_I) PublishStorageDeal(newStorageDeals []deal.StorageDeal) []PublishStorageDealResponse {
	l := len(newStorageDeals)
	response := make([]PublishStorageDealResponse, l)
	for i, newDeal := range newStorageDeals {
		if sma.verifyStorageDeal(newDeal) {
			id := sma.generateStorageDealID(newDeal)
			sma.Deals()[id] = newDeal
			response[i] = PublishStorageDealSuccess
		} else {
			response[i] = PublishStorageDealError
		}
	}

	return response
}

func (sma *StorageMarketActor_I) verifyStorageDeal(d deal.StorageDeal) bool {
	// TODO verify proposal or deal has not expired
	// TODO verify client and provider signature
	// TODO verfiy minimum StoragePrice and StorageCollateral
	p := d.Proposal()
	clientBalanceA := sma.Balances()[p.Client()].Available()
	providerBalanceA := sma.Balances()[p.Provider()].Available()

	if clientBalanceA < p.StoragePrice() ||
		providerBalanceA < p.StorageCollateral() {
		return false
	}
	return true
}

func (sma *StorageMarketActor_I) generateStorageDealID(storageDeal deal.StorageDeal) deal.DealID {
	panic("TODO")
}

func (sma *StorageMarketActor_I) HandleCronAction() {
	panic("TODO")
}

func (sma *StorageMarketActor_I) SettleExpiredDeals(storageDealIDs []deal.DealID) {
	// for dealID := range storageDealIDs {
	// Return the storage collateral
	// storageDeal := sma.Deals()[dealID]
	// storageCollateral := storageDeal.StorageCollateral()
	// provider := storageDeal.Provider()
	// assert(sma.Balances()[provider].Locked() >= storageCollateral)

	// // Move storageCollateral from locked to available
	// balance := sma.Balances()[provider]

	// sma.Balances()[provider] = &StorageParticipantBalance_I{
	// 	Locked_:    balance.Locked() - storageCollateral,
	// 	Available_: balance.Available() + storageCollateral,
	// }

	// // Delete reference to the deal
	// delete(sma.Deals_, dealID)
	// }
	panic("TODO")
}

func (sma *StorageMarketActor_I) handleStorageDealPayment(storageDealIDs []deal.DealID) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) slashStorageDealCollateral(storageDealIDs []deal.DealID) {
	// for _, dealID := range storageDealIDs {
	// 	faultStorageDeal := sma.Deals()[dealID]
	// TODO remove locked funds and send slashed fund to TreasuryActor
	// TODO provider lose power for the FaultSet but not PledgeCollateral
	// }
	panic("TODO")
}
