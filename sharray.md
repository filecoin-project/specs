# Sharray

> Sharded IPLD Array

The Sharray is a fixed width IPLD tree structure used to store an array of items.

## Overview

Each node has a height, and a number of items. The number of items must not exceed the trees given width.

The tree must not be sparse, if the tree represents an array of N items, then the left `N/Width` leaves must contain the first `N` items. If `N` is not evenly divisible by `Width` then the final leaf must contain the final remainder. Every node in the tree must have its height set to one less than the node above it. Nodes with a height of 0 contain array values, and nodes with heights greater than zero contain the cids of their child nodes.

## IPLD Representation

Each sharray node is represented by an IPLD node of the following schema:

```
type Node struct {
  height Int
  items [Item]
} representation tuple
```

`Item` may be either a direct value, if `height == 0`, or the Cid of a child node if `height > 0`.

(For details on IPLD Schemas, see the [IPLD Schema Spec (draft)](https://github.com/ipld/specs/blob/dcbfb25468092be796bab90e90e3f2535fdeddc7/schema/representations.md))

We use DAG-CBOR for serialization, and blake2b-256 for hashing.

## Operations

#### `create(items)`

> Create a sharray from a given set of items

```go
func create(items []Item) Cid {
  var layer []cid.Cid

  for len(items) > 0 {
    // get the next 'Width' items from the input items
    vals, items = getNextN(Width, items)

    nd := Node{
      height: 0,
      items: vals,
    }

    layer.append(nd.Cid())
  }

  var nextLayer []cid.Cid
  for height := 1; len(layer) > 1; height++ {
    for len(layer) > 0 {
      vals, layer = getNextN(Width, layer)

      nd := Node{
        height: height,
        items: vals,
      }

      nextLayer.append(nd.Cid())
    }
    layer = nextLayer
    nextLayer = []
  }

  return nextLayer[0]
}
```



#### `get(i)`

> Get the element at index `i`

```go
func (n node) get(i int) Item {
  if n.Height == 0 {
    return n.Array[i]
  }

  childWidth := Pow(Width, n.Height)

  child := loadNode(n.Array[i / childWidth])
  return child.get( i % childWidth )
}
```
