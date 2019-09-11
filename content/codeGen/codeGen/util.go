package codeGen

import (
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

func Hash(role HashRole, x []byte) []byte {
	panic("TODO")
}

type BigInt interface{}  // TODO: import

type Fraction struct {
	n BigInt
	d BigInt
}

type Word int64

type Timestamp Word

func CurrentTime() Timestamp {
	panic("TODO")
}
