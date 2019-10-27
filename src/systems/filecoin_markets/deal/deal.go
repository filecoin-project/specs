package deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

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

func (p *StorageDealProposal_I) ClientBalanceRequirement() actor.TokenAmount {
	collateralPerEpoch := p.ClientCollateralPerEpoch() + p.ProviderCollateralPerEpoch()
	return actor.TokenAmount(uint64(collateralPerEpoch) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.ProviderCollateralPerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) CID() ProposalCID {
	// TODO: should be generated in codegen
	var cid ProposalCID
	return cid
}
