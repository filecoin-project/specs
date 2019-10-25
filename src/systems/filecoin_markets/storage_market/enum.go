package storage_market

type StorageDealStatus int

const (
	StorageDealProposalNotFound  StorageDealStatus = 0
	StorageDealProposalRejected  StorageDealStatus = 1
	StorageDealProposalAccepted  StorageDealStatus = 2
	StorageDealProposalSigned    StorageDealStatus = 3
	StorageDealPublished         StorageDealStatus = 4
	StorageDealCommitted         StorageDealStatus = 5
	StorageDealActive            StorageDealStatus = 6
	StorageDealFailing           StorageDealStatus = 7
	StorageDealRecovering        StorageDealStatus = 8
	StorageDealExpired           StorageDealStatus = 9
	StorageDealNotFound          StorageDealStatus = 10
)

type PublishStorageDealResponse int

const (
	PublishStorageDealError   PublishStorageDealResponse = 0
	PublishStorageDealSuccess PublishStorageDealResponse = 1
)
