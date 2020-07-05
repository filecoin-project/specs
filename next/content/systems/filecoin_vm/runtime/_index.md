---
menuTitle: Runtime
statusIcon: 🔁
title: VM Runtime Environment (Inside the VM)
entries:
- exitcode
- gascost
---

# VM Runtime Environment (Inside the VM)
---

# Receipts

A `MessageReceipt` contains the result of a top-level message execution.

A syntactically valid receipt has:

- a non-negative `ExitCode`,
- a non empty `ReturnValue` only if the exit code is zero,
- a non-negative `GasUsed`.

# `vm/runtime` interface

{{< hint danger >}}
Issue with readfile
{{< /hint >}}
{{/* < readfile file="/docs/actors/actors/runtime/runtime.go" code="true" lang="go" > */}}

# `vm/runtime` implementation

{{< hint danger >}}
Issue with readfile
{{< /hint >}}
{{/* < readfile file="impl/runtime.go" code="true" lang="go" > */}}

# Code Loading

{{< hint danger >}}
Issue with readfile
{{< /hint >}}
{{/* < readfile file="impl/codeload.go" code="true" lang="go" > */}}

# Exit codes

{{< hint danger >}}
Issue with label
{{< /hint >}}
{{/* < readfile file="/docs/actors/actors/runtime/exitcode/vm_exitcodes.go" code="true" lang="go" > */}}
