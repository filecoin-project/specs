## Proof-of-Spacetime

This document describes Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

## Rational PoSt

### Definitions

- **POST_PROVING_PERIOD**: The time interval in which a PoSt has to be submitted.
- **POST_CHALLENGE_TIME**: The time offset at which the actual work of generating the PoSt should be started. This is some delta before the end of the `Proving Period`, and as such less then a single `Proving Period`.

### High Level API

#### Fault Detection

Fault detection happens over the course of the life time of a sector. When the sector is for some reason unavailable, the miner is responsible to post an up to date `AddFaults` message to the chain. The PoSt generation then takes the latest available `faults` of the miner to generate a PoSt matching the committed sectors and faults.

At the beginning of a new proving period all faults are reset, and if they persist the miner needs to resubmit an `AddFaults` message.

#### Generation

`GeneratePoSt` generates a __*Proof of Spacetime*__ over all  __*sealed sectors*__ of a single miner— identified by their `commR` commitments. This is accomplished by performing a series of merkle inclusion proofs (__*Proofs of Retrievability*__). Each proof is of a challenged node in a challenged sector. The challenges are generated pseudo-randomly, based on the provided `seed`. At each time step, a number of __*Proofs of Retrievability*__ are performed.

```go
func GeneratePoSt(sectorSize BytesAmount, sectors []commR, seed Seed, faults FaultSet) PoStProof {
    // Generate the Merkle Inclusion Proofs + Faults

    inclusionProofs := []
	sectorsSorted := []
    challenges := DerivePoStChallenges(seed, faults, sectorSize, len(sectors))

    for challenge in challenges {
        // Leaf index of the selected sector
        inclusionProof, isFault := GenerateMerkleInclusionProof(challenge.Sector, challenge.Leaf)
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

`VerifyPoSt` is the functional counterpart to `GeneratePoSt`. It takes all of `GeneratePoSt`'s output, along with those of `GeneratePost`'s inputs required to identify the claimed proof. All inputs are required because verification requires sufficient context to determine not only that a proof is valid but also that the proof indeed corresponds to what it purports to prove.

```go
func VerifyPoSt(sectorSize BytesAmount, sectors []commR, seed Seed, proof PoStProof, faults FaultSet) bool {
    challenges := DerivePoStChallenges(seed, faults, sectorSize, len(sectors))
    sectorsSorted := []

    // Match up commitments with challenges
    for challenge in challenges {
        sectorsSorted[i] = sectors[challenge.Sector]
    }

    // Verify snark
    VerifyPoStSnark(sectorSize, challenges, sectorsSorted)
}
```


#### Types

```go
// The final proof stored on chain.
type PoStProof struct {
    snark []byte
}
```

```go
// The random challenge seed, provided by the chain.
Seed [32]byte
```

```go
type Challenge struct {
    Sector Uint
    Leaf Uint
}
```

#### Challenge Derivation

```go
// Derive the full set of challenges for PoSt.
func DerivePoStChallenges(seed Seed, faults FaultSet, sectorSize: Uint, sectorCount: Uint) [POST_CHALLENGES_COUNT]Challenge {
    challenges := []

    for n in 0..POST_CHALLENGES_COUNT {
        attempt := 0
        while challenges[n] == nil {
            challenge := DerivePoStChallenge(seed, n, faults, attempt, sectorSize, sectorCount)

            // check if we landed in a faulty sector
            if !faults.Contains(challenge.Sector) {
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
func DerivePoStChallenge(seed Seed, n Uint, attempt Uint, sectorSize Uint, sectorCount: Uint) Challenge {
    n_bytes := WriteUintToLittleEndian(n)
    data := concat(seed, n_bytes, WriteUintToLittleEndian(attempt))
    challenge_bytes := blake2b(data)

    sector_challenge := ReadUintLittleEndian(challenge_bytes[0..8])
    leaf_challenge := ReadUintLittleEndian(challenge_bytes[8..16])

    Challenge {
        Sector: sector_challenge % sectorCount,
        Leaf: leaf_challenge % (sectorSize / NODE_SIZE),
    }
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
