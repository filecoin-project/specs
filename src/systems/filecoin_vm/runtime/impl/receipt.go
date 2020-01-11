package impl

import (
	vmr "github.com/filecoin-project/specs/actors/runtime"
	exitcode "github.com/filecoin-project/specs/actors/runtime/exitcode"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
)

type MessageReceipt struct {
	ExitCode    exitcode.ExitCode
	ReturnValue Bytes
	GasUsed     msg.GasAmount
}

func MessageReceipt_Make(output vmr.InvocOutput, exitCode exitcode.ExitCode, gasUsed msg.GasAmount) MessageReceipt {
	return MessageReceipt{
		ExitCode:    exitCode,
		ReturnValue: output.ReturnValue(),
		GasUsed:     gasUsed,
	}
}

func MessageReceipt_MakeSystemError(errCode exitcode.SystemErrorCode, gasUsed msg.GasAmount) MessageReceipt {
	return MessageReceipt_Make(
		nil,
		exitcode.SystemError(errCode),
		gasUsed,
	)
}
