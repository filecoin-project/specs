---
title: Block
weight: 1
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# Block

The Block is the main unit of the Filecoin blockchain.

The Block structure in the Filecoin blockchain is composed of: i) the Block Header, ii) the list of messages inside the block, and iii) the signed messages. This is represented inside the `FullBlock` abstraction. The messages indicate the required set of changes to apply in order to arrive at a deterministic state of the chain.

The Lotus implementation of the block has the following `struct`:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/fullblock.go"  lang="go" symbol="FullBlock">}}

> **Note**  
> A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is large in size, the Filecoin blockchain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.

A `BlockHeader` is a canonical representation of a block. BlockHeaders are propagated between miner nodes. From the blockcheader message, a miner has all the required information to apply the associated `FullBlock`'s state and update the chain. In order to be able to do this, the minimum set of information items that need to be included in the `BlockHeader` are shown below and include among others: the miner's address, the Ticket, the [Proof of SpaceTime](post), the CID of the parents where this block evolved from in the IPLD DAG, as well as the messages' own CIDs.

The Lotus implementation of the block header has the following `struct`s:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/blockheader.go"  lang="go" symbol="BlockHeader">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/blockheader.go"  lang="go" symbol="Ticket">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/electionproof.go"  lang="go" symbol="ElectionProof">}}

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/blockheader.go"  lang="go" symbol="BeaconEntry">}}

The `BlockHeader` structure has to refer to the TicketWinner of the current round which ensures the correct winner is passed to [ChainSync](chainsync).

```go
func IsTicketWinner(vrfTicket []byte, mypow BigInt, totpow BigInt) bool
```

The `Message` structure has to include the source (`From`) and destination (`To`) addresses, a `Nonce` and the `GasPrice`.

The Lotus implementation of the message has the following structure:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/message.go"  lang="go" symbol="Message">}}


The message is also validated before it is passed to the [chain synchronization logic](chainsync):

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/message.go"  lang="go" symbol="ValidForBlockInclusion">}}


## Block syntax validation

Syntax validation refers to validation that should be performed on a block and its messages _without_ reference to outside information such as the parent state tree. This type of validation is sometimes called _static validation_.

An invalid block must not be transmitted or referenced as a parent.

A syntactically valid block header must decode into fields matching the definitions below,  must be a valid CBOR PubSub `BlockMsg` message and must have:

- between 1 and `5*ec.ExpectedLeaders` `Parents` CIDs if `Epoch` is greater than zero (else empty `Parents`),
- a non-negative `ParentWeight`,
- less than or equal to `BlockMessageLimit` number of messages,
- aggregate message CIDs, encapsulated in the `MsgMeta` structure, serialized to the `Messages` CID in the block header,
- a `Miner` address which is an ID-address. The Miner `Address` in the block header should be present and correspond to a public-key address in the current chain state.
- Block signature (`BlockSig`) that belongs to the public-key address retrieved for the Miner
- a non-negative `Epoch`,
- a positive `Timestamp`,
- a `Ticket` with non-empty `VRFResult`,
- `ElectionPoStOutput` containing:
   - a `Candidates` array with between 1 and `EC.ExpectedLeaders` values (inclusive),
   - a non-empty `PoStRandomness` field,
   - a non-empty `Proof` field,
- a non-empty `ForkSignal` field.
   
A syntactically valid full block must have:

- all referenced messages syntactically valid,
- all referenced parent receipts syntactically valid,
- the sum of the serialized sizes of the block header and included messages is no greater than `block.BlockMaxSize`,
- the sum of the gas limit of all explicit messages is no greater than `block.BlockGasLimit`.

Note that validation of the block signature requires access to the miner worker address and public key from the parent tipset state, so signature validation forms part of semantic validation. Similarly, message signature validation requires lookup of the public key associated with each message's `From` account actor in the block's parent state.

## Block semantic validation

Semantic validation refers to validation that requires reference to information outside the block header and messages themselves. Semantic validation relates to the parent tipset and state on which the block is built.

In order to proceed to semantic validation the `FullBlock` must be assembled from the received block header retrieving its Filecoin messages. Block message CIDs can be retrieved from the network and be decoded into valid CBOR `Message`/`SignedMessage`.

In the Lotus implementation the semantic validation of a block is carried out by the `Syncer` module:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/sync.go"  lang="go" symbol="ValidateBlock">}}


Messages are retrieved through the `Syncer`. There are the following two steps followed by the `Syncer`:
1) Assemble a `FullTipSet` populated with the single block received earlier. The Block's `ParentWeight` is greater than the one from the (first block of the) heaviest tipset.
2) Retrieve all tipsets from the received Block down to our chain. Validation is expanded to every block inside these tipsets. The validation should ensure that:
	- Beacon entires are ordered by their round number.
	- The Tipset `Parents` CIDs match the fetched parent tipset through BlockSync.


