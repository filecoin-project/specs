---
menuTitle: Interpreter
statusIcon: ⚠️
title: VM Interpreter - Message Invocation (Outside VM)
---

The VM interpreter orchestrates the execution of messages from a tipset on that tipset's parent state,
producing a new state and a sequence of message receipts. The CIDs of this new state and of the receipt
collection are included in blocks from the subsequent epoch, which must agree about those CIDs 
in order to form a new tipset.

The messages from all the blocks in a tipset must be executed in order to produce a next state.
All messages from the first block are executed before those of second and subsequent blocks in the
tipset. For each block, BLS-aggregated messages are executed first, then SECP signed messages.

In addition, for each block:

- the block reward is paid to the miner owner account, and 
- the block producer's election PoSt is processed by an implicit invocation on the associated actor  

via implicit messages before that block's explicit messages are executed. 

For each tipset:

- the `CronActor`'s tick method is invoked implicitly after all the blocks' messages.

These implicit messages specify a gas price of zero, but must be included in the computation.
They must succeed (have an exit code of zero) in order for the new state to be
computed. Receipts for them are excluded from the receipt list; only explicit messages have an 
explicit receipt. 

The gas payment for each message execution is paid to the miner owner account immediately after
that message is executed. There are no encumbrances to either the block reward or gas fees earned: 
both may be spent immediately, including by a message in the same block in which they are earned.  

Since different miners produce blocks in the same epoch, multiple blocks in a single tipset may 
include the same message (identified by the same CID). 
When this happens, the message is processed only the first time it is encountered in the tipset's
canonical order. Subsequent instances of the message are ignored and do not result in any 
state mutation or produce a receipt. 

The sequence of executions for a tipset is thus summarised:

- pay reward for first block
- process election post for first block
- messages for first block (BLS before SECP)
- pay reward for second block
- process election post for second block
- messages for second block (BLS before SECP, skipping any already encountered)
- [... subsequent blocks ...]
- cron tick 

(You can see the _old_ VM interpreter [here](docs/systems/filecoin_vm/vm_interpreter_old) )

# `vm/interpreter` interface

{{< readfile file="vm_interpreter.id" code="true" lang="go" >}}

# `vm/interpreter` implementation

{{< readfile file="vm_interpreter.go" code="true" lang="go" >}}

# `vm/interpreter/registry`

{{< readfile file="vm_registry.go" code="true" lang="go" >}}
