package interpreter

import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
import util "github.com/filecoin-project/specs/util"

var TODO = util.TODO
var IMPL_FINISH = util.IMPL_FINISH

func (vmi *VMInterpreter_I) ApplyMessageBatch(inTree st.StateTree, msgs []MessageRef) (outTree st.StateTree, ret []msg.MessageReceipt) {
	compTree := inTree
	for _, m := range msgs {
		oT, r := vmi.ApplyMessage(compTree, m.Message(), m.MessageSizeOrig(), m.Miner())
		compTree = oT        // assign the current tree. (this call always succeeds)
		ret = append(ret, r) // add message receipt
	}
	return compTree, ret
}

func _withTransferFundsAssert(tree st.StateTree, from addr.Address, to addr.Address, amount actor.TokenAmount) st.StateTree {
	TODO() // TODO: assert amount is nonnegative (should TokenAmount be an int or a BigInt?)

	retTree, err := tree.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		panic("Interpreter error: insufficient funds (or transfer error) despite checks")
	} else {
		return retTree
	}
}

func (vmi *VMInterpreter_I) ApplyMessage(
	inTree st.StateTree, message msg.UnsignedMessage, messageSizeOrig int, minerAddr addr.Address) (
	st.StateTree, msg.MessageReceipt) {

	gasUsed := msg.GasAmount_Zero()
	gasRemaining := message.GasLimit()

	_applyError := func(errCode exitcode.SystemErrorCode) (st.StateTree, msg.MessageReceipt) {
		panic("TODO: deduct gasUsed from miner block reward.")
		panic("TODO: if used gas goes to miner anyway, is this a no-op?")

		return inTree, msg.MessageReceipt_MakeSystemError(errCode, gasUsed)
	}

	_useInitGas := func(initGas msg.GasAmount) (ok bool) {
		gasUsed = gasUsed.Add(initGas)
		gasRemaining, ok = gasRemaining.SubtractWhileNonnegative(initGas)
		return
	}

	ok := _useInitGas(gascost.OnChainMessage(messageSizeOrig))
	if !ok {
		return _applyError(exitcode.OutOfGas)
	}

	compTree := inTree
	var outTree st.StateTree
	var toActor actor.ActorState
	var err error

	fromActor := compTree.GetActorState(message.From())
	if fromActor == nil {
		return _applyError(exitcode.ActorNotFound)
	}

	// make sure fromActor has enough money to run with the specified gas limit
	gasLimitCost := gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(gasLimitCost)
	if fromActor.Balance() < totalCost {
		return _applyError(exitcode.InsufficientFunds_System)
	}

	// make sure this is the right message order for fromActor
	// (this is protection against replay attacks, and useful sequencing)
	if message.CallSeqNum() != fromActor.CallSeqNum()+1 {
		return _applyError(exitcode.InvalidCallSeqNum)
	}

	// WithActorForAddress may create new account actors
	compTree, toActor = compTree.Impl().WithActorForAddress(message.To())
	if toActor == nil {
		return _applyError(exitcode.ActorNotFound)
	}

	// deduct gas limit funds first
	compTree = _withTransferFundsAssert(compTree, message.From(), addr.BurntFundsActorAddr, gasLimitCost)

	rt := vmr.VMContext_Make(
		message.From(),
		minerAddr,
		fromActor.CallSeqNum(),
		actor.CallSeqNum(0),
		compTree,
		message.From(),
		actor.TokenAmount(0),
		message.GasLimit(),
	)

	sendRet, sendRetStateTree := rt.SendToplevelFromInterpreter(
		msg.InvocInput_Make(
			message.To(),
			message.Method(),
			message.Params(),
			message.Value(),
		),
	)

	if !sendRet.ExitCode().AllowsStateUpdate() {
		// error -- revert all state changes -- ie drop updates.
		outTree = inTree

		TODO() // TODO: As currently specced, this will burn the entire gas limit,
		// not just used gas; is this desired? If not, should refund unused gas as below.

	} else {
		// success -- refund unused gas
		outTree = sendRetStateTree
		refundGas := message.GasLimit().Subtract(sendRet.GasUsed())
		if refundGas.LessThan(msg.GasAmount_Zero()) {
			panic("Interpreter error: consumed more gas than limit")
		}
		outTree = _withTransferFundsAssert(
			outTree,
			addr.BurntFundsActorAddr,
			message.From(),
			gasToFIL(refundGas, message.GasPrice()),
		)
	}

	outTree, err = outTree.Impl().WithIncrementedCallSeqNum(message.To())
	if err != nil {
		// Note: if actor deletion is possible at some point, may need to allow this case
		panic("Internal interpreter error: failed to increment call sequence number")
	}

	// reward miner gas fees

	TODO() // TODO: Is the burnt funds actor allowed to send money, or only to receive?
	// If the latter, need to find a different temporary holding place for this amount.

	outTree = _withTransferFundsAssert(
		outTree,
		addr.BurntFundsActorAddr,
		minerAddr,
		gasToFIL(sendRet.GasUsed(), message.GasPrice()),
	)

	return outTree, sendRet
}

func gasToFIL(gas msg.GasAmount, price msg.GasPrice) actor.TokenAmount {
	IMPL_FINISH()
	panic("") // BigInt arithmetic
	// return actor.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}
