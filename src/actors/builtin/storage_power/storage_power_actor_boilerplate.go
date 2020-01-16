package storage_power

import (
	actor "github.com/filecoin-project/specs/actors"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	autil "github.com/filecoin-project/specs/actors/util"
)

type BalanceTableHAMT = autil.BalanceTableHAMT
type SectorStorageWeightDesc = autil.SectorStorageWeightDesc
type SectorTerminationType = autil.SectorTerminationType

var SectorTerminationType_NormalExpiration = autil.SectorTerminationType_NormalExpiration

var RT_MinerEntry_ValidateCaller_DetermineFundsLocation = vmr.RT_MinerEntry_ValidateCaller_DetermineFundsLocation

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
//
// This boilerplate should be essentially identical for all actors, and
// conceptually belongs in the runtime/VM. It is only duplicated here as a
// workaround due to the lack of generics support in Go.
////////////////////////////////////////////////////////////////////////////////

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

var Assert = autil.Assert
var IMPL_FINISH = autil.IMPL_FINISH
var IMPL_TODO = autil.IMPL_TODO
var TODO = autil.TODO

func Release(rt Runtime, h vmr.ActorStateHandle, st StoragePowerActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(&st))
	h.Release(checkCID)
}

func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st StoragePowerActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(&st))
	h.UpdateRelease(newCID)
}
