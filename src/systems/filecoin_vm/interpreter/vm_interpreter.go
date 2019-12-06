package interpreter

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	storage_mining "github.com/filecoin-project/specs/systems/filecoin_mining/storage_mining"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

type Bytes = util.Bytes

var Assert = util.Assert
var TODO = util.TODO
var IMPL_FINISH = util.IMPL_FINISH

type SenderResolveSpec int

const (
	SenderResolveSpec_OK SenderResolveSpec = 1 + iota
	SenderResolveSpec_Invalid
)

// Applies all the message in a tipset, along with implicit block- and tipset-specific state
// transitions.
func (vmi *VMInterpreter_I) ApplyTipSetMessages(inTree st.StateTree, msgs TipSetMessages) (outTree st.StateTree, receipts []vmr.MessageReceipt) {
	outTree = inTree
	seenMsgs := make(map[ipld.CID]struct{}) // CIDs of messages already seen once.
	var receipt vmr.MessageReceipt

	for _, blk := range msgs.Blocks() {
		minerOwner := storage_mining.GetMinerOwnerAddress_Assert(inTree, blk.Miner())

		// Process block miner's Election PoSt.
		epostMessage := _makeElectionPoStMessage(outTree, blk.Miner(), msgs.Epoch(), blk.PoStProof())
		outTree = _applyMessageBuiltinAssert(outTree, epostMessage, blk.Miner())

		minerPenaltyTotal := actor.TokenAmount(0)
		var minerPenaltyCurr actor.TokenAmount

		// Process BLS messages from the block.
		for _, m := range blk.BLSMessages() {
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}
			onChainMessageLen := len(msg.Serialize_UnsignedMessage(m))
			outTree, receipt, minerPenaltyCurr = vmi.ApplyMessage(outTree, m, onChainMessageLen, blk.Miner())
			minerPenaltyTotal += minerPenaltyCurr
			receipts = append(receipts, receipt)
			seenMsgs[_msgCID(m)] = struct{}{}
		}

		// Process SECP messages from the block.
		for _, sm := range blk.SECPMessages() {
			m := sm.Message()
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}
			onChainMessageLen := len(msg.Serialize_SignedMessage(sm))
			outTree, receipt, minerPenaltyCurr = vmi.ApplyMessage(outTree, m, onChainMessageLen, blk.Miner())
			minerPenaltyTotal += minerPenaltyCurr
			receipts = append(receipts, receipt)
			seenMsgs[_msgCID(m)] = struct{}{}
		}

		// Pay block reward.
		rewardMessage := _makeBlockRewardMessage(outTree, minerOwner, minerPenaltyTotal)
		outTree = _applyMessageBuiltinAssert(outTree, rewardMessage, blk.Miner())
	}

	// Invoke cron tick, attributing it to the miner of the first block.

	TODO()
	// TODO: miners shouldn't be able to trigger cron by sending messages;
	// use ControlActor instead (https://github.com/filecoin-project/specs/issues/665)

	firstMiner := msgs.Blocks()[0].Miner()
	cronMessage := _makeCronTickMessage(outTree, firstMiner)
	outTree = _applyMessageBuiltinAssert(outTree, cronMessage, firstMiner)

	return
}

