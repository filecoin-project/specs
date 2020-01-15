package crypto

import (
	"bytes"
	"encoding/binary"
	abi "github.com/filecoin-project/specs/actors/abi"
	util "github.com/filecoin-project/specs/util"
	"math"
)

func RandomInt(randomness abi.Randomness, nonce int, limit int) int {
	nonceBytes := LittleEndianBytesFromInt(nonce)
	input := randomness
	input = append(input, nonceBytes...)
	ranHash := SHA256(input[:])
	hashInt := IntFromLittleEndianBytes(ranHash)
	num := int(math.Mod(float64(hashInt), float64(limit)))
	return num
}

func SHA256(util.Bytes) util.Bytes {
	util.IMPL_FINISH()
	return []byte{}
}

func IntFromLittleEndianBytes(bytes []byte) int {
	util.IMPL_FINISH()
	return -1
}

func LittleEndianBytesFromInt(x int) util.Bytes {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	err := binary.Write(buf, binary.LittleEndian, x)
	util.Assert(err == nil)
	return buf.Bytes()
}
