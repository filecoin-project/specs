package cron

import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
)

type CronActorState struct{}

type CronActor struct {
	Entries []CronTableEntry
}

type CronTableEntry struct {
	ToAddr    addr.Address
	MethodNum abi.MethodNum
}

func (a *CronActor) Constructor(rt vmr.Runtime) vmr.InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

func (a *CronActor) EpochTick(rt vmr.Runtime) vmr.InvocOutput {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

	// a.Entries is basically a static registry for now, loaded
	// in the interpreter static registry.
	for _, entry := range a.Entries {
		rt.SendCatchingErrors(&vmr.InvocInput_I{
			To_:     entry.ToAddr,
			Method_: entry.MethodNum,
			Params_: nil,
			Value_:  abi.TokenAmount(0),
		})
	}

	return rt.SuccessReturn()
}
