---
title: "Storage Client"
---

# Storage Client
---

The `StorageClient` is a module that discovers miners, determines their asks, and proposes deals to `StorageProviders`. It also tracks deals as they move through the deal flow. Note that any address registered as a `StorageMarketParticipant` with the `StorageMarketActor` can be used with the `StorageClient`. A single participant can be a client, provider, or both at the same time.

--

{{< hint danger >}}
Issue with readfile
{{< /hint >}}
{{/* < readfile file="storage_client.id" code="true" lang="go" > */}}

<!-- # Storage Client State Machine -->
