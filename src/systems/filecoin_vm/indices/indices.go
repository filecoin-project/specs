package indices

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	actor_util "github.com/filecoin-project/specs/actors/util"
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

func StorageDeal_ProviderInitTimedOutSlashAmount(deal deal.OnChainDeal) abi.TokenAmount {
	// placeholder
	PARAM_FINISH()
	return deal.Deal().Proposal().ProviderBalanceRequirement()
}

func (inds *Indices_I) StorageDeal_DurationBounds(
	pieceSize piece.PieceSize,
	startEpoch abi.ChainEpoch,
) (minDuration abi.ChainEpoch, maxDuration abi.ChainEpoch) {

	// placeholder
	PARAM_FINISH()
	minDuration = abi.ChainEpoch(0)
	maxDuration = abi.ChainEpoch(1 << 20)
	return
}

func (inds *Indices_I) StorageDeal_StoragePricePerEpochBounds(
	pieceSize piece.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minPrice abi.TokenAmount, maxPrice abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageDeal_ProviderCollateralBounds(
	pieceSize piece.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minProviderCollateral abi.TokenAmount, maxProviderCollateral abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageDeal_ClientCollateralBounds(
	pieceSize piece.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minClientCollateral abi.TokenAmount, maxClientCollateral abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) SectorWeight(
	sectorSize sector.SectorSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
	dealWeight deal.DealWeight,
) abi.SectorWeight {
	// for every sector, given its size, start, end, and deals within the sector
	// assign sector power for the duration of its lifetime
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) PledgeCollateralReq(minerNominalPower abi.StoragePower) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) SectorWeightProportion(minerActiveSectorWeight abi.SectorWeight) util.BigInt {
	// return proportion of SectorWeight for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) PledgeCollateralProportion(minerPledgeCollateral abi.TokenAmount) util.BigInt {
	// return proportion of Pledge Collateral for miner
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StoragePower(
	minerActiveSectorWeight abi.SectorWeight,
	minerInactiveSectorWeight abi.SectorWeight,
	minerPledgeCollateral abi.TokenAmount,
) abi.StoragePower {
	// return StoragePower based on inputs
	// StoragePower for miner = func(ActiveSectorWeight for miner, PledgeCollateral for miner, global indices)
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StoragePowerProportion(
	minerStoragePower abi.StoragePower,
) util.BigInt {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) CurrEpochBlockReward() abi.TokenAmount {
	// total block reward allocated for CurrEpoch
	// each expected winner get an equal share of this reward
	// computed as a function of NetworkKPI, LastEpochReward, TotalUnmminedFIL, etc
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) GetCurrBlockRewardRewardForMiner(
	minerStoragePower abi.StoragePower,
	minerPledgeCollateral abi.TokenAmount,
	// TODO extend or eliminate
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

// TerminationFault
func (inds *Indices_I) StoragePower_PledgeSlashForSectorTermination(
	storageWeightDesc actor_util.SectorStorageWeightDesc,
	terminationType actor_util.SectorTerminationType,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

// DetectedFault
func (inds *Indices_I) StoragePower_PledgeSlashForSurprisePoStFailure(
	minerClaimedPower abi.StoragePower,
	numConsecutiveFailures int,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageMining_PreCommitDeposit(
	sectorSize sector.SectorSize,
	expirationEpoch abi.ChainEpoch,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) StorageMining_TemporaryFaultFee(
	storageWeightDescs []actor_util.SectorStorageWeightDesc,
	duration abi.ChainEpoch,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *Indices_I) NetworkTransactionFee(
	toActorCodeID actor.CodeID,
	methodNum abi.MethodNum,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func ConsensusPowerForStorageWeight(
	storageWeightDesc actor_util.SectorStorageWeightDesc,
) abi.StoragePower {
	PARAM_FINISH()
	panic("")
}

func StoragePower_ConsensusMinMinerPower() abi.StoragePower {
	PARAM_FINISH()
	panic("")
}

func StorageMining_PoStNoChallengePeriod() abi.ChainEpoch {
	PARAM_FINISH()
	panic("")
}

func StorageMining_SurprisePoStProvingPeriod() abi.ChainEpoch {
	PARAM_FINISH()
	panic("")
}

func StoragePower_SurprisePoStMaxConsecutiveFailures() int {
	PARAM_FINISH()
	panic("")
}

func StorageMining_DeclaredFaultEffectiveDelay() abi.ChainEpoch {
	PARAM_FINISH()
	panic("")
}
