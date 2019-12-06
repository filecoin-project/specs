---
title: Block
---

{{<label block>}}
# Block
The Block is a unit of the Filecoin blockchain.

A block header contains information relevant to a particular point in time over which the network may achieve consensus.

{{% notice note %}}
**Note:** A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is quite large, our chain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.
{{% /notice %}}

{{< readfile file="block.id" code="true" lang="go" >}}

{{< readfile file="election.id" code="true" lang="go" >}}

# Block syntax validation

Syntax validation refers to validation that may be performed on a block and its messages 
without reference to outside information such as the parent state tree.

An invalid block must not be transmitted or referenced as a parent.

A syntactically valid block header must decode into fields matching the type definition below. 

A syntactically valid header must have:

- a non-empty array of `Parents` CIDs if `Epoch` is greater than zero,
- a non-negative `ParentWeight`,
- a `Miner` address which is an ID-address,
- a non-negative `Epoch`,
- a positive `Timestamp`,
- a `Ticket` with non-empty `VRFResult`,
- `ElectionPoStOutput` containing:
   - a non-empty `Candidates` array
   - a non-empty `PoStRandomness` field
   - a non-empty `Proof` field
   
A syntactically valid full block must have:

- all referenced messages syntactically valid,
- all referenced parent receipts syntactically valid,

Note that validation of the block signature requires access to the miner worker address and
public key from the parent tipset state, so signature validation forms part of semantic validation. 
Similarly, message signature validation requires lookup of the public key associated with 
each message's `From` account actor in the block's parent state.

# Block semantic validation

Semantic validation refers to validation that requires reference to information outside the block
header and messages themselves, in particular the parent tipset and state on which the block is built.

A semantically valid block must have:

- `Parents` listed in lexicographic order of their header's `Ticket`,
- `Parents` all reference valid blocks and form a valid {{<sref tipset>}},
- `ParentState` matching the state tree produced by executing the parent tipset's messages (as defined by the VM interpreter) against that tipset's parent state,
- `ParentMessageReceipts` identifying the receipt list produced by parent tipset execution, with one receipt for each unique message from the parent tipset, 
- `ParentWeight` matching the weight of the chain up to and including the parent tipset,
- `Epoch` greater than that of its parents, and not in the future according to the node's local clock reading of the current epoch,
- `Miner` that is active in the storage power table in the parent tipset state,  
- a `Ticket` derived from the minimum ticket from the parent tipset's block headers, 
    - `Ticket.VRFResult` validly signed by the `Miner` actor's worker account public key,
- `ElectionPoStOutput` yielding winning partial tickets that were generated validly, 
  - `ElectionPoSt.Randomness` is well formed and appropriately drawn from a past tipset according to the PoStLookback,
  - `ElectionPoSt.Proof` is a valid proof verifying the generation of the `ElectionPoSt.Candidates` from the `Miner`'s eligible sectors,
  - `ElectionPoSt.Candidates` contains well formed `PoStCandidate`s each of which has a `PartialTicket` yielding a winning `ChallengeTicket` in Expected Consensus.
- a `Timestamp` in seconds lying within the quantized epoch window implied by the genesis block's timestamp and the block's `Epoch`,
- all SECP messages correctly signed by their sending actor's worker account key,
- a `BLSAggregate` signature that signs the array of CIDs of the BLS messages referenced by the block 
with their sending actor's key.
- a valid `Signature` over the block header's fields from the block's `Miner` actor's worker account public key.

There is no semantic validation of the messages included in a block beyond validation of their signatures.
If all messages included in a block are syntactically valid then they may be executed and produce a receipt. 

A chain sync system may perform syntactic and semantic validation in stages in order to minimize unnecessary resource expenditure.



