package crypto

import (
	"bytes"
	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	util "github.com/filecoin-project/specs/util"
)

type DomainSeparationTag int
type Randomness = abi.Randomness
type Serialization = util.Serialization

const (
	DomainSeparationTag_TicketDrawing DomainSeparationTag = 1 + iota
	DomainSeparationTag_TicketProduction
	DomainSeparationTag_ElectionPoStChallengeSeed
	DomainSeparationTag_SurprisePoStChallengeSeed
	DomainSeparationTag_SurprisePoStSelectMiners
	DomainSeparationTag_SurprisePoStSampleSectors
)

// Derive a random byte string from a domain separation tag and the appropriate values
func DeriveRandWithMinerAddr(tag DomainSeparationTag, tix abi.RandomnessSeed, minerAddr addr.Address) Randomness {
	var addrBuf bytes.Buffer
	err := minerAddr.MarshalCBOR(&addrBuf)
	util.Assert(err == nil)

	return _deriveRandInternal(tag, tix, addrBuf.Bytes())
}

func DeriveRandWithEpoch(tag DomainSeparationTag, tix abi.RandomnessSeed, epoch int) Randomness {
	return _deriveRandInternal(tag, tix, BigEndianBytesFromInt(epoch))
}

func DeriveRand(tag DomainSeparationTag, tix abi.RandomnessSeed) Randomness {
	return _deriveRandInternal(tag, tix, nil)
}

func _deriveRandInternal(tag DomainSeparationTag, randSeed abi.RandomnessSeed, s Serialization) Randomness {
	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, util.Bytes(randSeed)...)
	buffer = append(buffer, s...)
	return Randomness(SHA256(buffer))
}
