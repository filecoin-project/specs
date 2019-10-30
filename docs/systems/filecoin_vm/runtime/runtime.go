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
type RuntimeError = exitcode.RuntimeError

var EnsureErrorCode = exitcode.EnsureErrorCode
var SystemError = exitcode.SystemError
var TODO = util.TODO

func ActorSubstateCID_Equals(x, y ActorSubstateCID) bool {
	panic("TODO")
}

// ActorCode is the interface that all actor code types should satisfy.
// It is merely a method dispatch interface.
type ActorCode interface {
	InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput
}

type ActorStateHandle struct {
	_initValue *ActorSubstateCID
	_rt        *VMContext
}

func (h *ActorStateHandle) UpdateRelease(newStateCID ActorSubstateCID) {
	h._rt._updateReleaseActorState(newStateCID)
}

func (h *ActorStateHandle) Release(checkStateCID ActorSubstateCID) {
	h._rt._releaseActorState(checkStateCID)
}

func (h *ActorStateHandle) Take() ActorSubstateCID {
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
	_globalStateInit        st.StateTree
	_globalStatePending     st.StateTree
	_running                bool
	_actorAddress           addr.Address
	_actorStateAcquired     bool
	_actorStateAcquiredInit actor.ActorSubstateCID

	_valueSupplied    actor.TokenAmount
	_gasRemaining     msg.GasAmount
	_numValidateCalls int
	_output           msg.InvocOutput
}

func VMContext_Make(
	globalState st.StateTree,
	actorAddress addr.Address,
	valueSupplied actor.TokenAmount,
	gasRemaining msg.GasAmount) *VMContext {

	actorStateInit := globalState.GetActor(actorAddress).State()

	return &VMContext{
		_globalStateInit:        globalState,
		_globalStatePending:     globalState,
		_running:                false,
		_actorAddress:           actorAddress,
		_actorStateAcquired:     false,
		_actorStateAcquiredInit: actorStateInit.State(),

		_valueSupplied:    valueSupplied,
		_gasRemaining:     gasRemaining,
		_numValidateCalls: 0,
		_output:           nil,
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

func (rt *VMContext) CreateActor(stateCID actor.StateCID, address addr.Address, constructorParams actor.MethodParams) Runtime_CreateActor_FunRet {
	rt.ValidateCallerIs(addr.InitActorAddr)

	// TODO: set actor state in global states
	// rt._globalStatePending.ActorStates()[address] = stateCID

	// TODO: call constructor
	// TODO: can constructors fail?
	// TODO: maybe do this directly form InitActor, and only do the StateTree.ActorStates() updating here?
	rt.Send(&msg.InvocInput_I{
		To_:     address,
		Method_: actor.MethodConstructor,
		Params_: constructorParams,
		Value_:  rt.ValueSupplied(),
	})

	// TODO: finish
	panic("TODO")
}

func (rt *VMContext) _updateReleaseActorState(newStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()
	newGlobalStatePending, err := rt._globalStatePending.Impl().WithActorState(rt._actorAddress, newStateCID)
	if err != nil {
		panic("Error in runtime implementation: failed to update actor state")
	}
	rt._globalStatePending = newGlobalStatePending
	rt._actorStateAcquired = false
}

func (rt *VMContext) _releaseActorState(checkStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()

	prevState := rt._globalStatePending.GetActor(rt._actorAddress).State()
	prevStateCID := prevState.State()
	if !ActorSubstateCID_Equals(prevStateCID, checkStateCID) {
		rt.Abort("State CID differs upon release call")
	}

	rt._actorStateAcquired = false
}

func (rt *VMContext) Assert(cond bool) Runtime_Assert_FunRet {
	if !cond {
		rt.Abort("Runtime check failed")
	}
	return &Runtime_Assert_FunRet_I{}
}

func (rt *VMContext) _checkActorStateAcquired() {
	if !rt._running {
		panic("Error in runtime implementation: actor interface invoked without running actor")
	}

	if !rt._actorStateAcquired {
		rt.Abort("Actor state not acquired")
	}
}

func (rt *VMContext) Abort(errMsg string) Runtime_Abort_FunRet {
	rt._throwErrorFull(exitcode.SystemError(exitcode.MethodAbort), errMsg)
	return &Runtime_Abort_FunRet_I{}
}

func (rt *VMContext) Caller() addr.Address {
	panic("TODO")
}

func (rt *VMContext) ValidateCallerMatches(callerExpectedPattern CallerPattern) Runtime_ValidateCallerMatches_FunRet {
	rt._checkRunning()
	rt._checkNumValidateCalls(0)
	caller := rt.Caller()
	if !callerExpectedPattern.Matches(caller) {
		rt.Abort("Method invoked by incorrect caller")
	}
	rt._numValidateCalls += 1
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

func (rt *VMContext) ValidateCallerIs(callerExpected addr.Address) Runtime_ValidateCallerIs_FunRet {
	rt.ValidateCallerMatches(CallerPattern_MakeSingleton(callerExpected))
	return &Runtime_ValidateCallerIs_FunRet_I{}
}

func (rt *VMContext) _checkNumValidateCalls(x int) {
	if rt._numValidateCalls != x {
		rt.Abort("Method must validate caller identity exactly once")
	}
}

func (rt *VMContext) _checkRunning() {
	if !rt._running {
		panic("Internal runtime error: actor API called with no actor code running")
	}
}
func (rt *VMContext) SuccessReturn() InvocOutput {
	return msg.InvocOutput_Make(exitcode.OK(), nil)
}

func (rt *VMContext) ValueReturn(value util.Bytes) InvocOutput {
	return msg.InvocOutput_Make(exitcode.OK(), value)
}

func (rt *VMContext) ErrorReturn(exitCode ExitCode) InvocOutput {
	exitCode = exitcode.EnsureErrorCode(exitCode)
	return msg.InvocOutput_Make(exitCode, nil)
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

func (rt *VMContext) _checkStateLock(expected bool) {
	if rt._actorStateAcquired != expected {
		rt._apiError("State update and message send blocks must be disjoint")
	}
}

func (rt *VMContext) _checkGasRemaining() {
	if rt._gasRemaining.LessThan(msg.GasAmount_Zero()) {
		rt._throwError(exitcode.SystemError(exitcode.OutOfGas))
	}
}

func (rt *VMContext) _deductGasRemaining(x msg.GasAmount) {
	// TODO: check x >= 0
	rt._checkGasRemaining()
	rt._gasRemaining = rt._gasRemaining.Subtract(x)
	rt._checkGasRemaining()
}

func (rt *VMContext) _refundGasRemaining(x msg.GasAmount) {
	// TODO: check x >= 0
	rt._checkGasRemaining()
	rt._gasRemaining = rt._gasRemaining.Add(x)
	rt._checkGasRemaining()
}

func (rt *VMContext) _transferFunds(from addr.Address, to addr.Address, amount actor.TokenAmount) error {
	rt._checkRunning()
	rt._checkStateLock(false)

	newGlobalStatePending, err := rt._globalStatePending.Impl().WithFundsTransfer(from, to, amount)
	if err != nil {
		return err
	}

	rt._globalStatePending = newGlobalStatePending
	return nil
}

// TODO: This function should be private (not intended to be exposed to actors).
// (merging runtime and interpreter packages should solve this)
func (rt *VMContext) SendToplevelFromInterpreter(input InvocInput, ignoreErrors bool) (
	msg.MessageReceipt, st.StateTree) {

	rt._running = true
	ret := rt._sendInternal(input, ignoreErrors)
	rt._running = false
	return ret, rt._globalStatePending
}

func _catchRuntimeErrors(f func() msg.InvocOutput) (output msg.InvocOutput) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				output = msg.InvocOutput_Make(EnsureErrorCode(r.(*RuntimeError).ExitCode), nil)
			default:
				output = msg.InvocOutput_Make(SystemError(exitcode.MethodPanic), nil)
			}
		}
	}()

	output = f()
	return
}

