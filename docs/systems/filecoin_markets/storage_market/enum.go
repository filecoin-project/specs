package storage_market

type StorageDealStatus int

const (
	StorageDealNotFound StorageDealStatus = 0
	StorageDealRejected StorageDealStatus = 1
	StorageDealProposed StorageDealStatus = 2
	StorageDealStaged   StorageDealStatus = 3
	StorageDealActive   StorageDealStatus = 4
)

type PublishStorageDealResponse int

const (
	PublishStorageDealError   PublishStorageDealResponse = 0
	PublishStorageDealSuccess PublishStorageDealResponse = 1
)
