---
title: Block Producer
weight: 4
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Block Producer

## Mining Blocks

A miner registered with the storage power actor may begin generating and checking election tickets if it has proven storage that meets the [Minimum Miner Size](storage_power_consensus#minimum-miner-size) threshold requirement. 

In order to do so, the miner must be running chain validation, and be keeping track of the most recent blocks received. A miner's new block will be based on parents from the previous epoch.

### Block Creation

Producing a block for epoch `H` requires waiting for the beacon entry for that epoch and using it to run `GenerateElectionProof`. If `WinCount` ≥ 1 (i.e., when the miner is elected), the same beacon entry is used to run `WinningPoSt`. Armed by the `ElectionProof` ticket (output of `GenerateElectionProof`) and the `WinningPoSt` proof, the miner can produce an new block.

See [VM Interpreter](interpreter) for details of parent tipset evaluation, and [Block](block) for constraints on valid block header values. 

To create a block, the eligible miner must compute a few fields:

- `Parents` - the CIDs of the parent tipset's blocks.
- `ParentWeight` - the parent chain's weight (see [Chain Selection](expected_consensus#chain-selection)).
- `ParentState` - the CID of the state root from the parent tipset state evaluation (see the [VM Interpreter](interpreter)).
- `ParentMessageReceipts` - the CID of the root of an AMT containing receipts produced while computing `ParentState`.
- `Epoch` - the block's epoch, derived from the `Parents` epoch and the number of epochs it took to generate this block.
- `Timestamp` - a Unix timestamp, in seconds, generated at block creation.
- `BeaconEntries` - a set of drand entries generated since the last block (see [Beacon Entries](storage_power_consensus#beacon-entries)).
- `Ticket` - a new ticket generated from that in the prior epoch (see [Ticket Generation](storage_power_consensus#randomness-ticket-generation)).
- `Miner` - the block producer's miner actor address.
- `Messages` - The CID of a `TxMeta` object containing message proposed for inclusion in the new block:
  - Select a set of messages from the mempool to include in the block, satisfying block size and gas limits
  - Separate the messages into BLS signed messages and secpk signed messages
  - `TxMeta.BLSMessages`: The CID of the root of an AMT comprising the bare `UnsignedMessage`s
  - `TxMeta.SECPMessages`: the CID of the root of an AMT comprising the `SignedMessage`s
- `BeaconEntries`: a list of beacon entries to derive randomness from
- `BLSAggregate` - The aggregated signature of all messages in the block that used BLS signing.
- `Signature` - A signature with the miner's worker account private key (must also match the ticket signature) over the block header's serialized representation (with empty signature).
- `ForkSignaling` - A uint64 flag used as part of signaling forks. Should be set to 0 by default.

Note that the messages to be included in a block need not be evaluated in order to produce a valid block.
A miner may wish to speculatively evaluate the messages anyway in order to optimize for including
messages which will succeed in execution and pay the most gas.

The block reward is not evaluated when producing a block. It is paid when the block is included in a tipset in the following epoch.

The block's signature ensures integrity of the block after propagation, since unlike many PoW blockchains, 
a winning ticket is found independently of block generation.

### Block Broadcast

An eligible miner propagates the completed block to the network using the [GossipSub](gossip_sub) `/fil/blocks` topic and, assuming everything was done correctly, 
the network will accept it and other miners will mine on top of it, earning the miner a block reward.

Miners should output their valid block as soon as it is produced, otherwise they risk other miners receiving the block after the EPOCH_CUTOFF and not including them in the current epoch.

## Block Rewards

Block rewards are handled by the [Reward Actor](sysactors#rewardactor).
Further details on the Block Reward are discussed in the [Filecoin Token](filecoin_token) section and details about the Block Reward Collateral are discussed in the [Miner Collaterals](miner_collaterals) section.
