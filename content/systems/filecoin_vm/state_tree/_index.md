---
title: State Tree
weight: 2
dashboardWeight: 1.5
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# State Tree
---

The State Tree is the output of the execution of any operation applied on the Filecoin Blockchain. The on-chain (i.e., VM) state data structure is a map (in the form of a Hash Array Mapped Trie - HAMT) that binds addresses to actor states.

The current State Tree function is called by the VM upon every actor method invocation.

```go
func (vm *VM) StateTree() types.StateTree {
	return vm.cstate
}
```

The state of actors is stored in the State Tree by the actor's ID.

```go
type StateTree struct {
	root  *hamt.Node
	Store cbor.IpldStore

	snaps *stateSnaps
}
```

The `LoadStateTree` provides the current State Tree

```go
func LoadStateTree(cst cbor.IpldStore, c cid.Cid) (*StateTree, error) {
	nd, err := hamt.LoadNode(context.Background(), cst, c, hamt.UseTreeBitWidth(5))
	if err != nil {
		log.Errorf("loading hamt node %s failed: %s", c, err)
		return nil, err
	}

	return &StateTree{
		root:  nd,
		Store: cst,
		snaps: newStateSnaps(),
	}, nil
}
```

`stateSnaps` are _snapshots_ of the current state for the corresponding actor that invokes some execution. Snapshots can be arranged in _layers_, which represent a sequence of snapshots in time, are kept in a cache and are mapped to actors' addresses.


```go
type stateSnaps struct {
	layers []*stateSnapLayer
}

type stateSnapLayer struct {
	actors       map[address.Address]streeOp
	resolveCache map[address.Address]address.Address
}

func newStateSnapLayer() *stateSnapLayer {
	return &stateSnapLayer{
		actors:       make(map[address.Address]streeOp),
		resolveCache: make(map[address.Address]address.Address),
	}
}
```

There are several functions that relate to state snapshots:

```go
// Add a new layer based on a State Tree snapshot
func (ss *stateSnaps) addLayer() {}

// Drop an older layer from the layers cache
func (ss *stateSnaps) dropLayer() {}

// Merge a snapshot to a previous layer
func (ss *stateSnaps) mergeLastLayer() {}

// Resolve an actor's address from a snapshot
func (ss *stateSnaps) resolveAddress(addr address.Address) (address.Address, bool) {}

// Cache actor's address
func (ss *stateSnaps) cacheResolveAddress(addr, resa address.Address) {}

// Get an actor's address
func (ss *stateSnaps) getActor(addr address.Address) (*types.Actor, error) {}

// Set an actor's address corresponding to a snapshot
func (ss *stateSnaps) setActor(addr address.Address, act *types.Actor) {}

// Delete an actor's address in a snapshot
func (ss *stateSnaps) deleteActor(addr address.Address) {}
```

These snapshot manipulation functions (above) are used to perform (or avoid performing, if not appropriate) several actions on the State Tree itself:

```go
// LookupID gets the ID address of this actor's `addr` stored in the `InitActor`.
func (st *StateTree) LookupID(addr address.Address) (address.Address, error) {}

// GetActor returns the actor from any type of `addr` provided.
func (st *StateTree) GetActor(addr address.Address) (*types.Actor, error) {}

// Merge State Tree snapshot
func (st *StateTree) ClearSnapshot() {
	st.snaps.mergeLastLayer()
}

func (st *StateTree) DeleteActor(addr address.Address) error {}

func (st *StateTree) Revert() error {
	st.snaps.dropLayer()
	st.snaps.addLayer()

	return nil
}

func (st *StateTree) MutateActor(addr address.Address, f func(*types.Actor) error) error {}

```

{{<hint warning>}}
TODO

- Add ConvenienceAPI state to provide more user-friendly views.
{{</hint>}}
