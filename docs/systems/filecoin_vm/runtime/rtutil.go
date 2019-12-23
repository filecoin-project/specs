package runtime

import (
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	ai "github.com/filecoin-project/specs/systems/filecoin_vm/actor_interfaces"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	util "github.com/filecoin-project/specs/util"
)

var Assert = util.Assert
var IMPL_TODO = util.IMPL_TODO

type Any = util.Any

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

func RT_AddressIsSignable_AbortIfNotFound(rt Runtime, addr addr.Address) bool {
	codeID, ok := rt.GetActorCodeID(addr)
	if !ok {
		rt.AbortArgMsg("Address not found")
	}
	return codeID.IsSignable()
}

func RT_ValidateImmediateCallerIsSignable(rt Runtime) {
	IMPL_TODO() // TODO: Update for MultiSig actors
	rt.ValidateImmediateCallerAcceptAnyOfTypes(actor.BuiltinActorID_SignableTypes())
}

func RT_Address_Is_StorageMiner(rt Runtime, minerAddr addr.Address) bool {
	codeID, ok := rt.GetActorCodeID(minerAddr)
	Assert(ok)
	if !codeID.IsBuiltin() {
		return false
	}
	return (codeID.As_Builtin() == actor.BuiltinActorID_StorageMiner)
}

func RT_GetMinerAccountsAssert(rt Runtime, minerAddr addr.Address) (ownerAddr addr.Address, workerAddr addr.Address) {
	ownerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetOwnerAddr, []util.Serialization{}))

	workerAddr = addr.Deserialize_Address_Compact_Assert(
		rt.SendQuery(minerAddr, ai.Method_StorageMinerActor_GetWorkerAddr, []util.Serialization{}))

	return
}

func RT_MinerEntry_ValidateCaller_DetermineFundsLocation(rt Runtime, entryAddr addr.Address, entrySpec MinerEntrySpec) addr.Address {
	if RT_Address_Is_StorageMiner(rt, entryAddr) {
		// Storage miner actor entry; implied funds recipient is the associated owner address.
		ownerAddr, workerAddr := RT_GetMinerAccountsAssert(rt, entryAddr)
		rt.ValidateImmediateCallerInSet([]addr.Address{ownerAddr, workerAddr})
		return ownerAddr
	} else {
		if entrySpec == MinerEntrySpec_MinerOnly {
			rt.AbortArgMsg("Only miner entries valid in current context")
		}
		// Ordinary account-style actor entry; funds recipient is just the entry address itself.
		RT_ValidateImmediateCallerIsSignable(rt)
		return entryAddr
	}
}

func RT_ConfirmFundsReceiptOrAbort_RefundRemainder(rt Runtime, fundsRequired actor.TokenAmount) {
	if rt.ValueReceived() < fundsRequired {
		rt.AbortFundsMsg("Insufficient funds received accompanying message")
	}

	if rt.ValueReceived() > fundsRequired {
		rt.SendFunds(rt.ImmediateCaller(), rt.ValueReceived()-fundsRequired)
	}
}

func RT_VerifySignature(rt Runtime, pk filcrypto.PublicKey, sig filcrypto.Signature, m filcrypto.Message) bool {
	ret := rt.Compute(ComputeFunctionID_VerifySignature, []Any{pk, sig, m})
	return ret.(bool)
}
