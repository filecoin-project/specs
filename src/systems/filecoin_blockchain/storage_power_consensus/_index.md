---
title: Storage Power Consensus
entries:
- expected_consensus
- storage_power_actor
---

{{<label storage_power_consensus>}}
The Storage Power Consensus subsystem is the main interface which enables Filecoin nodes to agree on the state of the system. SPC accounts for individual storage miners' effective power over consensus in given chains in its _Power Table_. It also runs _Expected Consensus_ (the underlying consensus algorithm in use by Filecoin), enabling storage miners to run leader election and generate new blocks updating the state of the Filecoin system.

Succinctly, the SPC subsystem offers the following services:
- Access to the _Power Table_ for every subchain, accounting for individual storage miner power and total power on-chain.
- Access to {{<sref expected_consensus>}} for individual storage miners, enabling:
    - Access to verifiable randomness {{<sref tickets>}} as needed in the rest of the protocol.
    - Running  {{<sref leader_election>}} to produce new blocks.
    - Running {{<sref chain_selection>}} across subchains using EC's weighting function. 
    - Identification of {{<sref finality "the most recently finalized tipset">}}, for use by all protocol participants.

Much of the Storage Power Consensus' subsystem functionality is detailed in the code below but we touch upon some of its behaviors in more detail.

{{< readfile file="storage_power_consensus_subsystem.id" code="true" lang="go" >}}

## Distinguishing between storage miners and block miners

There are two ways to earn Filecoin tokens in the Filecoin network: 
- By participating in the {{<sref storage_market>}} as a storage provider and being paid by clients for file storage deals.
- By mining new blocks on the network, helping modify system state and secure the Filecoin consensus mechanism.

We must distinguish between both types of "miners" (storage and block miners). {{<sref leader_election>}} in Filecoin is predicated on a miner's storage power. Thus, while all block miners will be storage miners, the reverse is not necessarily true.

However, given Filecoin's "useful Proof-of-Work" is achieved through file storage (PoRep and PoSt), there is little overhead cost for storage miners to participate in leader election. Such a {{<sref storage_miner_actor>}} need only register with the {{<sref storage_power_actor>}} in order to participate in Expected Consensus and mine blocks.

{{<label power_table>}}
## The Power Table

The portion of blocks a given miner generates through leader election in EC (and so the block rewards they earn) is proportional to their `Power Fraction` over time. That is, a miner whose storage represents 1% of total storage on the network should mine 1% of blocks on expectation.

SPC provides a power table abstraction which tracks miner power (i.e. miner storage in relation to network storage) over time. The power table is updated for new sector commitments (incrementing miner power), when PoSts fail to be put on-chain (decrementing miner power) or for other storage and consensus faults.

An invariant of the storage power consensus subsystem is that all storage in the power table must be verified. That is, miners can only derive power from storage they have already proven to the network. 

In order to achieve this, Filecoin delays updating power for new sector commitments until the first valid PoSt in the next proving period corresponding to that sector.
Conversely, storage faults only lead to power loss once they are detected (up to one proving period after the fault) so miners will mine with no more power than they have used to store data over time.

Put another way, power accounting in the SPC is delayed between storage being proven or faulted, and power being updated in the power table (and so for leader election). This ensures fairness over time.

## Repeated leader election attempts

In the case that no miner is eligible to produce a block in a given round of EC, the storage power consensus subsystem will be called by the block producer to attempt another leader election by incrementing the nonce appended to the ticket drawn from the past in order to attempt to craft a new valid `ElectionProof`