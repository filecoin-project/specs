# PoRep with RSA Vector Commitments

This document describes Proof of Replication using BBF-VC ([Boneh Bunz Fisch Vector Commitments](https://eprint.iacr.org/2018/1188.pdf)).
The motiviation for this is to remove Merkle Trees from our SNARK proofs.

The PoRep replication algorithm would remain the same (except it would substitute Merkle Trees with BBF-VC, which would have the same APIs).
In the rest of the doc, we describe how to use Vector Commitments.

## Overview

The prover will send to a verifier a proof with the three components:

- VC proof: Inclusion proof of the challenged replicated nodes, original data nodes and parent nodes to the accumulator `commAcc`.
- VC openings: Bits revealed from the VC proof used as public inputs to verify the SNARK proof
- SNARK proof: Proves that:
  - Challenged nodes are encoded from data nodes and parent nodes correctly
  - VC Opening Bits are the correct last bits of the data nodes, replica nodes and parent nodes

## Prover Algorithm

We use the following constants:

- `PARENTS_COUNT = BASEDEGREE + EXPANSIONDEGREE`
- `c_t` number of revealed bits for the original data node
- `c_d` number of revealed bits for each parent node

A proof is `(VCProof, (ParentBits, DataBits),  SNARKProof)`

### VC Inclusion proofs (`VCProof`)

- Output `VCProof`
  - Inclusion proof for: the challenged data nodes, the challenged replica nodes and the parents of the challenged replica nodes

### VC Openings (`(ParentBits, DataBits)`)

- For each `challenge` in `0..CHALLENGE_COUNT`:
  - `DataBits[challenge] = Replica[challenge][:c_t]`
  - For each `parent` in `0..PARENTS_COUNT`:
    - `ParentBits[challenge][parent] = Replica[GetParent(challenge)[parent]][:c_d]`
- Output `ParentBits`, `DataBits`

### SNARK circuit (`SNARKProof`)

> TODO: Check if we need `ReplicaBits`

- Private Inputs:
  - `Parents [CHALLENGES_COUNT][PARENTS_COUNT]Fr`
  - `Data [CHALLENGES_COUNT]Fr`
- Public Inputs:
  - `ParentBits [CHALLENGES_COUNT][PARENTS_COUNT][c_d]Bool`
  - `DataBits [CHALLENGES_COUNT][c_t]Bool`
- Computation
  - Run the encoding check (same as in drgporep circuit)
    - Compute `key = KDF(parents, node)`
    - Compute `obtained = VDE.Decode(node, key)`
    - Check that `Data === obtained`
  - Check q bits:
    - For each challenge `ch`:
      - For each parent `parent`:
        - Check that `parent[:c_d] == ParentBits[ch][parent]`
      - Check that `Data[:c_t] == DataBits[ch]`

## Verifier Algorithm

The verifier will perform the following operations:
- Verify that the `VCProof` is correct
- Verify that the (`ParentBits`, `DataBits`) are in the accumulator
- Verify that the `SNARKProof` is correct
