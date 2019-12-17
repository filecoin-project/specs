package storage_market

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
