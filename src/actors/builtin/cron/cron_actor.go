package cron

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import abi "github.com/filecoin-project/specs/actors/abi"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import vmr "github.com/filecoin-project/specs/actors/runtime"

const (
	Method_CronActor_EpochTick = actor.MethodPlaceholder + iota
)

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

func (a *CronActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	return rt.SuccessReturn()
}

func (a *CronActorCode_I) EpochTick(rt vmr.Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)

	// a.Entries is basically a static registry for now, loaded
	// in the interpreter static registry.
	for _, entry := range a.Entries() {
		rt.SendCatchingErrors(&vmr.InvocInput_I{
			To_:     entry.ToAddr(),
			Method_: entry.MethodNum(),
			Params_: nil,
			Value_:  abi.TokenAmount(0),
		})
	}

	return rt.SuccessReturn()
}
