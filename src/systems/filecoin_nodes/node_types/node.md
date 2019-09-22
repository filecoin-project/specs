---
title: Node Interface
---

```
type FilecoinNode interface {
  Repository() Repository
  Network() Network
  Clock() Clock
  FileStore() FileStore
}
```
