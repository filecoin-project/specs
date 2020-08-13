---
title: Sector
weight: 2
bookCollapseSection: true
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Sector
---

The `Sector` is a fundamental "storage container" abstraction used in Filecoin Storage Mining. It is the basic unit of storage,
and serves to make storage conform to a set of expectations.

New sectors are empty upon creation. As the miner receives client data, they fill or "pack" the piece(s) into an unsealed sector.

Once a sector is full, the unsealed sector is combined by a proving tree into a single root `UnsealedSectorCID`. The sealing process then encodes (using CBOR) an unsealed sector into a sealed sector, with the root `SealedSectorCID`.

This diagram shows the composition of an unsealed sector and a sealed sector.

{{<figure src="sectors.png" title="Unsealed Sectors and Sealed Sectors" >}}

{{<embed src="sector.id" lang="go" >}}


{{<hint warning >}}
TODO:

- describe sizing ranges of sectors
- describe "storage/shipping container" analogy
{{</hint >}}
