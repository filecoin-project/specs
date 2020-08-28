---
title: Virtual Machine
weight: 1

dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Virtual Machine
---

The *Filecoin Virtual Machine* executes messages that are ordered on the blockchain and mantains a single global state that is updated at each epoch.
The are two primitives in the VM runtime: *Actors* and *Messages*.

## Actors

An *Actor* is an entity that can receive and send *messages*, it has a private state and can alter another actor's state through messaging.
It has a set of callable methods that can lead to state updates or new messages sent to other actors.
The VM has singleton actors (e.g. Cron, Power, Market) and multiple instances of the same actor class (e.g. Account, Miner, MultiSig, PaymentChannel).
An actor is either existing at genesis, or it is create via another actor.
The code that implements actors' methods is built-in in the VM and currently custom actor code is not supported.

### Addresses

An *address* is an identified that refers to an actor in Filecoin.
See the Filecoin Address spec.

| Address prefix | Address type       |
|----------------|--------------------|
| `t0`           | Incremental ID     |
| `t1`           | SECP256K1 key pair |
| `t2`           | Actor Unique ID    |
| `t3`           | BLS key pair       |

#### `t0` addresses

Every actor has a `t0` address.

When new actors are created via the `Init` actor, they are assigned a sequential ID.
The incremental IDs start from `t0100`, since addresses from `t00` to `t099` are reserved.
Some of these reserved addresses are currently used for Genesis actors.

Since the network may have forks, on each fork the `Init` actor may assign a `t0` to a different address.
`t0` addresses should only be used after a large number of epochs has passed.

#### `t1` addresses
These addresses are valid for `Account` actors only.

#### `t2` addresses

Every user-created actor has a `t2` address (e.g. `Account`, `Miner`, `MultiSig`)

When new actors are created via the `Init` actor, they are assigned a unique Identifier

#### `t3` addresses
These addresses are valid for `Account` actors only

### Actor Methods

Actors have methods that can be called by the actor itself or via a message from another actor.
All method has an identifier.
All actors implement two methods

| Method name | Method ID |
|-------------|-----------|
| Send        |       `0` |
| Constructor |       `1` |

### Actor Types

There are four actor types (excluding classes used for singletons): `Account`, `Miner`, `MultiSig` and `PaymentChannel`.

### Genesis Actors

A set of actor is created at the Filecoin Genesis epoch.
All actors are singleton except for `BurntFunds` that is an `Account` actor.

| Actor name         |    Actor address |
|--------------------|------------------|
| `System`           |            `t00` |
| `Init`             |            `t01` |
| `Reward`           |            `t02` |
| `Cron`             |            `t03` |
| `Power`            |            `t04` |
| `Market`           |            `t05` |
| `VerifiedRegistry` |            `t06` |
| `BurntFunds`       |           `t099` |


## Messages

A *Message* is the unit of communication between two actors.
A message specifies the method to be executed by the receiving actor and its necessary inputs.

At each epoch, messages can originate by Filecoin nodes sending *signed messages* or by the `Cron` actor, or by actors that are receiving messages.

Messages that are syntactically correct can be included in the block even if on VM execution they abort or return errors.

### Syntax Validation

A message is *syntactically* valid if:
- `To` field is a non-empty, well-formed `Actor` address
- `From` field is a non-empty, well-formed `Actor` address
- `Value` field is a positive number, but smaller than total FIL supply
- `GasPrice` field is positive number
- `GasLimit` field is nt greater than block's gas limit.

Syntax validation is performed when receiving messages and blocks.

### Semantic Validation
TODO

## Gas

*Gas* is the denomination of a unit of computation performed and storage used by a message execution.

A user pays for the execution of a message. The user sets the *gas limit*, an upperbound of gas that they are willing to cover costs for and the *gas price* in FIL.

A miner chooses which messages to include in a block and is rewarded according to each message’s gas price and consumption, forming a market.

A message that is included in a block pays for gas fees even if the VM execution led to an error or it reached the gas limit.

## Execution Model

Messages in a tipset are ordered and processed synchronously.
An actor waits for a sent message to complete before resuming control.

## Environment

### State Tree

The State Tree is a Merkle tree over the entire Filecoin VM state and it is generated after all the messages for the previous epoch have been processed.

### Randomness

There are two source of randomness in Filecoin: Random Beacon, Tickets.

#### Sources

