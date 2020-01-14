package impl

import (
	"bytes"
	"encoding/binary"
	"fmt"

	addr "github.com/filecoin-project/go-address"
	actor "github.com/filecoin-project/specs/actors"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	acctact "github.com/filecoin-project/specs/actors/builtin/account"
	initact "github.com/filecoin-project/specs/actors/builtin/init"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	exitcode "github.com/filecoin-project/specs/actors/runtime/exitcode"
	indices "github.com/filecoin-project/specs/actors/runtime/indices"
	ipld "github.com/filecoin-project/specs/libraries/ipld"
	chain "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/chain"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
	cid "github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	mh "github.com/multiformats/go-multihash"
)

type ActorSubstateCID = actor.ActorSubstateCID
type ExitCode = exitcode.ExitCode
type CallerPattern = vmr.CallerPattern
type Runtime = vmr.Runtime
type InvocInput = vmr.InvocInput
type InvocOutput = vmr.InvocOutput
type ActorStateHandle = vmr.ActorStateHandle

var EnsureErrorCode = exitcode.EnsureErrorCode
var SystemError = exitcode.SystemError

type Bytes = util.Bytes

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH
var IMPL_TODO = util.IMPL_TODO
var TODO = util.TODO

var EmptyCBOR cid.Cid

type RuntimeError struct {
	ExitCode ExitCode
	ErrMsg   string
}

func init() {
	n, err := cbornode.WrapObject(map[string]struct{}{}, mh.SHA2_256, -1)
	Assert(err == nil)
	EmptyCBOR = n.Cid()
}

func (x *RuntimeError) String() string {
	ret := fmt.Sprintf("Runtime error: %v", x.ExitCode)
	if x.ErrMsg != "" {
		ret += fmt.Sprintf(" (\"%v\")", x.ErrMsg)
	}
	return ret
}

func RuntimeError_Make(exitCode ExitCode, errMsg string) *RuntimeError {
	exitCode = EnsureErrorCode(exitCode)
	return &RuntimeError{
		ExitCode: exitCode,
		ErrMsg:   errMsg,
	}
}

func ActorSubstateCID_Equals(x, y ActorSubstateCID) bool {
	IMPL_FINISH()
	panic("")
}

type ActorStateHandle_I struct {
	_initValue *ActorSubstateCID
	_rt        *VMContext
}

func (h *ActorStateHandle_I) UpdateRelease(newStateCID ActorSubstateCID) {
	h._rt._updateReleaseActorSubstate(newStateCID)
}

func (h *ActorStateHandle_I) Release(checkStateCID ActorSubstateCID) {
	h._rt._releaseActorSubstate(checkStateCID)
}

func (h *ActorStateHandle_I) Take() ActorSubstateCID {
	if h._initValue == nil {
		h._rt._apiError("Must call Take() only once on actor substate object")
	}
	ret := *h._initValue
	h._initValue = nil
	return ret
}

// Concrete instantiation of the Runtime interface. This should be instantiated by the
// interpreter once per actor method invocation, and responds to that method's Runtime
// API calls.
type VMContext struct {
	_store              ipld.GraphStore
	_globalStateInit    st.StateTree
	_globalStatePending st.StateTree
	_running            bool
	_chain              chain.Chain
	_actorAddress       addr.Address
	_actorStateAcquired bool
	// Tracks whether actor substate has changed in order to charge gas just once
	// regardless of how many times it's written.
	_actorSubstateUpdated bool

	_immediateCaller addr.Address
	// Note: This is the actor in the From field of the initial on-chain message.
	// Not necessarily the immediate caller.
	_toplevelSender      addr.Address
	_toplevelBlockWinner addr.Address
	// Top-level call sequence number of the "From" actor in the initial on-chain message.
	_toplevelSenderCallSeqNum actor.CallSeqNum
	// Sequence number representing the total number of calls (to any actor, any method)
	// during the current top-level message execution.
	// Note: resets with every top-level message, and therefore not necessarily monotonic.
	_internalCallSeqNum actor.CallSeqNum
	_valueReceived      abi.TokenAmount
	_gasRemaining       msg.GasAmount
	_numValidateCalls   int
	_output             vmr.InvocOutput
}

