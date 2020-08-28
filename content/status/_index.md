---
title: Status
weight: 9

dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Status
---

Each section of the spec must be stable and audited before it is considered done. The state of each section is tracked below. 

- The **State** column indicates the stability as defined in the legend. 
- The **Theory Audit** column shows the date of the last theory audit with a link to the report.
- The **Weight** column is used to highlight the relative criticality of a section against the others.

**Spec Status Legend**
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

**Spec Status Overview**

{{<dashboard-spec>}}

**Spec Stabilization Progess**

This progress bar shows what percentage of the spec sections are considered stable.

{{<dashboard-progress>}}


**Implementations Status**

Known implementations of the filecoin spec are tracked below, with their current CI build status, their test coverage as reported by [codecov.io](https://codecov.io), and a link to their last security audit report where one exists.

{{<dashboard-impl>}}
