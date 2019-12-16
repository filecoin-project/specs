package impl

import (
	"bytes"
	"encoding/binary"

	ipld "github.com/filecoin-project/specs/libraries/ipld"
	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
	gascost "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/gascost"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	sysactors "github.com/filecoin-project/specs/systems/filecoin_vm/sysactors"
	util "github.com/filecoin-project/specs/util"
)

type ActorSubstateCID = actor.ActorSubstateCID
type ExitCode = exitcode.ExitCode
type RuntimeError = exitcode.RuntimeError
type CallerPattern = vmr.CallerPattern
type Runtime = vmr.Runtime
type InvocInput = vmr.InvocInput
type InvocOutput = vmr.InvocOutput
type MessageReceipt = vmr.MessageReceipt
type ActorStateHandle = vmr.ActorStateHandle

var EnsureErrorCode = exitcode.EnsureErrorCode
var SystemError = exitcode.SystemError

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

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
	_globalStateInit      st.StateTree
	_globalStatePending   st.StateTree
	_running              bool
	_actorAddress         addr.Address
	_actorStateAcquired   bool
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
	_valueReceived      actor.TokenAmount
	_gasRemaining       msg.GasAmount
	_numValidateCalls   int
	_output             vmr.InvocOutput
}

func VMContext_Make(
	toplevelSender addr.Address,
	toplevelBlockWinner addr.Address,
	toplevelSenderCallSeqNum actor.CallSeqNum,
	internalCallSeqNum actor.CallSeqNum,
	globalState st.StateTree,
	actorAddress addr.Address,
	valueReceived actor.TokenAmount,
	gasRemaining msg.GasAmount) *VMContext {

	return &VMContext{
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

func (rt *VMContext) _rtAllocGasCreateActor() {
	if !rt._actorAddress.Equals(addr.InitActorAddr) {
		rt.AbortAPI("Only InitActor may call rt.CreateActor_DeductGas")
	}

	rt._rtAllocGas(gascost.ExecNewActor)
}

func (rt *VMContext) CreateActor(codeID actor.CodeID, address addr.Address) {
	if !rt._actorAddress.Equals(addr.InitActorAddr) {
		rt.AbortAPI("Only InitActor may call rt.CreateActor")
	}

	rt._createActor(codeID, address)
}

func (rt *VMContext) _createActor(codeID actor.CodeID, address addr.Address) {
	// Create empty actor state.
	actorState := &actor.ActorState_I{
		CodeID_:     codeID,
		State_:      actor.ActorSubstateCID(ipld.EmptyCID()),
		Balance_:    actor.TokenAmount(0),
		CallSeqNum_: 0,
	}

	// Put it in the state tree.
	actorStateCID := actor.ActorSystemStateCID(rt.IpldPut(actorState))
	rt._updateActorSystemStateInternal(address, actorStateCID)

	rt._rtAllocGasCreateActor()
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

func (rt *VMContext) ValidateImmediateCallerIs(callerExpected addr.Address) {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeSingleton(callerExpected))
}

func (rt *VMContext) ValidateImmediateCallerAcceptAnyOfType(type_ actor.BuiltinActorID) {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeAcceptAnyOfType(rt, type_))
}

func (rt *VMContext) ValidateImmediateCallerAcceptAny() {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeAcceptAny())
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
	panic(exitcode.RuntimeError_Make(exitCode, errMsg))
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

func (rt *VMContext) _transferFunds(from addr.Address, to addr.Address, amount actor.TokenAmount) error {
	rt._checkRunning()
	rt._checkActorStateNotAcquired()

	newGlobalStatePending, err := rt._globalStatePending.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		return err
	}

	rt._globalStatePending = newGlobalStatePending
	return nil
}

