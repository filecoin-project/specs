package codeGen

// VDFResult _
type VDFResult interface {
	Value() []byte
	Proof() []byte
}

// VDFEval _
func VDFEval(x []byte) VDFResult {
	panic("TODO")
	// return &VRFResultI{}
}
