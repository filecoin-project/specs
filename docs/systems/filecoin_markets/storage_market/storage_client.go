package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
)

func (client *StorageClient_I) generatePieceCID(payloadCID ipld.CID) piece.PieceCID {
	panic("TODO")
	var pieceCID piece.PieceCID
	return pieceCID
}

func (client *StorageClient_I) PullPayload(payloadCID ipld.CID) {
	panic("TODO")
}

func (client *StorageClient_I) NotifyOfStorageDealProposalStatus(pieceCID piece.PieceCID, status StorageDealProposalStatus) {
	panic("TODO")
}

func (client *StorageClient_I) NotifyOfStorageDealStatus(dealCID deal.DealCID, status StorageDealStatus) {
	panic("TODO")
}
