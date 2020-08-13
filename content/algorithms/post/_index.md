---
title: Proof-of-Spacetime
bookCollapseSection: true
weight: 3
dashboardWeight: 2
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Proof-of-Spacetime
---
_Proof-of-Storage_ schemes allow a user to check if a storage provider is storing the outsourced data at the time
of the challenge. How can we use **PoS** schemes to prove that some data was being stored throughout a period
of time? A natural answer to this question is to require the user to repeatedly (e.g. every minute) send
challenges to the storage provider. However, the communication complexity required in each interaction can
be the bottleneck in systems such as Filecoin, where storage providers are required to submit their proofs to
the blockchain network.

To address this question, we introduce a new proof, _Proof-of-Spacetime_, where a verifier can check if a prover
is storing her/his outsourced data for a range of time. The intuition is to require the prover to (1) generate
sequential Proofs-of-Storage (in our case Proof-of-Replication), as a way to determine time (2) recursively
compose the executions to generate a short proof.

Section 3.3 of the [Filecoin Paper](https://filecoin.io/filecoin.pdf) provides the original introduction to Proof-of-Spacetime.
