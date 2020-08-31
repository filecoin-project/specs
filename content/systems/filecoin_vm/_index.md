---
title: Virtual Machine
description: VM - Virtual Machine
bookCollapseSection: true
weight: 3
dashboardWeight: 2
dashboardState: wip
dashboardAudit: missing
dashboardTests: 0
---

# VM - Virtual Machine

An Actor in the Filecoin Blockchain is the equivalent of the smart contract in the Ethereum Virtual Machine. Actors carry the logic needed in order to submit transactions, proofs and blocks, among other things, to the Filecoin blockchain. Every actor is identified by a unique address.

The Filecoin Virtual Machine (VM) is the system component that is in charge of execution of all actors code. Execution of actors on the Filecoin VM (i.e., on-chain executions) incur a gas cost.

Any operation applied (i.e., executed) on the Filecoin VM produces an output in the form of a _State Tree_ (discussed below). The latest _State Tree_ is the current source of truth in the Filecoin Blockchain. The _State Tree_ is identified by a CID, which is stored in the IPLD store.

```go
type VM struct {
	## The current State Tree
	cstate      *state.StateTree
	
	## The CID (i.e., identifier) of the current State Tree
	base        cid.Cid
	
	## The IPLD store of the CID (of the current State Tree)
	cst         *cbor.BasicIpldStore
	
	## Current BlockStore
	buf         *bufbstore.BufferedBS
	
	## The epoch number where the VM has been called
	blockHeight abi.ChainEpoch
	
	## The method invoking a State Tree change
	inv         *Invoker
	
	## Randomness included in messages submitted to the VM
	rand        Rand

	Syscalls runtime.Syscalls
}
```