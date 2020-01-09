package storage_miner

import (
	actors "github.com/filecoin-project/specs/actors"
	actor_util "github.com/filecoin-project/specs/actors/util"
	libp2p "github.com/filecoin-project/specs/libraries/libp2p"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	indices "github.com/filecoin-project/specs/systems/filecoin_vm/indices"
)

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

func SectorNumberSetHAMT_Empty() SectorNumberSetHAMT {
	IMPL_FINISH()
	panic("")
}

func (st *StorageMinerActorState_I) GetStorageWeightDescForSectorMaybe(sectorNumber sector.SectorNumber) (ret SectorStorageWeightDesc, ok bool) {
	sectorInfo, found := st.Sectors()[sectorNumber]
	if !found {
		ret = nil
		ok = false
		return
	}

	ret = &actor_util.SectorStorageWeightDesc_I{
		SectorSize_: st.Info().SectorSize(),
		DealWeight_: sectorInfo.DealWeight(),
		Duration_:   sectorInfo.Info().Expiration() - sectorInfo.ActivationEpoch(),
	}
	ok = true
	return
}

func (st *StorageMinerActorState_I) _getStorageWeightDescForSector(sectorNumber sector.SectorNumber) SectorStorageWeightDesc {
	ret, found := st.GetStorageWeightDescForSectorMaybe(sectorNumber)
	Assert(found)
	return ret
}

func (st *StorageMinerActorState_I) _getStorageWeightDescsForSectors(sectorNumbers []sector.SectorNumber) []SectorStorageWeightDesc {
	ret := []SectorStorageWeightDesc{}
	for _, sectorNumber := range sectorNumbers {
		ret = append(ret, st._getStorageWeightDescForSector(sectorNumber))
	}
	return ret
}

func MinerPoStState_New_OK(lastSuccessfulPoSt actors.ChainEpoch) MinerPoStState {
	return MinerPoStState_Make_OK(&MinerPoStState_OK_I{
		LastSuccessfulPoSt_: lastSuccessfulPoSt,
	})
}

func MinerPoStState_New_Challenged(
	surpriseChallengeEpoch actors.ChainEpoch,
	challengedSectors []sector.SectorNumber,
	numConsecutiveFailures int,
) MinerPoStState {
	return MinerPoStState_Make_Challenged(&MinerPoStState_Challenged_I{
		SurpriseChallengeEpoch_: surpriseChallengeEpoch,
		ChallengedSectors_:      challengedSectors,
		NumConsecutiveFailures_: numConsecutiveFailures,
	})
}

func MinerPoStState_New_DetectedFault(numConsecutiveFailures int) MinerPoStState {
	return MinerPoStState_Make_DetectedFault(&MinerPoStState_DetectedFault_I{
		NumConsecutiveFailures_: numConsecutiveFailures,
	})
}

func (x *SectorOnChainInfo_I) Is_TemporaryFault() bool {
	ret := (x.State() == SectorState_TemporaryFault)
	if ret {
		Assert(x.DeclaredFaultEpoch() != epochUndefined)
		Assert(x.DeclaredFaultDuration() != epochUndefined)
	}
	return ret
}

func (x *SectorOnChainInfo_I) EffectiveFaultBeginEpoch() actors.ChainEpoch {
	Assert(x.Is_TemporaryFault())
	return x.DeclaredFaultEpoch() + indices.StorageMining_DeclaredFaultEffectiveDelay()
}

func (x *SectorOnChainInfo_I) EffectiveFaultEndEpoch() actors.ChainEpoch {
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

func (st *StorageMinerActorState_I) VerifySurprisePoStMeetsTargetReq(candidate sector.PoStCandidate) bool {
	// TODO: Determine what should be the acceptance criterion for sector numbers proven in SurprisePoSt proofs.
	TODO()
	panic("")
}

func SectorNumberSetHAMT_Items(x SectorNumberSetHAMT) []sector.SectorNumber {
	IMPL_FINISH()
	panic("")
}
