---
title: Audit Reports
weight: 4
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Audit Reports

Security is a critical component in ensuring Filecoin can fulfill its mission to be the storage network for humanity. In addition to robust secure development processes, trainings, theory audits, and investing in external security research, the Filecoin project has engaged reputable third party auditing specialists to ensure that the theory behind the protocol and its implementation delivers the intended value, enabling Filecoin to be a safe and secure network. This section covers a selection of audit reports that have been published on Filecoin's theory and implementation.

## Filecoin Virtual Machine

### `2023-03-09` Filecoin EVM (FEVM)

- Report: [Filecoin EVM Audit](<https://github.com/oak-security/audit-reports/blob/master/Filecoin%20Foundation/2023-03-09%20Audit%20Report%20-%20Filecoin%20EVM%20(FEVM)%20v1.1.pdf>)
- Audit conducted by **Oak Security**

The audit covers the implementation of:

- FEVM's [builtin actors](https://github.com/filecoin-project/builtin-actors/tree/1b11df4b399550753a4105f45f58bc07015af2a3/actors/evm) out of which only [actors/evm](https://github.com/filecoin-project/builtin-actors/tree/1b11df4b399550753a4105f45f58bc07015af2a3/actors/evm) and [actors/eam](https://github.com/filecoin-project/builtin-actors/tree/1b11df4b399550753a4105f45f58bc07015af2a3/actors/eam) were included in scope along with code base of [ref-fvm](https://github.com/filecoin-project/ref-fvm). The report included auditing EVM runtime action and implementation, correctness of EVM opcodes, including Ethereum Address Manager(EAM). The report also included issues and enhancements methods for gas model and F4 addresses. The audit team also reviewed the message execution flow and kernel setup, WASM integration and FVM logs. All the valid issues raised by the audit were resolved and acknowledged including a few informational issues. More details on these issues are available in the report.

## Lotus

### `2020-10-20` Lotus Mainnet Ready Security Audit

- Report: [Lotus Security Assessment](https://drive.google.com/file/d/1pJnvxlz4ie9oB4NyzPRsTfs0QgYWaDdW/view)
- Audit conducted by: **Sigma Prime**

The scope of this audit covered:

- The Lotus Daemon: Core component responsible for handling the Blockchain node logic by handling peer- to-peer networking, chain syncing, block validation, data retrieval and transfer, etc.
- The Lotus Storage Miner: Mining component used to manage a single storage miner by contributing to the network through Sector commitments and Proofs-of-Spacetime data proving it is storing the sectors it has committed to. This component communicates with the Lotus daemon via JSON-RPC API calls.

## Venus

### `2021-06-29` Venus Security Audit

- Report: [Venus Security Assessment](https://leastauthority.com/static/publications/LeastAuthority_Filecoin_Foundation_Venus_Final_Audit_Report.pdf)
- Audit conducted by: **Least Authority**

The scope of this audit covered:

- The Venus Daemon: Core component responsible for handling the Filecoin node logic by handling peer-to-peer networking, chain syncing, block validation, etc.

## Actors

### `2020-10-19` Actors Mainnet Ready Security Audit

- Report: [**Filecoin Actors Audit**](https://diligence.consensys.net/audits/2020/09/filecoin-actors/)
- Audit conducted by: **Consensys Diligence**

This audit covers the implementation of Filecoin's builtin Actors, focusing on the role of Actors as a core component in the business logic of the Filecoin storage network. The audit process involved a manual review of the Actors code and conducting ongoing reviews of changes to the code during the course of the engagement. Issues uncovered through this process are all tracked in the GitHub repository. All Priority 1 issues have been resolved. Most Priority 2 issues have been resolved - ones that are still open have been determined to not be a risk for the Filecoin network or miner experience. Further details on these and all other issues raised available in the report.

## Proofs

### `2021-05-31` SnarkPack audit

An audit was conducted on the cryptographic part of [SnarkPack](https://eprint.iacr.org/2021/529.pdf), that is used in the [FIP0009](https://github.com/filecoin-project/FIPs/blob/master/FIPS/fip-0009.md):

- [Report](https://hackmd.io/@LIRa8YONSwKxiRz3cficng/B105no8w_) from Matteo Campanelli, a well known cryptography [researcher](https://www.binarywhales.com/)

One major issue was found in the report by Campanelli where the challenges of each prove commits were not tied to the aggregated proof; this could have led up to malicious miner forge valid aggregated proofs without the individual prove commits. The rest of the issues were of medium to informal severity.

### `2020-10-20` Filecoin Bellman and BLS Signatures

- Report: [**Filecoin Bellman/BLS Signatures Cryptography Review**](https://research.nccgroup.com/wp-content/uploads/2020/10/NCC_Group_ProtocolLabs_PRLB007_Report_2020-10-20_v1.0.pdf)
- Audit conducted by: **NCC Group**

This audit covers the core cryptographic primitives used by the Filecoin Proving subsystem, including BLS signatures, cryptographic arithmetic, pairings, and zk-SNARK operations. The scope of the audit included several repositories (most code is written in rust) - [bls-signatures](https://github.com/filecoin-project/bls-signatures/), Filecoin's [bellman](https://github.com/filecoin-project/bellman/), [ff](https://github.com/filecoin-project/ff), [group](https://github.com/filecoin-project/group), [paired](https://github.com/filecoin-project/paired), and [rush-sha2ni](https://github.com/filecoin-project/rust-sha2ni).The audit uncovered 1 medium severity issue which has been fixed, and a few other low severity/informational issues (the details of all issues raised and their status at time of publishing are available in the report).

### `2020-07-28` Filecoin Proving Subsystem

- Report: [**Security Assessment - Filecoin Proving Subsystem**](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/Sigma-Prime-Protocol-Labs-Filecoin-Proofs-Security-Review-v2.1.pdf)
- Audit conducted by: **Sigma Prime**

This audit covers the full Proving subsystem, including [rust-fil-proofs](https://github.com/filecoin-project/rust-fil-proofs) and [filecoin-ffi](https://github.com/filecoin-project/filecoin-ffi), through which Proof of Space-Time (PoSt), Proof of Retrievability (PoR), and Proof of Replication (PoRep) are implemented. The audit process included using fuzzing to identify potential vulnerabilities in the subsystem, each of which was resolved (the details of all issues raised and their resolutions are available in the report).

### `2020-07-28` zk-SNARK proofs

- Report: [zk-SNARK Proofs Audit](https://github.com/filecoin-project/rust-fil-proofs/blob/master/audits/protocolai-audit-20200728.pdf)
- Audit conducted by: **Dr. Jean-Philippe Aumasson and Antony Vennard**

This audit covers the core logic and implementation of the zk-SNARK tree-based proofs-of-replication (including the [fork of bellman](https://github.com/filecoin-project/bellman)), as well as the SNARK circuits creation. All issues raised by the audit were resolved.

## GossipSub

### `2020-06-03` GossipSub Design and Implementation

- Report: [**GossipSub v1.1 Protocol Design + Implementation**](https://gateway.ipfs.io/ipfs/QmWR376YyuyLewZDzaTHXGZr7quL5LB13HRFnNdSJ3CyXu/Least%20Authority%20-%20Gossipsub%20v1.1%20Final%20Audit%20Report%20%28v2%29.pdf)
- Audit conducted by: **Least Authority**

This audit focused specifically on GossipSub, a pubsub protocol built on libp2p, version 1.1, which includes a peer scoring layer to mitigate certain types of attacks that could compromise a network. The audit covered the [spec](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub), [go-libp2p-pubsub](https://github.com/libp2p/go-libp2p-pubsub) and [gossipsub-hardening](https://github.com/libp2p/gossipsub-hardening/). The report found 4 issues, primarily in the Peer Scoring that was introduced in v1.1, and includes additional suggestions. All the issues raised in the report have been resolved, and additional details are available in the report linked above.

### `2020-04-18` GossipSub Evaluation

- Report: [**GossipSub-v1.1 Evaluation Report**](https://gateway.ipfs.io/ipfs/QmRAFP5DBnvNjdYSbWhEhVRJJDFCLpPyvew5GwCCB4VxM4)
- Evaluation by: **ResNetLab @ Protocol Labs**

This evaluation focused on demonstrating that GossipSub is resilient against a range of attacks, capable of recovering the mesh, and can meet the message delivery requirements for Filecoin. Attacks used in testing include the Sybil, Eclipse, Degredation, Censorship, Attack at Dawn, "Cover Flash", and "Cold Boot" attacks. The spec for [v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md), [v1.0](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.0.md) and the [reference implementation](https://github.com/libp2p/go-libp2p-pubsub) were in scope for this audit.

## Drand

### `2020-08-09` drand reference implementation Security Audit

- Report: [**Drand Security Assessment**](https://drive.google.com/file/d/1fCy1ynO78gJLCNbqBruzHx7bh72Tu-q2/view)
- Audit conducted by: **Sigma Prime**

This report covers the end-to-end audit carried out on drand, including the implementations found in [drand/drand](https://github.com/drand/drand), [drand/bls12-381](https://github.com/drand/bls12-381) and [drand/kyber](https://github.com/drand/kyber). The audit assessed drand's ability to securely provide a distributed, continuous source of entropy / randomness for Filecoin, and included using fuzzing to find potential leaks, errors, or other panics. A handful of issues were found, 14 of which were marked as issues ranging from low to high risk, all of which have been resolved (the details of all issues raised and their resolutions are available in the report).
