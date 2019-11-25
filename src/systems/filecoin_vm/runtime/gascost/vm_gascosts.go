package runtime

import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"

// TODO: assign all of these.
var (
	// SimpleValueSend is the amount of gas charged for sending value from one
	// contract to another, without executing any other code.
	SimpleValueSend = msg.GasAmount(1)

	// // ActorLookupFail is the amount of gas charged for a failure to lookup
	// // an actor
	// ActorLookupFail = msg.GasAmount(1)

	// CodeLookupFail is the amount of gas charged for a failure to lookup
	// code in the VM's code registry.
	CodeLookupFail = msg.GasAmount(1)

	// ApplyMessageFail represents the gas cost for failures to apply message.
	// These failures are basic failures encountered at first application.
	ApplyMessageFail = msg.GasAmount(1)

	// TODO: determine these costs
	PublicKeyCryptoOp = msg.GasAmount(50)
)
