package interpreter

import "errors"
import msg "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/actor"
import st "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/state_tree"
import vmr "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime/exitcode"
import gascost "github.com/filecoin-project/specs/systems/filecoin_blockchain/vm/runtime/gascost"
import util "github.com/filecoin-project/specs/util"

func (vmi *VMInterpreter_I) ApplyMessageBatch(inTree st.StateTree, msgs []MessageRef) (outTree st.StateTree, ret []msg.MessageReceipt) {
	compTree := inTree
	for _, m := range msgs {
		oT, r := vmi.ApplyMessage(compTree, m.Message(), m.Miner())
		compTree = oT        // assign the current tree. (this call always succeeds)
		ret = append(ret, r) // add message receipt
	}
	return compTree, ret
}

func (vmi *VMInterpreter_I) ApplyMessage(inTree st.StateTree, message msg.UnsignedMessage, minerAddr addr.Address) (outTree st.StateTree, ret msg.MessageReceipt) {

	compTree := inTree
	fromActor := compTree.GetActor(message.From())
	if fromActor == nil {
		return applyError(inTree, exitcode.InvalidMethod, gascost.ApplyMessageFail)
	}

	// make sure fromActor has enough money to run the max invocation
	maxGasCost := gasToFIL(message.GasLimit(), message.GasPrice())
	totalCost := message.Value() + actor.TokenAmount(maxGasCost)
	if fromActor.State().Balance() < totalCost {
		return applyError(inTree, exitcode.InsufficientFunds, gascost.ApplyMessageFail)
	}

	// make sure this is the right message order for fromActor
	// (this is protection against replay attacks, and useful sequencing)
	if message.CallSeqNum() != fromActor.State().CallSeqNum()+1 {
		return applyError(inTree, exitcode.InvalidCallSeqNum, gascost.ApplyMessageFail)
	}

	// may return a different tree on succeess.
	// this MUST get rolled back if the invocation fails.
	var toActor actor.Actor
	var err error
	compTree, toActor, err = treeGetOrCreateAccountActor(compTree, message.To())
	if err != nil {
		return applyError(inTree, exitcode.ActorNotFound, gascost.ApplyMessageFail)
	}

	// deduct maximum expenditure gas funds first
	// TODO: use a single "transfer"
	compTree = treeDeductFunds(compTree, fromActor, maxGasCost)

	// transfer funds fromActor -> toActor
	// (yes deductions can be combined, spelled out here for clarity)
	// TODO: use a single "transfer"
	compTree = treeDeductFunds(compTree, fromActor, message.Value())
	compTree = treeDepositFunds(compTree, toActor, message.Value())

	// Prepare invocInput.
	invocInput := vmr.InvocInput{
		InTree:    compTree,
		FromActor: fromActor,
		ToActor:   toActor,
		Method:    message.Method(),
		Params:    message.Params(),
		Value:     message.Value(),
		GasLimit:  message.GasLimit(),
	}
	// TODO: this is mega jank. need to rework invocationInput + runtime boundaries.
	invocInput.Runtime = makeRuntime(compTree, invocInput)

	// perform the method call to the actor
	// TODO: eval if we should lift gas tracking and calc to the beginning of invocation
	// (ie, include account creation, gas accounting itself)
	out := invocationMethodDispatch(invocInput)

	// var outTree StateTree
	if out.ExitCode != 0 {
		// error -- revert all state changes -- ie drop updates. burn used gas.
		outTree = inTree // wipe!
		outTree = treeDeductFunds(outTree, fromActor, gasToFIL(out.GasUsed, message.GasPrice()))

	} else {
		// success -- refund unused gas
		outTree = out.OutTree // take the state from the invocation output
		refundGas := message.GasLimit() - out.GasUsed
		outTree = treeDepositFunds(outTree, fromActor, gasToFIL(refundGas, message.GasPrice()))
		outTree = treeIncrementActorSeqNo(outTree, fromActor)
	}

	// reward miner gas fees
	minerActor := compTree.GetActor(minerAddr) // TODO: may be nil.
	outTree = treeDepositFunds(outTree, minerActor, gasToFIL(out.GasUsed, message.GasPrice()))

	return outTree, &msg.MessageReceipt_I{
		ExitCode_:    out.ExitCode,
		ReturnValue_: out.ReturnValue,
		GasUsed_:     out.GasUsed,
	}
}

func invocationMethodDispatch(input vmr.InvocInput) vmr.InvocOutput {
	if input.Method == 0 {
		// just sending money. move along.
		return vmr.InvocOutput{
			OutTree:     input.InTree,
			GasUsed:     gascost.SimpleValueSend,
			ExitCode:    exitcode.OK,
			ReturnValue: nil,
		}
	}

	//TODO: actually invoke the funtion here.
	// put any vtable lookups in this function.

	actorCode, err := loadActorCode(input, input.ToActor.State().CodeCID())
	if err != nil {
		return vmr.InvocOutput{
			OutTree:     input.InTree,
			GasUsed:     gascost.ApplyMessageFail,
			ExitCode:    exitcode.ActorCodeNotFound,
			ReturnValue: nil, // TODO: maybe: err
		}
	}

	return actorCode.InvokeMethod(input, input.Method, input.Params)
}

func treeIncrementActorSeqNo(inTree st.StateTree, a actor.Actor) (outTree st.StateTree) {
	panic("todo")
}

func treeDeductFunds(inTree st.StateTree, a actor.Actor, amt actor.TokenAmount) (outTree st.StateTree) {
	// TODO: turn this into a single transfer call.
	panic("todo")
}

func treeDepositFunds(inTree st.StateTree, a actor.Actor, amt actor.TokenAmount) (outTree st.StateTree) {
	// TODO: turn this into a single transfer call.
	panic("todo")
}

func treeGetOrCreateAccountActor(inTree st.StateTree, a addr.Address) (outTree st.StateTree, _ actor.Actor, err error) {

	toActor := inTree.GetActor(a)
	if toActor != nil { // found
		return inTree, toActor, nil
	}

	switch a.Type().Which() {
	case addr.Address_Type_Case_BLS:
		return treeNewBLSAccountActor(inTree, a)
	case addr.Address_Type_Case_Secp256k1:
		return treeNewSecp256k1AccountActor(inTree, a)
	case addr.Address_Type_Case_ID:
		return inTree, nil, errors.New("no actor with given ID")
	case addr.Address_Type_Case_Actor:
		return inTree, nil, errors.New("no such actor")
	default:
		return inTree, nil, errors.New("unknown address type")
	}
}

func treeNewBLSAccountActor(inTree st.StateTree, addr addr.Address) (outTree st.StateTree, _ actor.Actor, err error) {
	panic("todo")
}

func treeNewSecp256k1AccountActor(inTree st.StateTree, addr addr.Address) (outTree st.StateTree, _ actor.Actor, err error) {
	panic("todo")
}

func applyError(tree st.StateTree, exitCode msg.ExitCode, gasUsed msg.GasAmount) (outTree st.StateTree, ret msg.MessageReceipt) {
	return outTree, &msg.MessageReceipt_I{
		ExitCode_:    exitCode,
		ReturnValue_: nil,
		GasUsed_:     gasUsed,
	}
}

func gasToFIL(gas msg.GasAmount, price msg.GasPrice) actor.TokenAmount {
	return actor.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}

func makeRuntime(tree st.StateTree, input vmr.InvocInput) vmr.Runtime {
	return &vmr.Runtime_I{
		Invocation_: input,
		State_: &vmr.VMState_I{
			StateTree_: tree, // TODO: also in input.InTree.
			Storage_:   &vmr.VMStorage_I{},
		},
	}
}
