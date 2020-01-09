package crypto

import util "github.com/filecoin-project/specs/util"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
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
func DeriveRandWithMinerAddr(tag DomainSeparationTag, tix abi.Randomness, minerAddr addr.Address) Randomness {
	buffer := _deriveRandInternal(tag, tix, -1)
	buffer = append(buffer, addr.Serialize_Address_Compact(minerAddr)...)
	ret := SHA256(buffer)
	return Randomness(ret)
}

func DeriveRandWithEpoch(tag DomainSeparationTag, tix abi.Randomness, epoch int) Randomness {
	buffer := _deriveRandInternal(tag, tix, -1)
	buffer = append(buffer, LittleEndianBytesFromInt(epoch)...)
	ret := SHA256(buffer)
	return Randomness(ret)
}

func _deriveRandInternal(tag DomainSeparationTag, tix abi.Randomness, index int) util.Bytes {
	buffer := []byte{}
	buffer = append(buffer, LittleEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, LittleEndianBytesFromInt(int(index))...)
	buffer = append(buffer, util.Bytes(tix)...)
	return buffer
}
