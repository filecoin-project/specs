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

func (a *InitActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, InitActorState) {
	h := rt.AcquireState()
	stateCID := ipld.CID(h.Take())
	if ipld.CID_Equals(stateCID, ipld.EmptyCID()) {
		if rt.CurrMethodNum() != actor.MethodConstructor {
			rt.AbortAPI("Actor state not initialized")
		}
		return h, nil
	}
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.AbortAPI("IPLD lookup error")
	}
	state, err := Deserialize_InitActorState(Serialization(stateBytes.As_Bytes()))
	if err != nil {
		rt.AbortAPI("State deserialization error")
	}
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

////////////////////////////////////////////////////////////////////////////////

func (a *InitActorCode_I) Constructor(rt Runtime) InvocOutput {
	h, st := a.State(rt)
	st = &InitActorState_I{
		AddressMap_: map[addr.Address]addr.ActorID{}, // TODO: HAMT
		NextID_:     addr.ActorID(0),
	}
	UpdateRelease(rt, h, st)
	return rt.ValueReturn(nil)
}

func (a *InitActorCode_I) Exec(rt Runtime, codeID actor.CodeID, constructorParams actor.MethodParams) InvocOutput {
	if !_codeIDSupportsExec(codeID) {
		rt.AbortArgMsg("cannot exec an actor of this type")
	}

	newAddr := rt.NewActorAddress()

	actorState := &actor.ActorState_I{
		CodeID_:     codeID,
		State_:      actor.ActorSubstateCID(ipld.EmptyCID()),
		Balance_:    actor.TokenAmount(0),
		CallSeqNum_: 0,
	}

	actorStateCID := actor.ActorSystemStateCID(rt.IpldPut(actorState))

	// Get the actor ID for this actor.
	h, st := a.State(rt)
	actorID := st._assignNextID()

	// Store the mappings of address to actor ID.
	st.AddressMap()[newAddr] = actorID

	UpdateRelease(rt, h, st)

	// Note: the following call may fail (e.g., if the actor already exists, or the actor's own
	// constructor call fails). In this case, an error should propagate up and cause Exec to fail
	// as well.
	rt.CreateActor(actorStateCID, newAddr, rt.ValueReceived(), constructorParams)

	return rt.ValueReturn(
		Bytes(addr.Serialize_Address_Compact(newAddr)))
}

func (s *InitActorState_I) _assignNextID() addr.ActorID {
	actorID := s.NextID_
	s.NextID_++
	return actorID
}

func (a *InitActorCode_I) GetActorIDForAddress(rt Runtime, address addr.Address) InvocOutput {
	h, st := a.State(rt)
	actorID := st.AddressMap()[address]
	Release(rt, h, st)
	return rt.ValueReturn(Bytes(addr.Serialize_ActorID(actorID)))
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

	case Method_InitActor_GetActorIDForAddress:
		address, err := addr.Deserialize_Address(ArgPop(&params, rt))
		CheckArgs(&params, rt, err == nil)
		ArgEnd(&params, rt)
		return a.GetActorIDForAddress(rt, address)

	default:
		rt.Abort(exitcode.SystemError(exitcode.InvalidMethod), "Invalid method")
		panic("")
	}
}
