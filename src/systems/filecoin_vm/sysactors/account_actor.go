package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////

func (a *AccountActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, AccountActorState) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state := AccDeserializeState(stateBytes.As_Bytes())
	return h, state
}

func AccRelease(rt Runtime, h vmr.ActorStateHandle, st AccountActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func AccUpdateRelease(rt Runtime, h vmr.ActorStateHandle, st AccountActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *AccountActorState_I) CID() ipld.CID {
	panic("TODO")
}
func AccDeserializeState(x Bytes) AccountActorState {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	return rt.SuccessReturn()
}

func (a *AccountActorCode_I) InvokeMethod(rt vmr.Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	switch method {
	default:
		rt.Abort(exitcode.SystemError(exitcode.InvalidMethod), "Invalid method")
		panic("")
	}
}
