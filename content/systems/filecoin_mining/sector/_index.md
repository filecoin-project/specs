---
title: Sector
weight: 1
bookCollapseSection: true
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Sector

Sectors are the basic units of storage on Filecoin. They have standard sizes, as well as well-defined time-increments for commitments. The size of a sector balances security concerns against usability. A sectorʼs lifetime is determined in the storage market, and sets the promised duration of the sector.

In the first iteration of the protocol, 32GiB and 64GiB sectors are supported. Maximum sector lifetime is determined by the proof algorithm. Maximum sector lifetime is initially 18 months. A sector naturally expires when it reaches the end of its lifetime. Additionally, the miner can extend the lifetime of their sectors. Rewards are earned and collaterals recovered when the miner fulfils their commitment.

Individual deals are formed when a storage miner and client are matched on Filecoinʼs storage market. The protocol does not distinguish miners matching with real clients from miners generating self-deals. However, **committed capacity** is a construction that is introduced to make self-dealing unnecessary and economically irrational. In earlier designs of the network, only sectors filled with deals increased the minerʼs likelihood of winning the block reward. This led to the expectation that miners would attack and exploit the network by playing the role of both storage provider and client, creating a malicious self-deal.

If a sector is only partially full of deals, the network considers the remainder to be _committed capacity_. Similarly, sectors with no deals are called committed capacity sectors; miners are rewarded for proving to the network that they are pledging storage capacity and are encouraged to find clients who need storage. When a miner finds storage demand, they can upgrade their committed capacity sectors to earn additional revenue in the form of a deal fee from paying clients. More details on how to add storage and upgrade sectors in [Adding Storage](adding_storage).

Committed capacity sectors improve minersʼ incentives to store client data, but they donʼt solve the problem entirely. Storing real client files adds some operational overhead for storage miners. In certain circumstances – for example, if a miner values block rewards far more than deal fees – miners might still choose to ignore client data entirely and simply store committed capacity to increase their storage power as rapidly as possible in pursuit of block rewards. This would make Filecoin less useful and limit clientsʼ ability to store data on the network. Filecoin addresses this issue by introducing the concept of verified clients. Verified clients are certified by a decentralized network of verifiers. Once verified, they can post a predetermined amount of verified client deal data to the storage market, set by the size of their DataCap. Sectors with verified client deals are awarded more storage power – and therefore more block rewards – than sectors without. This provides storage miners with an additional incentive to store client data.

Verification is not intended to be scarce – it will be very easy to acquire for anyone with real data to store on Filecoin. Even though verifiers may allocate verified client DataCaps liberally (yet responsibly and transparently) to make onboarding easier, the overall effect should be a dramatic increase in the proportion of useful data stored on Filecoin.

Once a sector is full (either with client data or as committed capacity), the unsealed sector is combined by a proving tree into a single root `UnsealedSectorCID`. The sealing process then encodes (using CBOR) an unsealed sector into a sealed sector, with the root `SealedSectorCID`.

This diagram shows the composition of an unsealed sector and a sealed sector.

![Unsealed Sectors and Sealed Sectors](sectors.png)

**Sector Storage & Window PoSt**

The Lotus implementation of the Window PoSt scheduler can be found [here](https://github.com/filecoin-project/lotus/blob/master/storage/wdpost/wdpost_sched.go) and the actual execution of Window PoSt on a sector can be found [here](https://github.com/filecoin-project/lotus/blob/master/storage/wdpost/wdpost_run.go).

The Lotus block store implementation for sectors can be found [here](https://github.com/filecoin-project/lotus/blob/master/storage/sectorblocks/blocks.go).
