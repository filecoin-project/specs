---
title: "Retrieval Client"
---

### Client Dependencies

The Retrieval Client Depends On The Following Dependencies

- **Host**: A libp2p host (set setup the libp2p protocols)
- **Filecoin Node**: A node implementation to query the chain for pieces and to setup and manage payment channels
- **BlockStore**: Same as one used by data transfer module
- **Data Transfer**: V1 only --Module used for transferring payload. Writes to the blockstore.

### API

{{< readfile file="retrieval_client.id" code="true" lang="go" >}}
