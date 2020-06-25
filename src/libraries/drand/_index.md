---
menuTitle: drand
title: drand - Distributed Randomness
---

DRand (Distributed Randomness) is a publicly verifiable random beacon protocol Filecoin relies on as a source of unbiasable entropy for leader election (see {{<sref leader_election>}}).

At a high-level, the drand protocol runs a series of MPCs (Multi-Party Computations) in order to produce a series of deterministic, verifiable random values. Specifically, after a trusted setup, a known (to each other) group of n drand nodes sign a given message using t-of-n threshold BLS signatures in a series of successive rounds occuring at regular intervals (the drand round time).
Any node that has gathered t of the signatures can reconstruct the full BLS signature. This signature can then be hashed in order to produce a collective random value which can be verified against the collective public key generated during the trusted setup. Note that while this can be done by the drand node, the random value (i.e. hashed value) should be checked by the consumer of the beacon. In Filecoin, we hash it using blake2b in order to obtain a 256 bit output.

drand assumes that at least t of the n nodes are honest (and online -- for liveness). If this threshold is broken, the adversary can permanently halt randomness production but cannot otherwise bias the randomness.

You can learn more about how drand works, by visiting its [repository](https://github.com/drand/drand), or reading its [spec](https://github.com/drand/drand/blob/master/docs/SPECS.md).

In the following sections we look in turn at how the Filecoin protocol makes use of drand randomness, and at some of the characteristics of the specific drand network Filecoin uses.

### Drand randomness outputs

By polling the appropriate endpoint (see below for specifics on the drand network Filecoin uses), a Filecoin node will get back a drand value formatted as follows (e.g.):

```
{
  "round": 367,
  "signature": "b62dd642e939191af1f9e15bef0f0b0e9562a5f570a12a231864afe468377e2a6424a92ccfc34ef1471cbd58c37c6b020cf75ce9446d2aa1252a090250b2b1441f8a2a0d22208dcc09332eaa0143c4a508be13de63978dbed273e3b9813130d5",
  "previous_signature": "afc545efb57f591dbdf833c339b3369f569566a93e49578db46b6586299422483b7a2d595814046e2847494b401650a0050981e716e531b6f4b620909c2bf1476fd82cf788a110becbc77e55746a7cccd47fb171e8ae2eea2a22fcc6a512486d"
}
```

Specifically, we have:

- `Signature`           -- the threshold BLS signature on the previous signature value `Previous` and the current round number `round`.
- `PreviousSignature`   -- the threshold BLS signature from the previous drand round.
- `Round`               -- the index of Randomness in the sequence of all random values produced by this drand network.

Specifically, the message signed is the concatenation of the round number treated as a uint64 and the previous signature. At the moment, drand uses BLS signatures on the BLS12-381 curve with the latest v7 RFC of hash-to-curve and the signature is made over G1 (for more see the [drand spec](https://github.com/drand/drand/blob/master/docs/SPECS.md#cryptographic-specification)).

### Polling the drand network

Filecoin nodes can make use of [drand endpoints](https://github.com/drand/drand/blob/master/client/client.go) in working with a drand beacon.

To start, a node must store a set of drand peer servers it will connect to to poll for values and the shared public key it expects them to have. In the case, of Filecoin, the node will be polling "relay nodes" rather than drand nodes directly. Simply put, the drand network will not be directly accessible by consumers, rather, highly-available relays will be set up to serve drand values over HTTP or gossipsub. See below section for more on the drand network configuration.

On initialization, the Filecoin node can call the `Group` endpoint in order to obtain the beacon's [group file](https://github.com/drand/drand/blob/57a6056a24d4ebaa27a44852636807364624b9fc/key/group.go). The client (node) should then have the hash of the group file (or the full file itself if it already has it) cached to verify the group it receives corresponds to the expected one. It should also cache:

- cache the beacon's `Period`                           -- the period of time between each drand randomness generation
- cache the beacon's `GenesisTime`                      -- at which the first round in the drand randomness chain is created
- verify that the beacon's `PublicKey` is appropriate   -- ie that the filecoin node connected to the right drand beacon

Thereafter, the Filecoin client can call drand's endpoints:

- `LastPublic` to get the latest randomness value produced by the beacon
- `Public` to get the randoomness value produced by the beacon at a given index

{{<label drand>}}
### Using drand in Filecoin

Drand is used as a randomness beacon for leader election in Filecoin. You can read more about that in {{<sref leader_election>}}. See drand used in the Filecoin lotus implementation [here](https://github.com/filecoin-project/lotus/blob/master/chain/beacon/drand/drand.go).

While drand returns multiple values with every call to the beacon (see above), Filecoin blocks need only store a subset of these in order to track a full drand chain. This information can then be mixed with on-chain data for use in Filecoin. See {{<sref randomness>}} for more.

#### Verifying an incoming drand value

Upon receiving a new drand randomness value from a beacon, a Filecoin node should immediately verify its validity. That is, it should verify:

- that the `Signature` field is verified by the beacon's `PublicKey` as the beacon's signature of `SHA256(PreviousSignature || Round)`.
- that the `Randomness` field is `SHA256(Signature)`.

See [drand](https://github.com/drand/drand/blob/master/beacon/beacon.go#L63) for an example.

#### Fetching the appropriate drand value while mining

There is a deterministic mapping between a needed drand round number and a Filecoin epoch number.

After initializing access to a drand beacon, a Filecoin node should have access to the following values:

- `filEpochDuration`    -- the Filecoin network's epoch duration (between any two leader elections)
- `filGenesisTime`      -- the Filecoin genesis timestamp
- `filEpoch`            -- the current Filecoin epoch
- `drandGenesisTime`    -- drand's genesis timestamp
- `drandPeriod`         -- drand's epoch duration (between any two randomness creations)

Using the above, a Filecoin node can determine the appropriate drand round value to be used for use in {{<sref leader_election>}} in an epoch using both networks' reliance on real time as follows:

```go
MaxBeaconRoundForEpoch(filEpoch) {
    // determine what the latest Filecoin timestamp was from the current epoch number
    var latestTs
    if filEpoch == 0 {
        latestTs = filGenesisTime
    } else {
        latestTs = ((uint64(filEpoch) * filEpochDuration) + filGenesisTime) - filEpochDuration
    }
    // determine the drand round number corresponding to this timestamp
    // keeping in mind that drand emits round 1 at the drandGenesisTime
    dround := (latestTs - drandGenesisTime) / uint64(drandPeriod) + 1
    return dround
}
```

### Edge cases and dealing with a drand outage

It is important to note that any drand beacon outage will effectively halt Filecoin block production. Given that new randomness is not produced, Filecoin miners cannot generate new blocks.

After a beacon downtime, drand nodes will work to quickly catch up to the current round, as defined by wall clock time. In this way, the above time-to-round mapping in drand (see above) used by Filecoin remains an invariant after this catch-up following downtime.

So while Filecoin miners were not able to mine during the drand outage, they will quickly be able to run leader election thereafter, given a rapid production of drand values. We call this a "catch up" period.

During the catch up period, Filecoin nodes will backdate their blocks in order to continue using the same time-to-round mapping to determine which drand round should be integrated according to the time. Miners can then choose to publish their null blocks for the outage period (including the appropriate drand entries throughout the blocks, per the time-to-round mapping), or (as is more likely) try to craft valid blocks that might have been created during the outage.

Note that based on the level of decentralization of the Filecoin network, we expect to see varying levels of miner collaboration during this period. This is because there are two incentives at play: trying to mine valid blocks from during the outage to collect block rewards, not falling behind a heavier chain being mined by a majority of miners that may or may not have ignored a portion of these blocks.

In any event, a heavier chain will emerge after the catch up period and mining can resume as normal.

### drand network specification

TODO once ready: @nikkolasg
- Filecoin node access to randomness (how to connect, poll, etc)
- Drand Node composition and governance
- Network characteristics (DDoS resistance, layer description etc)
