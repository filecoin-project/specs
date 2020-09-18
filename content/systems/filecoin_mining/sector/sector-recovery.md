---
title: Sector Recovery
weight: 6
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Sector Recovery

Miners should try to recover faulty sectors in order to avoid paying the penalty, which is approximately equal to the block reward that the miner would receive from that sector. After fixing technical issues, the miner should call `RecoveryDeclaration` and produce the WindowPoSt challenge in order to regain the power from that sector.

Note that if a sector is in a faulty state for 14 consecutive days it will be terminated and the miner will receive a penalty. The miner can terminate the sector themselves by calling `TerminationDeclaration`, if they know that they cannot recover it, in which case they will receive a smaller penalty fee.

Both the `RecoveryDeclaration` and the `TerminationDeclaration` can be found in the [miner actor implementation](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/miner/miner_actor.go).