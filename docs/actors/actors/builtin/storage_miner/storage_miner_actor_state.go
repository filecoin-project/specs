package storage_miner

import (
	"math/big"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	indices "github.com/filecoin-project/specs-actors/actors/runtime/indices"
	autil "github.com/filecoin-project/specs-actors/actors/util"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// Balance of a StorageMinerActor should equal exactly the sum of PreCommit deposits
// that are not yet returned or burned.
type StorageMinerActorState struct {
	Sectors    SectorsAMT
	PoStState  MinerPoStState
	ProvingSet SectorNumberSetHAMT
	Info       MinerInfo
}

type MinerPoStState struct {
	// Epoch of the last succesful PoSt, either election post or surprise post.
	LastSuccessfulPoSt abi.ChainEpoch

	// If >= 0 miner has been challenged and not yet responded successfully.
	// SurprisePoSt challenge state: The miner has not submitted timely ElectionPoSts,
	// and as a result, the system has fallen back to proving storage via SurprisePoSt.
	//  `epochUndefined` if not currently challeneged.
	SurpriseChallengeEpoch abi.ChainEpoch

	// Not empty iff the miner is challenged.
	ChallengedSectors []abi.SectorNumber

	// Number of surprised post challenges that have been failed since last successful PoSt.
	// Indicates that the claimed storage power may not actually be proven. Recovery can proceed by
	// submitting a correct response to a subsequent SurprisePoSt challenge, up until
	// the limit of number of consecutive failures.
	NumConsecutiveFailures int64
}

func (mps *MinerPoStState) Is_Challenged() bool {
	return mps.SurpriseChallengeEpoch != epochUndefined
}

func (mps *MinerPoStState) Is_OK() bool {
	return !mps.Is_Challenged() && !mps.Is_DetectedFault()
}

func (mps *MinerPoStState) Is_DetectedFault() bool {
	return mps.NumConsecutiveFailures > 0
}

type SectorState int64

const (
	PreCommit SectorState = iota
	Active
	TemporaryFault
)

type SectorOnChainInfo struct {
	State                 SectorState
	Info                  SectorPreCommitInfo // Also contains Expiration field.
	PreCommitDeposit      abi.TokenAmount
	PreCommitEpoch        abi.ChainEpoch
	ActivationEpoch       abi.ChainEpoch // -1 if still in PreCommit state.
	DeclaredFaultEpoch    abi.ChainEpoch // -1 if not currently declared faulted.
	DeclaredFaultDuration abi.ChainEpoch // -1 if not currently declared faulted.
	DealWeight            big.Int        // -1 if not yet validated with StorageMarketActor.
}

type SectorPreCommitInfo struct {
	SectorNumber abi.SectorNumber
	SealedCID    abi.SealedSectorCID // CommR
	SealEpoch    abi.ChainEpoch
	DealIDs      abi.DealIDs
	Expiration   abi.ChainEpoch
}

type SectorProveCommitInfo struct {
	SectorNumber     abi.SectorNumber
	RegisteredProof  abi.RegisteredProof
	Proof            abi.SealProof
	InteractiveEpoch abi.ChainEpoch
	Expiration       abi.ChainEpoch
}

// TODO AMT
type SectorsAMT map[abi.SectorNumber]SectorOnChainInfo

// TODO HAMT
type SectorNumberSetHAMT map[abi.SectorNumber]bool

type MinerInfo struct {
	// Account that owns this miner.
	// - Income and returned collateral are paid to this address.
	// - This address is also allowed to change the worker address for the miner.
	Owner addr.Address // Must be an ID-address.

	// Worker account for this miner.
	// This will be the key that is used to sign blocks created by this miner, and
	// sign messages sent on behalf of this miner to commit sectors, submit PoSts, and
	// other day to day miner activities.
	Worker       addr.Address // Must be an ID-address.
	WorkerVRFKey addr.Address // Must be a SECP or BLS address

	// Libp2p identity that should be used when connecting to this miner.
	PeerId peer.ID

	// Amount of space in each sector committed to the network by this miner.
	SectorSize             abi.SectorSize
	SealPartitions         int64
	ElectionPoStPartitions int64
	SurprisePoStPartitions int64
}

func (st *StorageMinerActorState) CID() cid.Cid {
	panic("TODO")
}

func (st *StorageMinerActorState) _getSectorOnChainInfo(sectorNo abi.SectorNumber) (info SectorOnChainInfo, ok bool) {
	sectorInfo, found := st.Sectors[sectorNo]
	if !found {
		return SectorOnChainInfo{}, false
	}
	return sectorInfo, true
}

func (st *StorageMinerActorState) _getSectorDealIDsAssert(sectorNo abi.SectorNumber) abi.DealIDs {
	sectorInfo, found := st._getSectorOnChainInfo(sectorNo)
	Assert(found)
	return sectorInfo.Info.DealIDs
}

func SectorsAMT_Empty() SectorsAMT {
	IMPL_FINISH()
	panic("")
}

func SectorNumberSetHAMT_Empty() SectorNumberSetHAMT {
	IMPL_FINISH()
	panic("")
}

func (st *StorageMinerActorState) GetStorageWeightDescForSectorMaybe(sectorNumber abi.SectorNumber) (ret SectorStorageWeightDesc, ok bool) {
	sectorInfo, found := st.Sectors[sectorNumber]
	if !found {
		ret = autil.SectorStorageWeightDesc{}
		ok = false
		return
	}

	ret = autil.SectorStorageWeightDesc{
		SectorSize: st.Info.SectorSize,
		DealWeight: sectorInfo.DealWeight,
		Duration:   sectorInfo.Info.Expiration - sectorInfo.ActivationEpoch,
	}
	ok = true
	return
}

func (st *StorageMinerActorState) _getStorageWeightDescForSector(sectorNumber abi.SectorNumber) SectorStorageWeightDesc {
	ret, found := st.GetStorageWeightDescForSectorMaybe(sectorNumber)
	Assert(found)
	return ret
}

func (st *StorageMinerActorState) _getStorageWeightDescsForSectors(sectorNumbers []abi.SectorNumber) []SectorStorageWeightDesc {
	ret := []SectorStorageWeightDesc{}
	for _, sectorNumber := range sectorNumbers {
		ret = append(ret, st._getStorageWeightDescForSector(sectorNumber))
	}
	return ret
}

func (x *SectorOnChainInfo) Is_TemporaryFault() bool {
	ret := (x.State == TemporaryFault)
	if ret {
		Assert(x.DeclaredFaultEpoch != epochUndefined)
		Assert(x.DeclaredFaultDuration != epochUndefined)
	}
	return ret
}

// Must be significantly larger than DeclaredFaultEpoch, since otherwise it may be possible
// to declare faults adaptively in order to exempt challenged sectors.
func (x *SectorOnChainInfo) EffectiveFaultBeginEpoch() abi.ChainEpoch {
	Assert(x.Is_TemporaryFault())
	return x.DeclaredFaultEpoch + indices.StorageMining_DeclaredFaultEffectiveDelay()
}

func (x *SectorOnChainInfo) EffectiveFaultEndEpoch() abi.ChainEpoch {
	Assert(x.Is_TemporaryFault())
	return x.EffectiveFaultBeginEpoch() + x.DeclaredFaultDuration
}

func MinerInfo_New(
	ownerAddr addr.Address, workerAddr addr.Address, sectorSize abi.SectorSize, peerId peer.ID) MinerInfo {

	ret := &MinerInfo{
		Owner:      ownerAddr,
		Worker:     workerAddr,
		PeerId:     peerId,
		SectorSize: sectorSize,
	}

	TODO() // TODO: determine how to generate/validate VRF key and initialize other fields

	return *ret
}

func (st *StorageMinerActorState) VerifySurprisePoStMeetsTargetReq(candidate abi.PoStCandidate) bool {
	// TODO: Determine what should be the acceptance criterion for sector numbers proven in SurprisePoSt proofs.
	TODO()
	panic("")
}

func SectorNumberSetHAMT_Items(x SectorNumberSetHAMT) []abi.SectorNumber {
	IMPL_FINISH()
	panic("")
}
