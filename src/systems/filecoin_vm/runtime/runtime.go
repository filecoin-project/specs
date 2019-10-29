package runtime

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"

type ActorSubstateCID = actor.ActorSubstateCID
type InvocInput = msg.InvocInput
type InvocOutput = msg.InvocOutput
type ExitCode = exitcode.ExitCode

var TODO = util.TODO

func ActorSubstateCID_Equals(x, y ActorSubstateCID) bool {
	panic("TODO")
}

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams)
}

const (
	Reserved_NoopMethod        actor.MethodNum = 0
	Reserved_CronMethod        actor.MethodNum = 1
	Reserved_ConstructorMethod actor.MethodNum = 2
)

func Runtime_Make(
	globalState st.StateTree,
	actorAddress addr.Address,
	valueSupplied actor.TokenAmount,
	gasRemaining msg.GasAmount) Runtime {

	actorStateInit := globalState.GetActor(actorAddress).State()

	return &Runtime_I{
		_globalStateInit_:        globalState,
		_globalStatePending_:     globalState,
		_running_:                false,
		_actorAddress_:           actorAddress,
		_actorStateAcquired_:     false,
		_actorStateAcquiredInit_: actorStateInit.State(),

		_valueSupplied_:    valueSupplied,
		_gasRemaining_:     gasRemaining,
		_numValidateCalls_: 0,
		_numReturnCalls_:   0,
		_messageQueue_:     []MessageQueueItem{},
		_output_:           nil,
	}
}

func _generateActorAddress(creator addr.Address, nonce actor.CallSeqNum) addr.Address {
	// _generateActorAddress computes the address of the contract,
	// based on the creator (invoking address) and nonce given.
	// TODO: why is this needed? -- InitActor
	// TODO: this has to be the origin call. and it's broken: could yield the same address
	//       need a truly unique way to assign an address.
	panic("TODO")
}

func (rt *Runtime_I) CreateActor(codeCID actor.CodeCID, constructorParams actor.MethodParams) Runtime_CreateActor_FunRet {
	rt.ValidateCallerIs(InitActorAddr)
	// TODO: _generateActorAddress
	// TODO: finish
	panic("TODO")
}

func (h *ActorStateHandle_I) UpdateRelease(newStateCID ActorSubstateCID) {
	h._rt().Impl()._updateReleaseActorState(newStateCID)
}

func (h *ActorStateHandle_I) Release(checkStateCID ActorSubstateCID) {
	h._rt().Impl()._releaseActorState(checkStateCID)
}

func (h *ActorStateHandle_I) Take() ActorSubstateCID {
	if h._initValue().Which() == ActorStateHandle__initValue_Case_None {
		h._rt().Impl()._apiError("Must call Take() only once on actor substate object")
	}
	ret := ActorSubstateCID(h._initValue().As_Some())
	h._initValue_ = ActorStateHandle__initValue_Make_None(&ActorStateHandle__initValue_None_I{})
	return ret
}

func (rt *Runtime_I) _updateReleaseActorState(newStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()
	newGlobalStatePending, err := rt._globalStatePending().Impl().WithActorState(rt._actorAddress(), newStateCID)
	if err != nil {
		panic("Error in runtime implementation: failed to update actor state")
	}
	rt._globalStatePending_ = newGlobalStatePending
	rt._actorStateAcquired_ = false
}

func (rt *Runtime_I) _releaseActorState(checkStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()

	prevState := rt._globalStatePending().GetActor(rt._actorAddress()).State()
	prevStateCID := prevState.State()
	if !ActorSubstateCID_Equals(prevStateCID, checkStateCID) {
		rt.Abort("State CID differs upon release call")
	}

	rt._actorStateAcquired_ = false
}

func (rt *Runtime_I) Check(cond bool) Runtime_Check_FunRet {
	if !cond {
		rt.Abort("Runtime check failed")
	}
	return &Runtime_Check_FunRet_I{}
}

func (rt *Runtime_I) _checkActorStateAcquired() {
	if !rt._running() {
		panic("Error in runtime implementation: actor interface invoked without running actor")
	}

	if !rt._actorStateAcquired() {
		rt.Abort("Actor state not acquired")
	}
}

func (rt *Runtime_I) Abort(errMsg string) Runtime_Abort_FunRet {
	rt._throwErrorFull(exitcode.SystemError(exitcode.MethodAbort), errMsg)
	return &Runtime_Abort_FunRet_I{}
}

func (rt *Runtime_I) Caller() addr.Address {
	panic("TODO")
}

func (rt *Runtime_I) ValidateCallerMatches(callerExpectedPattern CallerPattern) Runtime_ValidateCallerMatches_FunRet {
	rt._checkRunning()
	rt._checkNumValidateCalls(0)
	caller := rt.Caller()
	if !callerExpectedPattern.Matches(caller) {
		rt.Abort("Method invoked by incorrect caller")
	}
	rt._numValidateCalls_ += 1
	return &Runtime_ValidateCallerMatches_FunRet_I{}
}

