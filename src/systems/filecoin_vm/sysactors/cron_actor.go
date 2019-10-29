package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"

type Runtime = vmr.Runtime

func (a *CronActorCode_I) Constructor(rt vmr.Runtime) {
	// Nothing. intentionally left blank.
}

func (a *CronActorCode_I) EpochTick(rt vmr.Runtime) {
	// Hook period actions in here.

	// a.actors is basically a static registry for now, loaded
	// in the interpreter static registry.
	for _, a := range a.Actors() {
		rt.SendAllowingErrors(msg.InvocInput_Make(
			a,
			vmr.Reserved_CronMethod,
			[]actor.MethodParam{},
			actor.TokenAmount(0)),
		)
	}

	rt.ReturnSuccess()
}

func (a *CronActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) {
	switch method {
	case vmr.Reserved_CronMethod:
		rt.Assert(len(params) == 0)
		a.EpochTick(rt)
	default:
		rt.ReturnError(exitcode.SystemError(exitcode.InvalidMethod))
	}
}
