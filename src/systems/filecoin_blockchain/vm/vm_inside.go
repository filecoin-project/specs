package vm

import (
	base "github.com/filecoin-project/specs/systems/filecoin_blockchain"
    ipld "github.com/filecoin-project/specs/libraries/ipld"
)

type CID = ipld.CID

type VMMessage struct {
	To     base.Address
	Method string
	Value  TokenAmount
	Params []interface{}
}

type VMMessageReceipt struct {
}

type MessageNumber UVarint

// VMSyscalls is the interface of functions callable from within an actor.
// This defines all the functions an actor may call on the VM.
type VMSyscalls interface {
	// Fatal is a fatal error, and halts VM execution.
	// This is the normal error condition of actor execution.
	// On Fatal errors, the VM simply does not apply the state transition.
	// This is atomic across the entire contract execution.
	Fatal(string)

	// ComputeActorAddress computes the address of the contract,
	// based on the creator (invoking address) and nonce given.
	ComputeActorAddress(creator base.Address, nonce CallSeqNo) base.Address

	// Storage provides access to the VM storage layer
	Storage() VMStorage

	// ChainState is statetree accessible to all the actors
	ChainState() VMChainState

	// Send allows the current execution context to invoke methods on other actors in the system
	//
	// TODO: this should change to async -- put the message on the queue.
	//       do definied callback methods, with maybe a glue closure to align params or carry intermediate state.
	//
	// TODO: what are the return values here?
	SendMsg(to base.Address, method string, value TokenAmount, params []interface{}) ([][]byte, uint8, error)
}

// VMContext is the old syscalls version. not sure we need it.
type VMContext interface {
	// Message is the message that kicked off the current invocation
	Message() base.Message

	// Origin is the address of the account that initiated the top level invocation
	Origin() base.Address

	// Storage provides access to the VM storage layer
	Storage() Storage

	// Send allows the current execution context to invoke methods on other actors in the system
	Send(to base.Address, method string, value AttoFIL, params []interface{}) ([][]byte, uint8, error)

	// BlockHeight returns the height of the block this message was added to the chain in
	BlockHeight() base.BlockHeight
}

type VMStorage interface {
	// Put writes the given object to the storage staging area and returns its CID
	Put(interface{}) (CID, error)

	// Get fetches the given object from storage (either staging, or local) and returns
	// the serialized data.
	Get(CID) ([]byte, error)

	// Commit updates the actual stored state for the actor. This is a compare and swap
	// operation, and will fail if 'old' is not equal to the current return value of `Head`.
	// This functionality is used to prevent issues with re-entrancy
	//
	// TODO: YIKES i dont think we need commit to prevent re-entrancy. if we do, the model
	// is wrong.
	Commit(old CID, new CID) error

	// Head returns the CID of the current actor state
	Head() CID
}

// VMChainState is Chain state accessible to all contracts via the VM interface
type VMChainState interface {
	// BlockHeight returns the height of the block this message was added to the chain in
	BlockHeight() base.BlockHeight

	// RoundMessageNumber returns the number this message is in this round.
	// (BlockHeight, RoundMessageNumber) is a unique tuple per message invocation
	RoundMessageNumber() MessageNumber

	// ChainMessageNumber returns the number of messages in the chain so far.
	// TODO: this probably should be in a bunch of state accessible to the contracts.
	ChainMessageNumber() MessageNumber

	// MinerAddress who should earn this message txn id
	MinerAddress()
}
