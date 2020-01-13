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
	return Randomness(ret)
}

func DeriveRandWithEpoch(tag DomainSeparationTag, tix abi.RandomnessSeed, epoch int) Randomness {
	buffer := _deriveRandInternal(tag, tix, -1)
	buffer = append(buffer, BigEndianBytesFromInt(epoch)...)
	ret := SHA256(buffer)
	return Randomness(ret)
}

func _deriveRandInternal(tag DomainSeparationTag, s abi.RandomnessSeed, index int) util.Bytes {
	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, BigEndianBytesFromInt(int(index))...)
	buffer = append(buffer, util.Bytes(s)...)
	return buffer
}
