package indices

import (
	actor_util "github.com/filecoin-project/specs/actors/util"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

var PARAM_FINISH = util.PARAM_FINISH

func Indices_FromStateTree(tree st.StateTree) Indices {
	PARAM_FINISH()
	panic("")
}

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

func (inds *Indices_I) SectorWeight(
	sectorSize sector.SectorSize,
	startEpoch block.ChainEpoch,
	endEpoch block.ChainEpoch,
	dealWeight deal.DealWeight,
) block.SectorWeight {
	// for every sector, given its size, start, end, and deals within the sector
	// assign sector power for the duration of its lifetime
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) PledgeCollateralReq(minerNominalPower block.StoragePower) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) SectorWeightProportion(minerActiveSectorWeight block.SectorWeight) util.BigInt {
	// return proportion of SectorWeight for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) PledgeCollateralProportion(minerPledgeCollateral actor.TokenAmount) util.BigInt {
	// return proportion of Pledge Collateral for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StoragePower(
	minerActiveSectorWeight block.SectorWeight,
	minerInactiveSectorWeight block.SectorWeight,
	minerPledgeCollateral actor.TokenAmount,
) block.StoragePower {
	// return StoragePower based on inputs
	// StoragePower for miner = func(ActiveSectorWeight for miner, PledgeCollateral for miner, global indices)
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StoragePowerProportion(
	minerStoragePower block.StoragePower,
) util.BigInt {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) CurrEpochBlockReward() actor.TokenAmount {
	// total block reward allocated for CurrEpoch
	// each expected winner get an equal share of this reward
	// computed as a function of NetworkKPI, LastEpochReward, TotalUnmminedFIL, etc
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) GetCurrBlockRewardRewardForMiner(
	minerStoragePower block.StoragePower,
	minerPledgeCollateral actor.TokenAmount,
	// TODO extend or eliminate
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) GetPledgeSlashForStorageFault(
	affectedPower block.StoragePower,
	newActiveSectorWeight block.SectorWeight,
	newInactiveSectorWeight block.SectorWeight,
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

func (inds *Indices_I) NetworkTransactionFee(
	toActorCodeID actor.CodeID,
	methodNum actor.MethodNum,
) actor.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func ConsensusPowerForStorageWeight(
	storageWeightDesc actor_util.SectorStorageWeightDesc,
) block.StoragePower {
	PARAM_FINISH()
	panic("")
}

func StoragePower_ConsensusMinMinerPower() block.StoragePower {
	PARAM_FINISH()
	panic("")
}

func StorageMining_PoStNoChallengePeriod() block.ChainEpoch {
	PARAM_FINISH()
	panic("")
}

func StorageMining_SurprisePoStProvingPeriod() block.ChainEpoch {
	PARAM_FINISH()
	panic("")
}

func StoragePower_SurprisePoStMaxConsecutiveFailures() int {
	PARAM_FINISH()
	panic("")
}

func StorageMining_DeclaredFaultEffectiveDelay() block.ChainEpoch {
	PARAM_FINISH()
	panic("")
}
