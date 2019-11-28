package interpreter

import (
	ipld "github.com/filecoin-project/specs/libraries/ipld"
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

var TODO = util.TODO
var IMPL_FINISH = util.IMPL_FINISH

const (
	// TODO: reduce these when expected gas costs are known
	sumbitElectionPostGasLimit = 1e18
	payBlockRewardGasLimit     = 1e1
	cronTickGasLimit           = 1e18
)

// Applies all the message in a tipset, along with implicit block- and tipset-specific state
// transitions.
func (vmi *VMInterpreter_I) ApplyTipSetMessages(inTree st.StateTree, msgs TipSetMessages) (outTree st.StateTree, receipts []msg.MessageReceipt) {
	outTree = inTree
	seenMsgs := make(map[ipld.CID]struct{}) // CIDs of messages already seen once.
	var r msg.MessageReceipt
	for _, blk := range msgs.Blocks() {
		// Pay block reward.
		reward := _makeBlockRewardMessage(outTree, blk.Miner())
		outTree, r = vmi.ApplyMessage(outTree, reward, 0, blk.Miner())
		if r.ExitCode() != exitcode.OK() {
			panic("block reward failed")
		}

		// Process block miner's Election PoSt.
		epost := _makeElectionPoStMessage(outTree, blk.Miner(), msgs.Epoch(), blk.PoStProof())
		outTree, r = vmi.ApplyMessage(outTree, epost, 0, blk.Miner())
		if r.ExitCode() != exitcode.OK() {
			panic("election post failed")
		}

		// Process messages from the block.
		for _, m := range blk.Messages() {
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}

			TODO() // TODO: include signature length?
			onChainMessageLen := len(msg.Serialize_UnsignedMessage(m))

			outTree, r = vmi.ApplyMessage(outTree, m, onChainMessageLen, blk.Miner())
			receipts = append(receipts, r)
			seenMsgs[_msgCID(m)] = struct{}{}
		}
	}

	// Invoke cron tick, attributing it to the miner of the first block.

	TODO()
	// TODO: miners shouldn't be able to trigger cron by sending messages.
	// Maybe we need a separate ControlActor for this?

	cronSender := msgs.Blocks()[0].Miner()
	cron := _makeCronTickMessage(outTree, cronSender)
	outTree, r = vmi.ApplyMessage(outTree, cron, 0, cronSender)
	if r.ExitCode() != exitcode.OK() {
		panic("cron tick failed")
	}

	return
}

func StateTree_WithGasUsedFundsTransfer(
	tree st.StateTree, gasUsed msg.GasAmount, message msg.UnsignedMessage,
	fundsFrom addr.Address, fundsTo addr.Address) st.StateTree {

	TODO() // TODO: what if the miner has insufficient funds?
	// Should we instead aggregate these deltas, and SubtractWhileNonnegative at the end of the tipset?

	return _withTransferFundsAssert(
		tree,
		fundsFrom,
		fundsTo,
		_gasToFIL(gasUsed, message.GasPrice()),
	)
}

