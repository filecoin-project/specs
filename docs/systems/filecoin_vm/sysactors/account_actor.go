package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
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
