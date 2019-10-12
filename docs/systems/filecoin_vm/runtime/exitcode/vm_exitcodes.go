package runtime

import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"

// TODO: assign all of these.
var (
	// OK is the success return value, similar to unix exit code 0.
	OK = msg.ExitCode(0)

	// ActorNotFound represents a failure to find an actor.
	ActorNotFound = msg.ExitCode(1)

	// ActorCodeNotFound represents a failure to find the code for a
	// particular actor in the VM registry.
	ActorCodeNotFound = msg.ExitCode(2)

	// InvalidMethod represents a failure to find a method in
	// an actor
	InvalidMethod = msg.ExitCode(3)

	// InsufficientFunds represents a failure to apply a message, as
	// it did not carry sufficient funds for its application.
	InsufficientFunds = msg.ExitCode(4)

	// InvalidCallSeqNum represents a message invocation out of sequence.
	// This happens when message.CallSeqNum is not exactly actor.CallSeqNum + 1
	InvalidCallSeqNum = msg.ExitCode(5)
)
