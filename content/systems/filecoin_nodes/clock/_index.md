---
title: Clock
weight: 4
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# Clock

Filecoin assumes weak clock synchrony amongst participants in the system. That is, the system relies on participants having access to a globally synchronized clock (tolerating some bounded offset).

Filecoin relies on this system clock in order to secure consensus. Specifically, the clock is necessary to support validation rules that prevent block producers from mining blocks with a future timestamp and running leader elections more frequently than the protocol allows.

## Clock uses

The Filecoin system clock is used:

- by syncing nodes to validate that incoming blocks were mined in the appropriate epoch given their timestamp (see [Block Validation](block#block-syntax-validation)). This is possible because the system clock maps all times to a unique epoch number totally determined by the start time in the genesis block.
- by syncing nodes to drop blocks coming from a future epoch
- by mining nodes to maintain protocol liveness by allowing participants to try leader election in the next round if no one has produced a block in the current round (see [Storage Power Consensus](storage_power_consensus)).

In order to allow miners to do the above, the system clock must:

1. Have low enough offset relative to other nodes so that blocks are not mined in epochs considered future epochs from the perspective of other nodes (those blocks should not be validated until the proper epoch/time as per [validation rules](block#block-semantic-validation)).
2. Set epoch number on node initialization equal to `epoch = Floor[(current_time - genesis_time) / epoch_time]`

It is expected that other subsystems will register to a `NewRound()` event from the clock subsystem.

## Clock Requirements

Clocks used as part of the Filecoin protocol should be kept in sync, with offset less than 1 second so as to enable appropriate validation.

Computer-grade crystals can be expected to deviate by [1ppm](https://www.hindawi.com/journals/jcnc/2008/583162/) (i.e. 1 microsecond every second, or 0.6 seconds per week). Therefore, in order to respect the requirement above:

- Nodes SHOULD run an NTP daemon (e.g. timesyncd, ntpd, chronyd) to keep their clocks synchronized to one or more reliable external references.
  - We recommend the following sources:
    - **`pool.ntp.org`** ([details](https://www.ntppool.org/en/use.html))
    - `time.cloudflare.com:1234` ([details](https://www.cloudflare.com/time/))
    - `time.google.com` ([details](https://developers.google.com/time))
    - `time.nist.gov` ([details](https://tf.nist.gov/tf-cgi/servers.cgi))
- Larger mining operations MAY consider using local NTP/PTP servers with GPS references and/or frequency-stable external clocks for improved timekeeping.

Mining operations have a strong incentive to prevent their clock skewing ahead more than one epoch to keep their block submissions from being rejected. Likewise they have an incentive to prevent their clocks skewing behind more than one epoch to avoid partitioning themselves off from the synchronized nodes in the network.