##### Beacon Randomness
The *Beacon Randomness* is generated by an external *Random Beacon* called drand and emits randomness at each Filecoin epoch, and miners use this randomness in order to mine blocks and must include it in their mined blocks.

The *Beacon Randomness* is used to seed randomness generation for values that need to be unbiasable or unpredictable (e.g. to generate challenges for Proof-of-SpaceTime).


Properties:
- Randomness cannot be known earlier than its release time.
- Randomness cannot be biased.
- Randomness of an epoch is the same across forks.

###### Randomn Beacon Outages

A Random Beacon outage is a period in which the Random Beacon is not available to nodes in the network; it could be caused by software bugs, network partitions and other attacks.

During an outage Filecoin cannot generate new blocks, since the Random Beacon provides the randomness that for generating election proofs. After an outage, the Random Beacon used by Filecoin will go in a catch-up mode.

During catchu-up, the Random Beacon will restart to emit the randomness from the last emitted round and it will emit it the randomness at a faster pace until it has finally caught up; Filecoin miners will be generating blocks as soon as new randomness is released.

The catch-up round time is chosen such that blocks can be propagated to a large portion of the network and proofs can be generated and successfully submitted to the blockchain.

##### Ticket Randomness
The *Ticket Randomness* is generated by miners when winning blocks by computing a VRF on the previous ticket and the current drand randomness. The ticket of a tipset is the smallest ticket across blocks in that tipset.

The *Ticket Randomness* is used to tie values to a specific fork (e.g. to tie sealed sector to a chain).

Properties:
- Randomness of an epoch is different across forks.
- Randomness released by other miners cannot be known ahead of time. However, miners can predict a possible randomness outcome by selfish mining.
- Randomness could be biased by miners that choose to withhold their own blocks or that choose to exclude parent blocks.


#### Usage

Filecoin uses the randomness sources to generate randomness that is used across the protocol.

| Domain Separation Tag            | Randomness used       |
|----------------------------------|-----------------------|
| `SealRandomness`                 | `Ticket`              |
| `PoStChainCommit`                | `Ticket`              |
| `WindowedPoStDeadlineAssignment` | `Ticket`              |
| `TicketProduction`               | `Ticket` and `Beacon` |
| `ElectionProofProduction`        | `Beacon`              |
| `WinningPoStChallengeSeed`       | `Beacon`              |
| `WindowedPoStChallengeSeed`      | `Beacon`              |
| `InteractiveSealChallengeSeed`   | `Beacon`              |
| `MarketDealCronSeed`             | `Beacon`              |

##### ElectionProofProduction
Input to the VRF computation that generates an election proof.

Seeded with Beacon randomness at the election epoch.

##### InteractiveSealChallengeSeed
Seed for challenge generation for proving and verifying `miner.ProveCommitSector` messages.

Seeded with Beacon randomness at `InteractiveEpoch` epoch.

##### MarketDealCronSeed
Seed to the schedule the deal processing cron event in order to prevent DoS attacks and distribute the `Cron` actor's load.

Seeded with Beacon randomness at the epoch in which deals are published.

##### PoStChainCommit
Reference that ties a `miner.SubmittedWindowPoSt` to a tipset in a fork to avoid replay attacks of proofs submission in long-range or selfish mining attacks.

Seeded with Ticket randomness at the challenge epoch of the current deadline.

##### SealRandomness

Seed to the sealing process that ties a sealed sector to a fork in order to avoid `miner.PreCommitSectors` message in long-range or selfish mining attacks.

Seeded with Ticket randomness at the `SealEpoch` epoch.

##### TicketProduction
Input to the VRF computation to generate a new Ticket.

Seeded with Ticket randomness of the previous epoch and Beacon of the current epoch. TODO: specify what happens with null blocks.

The Beacon randomness guarantees that it cannot be predicted ahead of time. This is relevant for genesis block. TODO: note that this could be simplified by having TicketRandomness to be a Beacon randomness at genesis block.

##### WindowedPoStChallengeSeed
Seed for challenge generation for proving and verifying `miner.SubmittedWindowPoSt` messages.

Seeded with Beacon randomness at the current deadline's challenge epoch.

##### WindowedPoStDeadlineAssignment
TOD

##### WinningPoStChallengeSeed

Seed for challenge generation for proving and verifying `miner.SubmittedWindowPoSt` messages.

Seeded with Beacon randomness at the current election epoch.
