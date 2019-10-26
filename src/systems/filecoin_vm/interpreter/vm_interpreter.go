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
	transferRet := tree.WithFundsTransfer(from, to, amount)
	if transferRet.Which() == st.StateTree_WithFundsTransfer_FunRet_Case_error {
		panic("Interpreter error: insufficient funds (or transfer error) despite checks")
	} else {
		return transferRet.As_StateTree()
	}
}

func (vmi *VMInterpreter_I) ApplyMessage(inTree st.StateTree, message msg.UnsignedMessage, minerAddr addr.Address) (
	st.StateTree, msg.MessageReceipt) {

	compTree := inTree
	var outTree st.StateTree

	fromActor := compTree.GetActor(message.From())
	if fromActor == nil {
		// TODO: This was originally exitcode.InvalidMethod; which is correct?
		return inTree, _applyError(exitcode.ActorNotFound)
	}

	// make sure fromActor has enough money to run the max invocation
	maxGasCost := gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(maxGasCost)
	if fromActor.State().Balance() < totalCost {
		return inTree, _applyError(exitcode.InsufficientFunds)
	}

	// make sure this is the right message order for fromActor
	// (this is protection against replay attacks, and useful sequencing)
	if message.CallSeqNum() != fromActor.State().CallSeqNum()+1 {
		return inTree, _applyError(exitcode.InvalidCallSeqNum)
	}

	// TODO: Should we allow messages to implicitly create actors?
	// - If so, need to specify how the actor code is determined
	// - We should be able to disallow this and still preserve functionality by
	//   requiring two separate messages (create actor, transfer initial funds)
	createRet := compTree.WithNewActorIfMissing(message.To())
	if createRet.Which() == st.StateTree_WithNewActorIfMissing_FunRet_Case_error {
		return inTree, _applyError(exitcode.ActorNotFound)
	} else {
		compTree = createRet.As_StateTree()
	}

	toActor := compTree.GetActor(message.To())
	if toActor == nil {
		panic("Interpreter error: actor present (or created) but not retrieved")
	}

	// deduct maximum expenditure gas funds first
	compTree = _withTransferFundsAssert(compTree, message.From(), vmr.BurntFundsActorAddr, maxGasCost)

	rt := vmr.Runtime_Make(
		compTree,
		message.From(),
		actor.TokenAmount(0),
		message.GasLimit(),
	)

	sendRet, sendRetStateTree := rt.Impl().SendToplevelFromInterpreter(
		message.To(),
		msg.InvocInput_Make(
			message.Method(),
			message.Params(),
			message.Value(),
		),
		false,
	)

	if sendRet.ExitCode().AllowsStateUpdate() {
		// success -- refund unused gas
		outTree = sendRetStateTree
		refundGas := message.GasLimit() - sendRet.GasUsed()
		TODO() // TODO: assert refundGas is nonnegative
		outTree = _withTransferFundsAssert(
			outTree,
			vmr.BurntFundsActorAddr,
			message.From(),
			gasToFIL(refundGas, message.GasPrice()),
		)
	} else {
		// error -- revert all state changes -- ie drop updates. burn used gas.
		outTree = inTree
		outTree = _withTransferFundsAssert(
			outTree,
			message.From(),
			vmr.BurntFundsActorAddr,
			gasToFIL(sendRet.GasUsed(), message.GasPrice()),
		)
		TODO() // TODO: still increment fromActor sequence number on failure?
	}

	// reward miner gas fees
	outTree = _withTransferFundsAssert(
		outTree,
		vmr.BurntFundsActorAddr,
		minerAddr, // TODO: may not exist
		gasToFIL(sendRet.GasUsed(), message.GasPrice()),
	)

	return outTree, sendRet
}

func gasToFIL(gas msg.GasAmount, price msg.GasPrice) actor.TokenAmount {
	return actor.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}
