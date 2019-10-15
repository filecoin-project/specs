---
title: "Architecture Diagrams"
---


# Filecoin Systems

{{< incTocMap "/docs/systems" 2 "colorful" >}}


# Overview Diagram

TODO:

- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token

{{< diagram src="../diagrams/overview1/overview.dot.svg" title="Protocol Overview Diagram" >}}


# Protocol Flow Diagram -- deals off chain

{{< diagram src="../diagrams/sequence/full-deals-off-chain.mmd.svg" title="Protocol Sequence Diagram - Deals off Chain" >}}

# Protocol Flow Diagram -- deals on chain

{{< diagram src="../diagrams/sequence/full-deals-on-chain.mmd.svg" title="Protocol Sequence Diagram - Deals on Chain" >}}

# Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{< diagram src="../diagrams/orient/filecoin.dot.svg" title="Parameter Calculation Dependency Graph" >}}

