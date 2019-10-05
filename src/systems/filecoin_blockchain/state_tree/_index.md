---
title: VM State Tree
---

The State Tree is the output of applying operations on the Filecoin Blockchain.

```go
type StateTree struct {
  Actors map[ActorID]ActorStorage
}
```

TODO

- Add ConvenienceAPI state to provide more user-friendly views.
