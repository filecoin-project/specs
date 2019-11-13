package crypto

// Sign generates a proof that miner `M` generate message `m`
//
// Out:
//    sig - a signature
//    err - a standard error message indicating any process issues
// In:
//    m - a series of bytes representing a message to be signed
//
func Sign(keyPair SigKeyPair, m Message) (Signature, error) {
	panic("TODO")
}

// Verify validates the statement: only `M` could have generated `sig`
// given the validator has a message `m`, a signature `sig`, and a
// public key `pk`.
//
// Out:
//    valid - a boolean value indicating the signature is valid
//    err - a standard error message indicating any process issues
// In:
//    pk - the public key belonging to the signer `M`
//    m - a series of bytes representing the signed message
//
func Verify(pk PublicKey, sig Signature, m Message) (valid bool, err error) {
	panic("TODO")
}
