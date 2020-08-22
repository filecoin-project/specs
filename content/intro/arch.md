---
title: "Architecture Diagrams"
audit: 1
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Architecture Diagrams
---

## Overview Diagram
{{< details title="TODO" >}}
- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token
{{< /details >}}


{{< svg src="diagrams/overview1/overview.dot.svg" title="Protocol Overview Diagram" >}}

## Protocol Flow Diagram

{{< svg src="diagrams/sequence/full-deals-on-chain.mmd.svg" title="Deals on Chain" >}}

## Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{< svg src="diagrams/orient/filecoin.dot.svg" title="Protocol Overview Diagram" >}}
