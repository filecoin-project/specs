package codeGen

type SignatureType int

const (
	SignatureType_Secp256k1 SignatureType = 0
	SignatureType_BLS       SignatureType = 1
)

type Signature interface {
	Type() SignatureType
}

type SigKeyPair interface {
	Sign(message []byte) Signature
	Verify(message []byte, signature Signature) bool
}
