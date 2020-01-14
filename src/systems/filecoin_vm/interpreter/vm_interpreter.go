package interpreter

import (
	addr "github.com/filecoin-project/go-address"
	actor "github.com/filecoin-project/specs/actors"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	initact "github.com/filecoin-project/specs/actors/builtin/init"
	sminact "github.com/filecoin-project/specs/actors/builtin/storage_miner"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	exitcode "github.com/filecoin-project/specs/actors/runtime/exitcode"
	indices "github.com/filecoin-project/specs/actors/runtime/indices"
	serde "github.com/filecoin-project/specs/actors/serde"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
	vmri "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/impl"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
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
func (vmi *VMInterpreter_I) ApplyTipSetMessages(inTree st.StateTree, tipset chain.Tipset, msgs TipSetMessages) (outTree st.StateTree, receipts []vmri.MessageReceipt) {
	outTree = inTree
	seenMsgs := make(map[cid.Cid]struct{}) // CIDs of messages already seen once.
	var receipt vmri.MessageReceipt
	store := vmi.Node().Repository().StateStore()
	// get chain from Tipset
	chainRand := &chain.Chain_I{
		HeadTipset_: tipset,
	}

	for _, blk := range msgs.Blocks() {
		minerAddr := blk.Miner()
		util.Assert(minerAddr.Protocol() == addr.ID) // Block syntactic validation requires this.

		// Process block miner's Election PoSt.
		epostMessage := _makeElectionPoStMessage(outTree, minerAddr)
		outTree = _applyMessageBuiltinAssert(store, outTree, chainRand, epostMessage, minerAddr)

		minerPenaltyTotal := abi.TokenAmount(0)
		var minerPenaltyCurr abi.TokenAmount

		// Process BLS messages from the block.
		for _, m := range blk.BLSMessages() {
			_, found := seenMsgs[_msgCID(m)]
			if found {
				continue
			}
			onChainMessageLen := len(msg.Serialize_UnsignedMessage(m))
			outTree, receipt, minerPenaltyCurr = vmi.ApplyMessage(outTree, chainRand, m, onChainMessageLen, minerAddr)
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
			outTree, receipt, minerPenaltyCurr = vmi.ApplyMessage(outTree, chainRand, m, onChainMessageLen, minerAddr)
			minerPenaltyTotal += minerPenaltyCurr
			receipts = append(receipts, receipt)
			seenMsgs[_msgCID(m)] = struct{}{}
		}

		// Pay block reward.
		rewardMessage := _makeBlockRewardMessage(outTree, minerAddr, minerPenaltyTotal)
		outTree = _applyMessageBuiltinAssert(store, outTree, chainRand, rewardMessage, minerAddr)
	}

	// Invoke cron tick.
	// Since this is outside any block, the top level block winner is declared as the system actor.
	cronMessage := _makeCronTickMessage(outTree)
	outTree = _applyMessageBuiltinAssert(store, outTree, chainRand, cronMessage, builtin.SystemActorAddr)

	return
}

func (vmi *VMInterpreter_I) ApplyMessage(inTree st.StateTree, chain chain.Chain, message msg.UnsignedMessage, onChainMessageSize int, minerAddr addr.Address) (
	retTree st.StateTree, retReceipt vmri.MessageReceipt, retMinerPenalty abi.TokenAmount) {

	store := vmi.Node().Repository().StateStore()
	senderAddr := _resolveSender(store, inTree, message.From())
	minerOwner := _lookupMinerOwner(store, inTree, minerAddr)
	Assert(minerOwner.Protocol() == addr.ID)

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
			tree = _withTransferFundsAssert(tree, builtin.BurntFundsActorAddr, senderAddr, vmiGasRemainingFIL)
			tree = _withTransferFundsAssert(tree, builtin.BurntFundsActorAddr, minerOwner, vmiGasUsedFIL)
			retMinerPenalty = abi.TokenAmount(0)

		case SenderResolveSpec_Invalid:
			retMinerPenalty = vmiGasUsedFIL

		default:
			Assert(false)
		}

		retTree = tree
		retReceipt = vmri.MessageReceipt_Make(invocOutput, exitCode, vmiGasUsed)
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

	fromActor, ok := inTree.GetActor(senderAddr)
	if !ok {
		// Execution error; sender does not exist at time of message execution.
		_applyError(inTree, exitcode.ActorNotFound, SenderResolveSpec_Invalid)
		return
	}

	// make sure this is the right message order for fromActor
	if message.CallSeqNum() != fromActor.CallSeqNum() {
		_applyError(inTree, exitcode.InvalidCallSeqNum, SenderResolveSpec_Invalid)
		return
	}

	// Check sender balance.
	gasLimitCost := _gasToFIL(message.GasLimit(), message.GasPrice())
	networkTxnFee := indicesFromStateTree(inTree).NetworkTransactionFee(
		inTree.GetActorCodeID_Assert(message.To()), message.Method())
	totalCost := message.Value() + gasLimitCost + networkTxnFee
	if fromActor.Balance() < totalCost {
		// Execution error; sender does not have sufficient funds to pay for the gas limit.
		_applyError(inTree, exitcode.InsufficientFunds_System, SenderResolveSpec_Invalid)
		return
	}

	// At this point, construct compTreePreSend as a state snapshot which includes
	// the sender paying gas, and the sender's CallSeqNum being incremented;
	// at least that much state change will be persisted even if the
	// method invocation subsequently fails.
	compTreePreSend := _withTransferFundsAssert(inTree, senderAddr, builtin.BurntFundsActorAddr, gasLimitCost+networkTxnFee)
	compTreePreSend = compTreePreSend.Impl().WithIncrementedCallSeqNum_Assert(senderAddr)

	// Reload sender actor state after the transfer and CallSeqNum increment.
	sender, ok := compTreePreSend.GetActor(senderAddr)
	Assert(ok)
	invoc := _makeInvocInput(message)
	sendRet, compTreePostSend := _applyMessageInternal(store, compTreePreSend, chain, sender, senderAddr, invoc, vmiGasRemaining, minerAddr)

	ok = _vmiBurnGas(sendRet.GasUsed)
	if !ok {
		panic("Interpreter error: runtime execution used more gas than provided")
	}

	ok = _vmiAllocGas(gascost.OnChainReturnValue(sendRet.ReturnValue))
	if !ok {
		// Insufficient gas remaining to cover the on-chain return value; proceed as in the case
		// of method execution failure.
		_applyError(compTreePreSend, exitcode.OutOfGas, SenderResolveSpec_OK)
		return
	}

	compTreeRet := compTreePreSend
	if sendRet.ExitCode.AllowsStateUpdate() {
		compTreeRet = compTreePostSend
	}

	_applyReturn(
		compTreeRet, vmr.InvocOutput_Make(sendRet.ReturnValue), sendRet.ExitCode, SenderResolveSpec_OK)
	return
}

