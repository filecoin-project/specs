---
title: Storage Market On-Chain Components
weight: 2
bookCollapseSection: true
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Storage Market On-Chain Components

## Storage Deals

There are two types of deals in Filecoin markets, storage deals and retrieval deals. Storage deals are recorded on the blockchain and enforced by the protocol. Retrieval deals are off chain and enabled by a micropayment channel between transacting parties (see [Retrieval Market](retrieval_market) for more information).

The lifecycle of a Storage Deal touches several major subsystems, components, and protocols in Filecoin.

This section describes the storage deal data type and provides a technical outline of the deal flow in terms of how all the components interact with each other, as well as the functions they call. For more detail on the off-chain parts of the storage market see the [Storage Market section](storage_market).

## Data Types

{{<embed src="https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/market/deal.go" lang="go">}}
