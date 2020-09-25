---
title: Tipset
weight: 2
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# Tipset

Expected Consensus probabilistically elects multiple leaders in each epoch meaning a Filecoin chain may contain zero or multiple blocks at each epoch (one per elected miner). Blocks from the same epoch are assembled into tipsets. The [VM Interpreter](interpreter) modifies the Filecoin state tree by executing all messages in a tipset (after de-duplication of identical messages included in more than one block). 

Each block references a parent tipset and validates _that tipset's state_, while proposing messages to be included for the current epoch. The state to which a new block's messages apply cannot be known until that block is incorporated into a tipset. It is thus not meaningful to execute the messages from a single block in isolation: a new state tree is only known once all messages in that block's tipset are executed. 
 
A valid tipset contains a non-empty collection of blocks that have distinct miners and all specify identical:

- `Epoch` 
- `Parents`
- `ParentWeight`
- `StateRoot`
- `ReceiptsRoot`

The blocks in a tipset are canonically ordered by the lexicographic ordering of the bytes in each block's ticket, breaking ties with the bytes of the CID of the block itself.

Due to network propagation delay, it is possible for a miner in epoch N+1 to omit valid blocks mined at epoch N from their parent tipset. This does not make the newly generated block invalid, it does however reduce its weight and chances of being part of the canonical chain in the protocol as defined by EC's [Chain Selection](expected_consensus#chain-selection) function.

Block producers are expected to coordinate how they select messages for inclusion in blocks in order to avoid duplicates and thus maximize their expected earnings from transaction fees (see [Message Pool](message_pool)).

The main Tipset structure in the Lotus implementation includes the following:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/tipset.go"  lang="go" symbol="TipSet">}}


Semantic validation of a Tipset is performed in:

{{<embed src="https://github.com/filecoin-project/lotus/blob/master/chain/types/tipset.go"  lang="go" symbol="NewTipSet">}}

Semantic validation includes the following checks:
- The tipset should be composed of at least one block
- All blocks within a tipset should have the same height
- All blocks within a tipset have the same parents, both in terms of number of parents and in terms of CIDs.
