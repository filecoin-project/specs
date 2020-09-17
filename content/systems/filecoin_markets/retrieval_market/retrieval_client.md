---
title: "Retrieval Client"
weight: 3
dashboardWeight: 2
dashboardState: stable
dashboardAudit: n/a
dashboardTests: 0
---

# Retrieval Client

## Client Dependencies

The Retrieval Client Depends On The Following Dependencies

- **Host**: A libp2p host (set setup the libp2p protocols)
- **Filecoin Node**: A node implementation to query the chain for pieces and to setup and manage payment channels
- **BlockStore**: Same as one used by data transfer module
- **Data Transfer**: Module used for transferring payload. Writes to the blockstore.

{{<embed src="github:filecoin-project/go-fil-markets/retrievalmarket/client.go"  lang="go" title="Retrieval Client API">}}