package storage_market

import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"

type StorageDealStatus int

const (
	StorageDealProposalNotFound StorageDealStatus = 1 + iota
	StorageDealProposalRejected
	StorageDealProposalAccepted
	StorageDealProposalSigned
	StorageDealPublished
	StorageDealCommitted
	StorageDealActive
	StorageDealFailing
	StorageDealRecovering
	StorageDealExpired
	StorageDealNotFound
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
