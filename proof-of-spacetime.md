## Proof-of-Spacetime

This document descibes Rational-PoSt, the Proof of Spacetime used in Filecoin.


## Rational PoSt

### Definitions

- **PoSt Proving Period**: The time interval in which a PoSt has to be submitted.
- **POST_CHALLENGE_TIME**: The time offset at which the actual work of generating the PoSt should be started. This is some delta befre the end of the `Proving Period`, and as such less then a single `Proving Period`.

- start = height x
- challenge time = height x + y
- on post submission: verify that challenge comes from block(x + y)

### Execution Flow

```
    ■──────────────────────── Proving Period ─────────────────────────■


    ▲───────────────────────────────────────────────────▲─────────────▲
    │                                                   │             │
    │                                                   │             │
    │                                                   │             │
┌───────┐                                     ┌──────────────────┐┌───────┐
│ Start │                                     │ Generation Start ││  End  │
└───────┘                                     └──────────────────┘└───────┘
                                                        ▲

                                                        │
                                                                       ┌ ─ ─ ─ ─ ─ ─ ─ ─ ┐
                                                        └ ─ ─ ─ ─ ─ ─ ─ Randomness Input
                                                                       └ ─ ─ ─ ─ ─ ─ ─ ─ ┘
```

TODO: Add post submission to the diagram.




### High Level API

#### Generation

```go
func GeneratePoSt(sectorSize BytesAmount, sectors []commR) (PoStProof, FaultSet) {
    // Generate the Merkle Inclusion Proofs + Faults

    inclusionProofs := []
    faults := NewFaultSet()

    for n in 0..POST_CHALLENGES_COUNT {
        attempt := 0
        'inner: for {
            challenge := DerivePoStChallenge(seed, n, faults, attempt)
            sector := challenge % sectorSize
            // check if we landed in a previously marked faulty sector
            if !faults.Contains(sector) {
                attempt += 1
                continue
            }

            challenge_value = challenge / sectorSize
            inclusionProof, isFault := GenerateMerkleInclusionProof(sector, challenge_value)
            if isFault {
                // faulty sector, generate a new challenge
                faults.Add(sector)
                attempt += 1
            } else {
                // no fault, move on to the next challenge
                inclusionProofs[n] = inclusionProof
                break 'inner
            }
        }
    }

    // Generate the snark
    challenges := DerivePoStChallenges(sectorSize, seed, faults)

    snark_proof := GeneratePoStSnark(sectorSize, challenges, sectors, inclusionProofs)

    proof := PoStProof {
        snark: snark_proof
    }

    return proof, faults
}
```

#### Verification


```go
func VerifyPoSt(sectorSize BytesAmount, sectors []commR, seed []byte, proof PoStProof, faults FaultSet) bool {
    challenges := DerivePoStChallenges(sectorSize, seed, faults)

    // Verify snark
    VerifyPoStSnark(sectorSize, challenges, sectors)
}
```


#### Types

```go
type PoStProof struct {
    snark []byte
}
```

#### Challenge Derivation

```go
// Derive the full set of challenges for PoSt.
func DerivePoStChallenges(sectorSize BytesAmount, seed []byte, faults FaultSet) [POST_CHALLENGES_COUNT][]byte {
    challenges := []

    for n in 0..POST_CHALLENGES_COUNT {
        attempt := 0
        'inner: for {
            challenge := DerivePoStChallenge(seed, n, faults, attempt)

            // check if we landed in a faulty sector
            sector := challenge % sectorSize
            if !faults.Contains(sector) {
                // Valid challenge
                challenges[n] = challenge
                break 'inner
            }
            // invalid challenge, regenerate
            attempt += 1
        }
    }

    return challenges
}

// Derive a single challenge for PoSt.
func DerivePoStChallenge(seed []byte, n Uint, attempt Uint) []byte {
    n_bytes := WriteUintToLittleEndian(n)
    data := concat(seed, n_bytes, WriteUintToLittleEndian(attempt))
    challenge := blake2b(data)
}
```


### PoSt Circuit

#### Public Parameters

*Parameters that are embeded in the circuits or used to generate the circuit*

- `POST_CHALLENGES_COUNT: UInt`: Number of challenges.
- `POST_TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is `(log_2(Size of original data in bytes/32 bytes per leaf))`.
- `SECTOR_SIZE: UInt`:
- `MAX_SECTORS_COUNT`: maximum number of sectors that can be proven with a single post

#### Public Inputs

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `CommRs: [POST_CHALLENGES_COUNT]Fr`: The Merkle tree root hashe of all CommRs
- `InclusionPaths: [POST_CHALLENGES_COUNT]Fr`: Inclusion paths for the replica leafs. (Binary packed bools)

{{% notice todo %}}
**Todo**: `CommRs` should be optimized, by combining them into a single merkle tree, with a single root `CommA`.
Benchmark this first before commiting.
{{% /notice %}}

#### Private Inputs

*Inputs that the prover uses to generate a SNARK proof, these are not needed by the verifier to verify the proof*

- `InclusionProofs: [POST_CHALLENGES_COUNT][TREE_DEPTH]Fr`: Merkle tree inclusion proofs.
- `InclusionValues: [POST_CHALLENGES_COUNT]Fr`: Value of the encoded leaves for each challenge.


#### Circuit

##### High Level

In high level, we do 1 check:

1. **Inclusion Proofs Checks**: Check the inclusion proofs

##### Details

```go
for c in range POST_CHALLENGES_COUNT {
  // Inclusion Proofs Checks
  assert(MerkleTreeVerify(CommRs[c], InclusionPath[c], IncludionProof[c], InclusionValue[c]))
}
```

#### Verification of PoSt proof

- SNARK proof check: **Check** that given the SNARK proof and the public inputs, the SNARK verification outputs true
- Challenges check: rederive the challenges, based on the seed, and check that they are equal.
