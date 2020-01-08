package util

import (
	"bytes"
	"io"
	"log"
	big "math/big"
	"unsafe"
)

// Check that sizeof(int) == 8 (we only support 64-bit builds for now)
type _UnusedCompileAssert1 = [unsafe.Sizeof(int(0)) - 8]byte
type _UnusedCompileAssert2 = [8 - unsafe.Sizeof(int(0))]byte

type Bool bool
type Int int
type Any = interface{}
type String string
type Bytes = []byte
type Serialization Bytes

func (x Bool) Native() bool {
	return bool(x)
}

func (x Int) Native() int {
	return int(x)
}

func (x String) Native() string {
	return string(x)
}

func Bool_FromNative(x bool) Bool {
	return Bool(x)
}

func Int_FromNative(x int) Int {
	return Int(x)
}

func String_FromNative(x string) String {
	return String(x)
}

// TODO
func BigFromBytes(x []byte) *BigInt {
	return big.NewInt(0)
}

// TODO
func BigFromUint64(x uint64) *BigInt {
	return big.NewInt(0)
}

func BigFromInt(int) BigInt {
	IMPL_FINISH()
	panic("")
}

func BigFromUInt(uint) BigInt {
	IMPL_FINISH()
	panic("")
}

func Serialize_Int(int) Serialization {
	IMPL_FINISH()
	panic("")
}

func Serialize_BigInt(BigInt) Serialization {
	IMPL_FINISH()
	panic("")
}

func Deserialize_BigInt(Serialization) (ret BigInt, ok bool) {
	IMPL_FINISH()
	panic("")
}

func BigInt_Add(BigInt, BigInt) BigInt {
	IMPL_FINISH()
	panic("")
}

// Indicating behavior not yet specified, and may require other spec changes.
func IsBLS(x Bytes) bool {
	IMPL_FINISH()
	return false
}

func IsSECP(x Bytes) bool {
	IMPL_FINISH()
	return false
}

// Indicating behavior not yet specified, and may require other spec changes.
func TODO(...interface{}) {
	// Indirection to prevent the compiler from ignoring unreachable code
	panic("TODO")
}

// Version of TODO() indicating that the operation is clearly implementable,
// but some details remain to specify.
func IMPL_TODO(...interface{}) {
	panic("Not yet implemented in the spec")
}

// Version of TODO() indicating that the operation is believed to be unambiguous,
// but is not yet implemented as code in the spec repository.
func IMPL_FINISH(...interface{}) {
	panic("Not yet implemented in the spec")
}

// Version of TODO() indicating that the operation is believed to be unambiguous,
// except for some potential parameter settings (numerical or simple function
// definitions) which have been deferred for the moment.
//
// Note: accepts variadic arguments, which can be used to indicate that certain
// values are required for potential future completions of the parameterization,
// and the caller must ensure they are available.
func PARAM_FINISH(...interface{}) {
	panic("Not yet implemented in the spec")
}

type UVarint = uint64
type Varint = int64
type UInt = uint64
type Float = float64
type BigInt = big.Int
type BytesKey = string // to use Bytes in map keys.
type BytesAmount = UVarint
type T = struct{}     // For use in generic definitions.
type Randomness Bytes // Randomness is a string of random bytes

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

func TextAbbrev(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func SerializeBytes(b Bytes) Serialization {
	panic("TODO")
	return b
}

func SerializeBool(b bool) Serialization {
	panic("TODO")
	var byt []byte
	return byt
}

func DeserializeBool(s Serialization) bool {
	panic("TODO")
	var b bool
	return b
}

func IntMin(x, y int) int {
	if y < x {
		return y
	}
	return x
}

func IntMax(x, y int) int {
	if y > x {
		return y
	}
	return x
}

type IntOption struct {
	isSome bool
	value  int
}

func IntOptionSome(x int) IntOption {
	return IntOption{
		isSome: true,
		value:  x,
	}
}

func IntOptionNone() IntOption {
	return IntOption{
		isSome: false,
		value:  -1,
	}
}

func IntOptionMin(a, b IntOption) IntOption {
	if a.IsNone() || b.IsNone() {
		return IntOptionNone()
	} else {
		return IntOptionSome(IntMin(a.Get(), b.Get()))
	}
}

func IntOptionMax(a, b IntOption) IntOption {
	if a.IsNone() || b.IsNone() {
		return IntOptionNone()
	} else {
		return IntOptionSome(IntMax(a.Get(), b.Get()))
	}
}

func IntOptionAdd(a, b IntOption) IntOption {
	if a.IsNone() || b.IsNone() {
		return IntOptionNone()
	} else {
		return IntOptionSome(a.Get() + b.Get())
	}
}

func (x IntOption) IsSome() bool {
	return x.isSome
}

func (x IntOption) IsNone() bool {
	return !x.isSome
}

func (x IntOption) Get() int {
	Assert(x.IsSome())
	return x.value
}

func WriteRepeat(dst io.Writer, x string, n int) {
	for i := 0; i < n; i++ {
		dst.Write([]byte(x))
	}
}

func WriteRepeatString(x string, n int) string {
	buf := bytes.NewBuffer([]byte{})
	WriteRepeat(buf, x, n)
	return buf.String()
}

func SliceContainsString(s []string, x string) bool {
	for _, si := range s {
		if x == si {
			return true
		}
	}
	return false
}

func DerefCheckString(s *string) string {
	Assert(s != nil)
	return *s
}

func RefString(s string) *string {
	return &s
}
