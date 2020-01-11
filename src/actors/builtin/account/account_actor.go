package account

import (
	vmr "github.com/filecoin-project/specs/actors/runtime"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
)

type InvocOutput = vmr.InvocOutput

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	return rt.SuccessReturn()
}

func (st *AccountActorState_I) CID() ipld.CID {
	panic("TODO")
}
