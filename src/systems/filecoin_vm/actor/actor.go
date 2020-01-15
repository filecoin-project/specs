package actor

import (
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
)

var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO
var TODO = util.TODO

type Serialization = util.Serialization

func (st *ActorState_I) CID() cid.Cid {
	panic("TODO")
}
