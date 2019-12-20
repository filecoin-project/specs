package storage_mining

import (
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert

func (st *StorageMinerActorState_I) _getSectorOnChainInfo(sectorNo sector.SectorNumber) (info SectorOnChainInfo, ok bool) {
	sectorInfo, found := st.Sectors()[sectorNo]
	if !found {
		return nil, false
	}
	return sectorInfo, true
}

func (st *StorageMinerActorState_I) _getSectorDealIDsAssert(sectorNo sector.SectorNumber) deal.DealIDs {
	sectorInfo, found := st._getSectorOnChainInfo(sectorNo)
	Assert(found)
	return sectorInfo.Info().DealIDs()
}

func SectorsAMT_Empty() SectorsAMT {
	IMPL_FINISH()
	panic("")
}

func (st *StorageMinerActorState_I) _getStorageWeightForSector(sectorNumber sector.SectorNumber) SectorStorageWeight {
	sectorInfo, found := st.Sectors()[sectorNumber]
	Assert(found)

	return &actor_util.SectorStorageWeight_I{
		SectorSize_: st.Info().SectorSize(),
		DealWeight_: sectorInfo.DealWeight(),
		Duration_:   sectorInfo.Info().Expiration() - sectorInfo.ProveCommitEpoch(),
	}
}

func (st *StorageMinerActorState_I) _getStorageWeightsForSectors(sectorNumbers []sector.SectorNumber) []SectorStorageWeight {
	ret := []SectorStorageWeight{}
	for _, sectorNumber := range sectorNumbers {
		ret = append(ret, st._getStorageWeightForSector(sectorNumber))
	}
	return ret
}

func MinerPoStState_New_OK(lastSuccessfulPoSt block.ChainEpoch) MinerPoStState {
	return MinerPoStState_Make_OK(&MinerPoStState_OK_I{
		LastSuccessfulPoSt_: lastSuccessfulPoSt,
	})
}

func MinerPoStState_New_Challenged(surpriseChallengeEpoch block.ChainEpoch, numConsecutiveFailures int) MinerPoStState {
	return MinerPoStState_Make_Challenged(&MinerPoStState_Challenged_I{
		SurpriseChallengeEpoch_: surpriseChallengeEpoch,
		NumConsecutiveFailures_: numConsecutiveFailures,
	})
}

func MinerPoStState_New_Failing(numConsecutiveFailures int) MinerPoStState {
	return MinerPoStState_Make_Failing(&MinerPoStState_Failing_I{
		NumConsecutiveFailures_: numConsecutiveFailures,
	})
}

func (x *SectorOnChainInfo_I) Is_DeclaredFault() bool {
	ret := (x.State() == SectorState_DeclaredFault)
	Assert(ret == (x.DeclaredFaultEpoch() == block.ChainEpoch_None))
	Assert(ret == (x.DeclaredFaultDuration() == block.ChainEpoch_None))
	return ret
}

func (x *SectorOnChainInfo_I) EffectiveFaultBeginEpoch() block.ChainEpoch {
	Assert(x.Is_DeclaredFault())
	return x.DeclaredFaultEpoch() + DECLARED_FAULT_EFFECTIVE_DELAY
}

func (x *SectorOnChainInfo_I) EffectiveFaultEndEpoch() block.ChainEpoch {
	Assert(x.Is_DeclaredFault())
	return x.EffectiveFaultBeginEpoch() + x.DeclaredFaultDuration()
}
