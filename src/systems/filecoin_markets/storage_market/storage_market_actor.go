package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (sma *StorageMarketActor_I) WithdrawBalance(balance actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this

	if balance < 0 {
		panic("Error")
	}

	currentBalance, found := sma.Balances()[msgSender]
	if !found {
		panic("Error") // TODO replace with fatal
	}

	if currentBalance.Available < balance {
		panic("Error") // TODO replace with fatal
	}

	sma.Balances()[msgSender].Available -= balance
	// TODO send funds to msgSender
}

func (sma *StorageMarketActor_I) AddBalance(balance actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this

	// TODO subtract balance from msgSender
	// TODO add balance to StorageMarketActor
	if balance < 0 {
		panic("Error")
	}

	sma.Balances()[msgSender].Available += balance
}

func (sma *StorageMarketActor_I) CheckLockedBalance(participantAddr addr.Address) actor.TokenAmount {
	var msgSender addr.Address // TODO replace this

	return sma.Balances()[msgSender].Locked()
}

func (sma *StorageMarketActor_I) PublishStorageDeal(newStorageDeals []deal.StorageDeal) []deal.PublishStorageDealResponse {
	l := len(newStorageDeals)
	var response [l]deal.PublishStorageDealResponse
	for i, newDeal := range newStorageDeals {
		var publishResponse deal.PublishStorageDealResponse
		if sma.verifyStorageDeal(newDeal) {
			id := sma.generateStorageDealID(newDeal)
			sma.Deals()[id] = newDeal
			publishResponse := deal.PublishStorageDealSuccess{id, newDeal.PieceRef()}
		} else {
			publishResponse := deal.PublishStorageDealError{}
		}

		response[i] = publishResponse
	}

	return response
}

func (sma *StorageMarketActor_I) verifyStorageDeal(newStorageDeal deal.StorageDeal) bool {
	// TODO verify proposal or deal has not expired
	// TODO verify client and provider signature
	// TODO verfiy minimum StoragePrice and StorageCollateral
	if sma.Balances()[newStorageDeal.Proposal().Client()].Available() < newStorageDeal.Proposal().StoragePrice() ||
		sma.Balances()[newStorageDeal.Proposal().Provider()].Available() < newStorageDeal.Proposal().StorageCollateral() {
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

func (sma *StorageMarketActor_I) closeExpiredStorageDeal(storageDealIDs []deal.DealID) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) handleStorageDealPayment(storageDealIDs []deal.DealID) {
	panic("TODO")
}

func (sma *StorageMarketActor_I) slashStorageDealCollateral(storageDealIDs []deal.DealID) {
	for i, dealID := range storageDealIDs {
		faultStorageDeal := sma.Deals()[dealID]
		// TODO remove locked funds and send slashed fund to TreasuryActor
		// TODO provider lose power for the FaultSet but not PledgeCollateral is slashed
	}
}
