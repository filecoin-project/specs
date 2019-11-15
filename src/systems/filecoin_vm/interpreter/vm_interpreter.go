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

func (vmi *VMInterpreter_I) ApplyMessageBatch(inTree st.StateTree, msgs []MessageRef) (outTree st.StateTree, ret []msg.MessageReceipt) {
	compTree := inTree
	for _, m := range msgs {
		oT, r := vmi.ApplyMessage(compTree, m.Message(), m.Miner())
		compTree = oT        // assign the current tree. (this call always succeeds)
		ret = append(ret, r) // add message receipt
	}
	return compTree, ret
}

func _applyError(errCode exitcode.SystemErrorCode) msg.MessageReceipt {
	// TODO: should this gasUsed value be zero?
	// If nonzero, there is not guaranteed to be a nonzero gas balance from which to deduct it.
	gasUsed := gascost.ApplyMessageFail
	TODO()
	return msg.MessageReceipt_MakeSystemError(errCode, gasUsed)
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

func (vmi *VMInterpreter_I) ApplyMessage(inTree st.StateTree, message msg.UnsignedMessage, minerAddr addr.Address) (
	st.StateTree, msg.MessageReceipt) {

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
	maxGasCost := gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(maxGasCost)
	if fromActor.Balance() < totalCost {
		return inTree, _applyError(exitcode.InsufficientFunds)
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
		msg.InvocInput_Make(
			message.To(),
			message.Method(),
			message.Params(),
			message.Value(),
		),
		false,
	)

	if !sendRet.ExitCode().AllowsStateUpdate() {
		// error -- revert all state changes -- ie drop updates. burn used gas.
		outTree = inTree
		outTree = _withTransferFundsAssert(
			outTree,
			message.From(),
			addr.BurntFundsActorAddr,
			gasToFIL(sendRet.GasUsed(), message.GasPrice()),
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
			gasToFIL(refundGas, message.GasPrice()),
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
		gasToFIL(sendRet.GasUsed(), message.GasPrice()),
	)

	return outTree, sendRet
}

func gasToFIL(gas msg.GasAmount, price msg.GasPrice) actor.TokenAmount {
	return actor.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}
