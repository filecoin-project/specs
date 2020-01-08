package storage_deal

import actors "github.com/filecoin-project/specs/actors"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

import util "github.com/filecoin-project/specs/util"

var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

func (d *StorageDeal_I) Proposal() StorageDealProposal {
	// extract from d.ProposalMessage
	TODO()
	panic("")
}

func (p *StorageDealProposal_I) Duration() block.ChainEpoch {
	return (p.EndEpoch() - p.StartEpoch())
}

func (p *StorageDealProposal_I) TotalStorageFee() actors.TokenAmount {
	return actors.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ClientBalanceRequirement() actors.TokenAmount {
	return (p.ClientCollateral() + p.TotalStorageFee())
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actors.TokenAmount {
	return p.ProviderCollateral()
}
