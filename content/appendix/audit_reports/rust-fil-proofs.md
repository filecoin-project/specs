---
title: rust-fil-proofs
weight: 1
dashboardState: wip
dashboardAudit: n/a
---

# rust-fil-proofs

## 2020-07-28: Filecoin Proving Subsystem

Audit report: [Security Assessment - Filecoin Proving Subsystem](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/Sigma-Prime-Protocol-Labs-Filecoin-Proofs-Security-Review-v2.1.pdf)

This audit covers the full Proving Subsystem, including [rust-fil-proofs](https://github.com/filecoin-project/rust-fil-proofs) and [filecoin-ffi](https://github.com/filecoin-project/filecoin-ffi), through which Proof of Space-Time (PoSt), Proof of Retrievability (PoR), and Proof of Replication (PoRep) are implemented. The audit process included using fuzzing to identify potential vulnerabilities in the subsystem, each of which was resolved (the details of all issues raised and their resolutions are available in the report).

## 2020-07-28: zk-SNARK proofs

Audit report: [zk-SNARK Proofs Audit](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/protocolai-audit-20200728.pdf)

This audit covers the core logic and implementation of the zk-SNARK tree-based proofs-of-replication (including the [fork of bellman](https://github.com/filecoin-project/bellman)), as well as the SNARK circuits creation. All issues raised by the audit were resolved.
