---
title: Storage Deal
weight: 1
bookCollapseSection: true
---

# Storage Deals
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties (see {{<link "retrieval_market">}} for more information). 

The lifecycle of a Storage Deal touches several major subsystems, components, and protocols in Filecoin.

This section describes the storage deal data type and a technical outline for deal flow in terms of how all the components involved and the functions they call on each other. For more detail and prose explanations, see {{<link "storage_market">}} and {{<link storage_mining>}}.

## Data Types

{{< embed src="/specs-actors/actors/builtin/market/deal.go" lang="go" >}}

