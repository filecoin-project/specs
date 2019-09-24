package util

import (
	"unsafe"
	"log"
)

// TODO: finish
type CID struct {
}

// TODO: finish
type BigInt struct {
}

type Word = int

// Check that sizeof(Word) == 8 (we only support 64-bit builds for now)
type _UnusedCompileAssert1 = [unsafe.Sizeof(Word(0))-8]byte
type _UnusedCompileAssert2 = [8-unsafe.Sizeof(Word(0))]byte

type UVarint = uint64
type Varint = int64
type UInt = uint64
type Float = float64
type Bytes = []byte
type String = string

func Assert(b bool) {
	if !b {
		panic("Assertion failed")
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic("Error check failed")
	}
}

func CompareBytesStrict(x []byte, y []byte) int {
	panic("TODO")
}

func HashBlake2bInternal(x []byte) []byte {
	panic("TODO")
}

type Timestamp UVarint
func CurrentTime() Timestamp {
	panic("TODO")
}
