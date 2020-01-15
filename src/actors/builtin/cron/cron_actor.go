package cron

import (
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
)

type CronActorState interface {
	Impl() *CronActorState_I
}

type CronTableEntry interface {
	ToAddr() addr.Address
	MethodNum() abi.MethodNum
	Impl() *CronTableEntry_I
}

type CronActorCode interface {
	Entries() []CronTableEntry
	EpochTick(r vmr.Runtime)
	Impl() *CronActorCode_I
}

func (a *CronActorCode_I) Constructor(rt vmr.Runtime) vmr.InvocOutput {
	// Nothing. intentionally left blank.
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)
	return rt.SuccessReturn()
}

func (a *CronActorCode_I) EpochTick(rt vmr.Runtime) vmr.InvocOutput {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

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

type CronActorState_I struct {
	cached_cid []byte
}

func (c *CronActorState_I) Impl() *CronActorState_I {
	return c
}

type CronActorState_R struct {
	cid         []byte
	cached_impl *CronActorState_I
}

func (c *CronActorState_R) Impl() *CronActorState_I {
	return c.cached_impl
}

type CronTableEntry_I struct {
	ToAddr_    addr.Address
	MethodNum_ abi.MethodNum
	cached_cid []byte
}

func (c *CronTableEntry_I) ToAddr() addr.Address {
	return c.ToAddr_
}

func (c *CronTableEntry_I) MethodNum() abi.MethodNum {
	return c.MethodNum_
}

func (c *CronTableEntry_I) Impl() *CronTableEntry_I {
	return c
}

type CronTableEntry_R struct {
	cid         []byte
	cached_impl *CronTableEntry_I
}

func (c *CronTableEntry_R) ToAddr() addr.Address {
	return c.Impl().ToAddr_
}

func (c *CronTableEntry_R) MethodNum() abi.MethodNum {
	return c.Impl().MethodNum_
}

func (c *CronTableEntry_R) Impl() *CronTableEntry_I {
	return c.cached_impl
}

type CronActorCode_I struct {
	Entries_   []CronTableEntry
	cached_cid []byte
}

func (c *CronActorCode_I) Impl() *CronActorCode_I {
	return c
}

func (c *CronActorCode_I) Entries() []CronTableEntry {
	return c.Entries_
}

type CronActorCode_R struct {
	cid         []byte
	cached_impl *CronActorCode_I
}

func (c *CronActorCode_R) Entries() []CronTableEntry {
	return c.Impl().Entries_
}

func (c *CronActorCode_R) Impl() *CronActorCode_I {
	return c.cached_impl
}