func (vmi *VMInterpreter_I) ApplyMessage(
	inTree st.StateTree, message msg.UnsignedMessage, onChainMessageSize int, minerAddr addr.Address) (
	retTree st.StateTree, retReceipt vmr.MessageReceipt, retMinerPenalty actor.TokenAmount) {

	minerOwner := storage_mining.GetMinerOwnerAddress_Assert(inTree, minerAddr)

	vmiGasRemaining := message.GasLimit()
	vmiGasUsed := msg.GasAmount_Zero()

	_applyReturn := func(
		tree st.StateTree, invocOutput vmr.InvocOutput, exitCode exitcode.ExitCode,
		senderResolveSpec SenderResolveSpec) {

		vmiGasRemainingFIL := _gasToFIL(vmiGasRemaining, message.GasPrice())
		vmiGasUsedFIL := _gasToFIL(vmiGasUsed, message.GasPrice())

		switch senderResolveSpec {
		case SenderResolveSpec_OK:
			// In this case, the sender is valid and has already transferred funds to the burnt funds actor
			// sufficient for the gas limit. Thus, we may refund the unused gas funds to the sender here.
			Assert(!message.GasLimit().LessThan(vmiGasUsed))
			Assert(message.GasLimit().Equals(vmiGasUsed.Add(vmiGasRemaining)))
			tree = _withTransferFundsAssert(tree, addr.BurntFundsActorAddr, message.From(), vmiGasRemainingFIL)
			tree = _withTransferFundsAssert(tree, addr.BurntFundsActorAddr, minerOwner, vmiGasUsedFIL)
			retMinerPenalty = actor.TokenAmount(0)

		case SenderResolveSpec_Invalid:
			retMinerPenalty = vmiGasUsedFIL

		default:
			Assert(false)
		}

		retTree = tree
		retReceipt = vmr.MessageReceipt_Make(invocOutput, exitCode, vmiGasUsed)
	}

	_applyError := func(tree st.StateTree, errExitCode exitcode.SystemErrorCode, senderResolveSpec SenderResolveSpec) {
		_applyReturn(tree, vmr.InvocOutput_Make(nil), exitcode.SystemError(errExitCode), senderResolveSpec)
	}

	// Deduct an amount of gas corresponding to cost about to be incurred, but not necessarily
	// incurred yet.
	_vmiAllocGas := func(amount msg.GasAmount) (vmiAllocGasOK bool) {
		vmiGasRemaining, vmiAllocGasOK = vmiGasRemaining.SubtractIfNonnegative(amount)
		vmiGasUsed = message.GasLimit().Subtract(vmiGasRemaining)
		Assert(!vmiGasRemaining.LessThan(msg.GasAmount_Zero()))
		Assert(!vmiGasUsed.LessThan(msg.GasAmount_Zero()))
		return
	}

	// Deduct an amount of gas corresponding to costs already incurred, and for which the
	// gas cost must be paid even if it would cause the gas used to exceed the limit.
	_vmiBurnGas := func(amount msg.GasAmount) (vmiBurnGasOK bool) {
		vmiGasUsedPre := vmiGasUsed
		vmiBurnGasOK = _vmiAllocGas(amount)
		if !vmiBurnGasOK {
			vmiGasRemaining = msg.GasAmount_Zero()
			vmiGasUsed = vmiGasUsedPre.Add(amount)
		}
		return
	}

	ok := _vmiBurnGas(gascost.OnChainMessage(onChainMessageSize))
	if !ok {
		// Invalid message; insufficient gas limit to pay for the on-chain message size.
		_applyError(inTree, exitcode.OutOfGas, SenderResolveSpec_Invalid)
		return
	}

	fromActor := inTree.GetActorState(message.From())
	if fromActor == nil {
		// Execution error; sender does not exist at time of message execution.
		_applyError(inTree, exitcode.ActorNotFound, SenderResolveSpec_Invalid)
		return
	}

	// Check sender balance.
	gasLimitCost := _gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(gasLimitCost)
	if fromActor.Balance() < totalCost {
		// Execution error; sender does not have sufficient funds to pay for the gas limit.
		_applyError(inTree, exitcode.InsufficientFunds_System, SenderResolveSpec_Invalid)
		return
	}

	// make sure this is the right message order for fromActor
	if message.CallSeqNum() != fromActor.CallSeqNum() {
		_applyError(inTree, exitcode.InvalidCallSeqNum, SenderResolveSpec_Invalid)
		return
	}

	compTreePreSend := inTree

	// Deduct gas limit funds from sender first.
	// (This should always succeed, due to the sender balance check above.)
	compTreePreSend = _withTransferFundsAssert(
		compTreePreSend, message.From(), addr.BurntFundsActorAddr, gasLimitCost)

	// Increment sender CallSeqNum.
	compTreePreSend = compTreePreSend.Impl().WithIncrementedCallSeqNum_Assert(message.From())

	// WithActorForAddress may create new account actors
	var toActor actor.ActorState
	compTreePreSend, toActor = compTreePreSend.Impl().WithActorForAddress(message.To())
	if toActor == nil {
		// Execution error; receiver actor does not exist (and could not be implicitly created)
		// at time of message execution.
		_applyError(compTreePreSend, exitcode.ActorNotFound, SenderResolveSpec_OK)
		return
	}

	sendRet, compTreePostSend := _applyMessageInternal(compTreePreSend, message, vmiGasRemaining, minerAddr)

	ok = _vmiBurnGas(sendRet.GasUsed())
	if !ok {
		panic("Interpreter error: runtime execution used more gas than provided")
	}

	ok = _vmiAllocGas(gascost.OnChainReturnValue(sendRet.ReturnValue()))
	if !ok {
		// Insufficient gas remaining to cover the on-chain return value; proceed as in the case
		// of method execution failure.
		_applyError(compTreePreSend, exitcode.OutOfGas, SenderResolveSpec_OK)
		return
	}

	compTreeRet := compTreePreSend
	if sendRet.ExitCode().AllowsStateUpdate() {
		compTreeRet = compTreePostSend
	}

	_applyReturn(
		compTreeRet, vmr.InvocOutput_Make(sendRet.ReturnValue()), sendRet.ExitCode(), SenderResolveSpec_OK)
	return
}

