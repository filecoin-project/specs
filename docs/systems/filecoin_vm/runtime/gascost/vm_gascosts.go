package runtime

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	actor "github.com/filecoin-project/specs/actors/builtin"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	util "github.com/filecoin-project/specs/util"
)

type Bytes = util.Bytes

var TODO = util.TODO

var (
	// TODO: assign all of these.
	GasAmountPlaceholder                 = msg.GasAmount_FromInt(1)
	GasAmountPlaceholder_UpdateStateTree = GasAmountPlaceholder
)

var (
	///////////////////////////////////////////////////////////////////////////
	// System operations
	///////////////////////////////////////////////////////////////////////////

	// Gas cost charged to the originator of an on-chain message (regardless of
	// whether it succeeds or fails in application) is given by:
	//   OnChainMessageBase + len(serialized message)*OnChainMessagePerByte
	// Together, these account for the cost of message propagation and validation,
	// up to but excluding any actual processing by the VM.
	// This is the cost a block producer burns when including an invalid message.
	OnChainMessageBase    = GasAmountPlaceholder
	OnChainMessagePerByte = GasAmountPlaceholder

	// Gas cost charged to the originator of a non-nil return value produced
	// by an on-chain message is given by:
	//   len(return value)*OnChainReturnValuePerByte
	OnChainReturnValuePerByte = GasAmountPlaceholder

	// Gas cost for any message send execution(including the top-level one
	// initiated by an on-chain message).
	// This accounts for the cost of loading sender and receiver actors and
	// (for top-level messages) incrementing the sender's sequence number.
	// Load and store of actor sub-state is charged separately.
	SendBase = GasAmountPlaceholder

	// Gas cost charged, in addition to SendBase, if a message send
	// is accompanied by any nonzero currency amount.
	// Accounts for writing receiver's new balance (the sender's state is
	// already accounted for).
	SendTransferFunds = GasAmountPlaceholder

	// Gas cost charged, in addition to SendBase, if a message invokes
	// a method on the receiver.
	// Accounts for the cost of loading receiver code and method dispatch.
	SendInvokeMethod = GasAmountPlaceholder

	// Gas cost (Base + len*PerByte) for any Get operation to the IPLD store
	// in the runtime VM context.
	IpldGetBase    = GasAmountPlaceholder
	IpldGetPerByte = GasAmountPlaceholder

	// Gas cost (Base + len*PerByte) for any Put operation to the IPLD store
	// in the runtime VM context.
	//
	// Note: these costs should be significantly higher than the costs for Get
	// operations, since they reflect not only serialization/deserialization
	// but also persistent storage of chain data.
	IpldPutBase    = GasAmountPlaceholder
	IpldPutPerByte = GasAmountPlaceholder

	// Gas cost for updating an actor's substate (i.e., UpdateRelease).
	// This is in addition to a per-byte fee for the state as for IPLD Get/Put.
	UpdateActorSubstate = GasAmountPlaceholder_UpdateStateTree

	// Gas cost for creating a new actor (via InitActor's Exec method).
	// Actor sub-state is charged separately.
	ExecNewActor = GasAmountPlaceholder

	// Gas cost for deleting an actor.
	DeleteActor = GasAmountPlaceholder

	///////////////////////////////////////////////////////////////////////////
	// Pure functions (VM ABI)
	///////////////////////////////////////////////////////////////////////////

	// Gas cost charged per public-key cryptography operation (e.g., signature
	// verification).
	PublicKeyCryptoOp = GasAmountPlaceholder
)

func OnChainMessage(onChainMessageLen int) msg.GasAmount {
	return msg.GasAmount_Affine(OnChainMessageBase, onChainMessageLen, OnChainMessagePerByte)
}

func OnChainReturnValue(returnValue Bytes) msg.GasAmount {
	retLen := 0
	if returnValue != nil {
		retLen = len(returnValue)
	}

	return msg.GasAmount_Affine(msg.GasAmount_Zero(), retLen, OnChainReturnValuePerByte)
}

func IpldGet(dataSize int) msg.GasAmount {
	return msg.GasAmount_Affine(IpldGetBase, dataSize, IpldGetPerByte)
}

func IpldPut(dataSize int) msg.GasAmount {
	return msg.GasAmount_Affine(IpldPutBase, dataSize, IpldPutPerByte)
}

func InvokeMethod(value abi.TokenAmount, method abi.MethodNum) msg.GasAmount {
	ret := SendBase
	if value != abi.TokenAmount(0) {
		ret = ret.Add(SendTransferFunds)
	}
	if method != actor.MethodSend {
		ret = ret.Add(SendInvokeMethod)
	}
	return ret
}
