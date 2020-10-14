---
title: Proof-of-Storage
bookCollapseSection: true
weight: 2
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Proof-of-Storage

## Preliminaries

Storage miners in the Filecoin network have to prove that they hold a copy of the data at any given point in time. This is realised through the [Storage Miner Actor](storage_mining#storage_miner_actor) who is the main player in the [Storage Mining](filecoin_mining#storage_mining) subsystem. The proof that a storage miner indeed keeps a copy of the data they have promised to store is achieved through "challenges", that is, by providing answers to specific questions posed by the system. In order for the system to be able to prove that a challenge indeed proves that the miner stores the data, the challenge has to: i) target a random part of the data and ii) be requested at a time interval such that it is not possible, profitable, or rational for the miner to discard the copy of data and refetch it when challenged.

General _Proof-of-Storage_ schemes allow a user to check if a storage provider is storing the outsourced data at the time
of the challenge. How can we use **PoS** schemes to prove that some data was being stored throughout a period
of time? A natural answer to this question is to require the user to repeatedly (e.g. every minute) send
challenges to the storage provider. However, the communication complexity required in each interaction can
be the bottleneck in systems such as Filecoin, where storage providers are required to submit their proofs to
the blockchain network.

To address this question, we introduce a new proof, called _Proof-of-Spacetime_, where a verifier can check if a prover
has indeed stored the outsourced data they committed to over (storage) Space and over a period of Time (hence, the name Spacetime). 

Recall that in the Filecoin network, miners are storing data in fixed-size [sectors](filecoin_mining#sector). Sectors are filled with client data agreed through regular deals in the [Storage Market](filecoin_markets#storage_market), through [verified deals](algorithms#verified_clients), or with random client data in case of [Committed Capacity sectors](filecoin_mining#sector).
