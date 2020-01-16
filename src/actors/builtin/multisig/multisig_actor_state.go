package multisig

import (
	abi "github.com/filecoin-project/specs/actors/abi"
	autil "github.com/filecoin-project/specs/actors/util"
	cid "github.com/ipfs/go-cid"
)

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
