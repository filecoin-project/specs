package storage_miner

import (
	actor_util "github.com/filecoin-project/specs/actors/util"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

type SectorStorageWeightDesc = actor_util.SectorStorageWeightDesc
type SectorTerminationType = actor_util.SectorTerminationType

var RT_ConfirmFundsReceiptOrAbort_RefundRemainder = vmr.RT_ConfirmFundsReceiptOrAbort_RefundRemainder

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
//
// This boilerplate should be essentially identical for all actors, and
// conceptually belongs in the runtime/VM. It is only duplicated here as a
// workaround due to the lack of generics support in Go.
////////////////////////////////////////////////////////////////////////////////

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO
var TODO = util.TODO

func (a *StorageMinerActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, StorageMinerActorState) {
	h := rt.AcquireState()
	stateCID := ipld.CID(h.Take())
	var state StorageMinerActorState_I
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, &state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st StorageMinerActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st StorageMinerActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StorageMinerActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) StorageMinerActorState {
	panic("TODO")
}
