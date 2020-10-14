---
title: "Verifiable Random Function"
weight: 2
dashboardWeight: 2
dashboardState: incorrect
dashboardAudit: n/a
dashboardTests: 0
---

# Verifiable Random Function

Filecoin uses the notion of a [Verifiable Random
Function](https://people.csail.mit.edu/silvio/Selected%20Scientific%20Papers/Pseudo%20Randomness/Verifiable_Random_Functions.pdf)
(VRF). A VRF uses a private key to produce a digest of
an arbitrary message such that the output is unique per signer and per message.
Any third party in possession of the corresponding public key, the message, and
the VRF output, can verify if the digest has been computed correctly and by the
correct signer. Using a VRF in the ticket generation process allows anyone to
verify if a block comes from an eligible block producer (see [Ticket Generation](storage_power_consensus#tickets) for more details).

BLS signature can be used as the basis to construct a VRF. Filecoin transforms
the BLS signature scheme it uses (see [Signatures](signatures) into a
VRF, Filecoin uses the random oracle model and deterministically hashes the
signature (using blake2b to produce a 256 bit output) to produce the final digest.

These digests are often used as entropy for randomness in the protocol (see [Randomness](randomness)).

{{<embed src="vrf.id" lang="go">}}
{{<embed src="vrf.go" lang="go">}}
