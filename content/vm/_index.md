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

A *Message* is always sent from and to actors.
A message specifies the method to be executed by the receiving actor and its necessary inputs.
Actors only send messages to others if they received a message in the current epoch.
At each epoch messages can originate by Filecoin nodes sending *signed messages* or by the `Cron` actor.

## Execution Model