func (rt *VMContext) _getActorCodeID(actorAddr addr.Address) (ret actor.CodeID, ok bool) {
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
	method actor.MethodNum,
	params actor.MethodParams) (
	ret InvocOutput, exitCode exitcode.ExitCode, internalCallSeqNumFinal actor.CallSeqNum) {

	if method == actor.MethodSend {
		ret = vmr.InvocOutput_Make(nil)
		return
	}

	rt._running = true
	ret, exitCode = _catchRuntimeErrors(func() InvocOutput {
		methodOutput := actorCode.InvokeMethod(rt, method, params)
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

	rtOuter._rtAllocGas(gascost.InvokeMethod(input.Value()))

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

	return vmr.MessageReceipt_Make(invocOutput, exitCode, gasUsed)
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

	if !targetRaw.IsKeyType() {
		// Don't implicitly create an account actor for an address without an associated key.
		rt._throwError(exitcode.SystemError(exitcode.ActorNotFound))
	}

	// Allocate an ID address from the init actor and map the pubkey To address to it.
	newIdAddr := initSubState.MapAddressToNewID(targetRaw)
	rt._saveInitActorState(initSubState)

	// Create new account actor (charges gas).
	rt._createActor(actor.CodeID(actor.CodeID_Make_Builtin(actor.BuiltinActorID_Account)), newIdAddr)

	// Initialize account actor substate with it's pubkey address.
	substate := &sysactors.AccountActorState_I{
		Address_: targetRaw,
	}
	rt._saveAccountActorState(newIdAddr, substate)
	act, _ = rt._globalStatePending.GetActor(newIdAddr)
	return act, newIdAddr
}

func (rt *VMContext) _loadInitActorState() sysactors.InitActorState {
	initState, ok := rt._globalStatePending.GetActor(addr.InitActorAddr)
	util.Assert(ok)
	initStateBytes := rt.IpldGet(ipld.CID(initState.State()))
	initSubState, err := sysactors.Deserialize_InitActorState(util.Serialization(initStateBytes.As_Bytes()))
	if err != nil {
		rt.AbortAPI("State deserialization error")
	}
	return initSubState
}

func (rt *VMContext) _saveInitActorState(state sysactors.InitActorState) {
	rt._updateActorSubstateInternal(addr.InitActorAddr, actor.ActorSubstateCID(rt.IpldPut(state.Impl())))
}

func (rt *VMContext) _saveAccountActorState(address addr.Address, state sysactors.AccountActorState) {
	rt._updateActorSubstateInternal(address, actor.ActorSubstateCID(rt.IpldPut(state.Impl())))
}

func (rt *VMContext) _sendInternalOutputs(input InvocInput, errSpec ErrorHandlingSpec) (InvocOutput, exitcode.ExitCode) {
	ret := rt._sendInternal(input, errSpec)
	return vmr.InvocOutput_Make(ret.ReturnValue()), ret.ExitCode()
}

func (rt *VMContext) SendPropagatingErrors(input InvocInput) InvocOutput {
	ret, _ := rt._sendInternalOutputs(input, PropagateErrors)
	return ret
}

func (rt *VMContext) SendCatchingErrors(input InvocInput) (InvocOutput, exitcode.ExitCode) {
	rt.ValidateImmediateCallerIs(addr.CronActorAddr)
	return rt._sendInternalOutputs(input, CatchErrors)
}

func (rt *VMContext) CurrentBalance() actor.TokenAmount {
	IMPL_FINISH()
	panic("")
}

func (rt *VMContext) ValueReceived() actor.TokenAmount {
	return rt._valueReceived
}

func (rt *VMContext) Randomness(e block.ChainEpoch, offset uint64) util.Randomness {
	// TODO: validate CurrEpoch() - K <= e <= CurrEpoch()?
	// TODO: finish
	TODO()
	panic("")
}

func (rt *VMContext) NewActorAddress() addr.Address {
	buf := new(bytes.Buffer)
	_, err := buf.Write(addr.Serialize_Address(rt._immediateCaller))
	util.Assert(err != nil)
	err = binary.Write(buf, binary.BigEndian, rt._toplevelSenderCallSeqNum)
	util.Assert(err != nil)
	err = binary.Write(buf, binary.BigEndian, rt._internalCallSeqNum)
	util.Assert(err != nil)

	hash := addr.ActorExecHash(buf.Bytes())
	return addr.Address_Make_ActorExec(addr.Address_NetworkID_Testnet, hash)
}

func (rt *VMContext) IpldPut(x ipld.Object) ipld.CID {
	var serializedSize int
	IMPL_FINISH()
	panic("") // compute serializedSize

	rt._rtAllocGas(gascost.IpldPut(serializedSize))

	IMPL_FINISH()
	panic("") // write to IPLD store
}

func (rt *VMContext) IpldGet(c ipld.CID) vmr.Runtime_IpldGet_FunRet {
	IMPL_FINISH()
	panic("") // get from IPLD store

	var serializedSize int
	IMPL_FINISH()
	panic("") // retrieve serializedSize

	rt._rtAllocGas(gascost.IpldGet(serializedSize))

	IMPL_FINISH()
	panic("") // return item
}

func (rt *VMContext) CurrEpoch() block.ChainEpoch {
	IMPL_FINISH()
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

func CallerPattern_MakeSingleton(x addr.Address) CallerPattern {
	return vmr.CallerPattern{
		Matches: func(y addr.Address) bool { return x == y },
	}
}

func CallerPattern_MakeAcceptAny() CallerPattern {
	return vmr.CallerPattern{
		Matches: func(addr.Address) bool { return true },
	}
}

func CallerPattern_MakeAcceptAnyOfType(rt *VMContext, type_ actor.BuiltinActorID) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool {
			codeID, ok := rt._getActorCodeID(y)
			if !ok {
				panic("Internal runtime error: actor not found")
			}
			Assert(codeID != nil)
			return (codeID.IsBuiltin() && (codeID.As_Builtin() == type_))
		},
	}
}
