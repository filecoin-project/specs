package deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

import util "github.com/filecoin-project/specs/util"

var IMPL_FINISH = util.IMPL_FINISH

var MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount_Placeholder()
var MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount_Placeholder()
var MIN_DEAL_DURATION = block.ChainEpoch(0) // TODO: Placeholder
var MIN_DEAL_PRICE = actor.TokenAmount_Placeholder()

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
