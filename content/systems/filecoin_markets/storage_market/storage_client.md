---
title: Storage Client
weight: 4
dashboardWeight: 2
dashboardState: incomplete
dashboardAudit: 1
dashboardTests: 0
---

# Storage Client
---

The `StorageClient` is a module that discovers miners, determines their asks, and proposes deals to `StorageProviders`. It also tracks deals as they move through the deal flow. Note that any address registered as a `StorageMarketParticipant` with the `StorageMarketActor` can be used with the `StorageClient`.

Recall that a single participant can be a `StorageClient`, `StorageProvider`, or both at the same time.

{{< embed src="storage_client.id" lang="go" >}}

<!-- # Storage Client State Machine -->
