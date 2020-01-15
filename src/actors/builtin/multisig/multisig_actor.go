package multisig

import (
	addr "github.com/filecoin-project/go-address"
	actor "github.com/filecoin-project/specs/actors"
	abi "github.com/filecoin-project/specs/actors/abi"
	builtin "github.com/filecoin-project/specs/actors/builtin"
	vmr "github.com/filecoin-project/specs/actors/runtime"
	autil "github.com/filecoin-project/specs/actors/util"
	cid "github.com/ipfs/go-cid"
)

type InvocOutput = vmr.InvocOutput
type Runtime = vmr.Runtime

var AssertMsg = autil.AssertMsg
var IMPL_FINISH = autil.IMPL_FINISH
var IMPL_TODO = autil.IMPL_TODO
var TODO = autil.TODO

type TxnID int

type MultiSigTransaction struct {
	Proposer   addr.Address
	Expiration abi.ChainEpoch

	To     addr.Address
	Method abi.MethodNum
	Params abi.MethodParams
	Value  abi.TokenAmount
}

func (txn *MultiSigTransaction) Equals(MultiSigTransaction) bool {
	IMPL_FINISH()
	panic("")
}

type MultiSigTransactionHAMT map[TxnID]MultiSigTransaction
type MultiSigApprovalSetHAMT map[TxnID]autil.ActorIDSetHAMT

func MultiSigTransactionHAMT_Empty() MultiSigTransactionHAMT {
	IMPL_FINISH()
	panic("")
}

func MultiSigApprovalSetHAMT_Empty() MultiSigApprovalSetHAMT {
	IMPL_FINISH()
	panic("")
}

type MultiSigActorState struct {
	// Linear unlock
	InitialBalance abi.TokenAmount
	StartEpoch     abi.ChainEpoch
	UnlockDuration abi.ChainEpoch

	AuthorizedParties     autil.ActorIDSetHAMT
	NumApprovalsThreshold int
	NextTxnID             TxnID
	PendingTxns           MultiSigTransactionHAMT
	PendingApprovals      MultiSigApprovalSetHAMT
}

func (st *MultiSigActorState) AmountLocked(elapsedEpoch abi.ChainEpoch) abi.TokenAmount {
	if elapsedEpoch >= st.UnlockDuration {
		return abi.TokenAmount(0)
	}

	TODO() // BigInt
	lockedProportion := (st.UnlockDuration - elapsedEpoch) / st.UnlockDuration
	return abi.TokenAmount(uint64(st.InitialBalance) * uint64(lockedProportion))
}

// return true if MultiSig maintains required locked balance after spending the amount
func (st *MultiSigActorState) _hasAvailable(currBalance abi.TokenAmount, amountToSpend abi.TokenAmount, currEpoch abi.ChainEpoch) bool {
	if amountToSpend < 0 || currBalance < amountToSpend {
		return false
	}

	if currBalance-amountToSpend < st.AmountLocked(currEpoch-st.StartEpoch) {
		return false
	}

	return true
}

func (st *MultiSigActorState) CID() cid.Cid {
	panic("TODO")
}

type MultiSigActor struct{}

func (a *MultiSigActor) State(rt Runtime) (vmr.ActorStateHandle, MultiSigActorState) {
	h := rt.AcquireState()
	stateCID := cid.Cid(h.Take())
	var state MultiSigActorState
	if !rt.IpldGet(stateCID, &state) {
		rt.AbortAPI("state not found")
	}
	return h, state
}

func (a *MultiSigActor) Propose(rt vmr.Runtime, txn MultiSigTransaction) TxnID {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	callerAddr := rt.ImmediateCaller()
	a._rtValidateAuthorizedPartyOrAbort(rt, callerAddr)

	h, st := a.State(rt)
	txnID := st.NextTxnID
	st.NextTxnID += 1
	st.PendingTxns[txnID] = txn
	st.PendingApprovals[txnID] = autil.ActorIDSetHAMT_Empty()
	UpdateRelease_MultiSig(rt, h, st)

	// Proposal implicitly includes approval of a transaction.
	a._rtApproveTransactionOrAbort(rt, callerAddr, txnID, txn)

	TODO() // Ensure stability across reorgs (consider having proposer supply ID?)
	return txnID
}

func (a *MultiSigActor) Approve(rt vmr.Runtime, txnID TxnID, txn MultiSigTransaction) {
	vmr.RT_ValidateImmediateCallerIsSignable(rt)
	callerAddr := rt.ImmediateCaller()
	a._rtValidateAuthorizedPartyOrAbort(rt, callerAddr)
	a._rtApproveTransactionOrAbort(rt, callerAddr, txnID, txn)
}

