package crypto

import util "github.com/filecoin-project/specs/util"

func (self *VRFResult_I) Verify(input util.Bytes, pk VRFPublicKey) (valid bool, err error) {
	return new(BLS).Verify(self.Output, pk.(*BLSPublicKey), input)
}

func (self *VRFKeyPair) Generate(input util.Bytes) VRFResult_I {
	sig := new(BLS).Sign(input, self.SecretKey)
	return VRFResult_I{
		Output: sig,
	}
}
