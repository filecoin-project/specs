---
title: "Architecture Diagrams"
audit: 1
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Architecture Diagrams

Actor State Diagram

{{< mermaid >}}
stateDiagram
    Null --> Precommitted: PreCommitSectors
    Precommitted --> Committed: CommitSectors
    Precommitted --> Deleted: CronPreCommitExpiry (PCD)
    Committed --> Active: SubmittedWindowPoSt
    Committed --> Faulty: DeclareFault\nSubmitWindowPoSt (SP)\nProvingDeadline (SP)
    Committed --> Terminated: TerminateSectors\n(TF)
    Faulty --> Active: SubmittedWindowPoSt (FF)
    Faulty --> Faulty: ProvingDeadline (FF)
    Faulty --> Recovering: DeclareFaultRecovered
    Faulty --> Terminated: EarlyExpiration (TF)\nTerminateSectors (TF)
    Recovering --> Active: SubmittedWindowPoSt (FF)
    Recovering --> Faulty: DeclareFault\nProvingDeadline (SP)
    Recovering --> Terminated: TerminateSectors (TF)
    Active --> Active: SubmittedWindowPoSt
    Active --> Faulty: DeclareFault\nSubmitWindowPoSt (SP)\nProvingDeadline (SP)
    Active --> Terminated: CronExpiration\nTerminateSectors (TF)
    Terminated --> Deleted: CompactSectors
{{</ mermaid >}}
