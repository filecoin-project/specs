package filproofs

import (
	"golang.org/x/crypto/blake2s"
)

// TODO/FIXME: Update to use uint64, not uint32.

func Permute(numElements uint32, index uint32, keys []uint32) uint32 {
	u := encode(numElements, index, keys)
	for u >= numElements {
		u = encode(numElements, u, keys)
	}

	return u
}

func InvertPermute(numElements uint32, index uint32, keys []uint32) uint32 {
	u := decode(numElements, index, keys)
	for u >= numElements {
		u = decode(numElements, u, keys)
	}

	return u
}

func encode(numElements uint32, index uint32, keys []uint32) uint32 {
	// find nextPow4
	nextPow4 := uint32(4)
	log4 := uint32(1)
	for nextPow4 < numElements {
		nextPow4 *= 4
		log4++
	}

	// left and right masks
	leftMask := ((uint32(1) << log4) - 1) << log4
	rightMask := (uint32(1) << log4) - 1
	halfBits := log4

	left := ((index & leftMask) >> halfBits)
	right := (index & rightMask)

	for i := 0; i < 4; i++ {
		left, right = right, left^feistel(right, keys[i], rightMask)
	}

	return (left << halfBits) | right
}

func decode(numElements uint32, index uint32, keys []uint32) uint32 {

	// find nextPow4
	nextPow4 := uint32(4)
	log4 := uint32(1)
	for nextPow4 < numElements {
		nextPow4 *= 4
		log4++
	}

	// left and right masks
	leftMask := ((uint32(1) << log4) - 1) << log4
	rightMask := (uint32(1) << log4) - 1
	halfBits := log4

	left := ((index & leftMask) >> halfBits)
	right := (index & rightMask)

	for i := 3; i > -1; i-- {
		left, right = right^feistel(left, keys[i], rightMask), left
	}

	return (left << halfBits) | right
}

func feistel(right uint32, key uint32, rightMask uint32) uint32 {
	var data [8]byte
	data[0] = byte(right >> 24)
	data[1] = byte(right >> 16)
	data[2] = byte(right >> 8)
	data[3] = byte(right)

	data[4] = byte(key >> 24)
	data[5] = byte(key >> 16)
	data[6] = byte(key >> 8)
	data[7] = byte(key)

	hash := blake2s.Sum256(data[:])

	r := uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
	return r & rightMask
}
