---
title: "Retrieval Provider (Miner)"
weight: 4
dashboardWeight: 2
dashboardState: stable
dashboardAudit: missing
dashboardTests: 0
---

# Retrieval Provider (Miner)

## Provider Dependencies

The Retrieval Provider depends on the following dependencies

- **Host**: A libp2p host (set setup the libp2p protocols)
- **Filecoin Node**: A node implementation to query the chain for pieces and to setup and manage payment channels
- **StorageMining Subsystem**: For unsealing sectors
- **BlockStore**: Same as one used by data transfer module
- **Data Transfer**: Module used for transferring payload. Reads from the blockstore.

{{<embed src="github:filecoin-project/go-fil-markets/retrievalmarket/provider.go"  lang="go" title="Retrieval Provider API">}}