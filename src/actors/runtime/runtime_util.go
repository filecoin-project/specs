package runtime

import (
	"bytes"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
	util "github.com/filecoin-project/specs/util"
)

// TODO: most of this file doesn't need to be part of runtime, just generic actor shared code.

var Assert = util.Assert
var IMPL_TODO = util.IMPL_TODO

type Any = util.Any

// Name should be set per unique filecoin network
var Name = "mainnet"

func NetworkName() string {
	return Name
}

type MinerEntrySpec int64

const (
	MinerEntrySpec_MinerOnly = iota
	MinerEntrySpec_MinerOrSignable
)

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	//InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput
	// IMPL_TODO: method dispatch mechanism is deferred to implementations.
	// When the executable actor spec is complete we can re-instantiate something here.
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

func InvocInput_Make(to addr.Address, method abi.MethodNum, params abi.MethodParams, value abi.TokenAmount) InvocInput {
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

func RT_ValidateImmediateCallerIsSignable(rt Runtime) {
	rt.ValidateImmediateCallerAcceptAnyOfTypes([]abi.ActorCodeID{
		builtin.AccountActorCodeID,
		builtin.MultisigActorCodeID,
	})
}

func RT_Address_Is_StorageMiner(rt Runtime, minerAddr addr.Address) bool {
	codeID, ok := rt.GetActorCodeID(minerAddr)
	Assert(ok)
	return codeID == builtin.StorageMinerActorCodeID
}

func RT_GetMinerAccountsAssert(rt Runtime, minerAddr addr.Address) (ownerAddr addr.Address, workerAddr addr.Address) {
	raw := rt.SendQuery(minerAddr, builtin.Method_StorageMinerActor_GetOwnerAddr, nil)
	r := bytes.NewReader(raw)
	err := ownerAddr.UnmarshalCBOR(r)
	util.Assert(err == nil)

	raw = rt.SendQuery(minerAddr, builtin.Method_StorageMinerActor_GetWorkerAddr, nil)
	r = bytes.NewReader(raw)
	err = workerAddr.UnmarshalCBOR(r)
	util.Assert(err == nil)

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

func RT_ConfirmFundsReceiptOrAbort_RefundRemainder(rt Runtime, fundsRequired abi.TokenAmount) {
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
