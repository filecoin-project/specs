package storage_deal

import abi "github.com/filecoin-project/specs/actors/abi"

import util "github.com/filecoin-project/specs/util"

var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

func (p *StorageDealProposal_I) Duration() abi.ChainEpoch {
	return (p.EndEpoch() - p.StartEpoch())
}

func (p *StorageDealProposal_I) TotalStorageFee() abi.TokenAmount {
	return abi.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(p.Duration()))
}

func (p *StorageDealProposal_I) ClientBalanceRequirement() abi.TokenAmount {
	return (p.ClientCollateral() + p.TotalStorageFee())
}

func (p *StorageDealProposal_I) ProviderBalanceRequirement() abi.TokenAmount {
	return p.ProviderCollateral()
}
