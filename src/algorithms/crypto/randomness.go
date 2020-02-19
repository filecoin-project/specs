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

// Randomness returns a (pseudo)random byte array drawing from a
// random beacon at a given epoch and incorporating requisite entropy
GetRandomness(chain personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte) abi.Randomness {

	randSeed := chain.RandomnessAtEpoch(randEpoch)

	buffer := []byte{}
	buffer = append(buffer, BigEndianBytesFromInt(int64(personalization))...)
	buffer = append(buffer, randSeed...)
	buffer = append(buffer, entropy...)
	bufHash := blake2b.Sum256(buffer)
	return bufHash[:]	
}