A semantically valid block must meet all of the following requirements.

**`Parents`-Related**
- `Parents` listed in lexicographic order of their header's `Ticket`.
- `ParentStateRoot` CID of the block matches the state CID computed from the parent [Tipset](tipset).
- `ParentState` matches the state tree produced by executing the parent tipset's messages (as defined by the VM interpreter) against that tipset's parent state.
- `ParentMessageReceipts` identifying the receipt list produced by parent tipset execution, with one receipt for each unique message from the parent tipset. In other words, the Block's `ParentMessageReceipts` CID matches the receipts CID computed from the parent tipset.
- `ParentWeight` matches the weight of the chain up to and including the parent tipset.

**Time-Related**
- `Epoch` is greater than that of its `Parents`, and 
    - not in the future according to the node's local clock reading of the current epoch,
        - blocks with future epochs should not be rejected, but should not be evaluated (validated or included in a tipset) until the appropriate epoch
    - not farther in the past than the soft finality as defined by SPC [Finality](expected_consensus#finality-in-ec),
        - this rule only applies when receiving new gossip blocks (i.e. from the current chain head), not when syncing to the chain for the first time.
- The `Timestamp` included is in seconds that:
  - must not be bigger than current time plus `ΑllowableClockDriftSecs`
  - must not be smaller than previous block's `Timestamp` plus `BlockDelay` (including null blocks)
  - is of the precise value implied by the genesis block's timestamp, the network's Βlock time and the Βlock's `Epoch`.

**`Miner`-Related**
- The `Miner` is active in the storage power table in the parent tipset state. The Miner's address is registered in the `Claims` HAMT of the Power Actor
- The `TipSetState` should be included for each tipset being validated.
	- Every Block in the tipset should belong to different a miner.
- The Actor associated with the message's `From` address exists, is an account actor and its Nonce matches the message Nonce.
- Valid proofs that the Miner proved access to sealed versions of the sectors it was challenged for are included. In order to achieve that:
	- draw randomness for current epoch with `WinningPoSt` domain separation tag.
	- get list of sectors challanged in this epoch for this miner, based on the randomness drawn.
- Miner is not slashed in `StoragePowerActor`.


**`Beacon`- & `Ticket`-Related**
- Valid `BeaconEntries` should be included:
	- Check that every one of the `BeaconEntries` is a signature of a message: `previousSignature || round` signed using DRAND's public key.
	- All entries between `MaxBeaconRoundForEpoch` down to `prevEntry` (from previous tipset) should be included.
- A `Ticket` derived from the minimum ticket from the parent tipset's block headers, 
    - `Ticket.VRFResult` validly signed by the `Miner` actor's worker account public key,
- `ElectionProof Ticket` is computed correctly by checking BLS signature using miner's key. The `ElectionProof` ticket should be a winning ticket.

**Message- & Signature-Related**
- `secp256k1` messages are correctly signed by their sending actor's (`From`) worker account key,
- A `BLSAggregate` signature is included that signs the array of CIDs of all the BLS messages referenced by the block with their sending actor's key.
- A valid `Signature` over the block header's fields from the block's `Miner` actor's worker account public key is included.
- For each message in `ValidForBlockInclusion()` the following hold:
	-  Message fields `Version`, `To`, `From`, `Value`, `GasPrice`, and `GasLimit` are correctly defined.
	- Message `GasLimit` is under the message minimum gas cost (derived from chain height and message length).
- For each message in `ApplyMessage` (that is before a message is executed), the following hold:
	- Basic gas and value checks in `checkMessage()`:
		- The Message `GasLimit` is bigger than zero.
		- The Message `GasPrice` and `Value` are set.
	- The Message's storage gas cost is under the message's `GasLimit`.
	- The Message's `Nonce` matches the nonce in the Actor retrieved from the message's `From` address.
	- The Message's maximum gas cost (derived from its `GasLimit`, `GasPrice`, and `Value`) is under the balance of the Actor retrieved from message's `From` address.
	- The Message's transfer `Value` is under the balance of the Actor retrieved from message's `From` address.

There is no semantic validation of the messages included in a block beyond validation of their signatures.
If all messages included in a block are syntactically valid then they may be executed and produce a receipt. 

A chain sync system may perform syntactic and semantic validation in stages in order to minimize unnecessary resource expenditure.

If all of the above tests are successful, the block is marked as validated. Ultimately, an invalid block must not be propagated further or validated as a parent node.
