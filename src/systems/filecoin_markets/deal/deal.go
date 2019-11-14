package deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

const MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1) // TODO: Placeholder
const MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1)   // TODO: Placeholder
const MIN_DEAL_DURATION = block.ChainEpoch(0)                       // TODO: Placeholder
const MIN_DEAL_PRICE = actor.TokenAmount(0)                         // TODO: Placeholder

func (d *StorageDeal_I) Proposal() StorageDealProposal {
	// TODO: extract from d.ProposalMessage
	var proposal StorageDealProposal
	return proposal
}

func (d *StorageDeal_I) CID() DealCID {
	// TODO: should be generated in codegen
	var cid DealCID
	return cid
}

func (p *StorageDealProposal_I) Duration() block.ChainEpoch {
	return (p.EndEpoch() - p.StartEpoch())
}

func (p *StorageDealProposal_I) TotalStorageFee() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) TotalClientCollateral() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.ClientCollateralPerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ClientBalanceRequirement() actor.TokenAmount {
	return (p.TotalClientCollateral() + p.TotalClientCollateral())
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.ProviderCollateralPerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) CID() ProposalCID {
	// TODO: should be generated in codegen
	var cid ProposalCID
	return cid
}

// move storage fee from locked to unlocked
func (d *StorageDealTally_I) UnlockStorageFee(fee actor.TokenAmount) bool {
	// if d.Deal().Proposal().TotalStorageFee() < fee {
	// cannot unlock more than total
	// return false
	// }

	if d.LockedStorageFee() < fee {
		// cannot unlock more than locked
		return false
	}

	d.LockedStorageFee_ -= fee
	d.UnlockedStorageFee_ += fee
	return true
}

func (q *DealExpirationQueue_I) Size() int {
	// TODO
	return 0
}

// return all dealIds in the expiration queue
func (q *DealExpirationQueue_I) ActiveDealIDs() CompactDealSet {
	// TODO
	ret := CompactDealSet(make([]byte, q.Size()))
	return ret
}

// return last item in the expiration queue
func (q *DealExpirationQueue_I) LastDealExpiration() block.ChainEpoch {
	// TODO
	ret := block.ChainEpoch(0)
	return ret
}
