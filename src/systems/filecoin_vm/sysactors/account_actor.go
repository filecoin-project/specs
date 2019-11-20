package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	return rt.SuccessReturn()
}

func (a *AccountActorCode_I) InvokeMethod(rt vmr.Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	// AccountActor has no methods.
	rt.Abort("Invalid method")
	panic("")
}