func VMContext_Make(
	store ipld.GraphStore,
	chain chain.Chain,
	toplevelSender addr.Address,
	toplevelBlockWinner addr.Address,
	toplevelSenderCallSeqNum actor.CallSeqNum,
	internalCallSeqNum actor.CallSeqNum,
	globalState st.StateTree,
	actorAddress addr.Address,
	valueReceived abi.TokenAmount,
	gasRemaining msg.GasAmount) *VMContext {

	return &VMContext{
		_store:                store,
		_chain:                chain,
		_globalStateInit:      globalState,
		_globalStatePending:   globalState,
		_running:              false,
		_actorAddress:         actorAddress,
		_actorStateAcquired:   false,
		_actorSubstateUpdated: false,

		_toplevelSender:           toplevelSender,
		_toplevelBlockWinner:      toplevelBlockWinner,
		_toplevelSenderCallSeqNum: toplevelSenderCallSeqNum,
		_internalCallSeqNum:       internalCallSeqNum,
		_valueReceived:            valueReceived,
		_gasRemaining:             gasRemaining,
		_numValidateCalls:         0,
		_output:                   nil,
	}
}

func (rt *VMContext) AbortArgMsg(msg string) {
	rt.Abort(exitcode.UserDefinedError(exitcode.InvalidArguments_User), msg)
}

func (rt *VMContext) AbortArg() {
	rt.AbortArgMsg("Invalid arguments")
}

func (rt *VMContext) AbortStateMsg(msg string) {
	rt.Abort(exitcode.UserDefinedError(exitcode.InconsistentState_User), msg)
}

func (rt *VMContext) AbortState() {
	rt.AbortStateMsg("Inconsistent state")
}

func (rt *VMContext) AbortFundsMsg(msg string) {
	rt.Abort(exitcode.UserDefinedError(exitcode.InsufficientFunds_User), msg)
}

func (rt *VMContext) AbortFunds() {
	rt.AbortFundsMsg("Insufficient funds")
}

func (rt *VMContext) AbortAPI(msg string) {
	rt.Abort(exitcode.SystemError(exitcode.RuntimeAPIError), msg)
}

func (rt *VMContext) CreateActor(codeID abi.ActorCodeID, address addr.Address) {
	if rt._actorAddress != builtin.InitActorAddr {
		rt.AbortAPI("Only InitActor may call rt.CreateActor")
	}
	if address.Protocol() != addr.ID {
		rt.AbortAPI("New actor adddress must be an ID-address")
	}

	rt._createActor(codeID, address)
}

func (rt *VMContext) _createActor(codeID abi.ActorCodeID, address addr.Address) {
	// Create empty actor state.
	actorState := &actor.ActorState_I{
		CodeID_:     codeID,
		State_:      actor.ActorSubstateCID(EmptyCBOR),
		Balance_:    abi.TokenAmount(0),
		CallSeqNum_: 0,
	}

	// Put it in the state tree.
	actorStateCID := actor.ActorSystemStateCID(rt.IpldPut(actorState))
	rt._updateActorSystemStateInternal(address, actorStateCID)

	rt._rtAllocGas(gascost.ExecNewActor)
}

func (rt *VMContext) DeleteActor(address addr.Address) {
	// Only a given actor may delete itself.
	if rt._actorAddress != address {
		rt.AbortAPI("Invalid actor deletion request")
	}

	rt._deleteActor(address)
}

func (rt *VMContext) _deleteActor(address addr.Address) {
	rt._globalStatePending = rt._globalStatePending.Impl().WithDeleteActorSystemState(address)
	rt._rtAllocGas(gascost.DeleteActor)
}

func (rt *VMContext) _updateActorSystemStateInternal(actorAddress addr.Address, newStateCID actor.ActorSystemStateCID) {
	newGlobalStatePending, err := rt._globalStatePending.Impl().WithActorSystemState(rt._actorAddress, newStateCID)
	if err != nil {
		panic("Error in runtime implementation: failed to update actor system state")
	}
	rt._globalStatePending = newGlobalStatePending
}

func (rt *VMContext) _updateActorSubstateInternal(actorAddress addr.Address, newStateCID actor.ActorSubstateCID) {
	newGlobalStatePending, err := rt._globalStatePending.Impl().WithActorSubstate(rt._actorAddress, newStateCID)
	if err != nil {
		panic("Error in runtime implementation: failed to update actor substate")
	}
	rt._globalStatePending = newGlobalStatePending
}

