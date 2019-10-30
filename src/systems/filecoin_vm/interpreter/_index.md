---
menuTitle: Interpreter
statusIcon: ⚠️
title: VM Interpreter - Message Invocation (Outside VM)
entries:
- vm_outside
- vm_inside
# suppressMenu: true
---

(You can see the _old_ VM interpreter [here](docs/systems/filecoin_vm/vm_interpreter_old) )

# `vm/interpreter` interface

{{< readfile file="vm_interpreter.id" code="true" lang="go" >}}

# `vm/interpreter` implementation

{{< readfile file="vm_interpreter.go" code="true" lang="go" >}}

# `vm/interpreter/registry`

{{< readfile file="vm_registry.go" code="true" lang="go" >}}
