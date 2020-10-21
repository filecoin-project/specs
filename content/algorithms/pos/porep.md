---
title: PoRep
weight: 1
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Proof-of-Replication (PoRep)

In order to register a sector with the Filecoin network, the sector has to be sealed. Sealing is a computation-heavy process that produces a unique representation of the data in the form of a proof, called **_Proof-of-Replication_** or PoRep.

The PoRep proof ties together: i) the data itself, ii) the miner actor that performs the sealing and iii) the time when the specific data has been sealed by the specific miner. In other words, if the same miner attempts to seal the same data at a later time, then this will result in a different PoRep proof. Time is included as the blockchain height when sealing took place and the corresponding chain reference is called `SealRandomness`.

Once the proof has been generated, the miner runs a SNARK on the proof in order to compress it and submits the result to the blockchain. This constitutes a certification that the miner has indeed replicated a copy of the data they agreed to store.

The PoRep process includes the following two phases:

- **Sealing preCommit phase 1.** In this phase, PoRep SDR [encoding](sdr#encoding) and [replication](sdr#replication) takes place.
- **Sealing preCommit phase 2.** In this phase, [Merkle proof and tree generation](sdr#merkle-proofs) is performed using the Poseidon hashing algorithm.