func (rt *VMContext) _updateReleaseActorSubstate(newStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()
	rt._updateActorSubstateInternal(rt._actorAddress, newStateCID)
	rt._actorSubstateUpdated = true
	rt._actorStateAcquired = false
}

func (rt *VMContext) _releaseActorSubstate(checkStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()

	prevState, ok := rt._globalStatePending.GetActor(rt._actorAddress)
	util.Assert(ok)
	prevStateCID := prevState.State()
	if !ActorSubstateCID_Equals(prevStateCID, checkStateCID) {
		rt.AbortAPI("State CID differs upon release call")
	}

	rt._actorStateAcquired = false
}

func (rt *VMContext) Assert(cond bool) {
	if !cond {
		rt.Abort(exitcode.SystemError(exitcode.RuntimeAssertFailure), "Runtime assertion failed")
	}
}

func (rt *VMContext) _checkActorStateAcquiredFlag(expected bool) {
	rt._checkRunning()
	if rt._actorStateAcquired != expected {
		rt._apiError("State updates and message sends must be disjoint")
	}
}

func (rt *VMContext) _checkActorStateAcquired() {
	rt._checkActorStateAcquiredFlag(true)
}

func (rt *VMContext) _checkActorStateNotAcquired() {
	rt._checkActorStateAcquiredFlag(false)
}

func (rt *VMContext) Abort(errExitCode exitcode.ExitCode, errMsg string) {
	errExitCode = exitcode.EnsureErrorCode(errExitCode)
	rt._throwErrorFull(errExitCode, errMsg)
}

func (rt *VMContext) ImmediateCaller() addr.Address {
	return rt._immediateCaller
}

func (rt *VMContext) CurrReceiver() addr.Address {
	return rt._actorAddress
}

func (rt *VMContext) ToplevelBlockWinner() addr.Address {
	return rt._toplevelBlockWinner
}

func (rt *VMContext) ValidateImmediateCallerMatches(
	callerExpectedPattern CallerPattern) {

	rt._checkRunning()
	rt._checkNumValidateCalls(0)
	caller := rt.ImmediateCaller()
	if !callerExpectedPattern.Matches(caller) {
		rt.AbortAPI("Method invoked by incorrect caller")
	}
	rt._numValidateCalls += 1
}

func CallerPattern_MakeAcceptAnyOfTypes(rt *VMContext, types []abi.ActorCodeID) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool {
			codeID, ok := rt.GetActorCodeID(y)
			if !ok {
				panic("Internal runtime error: actor not found")
			}

			for _, type_ := range types {
				if codeID == type_ {
					return true
				}
			}
			return false
		},
	}
}

func (rt *VMContext) ValidateImmediateCallerIs(callerExpected addr.Address) {
	rt.ValidateImmediateCallerMatches(vmr.CallerPattern_MakeSingleton(callerExpected))
}

func (rt *VMContext) ValidateImmediateCallerInSet(callersExpected []addr.Address) {
	rt.ValidateImmediateCallerMatches(vmr.CallerPattern_MakeSet(callersExpected))
}

func (rt *VMContext) ValidateImmediateCallerAcceptAnyOfType(type_ abi.ActorCodeID) {
	rt.ValidateImmediateCallerAcceptAnyOfTypes([]abi.ActorCodeID{type_})
}

func (rt *VMContext) ValidateImmediateCallerAcceptAnyOfTypes(types []abi.ActorCodeID) {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeAcceptAnyOfTypes(rt, types))
}

func (rt *VMContext) ValidateImmediateCallerAcceptAny() {
	rt.ValidateImmediateCallerMatches(vmr.CallerPattern_MakeAcceptAny())
}

func (rt *VMContext) _checkNumValidateCalls(x int) {
	if rt._numValidateCalls != x {
		rt.AbortAPI("Method must validate caller identity exactly once")
	}
}

func (rt *VMContext) _checkRunning() {
	if !rt._running {
		panic("Internal runtime error: actor API called with no actor code running")
	}
}
func (rt *VMContext) SuccessReturn() InvocOutput {
	return vmr.InvocOutput_Make(nil)
}

func (rt *VMContext) ValueReturn(value util.Bytes) InvocOutput {
	return vmr.InvocOutput_Make(value)
}

