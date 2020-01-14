package actor

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
)

var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO
var TODO = util.TODO

type Serialization = util.Serialization

const (
	MethodSend        = abi.MethodNum(0)
	MethodConstructor = abi.MethodNum(1)

	// TODO: remove this once canonical method numbers are finalized
	MethodPlaceholder = abi.MethodNum(1 << 30)
)

func (st *ActorState_I) CID() cid.Cid {
	panic("TODO")
}

func (x ActorSubstateCID) Ref() *ActorSubstateCID {
	return &x
}
