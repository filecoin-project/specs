---
title: "Reserved Ranges"
weight: 2
dashboardWeight: 0.2
dashboardState: incomplete
dashboardAudit: 0
---

# Reserved Ranges
---

## Actor ID Reserved Ranges

| Actor                | ID |
|---|---|
| SystemActor          | 0 |
| InitActor            | 1 |
| RewardActor          | 2 |
| CronActor            | 3 |
| StoragePowerActor    | 4 |
| StorageMarketActor   | 5 |
| BurntFundsActor       | 99 |

All values below 100 are reserved for singleton actors. The first non-singleton actor starts at 100.

## Method Reserved Ranges

| Method               | ID |
|---|---|
| value send           | 0 |
| constructor          | 1 |

All other positive values are free for actors to use. For the canonical list, see TODO LINK TO ACTOR CODE WHEN DONE.

## Error Codes

{{< hint warning >}}
TODO
{{< /hint >}}
