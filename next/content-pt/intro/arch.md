---
title: "Architecture Diagrams"
audit: 1
dashboardState: permanent
dashboardInterface: incorrect
---

# Architecture Diagramsss
---

## Filecoin Systems

Status Legend:

- 🛑 **Bare** - Very incomplete at this time.
  - **Implementors:** This is far from ready for you.
- ⚠️ **Rough** -- work in progress, heavy changes coming, as we put in place key functionality.
  - **Implementors:** This will be ready for you soon.
- 🔁 **Refining** - Key functionality is there, some small things expected to change. Some big things may change.
  - **Implementors:** Almost ready for you. You can start building these parts, but beware there may be changes still.
- ✅ **Stable** - Mostly complete, minor things expected to change, no major changes expected.
  - **Implementors:** Ready for you. You can build these parts.

*Note that the status relates to the state of the spec either written out either in english or in code. The goal is for the spec to eventually be fleshed out in both language-sets.*

## Overview Diagram
{{< details title="TODO" >}}
- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token
{{< /details >}}


{{< svg src="/intro/overview.dot.svg" title="Protocol Overview Diagram" />}}

## Protocol Flow Diagram

{{< mermaid src="/intro/full-deals-on-chain.mmd" title="Deals on Chain"/>}}

## Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{< svg src="/intro/filecoin.dot.svg" title="Protocol Overview Diagram" />}}
