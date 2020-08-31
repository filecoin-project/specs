---
title: "Sharded IPLD Array"
weight: 1
dashboardWeight: 0.2
dashboardState: wip
dashboardAudit: n/a
---

# Sharded IPLD Array

The Sharray is an IPLD tree structure used to store an array of items. It is designed for usecases that know all items at the time of creation and do not need insertion or deletion.

## IPLD Representation

Each sharray node is represented by an IPLD node of the following schema:

```text
type Node struct {
  height Int
  items [Item]
} representation tuple
```

`Item` may be either a direct value, if `height == 0`, or the Cid of a child node if `height > 0`.

(For details on IPLD Schemas, see the [IPLD Schema Spec (draft)](https://github.com/ipld/specs/blob/dcbfb25468092be796bab90e90e3f2535fdeddc7/schema/representations.md))

We use DAG-CBOR for serialization, and blake2b-256 for hashing.

## Construction

The tree must not be sparse.
Given an array of size `N` and a fixed width of `W`:

- The left `floor(N/W)` leaves contain the first `N` items.
- If `N % W != 0` the final leaf contains the final remainder.
- The tree is perfectly balanced.
- The height is the distance from the leaves, not the root.
- Leaves (nodes with a height of 0) contain array values.
- Inner nodes (nodes with height greater than zero) contain the cids of their child nodes.

## Operations

### `create(items)`

> Create a sharray from a given set of items

```go
func create(items []Item) Cid {
	var layer cidQueue

	itemQ := queue(items)
	for !itemQ.Empty() {
		// get the next 'Width' items from the input items
		vals := itemQ.PopN(width)

		nd := Node{
			height: 0,
			items:  vals,
		}

		// persist the node to the datastore
		storeNode(nd)

		layer.push(nd.Cid())
	}

	var nextLayer cidQueue
	for height := 1; layer.Len() > 1; height++ {
		for layer.Len() > 0 {
			vals := layer.PopN(width)

			nd := Node{
				height: height,
				items:  vals,
			}

			storeNode(nd)

			nextLayer.append(nd.Cid())
		}
		layer = nextLayer
		nextLayer.ClearItems()
	}

	return nextLayer.First()
}
```

### `get(i)`

> Get the element at index `i`

```go
func (n node) get(i int) Item {
	if n.Height == 0 {
		return n.Array[i]
	}

	childWidth := Pow(Width, n.Height)

	child := loadNode(n.Array[i/childWidth])
	return child.get(i % childWidth)
}
```
