package actors

import (
	cid "github.com/ipfs/go-cid"
)

type ActorSubstateCID cid.Cid

func (x ActorSubstateCID) Ref() ActorSubstateCID {
	return x
}
