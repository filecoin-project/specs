---
menuTitle: Storage Deal
statusIcon: 🔁
title: Storage Deals
entries:
- storage_deal_flow
- storage_deal_states
- faults
---

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties (see {{<sref retrieval_market>}} for more information). 

The lifecycle of a Storage Deal touches several major subsystems, components, and protocols in Filecoin. As such, documenting it in full can be difficult.

This section describes the storage deal data type and a technical outline for deal flow in terms of how all the components involved and the functions they call on each other. For more detail and prose explanations, see {{<sref storage_market>}} and {{<sref storage_mining_subsystem>}}

# Data Types

{{< readfile file="storage_deal.id" code="true" lang="go" >}}
