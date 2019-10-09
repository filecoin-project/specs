package crypto

import util "github.com/filecoin-project/specs/util"

func (self *BLS_I) Verify(input util.Bytes, pk PublicKey, sig util.Bytes) (valid bool, err error) {
	blsPk := pk.(*BLSPublicKey)
	// 1. Verify public key according to string_to_curve section 2.6.2.1. in
	// 	https://tools.ietf.org/html/draft-boneh-bls-signature-00#page-12
	// 2. Verify signature according to section 2.3
	// 	https://tools.ietf.org/html/draft-boneh-bls-signature-00#page-8
}

func (self *BLS) Sign(input util.Bytes, sk *PrivateKey) bool {
	panic("see 2.3 in https://tools.ietf.org/html/draft-boneh-bls-signature-00#page-8")

}

func (self *BLS) Aggregate(sig2 util.Bytes) util.Bytes {
	panic("see 2.5 in https://tools.ietf.org/html/draft-boneh-bls-signature-00#page-8")
}

func (self *BLS) VerifyAggregate(messages []util.Bytes, aggPk PublicKey, aggSig util.Bytes) (valid bool, err error) {
	panic("see 2.5.2 in https://tools.ietf.org/html/draft-boneh-bls-signature-00#page-9")
}
