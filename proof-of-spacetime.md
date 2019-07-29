## Proof-of-Spacetime

This document describes Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

## Rational PoSt

### Definitions

- **POST_PROVING_PERIOD**: The time interval in which a PoSt has to be submitted.
- **POST_CHALLENGE_TIME**: The time offset at which the actual work of generating the PoSt should be started. This is some delta before the end of the `Proving Period`, and as such less then a single `Proving Period`.

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

#### Fault Detection

Fault detection happens over the course of the life time of a sector. When the sector is for some reason unavailable, the miner is responsible to post an up to date `AddFaults` message to the chain. When recovering any faults, they need to submit a `RecoverFaults` message. The PoSt generation then takes the latest available `faults` of the miner to generate a PoSt matching the committed sectors and faults.

TODO: Is it okay to add a fault and recover from it in the same proving period?

#### Generation

```go
func GeneratePoSt(sectorSize BytesAmount, sectors []commR, seed []byte, faults FaultSet) PoStProof {
    // Generate the Merkle Inclusion Proofs + Faults

    inclusionProofs := []
	sectorsSorted := []
    challenges := DerivePoStChallenges(sectorSize, seed, faults)

    for n in 0..POST_CHALLENGES_COUNT {
        challenge := challenges[n]
        sector := challenge % len(sectors)

        // Leaf index of the selected sector
        challenge_value = challenge % sectorSize
        inclusionProof, isFault := GenerateMerkleInclusionProof(sector, challenge_value)
        if isFault {
            // faulty sector, need to post a fault to the chain and try to recover from it
            return Fatal("Detected late fault")
        }

        inclusionProofs[n] = inclusionProof
		sectorsSorted[i] = sectors[sector]
    }
	
    // Generate the snark
    snark_proof := GeneratePoStSnark(sectorSize, challenges, sectorsSorted, inclusionProofs)

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
    sectorsSorted := []

    // Match up commitments with challenges
    for i in 0..challenges {
        sector = challenges[i] % len(sectors)
        sectorsSorted[i] = sectors[sector]
    }

    // Verify snark
    VerifyPoStSnark(sectorSize, challenges, sectorsSorted)
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
func DerivePoStChallenges(sectorCount: Uint, seed []byte, faults FaultSet) [POST_CHALLENGES_COUNT]Uint {
    challenges := []

    for n in 0..POST_CHALLENGES_COUNT {
        attempt := 0
        while challenges[n] == nil {
            challenge := DerivePoStChallenge(seed, n, faults, attempt)

            // check if we landed in a faulty sector
            sector := challenge % sectorCount
            if !faults.Contains(sector) {
                // Valid challenge
                challenges[n] = challenge
            }
            // invalid challenge, regenerate
            attempt += 1
        }
    }

    return challenges
}

// Derive a single challenge for PoSt.
func DerivePoStChallenge(seed []byte, n Uint, attempt Uint) Uint {
    n_bytes := WriteUintToLittleEndian(n)
    data := concat(seed, n_bytes, WriteUintToLittleEndian(attempt))
    challenge := blake2b(data)
    ReadUintLittleEndain(challenge)
}
```


### PoSt Circuit

#### Public Parameters

*Parameters that are embeded in the circuits or used to generate the circuit*

- `POST_CHALLENGES_COUNT: UInt`: Number of challenges.
- `POST_TREE_DEPTH: UInt`: Depth of the Merkle tree. Note, this is `(log_2(Size of original data in bytes/32 bytes per leaf))`.
- `SECTOR_SIZE: UInt`:

#### Public Inputs

*Inputs that the prover uses to generate a SNARK proof and that the verifier uses to verify it*

- `CommRs: [POST_CHALLENGES_COUNT]Fr`: The Merkle tree root hashes of all CommRs.
- `InclusionPaths: [POST_CHALLENGES_COUNT]Fr`: Inclusion paths for the replica leafs. (Binary packed bools)

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