func (rt *VMContext) _throwError(exitCode ExitCode) {
	rt._throwErrorFull(exitCode, "")
}

func (rt *VMContext) _throwErrorFull(exitCode ExitCode, errMsg string) {
	panic(RuntimeError_Make(exitCode, errMsg))
}

func (rt *VMContext) _apiError(errMsg string) {
	rt._throwErrorFull(exitcode.SystemError(exitcode.RuntimeAPIError), errMsg)
}

func _gasAmountAssertValid(x msg.GasAmount) {
	if x.LessThan(msg.GasAmount_Zero()) {
		panic("Interpreter error: negative gas amount")
	}
}

// Deduct an amount of gas corresponding to cost about to be incurred, but not necessarily
// incurred yet.
func (rt *VMContext) _rtAllocGas(x msg.GasAmount) {
	_gasAmountAssertValid(x)
	var ok bool
	rt._gasRemaining, ok = rt._gasRemaining.SubtractIfNonnegative(x)
	if !ok {
		rt._throwError(exitcode.SystemError(exitcode.OutOfGas))
	}
}

func (rt *VMContext) _transferFunds(from addr.Address, to addr.Address, amount abi.TokenAmount) error {
	rt._checkRunning()
	rt._checkActorStateNotAcquired()

	newGlobalStatePending, err := rt._globalStatePending.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		return err
	}

	rt._globalStatePending = newGlobalStatePending
	return nil
}

func (rt *VMContext) GetActorCodeID(actorAddr addr.Address) (ret abi.ActorCodeID, ok bool) {
	IMPL_FINISH()
	panic("")
}

type ErrorHandlingSpec int

const (
	PropagateErrors ErrorHandlingSpec = 1 + iota
	CatchErrors
)

// TODO: This function should be private (not intended to be exposed to actors).
// (merging runtime and interpreter packages should solve this)
// TODO: this should not use the MessageReceipt return type, even though it needs the same triple
// of values. This method cannot compute the total gas cost and the returned receipt will never
// go on chain.
func (rt *VMContext) SendToplevelFromInterpreter(input InvocInput) (MessageReceipt, st.StateTree) {

	rt._running = true
	ret := rt._sendInternal(input, CatchErrors)
	rt._running = false
	return ret, rt._globalStatePending
}

func _catchRuntimeErrors(f func() InvocOutput) (output InvocOutput, exitCode exitcode.ExitCode) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				output = vmr.InvocOutput_Make(nil)
				exitCode = (r.(*RuntimeError).ExitCode)
			default:
				panic(r)
			}
		}
	}()

	output = f()
	exitCode = exitcode.OK()
	return
}

func _invokeMethodInternal(
	rt *VMContext,
	actorCode vmr.ActorCode,
	method abi.MethodNum,
	params abi.MethodParams) (
	ret InvocOutput, exitCode exitcode.ExitCode, internalCallSeqNumFinal actor.CallSeqNum) {

	if method == builtin.MethodSend {
		ret = vmr.InvocOutput_Make(nil)
		return
	}

	rt._running = true
	ret, exitCode = _catchRuntimeErrors(func() InvocOutput {
		IMPL_TODO("dispatch to actor code")
		var methodOutput vmr.InvocOutput // actorCode.InvokeMethod(rt, method, params)
		if rt._actorSubstateUpdated {
			rt._rtAllocGas(gascost.UpdateActorSubstate)
		}
		rt._checkActorStateNotAcquired()
		rt._checkNumValidateCalls(1)
		return methodOutput
	})
	rt._running = false

	internalCallSeqNumFinal = rt._internalCallSeqNum

	return
}

