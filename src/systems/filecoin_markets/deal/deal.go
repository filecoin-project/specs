package deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

const MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1) // TODO: Placeholder
const MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH = actor.TokenAmount(1)   // TODO: Placeholder

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

func (p *StorageDealProposal_I) ClientBalanceRequirement() actor.TokenAmount {
	balanceRequirementPerEpoch := p.ClientCollateralPerEpoch() + p.StoragePricePerEpoch()
	return actor.TokenAmount(uint64(balanceRequirementPerEpoch) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.ProviderCollateralPerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) CID() ProposalCID {
	// TODO: should be generated in codegen
	var cid ProposalCID
	return cid
}
