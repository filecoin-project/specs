package codeGen

type VRFPublicKey interface{}

type VRFSecretKey interface{}

type VRFKeyPair interface {
	SecretKey() VRFSecretKey
	PublicKey() VRFPublicKey
}

type VRFResult interface {
	implements_VRFResult() VRFResult

	PublicKey() VRFPublicKey
	ToBytesInternal() []byte
}

func VRFEval(k VRFKeyPair, x []byte) VRFResult {
	panic("TODO")
}

////////////////////
// Implementation //
////////////////////

type VRFResultI struct {
	publicKey VRFPublicKey
	rawValue  []byte
}

func (r *VRFResultI) PublicKey() VRFPublicKey {
	return r.publicKey
}

func (r *VRFResultI) ToBytesInternal() []byte {
	return r.rawValue
}

func (r *VRFResultI) implements_VRFResult() VRFResult {
	return r
}
