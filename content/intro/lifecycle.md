---
title: Lifecycle of Data
weight: 4
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Lifecycle of Data in Filecoin

## Preparing for storing data on Filecoin

1. Data stored on the Filecoin network is split in chunks, similarly to what is happening when adding files to IPFS. Once a file is split into chunks, the IPLD DAG is created, which has its own root CID, called *payload CID* in Filecoin.
2. After some conversions (described in more detail in the corresponding [section](piece)) the DAG is converted to what constitutes the *Filecoin Piece*, which is the main unit of negotiation from the client's point of view. In other words, the *piece* is what the client is negotiating a deal for later on in the process. Every piece has its own CID, called the *piece CID*, or *CommP* (Piece Commitment).

## Finding a deal and publishing the deal on the Filecoin blockchain.

1. A user that wants to store data in the Filecoin network needs to agree a deal with the network/a miner. Finding and agreeing on a deal is done out of band, i.e., without involving the blockchain.

2. The deal is negotiated and agreed between the user and a miner (out-of-band). The *piece CID* is augmented with deal details to create the *Deal Proposal* that has its own *deal CID*. Among other things, the *deal CID* contains the identities of both the client and the miner, as well as transaction details. Data is sent from the client to the miner.

3. Once the data transfer is complete (and the miner verifies that the data is the same as the ones negotiated), the deal is published on the blockchain and the system enters the phase of Storage Mining. Together with publishing the deal on the blockchain, the miner is submitting Initial Pledge collateral, so that it can be held accountable if it fails to prove storage later in the process.

4. The miner places incoming data in a Sector and starts the Sealing process. Once the sealing process is complete the deal is announced/recorded to the Storage Market Actor (through the `ProveCommitSector` function). At this point the miner also puts pledge collateral proportional to the amount of storage they have committed on chain. This triggers the Proof of Replication algorithm. Note that both a sealed copy of the data is _unique to all of the following at the same time_: i) the specific data, ii) the miner storing the data, but also iii) the time at which the miner has sealed the data and produced the _payload CIDs_. In other words, if a miner attempts to re-seal the same data that it has sealed before it will produce a different _payload CID_.

5. **Proof-of-Replication** is the process by which the miner proves that they store a unique copy of the data. This proof is taking place only once when the miner first receives and stores the data. The proof consists of the "Sealed Sector CID" (or CommR, Commitment of Replication?), which is unique both to the miner and the sealing process itself. That is, if the same miner seals the same data again, they will produce a different "Sealed Sector CID".

6. From that point on, the miner is expected to keep the data stored for as long as the deal specifies. The system checks that the miner is indeed storing the sealed data through **Proofs of SpaceTime**. The system is continuously (i.e., not only once as in PoRep) asking the miner to solve puzzles based on some random part of the data. This process proves that the miner indeed stores a unique copy of the data. If the miner fails to provide the Proof-of-SpaceTime, the miner loses their collateral, which they have submitted to the system earlier.

7. Throughout the process, miners might run into Faults, which keeps them back from providing the data to the client, but also from proving that they store the data (provifing Proofs of SpaceTime). Depending on when in the process a miner reports a fault, they receive the corresponding penalty: a miner that reports a fault immediately is penalized less than a miner who fails to provide a proof, when asked.
