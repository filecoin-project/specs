package filproofs

import util "github.com/filecoin-project/specs/util"

////////////////////////////////////////////////////////////////////////////////
/// Generic Hashing

/// Binary hash compression.
// BinaryHash<T>
func BinaryHash_T(left []byte, right []byte) util.T {
	var preimage = append(left, right...)
	return HashBytes_T(preimage)
}

func TernaryHash_T(a []byte, b []byte, c []byte) util.T {
	var preimage = append(a, append(b, c...)...)
	return HashBytes_T(preimage)
}

// BinaryHash<PedersenHash>
func BinaryHash_PedersenHash(left []byte, right []byte) PedersenHash {
	return PedersenHash{}
}

func TernaryHash_PedersenHash(a []byte, b []byte, c []byte) PedersenHash {
	return PedersenHash{}
}

// BinaryHash<SHA256Hash>
func BinaryHash_SHA256Hash(left []byte, right []byte) SHA256Hash {
	result := SHA256Hash{}
	return trimToFr32(result)
}

func TernaryHash_SHA256Hash(a []byte, b []byte, c []byte) SHA256Hash {
	return SHA256Hash{}
}

////////////////////////////////////////////////////////////////////////////////

/// Digest
// HashBytes<T>
func HashBytes_T(data []byte) util.T {
	return util.T{}
}

// HashBytes<PedersenHash>
func HashBytes_PedersenHash(data []byte) PedersenHash {
	return PedersenHash{}
}

// HashBytes<SHA256Hash.
func HashBytes_SHA256Hash(data []byte) SHA256Hash {
	// Digest is truncated to 254 bits.
	result := SHA256Hash{}

	return result
}

////////////////////////////////////////////////////////////////////////////////

func DigestSize_T() int {
	panic("Unspecialized")
}

func DigestSize_PedersenHash() int {
	return 32
}

func DigestSize_SHA256Hash() int {
	return 32
}
