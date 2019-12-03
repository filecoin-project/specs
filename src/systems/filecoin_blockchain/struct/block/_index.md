---
title: Block
---

{{<label block>}}
The Block is a unit of the Filecoin blockchain.

A block header contains information relevant to a particular point in time over which the network may achieve consensus.

{{% notice note %}}
**Note:** A block is functionally the same as a block header in the Filecoin protocol. While a block header contains Merkle links to the full system state, messages, and message receipts, a block can be thought of as the full set of this information (not just the Merkle roots, but rather the full data of the state tree, message tree, receipts tree, etc.). Because a full block is quite large, our chain consists of block headers rather than full blocks. We often use the terms `block` and `block header` interchangeably.
{{% /notice %}}

{{< readfile file="block.id" code="true" lang="go" >}}

{{< readfile file="block.go" code="true" lang="go" >}}

