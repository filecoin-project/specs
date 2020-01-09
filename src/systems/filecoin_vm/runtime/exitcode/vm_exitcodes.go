package exitcode

import (
	"fmt"

	util "github.com/filecoin-project/specs/util"
)

type SystemErrorCode int
type UserDefinedErrorCode int

const (
	// TODO: remove once canonical error codes are finalized
	SystemErrorCode_Placeholder      = SystemErrorCode(-(1 << 30))
	UserDefinedErrorCode_Placeholder = UserDefinedErrorCode(-(1 << 30))
)

var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

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

	// InvalidArgumentsSystem indicates that a method was called with the incorrect
	// number of arguments, or that its arguments did not satisfy its
	// preconditions
	InvalidArguments_System

	// InsufficientFunds represents a failure to apply a message, as
	// it did not carry sufficient funds for its application.
	InsufficientFunds_System

	// InvalidCallSeqNum represents a message invocation out of sequence.
	// This happens when message.CallSeqNum is not exactly actor.CallSeqNum + 1
	InvalidCallSeqNum

	// OutOfGas is returned when the execution of an actor method
	// (including its subcalls) uses more gas than initially allocated.
	OutOfGas

	// RuntimeAPIError is returned when an actor method invocation makes a call
	// to the runtime that does not satisfy its preconditions.
	RuntimeAPIError

	// RuntimeAssertFailure is returned when an actor method invocation calls
	// rt.Assert with a false condition.
	RuntimeAssertFailure

	// MethodSubcallError is returned when an actor method's Send call has
	// returned with a failure error code (and the Send call did not specify
	// to ignore errors).
	MethodSubcallError
)

const (
	InsufficientFunds_User = UserDefinedErrorCode_Placeholder + iota
	InvalidArguments_User
	InconsistentState_User

	InvalidSectorPacking
	SealVerificationFailed
	PoStVerificationFailed
	DeadlineExceeded
	InsufficientPledgeCollateral
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
	return x.IsSuccess()
}

func (x *ExitCode_I) Equals(ExitCode) bool {
	IMPL_FINISH()
	panic("")
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

func UserDefinedError(e UserDefinedErrorCode) ExitCode {
	return ExitCode_Make_UserDefinedError(ExitCode_UserDefinedError(e))
}
