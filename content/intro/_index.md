---
title: Introduction
weight: 1

dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Introduction
---

This is the spec for the implementation of the core contracts (called *actors*) Filecoin Actors that are implemented in the Filecoin Virtual Machine.

## What is Filecoin?

Filecoin is a *decentralized storage network*, a network of independent storage providers offering storage and retrieval services in a market operated on a blockchain with a native protocol token called FIL.

### Filecoin Market
The *Filecoin Market* is an algorithmic market for storage and retrieval services.
Miners offer their storage capacity in the market and make *storage deals* with clients.
The market is verifiable: storage providers must provide cryptographic proofs that guarantee persistent storage to their clients.

### Filecoin Blockchain
The *Filecoin Blockchain* is a distributed ledger that orders FIL transactions and executes the *Filecoin Virtual Machine*.
Miners mantain the blockchain by creating blocks and verifying transactions and earn block rewards.
Differently from other protocols based on computational resources, miners' influence in the network and block reward earnings are proportional to the amount of storage they prove.

### Filecoin Virtual Machine
The *Filecoin Virtual Machine* is a state machine that implements core functionalities to operate the Filecoin markets and token transactions.

