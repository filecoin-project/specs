package storage_market

import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

type StorageDealStatus int

const (
	StorageDealProposalNotFound StorageDealStatus = 0
	StorageDealProposalRejected StorageDealStatus = 1
	StorageDealProposalAccepted StorageDealStatus = 2
	StorageDealProposalSigned   StorageDealStatus = 3
	StorageDealPublished        StorageDealStatus = 4
	StorageDealCommitted        StorageDealStatus = 5
	StorageDealActive           StorageDealStatus = 6
	StorageDealFailing          StorageDealStatus = 7
	StorageDealRecovering       StorageDealStatus = 8
	StorageDealExpired          StorageDealStatus = 9
	StorageDealNotFound         StorageDealStatus = 10
)

type PublishStorageDealResponse struct {
	StatusCode uint8 // 0 for failure 1 for success
	DealID     deal.DealID
}

func PublishStorageDealError() PublishStorageDealResponse {
	return PublishStorageDealResponse{
		StatusCode: 0,
		DealID:     deal.DealID(0),
	}
}

func PublishStorageDealSuccess(dealID deal.DealID) PublishStorageDealResponse {
	return PublishStorageDealResponse{
		StatusCode: 1,
		DealID:     dealID,
	}
}
