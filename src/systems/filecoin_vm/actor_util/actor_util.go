package actor_util

import (
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	st "github.com/filecoin-project/specs/systems/filecoin_vm/state_tree"
	util "github.com/filecoin-project/specs/util"
)

type MethodParams = actor.MethodParams
type TokenAmount = actor.TokenAmount

type Serialization = util.Serialization

var Assert = util.Assert
var IMPL_FINISH = util.IMPL_FINISH

// Interface for runtime/VMContext functionality (to avoid circular dependency in Go imports)
type Has_AbortArg interface {
	AbortArg()
}

func CheckArgs(params *MethodParams, rt Has_AbortArg, cond bool) {
	if !cond {
		rt.AbortArg()
	}
}

func ArgPop(params *MethodParams, rt Has_AbortArg) Serialization {
	CheckArgs(params, rt, len(*params) > 0)
	ret := (*params)[0]
	*params = (*params)[1:]
	return ret
}

func ArgEnd(params *MethodParams, rt Has_AbortArg) {
	CheckArgs(params, rt, len(*params) == 0)
}

// Create a new entry in the balance table, with the specified initial balance.
// May fail if the specified address already exists in the table.
func BalanceTable_WithNewAddressEntry(table BalanceTableHAMT, address addr.Address, initBalance TokenAmount) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Delete the specified entry in the balance table.
// May fail if the specified address does not exist in the table.
func BalanceTable_WithDeletedAddressEntry(table BalanceTableHAMT, address addr.Address) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Add the given amount to the given address's balance table entry.
func BalanceTable_WithAdd(table BalanceTableHAMT, address addr.Address, amount TokenAmount) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Subtract the given amount (or as much as possible, without making the resulting amount negative)
// from the given address's balance table entry, independent of any minimum balance maintenance
// requirement.
// Note: ok should be set to true here, even if the operation caused the entry to hit zero.
// The only failure case is when the address does not exist in the table.
func BalanceTable_WithSubtractPreservingNonnegative(
	table BalanceTableHAMT, address addr.Address, amount TokenAmount) (
	ret BalanceTableHAMT, amountSubtracted TokenAmount, ok bool) {

	return BalanceTable_WithExtractPartial(table, address, amount, TokenAmount(0))
}

// Extract the given amount from the given address's balance table entry, subject to the requirement
// of a minimum balance `minBalanceMaintain`. If not possible to withdraw the entire amount
// requested, then the balance will remain unchanged.
func BalanceTable_WithExtract(
	table BalanceTableHAMT, address addr.Address, amount TokenAmount, minBalanceMaintain TokenAmount) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Extract as much as possible (may be zero) up to the specified amount from the given address's
// balance table entry, subject to the requirement of a minimum balance `minBalanceMaintain`.
func BalanceTable_WithExtractPartial(
	table BalanceTableHAMT, address addr.Address, amount TokenAmount, minBalanceMaintain TokenAmount) (
	ret BalanceTableHAMT, amountExtracted TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Determine whether the given address's entry in the balance table meets the required minimum
// `minBalanceMaintain`.
func BalanceTable_IsEntrySufficient(
	table BalanceTableHAMT, address addr.Address, minBalanceMaintain TokenAmount) (ret bool, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Retrieve the balance table entry corresponding to the given address.
func BalanceTable_GetEntry(
	table BalanceTableHAMT, address addr.Address) (
	ret TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

func BalanceTableHAMT_Empty() BalanceTableHAMT {
	IMPL_FINISH()
	panic("")
}

func IntToDealIDHAMT_Empty() IntToDealIDHAMT {
	IMPL_FINISH()
	panic("")
}

func DealIDSetHAMT_Empty() DealIDSetHAMT {
	IMPL_FINISH()
	panic("")
}

func DealIDQueue_Empty() DealIDQueue {
	return &DealIDQueue_I{
		Values_:     IntToDealIDHAMT_Empty(),
		StartIndex_: 0,
		EndIndex_:   0,
	}
}

func (x *DealIDQueue_I) Enqueue(dealID deal.DealID) {
	nextIndex := x.EndIndex()
	x.Values()[nextIndex] = dealID
	x.EndIndex_ = nextIndex + 1
}

func (x *DealIDQueue_I) Dequeue() (dealID deal.DealID, ok bool) {
	Assert(x.StartIndex() <= x.EndIndex())

	if x.StartIndex() == x.EndIndex() {
		dealID = deal.DealID(-1)
		ok = false
		return
	} else {
		dealID = x.Values()[x.StartIndex()]
		delete(x.Values(), x.StartIndex())
		x.StartIndex_ += 1
		ok = true
		return
	}
}

// Get the owner account address associated to a given miner actor.
func GetMinerOwnerAddress(tree st.StateTree, minerAddr addr.Address) (addr.Address, error) {
	IMPL_FINISH()
	panic("")
}

// Get the owner account address associated to a given miner actor.
func GetMinerOwnerAddress_Assert(tree st.StateTree, a addr.Address) addr.Address {
	ret, err := GetMinerOwnerAddress(tree, a)
	Assert(err == nil)
	return ret
}
