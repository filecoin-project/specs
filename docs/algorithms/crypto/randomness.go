package crypto

import (
	util "github.com/filecoin-project/specs/util"
)

type Randomness = util.Randomness
type Serialization = util.Serialization

// Derive a random byte string from a domain separation tag and an arbitrary
// serializable object.
//
// Note: to produce values of type Serialization, use the auto-generated method prototypes
//   Serialize_T(T) Serialization
// for each type T defined in .id.
//
// In order to derive randomness from a collection of objects, rather than just a single
// object, define a struct at the .id level that contains those objects as member fields.
// This will then cause a Serialize_*() method to be generated for the struct type.
func (tag DomainSeparationTag) DeriveRand(s Serialization) Randomness {
	return _deriveRandInternal(tag, s, -1)
}

// As in DeriveRand(), but additionally accepts an index into the implicit pseudorandom stream.
// Index must be strictly positive.
func (tag DomainSeparationTag) DeriveRandWithIndex(s Serialization, index int) Randomness {
	if index <= 0 {
		panic("DeriveRandWithIndex only accepts indices > 0")
	}
	return _deriveRandInternal(tag, s, index)
}

func _deriveRandInternal(tag DomainSeparationTag, s Serialization, index int) Randomness {
	buffer := []byte{}
	buffer = append(buffer, util.IntToBytesLittleEndian(int(tag))...)
	buffer = append(buffer, util.IntToBytesLittleEndian(int(index))...)
	buffer = append(buffer, Bytes(s)...)
	ret := SHA256(buffer)
	return Randomness(ret)
}

type DomainSeparationTag int

const (
	DomainSeparationTag_TicketDrawing DomainSeparationTag = 1 + iota
	DomainSeparationTag_TicketProduction
	DomainSeparationTag_PoSt
	// ...
)
