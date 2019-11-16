---
title: Rational-PoSt
---

This document describes Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

# High Level API

## Fault Detection

Fault detection happens over the course of the life time of a sector. When the sector is for some reason unavailable, the miner is responsible to submit the known `faults`, before the PoSt challenge begins. (Using the `AddFaults` message to the chain).
Only faults which have been reported at challenge time, will be accounted for. If any other faults have occured the miner can not submit a valid PoSt for this proving period.

The PoSt generation then takes the latest available `faults` of the miner to generate a PoSt matching the committed sectors and faults.

When a PoSt is successfully submitted all faults are reset and assumed to be recovered. A miner must either (1) resolve a faulty sector and accept challenges against it in the next proof submission, (2) report a sector faulty again if it persists but is eventually recoverable, (3) report a sector faulty *and done* if the fault cannot be recovered.

If the miner knows that the sectors are permanently lost, they can submit them as part of the `doneSet`, to ensure they are removed from the proving set.

{{% notice note %}}
**Note**: It is important that all faults are known (i.e submitted to the chain) prior to challenge generation, because otherwise it would be possible to know the challenge set, before the actual challenge time. This would allow a miner to report only faults on challenged sectors, with a gurantee that other faulty sectors would not be detected.
{{% /notice %}}


{{% notice todo %}}
**TODO**: The penalization for faults is not clear yet.
{{% /notice %}}

## Fault Penalization

Each reported fault carries a penality with it.

{{% notice todo %}}
**TODO**: Define the exact penality structure for this.
{{% /notice %}}

## Generation

`GeneratePoSt` generates a __*Proof of Spacetime*__ over all  __*sealed sectors*__ of a single miner— identified by their `commR` commitments. This is accomplished by performing a series of merkle inclusion proofs (__*Proofs of Retrievability*__). Each proof is of a challenged node in a challenged sector. The challenges are generated pseudo-randomly, based on the provided `seed`. At each time step, a number of __*Proofs of Retrievability*__ are performed.

```go
// Generate a new PoSt.
func GeneratePoSt(sectorSize BytesAmount, sectors SectorSet, seed Seed, faults FaultSet) PoStProof {
    // Generate the Merkle Inclusion Proofs + Faults

    challenges := DerivePoStChallenges(seed, faults, sectorSize, SortAsc(GetSectorIds(sectors)))
    challengedSectors := []
    inclusionProofs := []

    for i := 0; i < len(challenges); i++ {
        challenge := challenges[i]

        // Leaf index of the selected sector
        inclusionProof, isFault := GenerateMerkleInclusionProof(challenge.Sector, challenge.Leaf)
        if isFault {
            // faulty sector, need to post a fault to the chain and try to recover from it
            return Fatal("Detected late fault")
        }

        inclusionProofs[i] = inclusionProof
        challengedSectors[i] = sectors[challenge.Sector]
    }

    // Generate the snark
    snarkProof := GeneratePoStSnark(sectorSize, challenges, challengedSectors, inclusionProofs)

    return snarkProof
}
```

## Verification

`VerifyPoSt` is the functional counterpart to `GeneratePoSt`. It takes all of `GeneratePoSt`'s output, along with those of `GeneratePost`'s inputs required to identify the claimed proof. All inputs are required because verification requires sufficient context to determine not only that a proof is valid but also that the proof indeed corresponds to what it purports to prove.

```go
// Verify a PoSt.
func VerifyPoSt(sectorSize BytesAmount, sectors SectorSet, seed Seed, proof PoStProof, faults FaultSet) bool {
    challenges := DerivePoStChallenges(seed, faults, sectorSize, SortAsc(GetSectorIds(sectors)))
    challengedSectors := []

    // Match up commitments with challenges
    for i := 0; i < len(challenges); i++ {
        challengedSectors[i] = sectors[challenges[i].Sector]
    }

    // Verify snark
    return VerifyPoStSnark(sectorSize, challenges, challengedSectors)
}
```

## Types

```go
// The random challenge seed, provided by the chain.
Seed [32]byte
```

```go
type Challenge struct {
    Sector SectorID
    Leaf Uint
}
```

## Challenge Derivation

```go
// Derive the full set of challenges for PoSt.
func DerivePoStChallenges(seed Seed, faults FaultSet, sectorSize Uint, sortedSectors []SectorID) [POST_CHALLENGES_COUNT]Challenge {
    challenges := []

    for n := 0; n < POST_CHALLENGES_COUNT; n++ {
        attemptedSectors := {SectorID:bool}
        attempt := 0
        while challenges[n] == nil {
            challenge := DerivePoStChallenge(seed, n, attempt, sectorSize, sortedSectors)
            attempt ++
            
            // check if we landed in a faulty sector
            if !faults.Contains(challenge.Sector) {
                // Valid challenge
                challenges[n] = challenge
            }

            // invalid challenge, regenerate
            attemptedSectors[challenge.Sector] = true

            if len(attemptedSectors) >= len(sortedSectors) {
                Fatal("All sectors are faulty")
            }
        }
    }

    return challenges
}

// Derive a single challenge for PoSt.
func DerivePoStChallenge(seed Seed, n Uint, attempt Uint, sectorSize Uint, sortedSectors []SectorID) Challenge {
    nBytes := WriteUintToLittleEndian(n)
    data := concat(seed, nBytes, WriteUintToLittleEndian(attempt))
    challengeBytes := blake2b(data)

    sectorChallenge := ReadUintLittleEndian(challengeBytes[0..8])
    leafChallenge := ReadUintLittleEndian(challengeBytes[8..16])

    sectorIdx := sectorChallenge % len(sortedSector

    return Challenge {
        Sector: sortedSectors[sectorIdx],
        Leaf: leafChallenge % (sectorSize / NODE_SIZE),
    }
}
```


# PoSt Circuit

## Public Parameters

*Parameters that are embeded in the circuits or used to generate the circuit*

- `POST_CHALLENGES_COUNT: UInt`: Number of challenges.
- `POST_TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is `(log_2(Size of original data in bytes/32 bytes per leaf))`.
- `SECTOR_SIZE: UInt`: The size of a single sector in bytes.

## Public Inputs

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `CommRs: [POST_CHALLENGES_COUNT]Fr`: The Merkle tree root hashes of all replicas, ordered to match the inclusion paths and challenge order.
- `InclusionPaths: [POST_CHALLENGES_COUNT]Fr`: Inclusion paths for the replica leafs, ordered to match the `CommRs` and challenge order. (Binary packed bools)

## Private Inputs

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `InclusionProofs: [POST_CHALLENGES_COUNT][TREE_DEPTH]Fr`: Merkle tree inclusion proofs, ordered to match the challenge order.
- `InclusionValues: [POST_CHALLENGES_COUNT]Fr`: Value of the encoded leaves for each challenge, ordered to match challenge order.


## Circuit

### High Level

In high level, we do 1 check:

1. **Inclusion Proofs Checks**: Check the inclusion proofs

### Details

```go
for c in range POST_CHALLENGES_COUNT {
  // Inclusion Proofs Checks
  assert(MerkleTreeVerify(CommRs[c], InclusionPath[c], InclusionProof[c], InclusionValue[c]))
}
```

## Verification of PoSt proof

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