func (vmi *VMInterpreter_I) ApplyMessage(
	inTree st.StateTree, message msg.UnsignedMessage, onChainMessageSize int, minerAddr addr.Address) (
	st.StateTree, msg.MessageReceipt) {

	vmiGasRemaining := message.GasLimit()

	_applyReturn := func(
		tree st.StateTree, invocOutput msg.InvocOutput, exitCode exitcode.ExitCode,
		gasUsedFrom addr.Address, gasUsedTo addr.Address) (st.StateTree, msg.MessageReceipt) {

		vmiGasUsed := message.GasLimit().Subtract(vmiGasRemaining)
		tree = StateTree_WithGasUsedFundsTransfer(tree, vmiGasUsed, message, gasUsedFrom, gasUsedTo)
		return tree, msg.MessageReceipt_Make(invocOutput, exitCode, vmiGasUsed)
	}

	_applyError := func(
		tree st.StateTree, errExitCode exitcode.SystemErrorCode,
		gasUsedFrom addr.Address, gasUsedTo addr.Address) (st.StateTree, msg.MessageReceipt) {

		return _applyReturn(
			tree, msg.InvocOutput_Make(nil), exitcode.SystemError(errExitCode), gasUsedFrom, gasUsedTo)
	}

	_vmiUseGas := func(amount msg.GasAmount) (vmiUseGasOK bool) {
		vmiGasRemaining, vmiUseGasOK = vmiGasRemaining.SubtractWhileNonnegative(amount)
		return
	}

	ok := _vmiUseGas(gascost.OnChainMessage(onChainMessageSize))
	if !ok {
		return _applyError(inTree, exitcode.OutOfGas, minerAddr, addr.BurntFundsActorAddr)
	}

	fromActor := inTree.GetActorState(message.From())
	if fromActor == nil {
		return _applyError(inTree, exitcode.ActorNotFound, minerAddr, addr.BurntFundsActorAddr)
	}

	// make sure fromActor has enough money to run with the specified gas limit
	gasLimitCost := _gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(gasLimitCost)
	if fromActor.Balance() < totalCost {
		return _applyError(inTree, exitcode.InsufficientFunds_System, minerAddr, addr.BurntFundsActorAddr)
	}

	// make sure this is the right message order for fromActor
	// (this is protection against replay attacks, and useful sequencing)
	if message.CallSeqNum() != fromActor.CallSeqNum()+1 {
		return _applyError(inTree, exitcode.InvalidCallSeqNum, minerAddr, addr.BurntFundsActorAddr)
	}

	// WithActorForAddress may create new account actors
	compTreePreSend := inTree
	compTreePreSend, toActor := compTreePreSend.Impl().WithActorForAddress(message.To())
	if toActor == nil {
		return _applyError(inTree, exitcode.ActorNotFound, message.From(), addr.BurntFundsActorAddr)
	}

	// deduct gas limit funds from sender first
	compTreePreSend = _withTransferFundsAssert(
		compTreePreSend, message.From(), addr.BurntFundsActorAddr, gasLimitCost)

	// Increment sender call sequence number.
	var err error
	compTreePreSend, err = compTreePreSend.Impl().WithIncrementedCallSeqNum(message.From())
	if err != nil {
		// Note: if actor deletion is possible at some point, may need to allow this case
		panic("Internal interpreter error: failed to increment call sequence number")
	}

	rt := vmr.VMContext_Make(
		message.From(),
		minerAddr,
		fromActor.CallSeqNum(),
		actor.CallSeqNum(0),
		compTreePreSend,
		message.From(),
		actor.TokenAmount(0),
		vmiGasRemaining,
	)

	sendRet, compTreePostSend := rt.SendToplevelFromInterpreter(
		msg.InvocInput_Make(
			message.To(),
			message.Method(),
			message.Params(),
			message.Value(),
		),
	)

	ok = _vmiUseGas(sendRet.GasUsed())
	if !ok {
		panic("Interpreter error: runtime execution used more gas than provided")
	}

	ok = _vmiUseGas(gascost.OnChainReturnValue(sendRet.ReturnValue()))
	if !ok {
		return _applyError(compTreePreSend, exitcode.OutOfGas, addr.BurntFundsActorAddr, minerAddr)
	}

	compTreeRet := compTreePreSend
	if sendRet.ExitCode().AllowsStateUpdate() {
		compTreeRet = compTreePostSend
	}

	// Refund unused gas to sender.
	refundGas := vmiGasRemaining
	compTreeRet = _withTransferFundsAssert(
		compTreeRet,
		addr.BurntFundsActorAddr,
		message.From(),
		_gasToFIL(refundGas, message.GasPrice()),
	)

	return _applyReturn(
		compTreeRet, msg.InvocOutput_Make(sendRet.ReturnValue()), sendRet.ExitCode(),
		addr.BurntFundsActorAddr, minerAddr)
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
func _makeBlockRewardMessage(state st.StateTree, minerActorAddr addr.Address) msg.UnsignedMessage {
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
