package actors

import (
	cid "github.com/ipfs/go-cid"
)

func (st *ActorState_I) CID() cid.Cid {
	panic("TODO")
}

func (x ActorSubstateCID) Ref() *ActorSubstateCID {
	return &x
}
