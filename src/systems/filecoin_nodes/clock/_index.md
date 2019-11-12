---
title: Clock
statusIcon: ⚠️
---

{{< readfile file="clock_subsystem.id" code="true" lang="go" >}}
{{< readfile file="clock_subsystem.go" code="true" lang="go" >}}

Filecoin assumes weak clock synchrony amongst participants in the system. That is, the system relies on participants having access to a globally synchronized clock, tolerating bounded delay in honest clock lower than epoch time (more on this in a forthcoming paper).

Filecoin relies on this system clock in order to secure consensus, specifically ensuring that participants are only running leader elections once per epoch and enabling miners to catch such deviations from the protocol. Given a system start and epoch time by the genesis block, the system clock allows miners to associate epoch and wall clock time, thereby enabling them to reason about block validity and give the protocol liveness.

## Clock uses
Specifically, the Filecoin system clock is used:

- to validate incoming blocks and ensure they were mined in the appropriate round, looking at the wall clock time in conjunction with the block's `ElectionProof` (which contains the epoch number) (see {{<sref leader_election>}} and {{<sref block_validation>}}).
- to help protocol convergence by giving miners a specific cutoff after which to reject incoming blocks in this round (see {{<sref chain_sync>}}).
- to maintain protocol liveness by allowing participants to try leader election in the next round if no one has produced a block in this round (see {{<sref storage_power_consensus>}}).

In order to allow miners to do the above, the system clock must:

1. have low clock drift: at most on the order of 1s (i.e. markedly lower than epoch time) at any given time.
2. maintain accurate network time over many epochs: resyncing and enforcing accurate network time.
3. set epoch number on client initialization equal to `epoch ~= (current_time - genesis_time) / epoch_time`

It is expected that other subsystems will register to a NewRound() event from the clock subsystem.

## Clock Requirements

Computer-grade clock crystals can be expected to have drift rates on the order of [1ppm](https://www.hindawi.com/journals/jcnc/2008/583162/) (i.e. 1 microsecond every second or .6 seconds a week), therefore, in order to respect the first above-requirement,

- clients SHOULD query an NTP server (`pool.ntp.org` is recommended) on an hourly basis to adjust clock skew.
  - We recommend one of the following:
    - `pool.ntp.org` (can be catered to a [specific zone](https://www.ntppool.org/zone))
    - `time.cloudflare.com:1234` (more on [Cloudflare time services](https://www.cloudflare.com/time/))
    - `time.google.com` (more on [Google Public NTP](https://developers.google.com/time))
    - `ntp-b.nist.gov` ([NIST](https://tf.nist.gov/tf-cgi/servers.cgi) servers require registration)
  - We further recommend making 3 measurements in order to drop by using the network to drop outliers
  - See how [go-ethereum does this](https://github.com/ethereum/go-ethereum/blob/master/p2p/discv5/ntp.go) for inspiration
- clients CAN consider using cesium clocks instead for accurate synchrony within larger mining operations

Assuming a majority of rational participants, the above should lead to relatively low skew over time, with seldom more than 10-20% clock skew that should be rectified periodically by the network, as is the case in other networks. This assumption can be tested over time by ensuring that:

- (real-time) epoch time is as dictated by the protocol
- (historical) the current epoch number is as expected

## Future work

If either of the above metrics show significant network skew over time, future versions of Filecoin may include potential timestamp/epoch correction periods at regular intervals.

More generally, future versions of the Filecoin protocol will use Verifiable Delay Functions (VDFs) to strongly enforce block time and fulfill this leader election requirement; we choose to explicitly assume clock synchrony until hardware VDF security has been proven more extensively.

