---
title: "Retrieval Client"
---

# Dependencies

- **Host**: A libp2p host (set setup the libp2p protocols)
- **Filecoin Node**: A node implementation to query the chain for pieces and to setup and manage payment channels
- **BlockStore**: Same as one used by data transfer module
- **Datastore**: for storing deal state in a state machine
- **Data Transfer**: V1 only --Module used for transferring payload. Writes to the blockstore.
- **Stored counter**: to generate unique deal IDs
- **Peer resolver**: a content routing interface to discover retrieval miners that have a given Piece, described in the [Peer Resolver](./retrieval_peer_resolver.md) section.

# API

{{< readfile file="retrieval_client.id" code="true" lang="go" >}}
