type Signature union {
    | Secp256k1Signature 0
    | Bls12_381Signature 1
} // representation byteprefix

type Secp256k1Signature Bytes
type Bls12_381Signature Bytes
