package storage_mining

import (
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
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

func (st *StorageMinerActorState_I) _getStorageWeightDescForSector(sectorNumber sector.SectorNumber) SectorStorageWeightDesc {
	sectorInfo, found := st.Sectors()[sectorNumber]
	Assert(found)

	return &actor_util.SectorStorageWeightDesc_I{
		SectorSize_: st.Info().SectorSize(),
		DealWeight_: sectorInfo.DealWeight(),
		Duration_:   sectorInfo.Info().Expiration() - sectorInfo.ProveCommitEpoch(),
	}
}

func (st *StorageMinerActorState_I) _getStorageWeightDescsForSectors(sectorNumbers []sector.SectorNumber) []SectorStorageWeightDesc {
	ret := []SectorStorageWeightDesc{}
	for _, sectorNumber := range sectorNumbers {
		ret = append(ret, st._getStorageWeightDescForSector(sectorNumber))
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

func (x *SectorOnChainInfo_I) Is_TemporaryFault() bool {
	ret := (x.State() == SectorState_TemporaryFault)
	if ret {
		Assert(x.DeclaredFaultEpoch() != block.ChainEpoch_None)
		Assert(x.DeclaredFaultDuration() != block.ChainEpoch_None)
	}
	return ret
}

func (x *SectorOnChainInfo_I) EffectiveFaultBeginEpoch() block.ChainEpoch {
	Assert(x.Is_TemporaryFault())
	return x.DeclaredFaultEpoch() + DECLARED_FAULT_EFFECTIVE_DELAY
}

func (x *SectorOnChainInfo_I) EffectiveFaultEndEpoch() block.ChainEpoch {
	Assert(x.Is_TemporaryFault())
	return x.EffectiveFaultBeginEpoch() + x.DeclaredFaultDuration()
}

func MinerInfo_New(
	ownerAddr addr.Address, workerAddr addr.Address, sectorSize sector.SectorSize, peerId libp2p.PeerID) MinerInfo {

	ret := &MinerInfo_I{
		Owner_:      ownerAddr,
		Worker_:     workerAddr,
		PeerId_:     peerId,
		SectorSize_: sectorSize,
	}

	TODO() // TODO: determine how to generate/validate VRF key and initialize other fields

	return ret
}
