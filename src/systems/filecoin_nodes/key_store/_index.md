---
title: Key Store
statusIcon: 🛑
---

The `Key Store` is a fundamental abstraction in any full Filecoin node used to store the keypairs associated to a given miner's address and distinct workers (should the miner choose to run multiple workers).

Node security depends in large part on keeping these keys secure. To that end we recommend keeping keys separate from any given subsystem and using a separate key store to sign requests as required by subsystems as well as keeping those keys not used as part of mining in cold storage.

{{< readfile file="key_store.id" code="true" lang="go" >}}
{{< readfile file="key_store.go" code="true" lang="go" >}}

TODO:

- describe the different types of keys used in the protocol and their usage
- clean interfaces for getting signatures for full filecoin mining cycles
- potential reccomendations or clear disclaimers with regards to consequences of failed key security
- protocol for changing worker keys in filecoin