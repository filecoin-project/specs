package state_tree

import (
	addr "github.com/filecoin-project/go-address"
	actor "github.com/filecoin-project/specs/actors"
	"github.com/filecoin-project/specs/actors/abi"
	"github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
)

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO

func (st *StateTree_I) RootCID() cid.Cid {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) GetActor(a addr.Address) (actor.ActorState, bool) {
	as, found := st.ActorStates()[a]
	return as, found
}

func (st *StateTree_I) GetActorCodeID_Assert(a addr.Address) abi.ActorCodeID {
	ret, found := st.GetActor(a)
	Assert(found)
	return ret.CodeID()
}

func (st *StateTree_I) WithActorSubstate(a addr.Address, actorState actor.ActorSubstateCID) (StateTree, error) {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithDeleteActorSystemState(a addr.Address) StateTree {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithActorSystemState(a addr.Address, actorState actor.ActorSystemStateCID) (StateTree, error) {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithFundsTransfer(from addr.Address, to addr.Address, amount abi.TokenAmount) (StateTree, error) {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithNewAccountActor(a addr.Address) (StateTree, actor.ActorState, error) {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithIncrementedCallSeqNum(a addr.Address) (StateTree, error) {
	IMPL_FINISH()
	panic("")
}

func (st *StateTree_I) WithIncrementedCallSeqNum_Assert(a addr.Address) StateTree {
	ret, err := st.WithIncrementedCallSeqNum(a)
	if err != nil {
		panic("Error incrementing actor call sequence number")
	}
	return ret
}

/*
TODO: finish


func treeIncrementActorSeqNo(inTree StateTree, a actor.Actor) (outTree StateTree) {
	panic("todo")
}

func treeDeductFunds(inTree StateTree, a actor.Actor, amt abi.TokenAmount) (outTree StateTree) {
	// TODO: turn this into a single transfer call.
	panic("todo")
}

func treeDepositFunds(inTree StateTree, a actor.Actor, amt abi.TokenAmount) (outTree StateTree) {
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
