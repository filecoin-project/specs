---
title: Sector Lifecycle
weight: 2
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Sector Lifecycle

Once the sector has been generated and the deal has been incorporated into the Filecoin blockchain, the storage miner begins generating Proofs-of-Spacetime (PoSt) on the sector, starting to potentially win block rewards and also earn storage fees. Parameters are set so that miners generate and capture more value if they guarantee that their sectors will be around for the duration of the original contract. However, some bounds are placed on a sectorʼs lifetime to improve the network performance.

In particular, as sectors of shorter lifetime are added, the networkʼs capacity can be bottlenecked. The reason is that the chainʼs bandwidth is consumed with new sectors only replacing expiring ones. As a result, a minimum sector lifetime of six months was introduced to more effectively utilize chain bandwidth and miners have the incentive to commit to sectors of longer lifetime. The maximum sector lifetime is limited by the security of the present proofs construction. For a given set of proofs and parameters, the security of Filecoinʼs Proof-of-Replication (PoRep) is expected to decrease as sector lifetimes increase.

It is reasonable to assume that miners enter the network by adding Committed Capacity sectors, that is, sectors that do not contain user data. Once miners agree storage deals with clients, they upgrade their sectors to Regular Sectors. Alternatively, if they find Verified Clients and agree a storage deal with them, they upgrade their sector accordingly. Depending on whether or not a sector includes a (verified) deal, the miner acquires the corresponding storage power in the network.

All sectors are expected to remain live until the end of their sector lifetime and early dropping of sectors will result in slashing. This is done to provide clients a certain level of guarantee on the reliability of their hosted data. Sector termination comes with a corresponding _termination fee_.

As with every system it is expected that sectors will present faults. Although this might degrade the quality offered by the network, the reaction of the miner to the fault drives system decisions on whether or not the miner should be penalized. A miner can recover the faulty sector, let the system terminate the sector automatically after 42 days of faults, or proactively terminate the sector immediately in the case of unrecoverable data loss. In case of a faulty sector, a small penalty fee approximately equal to the block reward that the sector would win per day is applied. The fee is calculated per day of the sector being unavailable to the network, i.e. until the sector is recovered or terminated.

Miners can extend the lifetime of a sector at any time, though the sector will be expected to remain live until it has reached the end of the new sector lifetime. This can be done by submitting a `ExtendedSectorExpiration` message to the chain.

A sector can be in one of the following states.

| State          | Description                                                                                                                                           |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Precommitted` | Miner seals sector and submits `miner.PreCommitSector` or `miner.PreCommitSectorBatch`                                                                |
| `Committed`    | Miner generates a Seal proof and submits `miner.ProveCommitSector` or `miner.ProveCommitAggregate`                                                    |
| `Active`       | Miner generate valid PoSt proofs and timely submits `miner.SubmitWindowedPoSt`                                                                        |
| `Faulty`       | Miner fails to generate a proof (see Fault section)                                                                                                   |
| `Recovering`   | Miner declared a faulty sector via `miner.DeclareFaultRecovered`                                                                                      |
| `Terminated`   | Either sector is expired, or early terminated by a miner via `miner.TerminateSectors`, or was failed to be proven for 42 consecutive proving periods. |
