type Address union {
    | AddressId 0
    | AddressSecp256k1 1
    | AddressActor 2
    | AddressBLS12_381 3
} // representation byteprefix

// ID
type AddressId UInt

// Blake2b-160 Hash
type AddressSecp256k1 Bytes

// Blake2b-160 Hash
type AddressActor Bytes

// 48 byte PublicKey
type AddressBLS12_381 Bytes
