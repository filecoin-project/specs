package impl

import (
	abi "github.com/filecoin-project/specs-actors/actors/abi"
	vmr "github.com/filecoin-project/specs-actors/actors/runtime"
)

func loadActorCode(codeID abi.ActorCodeID) (vmr.ActorCode, error) {

	panic("TODO")
	// TODO: resolve circular dependency

	// // load the code from StateTree.
	// // TODO: this is going to be enabled in the future.
	// // code, err := loadCodeFromStateTree(input.InTree, codeCID)
	// return staticActorCodeRegistry.LoadActor(codeCID)
}
