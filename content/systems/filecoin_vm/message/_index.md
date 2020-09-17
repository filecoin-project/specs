---
title: Message
weight: 4
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# VM Message - Actor Method Invocation

A message is the unit of communication between two actors, and thus the primitive cause of changes
in state. A message combines:

- a token amount to be transferred from the sender to the receiver, and
- a method with parameters to be invoked on the receiver (optional/where applicable).

Actor code may send additional messages to other actors while processing a received message. 
Messages are processed synchronously, that is, an actor waits for a sent message to complete before resuming control.

The processing of a message consumes units of computation and storage, both of which are denominated in gas. 
A message's _gas limit_ provides an upper bound on the computation required to process it. The sender of a message pays 
for the gas units consumed by a message's execution (including all nested messages) at a
gas price they determine. A block producer chooses which messages to include in a block and is
rewarded according to each message's gas price and consumption, forming a market.

## Message syntax validation

A syntactically invalid message must not be transmitted, retained in a message pool, or
included in a block. If an invalid message is received, it should be dropped and not propagated further.

When transmitted individually (before inclusion in a block), a message is packaged as
`SignedMessage`, regardless of signature scheme used. A valid signed message has a total serialized size no greater than `message.MessageMaxSize`.

```go
type SignedMessage struct {
	Message   Message
	Signature crypto.Signature
}
```

A syntactically valid `UnsignedMessage`:

- has a well-formed, non-empty `To` address,
- has a well-formed, non-empty `From` address, 
- has `Value` no less than zero and no greater than the total token supply (`2e9 * 1e18`), and
- has non-negative `GasPrice`,
- has `GasLimit` that is at least equal to the gas consumption associated with the message's serialized bytes,
- has `GasLimit` that is no greater than the block gas limit network parameter.


```go
type Message struct {
	// Version of this message (has to be non-negative)
	Version uint64

	// Address of the receiving actor.
	To   address.Address
	// Address of the sending actor.
	From address.Address

	CallSeqNum uint64

	// Value to transfer from sender's to receiver's balance.
	Value BigInt

	// GasPrice is a Gas-to-FIL cost
	GasPrice BigInt
	// Maximum Gas to be spent on the processing of this message
	GasLimit int64

	// Optional method to invoke on receiver, zero for a plain value transfer.
	Method abi.MethodNum
	//Serialized parameters to the method.
	Params []byte
}
```

There should be several functions able to extract information from the `Message struct`, such as the sender and recipient addresses, the value to be transferred, the required funds to execute the message and the CID of the message.

Given that Transaction Messages should eventually be included in a Block and added to the blockchain, the validity of the message should be checked with regard to the sender and the receiver of the message, the value (which should be non-negative and always smaller than the circulating supply), the gas price (which again should be non-negative) and the `BlockGasLimit` which should not be greater than the block's gas limit.


## Message semantic validation

Semantic validation refers to validation requiring information outside of the message itself.

A semantically valid `SignedMessage` must carry a signature that verifies the payload as having
been signed with the public key of the account actor identified by the `From` address. 
Note that when the `From` address is an ID-address, the public key must be
looked up in the state of the sending account actor in the parent state identified by the block.

Note: the sending actor must exist _in the parent state identified by the block_ that includes the message.
This means that it is not valid for a single block to include a message that creates a new account 
actor and a message from that same actor. 
The first message from that actor must wait until a subsequent epoch.
Message pools may exclude messages from an actor that is not yet present in the chain state.

There is no further semantic validation of a message that can cause a block including the message 
to be invalid. Every syntactically valid and correctly signed message can be included in a block and 
will produce a receipt from execution. The `MessageReceipt sturct` includes the following:

```go
type MessageReceipt struct {
	ExitCode exitcode.ExitCode
	Return   []byte
	GasUsed  int64
}
```


However, a message may fail to execute to completion, in which case it will not trigger the desired state change.  

The reason for this "no message semantic validation" policy is that the state that a message will
be applied to cannot be known before the message is executed _as part of a tipset_. A block producer
does not know whether another block will precede it in the tipset, thus altering the state to
which the block's messages will apply from the declared parent state.


{{<embed src="github:filecoin-project/lotus/chain/types/message.go"  lang="go">}}