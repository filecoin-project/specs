package sysactors

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"

// TODO invariant: always use MethodNum 1 for cron?
var CronMethodNumber = actor.MethodNum(1)

func (a *CronActorCode_I) Constructor(rt vmr.Runtime) {
	// Nothing. intentionally left blank.
}

func (a *CronActorCode_I) Tick(rt vmr.Runtime) {
	// Hook period actions in here.

	// a.actors is basically a static registry for now, loaded
	// in the interpreter static registry.
	for _, a := range a.Actors() {
		rt.Send(a, CronMethodNumber, nil, actor.TokenAmount(0))
	}
}

func (a *CronActorCode_I) InvokeMethod(input vmr.InvocInput, method actor.MethodNum, params actor.MethodParams) vmr.InvocOutput {
	switch method {
	case 1:
		a.Tick(input.Runtime)
		return vmr.InvocOutput{} // todo: grab the updated state tree + return success
	default:
		return vmr.ErrorInvocOutput(input.InTree, exitcode.InvalidMethod)
	}
}
