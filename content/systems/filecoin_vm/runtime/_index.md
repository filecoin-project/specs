---
title: Runtime
weight: 5
bookCollapseSection: true
dashboardAudit: 0
dashboardState: wip
dashboardInterface: wip
---

# VM Runtime Environment (Inside the VM)
---

## Receipts

A `MessageReceipt` contains the result of a top-level message execution.

A syntactically valid receipt has:

- a non-negative `ExitCode`,
- a non empty `ReturnValue` only if the exit code is zero,
- a non-negative `GasUsed`.

## `vm/runtime` interface

{{<embed src="/specs-actors/actors/runtime/runtime.go" lang="go" >}}

## `vm/runtime` implementation

{{<embed src="impl/runtime.go" lang="go" >}}

## Code Loading

{{<embed src="impl/codeload.go" lang="go" >}}

## Exit codes

{{<embed src="/specs-actors/actors/runtime/exitcode/common.go" lang="go" >}}
