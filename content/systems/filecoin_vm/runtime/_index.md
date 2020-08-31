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

{{<embed src="/externals/specs-actors/actors/runtime/runtime.go" lang="go" >}}

## `vm/runtime` VM Implementation

{{<embed src="/externals/lotus/chain/vm/runtime.go" lang="go" >}}

## Exit Codes

There are some common runtime exit codes that are shared by different actors.

{{<embed src="/externals/specs-actors/actors/runtime/exitcode/common.go" lang="go" >}}

## VM Gas Cost Constants

{{<embed src="/externals/lotus/chain/vm/gas.go" lang="go" >}}
