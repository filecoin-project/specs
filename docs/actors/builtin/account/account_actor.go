package account

import (
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	cid "github.com/ipfs/go-cid"
)

type InvocOutput = vmr.InvocOutput

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

func (st *AccountActorState_I) CID() cid.Cid {
	panic("TODO")
}
