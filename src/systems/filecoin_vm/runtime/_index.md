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

{{< readfile file="/docs/actors/runtime/runtime.id" code="true" lang="go" >}}

# `vm/runtime` implementation

{{< readfile file="impl/runtime.go" code="true" lang="go" >}}

# Code Loading

{{< readfile file="impl/codeload.go" code="true" lang="go" >}}

# Exit codes

{{< readfile file="/docs/actors/runtime/vm_exitcodes.id" code="true" lang="go" >}}

{{< readfile file="/docs/actors/runtime/vm_exitcodes.go" code="true" lang="go" >}}
