---
title: "Faults"
weight: 3
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Faults
---

There are two main categories of faults in the Filecoin network. 

- ConsensusFaults
- StorageDealFaults

ConsensusFaults are faults that impact network consensus and StorageDealFaults are faults where data in a `StorageDeal` is not maintained by the providers pursuant to deal terms.

[Pledge Collateral](storage_power_actor#pledge-collateral) is slashed for ConsensusFaults and [Storage Deal Collateral](storage_deal) for StorageDealFaults.

Any misbehavior may result in more than one fault thus lead to slashing on both collaterals. For example, missing a `PoStProof` will incur a penalty on both `PledgeCollateral` and `StorageDealCollateral` given it impacts both a given `StorageDeal` and power derived from the sector commitments in [Storage Power Consensus](storage_power_consensus).

## Storage Faults

{{<hint warning>}}
TODO: complete this.
{{</hint>}}