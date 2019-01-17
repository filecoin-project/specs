## Proof-of-Spacetime
   
This document descibes 
 
- VDF-PoSt: a Proof-of-Spacetime using VDFs
- An Extension to PoSt to support multiple sectors
- An Extension to PoSt to support challenges taken from a Random Beacon 

## Syntax

- **PoSt Epoch**: The total time passing between Online PoReps in the PoSt Computation (in VDF, this interval is the time it takes to run a VDF and an Online PoRep prove step). We define the number of epochs as `POST_EPOCHS_COUNT`
- **PoSt Period**: The total time it takes to run a single PoSt. If a PoSt is repeated multiple times, we define the number of periods as `POST_PERIODS_COUNT`. This can be reasoned as: `PoSt Epoch * POST_EPOCHS_COUNT`. We assume that the best  Post Period time is `MIN_POST_PERIOD_TIME` 
- **Total proving time**: The time it takes to run a PoSt. Note that a PoSt could be a composition of multiple PoSt. This can be reasoned as: `PoSt Period * POST_PERIODS_COUNT` TODO Check with papers' syntax

## VDF-PoSt: Proof-of-Spacetime based on VDFs

VDF-PoSt is a Proof-of-Spacetime that hashes the input and the output of the VDFs it uses `H(Vdf(H(x)))`, hence `VDF`.

### Parameters

- **Setup Parameters** 
  - `CHALLENGE_COUNT`: number of challenges to be asked at each iteration
  - `SECTOR_SIZE`: size of the sealed sector in bytes
  - `POST_EPOCHS`: number of times we repeat an online Proof-of-Replication in one single PoSt
  - `vdf_params`: vdf public parameters
  - `sectors_count`: number of sectors over which the proof is performed
- **Public Parameters**
  - `CHALLENGE_COUNT`: number of challenges to be asked at each iteration
  - `SECTOR_SIZE`: size of the sealed sector in bytes
  - `POST_EPOCHS`: number of times we repeat an online Proof-of-Replication in one single PoSt.
  - `vdf_params`: vdf public parameters
  - `sectors_count`: number of sectors over which the proof is performed
  - `challenge_bits`: number of bits in one challenge (length of a merkle path)
  - `seed_bits`: number of bits in one challenge
- **Public Inputs**
  - `commR: Hash`: Root hash of the Merkle Tree of the sealed sec+tor
  - `challenge_seed` : [32]byte: initial randomness (in Filecoin taken from the chain) from which challenges will be generated.
- **Private Inputs** 
  - `replica: SealedSector`: sealed sector
- **Proof**
  - `ys: [POST_EPOCHS-1]Value`
  - `vdf_proofs: [POST_EPOCHS-1]VDFProof` 
  - `porep_proofs: [POST_EPOCHS]PorepProof`

### Methods

#### `Prove(Public Parameters, Public Inputs, Private Inputs) -> Proof`

- *Step 1*: Generate `POST_EPOCHS` proofs: 
 - `mix = challenge_seed`
 - `challenge_stream = NewChallengeStream(PublicParams)`
 - Repeat `POST_EPOCHS` times:
   - `(challenges, challenged_sectors) = challenge_stream(mix)` 
   - Generate proof: `porep_proof = OnlinePoRep.prove(challenges, challenged_sectors, commR, replica)`
     - Note: you can have the tree cached in memory
   - append `porep_proof` to `porep_proofs[]`
   - Add `porep_proof` to `porep_proofs`
   - Slow challenge generation from previous proof `porep_proof`:
     - Run VDF and generate a proof
        - `x = ExtractVDFInput(porep_proof))`
        - `y, vdf_proof = VDF.eval(x)`
        - Add `vdf_proof` to `vdf_proofs`
        - Add `y` to `ys`
        - `mix = y`
- Step 3: Output `porep_proofs`, `vdf_proofs`, `ys`

#### `Verify(Public Parameters, Public Inputs, Proof) -> bool`

- *VDF Output Verification*
  - For `i` in `0..POST_EPOCHS-1`
    - assert: `VDF.verify(pp_vdf, ExtractVDFInput(porep_proofs[i]), ys[i], vdf_proofs[i])`	
