package indices

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
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

func (inds *Indices_I) BlockReward_SectorPower(
	sectorSize sector.SectorSize,
	startEpoch block.ChainEpoch,
	endEpoch block.ChainEpoch,
	dealWeight deal.DealWeight,
) block.StoragePower {
	// for every sector, given its size, start, end, and deals within the sector
	// assign sector power for the duration of its lifetime
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_PledgeCollateralReq(
	minerActiveStoragePower block.StoragePower,
	minerInactiveStoragePower block.StoragePower,
	minerPledgeCollateral actor.TokenAmount,
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_StoragePowerProportion(minerActiveStoragePower block.StoragePower) util.BigInt {
	// return proportion of StoragePower for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_CollateralPowerProportion(minerPledgeCollateral actor.TokenAmount) util.BigInt {
	// return proportion of Pledge Collateral for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_EarningPowerProportion(
	minerActiveStoragePower block.StoragePower,
	minerPledgeCollateral actor.TokenAmount,
) util.BigInt {
	// Earning Power for miner = func(proportion of StoragePower for miner, proportion of CollateralPower for miner)
	// EarningPowerProportion is a normalized proportion of TotalEarningPower
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_CurrEpochReward() actor.TokenAmount {
	// total block reward allocated for CurrEpoch
	// each expected winner get a share of this reward
	// computed as a function of NetworkKPI, LastEpochReward, TotalUnmminedFIL, etc
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_GetCurrRewardForMiner(
	minerActiveStoragePower block.StoragePower,
	minerPledgeCollateral actor.TokenAmount,
	// TODO extend
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) BlockReward_GetPledgeSlashForStorageFault(
	affectedPower block.StoragePower,
	newActivePower block.StoragePower,
	newInactivePower block.StoragePower,
	currPledge actor.TokenAmount,
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageMining_PreCommitDeposit(
	sectorSize sector.SectorSize,
	expirationEpoch block.ChainEpoch,
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageMining_TemporaryFaultFee(
	storageWeightDescs []actor_util.SectorStorageWeightDesc,
	duration block.ChainEpoch,
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}
