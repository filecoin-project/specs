package crypto

import util "github.com/filecoin-project/specs/util"

func (self *VRFResult_I) ValidateSyntax() bool {
	panic("TODO")
	return false
}

func (self *VRFResult_I) Verify(input util.Bytes, pk VRFPublicKey) bool {
	// return new(BLS).Verify(self.Output, pk.(*BLSPublicKey), input)
	return false
}

func (self *VRFResult_I) MaxValue() util.Bytes {
	panic("")
	// return new(BLS).MaxSigValue()
}

func (self *VRFKeyPair_I) Generate(input util.Bytes) VRFResult {
	// sig := new(BLS).Sign(input, self.SecretKey)
	var blsSig util.Bytes
	ret := &VRFResult_I{
		Proof_: blsSig,
		Output_: SHA256(blsSig),
	}
	return ret
}
