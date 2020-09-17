---
title: State Tree
weight: 2
dashboardWeight: 1.5
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# State Tree

The State Tree is the output of the execution of any operation applied on the Filecoin Blockchain. The on-chain (i.e., VM) state data structure is a map (in the form of a Hash Array Mapped Trie - HAMT) that binds addresses to actor states. The current State Tree function is called by the VM upon every actor method invocation.

{{<embed src="github:filecoin-project/lotus/chain/state/statetree.go"  lang="go" symbol="StateTree">}}