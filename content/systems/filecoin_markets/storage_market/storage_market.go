package storage_market

const (
	StorageDealProposalNotFound StorageDealStatus = iota
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
