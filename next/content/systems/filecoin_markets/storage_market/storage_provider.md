---
title: "Storage Provider"
weight: 3
---

# Storage Provider
---

The `StorageProvider` is a module that handles incoming queries for Asks and proposals for Deals from a `StorageClient`. It also tracks deals as they move through the deal flow, handling off chain actions during the negotiation phases of the deal and ultimately telling the `StorageMarketActor` to publish on chain. The `StorageProvider`'s last action is to handoff a published deal for storage and sealing to the Storage Mining Subsystem. Note that any address registered as a `StorageMarketParticipant` with the `StorageMarketActor` can be used with the `StorageClient`. A single participant can be a client, provider, or both at the same time.

Because most of what a Storage Provider does is respond to actions initiated by a `StorageClient`, most of its public facing methods relate to getting current status on deals, as opposed to initiating new actions. However, a user of the `StorageProvider` module can update the current Ask for the provider.

{{< embed src="storage_provider.id" lang="go" >}}

<!-- # Storage Provider State Machine -->
