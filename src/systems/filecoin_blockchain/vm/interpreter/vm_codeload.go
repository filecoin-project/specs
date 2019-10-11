package interpreter

import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import runtime "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime"

func loadActorCode(input runtime.InvocInput, codeCID actor.CodeCID) (ActorCode, error) {

	// load the code from StateTree.
	// TODO: this is going to be enabled in the future.
	// code, err := loadCodeFromStateTree(input.InTree, codeCID)
	return staticActorCodeRegistry.LoadActor(codeCID)
}
