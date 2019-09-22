---
title: State Tree
---

The State Tree is the output of applying operations on the Filecoin Blockchain.

```go
type StateTree struct {
  Actors map[ActorID]ActorStorage
}
```
