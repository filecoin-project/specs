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

{{<label block_validation>}}
## Block Validation Rules
  - **BV0 - Syntactic Validation**: Validate data structure packing and ensure correct typing.
  - **BV1 - Future Epoch Validation**: Validate that block does not claim to be from a future epoch. 
  - **BV2 - Light Consensus State Checks**: Validate `b.ChainWeight`, `b.ChainEpoch`, `b.MinerAddress`, `b.Timestamp`, are plausible (some ranges of bad values can be detected easily, especially if we have the state of the chain at `b.ChainEpoch - consensus.LookbackParameter`. Eg Weight and Epoch have well defined valid ranges, and `b.MinerAddress`
  must exist in the lookback state). This requires some chain state, enough to establish plausibility levels of each of these values. A node can determine the exact valid range for `b.Timestamp`  based on the `LastTrustedCheckpoint`, `b.ChainEpoch`, and a synchronized clock (see {{<sref clock>}}).
  - **BV3 - Signature Validation**: Verify `b.BlockSig` is correct.
  - **BV4 - Verify ElectionPoSt**: Verify `b.ElectionPoStOutput` was correctly generated and yielding winning `PartialTickets`.
  - **BV5 - Verify Ancestry links to chain**: Verify ancestry links back to trusted blocks. If the ancestry forks off before finality, or does not connect at all, it is a bad block.
  - **BV6 - Verify MessageSigs**: Verify the signatures on messages
  - **BV7 - Verify StateTree**: Verify the application of `b.Parents.Messages()` correctly produces `b.StateTree` and `b.MessageReceipts`

{{< readfile file="block.id" code="true" lang="go" >}}

{{< readfile file="block.go" code="true" lang="go" >}}

