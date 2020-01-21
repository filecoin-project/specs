package util

import (
	"math/big"

	addr "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/specs-actors/actors/abi"
)

type BalanceTableHAMT map[addr.Address]abi.TokenAmount

type DealIDSetHAMT map[abi.DealID]bool
type IntToDealIDHAMT map[int64]abi.DealID

type DealIDQueue struct {
	Values     IntToDealIDHAMT
	StartIndex int64
	EndIndex   int64
}

func (x *DealIDQueue) Enqueue(dealID abi.DealID) {
	nextIndex := x.EndIndex
	x.Values[nextIndex] = dealID
	x.EndIndex = nextIndex + 1
}

func (x *DealIDQueue) Dequeue() (dealID abi.DealID, ok bool) {
	AssertMsg(x.StartIndex <= x.EndIndex, "index %d > end %d", x.StartIndex, x.EndIndex)

	if x.StartIndex == x.EndIndex {
		dealID = abi.DealID(-1)
		ok = false
		return
	} else {
		dealID = x.Values[x.StartIndex]
		delete(x.Values, x.StartIndex)
		x.StartIndex += 1
		ok = true
		return
	}
}

func DealIDQueue_Empty() DealIDQueue {
	return DealIDQueue{
		Values:     IntToDealIDHAMT_Empty(),
		StartIndex: 0,
		EndIndex:   0,
	}
}

type MinerSetHAMT map[addr.Address]bool

type ActorIDSetHAMT map[abi.ActorID]bool

type MinerEvent struct {
	MinerAddr addr.Address
	Sectors   []abi.SectorNumber // Empty for global events, such as SurprisePoSt expiration.
}

// TODO HAMT
type MinerEventSetHAMT []MinerEvent

type SectorStorageWeightDesc struct {
	SectorSize abi.SectorSize
	Duration   abi.ChainEpoch
	DealWeight big.Int
}

type SectorTermination int64

// Note: Detected fault termination (due to exceeding the limit of consecutive
// SurprisePoSt failures) is not listed here, since this does not terminate all
// sectors individually, but rather the miner as a whole.
const (
	NormalExpiration SectorTermination = iota
	UserTermination
)

// Create a new entry in the balance table, with the specified initial balance.
// May fail if the specified address already exists in the table.
func BalanceTable_WithNewAddressEntry(table BalanceTableHAMT, address addr.Address, initBalance abi.TokenAmount) (
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
func BalanceTable_WithAdd(table BalanceTableHAMT, address addr.Address, amount abi.TokenAmount) (
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
	table BalanceTableHAMT, address addr.Address, amount abi.TokenAmount) (
	ret BalanceTableHAMT, amountSubtracted abi.TokenAmount, ok bool) {

	return BalanceTable_WithExtractPartial(table, address, amount, abi.TokenAmount(0))
}

// Extract the given amount from the given address's balance table entry, subject to the requirement
// of a minimum balance `minBalanceMaintain`. If not possible to withdraw the entire amount
// requested, then the balance will remain unchanged.
func BalanceTable_WithExtract(
	table BalanceTableHAMT, address addr.Address, amount abi.TokenAmount, minBalanceMaintain abi.TokenAmount) (
	ret BalanceTableHAMT, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Extract as much as possible (may be zero) up to the specified amount from the given address's
// balance table entry, subject to the requirement of a minimum balance `minBalanceMaintain`.
func BalanceTable_WithExtractPartial(
	table BalanceTableHAMT, address addr.Address, amount abi.TokenAmount, minBalanceMaintain abi.TokenAmount) (
	ret BalanceTableHAMT, amountExtracted abi.TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Extract all available from the given address's balance table entry.
func BalanceTable_WithExtractAll(table BalanceTableHAMT, address addr.Address) (
	ret BalanceTableHAMT, amountExtracted abi.TokenAmount, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Determine whether the given address's entry in the balance table meets the required minimum
// `minBalanceMaintain`.
func BalanceTable_IsEntrySufficient(
	table BalanceTableHAMT, address addr.Address, minBalanceMaintain abi.TokenAmount) (ret bool, ok bool) {

	IMPL_FINISH()
	panic("")
}

// Retrieve the balance table entry corresponding to the given address.
func BalanceTable_GetEntry(
	table BalanceTableHAMT, address addr.Address) (
	ret abi.TokenAmount, ok bool) {

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

func MinerSetHAMT_Empty() MinerSetHAMT {
	IMPL_FINISH()
	panic("")
}

func ActorIDSetHAMT_Empty() ActorIDSetHAMT {
	IMPL_FINISH()
	panic("")
}

func MinerEventSetHAMT_Empty() MinerEventSetHAMT {
	IMPL_FINISH()
	panic("")
}
