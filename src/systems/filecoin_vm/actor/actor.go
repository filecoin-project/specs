package actor

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import util "github.com/filecoin-project/specs/util"

var IMPL_FINISH = util.IMPL_FINISH
var TODO = util.TODO

type Serialization = util.Serialization

const (
	MethodSend        = MethodNum(0)
	MethodConstructor = MethodNum(1)

	// TODO: remove this once canonical method numbers are finalized
	MethodPlaceholder = MethodNum(-(1 << 30))
)

func (st *ActorState_I) CID() ipld.CID {
	panic("TODO")
}

func (id *CodeID_I) IsBuiltin() bool {
	switch id.Which() {
	case CodeID_Case_Builtin:
		return true
	default:
		panic("Actor code ID case not supported")
	}
}

func (id *CodeID_I) IsSingleton() bool {
	if !id.IsBuiltin() {
		return false
	}

	for _, a := range []BuiltinActorID{
		BuiltinActorID_Init,
		BuiltinActorID_Cron,
		BuiltinActorID_Init,
		BuiltinActorID_StoragePower,
		BuiltinActorID_StorageMarket,
	} {
		if id.As_Builtin() == a {
			return true
		}
	}

	for _, a := range []BuiltinActorID{
		BuiltinActorID_Account,
		BuiltinActorID_PaymentChannel,
		BuiltinActorID_StorageMiner,
	} {
		if id.As_Builtin() == a {
			return false
		}
	}

	panic("Actor code ID case not supported")
}

func (x ActorSubstateCID) Ref() *ActorSubstateCID {
	return &x
}

func TokenAmount_Placeholder() TokenAmount {
	TODO()
	panic("")
}

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
