# Solving SEALSTACK

**The SEALSTACK attack**: A prover commits to replicate D into a replica R, and then commit to replicate R into a replica R'. The prover performs the second replication in-place, reusing R space for R'. During Proof of Spacetime, when the prover is challenged to reveal nodes in R, they efficiently decode from R' the challenged nodes.

## Status quo: Symmetric PoRep

This document is a proposal to mitigate the SEALSTACK attack by using a symmetric PoRep. Spec changes from this document should be propagated.

### Intuition

We mitigate SEALSTACK by making sealing and unsealing take the same amount of time. In this way, if the miner has stacked multiple replicas, it would take them too long to decode a layer to reply to the Proof of Spacetime challenges.

### Construction

**Replication**:

- Run a Proof of Space, PoS (ZigZag would work too) with input `replica_id`
- Take the Merkle tree hash of the PoS output (the PoS output is the size of a sector). We call the root hash `CommPoS`.
- XOR the data to seal with the PoS output
- Take the Merkle tree hash of the result, this is `CommR`

**SEAL Proof**:

- Prove that the challenged node in `CommR` is the XOR between data in `CommD` and in `CommPoS`
- Prove that the PoS initialization for the challenged node in `CommPoS` is correct

## Proposal 1: Asymmetric PoRep with phases

This document is a proposal to mitigate the SEALSTACK attack, yet keeping the PoRep asymmetric. This proposal is a protocol change.

Miners should be able to retrieve data from the sealed sectors very efficiently. To do so, the PoRep construction must be asymmetric: slow to decode, fast to encode. In order to do so, we must resolve SEALSTACK.

### Intuition

We mitigate SEALSTACK by introducing two phases: commitment phase, proving phase.

- Commitment phase: during this phase, every miner publishes a commitment to store multiple datasets.
- Proving phase: at the beginning of the phase, a random number is released (miners can't have access to it earlier) and must use this number to encode the committed data. During this period commitments to store more data are invalid.

**Why does this work?** This guarantees that miners cannot commit to store replicas (and hence SEALSTACK) since they would not know what the replica would look like until the proving phase, but at that point, they can't commit to store it.

**Cons?**

- The network is forced to operate in phases

## Proposal 2: Slow Decoding PoRep to solve SEALSTACK
A prover generates a response to a challenge in `k` steps and must reply before `T` time steps.

Assume a prover that is performing a 2-layers SEALSTACK attack. If the prover is challenged on the first layer, then they must decode the second layer which takes `t` time. A PoRep is secure against `n`-SEALSTACK if a malicious miner proof algorithm takes `t+k*n > T` steps.

In order to make PoRep secure, we can increase `t` by making honest proving taking a longer time (e.g. increase sequential steps during encoding) or by increasing the minimum decoding time (e.g. increase sequential steps during decoding).

Note: ZigZag is a favorable construction since in order to decode the graph for a single challenge is equivalent to decode the graph for multiple challenges. This means that a prover can only reply challenges for one layer per time, otherwise the prover would be using more storage than they have (so, not performing a SEALSTACK)

### Questions to answer
- What is `k` concretely?
- Calculate what is `k` in ZigZag (assuming bounded parallelization)
