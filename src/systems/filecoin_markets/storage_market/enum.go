package storage_market

type StorageDealProposalStatus int

const (
	StorageDealProposalNotFound StorageDealProposalStatus = 0
	StorageDealProposalRejected StorageDealProposalStatus = 1
	StorageDealProposalAccepted StorageDealProposalStatus = 2
	StorageDealProposalSigned   StorageDealProposalStatus = 3
)

type StorageDealStatus int

const (
	StorageDealNotFound   StorageDealStatus = 0
	StorageDealCreated    StorageDealStatus = 1
	StorageDealPublished  StorageDealStatus = 2
	StorageDealCommitted  StorageDealStatus = 3
	StorageDealActive     StorageDealStatus = 4
	StorageDealFailing    StorageDealStatus = 5
	StorageDealRecovering StorageDealStatus = 6
	StorageDealExpired    StorageDealStatus = 7
)

type PublishStorageDealResponse int

const (
	PublishStorageDealError   PublishStorageDealResponse = 0
	PublishStorageDealSuccess PublishStorageDealResponse = 1
)
