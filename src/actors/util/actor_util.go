package util

import (
	actors "github.com/filecoin-project/specs/actors"
	deal "github.com/filecoin-project/specs/systems/filecoin_markets/storage_market/storage_deal"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
)

// Create a new entry in the balance table, with the specified initial balance.
// May fail if the specified address already exists in the table.
func BalanceTable_WithNewAddressEntry(table BalanceTableHAMT, address addr.Address, initBalance actors.TokenAmount) (
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
func BalanceTable_WithAdd(table BalanceTableHAMT, address addr.Address, amount actors.TokenAmount) (
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
	table BalanceTableHAMT, address addr.Address, amount actors.TokenAmount) (
	ret BalanceTableHAMT, amountSubtracted actors.TokenAmount, ok bool) {

	return BalanceTable_WithExtractPartial(table, address, amount, actors.TokenAmount(0))
}

// Extract the given amount from the given address's balance table entry, subject to the requirement
// of a minimum balance `minBalanceMaintain`. If not possible to withdraw the entire amount
// requested, then the balance will remain unchanged.
func BalanceTable_WithExtract(
	table BalanceTableHAMT, address addr.Address, amount actors.TokenAmount, minBalanceMaintain actors.TokenAmount) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Extract as much as possible (may be zero) up to the specified amount from the given address's
// balance table entry, subject to the requirement of a minimum balance `minBalanceMaintain`.
func BalanceTable_WithExtractPartial(
	table BalanceTableHAMT, address addr.Address, amount actors.TokenAmount, minBalanceMaintain actors.TokenAmount) (
	ret BalanceTableHAMT, amountExtracted actors.TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Extract all available from the given address's balance table entry.
func BalanceTable_WithExtractAll(table BalanceTableHAMT, address addr.Address) (
	ret BalanceTableHAMT, amountExtracted actors.TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Determine whether the given address's entry in the balance table meets the required minimum
// `minBalanceMaintain`.
func BalanceTable_IsEntrySufficient(
	table BalanceTableHAMT, address addr.Address, minBalanceMaintain actors.TokenAmount) (ret bool, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Retrieve the balance table entry corresponding to the given address.
func BalanceTable_GetEntry(
	table BalanceTableHAMT, address addr.Address) (
	ret actors.TokenAmount, ok bool) {

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
	AssertMsg(x.StartIndex() <= x.EndIndex(), "index %d > end %d", x.StartIndex(), x.EndIndex())

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

func MinerSetHAMT_Empty() MinerSetHAMT {
	IMPL_FINISH()
	panic("")
}

func ActorIDSetHAMT_Empty() ActorIDSetHAMT {
	IMPL_FINISH()
	panic("")
}
