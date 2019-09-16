package codeGen

type CID interface {
	Equals(CID) bool
}

type CID_I struct {
	rawValue []byte
}

func (x *CID_I) Equals(y *CID_I) bool {
	return CompareBytesStrict(x.rawValue, y.rawValue) == 0
}
