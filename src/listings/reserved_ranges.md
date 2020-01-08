---
title: "Reserved Ranges"
---

# Actor Reserved Ranges

| Actor                | ID |
|---|---|
| SystemActor          | 1 |
| InitActor            | 2 |
| CronActor            | 3 |
| AccountActor         | 4 |
| StoragePowerActor    | 5 |
| StorageMinerActor    | 6 |
| StorageMarketActor   | 7 |
| PaymentChannelActor  | 8 |

# Method Reserved Ranges

## InitActor Methods
| Method                                    | ID |
|---|---|
| Method_InitActor_Exec                     | 1 |
| Method_InitActor_GetActorIDForAddress     | 2 |

## RewardActor Methods
| Method                                    | ID |
|---|---|
| Method_RewardActor_AwardBlockReward       | 1 |

## MultiSigActor Methods
| Method                                            | ID |
|---|---|
| Method_MultiSigActor_Propose                      | 1 |
| Method_MultiSigActor_Approve                      | 2 |
|	Method_MultiSigActor_AddAuthorizedParty           | 3 |
|	Method_MultiSigActor_RemoveAuthorizedParty        | 4 |
|	Method_MultiSigActor_SwapAuthorizedParty          | 5 |
|	Method_MultiSigActor_ChangeNumApprovalsThreshold  | 6 |

## StorageMinerActor Methods
| Method                                              | ID |
|---|---|
| Method_StorageMinerActor_OnDeferredCronEvent        | # |
|	Method_StorageMinerActor_PreCommitSector            | # |
|	Method_StorageMinerActor_ProveCommitSector          | # |
|	Method_StorageMinerActor_DeclareTemporaryFaults     | # |
|	Method_StorageMinerActor_RecoverTemporaryFaults     | # |
|	Method_StorageMinerActor_ExtendSectorExpiration     | # |
|	Method_StorageMinerActor_TerminateSector            | # |
|	Method_StorageMinerActor_SubmitSurprisePoStResponse | # |
|	Method_StorageMinerActor_OnVerifiedElectionPoSt     | # |
|	Method_StorageMinerActor_OnSurprisePoStChallenge    | # |
|	Method_StorageMinerActor_GetPoStState               | # |
|	Method_StorageMinerActor_GetOwnerAddr               | # |
|	Method_StorageMinerActor_GetWorkerAddr              | # |
|	Method_StorageMinerActor_GetWorkerVRFKey            | # |

## StorageMarketActor Methods
| Method                                                                  | ID |
|---|---|
| Method_StorageMarketActor_OnEpochTickEnd                                | 1 |
| Method_StorageMarketActor_AddBalance                                    | 2 |
|	Method_StorageMarketActor_WithdrawBalance                               | 3 |
|	Method_StorageMarketActor_PublishStorageDeals                           | 4 |
|	Method_StorageMarketActor_OnMinerSectorPreCommit_VerifyDealsOrAbort     | 5 |
|	Method_StorageMarketActor_OnMinerSectorProveCommit_VerifyDealsOrAbort   | 6 |
|	Method_StorageMarketActor_OnMinerSectorsTerminate                       | 7 |
|	Method_StorageMarketActor_GetPieceInfosForDealIDs                       | 8 |
|	Method_StorageMarketActor_GetWeightForDealSet                           | 9 |

## StoragePowerActor Methods
| Method                                                        | ID |
|---|---|
|	Method_StoragePowerActor_OnEpochTickEnd                       | 1 |
|	Method_StoragePowerActor_AddBalance                           | 2 |
|	Method_StoragePowerActor_WithdrawBalance                      | 3 |
|	Method_StoragePowerActor_CreateMiner                          | 4 |
|	Method_StoragePowerActor_DeleteMiner                          | 5 |
|	Method_StoragePowerActor_ReportConsensusFault                 | 6 |
|	Method_StoragePowerActor_OnSectorProveCommit                  | 7 |
|	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveBegin | 8 |
|	Method_StoragePowerActor_OnSectorTemporaryFaultEffectiveEnd   | 9 |
|	Method_StoragePowerActor_OnSectorModifyWeightDesc             | 10 |
|	Method_StoragePowerActor_OnSectorTerminate                    | 11 |
|	Method_StoragePowerActor_OnMinerSurprisePoStSuccess           | 12 |
|	Method_StoragePowerActor_OnMinerSurprisePoStFailure           | 13 |
|	Method_StoragePowerActor_OnMinerEnrollCronEvent               | 14 |
|	Method_StoragePowerActor_GetMinerConsensusPower               | 15 |
|	Method_StoragePowerActor_GetMinerUnmetPledgeCollateralRequirement | 16 |

# Error Codes

TODO
