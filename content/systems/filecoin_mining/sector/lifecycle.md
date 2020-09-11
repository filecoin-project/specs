---
title: Sector Lifecycle
weight: 6
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Sector Lifecycle Summary

Once the sector has been generated and the deal has been incorporated into the Filecoin blockchain, the storage miner begins generating Proofs-of-Spacetime (PoSt) on the sector, starting to potentially win block rewards and also earn storage fees. Parameters are set so that miners generate and capture more value if they guarantee that their sectors will be around for the duration of the original contract. However, some bounds are placed on a sectorʼs lifetime to improve the network performance. As sectors of shorter lifetime are added, the networkʼs capacity can be bottlenecked. The reason is that the chainʼs bandwidth is consumed with new sectors only replacing expiring ones. As a result, a minimum sector lifetime of six months was introduced to more effectively utilize chain bandwidth and miners have the incentive to commit to sectors of longer lifetime. The maximum sector lifetime is limited by the security of the present proofs construction. For a given set of proofs and parameters, the security of Filecoinʼs Proof-of-Replication (PoRep) is expected to decrease as sector lifetimes increase.
