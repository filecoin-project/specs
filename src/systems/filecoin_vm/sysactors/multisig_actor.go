package sysactors

import (
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
	vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
	util "github.com/filecoin-project/specs/util"
)

var IMPL_FINISH = util.IMPL_FINISH

////////////////////////////////////////////////////////////////////////////////
// Actor methods
////////////////////////////////////////////////////////////////////////////////

func (a *MultiSigActorCode_I) EndorseMessage(rt vmr.Runtime, message msg.UnsignedMessageBase) {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	callerAddr := rt.ImmediateCaller()

	h, st := a.State(rt)

	authorized := false
	for _, authorizedPartyAddr := range st.AuthorizedParties() {
		if authorizedPartyAddr.Equals(callerAddr) {
			authorized = true
			break
		}
	}

	if !authorized {
		rt.AbortArgMsg("Caller not among authorized parties")
	}

	if _, found := st.PendingMessages()[message]; !found {
		st.PendingMessages()[message] = EndorsementSetHAMT_Empty()
		st.PendingMessagesValue()[message] = actor.TokenAmount(0)
	}
	st.PendingMessages()[message][callerAddr] = true
	st.PendingMessagesValue()[message] += rt.ValueReceived()

	totalValue := st.PendingMessagesValue()[message]
	isFinalized := (len(st.PendingMessages()[message]) == st.NumEndorsementsThreshold())
	if isFinalized {
		// A sufficient number of endorsements have arrived: relay the message and delete from pending queue.
		delete(st.PendingMessages(), message)
		delete(st.PendingMessagesValue(), message)
	}

	UpdateRelease_MultiSig(rt, h, st)

	if isFinalized {
		rt.SendRelay(message, totalValue)
	}
}

func (a *MultiSigActorCode_I) Constructor(
	rt vmr.Runtime, authorizedParties []addr.Address, numEndorsementsThreshold int) {

	rt.ValidateImmediateCallerIs(addr.InitActorAddr)
	h := rt.AcquireState()

	st := &MultiSigActorState_I{
		AuthorizedParties_:        authorizedParties,
		NumEndorsementsThreshold_: numEndorsementsThreshold,
		PendingMessages_:          MultiSigPendingMessageHAMT_Empty(),
	}

	UpdateRelease_MultiSig(rt, h, st)
}

////////////////////////////////////////////////////////////////////////////////
// Method utility functions
////////////////////////////////////////////////////////////////////////////////

func EndorsementSetHAMT_Empty() EndorsementSetHAMT {
	IMPL_FINISH()
	panic("")
}

func MultiSigPendingMessageHAMT_Empty() MultiSigPendingMessageHAMT {
	IMPL_FINISH()
	panic("")
}

////////////////////////////////////////////////////////////////////////////////
// Dispatch table
////////////////////////////////////////////////////////////////////////////////

func (a *MultiSigActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	IMPL_FINISH()
	panic("")
}
