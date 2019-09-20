package codeGen

import (
	"unsafe"
	"log"
)

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

type Word = int

// Check that sizeof(Word) == 8 (we only support 64-bit builds for now)
type _UnusedCompileAssert1 = [unsafe.Sizeof(Word(0))-8]byte
type _UnusedCompileAssert2 = [8-unsafe.Sizeof(Word(0))]byte

type Timestamp Word

func CurrentTime() Timestamp {
	panic("TODO")
}
