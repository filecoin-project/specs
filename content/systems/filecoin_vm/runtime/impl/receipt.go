package impl

import (
	vmr "github.com/filecoin-project/specs-actors/actors/runtime"
	exitcode "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
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
		ReturnValue: output.ReturnValue,
		GasUsed:     gasUsed,
	}
}
