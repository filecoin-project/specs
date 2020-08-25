---
title: Introduction
weight: 1

dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Introduction
---

## What is Filecoin?

Filecoin is a *decentralized storage network*, a network of independent storage providers offering storage and retrieval services in a market operated on a blockchain with a native protocol token called FIL.

The *Filecoin Market* is an algorithmic market for storage and retrieval services.
Miners offer their storage capacity in the market and make *storage deals* with clients.
The market is verifiable: storage providers must provide cryptographic proofs that guarantee persistent storage to their clients.

The *Filecoin Blockchain* is a distributed ledger that orders FIL transactions and executes the *Filecoin Virtual Machine*, a state machine that implements core functionalities to operate the Filecoin markets and token transactions.
Miners mantain the blockchain by creating blocks and verifying transactions and earn block rewards.
Differently from other protocols based on computational resources, miners' influence in the network and block reward earnings are proportional to the amount of storage they prove.

This spec documents the logic implemented in the Filecoin Virtual Machine.

## Spec Status

Each section of the spec must be stable and audited before it is considered done. The state of each section is tracked below. 

- The **State** column indicates the stability as defined in the legend. 
- The **Theory Audit** column shows the date of the last theory audit with a link to the report.
- The **Weight** column is used to highlight the relative criticality of a section against the others.

### Spec Status Legend

<table class="Dashboard"">
  <thead>
    <tr>
      <th>Spec state</th>
      <th>Label</th>
    <tr>
  <thead>
  <tbody>
    <tr>
      <td>Final, will not change before mainnet launch</td>
      <td class="text-black bg-stable">Stable</td>
    </tr>
    <tr>
      <td>Correct, but some details are missing</td>
      <td class="text-black bg-incomplete">Incomplete</td>
    </tr>
    <tr>
      <td>Likely to change. Details still being finalised</td>
      <td class="text-black bg-wip">WIP</td>
    </tr>
    <tr>
      <td>Do not follow. Important things have changed</td>
      <td class="text-black bg-incorrect">Incorrect</td>
    </tr>
  </tbody>
</table>

### Spec Status Overview

{{<dashboard-spec>}}

### Spec Stabilization Progess

This progress bar shows what percentage of the spec sections are considered stable.

{{<dashboard-progress>}}


### Implementations Status

Known implementations of the filecoin spec are tracked below, with their current CI build status, their test coverage as reported by [codecov.io](https://codecov.io), and a link to their last security audit report where one exists.

{{<dashboard-impl>}}
