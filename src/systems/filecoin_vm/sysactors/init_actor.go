package sysactors

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

// import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"

// import st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
// import util "github.com/filecoin-project/specs/util"

func (a *InitActorCode_I) Constructor(rt Runtime) {
	panic("TODO")
}

func (a *InitActorCode_I) Exec(rt Runtime, codeCID actor.CodeCID, constructorParams actor.MethodParams) InvocOutput {
	rt.CreateActor(codeCID, constructorParams)

	panic("TODO")
	/*
		// TODO: update

		// Make sure that only the actors defined in the spec can be launched.
		if !a._isBuiltinActor(codeCID) {
			rt.Fatal("cannot launch actor instance that is not a builtin actor")
		}

		// Ensure that singeltons can only be launched once.
		// TODO: do we want to enforce this? If so how should actors be marked as such?
		if a._isSingletonActor(codeCID) {
			rt.Fatal("cannot launch another actor of this type")
		}

		// Get the actor ID for this actor.
		actorID := a._assignNextID(state)

		// This generates a unique address for this actor that is stable across message
		// reordering
		// TODO: where do `creator` and `nonce` come from?
		// TODO: CallSeqNum is not related to From -- it's related to Origin
		// addr := rt.ComputeActorAddress(rt.Invocation().FromActor(), rt.Invocation().CallSeqNum())
		var addr addr.Address // TODO



		var initBalance actor.TokenAmount
		panic("TODO")
		// TODO: initBalance := rt.Invocation().InitSendInput.Value()



		// Set up the actor itself
		actorState := &actor.ActorState_I{
			CodeCID_: codeCID,
			// State_:   nil, // TODO: do we need to init the state? probably not
			Balance_:    initBalance,
			CallSeqNum_: 0,
		}

		// The call to the actors constructor will set up the initial state
		// from the given parameters, setting `actor.Head` to a new value when successfull.
		// TODO: can constructors fail?
		// TODO: this needs to be written such that the specific type Constructor is called.
		//       right now actor.Constructor(p) calls the Actor type, not the concrete type.
		// a.Constructor(params) // TODO: uncomment this.

		// TODO: where is this VM.GlobalState?
		// TODO: do we need this?
		// runtime.State().Storage().Set(actorID, actor)

		// Store the mappings of address to actor ID.
		state.AddressMap()[addr] = actorID
		state.IDMap()[actorID] = addr

		// TODO: adjust this to be proper state setting.
		rt.State().StateTree().ActorStates()[addr] = actorState // atm it's nil

		return addr
	*/
}

func (_ *InitActorCode_I) _assignNextID(state InitActorState) actor.ActorID {
	stateI := state.Impl() // TODO: unwrapping like this is ugly.
	actorID := stateI.NextID_
	stateI.NextID_++
	return actorID
}

func (_ *InitActorCode_I) GetActorIDForAddress(state InitActorState, address addr.Address) actor.ActorID {
	return state.AddressMap()[address]
}

// TODO: derive this OR from a union type
func (_ *InitActorCode_I) _isSingletonActor(codeCID actor.CodeCID) bool {
	return true
	// TODO: uncomment this
	// return codeCID == StorageMarketActor ||
	// 	codeCID == StoragePowerActor ||
	// 	codeCID == CronActor ||
	// 	codeCID == InitActor
}

// TODO: derive this OR from a union type
func (_ *InitActorCode_I) _isBuiltinActor(codeCID actor.CodeCID) bool {
	return true
	// TODO: uncomment this
	// return codeCID == StorageMarketActor ||
	// 	codeCID == StoragePowerActor ||
	// 	codeCID == CronActor ||
	// 	codeCID == InitActor ||
	// 	codeCID == StorageMinerActor ||
	// 	codeCID == PaymentChannelActor
}

// TODO
func (a *InitActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	// TODO: load state
	// var state InitActorState
	// storage := input.Runtime().State().Storage()
	// err := loadActorState(storage, input.ToActor().State(), &state)

	switch method {
	// case 0: -- disable: value send
	// case 1: -- disable: cron. init has no cron action
	// case 2:
	// 	return a.InitConstructor(input, state)
	// case 3:
	// 	return a.Exec(input, state, params[0], params[1])
	// case 4:
	// 	return a.GetActorIDForAddress(input, state, params[0])
	default:
		return rt.ErrorReturn(exitcode.SystemError(exitcode.InvalidMethod))
	}
}
