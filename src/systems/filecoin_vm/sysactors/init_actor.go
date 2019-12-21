package sysactors

import (
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
)

func (a *InitActorCode_I) Constructor(rt Runtime) InvocOutput {
	h := rt.AcquireState()
	st := &InitActorState_I{
		AddressMap_:  map[addr.Address]addr.ActorID{}, // TODO: HAMT
		NextID_:      addr.ActorID(addr.FirstNonSingletonActorId),
		NetworkName_: vmr.NetworkName(),
	}
	UpdateRelease(rt, h, st)
	return rt.ValueReturn(nil)
}

func (a *InitActorCode_I) Exec(rt Runtime, execCodeID actor.CodeID, constructorParams actor.MethodParams) InvocOutput {
	rt.ValidateImmediateCallerAcceptAny()
	callerCodeID, ok := rt.GetActorCodeID(rt.ImmediateCaller())
	Assert(ok)
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
		return addr.Address_Make_ID(addr.Address_NetworkID_Testnet, actorID)
	}
	return address
}

func (s *InitActorState_I) MapAddressToNewID(address addr.Address) addr.Address {
	actorID := s.NextID_
	s.NextID_++
	s.AddressMap()[address] = actorID
	return addr.Address_Make_ID(addr.Address_NetworkID_Testnet, actorID)
}

func _codeIDSupportsExec(callerCodeID actor.CodeID, execCodeID actor.CodeID) bool {
	if !execCodeID.IsBuiltin() || execCodeID.IsSingleton() {
		return false
	}

	if execCodeID.As_Builtin() == actor.BuiltinActorID_Account {
		// Special case: account actors must be created implicitly by sending value;
		// cannot be created via exec.
		return false
	}

	if execCodeID.As_Builtin() == actor.BuiltinActorID_PaymentChannel {
		return true
	}

	if execCodeID.As_Builtin() == actor.BuiltinActorID_StorageMiner {
		if callerCodeID.Is_Builtin() && callerCodeID.As_Builtin() == actor.BuiltinActorID_StoragePower {
			return true
		}
	}

	return false
}

func (a *InitActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	switch method {
	case actor.MethodConstructor:
		ArgEnd(&params, rt)
		return a.Constructor(rt)

	case ai.Method_InitActor_Exec:
		codeId, err := actor.Deserialize_CodeID(ArgPop(&params, rt))
		CheckArgs(&params, rt, err == nil)
		// Note: do not call ArgEnd (params is forwarded to Exec)
		return a.Exec(rt, codeId, params)

	//case Method_InitActor_GetActorIDForAddress:
	//	address, err := addr.Deserialize_Address(ArgPop(&params, rt))
	//	CheckArgs(&params, rt, err == nil)
	//	ArgEnd(&params, rt)
	//	return a.GetActorIDForAddress(rt, address)

	default:
		rt.Abort(exitcode.SystemError(exitcode.InvalidMethod), "Invalid method")
		panic("")
	}
}
