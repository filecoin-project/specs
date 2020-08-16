---
title: "Market"
description: "Markets in Filecoin"
bookCollapseSection: true
weight: 7
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 0
dashboardTests: 0
---

# Markets
---

Filecoin is a consensus protocol, a data-storage platform, and a marketplace for storing and retrieving data. There are two major components to Filecoin markets, the storage market and the retrieval market. While storage and retrieval negotiations for both the storage and the retrieval markets are taking place primarily *off the blockchain* (at least in the current version of Filecoin), storage deals made in the storage market will be published on-chain and will be enforced by the protocol. Storage deal negotiation and order matching are expected to happen off-chain in the first version of Filecoin. Retrieval deals are also negotiated off-chain and executed with micropayments between transacting parties in payment channels.

Even though most of the market actions happen off the blockchain, there are on-chain invariants that create economic structure for network success and allow for positive emergent behavior. You can read more about the relationship between on-chain deals and storage power in [Storage Power Consensus](storage_power_consensus).

## Status Overview

{{< dashboard-level name="Market" open="true">}}