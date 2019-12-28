package init

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	autil "github.com/filecoin-project/specs/actors/util"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	node_base "github.com/filecoin-project/specs/systems/filecoin_nodes/node_base"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
)

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = abi.Bytes

var AssertMsg = autil.AssertMsg

func (a *InitActorCode_I) Constructor(rt Runtime) InvocOutput {
	rt.ValidateImmediateCallerIs(addr.SystemActorAddr)
	h := rt.AcquireState()
	st := &InitActorState_I{
		AddressMap_:  map[addr.Address]addr.ActorID{}, // TODO: HAMT
		NextID_:      addr.ActorID(addr.FirstNonSingletonActorId),
		NetworkName_: vmr.NetworkName(),
	}
	UpdateRelease(rt, h, st)
	return rt.ValueReturn(nil)
}

func (a *InitActorCode_I) Exec(rt Runtime, execCodeID abi.ActorCodeID, constructorParams abi.MethodParams) InvocOutput {
	rt.ValidateImmediateCallerAcceptAny()
	callerCodeID, ok := rt.GetActorCodeID(rt.ImmediateCaller())
	AssertMsg(ok, "no code for actor at %s", rt.ImmediateCaller())
	if !_codeIDSupportsExec(callerCodeID, execCodeID) {
		rt.AbortArgMsg("Caller type cannot create an actor of requested type")
	}

	// Compute a re-org-stable address.
	// This address exists for use by messages coming from outside the system, in order to
	// stably address the newly created actor even if a chain re-org causes it to end up with
	// a different ID.
	newAddr := rt.NewActorAddress()

	// Allocate an ID for this actor.
	// Store mapping of pubkey or actor address to actor ID
	h, st := _loadState(rt)
	idAddr := st.MapAddressToNewID(newAddr)
	UpdateRelease(rt, h, st)

	// Create an empty actor.
	rt.CreateActor(execCodeID, idAddr)

	// Invoke constructor. If construction fails, the error should propagate and cause
	// Exec to fail too.
	rt.SendPropagatingErrors(&vmr.InvocInput_I{
		To_:     idAddr,
		Method_: actor.MethodConstructor,
		Params_: constructorParams,
		Value_:  rt.ValueReceived(),
	})

	return rt.ValueReturn(
		Bytes(addr.Serialize_Address_Compact(idAddr)))
}

// This method is disabled until proven necessary.
//func (a *InitActorCode_I) GetActorIDForAddress(rt Runtime, address addr.Address) InvocOutput {
//	h, st := _loadState(rt)
//	actorID := st.AddressMap()[address]
//	Release(rt, h, st)
//	return rt.ValueReturn(Bytes(addr.Serialize_ActorID(actorID)))
//}

func (s *InitActorState_I) ResolveAddress(address addr.Address) addr.Address {
	actorID, ok := s.AddressMap()[address]
	if ok {
		return addr.Address_Make_ID(node_base.NETWORK, actorID)
	}
	return address
}

func (s *InitActorState_I) MapAddressToNewID(address addr.Address) addr.Address {
	actorID := s.NextID_
	s.NextID_++
	s.AddressMap()[address] = actorID
	return addr.Address_Make_ID(node_base.NETWORK, actorID)
}

func _codeIDSupportsExec(callerCodeID abi.ActorCodeID, execCodeID abi.ActorCodeID) bool {
	if execCodeID == builtin.AccountActorCodeID {
		// Special case: account actors must be created implicitly by sending value;
		// cannot be created via exec.
		return false
	}

	if execCodeID == builtin.PaymentChannelActorCodeID {
		return true
	}

	if execCodeID == builtin.StorageMinerActorCodeID {
		if callerCodeID == builtin.StoragePowerActorCodeID {
			return true
		}
	}

	return false
}

///// Boilerplate /////

func _loadState(rt Runtime) (vmr.ActorStateHandle, InitActorState) {
	h := rt.AcquireState()
	stateCID := ipld.CID(h.Take())
	var state InitActorState_I
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, &state
}

func Release(rt Runtime, h vmr.ActorStateHandle, st InitActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}

func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st InitActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}

func (st *InitActorState_I) CID() ipld.CID {
	panic("TODO")
}
