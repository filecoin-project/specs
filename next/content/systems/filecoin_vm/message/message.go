package message

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

type Serialization = util.Serialization

// The maximum serialized size of a SignedMessage.
const MessageMaxSize = 32 * 1024

func SignedMessage_Make(message UnsignedMessage, signature filcrypto.Signature) SignedMessage {
	return &SignedMessage_I{
		Message_:   message,
		Signature_: signature,
	}
}

func Sign(message UnsignedMessage, keyPair filcrypto.SigKeyPair) (SignedMessage, error) {
	sig, err := filcrypto.Sign(keyPair, util.Bytes(Serialize_UnsignedMessage(message)))
	if err != nil {
		return nil, err
	}
	return SignedMessage_Make(message, sig), nil
}

func SignatureVerificationError() error {
	IMPL_FINISH()
	panic("")
}

func Verify(message SignedMessage, publicKey filcrypto.PublicKey) (UnsignedMessage, error) {
	m := util.Bytes(Serialize_UnsignedMessage(message.Message()))
	sigValid, err := filcrypto.Verify(publicKey, message.Signature(), m)
	if err != nil {
		return nil, err
	}
	if !sigValid {
		return nil, SignatureVerificationError()
	}
	return message.Message(), nil
}

func (x *GasAmount_I) Add(y GasAmount) GasAmount {
	IMPL_FINISH()
	panic("")
}

func (x *GasAmount_I) Subtract(y GasAmount) GasAmount {
	IMPL_FINISH()
	panic("")
}

func (x *GasAmount_I) SubtractIfNonnegative(y GasAmount) (ret GasAmount, ok bool) {
	ret = x.Subtract(y)
	ok = true
	if ret.LessThan(GasAmount_Zero()) {
		ret = x
		ok = false
	}
	return
}

func (x *GasAmount_I) LessThan(y GasAmount) bool {
	IMPL_FINISH()
	panic("")
}

func (x *GasAmount_I) Equals(y GasAmount) bool {
	IMPL_FINISH()
	panic("")
}

func (x *GasAmount_I) Scale(count int) GasAmount {
	IMPL_FINISH()
	panic("")
}

func GasAmount_Affine(b GasAmount, x int, m GasAmount) GasAmount {
	return b.Add(m.Scale(x))
}

func GasAmount_Zero() GasAmount {
	return GasAmount_FromInt(0)
}

func GasAmount_FromInt(x int) GasAmount {
	IMPL_FINISH()
	panic("")
}

func GasAmount_SentinelUnlimited() GasAmount {
	// Amount of gas larger than any feasible execution; meant to indicated unlimited gas
	// (e.g., for builtin system method invocations).
	return GasAmount_FromInt(1).Scale(1e9).Scale(1e9) // 10^18
}
