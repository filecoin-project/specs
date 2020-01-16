package indices

import (
	"math/big"

	abi "github.com/filecoin-project/specs/actors/abi"
	actor_util "github.com/filecoin-project/specs/actors/util"
)

var PARAM_FINISH = actor_util.PARAM_FINISH

// Data in Indices are populated at instantiation with data from the state tree
// Indices itself has no state tree or access to the runtime
// it is a passive data structure that allows for convenience access to network indices
// and pure functions in implementing economic policies given states
type Indices interface {
	Epoch() abi.ChainEpoch
	NetworkKPI() big.Int
	TotalNetworkSectorWeight() abi.SectorWeight
	TotalPledgeCollateral() abi.TokenAmount
	TotalNetworkEffectivePower() abi.StoragePower // power above minimum miner size
	TotalNetworkPower() abi.StoragePower          // total network power irrespective of meeting minimum miner size

	TotalMinedFIL() abi.TokenAmount
	TotalUnminedFIL() abi.TokenAmount
	TotalBurnedFIL() abi.TokenAmount
	LastEpochReward() abi.TokenAmount

	StorageDeal_DurationBounds(
		pieceSize abi.PieceSize,
		startEpoch abi.ChainEpoch,
	) (minDuration abi.ChainEpoch, maxDuration abi.ChainEpoch)
	StorageDeal_StoragePricePerEpochBounds(
		pieceSize abi.PieceSize,
		startEpoch abi.ChainEpoch,
		endEpoch abi.ChainEpoch,
	) (minPrice abi.TokenAmount, maxPrice abi.TokenAmount)
	StorageDeal_ProviderCollateralBounds(
		pieceSize abi.PieceSize,
		startEpoch abi.ChainEpoch,
		endEpoch abi.ChainEpoch,
	) (minProviderCollateral abi.TokenAmount, maxProviderCollateral abi.TokenAmount)
	StorageDeal_ClientCollateralBounds(
		pieceSize abi.PieceSize,
		startEpoch abi.ChainEpoch,
		endEpoch abi.ChainEpoch,
	) (minClientCollateral abi.TokenAmount, maxClientCollateral abi.TokenAmount)
	SectorWeight(
		sectorSize abi.SectorSize,
		startEpoch abi.ChainEpoch,
		endEpoch abi.ChainEpoch,
		dealWeight abi.DealWeight,
	) abi.SectorWeight
	PledgeCollateralReq(minerNominalPower abi.StoragePower) abi.TokenAmount
	SectorWeightProportion(minerActiveSectorWeight abi.SectorWeight) big.Int
	PledgeCollateralProportion(minerPledgeCollateral abi.TokenAmount) big.Int
	StoragePower(
		minerActiveSectorWeight abi.SectorWeight,
		minerInactiveSectorWeight abi.SectorWeight,
		minerPledgeCollateral abi.TokenAmount,
	) abi.StoragePower
	StoragePowerProportion(
		minerStoragePower abi.StoragePower,
	) big.Int
	CurrEpochBlockReward() abi.TokenAmount
	GetCurrBlockRewardRewardForMiner(
		minerStoragePower abi.StoragePower,
		minerPledgeCollateral abi.TokenAmount,
		// TODO extend or eliminate
	) abi.TokenAmount
	StoragePower_PledgeSlashForSectorTermination(
		storageWeightDesc actor_util.SectorStorageWeightDesc,
		terminationType actor_util.SectorTermination,
	) abi.TokenAmount
	StoragePower_PledgeSlashForSurprisePoStFailure(
		minerClaimedPower abi.StoragePower,
		numConsecutiveFailures int64,
	) abi.TokenAmount
	StorageMining_PreCommitDeposit(
		sectorSize abi.SectorSize,
		expirationEpoch abi.ChainEpoch,
	) abi.TokenAmount
	StorageMining_TemporaryFaultFee(
		storageWeightDescs []actor_util.SectorStorageWeightDesc,
		duration abi.ChainEpoch,
	) abi.TokenAmount
	NetworkTransactionFee(
		toActorCodeID abi.ActorCodeID,
		methodNum abi.MethodNum,
	) abi.TokenAmount
	GetCurrBlockRewardForMiner(
		minerStoragePower abi.StoragePower,
		minerPledgeCollateral abi.TokenAmount,
	) abi.TokenAmount
}

