# Proofs Requirements: Security, Scaling, Costs

Author: @nicola

This document explores the design space for Proof of Replication and Proof of Spacetime and it describes the security requirements, the scaling requirements and the cost requirements.

------

## Context

### Miner profiles

We define two types of miners that we will use in our calculations.

#### Slowest Miner

The slowest miner is the slowest admissible miner; slower miners will be out of our optimizations spectrum and might not meet cost requirements and security requirements. In other words, it might be too expensive / impractical for slower miners to mine. 

Optimizing our proofs for this miner results in reducing the cost for mining.

Current Minimal required spec: 12 core CPU i9 5GHz.

#### Fastest miner

The fastest miner is the fastest tolerable miner; faster are miners will be able to perform generation attacks. The specification of the fastest miner are based on software and harware optimizations and must exceed what can be achieved in practice in the next 5-10 years.

Optimizing our proofs for this miner results in increasing security for Filecoin.

Current spec required: unbounded parallelization, ASICs hardware for PoRep.

### Types of requirements

#### Scaling requirements

Scaling requirements constrain the design space so that: Filecoin scales beyond hexabytes of storage.

Some examples of scaling requirements are: the on-chain proof footprint must be small, new storage must be on-boarded on practical times.

#### Security requirements

Security requirements constrain the design space so that:

- Filecoin power table is sound: no fastest miner can claim more than allowed epsilon storage.
- Filecoin storage market is sound: no fastest miner can lose more than epsilon users's data.

#### Cost requirements

Cost requirements constrain the design space so that mining Filecoin is not prohibitively expensive to mine.

## Filecoin metrics

The following are standard Filecoin metrics that shall be used in tests, performance evaluations and predictions. It's important that these metrics are used across the project in order to evaluate improvements/optimizations formally.

Table 1: Summary of the key metrics

| Metric                | Scaling Requirement    | Security Requirement | Cost requirement | Current | Target  |
| --------------------- | :--------------------- | -------------------- | ---------------- | ------- | ------- |
| `cpu_core_cost`       | *                      | *                    | ≤1GiBps          |         | ≤1GiBps |
| `repl_memory_slowest` | todo                   | *                    | *                |         | todo    |
| `repl_time_slowest`   | todo                   | *                    | *                |         | todo    |
| `polling_time`        | ≥ `late_post_blocks`** | todo                 |                  |         | todo    |
| `proving_period`      | todo                   |                      |                  |         | todo    |

\*: the requirement is explicit in other metrics

\**: TODO We might be too strict, and allow for space savings at the cost of losing collateral

### CPU Core Replication Costs

These metrics highlight the replication cost for 1 GiB per second for a single CPU core (with i9 5GHz spec).

#### `cpu_core_cost` 

Average GiB replicated per second on a single CPU core  

**Unit**: GiBps (on i9 5GHz)  

**Role**: Upper-bound the computation cost of replication. Optimizing towards this metric results in lower computational cost

**Requirements**:

- Cost requirements: The lower this number, the lower the computational cost of replication.

#### Other useful CPU core metrics for performance analysis

We differentiate across CPU core cost for snarks, merkle trees since each might be individually optimized.

- `cpu_core_snark_cost` (unit: GiBps on i9 5GHz): Average GiB proved during a SNARK computation per second on a single CPU core.

  Note that with current construction, this scales logarithmically in the size of the data and it should be compared on a logarithmic scale when evaluated on different sector sizes.

- `cpu_core_merkletree_cost` (unit: GiBps on i9 5GHz): Average GiB for which a Merkle Tree has been generated per second.

- `cpu_core_kdf_cost` (unit: GiBps on i9 5GHz): Average GiB for which a KDF have been performed per second.

### Wall-clock metrics

#### `repl_time_slowest` 

Replication time for the slowest miner for the current `sector_size`.

**Unit**: Seconds

**Role**: Estimate practical sector sizes

**Note**: this can be calculated by `repl_step_time_slowest`* size of the graph.

#### `proving_period`

The Proving Period is the time between two PoSt proof submissions

- Scaling requirements: The shorter the proving period, the higher the proof size per sector per year
- Security requirements: none

#### `polling_time`

Polling time is the time between two online porep proofs in a proving period, it is defined by time that it takes to the fastest prover to re-encode 1% of the data

**Requirements**:

- Scaling requirements:

  - With current implementation:
    - Proofs requirements: The shorter the polling time, the higher the number of proofs in a fixed proving period, this leads in a higher circuit size. There is a function that determines the smallest feasible polling time for a proving period (TODO: write this function). This will not be a problem with batching proofs.
  - Late submission requirement (at least X blocks): a prover might not manage to get their PoSt on time on chain, how many blocks guarantee that their PoSt goes on chain?

- Security requirements:

  - Storage saving attack: The polling time is constraint by the time that it takes the fastest miner to re-encode a percentage of the data (target 1%). 

    From current best known attacks, the fastest miner can delete 1% of the data and regenerate it in `steps = 1/4 * N` sequential steps, where N is the number of nodes in the graph.

    **Note 1:** We should be conservative and expect attacks to improve (hence `steps = (1/4 + drg_attack_tolerance)*n`
    **Note 2**: The min polling time should be defined as `steps * repl_step_time_fastest`

    **Scratch calculations: ** 7.3 * 10^-7  * 1/4 * 2^25*64 = 391s

**Notes**:

- What if polling time is too small?
  - Make sector size larger
  - Increase the number of VDF evaluations
  - Increase the percentage of data the malicious miner can save

#### Other useful metrics

##### `repl_step_time_slowest`

Average replication time for encoding a single node in the graph

**Unit**: seconds

**Role**: Compare different algorithm optimizations and space tradeoffs

**Requirements**:

- Scaling requirements: the higher the speed, the slower new storage is onboarded on Filecoin.
- Cost requirements: none, this metric is not used for costs, see `cpu_core_cost` instead.

##### `repl_step_time_fastest`

Estimated replication time for encoding a single node in the graph

**Unit**: seconds

**Role**: Set an upper bound for the polling time (see Polling time section)

**Requirements**:

- Security requirements: the smaller the fastest replication time for a single node, the shorter the polling period

### Memory metrics

#### `repl_memory_slowest`

Exact memory needed by the slowest miner to perform replication

In ZigZag, this is calculated by: (Layer+1) * Merkle tree + Original Data + Expansion parents 

- Merkle Tree = 2 * Original Data
- Base Parents = 2 * Base Degree * Original Data * 64bits
- Expansion Parents = Expansion Degree * Original Data * 64bits

#### Other useful metrics for performance analysis

##### `repl_memory_fastest`

Memory needed by the slowest miner to perform replication

In ZigZag, this is calculated by: (Layer+1) * Merkle tree + Original Data + Base Parents + Expansion parents 

- Merkle Tree = 2 * Original Data
- Base Parents = 2 * Base Degree * Original Data * 64bits
- Expansion Parents = Expansion Degree * Original Data * 64bits

TODO: this is not accounting for DRG saving at each layer

##### `space_advantage`

Space ratio: SM space / FM space (Security requirement)
