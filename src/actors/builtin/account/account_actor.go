package account

import (
	addr "github.com/filecoin-project/go-address"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	cid "github.com/ipfs/go-cid"
)

type InvocOutput = vmr.InvocOutput

type AccountActorCode interface {
	Impl() *AccountActorCode_I
}

type AccountActorState interface {
	Address() addr.Address
	Impl() *AccountActorState_I
}

type AccountActorCode_I struct {
	cached_cid []byte
}

func (a *AccountActorCode_I) Constructor(rt vmr.Runtime) InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

func (a *AccountActorCode_I) Impl() *AccountActorCode_I {
	return a
}

type AccountActorState_I struct {
	Address_   addr.Address
	cached_cid []byte
}

func (st *AccountActorState_I) CID() cid.Cid {
	panic("TODO")
}

func (a *AccountActorState_I) Address() addr.Address {
	return a.Address_
}

func (a *AccountActorState_I) Impl() *AccountActorState_I {
	return a
}

type AccountActorCode_R struct {
	cid         []byte
	cached_impl *AccountActorCode_I
}

func (a *AccountActorCode_R) Impl() *AccountActorCode_I {
	return a.cached_impl
}

type AccountActorState_R struct {
	cid         []byte
	cached_impl *AccountActorState_I
}

func (a *AccountActorState_R) Impl() *AccountActorState_I {
	return a.cached_impl
}

func (a *AccountActorState_R) Address() addr.Address {
	return a.Impl().Address_
}
