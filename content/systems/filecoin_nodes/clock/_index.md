---
title: Clock
weight: 4
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# Clock

Filecoin assumes weak clock synchrony amongst participants in the system. That is, the system relies on participants having access to a globally synchronized clock (tolerating some bounded drift).

Filecoin relies on this system clock in order to secure consensus.  Specifically the clock is necessary to support validation rules that prevent block producers from mining blocks with a future timstamp, and running leader elections more frequently than the protocol allows.


## Clock uses

The Filecoin system clock is used:

- by syncing nodes to validate that incoming blocks were mined in the appropriate epoch given their timestamp (see [Block Validation](block#block-syntax-validation)).  This is possible because the system clock maps all times to a unique epoch number totally determined by the start time in the genesis block.
- by syncing nodes to drop blocks coming from a future epoch
- by mining nodes to maintain protocol liveness by allowing participants to try leader election in the next round if no one has produced a block in the current round (see [Storage Power Consensus](storage_power_consensus)).

In order to allow miners to do the above, the system clock must:

1. Have low enough clock drift (sub 1s) relative to other nodes so that blocks are not mined in epochs considered future epochs from the persective of other nodes (those blocks should not be validated until the proper epoch/time as per [validation rules](block#block-semantic-validation)).
2. Set epoch number on node initialization equal to `epoch = Floor[(current_time - genesis_time) / epoch_time]`

It is expected that other subsystems will register to a `NewRound()` event from the clock subsystem.

## Clock Requirements

Clocks used as part of the Filecoin protocol should be kept in sync, with drift less than 1 second so as to enable appropriate validation.

Computer-grade clock crystals can be expected to have drift rates on the order of [1ppm](https://www.hindawi.com/journals/jcnc/2008/583162/) (i.e. 1 microsecond every second or .6 seconds a week), therefore, in order to respect the above-requirement,

- clients SHOULD query an NTP server (`pool.ntp.org` is recommended) on an hourly basis to adjust clock skew.
  - We recommend one of the following:
    - `pool.ntp.org` (can be catered to a [specific zone](https://www.ntppool.org/zone))
    - `time.cloudflare.com:1234` (more on [Cloudflare time services](https://www.cloudflare.com/time/))
    - `time.google.com` (more on [Google Public NTP](https://developers.google.com/time))
    - `ntp-b.nist.gov` ([NIST](https://tf.nist.gov/tf-cgi/servers.cgi) servers require registration)
  - We further recommend making three (3) measurements in order to drop outliers
- clients MAY consider using cesium clocks instead for accurate synchrony within larger mining operations

Mining operations have a strong incentive to prevent their clock from drifting ahead more than one epoch to keep their block submissions from being rejected.  Likewise they have an incentive to prevent their clocks from drifting behind more than one epoch to avoid partitioning themselves off from the synchronized nodes in the network.

