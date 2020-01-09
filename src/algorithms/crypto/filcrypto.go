package crypto

import (
	"bytes"
	"encoding/binary"
	"math"

	util "github.com/filecoin-project/specs/util"
)

func RandomInt(randomness util.Bytes, nonce int, limit int) int {
	nonceBytes := LittleEndianBytesFromInt(nonce)
	input := randomness
	input = append(input, nonceBytes...)
	ranHash := SHA256(input[:])
	hashInt := IntFromLittleEndianBytes(ranHash)
	num := int(math.Mod(float64(hashInt), float64(limit)))
	return num
}

func SHA256(util.Bytes) util.Bytes {
	panic("TODO")
}

func IntFromLittleEndianBytes(bytes []byte) int {
	panic("TODO")
	return -1
}

func LittleEndianBytesFromInt(x int) util.Bytes {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	err := binary.Write(buf, binary.LittleEndian, x)
	util.Assert(err == nil)
	return buf.Bytes()
}