// Resolves an address through the InitActor's map.
// Returns the resolved address (which will be an ID address) if found, else the original address.
func _resolveSender(store ipld.GraphStore, tree st.StateTree, address addr.Address) addr.Address {
	initState, ok := tree.GetActor(builtin.InitActorAddr)
	util.Assert(ok)
	serialized, ok := store.Get(cid.Cid(initState.State()))
	initSubState := initact.Deserialize_InitActorState_Assert(serialized)
	return initSubState.ResolveAddress(address)
}

// Get the owner account address associated to a given miner actor.
func _lookupMinerOwner(store ipld.GraphStore, tree st.StateTree, minerAddr addr.Address) addr.Address {
	initState, ok := tree.GetActor(minerAddr)
	util.Assert(ok)
	serialized, ok := store.Get(cid.Cid(initState.State()))
	// This tiny coupling between the VM and the storage miner actor is unfortunate.
	// It could be avoided by:
	// - paying gas rewards via the RewardActor, which can do the miner->owner lookup
	// - paying rewards to the miner actor instead of owner (with some miner->owner withdrawal mechanism)
	// - resolving all miner->owner mappings in the state before invoking interpreter (e.g. during validation) and passing them in
	// - including the owner address in block headers (and requiring it match the block's miner as part of semantic validation)
	minerSubState := sminact.Deserialize_StorageMinerActorState_Assert(serialized)
	return minerSubState.Info().Owner()
}

func _applyMessageBuiltinAssert(store ipld.GraphStore, tree st.StateTree, chain chain.Chain, message msg.UnsignedMessage, minerAddr addr.Address) st.StateTree {
	senderAddr := message.From()
	Assert(senderAddr == builtin.SystemActorAddr)
	Assert(senderAddr.Protocol() == addr.ID)
	// Note: this message CallSeqNum is never checked (b/c it's created in this file), but probably should be.
	// Since it changes state, we should be sure about the state transition.
	// Alternatively we could special-case the system actor and declare that its CallSeqNumber
	// never changes (saving us the state-change overhead).
	tree = tree.Impl().WithIncrementedCallSeqNum_Assert(senderAddr)

	sender, ok := tree.GetActor(senderAddr)
	Assert(ok)
	invoc := _makeInvocInput(message)
	retReceipt, retTree := _applyMessageInternal(store, tree, chain, sender, senderAddr, invoc, message.GasLimit(), minerAddr)
	if retReceipt.ExitCode != exitcode.OK() {
		panic("internal message application failed")
	}

	return retTree
}

