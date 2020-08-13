---
title: Actor
weight: 1
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# VM Actor Interface
---

As mentioned above, Actors are the Filecoin equivalent of smart contracts in the Ethereum Virtual Machine. As such, Actors are very central components of the system. Any change to the current state of the Filecoin blockchain has to be triggered through an actor method invocation.

This sub-section describes the _interface_ between Actors and the Filecoin Virtual Machine. This means that most of what is described below does not strictly belong to the VM. Instead it is logic that sits on the interface between the VM and Actors logic.

There are eleven (11) types of _builtin_ Actors in total, not all of which interact with the VM. Some Actors do not invoke changes to the StateTree of the blockchain and as such do not need to have an interface to the VM. We discuss the details of all System Actors later on in the System Actors subsection.

Every Actor is identified by a Code ID (CID), not to be confused with the traditional, IPFS-style Content Identifier (also CID). The builtin actor structure is composed of the `abi.Infokee` of the actor and the actor's CID.

```go
var _ abi.Invokee = BuiltinActor{}

type BuiltinActor struct {
	actor abi.Invokee
	code  cid.Cid
}

// Code is the CodeID (cid) of the actor.
func (b BuiltinActor) Code() cid.Cid {
	return b.code
}

// Exports returns a slice of callable Actor methods.
func (b BuiltinActor) Exports() []interface{} {
	return b.actor.Exports()
}
```

The `ActorState` structure is composed of the actor's balance, in terms of tokens held by this actor, as well as a group of state methors used to query, inspect and interact with chain state. All methods take a TipSetKey as a parameter. The state looked up is the state at that tipset. A nil TipSetKey can be provided as a param, this will cause the heaviest tipset in the chain to be used.

```go
type ActorState struct {
	Balance types.BigInt
	State   interface{}
}
```

```go

// FullNode API is a low-level interface to the Filecoin network full node
type FullNode interface {
	Common

...

	// StateCall runs the given message and returns its result without any persisted changes.
	StateCall(context.Context, *types.Message, types.TipSetKey) (*InvocResult, error)
	// StateReplay returns the result of executing the indicated message, assuming it was executed in the indicated tipset.
	StateReplay(context.Context, types.TipSetKey, cid.Cid) (*InvocResult, error)
	// StateGetActor returns the indicated actor's nonce and balance.
	StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error)
	// StateReadState returns the indicated actor's state.
	StateReadState(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*ActorState, error)
	// StateListMessages looks back and returns all messages with a matching to or from address, stopping at the given height.
	StateListMessages(ctx context.Context, match *types.Message, tsk types.TipSetKey, toht abi.ChainEpoch) ([]cid.Cid, error)
}
```
