---
title: Message Pool
entries:
- message_syncer
- message_storage
---

{{<label message_pool>}}
The Message Pool is a subsystem in the Filecoin blockchain system. The message pool is acts as the interface between Filecoin nodes and a peer-to-peer network used for off-chain message transmission. It is used by nodes to maintain a set of messages to transmit to the Filecoin VM (for "on-chain" execution).

{{< readfile file="message_pool_subsystem.id" code="true" lang="go" >}}

Clients that use a message pool include:

- storage market provider and client nodes - for transmission of deals on chain
- storage miner nodes - for transmission of PoSts, sector commitments, deals, and other operations tracked on chain
- verifier nodes - for transmission of potential faults on chain
- relayer nodes - for forwarding and discarding messages appropriately.

The message pool subsystem is made of two components:

- The message syncer {{<sref message_syncer>}} -- which receives and propagates messages.
- Message storage {{<sref message_storage>}} -- which caches messages according to a given policy.

TODOs:

- discuss how messages are meant to propagate slowly/async
- explain algorithms for choosing profitable txns
