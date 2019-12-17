package runtime

import (
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	util "github.com/filecoin-project/specs/util"
)

// Name should be set per unique filecoin network
var Name = "mainnet"

func NetworkName() string {
	return Name
}

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput
}

type CallerPattern struct {
	Matches func(addr.Address) bool
}

func CallerPattern_MakeSingleton(x addr.Address) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool { return x == y },
	}
}

func CallerPattern_MakeSet(x []addr.Address) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool {
			for _, xi := range x {
				if y == xi {
					return true
				}
			}
			return false
		},
	}
}

func CallerPattern_MakeAcceptAny() CallerPattern {
	return CallerPattern{
		Matches: func(addr.Address) bool { return true },
	}
}

func InvocInput_Make(to addr.Address, method actor.MethodNum, params actor.MethodParams, value actor.TokenAmount) InvocInput {
	return &InvocInput_I{
		To_:     to,
		Method_: method,
		Params_: params,
		Value_:  value,
	}
}

func InvocOutput_Make(returnValue util.Bytes) InvocOutput {
	return &InvocOutput_I{
		ReturnValue_: returnValue,
	}
}

func MessageReceipt_Make(output InvocOutput, exitCode exitcode.ExitCode, gasUsed msg.GasAmount) MessageReceipt {
	return &MessageReceipt_I{
		ExitCode_:    exitCode,
		ReturnValue_: output.ReturnValue(),
		GasUsed_:     gasUsed,
	}
}

func MessageReceipt_MakeSystemError(errCode exitcode.SystemErrorCode, gasUsed msg.GasAmount) MessageReceipt {
	return MessageReceipt_Make(
		nil,
		exitcode.SystemError(errCode),
		gasUsed,
	)
}
