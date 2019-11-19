---
title: Tipset
---

EC can elect multiple leaders in a given epoch meaning Filecoin chains can contain multiple blocks at each height (one per winning miner). This greatly increases chain throughput by allowing blocks to propagate through the network of nodes more efficiently but also means miners should coordinate how they select messages for inclusion in their blocks in order to avoid duplicates and maximize their earnings from transaction fees (see {{<sref message_pool>}}).

Accordingly, all valid blocks generated in a round form a `Tipset` that participants will attempt to mine off of in the subsequent round (see above). Tipsets are valid so long as:
- All blocks in a Tipset have the same parent Tipset
- All blocks in a Tipset are mined in the same Epoch

During state computation, blocks in a tipset are processed in order of block ticket, breaking ties with the block CID bytes. The Filecoin state tree is modified by the execution of all messages in a given Tipset.

Due to this fact, adding new blocks to the chain actually validates those blocks' parent Tipset, that is: executing the messages of a new block, a miner cannot know exactly what state tree this will yield. That state tree is only known once all messages in that block's Tipset have been executed. Accordingly, it is in the next round (and based on the number of blocks mined on a given Tipset) that a miner will be able to choose which state tree to extend.

Due to network propagation delay, it is possible for a miner in epoch N+1 to omit valid blocks mined at epoch N from their Tipset. This does not make the newly generated block invalid, it does however reduce its weight and chances of being part of the canonical chain in the protocol as defined by EC's {{<sref chain_selection>}} function.

{{<label tipset>}}
The Tipset is a group of blocks in the same exact round, that all share the exact same parents.

{{< readfile file="tipset.id" code="true" lang="go" >}}

{{< readfile file="tipset.go" code="true" lang="go" >}}