type CallerPattern struct {
	Matches func(addr.Address) bool
}

func CallerPattern_MakeSingleton(x addr.Address) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool {
			return x == y
		},
	}
}

func (rt *Runtime_I) ValidateCallerIs(callerExpected addr.Address) Runtime_ValidateCallerIs_FunRet {
	rt.ValidateCallerMatches(CallerPattern_MakeSingleton(callerExpected))
	return &Runtime_ValidateCallerIs_FunRet_I{}
}

func (rt *Runtime_I) _checkNumValidateCalls(x int) {
	if rt._numValidateCalls() != x {
		rt.Abort("Method must validate caller identity exactly once")
	}
}

func (rt *Runtime_I) _checkNumReturnCalls(x int) {
	if rt._numReturnCalls() != x {
		rt.Abort("Method must call return exactly once")
	}
}

func (rt *Runtime_I) _checkRunning() {
	if !rt._running() {
		panic("Internal runtime error: actor API called with no actor code running")
	}
}

func (rt *Runtime_I) _returnInternal(output InvocOutput) {
	rt._checkRunning()
	rt._checkStateLock(false)
	rt._checkNumReturnCalls(0)
	rt._output_ = output
	rt._numReturnCalls_ += 1
}

func (rt *Runtime_I) ReturnSuccess() Runtime_ReturnSuccess_FunRet {
	rt._returnInternal(msg.InvocOutput_Make(exitcode.OK(), nil))
	return &Runtime_ReturnSuccess_FunRet_I{}
}

func (rt *Runtime_I) ReturnValue(value util.Bytes) Runtime_ReturnValue_FunRet {
	rt._returnInternal(msg.InvocOutput_Make(exitcode.OK(), value))
	return &Runtime_ReturnValue_FunRet_I{}
}

func (rt *Runtime_I) ReturnError(exitCode ExitCode) Runtime_ReturnError_FunRet {
	exitCode = exitcode.EnsureErrorCode(exitCode)
	rt._returnInternal(msg.InvocOutput_Make(exitCode, nil))
	return &Runtime_ReturnError_FunRet_I{}
}

func (rt *Runtime_I) _throwError(exitCode ExitCode) {
	rt._throwErrorFull(exitCode, "")
}

func (rt *Runtime_I) _throwErrorFull(exitCode ExitCode, errMsg string) {
	panic(exitcode.RuntimeError_Make(exitCode, errMsg))
}

func (rt *Runtime_I) _apiError(errMsg string) {
	rt._throwErrorFull(exitcode.SystemError(exitcode.RuntimeAPIError), errMsg)
}

func (rt *Runtime_I) _checkStateLock(expected bool) {
	if rt._actorStateAcquired() != expected {
		rt._apiError("State update and message send blocks must be disjoint")
	}
}

func (rt *Runtime_I) _checkGasRemaining() {
	if rt._gasRemaining().LessThan(msg.GasAmount_Zero()) {
		rt._throwError(exitcode.SystemError(exitcode.OutOfGas))
	}
}

func (rt *Runtime_I) _deductGasRemaining(x msg.GasAmount) {
	// TODO: check x >= 0
	rt._checkGasRemaining()
	rt._gasRemaining_ = rt._gasRemaining().Subtract(x)
	rt._checkGasRemaining()
}

func (rt *Runtime_I) _refundGasRemaining(x msg.GasAmount) {
	// TODO: check x >= 0
	rt._checkGasRemaining()
	rt._gasRemaining_ = rt._gasRemaining().Add(x)
	rt._checkGasRemaining()
}

func (rt *Runtime_I) _transferFunds(from addr.Address, to addr.Address, amount actor.TokenAmount) error {
	rt._checkRunning()
	rt._checkStateLock(false)

	newGlobalStatePending, err := rt._globalStatePending().Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		return err
	}

	rt._globalStatePending_ = newGlobalStatePending
	return nil
}

// TODO: This function should be private (not intended to be exposed to actors).
// (merging runtime and interpreter packages should solve this)
func (rt *Runtime_I) SendToplevelFromInterpreter(input InvocInput, ignoreErrors bool) (
	msg.MessageReceipt, st.StateTree) {

	rt._running_ = true
	ret := rt._sendInternal(input, ignoreErrors)
	rt._running_ = false
	return ret, rt._globalStatePending()
}

