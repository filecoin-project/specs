package message

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"

func MessageReceipt_Make(output InvocOutput, gasUsed GasAmount) MessageReceipt {
	return &MessageReceipt_I{
		ExitCode_:    output.ExitCode(),
		ReturnValue_: output.ReturnValue(),
		GasUsed_:     gasUsed,
	}
}

func (x GasAmount) Add(y GasAmount) GasAmount {
	panic("TODO")
}

func (x GasAmount) Subtract(y GasAmount) GasAmount {
	panic("TODO")
}

func (x GasAmount) LessThan(y GasAmount) bool {
	panic("TODO")
}

func GasAmount_Zero() GasAmount {
	panic("TODO")
}

func InvocInput_Make(method actor.MethodNum, params actor.MethodParams, value actor.TokenAmount) InvocInput {
	return &InvocInput_I{
		Method_: method,
		Params_: params,
		Value_:  value,
	}
}

func InvocOutput_Make(exitCode exitcode.ExitCode, returnValue util.Bytes) InvocOutput {
	return &InvocOutput_I{
		ExitCode_:    exitCode,
		ReturnValue_: returnValue,
	}
}

func MessageReceipt_MakeSystemError(errCode exitcode.SystemErrorCode, gasUsed GasAmount) MessageReceipt {
	return MessageReceipt_Make(
		InvocOutput_Make(exitcode.SystemError(errCode), nil),
		gasUsed,
	)
}
