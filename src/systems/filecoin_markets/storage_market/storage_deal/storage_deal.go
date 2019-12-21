package storage_deal

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
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

func (p *StorageDealProposal_I) TotalStorageFee() actor.TokenAmount {
	return actor.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ClientBalanceRequirement() actor.TokenAmount {
	return (p.ClientCollateral() + p.TotalStorageFee())
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() actor.TokenAmount {
	return p.ProviderCollateral()
}
