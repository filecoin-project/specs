---
title: Runtime
weight: 5
bookCollapseSection: true
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: missing
dashboardTests: 0
---

# VM Runtime Environment (Inside the VM)

## Receipts

A `MessageReceipt` contains the result of a top-level message execution. Every syntactically valid and correctly signed message can be included in a block and will produce a receipt from execution. 

A syntactically valid receipt has:

- a non-negative `ExitCode`,
- a non empty `Return` value only if the exit code is zero, and
- a non-negative `GasUsed`.

```go
type MessageReceipt struct {
	ExitCode exitcode.ExitCode
	Return   []byte
	GasUsed  int64
}
```

## `vm/runtime` Actors Interface

The Actors Interface implementation can be found [here](https://github.com/filecoin-project/specs-actors/blob/master/actors/runtime/runtime.go)

## `vm/runtime` VM Implementation

The Lotus implementation of the Filecoin Virtual Machine runtime can be found [here](https://github.com/filecoin-project/lotus/blob/master/chain/vm/runtime.go)

## Exit Codes

There are some common runtime exit codes that are shared by different actors. Their definition can be found [here](https://github.com/filecoin-project/go-state-types/blob/master/exitcode/common.go).
