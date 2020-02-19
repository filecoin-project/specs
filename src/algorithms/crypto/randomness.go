package crypto

import (
	"bytes"
	"encoding/binary"

	addr "github.com/filecoin-project/go-address"
	"github.com/minio/blake2b-simd"

	abi "github.com/filecoin-project/specs-actors/actors/abi"
	big "github.com/filecoin-project/specs-actors/actors/abi/big"
	autil "github.com/filecoin-project/specs-actors/actors/util"
)

// Specifies a domain for randomness generation.
type DomainSeparationTag int

const (
	DomainSeparationTag_TicketProduction DomainSeparationTag = 1 + iota
	DomainSeparationTag_ElectionPoStChallengeSeed
	DomainSeparationTag_SurprisePoStChallengeSeed
	DomainSeparationTag_SurprisePoStSelectMiners
	DomainSeparationTag_SurprisePoStSampleSectors
)

// Derive a random byte string from a domain separation tag and the appropriate values
func DeriveRandWithMinerAddr(tag DomainSeparationTag, randSeed []byte, minerAddr addr.Address) abi.Randomness {
	var addrBuf bytes.Buffer
	err := minerAddr.MarshalCBOR(&addrBuf)
	autil.AssertNoError(err)

	return _deriveRandInternal(tag, randSeed, addrBuf.Bytes())
}

type AddressEpochEntropy struct {
	MinerAddress addr.Address // Must be an ID-addr
	Epoch        abi.ChainEpoch
}

func DeriveRandWithMinerAddrAndEpoch(tag DomainSeparationTag, randSeed abi.Randomness, minerAddr addr.Address, epoch abi.ChainEpoch) abi.Randomness {
	var st init_.State
	rt.GetState(&st)
	minerIDAddr, err := st.ResolveAddress(adt.AsStore(rt), uniqueAddr1)
	assert.NoError(t, err)
	assert.Equal(t, expectedIdAddr1, actualIdAddr)
	minerAddr := 


	entropy := &AddressEpochEntropy{
		MinerAddress: minerAddr,
		Epoch:        epoch,
	}
	var entrBuf bytes.Buffer
	err := entropy.MarshalCBOR(&entrBuf)
	autil.AssertNoError(err)

	return _deriveRandInternal(tag, randSeed, entrBuf.Bytes())
}

func _deriveRandInternal(tag DomainSeparationTag, randSeed abi.Randomness, serializedEntropy []byte) abi.Randomness {
	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int64(tag))...)
	buffer = append(buffer, randSeed...)
	buffer = append(buffer, serializedEntropy...)
	bufHash := blake2b.Sum256(buffer)
	return bufHash[:]
}

// TODO hs: remove once 148 lands
// Computes an unpredictable integer less than limit from inputs seed and nonce.
func RandomInt(seed abi.Randomness, nonce int64, limit int64) int64 {
	nonceBytes := BigEndianBytesFromInt(nonce)
	input := append(seed, nonceBytes...)
	ranHash := blake2b.Sum256(input)
	hashInt := big.PositiveFromUnsignedBytes(ranHash[:])

	num := big.Mod(hashInt, big.NewInt(limit))
	return num.Int64()
}

// Returns an 8-byte slice of the big-endian bytes of an integer.
func BigEndianBytesFromInt(x int64) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	err := binary.Write(buf, binary.BigEndian, x)
	autil.AssertNoError(err)
	return buf.Bytes()
}
