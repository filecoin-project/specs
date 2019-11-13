package exitcode

import util "github.com/filecoin-project/specs/util"

import (
	"fmt"
)

type SystemErrorCode util.UVarint

// TODO: assign all of these.
var (
	// // OK is the success return value, similar to unix exit code 0.
	// OK = SystemErrorCode(0)

	// ActorNotFound represents a failure to find an actor.
	ActorNotFound = SystemErrorCode(1)

	// ActorCodeNotFound represents a failure to find the code for a
	// particular actor in the VM registry.
	ActorCodeNotFound = SystemErrorCode(2)

	// InvalidMethod represents a failure to find a method in
	// an actor
	InvalidMethod = SystemErrorCode(3)

	// InsufficientFunds represents a failure to apply a message, as
	// it did not carry sufficient funds for its application.
	InsufficientFunds = SystemErrorCode(4)

	// InvalidCallSeqNum represents a message invocation out of sequence.
	// This happens when message.CallSeqNum is not exactly actor.CallSeqNum + 1
	InvalidCallSeqNum = SystemErrorCode(5)

	// OutOfGasError is returned when the execution of an actor method
	// (including its subcalls) uses more gas than initially allocated.
	OutOfGas = SystemErrorCode(6)

	// RuntimeAPIError is returned when an actor method invocation makes a call
	// to the runtime that does not satisfy its preconditions.
	RuntimeAPIError = SystemErrorCode(7)

	// MethodPanic is returned when an actor method invocation calls rt.Abort.
	MethodAbort = SystemErrorCode(8)

	// MethodSubcallError is returned when an actor method's Send call has
	// returned with a failure error code (and the Send call did not specify
	// to ignore errors).
	MethodSubcallError = SystemErrorCode(9)
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
