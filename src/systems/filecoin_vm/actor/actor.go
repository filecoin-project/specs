package actor

import ipld "github.com/filecoin-project/specs/libraries/ipld"

const (
	MethodSend        = MethodNum(0)
	MethodConstructor = MethodNum(1)
	MethodCron        = MethodNum(2)

	MethodGetUnsealedCIDForDealIDs = MethodNum(99)
)

func (st *ActorState_I) CID() ipld.CID {
	panic("TODO")
}