func _applyMessageBuiltinAssert(tree st.StateTree, message msg.UnsignedMessage, topLevelBlockWinner addr.Address) st.StateTree {
	gasRemainingInit := msg.GasAmount_SentinelUnlimited()
	Assert(gasRemainingInit.Equals(message.GasLimit()))

	TODO() // TODO: assert message.From() is ControlActor

	tree = tree.Impl().WithIncrementedCallSeqNum_Assert(message.From())

	retReceipt, retTree := _applyMessageInternal(tree, message, gasRemainingInit, topLevelBlockWinner)
	if retReceipt.ExitCode() != exitcode.OK() {
		panic("internal message application failed")
	}

	return retTree
}

func _applyMessageInternal(
	tree st.StateTree, message msg.UnsignedMessage, gasRemainingInit msg.GasAmount, topLevelBlockWinner addr.Address) (
	vmr.MessageReceipt, st.StateTree) {

	fromActor := tree.GetActorState(message.From())
	Assert(fromActor != nil)

	rt := vmr.VMContext_Make(
		message.From(),
		topLevelBlockWinner,
		fromActor.CallSeqNum(),
		actor.CallSeqNum(0),
		tree,
		message.From(),
		actor.TokenAmount(0),
		gasRemainingInit,
	)

	return rt.SendToplevelFromInterpreter(
		vmr.InvocInput_Make(
			message.To(),
			message.Method(),
			message.Params(),
			message.Value(),
		),
	)
}

func _withTransferFundsAssert(tree st.StateTree, from addr.Address, to addr.Address, amount actor.TokenAmount) st.StateTree {
	// TODO: assert amount nonnegative
	retTree, err := tree.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		panic("Interpreter error: insufficient funds (or transfer error) despite checks")
	} else {
		return retTree
	}
}

func _gasToFIL(gas msg.GasAmount, price msg.GasPrice) actor.TokenAmount {
	IMPL_FINISH()
	panic("") // BigInt arithmetic
	// return actor.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}

// Builds a message for paying block reward from the treasury account to a miner owner.
func _makeBlockRewardMessage(
	state st.StateTree, minerOwnerAddr addr.Address, minerPenaltyTotal actor.TokenAmount) msg.UnsignedMessage {

	var blockReward actor.TokenAmount
	TODO() // TODO: finish

	blockReward -= minerPenaltyTotal
	if blockReward < 0 {
		blockReward = 0
	}

	// networkTreasuryActor := loadActor(NetworkTreasuryActorAddress)
	// minerActor := loadActorSubstate(minerActorAddr)
	// minerWorker := loadActor(minerActor.Info.Worker)
	//return &UnsignedMessage_I{
	//	From_:       minerActor.Info.Worker,
	//	To_:         NetworkTreasuryActorAddress,
	//	Method_:     NetworkTreasury.PayBlockReward,
	//	Params_:     serialize([minerActorAddr]),
	//	CallSeqNum_: minerWorker.CallSeqNum,
	//	Value_:      0,
	//	GasPrice_:   0,
	//	GasLimit_:   payBlockRewardGasLimit,
	//}
	panic("TODO: implement when network treasury actor implemented")
}

// Builds a message for submitting ElectionPost on behalf of a miner actor.
func _makeElectionPoStMessage(state st.StateTree, minerActorAddr addr.Address, epoch UInt64, postProof Bytes) msg.UnsignedMessage {
	// minerActor := loadActorSubstate(minerActorAddr)
	// minerWorker := loadActor(minerActor.Info.Worker)
	//return &UnsignedMessage_I{
	//	From_:       minerActor.Info.Worker,
	//	To_:         minerActorAddr,
	//	Method_:     StorageMinerActor.SubmitElectionPoSt,
	//	Params_:     serialize([PoStSubmission{postProof, epoch}]),
	//	CallSeqNum_: minerWorker.CallSeqNum,
	//	Value_:      0,
	//	GasPrice_:   0,
	//	GasLimit_:   sumbitElectionPostGasLimit,
	//}
	panic("TODO: implement when necessary dependencies are importable")
}

// Builds a message for invoking the cron actor tick.
func _makeCronTickMessage(state st.StateTree, minerActorAddr addr.Address) msg.UnsignedMessage {
	// minerActor := loadActorSubstate(minerActorAddr)
	// minerWorker := loadActor(minerActor.Info.Worker)
	//return &UnsignedMessage_I{
	//	From_:       minerActor.Info.Worker,
	//	To_:         CronActorAddress,
	//	Method_:     CronActor.EpochTick,
	//	Params_:     nil,
	//	CallSeqNum_: minerWorker.CallSeqNum,
	//	Value_:      0,
	//	GasPrice_:   0,
	//	GasLimit_:   cronTickGasLimit,
	//}
	panic("TODO: implement when necessary dependencies are importable")
}

func _msgCID(msg msg.UnsignedMessage) ipld.CID {
	panic("TODO")
}

func _encodeParams(p []interface{}) Bytes {
	panic("TODO")
}
