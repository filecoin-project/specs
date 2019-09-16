---
title: Filecoin VM Actors
entries:
  - standard
  - singleton
# suppressMenu: true
---

{{<hd 1 "ActorState">}}

The following data structures use _kinded_ representations for their IPLD
encodings, since the types can be inferred from the context in which they are used
(`Actor` or `UnsignedMessage`).

{{<goFile ActorState>}}
{{<goFile ActorMethod>}}
