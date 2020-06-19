package crypto

import (
	util "github.com/filecoin-project/specs/util"
	"golang.org/x/crypto/blake2b"
)

func (self *VRFResult_I) ValidateSyntax() bool {
	panic("TODO")
	return false
}

func (self *VRFResult_I) Verify(input util.Bytes, pk VRFPublicKey) bool {
	// return new(BLS).Verify(self.Proof, pk.(*BLSPublicKey), input)
	return false
}

func (self *VRFResult_I) MaxValue() util.Bytes {
	panic("")
	// return new(BLS).MaxSigValue()
}

func (self *VRFKeyPair_I) Generate(input util.Bytes) VRFResult {
	// sig := new(BLS).Sign(input, self.SecretKey)
	var blsSig util.Bytes

	digest := blake2b.Sum256(blsSig)
	ret := &VRFResult_I{
		Proof_:  blsSig,
		Digest_: digest[:],
	}
	return ret
}
