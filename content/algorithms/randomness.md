---
title: "Randomness in FIL"
weight: 5
dashboardWeight: 1
dashboardState: stable
dashboardAudit: 0
dashboardTests: 0
---

# Randomness in Filecoin
---

Randomness is an important part of the healthy functioning of the Filecoin network. Randomness is primarily used in Leader Election of the Filecoin Expected Consensus algorithm, in order to choose the winning miner(s) of the next epoch. In order for the randomness to be unbiasable it has to be external to the system. Filecoin is using [DRAND](https://drand.love/), a verifiable, unpredictable and unbiased randomness beacon produced by a distributed network of servers operated by several independent organisations.

The design of DRAND is such that none of the organisations participating in the DRAND network can predict future beacons neither on its own, nor if they team up with other partners of the network. In order to influence randomness in DRAND, all partners would have to team up, a situation that is impossible to happen due to the vastly divergent interests of the different partners and the applications that are running on top. In case of outage of some of the partners, the DRAND network is still operational, as long as 70% of DRAND members are functional.

The DRAND network is producing a new randomness beacon every 30 sec.

For more details on DRAND, please refer to its [documentation](https://drand.love/docs/), [developer docs](https://drand.love/developer/) and [specification](https://drand.love/docs/specification/).

## GetRandomness(personalization, epoch, entropy)

The `GetRandomness` function is the main function used by miners to "stamp" a newly produced block with a unique identifier that links the block to the specific epoch when it was created. The function is also associating the identifier of the miner that produced the block.

The `GetRandomness` function should always return the same result, regardless of the content of the current, previous or future blocks. In other words, the function is independent of the block itself. Instead, the function is primarily dependent on the epoch given to it as input.

The outcome of `GetRandomness` is a hash that is unique to the miner and the epoch when the random value was obtained from DRAND and attached to the newly minted block.

If the epoch is in the future or DRAND is catching up after an outage, `GetRandomness` waits until the DRAND network outputs the value for the right epoch.


```text
Algorithm

##
- randSeed = chain.GetBeaconEntryForEpoch(epoch)
##
- H(personalization || randSeed || epoch || entropy)
```

## GetBeaconEntryForEpoch(epoch)

The `GetBeaconEntryForEpoch` function returns a _single beacon entry_ given an epoch provided as input. It is used by the `GetRandomness` function.

If there is no entry for this epoch, `GetBeaconEntryForEpoch` returns the latest entry for the previous epoch.

## GetBeaconEntriesForEpoch(epoch)

The `GetBeaconEntriesForEpoch` is a deterministic function that returns all the beacon entries expected for this epoch, i.e., all the randomness generated from the last epoch until now. This applies in case the epoch duration is larger than the randomness generation interval. In case of Filecoin and according to current settings the DRAND randomness generation interaval is equal to the Filecoin epoch - both set at 30 sec.

```text
Algorithm:

- Compute the drand round for this epoch `maxDrandRound`
- Compute the drand round for the previous epoch `prevMaxDrandRound`
- Get all entries in between: `drand.Public(prevMaxDrandRound+1)`, ..., `drand.Public(maxDrandRound)`

Edge cases:

- If `maxDrandRound == prevMaxDrandRound`, then return nil
```

### Explanations

#### _Relationship between `maxDrandRound` and `prevMaxDrandRound`_

The values of `maxDrandRound` and `prevMaxDrandRound` can be equal in case Drand did not produce any value in two subsequent Filecoin epochs, e.g., when a `DrandRound` has ended up being 61 seconds, while the default Filecoin epoch is 30 seconds.

#### _Randomness during the first epoch / Genesis Block_

The following is a sequence of things for the first few epochs after the Genesis Block is generated.

- There is no DRAND value in the genesis block. Randomness information is included in the genesis block, but is not drawn from DRAND.
- The first epoch block contains one _partial_ DRAND entry which is _not_ verifiable since the entry does not contain the previous signature.
- At the second epoch, the "main chain" 2-nd epoch block contains a DRAND entry that should point to the first DRAND entry.
- Miners will build on the block(s) produced during the first epoch without knowing if the randomness is valid or not. However, on the second epoch, miners will be able to verify if the randomness included in the first epoch is correct.

#### _DRAND Outage_

A DRAND outage is causing severe problems to the Filecoin blockchain, which cannot continue its normal operation, as miners cannot associate the creation time (i.e., epoch) of their blocks with the randomness value produced during that epoch.

- _When DRAND is down_ miners are not able to create blocks for the epochs for which a DRAND entry was not generated. Note that `GetBeaconEntryForEpoch` calls `drand.Public(..)` which will waits for DRAND to create a new random value. The miner should insist on fetching the same randomness (i.e., the randomness for the specific epoch), until it succeeds.
- _When DRAND catches up_ miners are able to re-create blocks in the past and the heaviest chain will eventually be chosen.
- _During DRAND catchup_ both DRAND and the Filecoin blockchain are trying to catch up and are therefore moving at a faster pace than in normal conditions. The _catch up_ mode continues until the DRAND round is the same as the one that would exist had the outage not happened. The DRAND round (i.e., the frequency at which a new randomness value is generated) and the Filecoin epoch are reduced to 15 sec during catch up mode, from 30 sec during normal operation. This has been decided in order to avoid WindowPoSt failures and Fake power generation (PoRep) which could be triggered by malicious miners after the DRAND outage finishes if the catch up takes too long.


## MaxBeaconRoundForEpoch(epoch)

The `MaxBeaconRoundForEpoch` is a deterministic function that returns the round at which DRAND is at (i.e., the latest round). This is the _maximum_ round meaning that in case the DRAND network is catching up, this will be the latest round up to which it has to catch up. 

```text
Algorithm:

- Compute UNIX timestamp of Filecoin for this `epoch`
- Compute the Drand round
```

### Explanations

#### _Filecoin Timestamp_

The Filecoin Timestamp in the first few epochs is generated as follows:

- Epoch 0 blocks has `Filecoin Genesis` timestamp
- Epoch 1 block has `Filecoin Genesis + Filecoin Epoch` timestamp
- Epoch 2 is `Filecoin Genesis time + 2*Filecoin Epoch` timestamp
- Epoch 3 is `Filecoin Genesis time + 3*Filecoin Epoch` timestamp

#### _Filecoin Timestamp Calculation_

When computing the Filecoin timestamp one epoch has to be subtracted according to:

```go
latestTs = ((uint64(filEpoch) * filEpochDuration) + filGenesisTime) - filEpochDuration
```
This is because Filecoin takes the DRAND value that corresponds to the previous Filecoin epoch of the one when the block is created. That allows the miner to have the DRAND value when it creates the block which is necessary for leader election.

#### _Update of `MaxBeaconRoundForEpoch` relative to DRAND round and Filecoin epoch_

The `MaxBeaconRoundForEpoch` function might return no updated values in case the DRAND round duration is longer than the Filecoin epoch. For instance, in case DRAND round is 60 sec and Filecoin epoch is 30 sec, then for two subsequent epochs the function will return the same value.

On the other hand, the `MaxBeaconRoundForEpoch` function might return multiple values in one Filecoin epoch in case the Filecoin epoch is longer than the DRAND round. For instance, in case the Filecoin epoch is 60 sec and DRAND round is 30 sec, then the function will return two values within the same Filecoin epoch.

NOTE: None of the above two cases is taking place according to current developments, as both the Filecoin epoch and the DRAND round are set to 30 sec. These values are unlikely to change.


### Note on discrepancy of DRAND round

The DRAND round value is permanently off by one, that is, for a given timestamp, the `MaxBeaconRoundForEpoch` function should return `current round +1`, while the actual computation should return `current round`:

```go
fromGenesis := now - genesis
dround = uint64(math.Floor(float64(fromGenesis)/period.Seconds())) + 1
```

The reason of the extra (`+1`) round in the DRAND round calculation is because the first DRAND value is pulled at the first block, which is `+1` block after the genesis block.
