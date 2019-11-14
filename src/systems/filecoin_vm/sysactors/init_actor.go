package sysactors

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

const (
	Method_InitActor_Exec = actor.MethodPlaceholder + iota
	Method_InitActor_GetActorIDForAddress
)

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = msg.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type Serialization = util.Serialization

func (a *InitActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, InitActorState) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.Abort("IPLD lookup error")
	}
	state := DeserializeState(stateBytes.As_Bytes())
	return h, state
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
func DeserializeState(x Bytes) InitActorState {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (a *InitActorCode_I) Constructor(rt Runtime) InvocOutput {
	panic("TODO")
}

func (a *InitActorCode_I) Exec(rt Runtime, codeID actor.CodeID, constructorParams actor.MethodParams) InvocOutput {
	if !codeID.IsBuiltin() {
		rt.Abort("cannot launch actor instance that is not a builtin actor")
	}

	if !_codeIDSupportsExec(codeID) {
		rt.Abort("cannot exec an actor of this type")
	}

	// Get the actor ID for this actor.
	h, st := a.State(rt)
	actorID := st._assignNextID()

	// This generates a unique address for this actor that is stable across message
	// reordering
	// TODO: where do `creator` and `nonce` come from?
	// TODO: CallSeqNum is not related to From -- it's related to Origin
	// addr := rt.ComputeActorAddress(rt.Invocation().FromActor(), rt.Invocation().CallSeqNum())
	addr := a._computeNewAddress(rt, actorID)

	initBalance := rt.ValueReceived()

	// Set up the actor itself
	actorState := &actor.ActorState_I{
		CodeID_: codeID,
		// State_:   nil, // TODO: do we need to init the state? probably not
		Balance_:    initBalance,
		CallSeqNum_: 0,
	}

	stateCid := actor.StateCID(rt.IpldPut(actorState))

	// runtime.State().Storage().Set(actorID, actor)

	// Store the mappings of address to actor ID.
	st.AddressMap()[addr] = actorID
	st.IDMap()[actorID] = addr

	// TODO: adjust this to be proper state setting.
	UpdateRelease(rt, h, st)

	// TODO: can this fail?
	rt.CreateActor(stateCid, addr, constructorParams)

	return rt.ValueReturn([]byte(addr.String()))
}

func (s *InitActorState_I) _assignNextID() actor.ActorID {
	actorID := s.NextID_
	s.NextID_++
	return actorID
}

func (_ *InitActorCode_I) _computeNewAddress(rt Runtime, id actor.ActorID) addr.Address {
	// assign an address based on some randomness
	// we use the current epoch, and the actor id. this should be a unique identifier,
	// stable across reorgs.
	//
	// TODO: do we really need this? it's pretty confusing...
	r := rt.Randomness(rt.CurrEpoch(), uint64(id))

	_ = r // TODO: use r in a
	// a := &addr.Address_Type_Actor_I{}
	// n := &addr.Address_NetworkID_Testnet_I{}
	// return addr.MakeAddress(n, a)
	panic("TODO")
	return nil
}

func (a *InitActorCode_I) GetActorIDForAddress(rt Runtime, address addr.Address) InvocOutput {
	h, st := a.State(rt)
	s := st.AddressMap()[address]
	Release(rt, h, st)
	// return rt.ValueReturn(s)
	// TODO
	_ = s
	return rt.ValueReturn(nil)
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
	argError := rt.ErrorReturn(exitcode.SystemError(exitcode.InvalidArguments))

	switch method {
	// case actor.MethodSend:
	//     disable

	case actor.MethodConstructor:
		if len(params) != 0 {
			return argError
		}
		return a.Constructor(rt)

	// case actor.MethodCron:
	//     disable. init has no cron action

	case Method_InitActor_Exec:
		if len(params) == 0 {
			return argError
		}
		codeId, err := actor.Deserialize_CodeID(Serialization(params[0]))
		if err != nil {
			return argError
		}
		params = params[1:]
		return a.Exec(rt, codeId, params)

	case Method_InitActor_GetActorIDForAddress:
		if len(params) == 0 {
			return argError
		}
		address, err := addr.Deserialize_Address(Serialization(params[0]))
		if err != nil {
			return argError
		}
		params = params[1:]

		if len(params) != 0 {
			return argError
		}
		return a.GetActorIDForAddress(rt, address)

	default:
		return rt.ErrorReturn(exitcode.SystemError(exitcode.InvalidMethod))
	}
}
