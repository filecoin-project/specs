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
The VM has singleton actors (e.g. Cron, Power, Market) and multiple instances of the same actor class (e.g. Account, Miner).
An actor is either existing at genesis, or it is create via another actor.
The code that implements actors' methods is built-in in the VM and currently custom actor code is not supported.

### Addresses

Each actor has an address.
Address from `t00` to `t099` are reserved to built-in singleton actors, the first non-singleton actor starts from `t0100`.

### Actor Methods

Actors have methods that can be called by the actor itself or via a message from another actor.
All method has an identifier.
All actors implement two methods

| Method name | Method ID |
|-------------|-----------|
| Send        |       `0` |
| Constructor |       `1` |

### Actor Types

There are two actor types (excluding classes used for singletons): `Account` and `Miner`.

### Genesis Actors

A set of actor is created at genesis.
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
