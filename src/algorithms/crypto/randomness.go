package crypto

import util "github.com/filecoin-project/specs/util"
import addr "github.com/filecoin-project/go-address"
import abi "github.com/filecoin-project/specs/actors/abi"

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
	buffer := _deriveRandInternal(tag, tix, -1)
	var serializedAddr abi.Bytes
	util.IMPL_FINISH() // serialize the address
	buffer = append(buffer, serializedAddr...)
	ret := SHA256(buffer)
	return Randomness(_deriveRandInternal(tag, tix, -1, minerAddr)))
}

func DeriveRandWithEpoch(tag DomainSeparationTag, tix abi.RandomnessSeed, epoch int) Randomness {
	buffer = append(buffer, BigEndianBytesFromInt(epoch)...)
	return Randomness(_deriveRandInternal(tag, tix, -1, epoch))
}

func _deriveRandInternal(tag DomainSeparationTag, randSeed abi.RandomnessSeed, index int, s Serialization) util.Bytes {
	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, BigEndianBytesFromInt(int(index))...)
	buffer = append(buffer, util.Bytes(randSeed)...)
	buffer = append(buffer, s...)
	return SHA256(buffer)
}
