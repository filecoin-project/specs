---
title: Message
weight: 4
dashboardWeight: 1.5
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# VM Message - Actor Method Invocation

A message is the unit of communication between two actors, and thus the primitive cause of changes
in state. A message combines:

- a token amount to be transferred from the sender to the receiver, and
- a method with parameters to be invoked on the receiver (optional).

Actor code may send additional messages to other actors while processing a received message. 
Messages are processed synchronously: an actor waits for a sent message to complete before resuming control.

The processing of a message consumes units of computation and storage denominated in gas. 
A message's gas limit provides an upper bound on its computation. The sender of a message pays 
for the gas units consumed by a message's execution (including all nested messages) at a
gas price they determine. A block producer chooses which messages to include in a block and is
rewarded according to each message's gas price and consumption, forming a market.

## Message syntax validation

A syntactically invalid message must not be transmitted, retained in a message pool, or
included in a block.

A syntactically valid `UnsignedMessage`:

- has a well-formed, non-empty `To` address,
- has a well-formed, non-empty `From` address, 
- has a non-negative `CallSeqNum`,
- has `Value` no less than zero and no greater than the total token supply (`2e9 * 1e18`), and
- has a non-negative `MethodNum`,
- has non-empty `Params` only if `MethodNum` is zero,
- has non-negative `GasPrice`,
- has `GasLimit` that is at least equal to the gas consumption associated with the message's serialized bytes,
- has `GasLimit` that is no greater than the block gas limit network parameter.

When transmitted individually (before inclusion in a block), a message is packaged as
`SignedMessage`, regardless of signature scheme used. A valid signed message:

- has a total serialized size no greater than `message.MessageMaxSize`.

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
will produce a receipt from execution. 
However, a message may fail to execute to completion, in which case it will not effect the desired state change.  

The reason for this "no message semantic validation" policy is that the state that a message will
be applied to cannot be known before the message is executed _as part of a tipset_. A block producer
does not know whether another block will precede it in the tipset, thus altering the state to
which the block's messages will apply from the declared parent state.

{{<embed src="message.id" lang="go" >}}

{{<embed src="message.go" lang="go" >}}
