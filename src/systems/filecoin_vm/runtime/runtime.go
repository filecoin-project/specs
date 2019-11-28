package runtime

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import filcrypto "github.com/filecoin-project/specs/algorithms/crypto"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import exitcode "github.com/filecoin-project/specs/systems/filecoin_vm/runtime/exitcode"
import util "github.com/filecoin-project/specs/util"

type ActorSubstateCID = actor.ActorSubstateCID
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
	h._rt._updateReleaseActorSubstate(newStateCID)
}

func (h *ActorStateHandle) Release(checkStateCID ActorSubstateCID) {
	h._rt._releaseActorSubstate(checkStateCID)
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
	_output             InvocOutput
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

	actorStateInit := globalState.GetActorState(actorAddress)

	return &VMContext{
		_globalStateInit:        globalState,
		_globalStatePending:     globalState,
		_running:                false,
		_actorAddress:           actorAddress,
		_actorStateAcquired:     false,
		_actorStateAcquiredInit: actorStateInit.State(),

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

func (rt *VMContext) CreateActor(
	stateCID actor.ActorSystemStateCID,
	address addr.Address,
	initBalance actor.TokenAmount,
	constructorParams actor.MethodParams) Runtime_CreateActor_FunRet {

	if !rt._actorAddress.Equals(addr.InitActorAddr) {
		rt.Abort("Only InitActor may call rt.CreateActor")
	}

	rt._updateActorSystemStateInternal(address, stateCID)

	rt.SendPropagatingErrors(&InvocInput_I{
		To_:     address,
		Method_: actor.MethodConstructor,
		Params_: constructorParams,
		Value_:  initBalance,
	})

	return &Runtime_CreateActor_FunRet_I{}
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
	rt._actorStateAcquired = false
}

func (rt *VMContext) _releaseActorSubstate(checkStateCID ActorSubstateCID) {
	rt._checkRunning()
	rt._checkActorStateAcquired()

	prevState := rt._globalStatePending.GetActorState(rt._actorAddress)
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

func (rt *VMContext) ImmediateCaller() addr.Address {
	return rt._immediateCaller
}

func (rt *VMContext) ToplevelBlockWinner() addr.Address {
	return rt._toplevelBlockWinner
}

func (rt *VMContext) ValidateImmediateCallerMatches(
	callerExpectedPattern CallerPattern) Runtime_ValidateImmediateCallerMatches_FunRet {

	rt._checkRunning()
	rt._checkNumValidateCalls(0)
	caller := rt.ImmediateCaller()
	if !callerExpectedPattern.Matches(caller) {
		rt.Abort("Method invoked by incorrect caller")
	}
	rt._numValidateCalls += 1
	return &Runtime_ValidateImmediateCallerMatches_FunRet_I{}
}

type CallerPattern struct {
	Matches func(addr.Address) bool
}

func CallerPattern_MakeSingleton(x addr.Address) CallerPattern {
	return CallerPattern{
		Matches: func(y addr.Address) bool { return x == y },
	}
}

func CallerPattern_MakeAcceptAny() CallerPattern {
	return CallerPattern{
		Matches: func(addr.Address) bool { return true },
	}
}

func (rt *VMContext) ValidateImmediateCallerIs(callerExpected addr.Address) Runtime_ValidateImmediateCallerIs_FunRet {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeSingleton(callerExpected))
	return &Runtime_ValidateImmediateCallerIs_FunRet_I{}
}

func (rt *VMContext) ValidateImmediateCallerAcceptAny() Runtime_ValidateImmediateCallerAcceptAny_FunRet {
	rt.ValidateImmediateCallerMatches(CallerPattern_MakeAcceptAny())
	return &Runtime_ValidateImmediateCallerAcceptAny_FunRet_I{}
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
	return InvocOutput_Make(exitcode.OK(), nil)
}

func (rt *VMContext) ValueReturn(value util.Bytes) InvocOutput {
	return InvocOutput_Make(exitcode.OK(), value)
}

func (rt *VMContext) ErrorReturn(exitCode ExitCode) InvocOutput {
	exitCode = exitcode.EnsureErrorCode(exitCode)
	return InvocOutput_Make(exitCode, nil)
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

type ErrorHandlingSpec int

const (
	PropagateErrors ErrorHandlingSpec = 1 + iota
	CatchErrors
)

// TODO: This function should be private (not intended to be exposed to actors).
// (merging runtime and interpreter packages should solve this)
func (rt *VMContext) SendToplevelFromInterpreter(input InvocInput) (
	MessageReceipt, st.StateTree) {

	rt._running = true
	ret := rt._sendInternal(input, CatchErrors)
	rt._running = false
	return ret, rt._globalStatePending
}

func _catchRuntimeErrors(f func() InvocOutput) (output InvocOutput) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				output = InvocOutput_Make(EnsureErrorCode(r.(*RuntimeError).ExitCode), nil)
			default:
				panic(r)
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
	params actor.MethodParams) (
	ret InvocOutput, gasUsed msg.GasAmount, internalCallSeqNumFinal actor.CallSeqNum) {

	if method == actor.MethodSend {
		ret = InvocOutput_Make(exitcode.OK(), nil)
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

	internalCallSeqNumFinal = rt._internalCallSeqNum

	return
}

func (rtOuter *VMContext) _sendInternal(input InvocInput, errSpec ErrorHandlingSpec) MessageReceipt {
	rtOuter._checkRunning()
	rtOuter._checkStateLock(false)

	toActor := rtOuter._globalStatePending.GetActorState(input.To())

	toActorCode, err := loadActorCode(toActor.CodeID())
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
		rtOuter._toplevelSender,
		rtOuter._toplevelBlockWinner,
		rtOuter._toplevelSenderCallSeqNum,
		rtOuter._internalCallSeqNum+1,
		rtOuter._globalStatePending,
		input.To(),
		input.Value(),
		rtOuter._gasRemaining,
	)

	invocOutput, gasUsed, internalCallSeqNumFinal := _invokeMethodInternal(
		rtInner,
		toActorCode,
		input.Method(),
		input.Params(),
	)

	rtOuter._internalCallSeqNum = internalCallSeqNumFinal

	rtOuter._refundGasRemaining(toActorMethodGasBound)
	rtOuter._deductGasRemaining(gasUsed)

	if errSpec == PropagateErrors && invocOutput.ExitCode().IsError() {
		rtOuter._throwError(exitcode.SystemError(exitcode.MethodSubcallError))
	}

	if invocOutput.ExitCode().AllowsStateUpdate() {
		rtOuter._globalStatePending = rtInner._globalStatePending
	}

	return MessageReceipt_Make(invocOutput, gasUsed)
}

func (rtOuter *VMContext) _sendInternalOutputOnly(input InvocInput, errSpec ErrorHandlingSpec) InvocOutput {
	ret := rtOuter._sendInternal(input, errSpec)
	return &InvocOutput_I{
		ExitCode_:    ret.ExitCode(),
		ReturnValue_: ret.ReturnValue(),
	}
}

func (rt *VMContext) SendPropagatingErrors(input InvocInput) InvocOutput {
	return rt._sendInternalOutputOnly(input, PropagateErrors)
}

func (rt *VMContext) SendCatchingErrors(input InvocInput) InvocOutput {
	return rt._sendInternalOutputOnly(input, CatchErrors)
}

func (rt *VMContext) CurrentBalance() actor.TokenAmount {
	panic("TODO")
}

func (rt *VMContext) ValueReceived() actor.TokenAmount {
	return rt._valueReceived
}

func (rt *VMContext) Randomness(e block.ChainEpoch, offset uint64) util.Randomness {
	// TODO: validate CurrEpoch() - K <= e <= CurrEpoch()?
	// TODO: finish
	panic("TODO")
}

func (rt *VMContext) NewActorAddress() addr.Address {
	seed := &ActorExecAddressSeed_I{
		creator_:            rt._immediateCaller,
		toplevelCallSeqNum_: rt._toplevelSenderCallSeqNum,
		internalCallSeqNum_: rt._internalCallSeqNum,
	}
	hash := addr.ActorExecHash(Serialize_ActorExecAddressSeed(seed))

	return addr.Address_Make_ActorExec(addr.Address_NetworkID_Testnet, hash)
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

func (rt *VMContext) CurrMethodNum() actor.MethodNum {
	panic("TODO")
}

func (rt *VMContext) VerifySignature(signerActor addr.Address, sig filcrypto.Signature, m filcrypto.Message) bool {
	st := rt._globalStatePending.Impl().GetActorState(signerActor)
	if st == nil {
		rt.Abort("VerifySignature: signer actor not found")
	}
	pk := st.GetSignaturePublicKey()
	if pk == nil {
		rt.Abort("VerifySignature: signer actor has no public key")
	}
	ret := rt.Compute(ComputeFunctionID_VerifySignature, []Any{pk, sig, m})
	return ret.(bool)
}

func (rt *VMContext) Compute(f ComputeFunctionID, args []Any) Any {
	def, found := _computeFunctionDefs[f]
	if !found {
		rt.Abort("Function definition in rt.Compute() not found")
	}
	gasCost := def.GasCostFn(args)
	rt._deductGasRemaining(gasCost)
	return def.Body(args)
}

func InvocInput_Make(to addr.Address, method actor.MethodNum, params actor.MethodParams, value actor.TokenAmount) InvocInput {
	return &InvocInput_I{
		To_:     to,
		Method_: method,
		Params_: params,
		Value_:  value,
	}
}

func InvocOutput_Make(exitCode exitcode.ExitCode, returnValue util.Bytes) InvocOutput {
	return &InvocOutput_I{
		ExitCode_:    exitCode,
		ReturnValue_: returnValue,
	}
}

func MessageReceipt_Make(output InvocOutput, gasUsed msg.GasAmount) MessageReceipt {
	return &MessageReceipt_I{
		ExitCode_:    output.ExitCode(),
		ReturnValue_: output.ReturnValue(),
		GasUsed_:     gasUsed,
	}
}

func MessageReceipt_MakeSystemError(errCode exitcode.SystemErrorCode, gasUsed msg.GasAmount) MessageReceipt {
	return MessageReceipt_Make(
		InvocOutput_Make(exitcode.SystemError(errCode), nil),
		gasUsed,
	)
}
