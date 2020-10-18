---
title: Audit Reports
weight: 4
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Audit Reports

Security is a critical component in ensuring Filecoin can fulfill its mission to be the storage network for humanity. In addition to robust secure development processes, trainings, theory audits, and investing in external security research, the Filecoin project has engaged reputable third party auditing specialists to ensure that the theory behind the protocol and its implementation delivers the intended value, enabling Filecoin to be a safe and secure network. This section covers a selection of audit reports that have been published on Filecoin's theory and implementation.

## GossipSub

### 2020-06-03: GossipSub Design and Implementation

Report: [GossipSub v1.1 Protocol Design + Implementation](https://gateway.ipfs.io/ipfs/QmWR376YyuyLewZDzaTHXGZr7quL5LB13HRFnNdSJ3CyXu/Least%20Authority%20-%20Gossipsub%20v1.1%20Final%20Audit%20Report%20%28v2%29.pdf)

Audit conducted by: Least Authority

This audit focused specifically on GossipSub, a pubsub protocol built on libp2p, version 1.1, which includes a peer scoring layer to mitigate certain types of attacks that could compromise a network. The audit covered the [spec](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub), [go-libp2p-pubsub](https://github.com/libp2p/go-libp2p-pubsub) and [gossipsub-hardening](https://github.com/libp2p/gossipsub-hardening/). The report found 4 issues, primarily in the Peer Scoring that was introduced in v1.1, and includes additional suggestions. All the issues raised in the report have been resolved, and additional details are available in the report linked above.

### 2020-04-18: GossipSub Evaluation

Report: [GossipSub-v1.1 Evaluation Report](https://gateway.ipfs.io/ipfs/QmRAFP5DBnvNjdYSbWhEhVRJJDFCLpPyvew5GwCCB4VxM4)

Evaluation by: Protocol Labs

This evaluation focused on demonstrating that GossipSub is resilient against a range of attacks, capable of recovering the mesh, and can meet the message delivery requirements for Filecoin. Attacks used in testing include the Sybil, Eclipse, Degredation, Censorship, Attack at Dawn, "Cover Flash", and "Cold Boot" attacks. The spec for [v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md), [v1.0](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.0.md) and the [reference implementation](https://github.com/libp2p/go-libp2p-pubsub) were in scope for this audit.

## Drand

### 2020-08-09: drand

Report: [Drand Security Assessment](https://drive.google.com/file/d/1fCy1ynO78gJLCNbqBruzHx7bh72Tu-q2/view)

Audit conducted by: Sigma Prime

This report covers the end-to-end audit carried out on drand, including the implementations found in [drand/drand](https://github.com/drand/drand), [drand/bls12-381](https://github.com/drand/bls12-381) and [drand/kyber](https://github.com/drand/kyber). The audit assessed drand's ability to securely provide a distributed, continuous source of entropy / randomness for Filecoin, and included using fuzzing to find potential leaks, errors, or other panics. A handful of issues were found, 14 of which were marked as issues ranging from low to high risk, all of which have been resolved (the details of all issues raised and their resolutions are available in the report).

## rust-fil-proofs

### 2020-07-28: Filecoin Proving Subsystem

Report: [Security Assessment - Filecoin Proving Subsystem](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/Sigma-Prime-Protocol-Labs-Filecoin-Proofs-Security-Review-v2.1.pdf)

Audit conducted by: Sigma Prime

This audit covers the full Proving Subsystem, including [rust-fil-proofs](https://github.com/filecoin-project/rust-fil-proofs) and [filecoin-ffi](https://github.com/filecoin-project/filecoin-ffi), through which Proof of Space-Time (PoSt), Proof of Retrievability (PoR), and Proof of Replication (PoRep) are implemented. The audit process included using fuzzing to identify potential vulnerabilities in the subsystem, each of which was resolved (the details of all issues raised and their resolutions are available in the report).

### 2020-07-28: zk-SNARK proofs

Report: [zk-SNARK Proofs Audit](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/protocolai-audit-20200728.pdf)

Audit conducted by: Dr. Jean-Philippe Aumasson and Antony Vennard

This audit covers the core logic and implementation of the zk-SNARK tree-based proofs-of-replication (including the [fork of bellman](https://github.com/filecoin-project/bellman)), as well as the SNARK circuits creation. All issues raised by the audit were resolved.