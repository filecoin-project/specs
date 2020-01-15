package storage_market

import (
	actor "github.com/filecoin-project/specs/actors"
	abi "github.com/filecoin-project/specs/actors/abi"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	autil "github.com/filecoin-project/specs/actors/util"
	cid "github.com/ipfs/go-cid"
)

type BalanceTableHAMT = autil.BalanceTableHAMT
type DealIDQueue = autil.DealIDQueue

var RT_MinerEntry_ValidateCaller_DetermineFundsLocation = vmr.RT_MinerEntry_ValidateCaller_DetermineFundsLocation
var RT_ValidateImmediateCallerIsSignable = vmr.RT_ValidateImmediateCallerIsSignable

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
//
// This boilerplate should be essentially identical for all actors, and
// conceptually belongs in the runtime/VM. It is only duplicated here as a
// workaround due to the lack of generics support in Go.
////////////////////////////////////////////////////////////////////////////////

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = abi.Bytes

var Assert = autil.Assert
var IMPL_FINISH = autil.IMPL_FINISH
var IMPL_TODO = autil.IMPL_TODO
var TODO = autil.TODO

func (a *StorageMarketActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, StorageMarketActorState) {
	h := rt.AcquireState()
	var state StorageMarketActorState_I
	stateCID := cid.Cid(h.Take())
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, &state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st StorageMarketActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st StorageMarketActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StorageMarketActorState_I) CID() cid.Cid {
	IMPL_FINISH()
	panic("")
}

func DealsAMT_Empty() DealsAMT {
	IMPL_FINISH()
	panic("")
}

func CachedDealIDsByPartyHAMT_Empty() CachedDealIDsByPartyHAMT {
	IMPL_FINISH()
	panic("")
}

func CachedExpirationsPendingHAMT_Empty() CachedExpirationsPendingHAMT {
	IMPL_FINISH()
	panic("")
}
