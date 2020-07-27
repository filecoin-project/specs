---
title: Block Producer
weight: 5
dashboardWeight: 1.5
dashboardState: incorrect
dashboardAudit: 0
dashboardTests: 0
---

# Block Producer
---

## Mining Blocks

A miner registered with the storage power actor may begin generating and checking election tickets if it has proven storage meeting the [Minimum Miner Size](storage_power_consensus#minimum-miner-size) threshold requirement. 

In order to do so, the miner must be running chain validation, and be keeping track of the most recent blocks received. A miner's new block will be based on parents from the previous epoch.

### Block Creation

Producing a block for epoch `H` requires computing a tipset for epoch `H-1` (or possibly a prior epoch,
if no blocks were received for that epoch). Using the state produced by this tipset, a miner can
scratch winning ElectionPoSt ticket(s). 
Armed with the requisite `ElectionPoStOutput`, as well as a new randomness ticket generated in this epoch, a miner can produce a new block.

See [VM Interpreter](interpreter) for details of parent tipset evaluation, and [Block](block) for constraints 
on valid block header values. 

To create a block, the eligible miner must compute a few fields:

- `Parents` - the CIDs of the parent tipset's blocks.
- `ParentWeight` - the parent chain's weight (see [Chain Selection](expected_consensus#chain-selection)).
- `ParentState` - the CID of the state root from the parent tipset state evaluation (see the [VM Interpreter](interpreter)).
- `ParentMessageReceipts` - the CID of the root of an AMT containing receipts produced while computing `ParentState`.
- `Epoch` - the block's epoch, derived from the `Parents` epoch and the number of epochs it took to generate this block.
- `Timestamp` - a Unix timestamp, in seconds, generated at block creation.
- `Ticket` - a new ticket generated from that in the prior epoch (see [Ticket Generation](storage_power_consensus#randomness-ticket-generation)).
- `Miner` - the block producer's miner actor address.
- `ElectionPoStVerifyInfo` - The byproduct of running an ElectionPoSt yielding requisite on-chain information (see [Election Post](election_post)), namely:
  - An array of `PoStCandidate` objects, all of which include a winning partial ticket used to run leader election.
  - `PoStRandomness` used to challenge the miner's sectors and generate the partial tickets.
  - A `PoStProof` snark output to prove that the partial tickets were correctly generated.
- `Messages` - The CID of a `TxMeta` object containing message proposed for inclusion in the new block:
  - Select a set of messages from the mempool to include in the block, satisfying block size and gas limits
  - Separate the messages into BLS signed messages and secpk signed messages
  - `TxMeta.BLSMessages`: The CID of the root of an AMT comprising the bare `UnsignedMessage`s
  - `TxMeta.SECPMessages`: the CID of the root of an AMT comprising the `SignedMessage`s
- `BLSAggregate` - The aggregated signature of all messages in the block that used BLS signing.
- `Signature` - A signature with the miner's worker account private key (must also match the ticket signature) over the block header's serialized representation (with empty signature). 

Note that the messages to be included in a block need not be evaluated in order to produce a valid block.
A miner may wish to speculatively evaluate the messages anyway in order to optimize for including
messages which will succeed in execution and pay the most gas.

The block reward is not evaluated when producing a block. It is paid when the block is included in a tipset in the following epoch.

The block's signature ensures integrity of the block after propagation, since unlike many PoW blockchains, 
a winning ticket is found independently of block generation.

### Block Broadcast

An eligible miner broadcasts the completed block to the network and, assuming everything was done correctly, 
the network will accept it and other miners will mine on top of it, earning the miner a block reward!

Miners should output their valid block as soon as it is produced, otherwise they risk other miners receiving the block after the EPOCH_CUTOFF and not including them.

## Block Rewards

{{< hint warning >}}
TODO: Rework this.
{{</ hint >}}
Over the entire lifetime of the protocol, 1,400,000,000 FIL (`TotalIssuance`) will be given out to miners. Each of the miners who produced a block in a tipset will receive a block reward. 

Note: Due to jitter in EC, and the gregorian calendar, there may be some error in the issuance schedule over time. This is expected to be small enough that it's not worth correcting for. Additionally, since the payout mechanism is transferring from the network account to the miner, there is no risk of minting *too much* FIL.

{{< hint warning >}}
TODO: Ensure that if a miner earns a block reward while undercollateralized, then `min(blockReward, requiredCollateral-availableBalance)` is garnished (transfered to the miner actor instead of the owner).
{{</ hint >}}