func _invokeMethodInternal(
	rt *Runtime_I,
	actorCode ActorCode,
	method actor.MethodNum,
	params actor.MethodParams) (ret InvocOutput, gasUsed msg.GasAmount) {

	if method == Reserved_NoopMethod {
		ret = msg.InvocOutput_Make(exitcode.OK(), nil)
		gasUsed = msg.GasAmount_Zero() // TODO: verify
		return
	}

	rt._running_ = true
	catchRetCode := exitcode.CatchRuntimeErrors(func() {
		actorCode.InvokeMethod(rt, method, params)
		rt._checkStateLock(false)
		rt._checkNumValidateCalls(1)
		rt._checkNumReturnCalls(1)
		for _, item := range rt._messageQueue() {
			rt._sendInternal(item.input, item.ignoreErrors)
		}
		rt._messageQueue_ = []MessageQueueItem{}
	})
	rt._running_ = false

	if catchRetCode.IsSuccess() {
		if rt._output() == nil {
			panic("Internal error: return call recorded but no output object returned")
		}
		ret = rt._output()
	} else {
		ret = msg.InvocOutput_Make(catchRetCode, nil)
	}

	// TODO: Update gasUsed
	TODO()

	return
}

func (rtOuter *Runtime_I) _sendInternal(input InvocInput, ignoreErrors bool) msg.MessageReceipt {
	rtOuter._checkRunning()
	rtOuter._checkStateLock(false)

	toActor := rtOuter._globalStatePending().GetActor(input.To()).State()

	toActorCode, err := loadActorCode(toActor.CodeCID())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.ActorCodeNotFound))
	}

	var toActorMethodGasBound msg.GasAmount
	TODO() // TODO: obtain from actor registry
	rtOuter._deductGasRemaining(toActorMethodGasBound)
	// TODO: gasUsed may be larger than toActorMethodGasBound if toActor itself makes sub-calls.
	// To prevent this, we would need to calculate the gas bounds recursively.

	err = rtOuter._transferFunds(rtOuter._actorAddress(), input.To(), input.Value())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.InsufficientFunds))
	}

	rtInner := Runtime_Make(
		rtOuter._globalStatePending(),
		input.To(),
		input.Value(),
		rtOuter._gasRemaining(),
	)

	invocOutput, gasUsed := _invokeMethodInternal(
		rtInner.Impl(),
		toActorCode,
		input.Method(),
		input.Params(),
	)

	rtOuter._refundGasRemaining(toActorMethodGasBound)
	rtOuter._deductGasRemaining(gasUsed)

	if !ignoreErrors && invocOutput.ExitCode().IsError() {
		rtOuter._throwError(exitcode.SystemError(exitcode.MethodSubcallError))
	}

	if invocOutput.ExitCode().AllowsStateUpdate() {
		rtOuter._globalStatePending_ = rtInner._globalStatePending()
	}

	return msg.MessageReceipt_Make(invocOutput, gasUsed)
}

func (rt *Runtime_I) Send(input InvocInput) msg.MessageReceipt {
	return rt._sendInternal(input, false)
}

func (rt *Runtime_I) SendAllowingErrors(input InvocInput) msg.MessageReceipt {
	return rt._sendInternal(input, true)
}

type MessageQueue []MessageQueueItem
type MessageQueueItem struct {
	input        InvocInput
	ignoreErrors bool
}

func MessageQueueItem_Make(input InvocInput, ignoreErrors bool) MessageQueueItem {
	return MessageQueueItem{
		input:        input,
		ignoreErrors: ignoreErrors,
	}
}

func (rt *Runtime_I) _deferredSendInternal(input InvocInput, ignoreErrors bool) {
	rt._checkRunning()
	rt._checkStateLock(false)
	rt._messageQueue_ = append(rt._messageQueue(), MessageQueueItem_Make(input, ignoreErrors))
}

func (rt *Runtime_I) DeferredSend(input InvocInput) Runtime_DeferredSend_FunRet {
	rt._deferredSendInternal(input, false)
	return &Runtime_DeferredSend_FunRet_I{}
}

func (rt *Runtime_I) DeferredSendAllowingErrors(input InvocInput) Runtime_DeferredSendAllowingErrors_FunRet {
	rt._deferredSendInternal(input, true)
	return &Runtime_DeferredSendAllowingErrors_FunRet_I{}
}

func (rt *Runtime_I) ValueSupplied() actor.TokenAmount {
	return rt._valueSupplied()
}

func (rt *Runtime_I) Randomness(offset uint64) Randomness {
	panic("TODO")
}

func (rt *Runtime_I) IpldPut(x ipld.Object) ipld.CID {
	panic("TODO")
}

func (rt *Runtime_I) IpldGet(c ipld.CID) Runtime_IpldGet_FunRet {
	panic("TODO")
}

func (rt *Runtime_I) CurrEpoch() block.ChainEpoch {
	panic("TODO")
}

func (rt *Runtime_I) AcquireState() ActorStateHandle {
	panic("TODO")
}
