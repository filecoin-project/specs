package builtin

import (
	abi "github.com/filecoin-project/specs/actors/abi"
)

const (
	MethodSend        = abi.MethodNum(0)
	MethodConstructor = abi.MethodNum(1)

	// TODO: remove this once canonical method numbers are finalized
	MethodPlaceholder = abi.MethodNum(1 << 30)
)

const (
	Method_InitActor_Exec = MethodPlaceholder + iota
	Method_InitActor_GetActorIDForAddress
)

const (
	Method_CronActor_EpochTick = MethodPlaceholder + iota
)

const (
	Method_RewardActor_AwardBlockReward = MethodPlaceholder + iota
)

const (
	Method_MultiSigActor_Propose = MethodPlaceholder + iota
	Method_MultiSigActor_Approve
	Method_MultiSigActor_AddAuthorizedParty
	Method_MultiSigActor_RemoveAuthorizedParty
	Method_MultiSigActor_SwapAuthorizedParty
	Method_MultiSigActor_ChangeNumApprovalsThreshold
)

const (
	// Proxy cron tick method (via StoragePowerActor)
	Method_StorageMinerActor_OnDeferredCronEvent = MethodPlaceholder + iota

	// User-callable methods
	Method_StorageMinerActor_PreCommitSector
	Method_StorageMinerActor_ProveCommitSector
	Method_StorageMinerActor_DeclareTemporaryFaults
	Method_StorageMinerActor_RecoverTemporaryFaults
	Method_StorageMinerActor_ExtendSectorExpiration
	Method_StorageMinerActor_TerminateSector
	Method_StorageMinerActor_SubmitSurprisePoStResponse

	// Internal mechanism events
	Method_StorageMinerActor_OnVerifiedElectionPoSt
	Method_StorageMinerActor_OnSurprisePoStChallenge
	Method_StorageMinerActor_OnDeleteMiner

	// State queries
	Method_StorageMinerActor_GetPoStState
	Method_StorageMinerActor_GetOwnerAddr
	Method_StorageMinerActor_GetWorkerAddr
	Method_StorageMinerActor_GetWorkerVRFKey
)

const (
	// Cron tick method
	Method_StorageMarketActor_OnEpochTickEnd = MethodPlaceholder + iota

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
	Method_StoragePowerActor_OnEpochTickEnd = MethodPlaceholder + iota

	// User-callable methods
	Method_StoragePowerActor_AddBalance
	Method_StoragePowerActor_WithdrawBalance
	Method_StoragePowerActor_CreateMiner
	Method_StoragePowerActor_DeleteMiner
	Method_StoragePowerActor_ReportConsensusFault

	// Internal mechanism events
	Method_StoragePowerActor_OnSectorProveCommit
	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveBegin
	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd
	Method_StoragePowerActor_OnSectorModifyWeightDesc
	Method_StoragePowerActor_OnSectorTerminate
	Method_StoragePowerActor_OnMinerSurprisePoStSuccess
	Method_StoragePowerActor_OnMinerSurprisePoStFailure
	Method_StoragePowerActor_OnMinerEnrollCronEvent

	// State queries
	Method_StoragePowerActor_GetMinerConsensusPower
	Method_StoragePowerActor_GetMinerUnmetPledgeCollateralRequirement
)
