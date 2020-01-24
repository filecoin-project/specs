---
title: Key Store
statusIcon: 🛑
---

The `Key Store` is a fundamental abstraction in any full Filecoin node used to store the keypairs associated to a given miner's address and distinct workers (should the miner choose to run multiple workers).

Node security depends in large part on keeping these keys secure. To that end we recommend keeping keys separate from any given subsystem and using a separate key store to sign requests as required by subsystems as well as keeping those keys not used as part of mining in cold storage.

{{< readfile file="key_store.id" code="true" lang="go" >}}

Filecoin storage miners rely on three main components:

- **The miner address** uniquely assigned to a given storage miner actor upon calling `registerMiner()` in the Storage Power Consensus Subsystem. It is a unique identifier for a given storage miner to which its power and other keys will be associated.
- **The owner keypair** is provided by the miner ahead of registration and its public key associated to the miner address. Block rewards and other payments are made to the ownerAddress.
- **The worker keypair** can be chosen and changed by the miner, its public key associated to the miner address. It is used to sign transactions, signatures, etc. It must be a BLS keypair given its use as part of the {{<sref vrf>}}.

While miner addresses are unique, multiple storage miner actors can share an owner public key or likewise a worker public key.

The process for changing the worker keypairs on-chain (i.e. the workerKey associated to a storage miner) is specified in {{<sref storage_miner_actor>}}. Note that this is a two-step process. First a miner stages a change by sending a message to the chain. When received, the key change is staged to occur in twice the randomness lookback parameter number of epochs, to prevent adaptive key selection attacks. 
Every time a worker key is queried, a pending change is lazily checked and state is potentially updated as needed.

TODO:

- potential reccomendations or clear disclaimers with regards to consequences of failed key security