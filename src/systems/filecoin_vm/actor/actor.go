package actor

import filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
import ipld "github.com/filecoin-project/specs/libraries/ipld"

const (
	MethodSend        = MethodNum(0)
	MethodConstructor = MethodNum(1)
	MethodCron        = MethodNum(2)

	// TODO: remove this once canonical method numbers are finalized
	MethodPlaceholder = MethodNum(-(1 << 30))
)

func (st *ActorState_I) CID() ipld.CID {
	panic("TODO")
}

// Note: may be nil if actor has no public key
func (st *ActorState_I) GetSignaturePublicKey() filcrypto.PublicKey {
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