func (rtOuter *VMContext) _sendInternal(input InvocInput, errSpec ErrorHandlingSpec) MessageReceipt {
	rtOuter._checkRunning()
	rtOuter._checkActorStateNotAcquired()

	initGasRemaining := rtOuter._gasRemaining

	rtOuter._rtAllocGas(gascost.InvokeMethod(input.Value(), input.Method()))

	receiver, receiverAddr := rtOuter._resolveReceiver(input.To())
	receiverCode, err := loadActorCode(receiver.CodeID())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.ActorCodeNotFound))
	}

	err = rtOuter._transferFunds(rtOuter._actorAddress, receiverAddr, input.Value())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.InsufficientFunds_System))
	}

	rtInner := VMContext_Make(
		rtOuter._store,
		rtOuter._chain,
		rtOuter._toplevelSender,
		rtOuter._toplevelBlockWinner,
		rtOuter._toplevelSenderCallSeqNum,
		rtOuter._internalCallSeqNum+1,
		rtOuter._globalStatePending,
		receiverAddr,
		input.Value(),
		rtOuter._gasRemaining,
	)

	invocOutput, exitCode, internalCallSeqNumFinal := _invokeMethodInternal(
		rtInner,
		receiverCode,
		input.Method(),
		input.Params(),
	)

	_gasAmountAssertValid(rtOuter._gasRemaining.Subtract(rtInner._gasRemaining))
	rtOuter._gasRemaining = rtInner._gasRemaining
	gasUsed := initGasRemaining.Subtract(rtOuter._gasRemaining)
	_gasAmountAssertValid(gasUsed)

	rtOuter._internalCallSeqNum = internalCallSeqNumFinal

	if exitCode.Equals(exitcode.SystemError(exitcode.OutOfGas)) {
		// OutOfGas error cannot be caught
		rtOuter._throwError(exitCode)
	}

	if errSpec == PropagateErrors && exitCode.IsError() {
		rtOuter._throwError(exitcode.SystemError(exitcode.MethodSubcallError))
	}

	if exitCode.AllowsStateUpdate() {
		rtOuter._globalStatePending = rtInner._globalStatePending
	}

	return MessageReceipt_Make(invocOutput, exitCode, gasUsed)
}

// Loads a receiving actor state from the state tree, resolving non-ID addresses through the InitActor state.
// If it doesn't exist, and the message is a simple value send to a pubkey-style address,
// creates the receiver as an account actor in the returned state.
// Aborts otherwise.
func (rt *VMContext) _resolveReceiver(targetRaw addr.Address) (actor.ActorState, addr.Address) {
	// Resolve the target address via the InitActor, and attempt to load state.
	initSubState := rt._loadInitActorState()
	targetIdAddr := initSubState.ResolveAddress(targetRaw)
	act, found := rt._globalStatePending.GetActor(targetIdAddr)
	if found {
		return act, targetIdAddr
	}

	if targetRaw.Protocol() != addr.SECP256K1 && targetRaw.Protocol() != addr.BLS {
		// Don't implicitly create an account actor for an address without an associated key.
		rt._throwError(exitcode.SystemError(exitcode.ActorNotFound))
	}

	// Allocate an ID address from the init actor and map the pubkey To address to it.
	newIdAddr := initSubState.MapAddressToNewID(targetRaw)
	rt._saveInitActorState(initSubState)

	// Create new account actor (charges gas).
	rt._createActor(builtin.AccountActorCodeID, newIdAddr)

	// Initialize account actor substate with it's pubkey address.
	substate := &acctact.AccountActorState_I{
		Address_: targetRaw,
	}
	rt._saveAccountActorState(newIdAddr, substate)
	act, _ = rt._globalStatePending.GetActor(newIdAddr)
	return act, newIdAddr
}

func (rt *VMContext) _loadInitActorState() initact.InitActorState {
	initState, ok := rt._globalStatePending.GetActor(builtin.InitActorAddr)
	util.Assert(ok)
	var initSubState initact.InitActorState_I
	ok = rt.IpldGet(cid.Cid(initState.State()), &initSubState)
	util.Assert(ok)
	return &initSubState
}

func (rt *VMContext) _saveInitActorState(state initact.InitActorState) {
	// Gas is charged here separately from _actorSubstateUpdated because this is a different actor
	// than the receiver.
	rt._rtAllocGas(gascost.UpdateActorSubstate)
	rt._updateActorSubstateInternal(builtin.InitActorAddr, actor.ActorSubstateCID(rt.IpldPut(state.Impl())))
}

func (rt *VMContext) _saveAccountActorState(address addr.Address, state acctact.AccountActorState) {
	// Gas is charged here separately from _actorSubstateUpdated because this is a different actor
	// than the receiver.
	rt._rtAllocGas(gascost.UpdateActorSubstate)
	rt._updateActorSubstateInternal(address, actor.ActorSubstateCID(rt.IpldPut(state.Impl())))
}

