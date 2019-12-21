package actor_interfaces

import (
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
)

const (
	Method_InitActor_Exec = actor.MethodPlaceholder + iota
	Method_InitActor_GetActorIDForAddress
)

const (
	Method_RewardActor_AwardBlockReward = actor.MethodPlaceholder + iota
)

const (
	// Proxy cron tick method (via StoragePowerActor)
	Method_StorageMinerActor_OnDeferredCronEvent = actor.MethodPlaceholder + iota

	// User-callable methods
	Method_StorageMinerActor_PreCommitSector
	Method_StorageMinerActor_ProveCommitSector
	Method_StorageMinerActor_DeclareTemporaryFaults
	Method_StorageMinerActor_RecoverTemporaryFaults
	Method_StorageMinerActor_TerminateSector
	Method_StorageMinerActor_SubmitSurprisePoStResponse

	// Internal mechanism events
	Method_StorageMinerActor_OnVerifiedElectionPoSt
	Method_StorageMinerActor_OnSurprisePoStChallenge
	Method_StorageMinerActor_OnSurprisePoStExpiryCheck
	Method_StorageMinerActor_OnSectorExpiryCheck

	// State queries
	Method_StorageMinerActor_GetPoStState
	Method_StorageMinerActor_GetOwnerAddr
	Method_StorageMinerActor_GetWorkerAddr
	Method_StorageMinerActor_GetWorkerVRFKey
)

const (
	// Cron tick method
	Method_StorageMarketActor_OnEpochTickEnd = actor.MethodPlaceholder + iota

	// User-callable methods
	Method_StorageMarketActor_AddBalance
	Method_StorageMarketActor_WithdrawBalance
	Method_StorageMarketActor_PublishStorageDeals

	// Internal mechanism events
	Method_StorageMarketActor_OnMinerSectorPreCommit_VerifyDealsOrAbort
	Method_StorageMarketActor_OnMinerSectorProveCommit_VerifyDealsOrAbort
	Method_StorageMarketActor_OnMinerSectorsTerminate

	// State queries
	Method_StorageMarketActor_GetPieceInfosForDealIDs
	Method_StorageMarketActor_GetWeightForDealSet
)

const (
	// Cron tick method
	Method_StoragePowerActor_OnEpochTickEnd = actor.MethodPlaceholder + iota

	// User-callable methods
	Method_StoragePowerActor_AddBalance
	Method_StoragePowerActor_WithdrawBalance
	Method_StoragePowerActor_CreateMiner
	Method_StoragePowerActor_DeleteMiner

	// Internal mechanism events
	Method_StoragePowerActor_OnSectorProveCommit
	Method_StoragePowerActor_OnSectorTerminate
	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveBegin
	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd
	Method_StoragePowerActor_OnMinerSurprisePoStSuccess
	Method_StoragePowerActor_OnMinerSurprisePoStFailure
	Method_StoragePowerActor_OnMinerEnrollCronEvent

	// State queries
	Method_StoragePowerActor_GetMinerConsensusPower
	Method_StoragePowerActor_GetMinerUnmetPledgeCollateralRequirement
)
