---
title: Storage Market On-Chain Components
weight: 2
bookCollapseSection: true
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Storage Market On-Chain Components

## Storage Deals

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by micropayment channel by transacting parties (see [Retrieval Market](retrieval_market) for more information). 

The lifecycle of a Storage Deal touches several major subsystems, components, and protocols in Filecoin.

This section describes the storage deal data type and a technical outline for deal flow in terms of how all the components involved and the functions they call on each other. For more detail and prose explanations, see [Storage Market](storage_market) and [Storage Mining](storage_mining).

## Data Types

{{< embed src="github:filecoin-project/specs-actors/actors/builtin/market/deal.go" lang="go" >}}
