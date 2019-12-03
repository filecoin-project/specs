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
func (vmi *VMInterpreter_I) ApplyTipSetMessages(inTree st.StateTree, msgs TipSetMessages) (outTree st.StateTree, receipts []vmr.MessageReceipt) {
	outTree = inTree
	seenMsgs := make(map[ipld.CID]struct{}) // CIDs of messages already seen once.
	var r vmr.MessageReceipt
	for _, blk := range msgs.Blocks() {
		// Pay block reward.
		reward := _makeBlockRewardMessage(outTree, blk.Miner())
		outTree, r = vmi.ApplyMessage(outTree, reward, blk.Miner())
		if r.ExitCode() != exitcode.OK() {
			panic("block reward failed")
		}

		// Process block miner's Election PoSt.
		epost := _makeElectionPoStMessage(outTree, blk.Miner(), msgs.Epoch(), blk.PoStProof())
		outTree, r = vmi.ApplyMessage(outTree, epost, blk.Miner())
		if r.ExitCode() != exitcode.OK() {
			panic("election post failed")
		}

		// Process BLS messages from the block.
		for _, m := range blk.BLSMessages() {
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}
			outTree, r = vmi.ApplyMessage(outTree, m, blk.Miner())
			receipts = append(receipts, r)
			seenMsgs[_msgCID(m)] = struct{}{}
		}
		// Process SECP messages from the block.
		for _, sm := range blk.SECPMessages() {
			m := sm.Message()
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}
			outTree, r = vmi.ApplyMessage(outTree, m, blk.Miner())
			receipts = append(receipts, r)
			seenMsgs[_msgCID(m)] = struct{}{}
		}
	}

	// Invoke cron tick, attributing it to the miner of the first block.
	cronSender := msgs.Blocks()[0].Miner()
	cron := _makeCronTickMessage(outTree, cronSender)
	outTree, r = vmi.ApplyMessage(outTree, cron, cronSender)
	if r.ExitCode() != exitcode.OK() {
		panic("cron tick failed")
	}

	return
}

func (vmi *VMInterpreter_I) ApplyMessage(inTree st.StateTree, message msg.UnsignedMessage, minerAddr addr.Address) (
	st.StateTree, vmr.MessageReceipt) {

	compTree := inTree
	var outTree st.StateTree
	var toActor actor.ActorState
	var err error

	fromActor := compTree.GetActorState(message.From())
	if fromActor == nil {
		// TODO: This was originally exitcode.InvalidMethod; which is correct?
		return inTree, _applyError(exitcode.ActorNotFound)
	}

	// make sure fromActor has enough money to run the max invocation
	maxGasCost := _gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(maxGasCost)
	if fromActor.Balance() < totalCost {
		return inTree, _applyError(exitcode.InsufficientFunds_System)
	}

	// make sure this is the right message order for fromActor
	// (this is protection against replay attacks, and useful sequencing)
	if message.CallSeqNum() != fromActor.CallSeqNum()+1 {
		return inTree, _applyError(exitcode.InvalidCallSeqNum)
	}

	// WithActorForAddress may create new account actors
	compTree, toActor = compTree.Impl().WithActorForAddress(message.To())
	if toActor == nil {
		return inTree, _applyError(exitcode.ActorNotFound)
	}

	// deduct maximum expenditure gas funds first
	compTree = _withTransferFundsAssert(compTree, message.From(), addr.BurntFundsActorAddr, maxGasCost)

	rt := vmr.VMContext_Make(
		message.From(),
		minerAddr, // TODO: may not exist? (see below)
		fromActor.CallSeqNum(),
		actor.CallSeqNum(0),
		compTree,
		message.From(),
		actor.TokenAmount(0),
		message.GasLimit(),
	)

	sendRet, sendRetStateTree := rt.SendToplevelFromInterpreter(
		vmr.InvocInput_Make(
			message.To(),
			message.Method(),
			message.Params(),
			message.Value(),
		),
	)

	if !sendRet.ExitCode().AllowsStateUpdate() {
		// error -- revert all state changes -- ie drop updates. burn used gas.
		outTree = inTree
		outTree = _withTransferFundsAssert(
			outTree,
			message.From(),
			addr.BurntFundsActorAddr,
			_gasToFIL(sendRet.GasUsed(), message.GasPrice()),
		)
	} else {
		// success -- refund unused gas
		outTree = sendRetStateTree
		refundGas := message.GasLimit() - sendRet.GasUsed()
		TODO() // TODO: assert refundGas is nonnegative
		outTree = _withTransferFundsAssert(
			outTree,
			addr.BurntFundsActorAddr,
			message.From(),
			_gasToFIL(refundGas, message.GasPrice()),
		)
	}

	outTree, err = outTree.Impl().WithIncrementedCallSeqNum(message.To())
	if err != nil {
		// TODO: if actor deletion is possible at some point, may need to allow this case
		panic("Internal interpreter error: failed to increment call sequence number")
	}

	// reward miner gas fees
	outTree = _withTransferFundsAssert(
		outTree,
		addr.BurntFundsActorAddr,
		minerAddr, // TODO: may not exist
		_gasToFIL(sendRet.GasUsed(), message.GasPrice()),
	)

	return outTree, sendRet
}

func _applyError(errCode exitcode.SystemErrorCode) vmr.MessageReceipt {
	// TODO: should this gasUsed value be zero?
	// If nonzero, there is not guaranteed to be a nonzero gas balance from which to deduct it.
	gasUsed := gascost.ApplyMessageFail
	TODO()
	return vmr.MessageReceipt_MakeSystemError(errCode, gasUsed)
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
