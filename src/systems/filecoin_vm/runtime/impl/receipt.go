package impl

import (
	exitcode "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
)

type MessageReceipt struct {
	ExitCode    exitcode.ExitCode
	ReturnValue []byte
	GasUsed     msg.GasAmount
}

func MessageReceipt_Make(output []byte, exitCode exitcode.ExitCode, gasUsed msg.GasAmount) MessageReceipt {
	return MessageReceipt{
		ExitCode:    exitCode,
		ReturnValue: output,
		GasUsed:     gasUsed,
	}
}
