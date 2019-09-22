---
title: Network Interface
---

# Filecoin Network Interface

```
type NetworkInterface struct {
  libp2p.Node

  MountProtocol() // ...
}
```

TODO:
- explain how we use libp2p (very briefly)
- explain how other protocols mount on top of libp2p
- mount all filecoin protocols under `/filecoin/...`
- explain what libp2p protocols we use, and what for
  - graphsync
  - bitswap
  - gossipsub
  - kad-dht
  - bootstrap
