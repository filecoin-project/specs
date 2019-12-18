package indices

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	util "github.com/filecoin-project/specs/util"
)

var PARAM_FINISH = util.PARAM_FINISH

func StorageDeal_ProviderInitTimedOutSlashAmount(deal deal.OnChainDeal) actor.TokenAmount {
	// placeholder
	PARAM_FINISH()
	return deal.Deal().Proposal().ProviderBalanceRequirement()
}

func (inds *Indices_I) StorageDeal_DurationBounds(
	pieceSize piece.PieceSize,
	startEpoch block.ChainEpoch,
) (minDuration block.ChainEpoch, maxDuration block.ChainEpoch) {

	// placeholder
	PARAM_FINISH()
	minDuration = block.ChainEpoch(0)
	maxDuration = block.ChainEpoch(1 << 20)
	return
}

func (inds *Indices_I) StorageDeal_StoragePricePerEpochBounds(
	pieceSize piece.PieceSize,
	startEpoch block.ChainEpoch,
	endEpoch block.ChainEpoch,
) (minPrice actor.TokenAmount, maxPrice actor.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageDeal_ProviderCollateralBounds(
	pieceSize piece.PieceSize,
	startEpoch block.ChainEpoch,
	endEpoch block.ChainEpoch,
) (minProviderCollateral actor.TokenAmount, maxProviderCollateral actor.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageDeal_ClientCollateralBounds(
	pieceSize piece.PieceSize,
	startEpoch block.ChainEpoch,
	endEpoch block.ChainEpoch,
) (minClientCollateral actor.TokenAmount, maxClientCollateral actor.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}