func _invokeMethodInternal(
	rt *VMContext,
	actorCode ActorCode,
	method actor.MethodNum,
	params actor.MethodParams) (ret InvocOutput, gasUsed msg.GasAmount) {

	if method == actor.MethodSend {
		ret = msg.InvocOutput_Make(exitcode.OK(), nil)
		gasUsed = msg.GasAmount_Zero() // TODO: verify
		return
	}

	rt._running = true
	ret = _catchRuntimeErrors(func() InvocOutput {
		methodOutput := actorCode.InvokeMethod(rt, method, params)
		rt._checkStateLock(false)
		rt._checkNumValidateCalls(1)
		return methodOutput
	})
	rt._running = false

	// TODO: Update gasUsed
	TODO()

	return
}

func (rtOuter *VMContext) _sendInternal(input InvocInput, ignoreErrors bool) msg.MessageReceipt {
	rtOuter._checkRunning()
	rtOuter._checkStateLock(false)

	toActor := rtOuter._globalStatePending.GetActor(input.To()).State()

	toActorCode, err := loadActorCode(toActor.CodeCID())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.ActorCodeNotFound))
	}

	var toActorMethodGasBound msg.GasAmount
	TODO() // TODO: obtain from actor registry
	rtOuter._deductGasRemaining(toActorMethodGasBound)
	// TODO: gasUsed may be larger than toActorMethodGasBound if toActor itself makes sub-calls.
	// To prevent this, we would need to calculate the gas bounds recursively.

	err = rtOuter._transferFunds(rtOuter._actorAddress, input.To(), input.Value())
	if err != nil {
		rtOuter._throwError(exitcode.SystemError(exitcode.InsufficientFunds))
	}

	rtInner := VMContext_Make(
		rtOuter._globalStatePending,
		input.To(),
		input.Value(),
		rtOuter._gasRemaining,
	)

	invocOutput, gasUsed := _invokeMethodInternal(
		rtInner,
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
		rtOuter._globalStatePending = rtInner._globalStatePending
	}

	return msg.MessageReceipt_Make(invocOutput, gasUsed)
}

func (rt *VMContext) Send(input InvocInput) msg.MessageReceipt {
	return rt._sendInternal(input, false)
}

func (rt *VMContext) SendAllowingErrors(input InvocInput) msg.MessageReceipt {
	return rt._sendInternal(input, true)
}

func (rt *VMContext) ValueSupplied() actor.TokenAmount {
	return rt._valueSupplied
}

func (rt *VMContext) Randomness(e block.ChainEpoch, offset uint64) Randomness {
	// TODO: validate CurrEpoch() - K <= e <= CurrEpoch()?
	// TODO: finish
	panic("TODO")
}

func (rt *VMContext) IpldPut(x ipld.Object) ipld.CID {
	panic("TODO")
}

func (rt *VMContext) IpldGet(c ipld.CID) Runtime_IpldGet_FunRet {
	panic("TODO")
}

func (rt *VMContext) CurrEpoch() block.ChainEpoch {
	panic("TODO")
}

func (rt *VMContext) AcquireState() ActorStateHandle {
	panic("TODO")
}
