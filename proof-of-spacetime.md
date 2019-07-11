## Proof-of-Spacetime

This document descibes Rational-PoSt, the Proof of Spacetime used in Filecoin.


## Rational PoSt

### Definitions

- **PoSt Proving Period**: The time interval in which a PoSt has to be submitted.
- **PoSt Proving Time**: The time it takes to actually run a single PoSt.
- **PoSt Generation Start Time**: The time at which the actual work of generating the PoSt should be started. This is some delta befre the end of the `Proving Period`, and as such less then a single `Proving Period`.


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

- **Setup Parameters**
  - Same as Public parameters.
- **Public Parameters**
  - `challenge_count`: Number of challenges to be asked at each iteration.
  - `sector_size`: Size of the sealed sector in bytes.
  - `sectors_count`: Number of sectors over which the proof is performed.
  - `challenge_bits`: Number of bits in one challenge (length of a merkle path)
  - `seed_bits`: ?
- **Public Inputs**
  - `CommRs : [sectors_count]Hash` CommRs must be ordered by their timestamp on the blockchain.
- **Private Inputs**
  - `replicas: [sectors_count]SealedSector`: sealed sectors
- **Proof**
  - `post_proofs [post_periods_count]PoRepProof`

### Methods

#### `Prove(ChallengeSeed, PublicParameters, PublicInputs, PrivateInputs) -> (Proof, Faults)`

**TODO**: describe faults

`Prove` is gets called when the `Generation Start Time` is hit (in every `Proving Period`).

- `(challenges, challenged_sectors) = derive_challenges(challenge_count, challenge_seed)`
- `porep_proof = OnlinePoRep.prove(challenges, challenged_sectors, commR, replica)`
- Output `(porep_proofs, faults)`

#### `Verify(PublicParameters, PublicInputs, Proof, Faults, ChallengeSeed) -> bool`

**TODO:** Handle the passed in `Faults**
**TODO:** Verify integration with challenges and verification

- `(challenges, challenged_sectors) = derive_challenges(challenge_count, challenge_seed)`
- Assert: `verify_final_challenge_derivation(challenges, partial_challenge, challenge_seed, challenge_bits)`
- Assert: `PoSt.Verify(CommRs, challenges, post_proof)`


### `GenerateChallengeSeed(t: Height) -> ChallengeSeed`

Before calling `Prove`, first this executed.

- `lookback_ticket = minTicket(RandomnessLookback(blk(t)))`
- `challenge_seed = blake2s(lookbackTicket)`
- Output `challenge_seed`

### Challenge Stream

In order to reduce verification costs inside the circuit, challenge generation is split into two parts, partial challenge generation, and final challenge generation.

#### `derive_challenges(challenge_count, seed)`

- `partial_challenge = derive_partial_challenge(seed)`
- `while all_challenges.len() < count`
  - `(challenges, challenged_sectors)`
  - `all_challenges.extend(challenges)`
  - `all_challenged_sectors.push(challenged_sectors)`
- Output (`all_challenges[0..challenge_count]`, `all_challenged_sectors[0..challenge_count]`)

#### `derive_partial_challenge(seed)`

- `partial_challenge = H(seed | 0)`

#### `derive_final_challenges(partial_challenge, seed, sectors_count, challenge_bits)`

- `mixed = partial_challenge - seed`
- `mixed_bytes = fr_to_bytes(bits_to_fr(mixed))`
- `for chunk in mixed_bytes.chunks(challenge_bits)`
  - `challenge = 0`
  - `place = 1`
  - `for bit in chunk`
    - `if bit`
      - `challenge += place`
    - `place = place << 1`
  - `challenged_sector = ???` **FIXME**
  - `challenges.push(challenge)`
  - `challenged_sectors.push(challenged_sectors)`
- Output `(challenges, challenged_sectors)`

#### `verify_final_challenge_derivation(challenges, partial_challenge, seed, challenge_bits) -> bool`

- `shift_factor = bytes_to_fr(1 << challenge_bits)`
- `packed = Fr(0)`
- `for challenge in challenges`
  - `fr_challenge = bytes_to_fr(challenge)`
  - `packed = packed * shift_factor`
  - `packed = packed + fr_challenge`
- `fr_seed = bytes_to_fr(seed)`
- `fr_partial = bytes_to_fr(partial_challenge)`
- `fr_mixed = fr_mixed + packed`

- Output `fr_partial == fr_mixed`
