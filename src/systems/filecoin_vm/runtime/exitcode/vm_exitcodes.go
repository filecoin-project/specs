package exitcode

import util "github.com/filecoin-project/specs/util"

import (
	"fmt"
)

type SystemErrorCode util.Int

const (
	// TODO: remove once canonical error codes are finalized
	SystemErrorCode_Placeholder = SystemErrorCode(-(1 << 30))
)

// TODO: assign all of these.
const (
	// ActorNotFound represents a failure to find an actor.
	ActorNotFound = SystemErrorCode_Placeholder + iota

	// ActorCodeNotFound represents a failure to find the code for a
	// particular actor in the VM registry.
	ActorCodeNotFound

	// InvalidMethod represents a failure to find a method in
	// an actor
	InvalidMethod

	// InvalidArguments indicates that a method was called with the incorrect
	// number of arguments, or that its arguments did not satisfy its
	// preconditions
	InvalidArguments

	// InsufficientFunds represents a failure to apply a message, as
	// it did not carry sufficient funds for its application.
	InsufficientFunds

	// InvalidCallSeqNum represents a message invocation out of sequence.
	// This happens when message.CallSeqNum is not exactly actor.CallSeqNum + 1
	InvalidCallSeqNum

	// OutOfGasError is returned when the execution of an actor method
	// (including its subcalls) uses more gas than initially allocated.
	OutOfGas

	// RuntimeAPIError is returned when an actor method invocation makes a call
	// to the runtime that does not satisfy its preconditions.
	RuntimeAPIError

	// MethodPanic is returned when an actor method invocation calls rt.Abort.
	MethodAbort

	// MethodSubcallError is returned when an actor method's Send call has
	// returned with a failure error code (and the Send call did not specify
	// to ignore errors).
	MethodSubcallError
)

var (
	InvalidSectorPacking = UserDefinedError(1)
)

func OK() ExitCode {
	return ExitCode_Make_Success(&ExitCode_Success_I{})
}

func SystemError(x SystemErrorCode) ExitCode {
	return ExitCode_Make_SystemError(ExitCode_SystemError(x))
}

func (x *ExitCode_I) IsSuccess() bool {
	return x.Which() == ExitCode_Case_Success
}

func (x *ExitCode_I) IsError() bool {
	return !x.IsSuccess()
}

func (x *ExitCode_I) AllowsStateUpdate() bool {
	// TODO: Confirm whether this is the desired behavior

	// return x.IsSuccess() || x.Which() == ExitCode_Case_UserDefinedError
	return x.IsSuccess()
}

func EnsureErrorCode(x ExitCode) ExitCode {
	if !x.IsError() {
		// Throwing an error with a non-error exit code is itself an error
		x = SystemError(RuntimeAPIError)
	}
	return x
}

type RuntimeError struct {
	ExitCode ExitCode
	ErrMsg   string
}

func (x *RuntimeError) String() string {
	ret := fmt.Sprintf("Runtime error: %v", x.ExitCode)
	if x.ErrMsg != "" {
		ret += fmt.Sprintf(" (\"%v\")", x.ErrMsg)
	}
	return ret
}

func RuntimeError_Make(exitCode ExitCode, errMsg string) *RuntimeError {
	exitCode = EnsureErrorCode(exitCode)
	return &RuntimeError{
		ExitCode: exitCode,
		ErrMsg:   errMsg,
	}
}

func UserDefinedError(e util.UVarint) ExitCode {
	return ExitCode_Make_UserDefinedError(ExitCode_UserDefinedError(e))
}
