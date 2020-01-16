package crypto

type SigType int64

const (
	SigTypeUnknown = SigType(-1)

	SigTypeSecp256k1 = SigType(iota)
	SigTypeBLS
)

type Signature struct {
	Type SigType
	Data []byte
}
