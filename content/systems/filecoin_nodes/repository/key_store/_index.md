---
title: Key Store
weight: 2
dashboardWeight: 1
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# Key Store

The `Key Store` is a fundamental abstraction in any full Filecoin node used to store the keypairs associated with a given miner's address (see actual definition further down) and distinct workers (should the miner choose to run multiple workers).

Node security depends in large part on keeping these keys secure. To that end we strongly recommend: 1) keeping keys separate from all subsystem, ii) using a separate key store to sign requests as required by other subsystems, and iii) keeping those keys that are not used as part of mining in cold storage.

<!--
{{<embed src="key_store.id" lang="go" >}}
-->

Filecoin storage miners rely on three main components:

- **The storage miner _actor_ address** is uniquely assigned to a given storage miner actor upon calling `registerMiner()` in the Storage Power Consensus Subsystem. In effect, the storage miner does not have an address itself, but is rather identified by the address of the actor it is tied to. This is a unique identifier for a given storage miner to which its power and other keys will be associated. The `actor value` specifies the address of an already created miner actor.
- **The owner keypair** is provided by the miner ahead of registration and its public key associated with the miner address. Block rewards and other payments are made to the `ownerAddress`.
- **The worker keypair** is the public key associated with the storage miner actor address. It can be chosen and changed by the miner. The worker keypair is used to sign transactions and signatures. It must be a BLS keypair given its use as part of the [Verifiable Random Function](vrf).

Multiple storage miner actors can share one owner public key or likewise a worker public key.

The process for changing the worker keypairs on-chain (i.e. the worker Key associated with a storage miner actor) is specified in [Storage Miner Actor](storage_miner_actor). Note that this is a two-step process. First, a miner stages a change by sending a message to the chain. When received, the key change is staged to occur in twice the randomness lookback parameter number of epochs, to prevent adaptive key selection attacks. 
Every time a worker key is queried, a pending change is lazily checked and state is potentially updated as needed.

Key security is of utmost importance in Filecoin, as is also the case with keys in every blockchain. **Failure to securely store and use keys or exposure of private keys to adversaries can result in the adversary having access to the miner's funds.**
