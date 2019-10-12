package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (sma *StorageMarketActor_I) WithdrawBalance(balance actor.TokenAmount) {
	var msgSender addr.Address // TODO replace this

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

	newAvailableBalance := senderBalance.Available() - balance
	sma.Balances()[msgSender] = &StorageParticipantBalance_I{
		Locked_:    senderBalance.Locked(),
		Available_: newAvailableBalance,
	}

	// TODO send funds to msgSender
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
		newAvailableBalance := senderBalance.Available() + balance
		sma.Balances()[msgSender] = &StorageParticipantBalance_I{
			Locked_:    senderBalance.Locked(),
			Available_: newAvailableBalance,
		}
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
	// for _, dealID := range storageDealIDs {
	// 	faultStorageDeal := sma.Deals()[dealID]
	// TODO remove locked funds and send slashed fund to TreasuryActor
	// TODO provider lose power for the FaultSet but not PledgeCollateral
	// }
	panic("TODO")
}
