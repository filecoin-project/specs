package account

import (
	addr "github.com/filecoin-project/go-address"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	cid "github.com/ipfs/go-cid"
)

type InvocOutput = vmr.InvocOutput

type AccountActor struct{}

func (a *AccountActor) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

type AccountActorState struct {
	Address addr.Address
}

func (AccountActorState) CID() cid.Cid {
	panic("TODO")
}
