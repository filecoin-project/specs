package util

import "fmt"

// Indicates a condition that should never happen. If encountered, execution will halt and the
// resulting state is undefined.
func AssertMsg(b bool, format string, a ...interface{}) {
	if !b {
		panic(fmt.Sprintf(format, a...))
	}
}

func Assert(b bool) {
	AssertMsg(b, "assertion failed")
}

func AssertNoError(e error) {
	AssertMsg(e == nil, e.Error())
}

// Indicating behavior not yet specified, and may require other spec changes.
func TODO(...interface{}) {
	// Indirection to prevent the compiler from ignoring unreachable code
	panic("TODO")
}

// Version of TODO() indicating that the operation is clearly implementable,
// but some details remain to be filled in during implementation.
func IMPL_TODO(...interface{}) {
	panic("Not yet implemented in the spec")
}

// Version of TODO() indicating that the operation is believed to be unambiguous,
// but is not yet implemented as code in the spec repository.
func IMPL_FINISH(...interface{}) {
	panic("Not yet implemented in the spec")
}

// Version of TODO() indicating that the operation is believed to be unambiguous,
// but is not yet implemented as code in the spec repository.
func PARAM_FINISH(...interface{}) {
	panic("Not yet implemented in the spec")
}
