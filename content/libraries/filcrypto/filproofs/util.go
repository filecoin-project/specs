package filproofs

import (
	"encoding/binary"
	big "math/big"

	abi "github.com/filecoin-project/specs-actors/actors/abi"
	file "github.com/filecoin-project/specs/systems/filecoin_files/file"
	util "github.com/filecoin-project/specs/util"
)

// Utilities

func reverse(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}

func bigIntFromLittleEndianBytes(bytes []byte) *big.Int {
	reverse(bytes)
	return new(big.Int).SetBytes(bytes)
}

func bigIntFromBigEndianBytes(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// size is number of bytes to return
func littleEndianBytesFromBigInt(z *big.Int, size int) []byte {
	bytes := z.Bytes()[0:size]
	reverse(bytes)

	return bytes
}

// size is number of bytes to return
func bigEndianBytesFromBigInt(z *big.Int, size int) []byte {
	return z.Bytes()[0:size]
}

func littleEndianBytesFromInt(n int, size int) []byte {
	z := new(big.Int)
	z.SetInt64(int64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func bigEndianBytesFromInt(n int, size int) []byte {
	z := new(big.Int)
	z.SetInt64(int64(n))
	return bigEndianBytesFromBigInt(z, size)
}

func littleEndianBytesFromUInt(n UInt, size int) []byte {
	z := new(big.Int)
	z.SetUint64(uint64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func bigEndianBytesFromUInt(n UInt, size int) []byte {
	z := new(big.Int)
	z.SetUint64(uint64(n))
	return bigEndianBytesFromBigInt(z, size)
}

func AsBytes_T(t util.T) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func AsBytes_UnsealedSectorCID(cid abi.UnsealedSectorCID) []byte {
	panic("Unimplemented for UnsealedSectorCID")

	return []byte{}
}

func AsBytes_SealedSectorCID(CID abi.SealedSectorCID) []byte {
	panic("Unimplemented for SealedSectorCID")

	return []byte{}
}

func AsBytes_PieceCID(CID abi.PieceCID) []byte {
	panic("Unimplemented for PieceCID")

	return []byte{}
}

func fromBytes_T(_ interface{}) util.T {
	panic("Unimplemented for T")
	return util.T{}
}

func fromBytes_PieceCID(_ interface{}) abi.PieceCID {
	panic("Unimplemented for PieceCID")
}

func isPow2(n int) bool {
	return n != 0 && n&(n-1) == 0
}

// FIXME: This does not belong in filproofs, and no effort is being made to ensure it has any particular properties.
func RandomInt(randomness util.Randomness, nonce int, limit *big.Int) *big.Int {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, uint64(nonce))
	input := randomness
	input = append(input, nonceBytes...)
	ranHash := HashBytes_SHA256Hash(input[:])
	hashInt := bigIntFromLittleEndianBytes(ranHash)
	num := hashInt.Mod(hashInt, limit)
	return num
}

////////////////////////////////////////////////////////////////////////////////

// Destructively trim data so most significant two bits of last byte are 0.
// This ensure data interpreted as little-endian will not exceed a field with 254-bit capacity.
// NOTE: 254 bits is the capacity of BLS12-381, but other curves with ~32-byte field elements
// may have a different capacity. (Example: BLS12-377 has a capacity of 252 bits.)
func trimToFr32(data []byte) []byte {
	util.Assert(len(data) == 32)
	data[31] &= 0x3f // 0x3f = 0b0011_1111
	return data
}

func UnsealedSectorCID(h SHA256Hash) abi.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h PedersenHash) abi.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func Commitment_UnsealedSectorCID(cid abi.UnsealedSectorCID) Commitment {
	panic("not implemented -- re-arrange bits")
}

func Commitment_SealedSectorCID(cid abi.SealedSectorCID) Commitment {
	panic("not implemented -- re-arrange bits")
}

func ComputeDataCommitment(data []byte) (abi.UnsealedSectorCID, file.Path) {
	// TODO: make hash parameterizable
	hash, path := BuildTree_SHA256Hash(data)
	return UnsealedSectorCID(hash), path
}

// Compute CommP or CommD.
func ComputeUnsealedSectorCID(data []byte) (abi.UnsealedSectorCID, file.Path) {
	// TODO: check that len(data) > minimum piece size and is a power of 2.
	hash, treePath := BuildTree_SHA256Hash(data)
	return UnsealedSectorCID(hash), treePath
}
