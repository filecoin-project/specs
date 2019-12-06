package runtime

import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import util "github.com/filecoin-project/specs/util"

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
	OnChainMessageBase    = GasAmountPlaceholder
	OnChainMessagePerByte = GasAmountPlaceholder

	// Gas cost charged to the originator of a non-nil return value produced
	// by an on-chain message is given by:
	//   len(return value)*OnChainReturnValuePerByte
	OnChainReturnValuePerByte = GasAmountPlaceholder

	// Gas cost for any method invocation (including the original one initiated
	// by an on-chain message).
	InvokeMethodBase = GasAmountPlaceholder

	// Gas cost charged, in addition to InvokeMethodBase, if a method invocation
	// is accompanied by any nonzero currency amount.
	InvokeMethodTransferFunds = GasAmountPlaceholder_UpdateStateTree

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
	UpdateActorSubstate = GasAmountPlaceholder_UpdateStateTree

	// Gas cost for creating a new actor (via InitActor's Exec method).
	ExecNewActor = GasAmountPlaceholder

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

func InvokeMethod(valueSent actor.TokenAmount) msg.GasAmount {
	ret := InvokeMethodBase

	TODO() // TODO: BigInt
	if valueSent > actor.TokenAmount(0) {
		ret = ret.Add(InvokeMethodTransferFunds)
	}
	return ret
}
