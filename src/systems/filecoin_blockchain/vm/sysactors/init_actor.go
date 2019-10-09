package sysactors

import addr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import vmr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime"

// import st "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/state_tree"
// import util "github.com/filecoin-project/specs/util"

func (self *InitActor_I) Constructor() {

}

func (self *InitActor_I) Exec(rt vmr.Runtime, codeCID actor.CodeCID, method actor.MethodNum, params actor.MethodParams) addr.Address {
	// Make sure that only the actors defined in the spec can be launched.
	if !self._isBuiltinActor(codeCID) {
		rt.Fatal("cannot launch actor instance that is not a builtin actor")
	}

	// Get the actor ID for this actor.
	actorID := self._assignNextID()

	// Ensure that singeltons can only be launched once.
	// TODO: do we want to enforce this? If so how should actors be marked as such?
	if self._isSingletonActor(codeCID) {
		rt.Fatal("cannot launch another actor of this type")
	}

	// This generates a unique address for this actor that is stable across message
	// reordering
	// TODO: where do `creator` and `nonce` come from?
	addr := rt.ComputeActorAddress(rt.Message().From(), rt.Message().CallSeqNum())

	// Set up the actor itself
	actor := actor.Actor_I{
		CodeCID_: codeCID,
		Balance_: rt.Message().Value(),
		State_:   nil,
		SeqNum_:  0,
	}

	// The call to the actors constructor will set up the initial state
	// from the given parameters, setting `actor.Head` to a new value when successfull.
	// TODO: can constructors fail?
	// TODO: this needs to be written such that the specific type Constructor is called.
	//       right now actor.Constructor(p) calls the Actor type, not the concrete type.
	// actor.Constructor(params) // TODO: uncomment this.

	// TODO: where is this VM.GlobalState?
	// TODO: do we need this?
	// self.VM.GlobalState.Set(actorID, actor)

	// Store the mappings of address to actor ID.
	self.AddressMap_[addr] = actorID
	self.IDMap_[actorID] = addr

	// TODO: adjust this to be proper state setting.
	rt.State().StateTree().ActorStates()[addr] = actor.State_ // atm it's nil

	return addr
}

func (self *InitActor_I) _assignNextID() actor.ActorID {
	actorID := self.NextID_
	self.NextID_++
	return actorID
}

func (self *InitActor_I) GetActorIDForAddress(address addr.Address) actor.ActorID {
	return self.AddressMap_[address]
}

// TODO: derive this OR from a union type
func (_ *InitActor_I) _isSingletonActor(codeCID actor.CodeCID) bool {
	return true
	// TODO: uncomment this
	// return codeCID == StorageMarketActor ||
	// 	codeCID == StoragePowerActor ||
	// 	codeCID == CronActor ||
	// 	codeCID == InitActor
}

// TODO: derive this OR from a union type
func (_ *InitActor_I) _isBuiltinActor(codeCID actor.CodeCID) bool {
	return true
	// TODO: uncomment this
	// return codeCID == StorageMarketActor ||
	// 	codeCID == StoragePowerActor ||
	// 	codeCID == CronActor ||
	// 	codeCID == InitActor ||
	// 	codeCID == StorageMinerActor ||
	// 	codeCID == PaymentChannelActor
}
