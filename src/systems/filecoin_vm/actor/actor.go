package actor

import (
	actors "github.com/filecoin-project/specs/actors"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO
var TODO = util.TODO

type Serialization = util.Serialization

const (
	MethodSend        = actors.MethodNum(0)
	MethodConstructor = actors.MethodNum(1)

	// TODO: remove this once canonical method numbers are finalized
	MethodPlaceholder = actors.MethodNum(1 << 30)
)

func (st *ActorState_I) CID() ipld.CID {
	panic("TODO")
}

func (id *CodeID_I) IsBuiltin() bool {
	switch id.Which() {
	case CodeID_Case_Builtin:
		return true
	default:
		panic("Actor code ID case not supported")
	}
}

func (id *CodeID_I) IsSingleton() bool {
	if !id.IsBuiltin() {
		return false
	}

	for _, a := range []BuiltinActorID{
		BuiltinActorID_Init,
		BuiltinActorID_Cron,
		BuiltinActorID_Init,
		BuiltinActorID_StoragePower,
		BuiltinActorID_StorageMarket,
	} {
		if id.As_Builtin() == a {
			return true
		}
	}

	for _, a := range []BuiltinActorID{
		BuiltinActorID_Account,
		BuiltinActorID_PaymentChannel,
		BuiltinActorID_StorageMiner,
	} {
		if id.As_Builtin() == a {
			return false
		}
	}

	panic("Actor code ID case not supported")
}

func BuiltinActorID_SignableTypes() []BuiltinActorID {
	IMPL_TODO() // TODO: Update for MultiSig actors

	return []BuiltinActorID{
		BuiltinActorID_Account,
	}
}

func (id *CodeID_I) IsSignable() bool {
	if !id.IsBuiltin() {
		return false
	}

	for _, a := range BuiltinActorID_SignableTypes() {
		if id.As_Builtin() == a {
			return true
		}
	}

	return false
}

func (x ActorSubstateCID) Ref() *ActorSubstateCID {
	return &x
}
