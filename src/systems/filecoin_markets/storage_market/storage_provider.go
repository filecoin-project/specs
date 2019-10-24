package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
)

func (provider *StorageProvider_I) HandleStorageDealProposal(proposal deal.StorageDealProposal, payloadCID ipld.CID) {

	_, found := provider.ProposalStatus()[payloadCID]
	if found {
		// TODO: may want to throw here
		return
	}

	var shouldReject bool // specified by StorageProvider
	if shouldReject {
		provider.rejectStorageDealProposal(payloadCID)
		return
	}

	if provider.verifyStorageDealProposal(proposal, payloadCID) {
		provider.acceptStorageDealProposal(payloadCID)
		// TODO: PullPayload(proposal.PieceCID) over libp2p
	} else {
		provider.rejectStorageDealProposal(payloadCID)
		return
	}

}

func (provider *StorageProvider_I) OnReceivingPayload(payloadCID ipld.CID) {
	_, found := provider.ProposalStatus()[payloadCID]
	if !found {
		// TODO: error here
	}

	// TODO: may want to check on proposalStatus here

	// TODO: do we need to know storage client somewhere or can get this from libp2p
	// TODO: get proposal
	var proposal deal.StorageDealProposal
	isProposalVerified := provider.verifyStorageDealProposal(proposal, payloadCID)

	if !isProposalVerified {
		provider.rejectStorageDealProposal(payloadCID)
		return
	}

	provider.signStorageDealProposal(proposal)

}

func (provider *StorageProvider_I) signStorageDealProposal(proposal deal.StorageDealProposal) deal.StorageDeal {
	// TODO add signature to the proposal
	// TODO notify StorageClient StorageDealProposalSigned
	// TODO notify StorageClient StorageDealCreated
	panic("TODO")
}

func (provider *StorageProvider_I) acceptStorageDealProposal(payloadCID ipld.CID) {
	provider.ProposalStatus()[payloadCID] = StorageDealProposalAccepted
	// TODO: notify StorageClient StorageDealProposalAccepted
}

func (provider *StorageProvider_I) rejectStorageDealProposal(payloadCID ipld.CID) {
	provider.ProposalStatus()[payloadCID] = StorageDealProposalRejected
	// TODO: notify StorageClient StorageDealProposalRejected
}

func (provider *StorageProvider_I) verifyStorageDealProposal(proposal deal.StorageDealProposal, payloadCID ipld.CID) bool {
	// TODO make call to StorageMarketActor
	// balance, found := StorageMarketActor.Balances()[address]

	// if !found {
	// 	return false
	// }

	// if balance < price {
	// 	return false
	// }

	isPieceCIDVerified := provider.verifyPieceCID(proposal.PieceCID(), payloadCID)
	if !isPieceCIDVerified {
		// TODO: error out
		panic("TODO")
	}

	// TODO Check on Signature
	// return true
	panic("TODO")
}

func (provider *StorageProvider_I) verifyPieceCID(pieceCID piece.PieceCID, payloadCID ipld.CID) bool {
	panic("TODO")
	return true
}



func (provider *StorageProvider_I) NotifyOfOnChainDealStatus(dealID deal.DealID, newStatus StorageDealStatus) {
	_, found := provider.DealStatus()[dealID]
	if found {
		if newStatus == StorageDealCreated || newStatus == StorageDealNotFound {
			// TODO error out
			panic("invalid onchain deal status")
		} else {
			provider.DealStatus()[dealID] = newStatus
		}
	}
}

func (provider *StorageProvider_I) HandleStorageDealProposalQuery(payloadCID ipld.CID) StorageDealProposalStatus {
	proposalStatus, found := provider.ProposalStatus()[payloadCID]

	if found {
		return proposalStatus
	}

	return StorageDealProposalNotFound
}

func (provider *StorageProvider_I) HandleStorageDealQuery(dealID deal.DealID) StorageDealStatus {
	dealStatus, found := provider.DealStatus()[dealID]

	if found {
		return dealStatus
	}

	return StorageDealNotFound
}
