package sysactors

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate: System singleton actors
//
// This boilerplate should be essentially identical for all actors, and
// conceptually belongs in the runtime/VM. It is only duplicated here as a
// workaround due to the lack of generics support in Go.
////////////////////////////////////////////////////////////////////////////////

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type Serialization = util.Serialization

var CheckArgs = actor_util.CheckArgs
var ArgPop = actor_util.ArgPop
var ArgEnd = actor_util.ArgEnd

////////////////////////////////////////////////////////////////////////////////
// -- InitActor
////////////////////////////////////////////////////////////////////////////////

func _loadState(rt Runtime) (vmr.ActorStateHandle, InitActorState) {
	h := rt.AcquireState()
	stateCID := ipld.CID(h.Take())
	if ipld.CID_Equals(stateCID, ipld.EmptyCID()) {
		rt.AbortAPI("Actor state not initialized")
	}
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state, err := Deserialize_InitActorState(Serialization(stateBytes.As_Bytes()))
	if err != nil {
		rt.AbortAPI("State deserialization error")
	}
	return h, state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st InitActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st InitActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *InitActorState_I) CID() ipld.CID {
	panic("TODO")
}
