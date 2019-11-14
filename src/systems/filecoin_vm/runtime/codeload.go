package runtime

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"

func loadActorCode(codeID actor.CodeID) (ActorCode, error) {

	panic("TODO")
	// TODO: resolve circular dependency

	// // load the code from StateTree.
	// // TODO: this is going to be enabled in the future.
	// // code, err := loadCodeFromStateTree(input.InTree, codeCID)
	// return staticActorCodeRegistry.LoadActor(codeCID)
}