func (a *MultiSigActor) AddAuthorizedParty(rt vmr.Runtime, actorID abi.ActorID) {
	// Can only be called by the multisig wallet itself.
	rt.ValidateImmediateCallerIs(rt.CurrReceiver())

	h, st := a.State(rt)
	st.AuthorizedParties[actorID] = true
	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) RemoveAuthorizedParty(rt vmr.Runtime, actorID abi.ActorID) {
	// Can only be called by the multisig wallet itself.
	rt.ValidateImmediateCallerIs(rt.CurrReceiver())

	h, st := a.State(rt)

	if _, found := st.AuthorizedParties[actorID]; !found {
		rt.AbortStateMsg("Party not found")
	}

	delete(st.AuthorizedParties, actorID)

	if len(st.AuthorizedParties) < st.NumApprovalsThreshold {
		rt.AbortStateMsg("Cannot decrease authorized parties below threshold")
	}

	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) SwapAuthorizedParty(rt vmr.Runtime, oldActorID abi.ActorID, newActorID abi.ActorID) {
	// Can only be called by the multisig wallet itself.
	rt.ValidateImmediateCallerIs(rt.CurrReceiver())

	h, st := a.State(rt)

	if _, found := st.AuthorizedParties[oldActorID]; !found {
		rt.AbortStateMsg("Party not found")
	}

	if _, found := st.AuthorizedParties[oldActorID]; !found {
		rt.AbortStateMsg("Party already present")
	}

	delete(st.AuthorizedParties, oldActorID)
	st.AuthorizedParties[newActorID] = true

	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) ChangeNumApprovalsThreshold(rt vmr.Runtime, newThreshold int) {
	// Can only be called by the multisig wallet itself.
	rt.ValidateImmediateCallerIs(rt.CurrReceiver())

	h, st := a.State(rt)

	if newThreshold <= 0 || newThreshold > len(st.AuthorizedParties) {
		rt.AbortStateMsg("New threshold value not supported")
	}

	st.NumApprovalsThreshold = newThreshold

	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) Constructor(rt vmr.Runtime, authorizedParties autil.ActorIDSetHAMT, numApprovalsThreshold int) {

	rt.ValidateImmediateCallerIs(builtin.InitActorAddr)
	h := rt.AcquireState()

	st := MultiSigActorState{
		AuthorizedParties:     authorizedParties,
		NumApprovalsThreshold: numApprovalsThreshold,
		PendingTxns:           MultiSigTransactionHAMT_Empty(),
		PendingApprovals:      MultiSigApprovalSetHAMT_Empty(),
	}

	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) _rtApproveTransactionOrAbort(rt Runtime, callerAddr addr.Address, txnID TxnID, txn MultiSigTransaction) {

	h, st := a.State(rt)

	txnCheck, found := st.PendingTxns[txnID]
	if !found || !txnCheck.Equals(txn) {
		rt.AbortStateMsg("Requested transcation not found or not matched")
	}

	expirationExceeded := (rt.CurrEpoch() > txn.Expiration)
	if expirationExceeded {
		rt.AbortStateMsg("Transaction expiration exceeded")

		TODO()
		// Determine what to do about state accumulation over time.
		// Cannot rely on proposer to delete unexecuted transactions;
		// there is no incentive (in fact, this costs gas).
		// Could potentially amortize cost of cleanup via Cron.
	}

	AssertMsg(callerAddr.Protocol() == addr.ID, "caller address does not have ID")
	actorID, err := addr.IDFromAddress(callerAddr)
	autil.Assert(err == nil)

	st.PendingApprovals[txnID][abi.ActorID(actorID)] = true
	thresholdMet := (len(st.PendingApprovals[txnID]) == st.NumApprovalsThreshold)

	UpdateRelease_MultiSig(rt, h, st)

	if thresholdMet {
		if !st._hasAvailable(rt.CurrentBalance(), txn.Value, rt.CurrEpoch()) {
			rt.AbortArgMsg("insufficient funds unlocked")
		}

		// A sufficient number of approvals have arrived and sufficient funds have been unlocked: relay the message and delete from pending queue.
		rt.Send(
			txn.To,
			txn.Method,
			txn.Params,
			txn.Value,
		)
		a._rtDeletePendingTransaction(rt, txnID)
	}
}

func (a *MultiSigActor) _rtDeletePendingTransaction(rt Runtime, txnID TxnID) {
	h, st := a.State(rt)
	delete(st.PendingTxns, txnID)
	delete(st.PendingApprovals, txnID)
	UpdateRelease_MultiSig(rt, h, st)
}

func (a *MultiSigActor) _rtValidateAuthorizedPartyOrAbort(rt Runtime, address addr.Address) {
	AssertMsg(address.Protocol() == addr.ID, "caller address does not have ID")
	actorID, err := addr.IDFromAddress(address)
	autil.Assert(err == nil)

	h, st := a.State(rt)
	if _, found := st.AuthorizedParties[abi.ActorID(actorID)]; !found {
		rt.AbortArgMsg("Party not authorized")
	}
	Release_MultiSig(rt, h, st)
}

func Release_MultiSig(rt Runtime, h vmr.ActorStateHandle, st MultiSigActorState) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(&st))
	h.Release(checkCID)
}

func UpdateRelease_MultiSig(rt Runtime, h vmr.ActorStateHandle, st MultiSigActorState) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(&st))
	h.UpdateRelease(newCID)
}
