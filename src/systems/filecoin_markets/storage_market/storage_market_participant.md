---
title: "Storage Market Participant"
---

Both `StorageProvider` and `StorageClient` are `StorageMarketParticipant` and the protocol does not explicitly differentiate the two roles. Any party can be a storage provider or client or both at the same time. Storage deal negotiation is expected to happen completely off chain and the request-response style storage deal protocol is to submit agreed-upon storage deal onto the network and gain storage power on chain.

{{< readfile file="storage_market_participant.id" code="true" lang="go" >}}
