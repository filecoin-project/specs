---
title: "Faults"
---

There are two main categories of faults in the Filecoin network. 

- ConsensusFaults
- StorageDealFaults

ConsensusFaults are faults that hurt network consensus and StorageDealFaults are faults where data in a `StorageDeal` is not maintained by the providers. `PledgeCollateral` is slashed for ConsensusFaults and `StorageDealCollateral` for StorageDealFaults.

Any misbehavior may result in more than one fault and can lead to slashing on both collaterals. For example, missing a `PoStProof` will incur a penalty on both `PledgeCollateral` and `StorageDealCollateral` if there is the data is stored in a `StorageDeal`.