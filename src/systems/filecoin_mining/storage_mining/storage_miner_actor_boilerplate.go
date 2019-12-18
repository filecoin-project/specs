package storage_mining

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type Any = util.Any
type Bool = util.Bool
type Bytes = util.Bytes
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

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
