---
title: "Architecture Diagrams"
---


# Filecoin Systems

<script type="text/javascript">

function statusIndicatorsShow() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.remove('hidden')
  })
  return false; // stop click event
}

function statusIndicatorsHide() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.add('hidden')
  })
  return false; // stop click event
}

</script>


Status Legend:

- 🛑 **Bare** - Very incomplete at this time.
  - **Implementors:** This is far from ready for you.
- ⚠️ **Rough** -- work in progress, heavy changes coming, as we put in place key functionality.
  - **Implementors:** This will be ready for you soon.
- 🔁 **Refining** - Key functionality is there, some small things expected to change. Some big things may change.
  - **Implementors:** Almost ready for you. You can start building these parts, but beware there may be changes still.
- ✅ **Stable** - Mostly complete, minor things expected to change, no major changes expected.
  - **Implementors:** Ready for you. You can build these parts.

[<a href="#" onclick="return statusIndicatorsShow();">Show</a> / <a href="#" onclick="return statusIndicatorsHide();">Hide</a> ] status indicators


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

