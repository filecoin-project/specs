---
menuTitle: Storage Deal
statusIcon: 🔁
title: Storage Deals
entries:
- storage_deal_flow
- storage_deal_states
- faults
---

# Storage Deal
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties (see [Retrieval Market](\missing-link) for more information). 

The lifecycle of a Storage Deal touches several major subsystems, components, and protocols in Filecoin.

This section describes the storage deal data type and a technical outline for deal flow in terms of how all the components involved and the functions they call on each other. For more detail and prose explanations, see [Storage Market](\missing-link) and [Storage Mining Subsystem](\missing-link).

# Data Types

{{< hint danger >}}
Issue with readfile
{{< /hint >}}

{{/* < readfile file="/docs/actors/actors/builtin/storage_market/storage_deal.go" code="true" lang="go" > */}}

