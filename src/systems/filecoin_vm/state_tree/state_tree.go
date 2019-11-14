package state_tree

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func (inTree *StateTree_I) WithActorForAddress(a addr.Address) (StateTree, actor.Actor) {
	var err error
	var actor actor.Actor
	var compTree StateTree

	actor = inTree.GetActor(a)
	if actor != nil {
		return inTree, actor // done
	}

	if !a.IsKeyType() { // BLS or Secp
		return inTree, nil // not a key type, done.
	}

	compTree, actor, err = inTree.Impl().WithNewAccountActor(a)
	if err != nil {
		return inTree, nil
	}

	return compTree, actor
}

func (st *StateTree_I) GetActor(a addr.Address) actor.Actor {
	panic("TODO")
}

func (st *StateTree_I) Balance(a addr.Address) actor.TokenAmount {
	panic("TODO")
}

func (st *StateTree_I) WithActorSubstate(a addr.Address, actorState actor.ActorSubstateCID) (StateTree, error) {
	panic("TODO")
}

func (st *StateTree_I) WithActorSystemState(a addr.Address, actorState actor.ActorSystemStateCID) (StateTree, error) {
	panic("TODO")
}

func (st *StateTree_I) WithFundsTransfer(from addr.Address, to addr.Address, amount actor.TokenAmount) (StateTree, error) {
	panic("TODO")
}

func (st *StateTree_I) WithNewAccountActor(a addr.Address) (StateTree, actor.Actor, error) {
	panic("TODO")
}

func (st *StateTree_I) WithIncrementedCallSeqNum(a addr.Address) (StateTree, error) {
	panic("TODO")
}

/*
TODO: finish


func treeIncrementActorSeqNo(inTree StateTree, a actor.Actor) (outTree StateTree) {
	panic("todo")
}

func treeDeductFunds(inTree StateTree, a actor.Actor, amt actor.TokenAmount) (outTree StateTree) {
	// TODO: turn this into a single transfer call.
	panic("todo")
}

func treeDepositFunds(inTree StateTree, a actor.Actor, amt actor.TokenAmount) (outTree StateTree) {
	// TODO: turn this into a single transfer call.
	panic("todo")
}

func treeGetOrCreateAccountActor(inTree StateTree, a addr.Address) (outTree StateTree, _ actor.Actor, err error) {

	toActor := inTree.GetActor(a)
	if toActor != nil { // found
		return inTree, toActor, nil
	}

	switch a.Type().Which() {
	case addr.Address_Type_Case_BLS:
		return treeNewBLSAccountActor(inTree, a)
	case addr.Address_Type_Case_Secp256k1:
		return treeNewSecp256k1AccountActor(inTree, a)
	case addr.Address_Type_Case_ID:
		return inTree, nil, errors.New("no actor with given ID")
	case addr.Address_Type_Case_Actor:
		return inTree, nil, errors.New("no such actor")
	default:
		return inTree, nil, errors.New("unknown address type")
	}
}

func treeNewBLSAccountActor(inTree StateTree, addr addr.Address) (outTree StateTree, _ actor.Actor, err error) {
	panic("todo")
}

func treeNewSecp256k1AccountActor(inTree StateTree, addr addr.Address) (outTree StateTree, _ actor.Actor, err error) {
	panic("todo")
}
*/
