---
title: Proof-of-Replication
bookCollapseSection: true
weight: 2
bookhidden: true
---

# Proof-of-Replication

_Proof-of-Replication(PoRep)_, is a new kind of _Proof-of-Storage_, that can be used to prove that some data _D_ has been replicated to its own uniquely dedicated physical storage. Enforcing unique physical copies enables a verifier to check that a prover is not deduplicating multiple copies of _D_ into the same storage space. This construction is particularly useful in Cloud Computing and Decentralized Storage Networks, which must be transparently verifiable, resistant to Sybil attacks, and unfriendly to outsourcing.

Section 3.2 of the [Filecoin Paper](https://filecoin.io/filecoin.pdf) provides the original introduction to Proof-of-Replication.
