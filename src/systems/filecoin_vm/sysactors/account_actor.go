package sysactors

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
)

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	return rt.SuccessReturn()
}

func (a *AccountActorCode_I) InvokeMethod(rt vmr.Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	switch method {
	case actor.MethodConstructor:
		rt.Assert(len(params) == 0)
		return a.Constructor(rt)

	default:
		// AccountActor has no methods.
		rt.Abort(exitcode.SystemError(exitcode.InvalidMethod), "Invalid method")
		panic("")
	}
}

func (st *AccountActorState_I) CID() ipld.CID {
	panic("TODO")
}
