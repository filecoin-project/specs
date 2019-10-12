---
title: "Storage Provider"
---

Both `StorageProvider` and `StorageClient` are `StorageMarketParticipant`. Any party can be a storage provider or client or both at the same time. Storage deal negotiation is expected to happen completely off chain and the request-response style storage deal protocol is to submit agreed-upon storage deal onto the network and gain storage power on chain. `StorageClient` will initiate the storage deal protocol by submitting a `StorageDealProposal` to the `StorageProvider` who will then add the deal data to a `Sector` and commit the sector onto the blockchain. 

{{< readfile file="storage_provider.id" code="true" lang="go" >}}
