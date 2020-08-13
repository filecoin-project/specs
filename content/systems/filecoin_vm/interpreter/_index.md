---
title: Interpreter
weight: 7
dashboardWeight: 1
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# VM Interpreter - Message Invocation (Outside VM)
---

The VM interpreter orchestrates the execution of messages from a tipset on that tipset's parent state,
producing a new state and a sequence of message receipts. The CIDs of this new state and of the receipt
collection are included in blocks from the subsequent epoch, which must agree about those CIDs 
in order to form a new tipset.

Every state change is driven by the execution of a message.
The messages from all the blocks in a tipset must be executed in order to produce a next state.
All messages from the first block are executed before those of second and subsequent blocks in the
tipset. For each block, BLS-aggregated messages are executed first, then SECP signed messages.

## Implicit messages

In addition to the messages explicitly included in each block, a few state changes at each epoch
are made by implicit messages. Implicit messages are not transmitted between nodes, but constructed
by the interpreter at evaluation time.

For each block in a tipset, an implicit message:

- invokes the block producer's miner actor to process the (already-validated) election PoSt submission,
as the first message in the block;
- invokes the reward actor to pay the block reward to the miner's owner account, as the final message in the block;

For each tipset, an implicit message:

- invokes the cron actor to process automated checks and payments, as the final message in the tipset.

All implicit messages are constructed with a `From` address being the distinguished system account actor.
They specify a gas price of zero, but must be included in the computation.
They must succeed (have an exit code of zero) in order for the new state to be computed.
Receipts for implicit messages are not included in the receipt list; only explicit messages have an
explicit receipt. 

## Gas payments

In most cases, the sender of a message pays the miner which produced the block including that message
a gas fee for its execution.

The gas payments for each message execution are paid to the miner owner account immediately after
that message is executed. There are no encumbrances to either the block reward or gas fees earned: 
both may be spent immediately.  

## Duplicate messages

Since different miners produce blocks in the same epoch, multiple blocks in a single tipset may 
include the same message (identified by the same CID). 
When this happens, the message is processed only the first time it is encountered in the tipset's
canonical order. Subsequent instances of the message are ignored and do not result in any 
state mutation, produce a receipt, or pay gas to the block producer. 

The sequence of executions for a tipset is thus summarised:

- pay reward for first block
- process election post for first block
- messages for first block (BLS before SECP)
- pay reward for second block
- process election post for second block
- messages for second block (BLS before SECP, skipping any already encountered)
- [... subsequent blocks ...]
- cron tick 

## Message validity and failure
Every message in a valid block can be processed and produce a receipt (note that block validity
implies all messages are syntactically valid -- see [Message Syntax](message#message-syntax-validation) -- and correctly signed).
However, execution may or may not succeed, depending on the state to which the message is applied. If the execution
of a message fails, the corresponding receipt will carry a non-zero exit code. 

If a message fails due to a reason that can reasonably be attributed to the miner including a
message that could never have succeeded in the parent state, or because the sender lacks funds
to cover the maximum message cost, then the miner pays a penalty by burning the gas fee 
(rather than the sender paying fees to the block miner).

The only state changes resulting from a message failure are either:

- incrementing of the sending actor's `CallSeqNum`, and payment of gas fees from the sender to the owner of the miner of the block including the message; or
- a penalty equivalent to the gas fee for the failed message, burnt by the miner (sender's `CallSeqNum` unchanged).
 
A message execution will fail if, in the immediately preceding state:

- the `From` actor does not exist in the state (miner penalized),
- the `From` actor is not an account actor (miner penalized),
- the `CallSeqNum` of the message does not match the `CallSeqNum` of the `From` actor (miner penalized),
- the `From` actor does not have sufficient balance to cover the sum of the message `Value` plus the
maximum gas cost, `GasLimit * GasPrice` (miner penalized),
- the `To` actor does not exist in state and the `To` address is not a pubkey-style address,
- the `To` actor exists (or is implicitly created as an account) but does not have a method corresponding to the non-zero `MethodNum`,
- deserialized `Params` is not an array of length matching the arity of the `To` actor's `MethodNum` method,
- deserialized `Params` are not valid for the types specified by the `To` actor's `MethodNum` method,
- the invoked method consumes more gas than the `GasLimit` allows,
- the invoked method exits with a non-zero code (via `Runtime.Abort()`), or
- any inner message sent by the receiver fails for any of the above reasons.

Note that if the `To` actor does not exist in state and the address is a valid `H(pubkey)` address, 
it will be created as an account actor.

(You can see the _old_ VM interpreter [here](vm_interpreter_old) )

## `vm/interpreter` interface

{{<embed src="vm_interpreter.id"  lang="go" >}}

## `vm/interpreter` implementation

{{<embed src="vm_interpreter.go"  lang="go" >}}

## `vm/interpreter/registry`

{{<embed src="vm_registry.go"  lang="go" >}}
