package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
)

func (provider *StorageProvider_I) OnNewStorageDealProposal(proposal deal.StorageDealProposal, payloadCID ipld.CID) {

	_, found := provider.ProposalStatus()[proposal.CID()]
	if found {
		// TODO: return error
		return
	}

	var shouldReject bool // specified by StorageProvider
	if shouldReject {
		provider.rejectStorageDealProposal(proposal)
		return
	}

	if provider.verifyStorageDealProposal(proposal, payloadCID) {
		provider.acceptStorageDealProposal(proposal)
	} else {
		provider.rejectStorageDealProposal(proposal)
		return
	}

}

func (provider *StorageProvider_I) signStorageDealProposal(proposal deal.StorageDealProposal) msg.Message {
	// TODO: construct StorageDeal Message
	var storageDealMessage msg.Message

	// TODO: notify StorageClient StorageDealSigned
	return storageDealMessage
}

func (provider *StorageProvider_I) publishStorageDealMessage(message msg.Message) deal.StorageDeal {
	// TODO: send message to StorageMarketActor.PublishStorageDeal and get back DealID
	var dealID deal.DealID
	var dealCID deal.DealCID

	storageDeal := &deal.StorageDeal_I{
		ProposalMessage_: message,
		ID_:              dealID,
	}

	provider.DealStatus()[dealCID] = StorageDealPublished

	// TODO: notify StorageClient StorageDealPublished
	return storageDeal
}

func (provider *StorageProvider_I) acceptStorageDealProposal(proposal deal.StorageDealProposal) {
	provider.ProposalStatus()[proposal.CID()] = StorageDealProposalAccepted
	// TODO: notify StorageClient StorageDealAccepted
}

func (provider *StorageProvider_I) rejectStorageDealProposal(proposal deal.StorageDealProposal) {
	provider.ProposalStatus()[proposal.CID()] = StorageDealProposalRejected
	// TODO: notify StorageClient StorageDealRejected
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
	return false
}

func (provider *StorageProvider_I) NotifyOfOnChainDealStatus(dealCID deal.DealCID, newStatus StorageDealStatus) {
	_, found := provider.DealStatus()[dealCID]
	if found {
		provider.DealStatus()[dealCID] = newStatus
	}
}

// the entire payload graph is now in local IPLD store
// TODO: integrate with Data Transfer
func (provider *StorageProvider_I) OnReceivingPayload(payloadCID ipld.CID) {
	// TODO: get proposalCID from local storage
	var proposalCID deal.ProposalCID

	_, found := provider.ProposalStatus()[proposalCID]
	if !found {
		// TODO: error here
	}

	// TODO: get client addr from libp2p
	// TODO: get proposal from local storage
	var proposal deal.StorageDealProposal
	isProposalVerified := provider.verifyStorageDealProposal(proposal, payloadCID)

	if !isProposalVerified {
		provider.rejectStorageDealProposal(proposal)
		return
	}

	// StorageProvider can decide what to do here
	provider.signStorageDealProposal(proposal)

}

func (provider *StorageProvider_I) OnStorageDealProposalQuery(proposalCID deal.ProposalCID) StorageDealStatus {
	proposalStatus, found := provider.ProposalStatus()[proposalCID]

	if found {
		return proposalStatus
	}

	return StorageDealProposalNotFound
}

func (provider *StorageProvider_I) OnStorageDealQuery(dealCID deal.DealCID) StorageDealStatus {
	dealStatus, found := provider.DealStatus()[dealCID]

	if found {
		return dealStatus
	}

	return StorageDealNotFound
}
