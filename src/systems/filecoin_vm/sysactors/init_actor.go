package sysactors

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

const (
	Method_InitActor_Exec = actor.MethodPlaceholder + iota
	Method_InitActor_GetActorIDForAddress
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type Serialization = util.Serialization

var CheckArgs = actor.CheckArgs
var ArgPop = actor.ArgPop
var ArgEnd = actor.ArgEnd

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

////////////////////////////////////////////////////////////////////////////////

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

func (a *InitActorCode_I) Exec(rt Runtime, codeID actor.CodeID, constructorParams actor.MethodParams) InvocOutput {
	if !_codeIDSupportsExec(codeID) {
		rt.AbortArgMsg("cannot exec an actor of this type")
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
	rt.CreateActor(codeID, idAddr)

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

func _codeIDSupportsExec(codeID actor.CodeID) bool {
	if !codeID.IsBuiltin() || codeID.IsSingleton() {
		return false
	}

	which := codeID.As_Builtin()

	if which == actor.BuiltinActorID_Account {
		// Special case: account actors must be created implicitly by sending value;
		// cannot be created via exec.
		return false
	}

	util.Assert(
		which == actor.BuiltinActorID_PaymentChannel ||
			which == actor.BuiltinActorID_StorageMiner)

	return true
}

func (a *InitActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	switch method {
	case actor.MethodConstructor:
		ArgEnd(&params, rt)
		return a.Constructor(rt)

	case Method_InitActor_Exec:
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