func _applyMessageInternal(store ipld.GraphStore, tree st.StateTree, chain chain.Chain, sender actor.ActorState, senderAddr addr.Address, invoc vmr.InvocInput,
	gasRemainingInit msg.GasAmount, topLevelBlockWinner addr.Address) (vmri.MessageReceipt, st.StateTree) {

	rt := vmri.VMContext_Make(
		store,
		chain,
		senderAddr,
		topLevelBlockWinner,
		sender.CallSeqNum(),
		actor.CallSeqNum(0),
		tree,
		senderAddr,
		abi.TokenAmount(0),
		gasRemainingInit,
	)

	return rt.SendToplevelFromInterpreter(invoc)
}

func _withTransferFundsAssert(tree st.StateTree, from addr.Address, to addr.Address, amount abi.TokenAmount) st.StateTree {
	// TODO: assert amount nonnegative
	retTree, err := tree.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		panic("Interpreter error: insufficient funds (or transfer error) despite checks")
	} else {
		return retTree
	}
}

func indicesFromStateTree(st st.StateTree) indices.Indices {
	TODO()
	panic("")
}

func _gasToFIL(gas msg.GasAmount, price abi.TokenAmount) abi.TokenAmount {
	IMPL_FINISH()
	panic("") // BigInt arithmetic
	// return abi.TokenAmount(util.UVarint(gas) * util.UVarint(price))
}

func _makeInvocInput(message msg.UnsignedMessage) vmr.InvocInput {
	return &vmr.InvocInput_I{
		To_:     message.To(), // Receiver address is resolved during execution.
		Method_: message.Method(),
		Params_: message.Params(),
		Value_:  message.Value(),
	}
}

// Builds a message for paying block reward to a miner's owner.
func _makeBlockRewardMessage(state st.StateTree, minerAddr addr.Address, penalty abi.TokenAmount) msg.UnsignedMessage {
	params := serde.MustSerializeParams(minerAddr, penalty)
	TODO() // serialize other inputs to BlockRewardMessage or get this from query in RewardActor

	sysActor, ok := state.GetActor(builtin.SystemActorAddr)
	Assert(ok)
	return &msg.UnsignedMessage_I{
		From_:       builtin.SystemActorAddr,
		To_:         builtin.RewardActorAddr,
		Method_:     builtin.Method_RewardActor_AwardBlockReward,
		Params_:     params,
		CallSeqNum_: sysActor.CallSeqNum(),
		Value_:      0,
		GasPrice_:   0,
		GasLimit_:   msg.GasAmount_SentinelUnlimited(),
	}
}

// Builds a message for submitting ElectionPost on behalf of a miner actor.
func _makeElectionPoStMessage(state st.StateTree, minerActorAddr addr.Address) msg.UnsignedMessage {
	sysActor, ok := state.GetActor(builtin.SystemActorAddr)
	Assert(ok)
	return &msg.UnsignedMessage_I{
		From_:       builtin.SystemActorAddr,
		To_:         minerActorAddr,
		Method_:     builtin.Method_StorageMinerActor_OnVerifiedElectionPoSt,
		Params_:     nil,
		CallSeqNum_: sysActor.CallSeqNum(),
		Value_:      0,
		GasPrice_:   0,
		GasLimit_:   msg.GasAmount_SentinelUnlimited(),
	}
}

// Builds a message for invoking the cron actor tick.
func _makeCronTickMessage(state st.StateTree) msg.UnsignedMessage {
	sysActor, ok := state.GetActor(builtin.SystemActorAddr)
	Assert(ok)
	return &msg.UnsignedMessage_I{
		From_:       builtin.SystemActorAddr,
		To_:         builtin.CronActorAddr,
		Method_:     builtin.Method_CronActor_EpochTick,
		Params_:     nil,
		CallSeqNum_: sysActor.CallSeqNum(),
		Value_:      0,
		GasPrice_:   0,
		GasLimit_:   msg.GasAmount_SentinelUnlimited(),
	}
}

func _msgCID(msg msg.UnsignedMessage) cid.Cid {
	panic("TODO")
}
