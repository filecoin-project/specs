package storage_market

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	actor_util "github.com/filecoin-project/specs/systems/filecoin_vm/actor_util"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

type BalanceTableHAMT = actor_util.BalanceTableHAMT
type DealIDQueue = actor_util.DealIDQueue

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
//
// This boilerplate should be essentially identical for all actors, and
// conceptually belongs in the runtime/VM. It is only duplicated here as a
// workaround due to the lack of generics support in Go.
////////////////////////////////////////////////////////////////////////////////

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes

var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO

func (a *StorageMarketActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, StorageMarketActorState) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state := Deserialize_StorageMarketActorState_Assert(stateBytes.As_Bytes())
	return h, state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st StorageMarketActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st StorageMarketActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StorageMarketActorState_I) CID() ipld.CID {
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