- *Sequential Online PoRep Verification*
  - assert: `OnlinePoRep.verify(commR, challenges_0, porep_proofs[0])`
  - for `i` in `1..POST_EPOCHS`
    - Generate challenges  `for j in 0..CHALLENGE_COUNT: challenges[j] = H(H(ys[i-1])|j)` `
    - assert: `OnlinePoRep.verify(commR, challenges, porep_proofs[i])`

## EVDF-PoSt: Extending a single CommR PoSt to multiple CommRs

**Problem**: A PoSt proves space on a single sector (whose Merkle root hash is `CommR`). In order to prove space over multiple sectors, we can either:

- Run an PoSt for each sector (PoRep guarantees): this means running `Prove` `SECTORS_COUNT` times and have the proof size to be  `SECTORS_COUNT`*`PROOF_SIZE`
- Extend a single PoSt to run sectors (PoS guarantees): this means security is not defined per sector, but across sector. For example: Assume PoSt guarantees 99% of the data being stored. A miner has 100 sectors and runs a single PoSt per sector. The worst that can happen is that the miner loses 1% of each sector. If the miner runs a single PoSt across all the sectors, then, the worst it can happen is that the miner loses 1% of all the sectors.

**Filecoin note**: In Filecoin, we use the second strategy in order to have shorter proofs. It is worth mentioning that misbehaving miners have an economic incentive not to misbehave in Filecoin. This section documents how to extend VDF-PoSt over a single sector, into a VDF-PoSt over multiple sectors.

### Difference between standard VDF-PoSt and the extension

- **Public Parameters** & **Setup Parameters**
  - `SECTORS_COUNT` which is the number of sectors over which we are running PoSt
- **Public Inputs**
  - `CommRs : [SECTORS_COUNT]Hash` instead of `CommR : Hash`. CommRs must have a specific order (e.g. lexographical order, order of timestamps on the blockchain)
  - `challenge_seed` : [32]byte: initial randomness (in Filecoin taken from the chain) from which challenges will be generated.
- **Private Inputs**
  - `replicas: [SECTORS_COUNT]replica` instead of `replica`. (same order as the `CommRs`)
- **Prove** & **Verify** Computation
  - A challenge in `challenges` points to a leaf in one of the sectors
    - Sector is chosen by `challenge % SECTORS_COUNT` (TODO check if this is fine)
    - Leaf is chosen in the same way as in Online Porep (`challenge % SECTORS_SIZE/32`)

### Methods

#### `Prove(Public Parameters, Public Inputs, Private Inputs) -> Proof`

- *Step 1:* Generate first proof
  - Generate proof `pos_proof = OnlinePoS.prove(commRs, challenges, replicas)`
  - Add `porep_proof` to `pos_proofs`
- *Step 2:* Generate `POST_EPOCHS - 1` remaining proofs:
  - Repeat `POST_EPOCHS - 1` times:
    - Slow challenge generation from previous proof pos_proof`:
      - Run VDF and generate a proof
        - `x = ExtractVDFInput(pos_proof))`
        - `y, vdf_proof = VDF.eval(x)`
        - Add `vdf_proof` to `vdf_proofs`
        - Add `y` to `ys`
        - `r = H(y)`
      - Generate challenges  `for i in 0..CHALLENGE_COUNT: challenges[i] = H(r|i)`
    - Generate a proof as done in Step 1
- Step 3: Output `pos_proofs`, `vdf_proofs`, `ys`

#### `Verify(Public Parameters, Public Inputs, Proof) -> bool`

- *VDF Output Verification*
  - For `i` in `0..POST_EPOCHS-1`
    - assert: `VDF.verify(pp_vdf, ExtractVDFInput(pos_proofs[i]), ys[i], vdf_proofs[i])`	
