package account

import (
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
)

type InvocOutput = vmr.InvocOutput

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

func (st *AccountActorState_I) CID() ipld.CID {
	panic("TODO")
}
