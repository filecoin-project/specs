package sysactors

import msg "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import st "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/state_tree"
import util "github.com/filecoin-project/specs/util"

func (self *InitActor) InitConstructor() {

}

func (self *InitActor) Exec(codeCID actor.CodeCID, method actor.MethodNum) addr.Address {
    // Make sure that only the actors defined in the spec can be launched.
    if !self._isBuiltinActor(code) {
        self.VM.Fatal("cannot launch actor instance that is not a builtin actor")
    }

    // Get the actor ID for this actor.
    actorID := self._assignNextID()

    // Ensure that singeltons can only be launched once.
    // TODO: do we want to enforce this? If so how should actors be marked as such?
    if self._isSingletonActor(codeCID) {
        Fatal("cannot launch another actor of this type")
    }

    // This generates a unique address for this actor that is stable across message
    // reordering
    // TODO: where do `creator` and `nonce` come from?
    addr := self.VM.ComputeActorAddress(creator, nonce)

    // Set up the actor itself
    actor := actor.Actor{
        CodeCid: codeCID,
        Balance: msg.Value,
        State:   nil,
        Nonce:   0,
    }

    // The call to the actors constructor will set up the initial state
    // from the given parameters, setting `actor.Head` to a new value when successfull.
    // TODO: can constructors fail?
    actor.Constructor(params)

    // TODO: where is this VM.GlobalState?
    // TODO: do we need this?
    // self.VM.GlobalState.Set(actorID, actor)

    // Store the mappings of address to actor ID.
    self.AddressMap[addr] = actorID
    self.IDMap[actorID] = addr

    return addr
}

func (a *InitActor) _assignNextID() ActorID {
    actorID := self.State.NextID
    self.State.NextID++
    return actorID
}

func (a *InitActor) GetActorIDForAddress(address Address) ActorID {
    // go code
}

// TODO: derive this OR from a union type
func (a *InitActor) _isSingletonActor(code CID) bool {
    return code == StorageMarketActor ||
        code == StoragePowerActor ||
        code == CronActor ||
        code == InitActor
}

// TODO: derive this OR from a union type
func (a *InitActor) _isBuiltinActor(code CID) bool {
    return code == StorageMarketActor ||
        code == StoragePowerActor ||
        code == CronActor ||
        code == InitActor ||
        code == StorageMinerActor ||
        code == PaymentChannelActor
}
