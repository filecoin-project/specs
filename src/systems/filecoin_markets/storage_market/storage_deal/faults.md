---
title: "Faults"
---

{{<label storage_faults>}}
There are two main categories of faults in the Filecoin network. 

- ConsensusFaults
- StorageDealFaults

ConsensusFaults are faults that impact network consensus and StorageDealFaults are faults where data in a `StorageDeal` is not maintained by the providers pursuant to deal terms.

{{<sref pledge_collateral>}} is slashed for ConsensusFaults and {{<sref storage_deal_collateral>}} for StorageDealFaults.

Any misbehavior may result in more than one fault thus lead to slashing on both collaterals. For example, missing a `PoStProof` will incur a penalty on both `PledgeCollateral` and `StorageDealCollateral` given it impacts both a given `StorageDeal` and power derived from the sector commitments in {{<sref storage_power_consensus>}}.

## Storage Faults
TODO: complete this.