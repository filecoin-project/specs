---
menuTitle: Runtime
statusIcon: ⚠️
title: VM Runtime Environment (Inside the VM)
entries:
- exitcode
- gascost
---

# Receipts

A `MessageReceipt` contains the result of a top-level message execution.

A syntactically valid receipt has:

- a non-negative `ExitCode`,
- a non empty `ReturnValue` only if the exit code is zero,
- a non-negative `GasUsed`.

# `vm/runtime` interface

{{< readfile file="runtime.id" code="true" lang="go" >}}

# `vm/runtime` implementation

{{< readfile file="runtime.go" code="true" lang="go" >}}

# Code Loading

{{< readfile file="codeload.go" code="true" lang="go" >}}