- *Sequential Online PoRep Verification*
  - assert: `OnlinePoS.verify(commR, challenges_0, pos_proofs[0])`
  - for `i` in `1..POST_EPOCHS`
    - Generate challenges  `for j in 0..CHALLENGE_COUNT: challenges[j] = H(H(ys[i-1])|j)` `
    - assert: `OnlinePoS.verify(commR, challenges, pos_proofs[i])`

### Security note

- **Avoiding grinding**: If the prover can choose arbitrary `SECTORS_COUNT`, after receiving a challenge, they can try different sector sizes to have more favourable challenged leaves. In order to avoid this, the prover commit to the `SECTORS_COUNT`, and the `CommR`s before receiving the challenges. In Filecoin, we get this for free, since all the sectors to be proven are committed on chain and the `SECTORS_COUNT` can't be altered.
- **Storage security**: An VDF-PoSt with a single CommR inherits the Online PoRep security guarentees, while this extension does not. In VDF-PoSt, the prover answer `CHALLENGES_COUNT` challenges on a single sector, in this extension, the prover answers `CHALLENGES_COUNT` across multiple sectors. 

## BeaconPost: Taking challenges over time via a Random Beacon

**Problem with large `POST_EPOCH_COUNTS`**: Different VDF hardware run at different speed. A small percentage of gain in a `PoSt Epoch` would result in a large time difference in `Total Proving Time` between the fastest and the slowest prover. We call the difference between fastest and average prover `VDF speedup gap`. We define a VDF Speedup gap as a percentage (0-1) and we assume a concrete gap for a PoSt Period between the assumed fastest and the best known prover. We define this gap as `VDF_SPEEDUP_GAP`.

**Mitigating VDF Speedups**: We break up a PoSt into multiple PoSt Periods. Each period must take challenges from a Random Beacon which outputs randomness every interval  `MIN_POST_PERIOD_TIME` . In this way, the faster prover can be  `VDF_SPEEDUP_GAP` faster in each PoSt Period, but cannot be `VDF_SPEEDUP_GAP` faster over the Total Proving Period. In other words, the fastest prover cannot accumulate the gains at each PoSt period because, they have to wait for the new challenges from the Random Beacon. In the case of Filecoin, the blockchain acts as a Random Beacon).

- **Setup Parameters** 
  - Same as Public parameters.
- **Public Parameters**
  - `POST_PUBLIC_PARAMS`: Public Parameters as defined for VDF PoSt.
  - `POST_PERIODS_COUNT: uint`
- **Public Inputs**
  - `	CommRs : [SECTORS_COUNT]Hash` instead of `CommR : Hash`. CommRs must have a specific order (e.g. lexographical order, order of timestamps on the blockchain)
- **Private Inputs** 
  - `replicas: [SECTORS_COUNT]SealedSector`: sealed sectors
- **Proof**
  - `post_proofs [POST_PERIODS_COUNT]VDFProof`

### Methods

#### `Prove(Public Parameters, Public Inputs, Private Inputs) -> Proof`

Prove is a process that has access to a Random Beacon functionality that outputs new randomness every `MIN_POST_PERIOD_TIME`:

- `t = 0`:
- For `t = 0..POST_PERIODS_COUNT`:
  - Query Random Beacon: `challenge_seed = RandomBeacon(t)`
  - Compute a `post_proofs[t] = PoSt.Prove(CommRs, challenge_seed, replicas)` 
- Outputs `post_proofs`

#### `Verify(Public Parameters, Public Inputs, Proof) -> bool`

- `t = 0`:
  - Query Random Beacon: `r = RandomBeacon(t)`
  - Generate challenges: `for i=0..CHALLENGES_COUNT: challenges[i] = H( r | t | i)`
  - Assert: `PoSt.Verify(CommRs, challenges, post_proofs[t])` 
- For `t = 1..POST_PERIODS_COUNT`:
  - Query Random Beacon: `r = RandomBeacon(t)`
  - Generate challenges: `for i=0..CHALLENGES_COUNT: challenges[i] = H(ExtractPoStInput(post_proofs[t-1]) | r | t | i )`
  - Assert: `PoSt.Verify(CommRs, challenges, post_proofs[t])` TODO check

### Random Beacon functionality

A Random Beacon outputs a single randomness every `MIN_POST_PERIOD_TIME`.

## Other Functions used

### ExtractVDFInput

##### Inputs

- `porep_proof PoRep.Proof`

##### Computation

- Hash the concatenation of the leaves of each tree in `OnlinePoRep.Proof`

### VDF

- `VDF.setup() -> VDFPublicParams`
- `VDF.eval(pp: VDFPublicParams, x: Value) -> (Value, VDFProof)`
- `VDF.verify(pp: VDFPublicParams, x: Value, y: Value, proof: VDFProof) -> bool `
