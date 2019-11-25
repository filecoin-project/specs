package deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

const MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1) // TODO: Placeholder
const MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1)   // TODO: Placeholder
const MIN_DEAL_DURATION = block.ChainEpoch(0)                       // TODO: Placeholder
const MIN_DEAL_PRICE = actor.TokenAmount(0)                         // TODO: Placeholder

func (d *StorageDeal_I) Proposal() StorageDealProposal {
	// extract from d.ProposalMessage
	var proposal StorageDealProposal
	return proposal
}

func (d *StorageDeal_I) CID() DealCID {
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
	return (p.TotalClientCollateral() + p.TotalStorageFee())
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.ProviderCollateralPerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) CID() ProposalCID {
	var cid ProposalCID
	return cid
}

type StorageDealSlashAction = int

const (
	SlashTerminatedFaults StorageDealSlashAction = 0
)

func (amt *DealExpirationAMT_I) Size() int {
	return 0
}

func (amt *DealExpirationAMT_I) Add(key block.ChainEpoch, value DealExpirationValue) {
	// helper function to add entry into the AMT
}

func (amt *DealExpirationAMT_I) ActiveDealIDs() []DealID {
	ret := make([]DealID, 0)
	return ret
}

// return last item in the expiration amt
func (q *DealExpirationAMT_I) LastDealExpiration() block.ChainEpoch {
	ret := block.ChainEpoch(0)
	return ret
}

// return deal IDs expiring in epoch range
func (q *DealExpirationAMT_I) ExpiredDealsInRange(start block.ChainEpoch, end block.ChainEpoch) []DealExpirationValue {
	ret := make([]DealExpirationValue, 0)
	return ret
}
