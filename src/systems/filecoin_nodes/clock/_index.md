---
title: Clock
---

```go
type Time string // ISO nano timestamp
type UnixTime int // unix timestamp
type Round int // Blockchain round

type Clock interface {
  UTCTime() Time
  UnixNano() UnixTime

  CurrentRound() Round
  LastRoundObserved() Round
}
```

TODO:

- explain why we need a system clock
- explain where it is used
  - for rejecting/accepting blocks
- explain synchrony requirements
  - small clock drift -- <2s
  - very important to have accurate time
- explain how we can resync
  - Network
    - recommend various NTP servers
  - Cesium clocks
- Future work:
  - VDF Clocks
