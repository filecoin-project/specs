---
menuTitle: "Change Log"
title: "Change Log - Version History"
---

# v1.1 - 2019-10-30 - `c3f6a6dd`

- **Deals on chain**
    - Storage Deals
    - Full `StorageMarketActor` logic:
        - client and miner balances: deposits, locking, charges, and withdrawls
        - collateral slashing
    - Full `StorageMinerActor` logic:
        - sector states, state transitions, state accounting, power accounting
        - DeclareFaults + RecoverSectors flow
        - `CommitSector` flow
        - `SubmitPost` flow
            - Sector proving, faults, recovery, and expiry
        - `OnMissedPost` flow
            - Fault sectors, drop power, expiry, and more
    - `StoragePowerActor`
        - power accounting based on `StorageMinerActor` state changes
        - Collaterals: deposit, locking, withdrawal
        - Slashing collaerals
    - Interactive-Post
        - `StorageMinerActor`: `PrecommitSector` and `CommitSector`
    - Surprise-Post
        - Challenge flow through `CronActor -> StoragePowerActor -> StorageMiner`
- **Virtual Machine**
    - Extracted VM system out of blockchain
    - Addresses
    - Actors
        - Separation of code and state
    - Messages
        - Method invocation representation
    - Runtime
        - Slimmed down interface
        - Safer state Acquire, Release, Commit flow
        - Exit codes
        - Full invocation flow
        - Safer recursive context construction
        - Error levels and handling
        - Detecting and handling out of gas errors
    - Interpreter
        - `ApplyMessage`
        - `{Deduct,Deposit} -> Transfer` - safer
        - Gas accounting
    - VM system actors
        - `InitActor` basic flow, plug into Runtime
        - `CronActor` full flow, static registry
    - `AccountActor` basic flow
- **Data Transfer**
    - Full Data Transfer flows
        - push, pull, 1-RTT pull
    - protocol, data structures, interface
    - diagrams
- **blockchain/ChainSync:**
    - first version of ChainSync protocol description
    - Includes protocol state machine description
    - Network bootstrap -- connectivity and state
    - Progressive Block Validation
    - Progressive Block Propagation
- **Other**
    - Spec section status indicators
    - Changelog

# v1.0 - 2019-10-07 - `583b1d06`

- **Full spec reorganization**
- **Tooling**
    - added a build system to compile tools
    - added diagraming tools (dot, mermaid, etc)
    - added dependency installation
    - added Orient to calculate protocol parameters
- **Content**
    - **filecoin_nodes**
        - types - an overview of different filecoin node types
        - repository - local data-structure storage
        - network interface - connecting to libp2p
        - clock - a wall clock
    - **files & data**
        - file - basic representation of data
        - piece - representation of data to store in filecoin
    - **blockchain**
        - blocks - basic blockchain data structures (block, tipset, chain, etc)
        - storage power consensus - basic algorithms and crypto artifacts for SPC
        - `StoragePowerActor` basics
    - **token**
        - skeleton of sections
    - **storage mining**
        - storage miner: module that controls and coordinates storage mining
        - sector: unit of storage, sealing, crypto artifacts, etc.
        - sector index: accounting sectors and metadata
        - storage proving: seals, posts, and more
    - **market**
        - deals: storage market deal basics
        - storage market: `StorageMarketActor` basics
    - **orient**
        - orient models for proofs and block sizes
    - **libraries**
        - filcrypto - sealing, PoRep, PoSt algorithms
        - ipld - cids, ipldstores
        - libp2p - host/node representation
        - ipfs - graphsync and bitswap
        - multiformats - multihash, multiaddr
    - **diagrams**
        - system overview
        - full protocol mermaid flow


# pre v1.0

- Extensive write up of the filecoin protocol - visible [here](https://github.com/filecoin-project/specs/tree/prevspec)
- See full changelog: https://github.com/filecoin-project/specs/commits/prevspec

