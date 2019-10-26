package state_tree

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
