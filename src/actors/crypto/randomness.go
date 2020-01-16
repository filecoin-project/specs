package crypto

import (
	"bytes"
	"encoding/binary"
	"math"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	autil "github.com/filecoin-project/specs/actors/util"
)

// Specifies a domain for randomness generation.
type DomainSeparationTag int

const (
	DomainSeparationTag_TicketDrawing DomainSeparationTag = 1 + iota
	DomainSeparationTag_TicketProduction
	DomainSeparationTag_ElectionPoStChallengeSeed
	DomainSeparationTag_SurprisePoStChallengeSeed
	DomainSeparationTag_SurprisePoStSelectMiners
	DomainSeparationTag_SurprisePoStSampleSectors
)

// Derive a random byte string from a domain separation tag and the appropriate values
func DeriveRandWithMinerAddr(tag DomainSeparationTag, tix abi.RandomnessSeed, minerAddr addr.Address) abi.Randomness {
	var addrBuf bytes.Buffer
	err := minerAddr.MarshalCBOR(&addrBuf)
	autil.AssertNoError(err)

	return _deriveRandInternal(tag, tix, -1, addrBuf.Bytes())
}

func DeriveRandWithEpoch(tag DomainSeparationTag, tix abi.RandomnessSeed, epoch int) abi.Randomness {
	return _deriveRandInternal(tag, tix, -1, BigEndianBytesFromInt(epoch))
}

func _deriveRandInternal(tag DomainSeparationTag, randSeed abi.RandomnessSeed, index int, s []byte) abi.Randomness {
	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int(tag))...)
	buffer = append(buffer, BigEndianBytesFromInt(int(index))...)
	buffer = append(buffer, abi.Bytes(randSeed)...)
	buffer = append(buffer, s...)
	return abi.Randomness(SHA256(buffer))
}

func RandomInt(randomness abi.Randomness, nonce int, limit int) int {
	nonceBytes := BigEndianBytesFromInt(nonce)
	input := randomness
	input = append(input, nonceBytes...)
	ranHash := SHA256(abi.Bytes(input[:]))
	hashInt := IntFromBigEndianBytes(ranHash)
	num := int(math.Mod(float64(hashInt), float64(limit)))
	return num
}

func BigEndianBytesFromInt(x int) abi.Bytes {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	err := binary.Write(buf, binary.BigEndian, x)
	autil.AssertNoError(err)
	return buf.Bytes()
}

func SHA256(abi.Bytes) abi.Bytes {
	autil.IMPL_FINISH()
	return []byte{}
}

func IntFromBigEndianBytes(bytes []byte) int {
	autil.IMPL_FINISH()
	return -1
}