func (rt *VMContext) _sendInternalOutputs(input InvocInput, errSpec ErrorHandlingSpec) (InvocOutput, exitcode.ExitCode) {
	ret := rt._sendInternal(input, errSpec)
	return vmr.InvocOutput_Make(ret.ReturnValue), ret.ExitCode
}

func (rt *VMContext) Send(
	toAddr addr.Address, methodNum abi.MethodNum, params abi.MethodParams, value abi.TokenAmount) InvocOutput {

	return rt.SendPropagatingErrors(vmr.InvocInput_Make(toAddr, methodNum, params, value))
}

func (rt *VMContext) SendQuery(toAddr addr.Address, methodNum abi.MethodNum, params abi.MethodParams) util.Serialization {
	invocOutput := rt.Send(toAddr, methodNum, params, abi.TokenAmount(0))
	ret := invocOutput.ReturnValue()
	Assert(ret != nil)
	return ret
}

func (rt *VMContext) SendFunds(toAddr addr.Address, value abi.TokenAmount) {
	rt.Send(toAddr, builtin.MethodSend, nil, value)
}

func (rt *VMContext) SendPropagatingErrors(input InvocInput) InvocOutput {
	ret, _ := rt._sendInternalOutputs(input, PropagateErrors)
	return ret
}

func (rt *VMContext) SendCatchingErrors(input InvocInput) (InvocOutput, exitcode.ExitCode) {
	rt.ValidateImmediateCallerIs(builtin.CronActorAddr)
	return rt._sendInternalOutputs(input, CatchErrors)
}

func (rt *VMContext) CurrentBalance() abi.TokenAmount {
	IMPL_FINISH()
	panic("")
}

func (rt *VMContext) ValueReceived() abi.TokenAmount {
	return rt._valueReceived
}

func (rt *VMContext) GetRandomness(epoch abi.ChainEpoch) abi.RandomnessSeed {
	return rt._chain.RandomnessAtEpoch(epoch)
}

func (rt *VMContext) NewActorAddress() addr.Address {
	addrBuf := new(bytes.Buffer)

	err := rt._immediateCaller.MarshalCBOR(addrBuf)
	util.Assert(err == nil)
	err = binary.Write(addrBuf, binary.BigEndian, rt._toplevelSenderCallSeqNum)
	util.Assert(err != nil)
	err = binary.Write(addrBuf, binary.BigEndian, rt._internalCallSeqNum)
	util.Assert(err != nil)

	newAddr, err := addr.NewActorAddress(addrBuf.Bytes())
	util.Assert(err == nil)
	return newAddr
}

func (rt *VMContext) IpldPut(x ipld.Object) cid.Cid {
	IMPL_FINISH() // Serialization
	serialized := []byte{}
	cid := rt._store.Put(serialized)
	rt._rtAllocGas(gascost.IpldPut(len(serialized)))
	return cid
}

func (rt *VMContext) IpldGet(c cid.Cid, o ipld.Object) bool {
	serialized, ok := rt._store.Get(c)
	if ok {
		rt._rtAllocGas(gascost.IpldGet(len(serialized)))
	}
	IMPL_FINISH() // Deserialization into o
	return ok
}

func (rt *VMContext) CurrEpoch() abi.ChainEpoch {
	IMPL_FINISH()
	panic("")
}

func (rt *VMContext) CurrIndices() indices.Indices {
	// TODO: compute from state tree (rt._globalStatePending), using individual actor
	// state helper functions when possible
	TODO()
	panic("")
}

func (rt *VMContext) AcquireState() ActorStateHandle {
	rt._checkRunning()
	rt._checkActorStateNotAcquired()
	rt._actorStateAcquired = true

	state, ok := rt._globalStatePending.GetActor(rt._actorAddress)
	util.Assert(ok)
	return &ActorStateHandle_I{
		_initValue: state.State().Ref(),
		_rt:        rt,
	}
}

func (rt *VMContext) Compute(f ComputeFunctionID, args []util.Any) Any {
	def, found := _computeFunctionDefs[f]
	if !found {
		rt.AbortAPI("Function definition in rt.Compute() not found")
	}
	gasCost := def.GasCostFn(args)
	rt._rtAllocGas(gasCost)
	return def.Body(args)
}