type IndicesImpl struct {
	// these fields are computed from StateTree upon construction
	// they are treated as globally available states
	Epoch                      abi.ChainEpoch
	NetworkKPI                 big.Int
	TotalNetworkSectorWeight   abi.SectorWeight
	TotalPledgeCollateral      abi.TokenAmount
	TotalNetworkEffectivePower abi.StoragePower // power above minimum miner size
	TotalNetworkPower          abi.StoragePower // total network power irrespective of meeting minimum miner size

	TotalMinedFIL   abi.TokenAmount
	TotalUnminedFIL abi.TokenAmount
	TotalBurnedFIL  abi.TokenAmount
	LastEpochReward abi.TokenAmount
}

func (inds *IndicesImpl) StorageDeal_DurationBounds(
	pieceSize abi.PieceSize,
	startEpoch abi.ChainEpoch,
) (minDuration abi.ChainEpoch, maxDuration abi.ChainEpoch) {

	// placeholder
	PARAM_FINISH()
	minDuration = abi.ChainEpoch(0)
	maxDuration = abi.ChainEpoch(1 << 20)
	return
}

func (inds *IndicesImpl) StorageDeal_StoragePricePerEpochBounds(
	pieceSize abi.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minPrice abi.TokenAmount, maxPrice abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StorageDeal_ProviderCollateralBounds(
	pieceSize abi.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minProviderCollateral abi.TokenAmount, maxProviderCollateral abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StorageDeal_ClientCollateralBounds(
	pieceSize abi.PieceSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
) (minClientCollateral abi.TokenAmount, maxClientCollateral abi.TokenAmount) {

	// placeholder
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) SectorWeight(
	sectorSize abi.SectorSize,
	startEpoch abi.ChainEpoch,
	endEpoch abi.ChainEpoch,
	dealWeight abi.DealWeight,
) abi.SectorWeight {
	// for every sector, given its size, start, end, and deals within the sector
	// assign sector power for the duration of its lifetime
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) PledgeCollateralReq(minerNominalPower abi.StoragePower) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) SectorWeightProportion(minerActiveSectorWeight abi.SectorWeight) big.Int {
	// return proportion of SectorWeight for miner
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) PledgeCollateralProportion(minerPledgeCollateral abi.TokenAmount) big.Int {
	// return proportion of Pledge Collateral for miner
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StoragePower(
	minerActiveSectorWeight abi.SectorWeight,
	minerInactiveSectorWeight abi.SectorWeight,
	minerPledgeCollateral abi.TokenAmount,
) abi.StoragePower {
	// return StoragePower based on inputs
	// StoragePower for miner = func(ActiveSectorWeight for miner, PledgeCollateral for miner, global indices)
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StoragePowerProportion(
	minerStoragePower abi.StoragePower,
) big.Int {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) CurrEpochBlockReward() abi.TokenAmount {
	// total block reward allocated for CurrEpoch
	// each expected winner get an equal share of this reward
	// computed as a function of NetworkKPI, LastEpochReward, TotalUnmminedFIL, etc
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) GetCurrBlockRewardRewardForMiner(
	minerStoragePower abi.StoragePower,
	minerPledgeCollateral abi.TokenAmount,
	// TODO extend or eliminate
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

// TerminationFault
func (inds *IndicesImpl) StoragePower_PledgeSlashForSectorTermination(
	storageWeightDesc actor_util.SectorStorageWeightDesc,
	terminationType actor_util.SectorTermination,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

// DetectedFault
func (inds *IndicesImpl) StoragePower_PledgeSlashForSurprisePoStFailure(
	minerClaimedPower abi.StoragePower,
	numConsecutiveFailures int64,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StorageMining_PreCommitDeposit(
	sectorSize abi.SectorSize,
	expirationEpoch abi.ChainEpoch,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) StorageMining_TemporaryFaultFee(
	storageWeightDescs []actor_util.SectorStorageWeightDesc,
	duration abi.ChainEpoch,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) NetworkTransactionFee(
	toActorCodeID abi.ActorCodeID,
	methodNum abi.MethodNum,
) abi.TokenAmount {
	PARAM_FINISH()
	panic("")
}

func (inds *IndicesImpl) GetCurrBlockRewardForMiner(
	minerStoragePower abi.StoragePower,
	minerPledgeCollateral abi.TokenAmount,
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

func StorageDeal_ProviderInitTimedOutSlashAmount(providerCollateral abi.TokenAmount) abi.TokenAmount {
	// placeholder
	PARAM_FINISH()
	return providerCollateral
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

func StoragePower_SurprisePoStMaxConsecutiveFailures() int64 {
	PARAM_FINISH()
	panic("")
}

func StorageMining_DeclaredFaultEffectiveDelay() abi.ChainEpoch {
	PARAM_FINISH()
	panic("")
}
