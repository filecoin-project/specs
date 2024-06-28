---
title: Adding Storage
weight: 7
dashboardWeight: 2
dashboardState: stable
dashboardAudit: wip
dashboardTests: 0
---

# Adding Storage

A Miner adds more storage in the form of Sectors. Adding more storage is a two-step process:

1. **PreCommitting a Sector**: A Miner publishes a Sector's SealedCID, through `miner.PreCommitSector` of `miner.PreCommitSectorBatch`, and makes a deposit. The Sector is now registered to the Miner, and the Miner must ProveCommit the Sector or lose their deposit.
2. **ProveCommitting a Sector**: The Miner provides a Proof of Replication (PoRep) for the Sector through miner.ProveCommitSector or miner.ProveCommitAggregate. This proof must be submitted AFTER a delay (the InteractiveEpoch), and BEFORE PreCommit expiration.

This two-step process provides assurance that the Miner's PoRep _actually proves_ that the Miner has replicated the Sector data and is generating proofs from it:

- ProveCommitments must happen AFTER the InteractiveEpoch (150 blocks after Sector PreCommit), as the randomness included at that epoch is used in the PoRep.
- ProveCommitments must happen BEFORE the PreCommit expiration, which is a boundary established to make sure Miners don't have enough time to "fake" PoRep generation.

For each Sector successfully ProveCommitted, the Miner becomes responsible for continuously proving the existence of their Sectors' data. In return, the Miner is awarded storage power.

# Upgrading Sectors

Miners are granted storage power in exchange for the storage space they dedicate to Filecoin. Ideally, this storage space is used to store data on behalf of Clients, but there may not always be enough Clients to utilize all the space a Miner has to offer.

In order for a Miner to maximize storage power (and profit), they should take advantage of all available storage space immediately, _even before they find enough Clients to use this space_.

To facilitate this, there are _two types_ of Sectors that may be sealed and ProveCommitted:

- **Regular Sector**: A Sector that contains Client data
- **Committed Capacity (CC) Sector**: A Sector with no data (all zeroes)

Miners are free to choose which types of Sectors to store. CC sectors, in particular, allow Miners to immediately make use of existing disk space: earning storage power and a higher chance at producing a block. Miners can decide if they should upgrade their CC sectors to take client deals or continue proving CC sectors. Currently, CC sectors store randomness by default in client implementation, but this does not preclude miners from storing any type of useful data that increase their private utility in CC sectors (as long as it is legal). The protocol expects that new use-cases and diversity will emerge out of such behaviour.

To incentivize Miners to hoard storage space and dedicate it to Filecoin, CC Sectors have a unique capability: **they can be "upgraded" to Regular Sectors** (also called "replacing a CC Sector").

Miners upgrade their ProveCommitted CC Sectors by PreCommitting a Regular Sector, and specifying that it should replace an existing CC Sector. Once the Regular Sector is successfully ProveCommitted, it will replace the existing CC Sector. If the newly ProveCommitted Regular sector contains a Verified Client deal, i.e., a deal with higher Sector Quality, then the miner's storage power will increase accordingly.

Upgrading capacity currently involves resealing, that is, creating a unique representation of the new data included in the Sector through a computationally intensive process. Looking ahead, committed capacity upgrades should eventually be possible without a reseal. A succinct and publicly verifiable proof that the committed capacity has been correctly replaced with replicated data should achieve this goal. However, this mechanism must be fully specified to preserve the security and incentives of the network before it can be implemented and is, therefore, left as a future improvement.
