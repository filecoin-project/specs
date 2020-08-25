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


![Protocol Overview Diagram](diagrams/overview1/overview.dot)

## Protocol Flow Diagram

![Deals on Chain](diagrams/sequence/full-deals-on-chain.mmd)

## Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

![Protocol Overview Diagram](diagrams/orient/filecoin.dot)