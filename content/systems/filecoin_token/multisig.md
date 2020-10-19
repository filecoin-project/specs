---
title: Multisig Wallet
weight: 4
bookCollapseSection: true
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: done
dashboardAuditURL: /#section-appendix.audit_reports.specs-actors
dashboardAuditDate: '2020-10-19'
dashboardTests: 0
---

# Multisig Wallet & Actor


The Multisig actor is a single actor representing a group of Signers. Signers may be external users, other Multisigs, or even the Multisig itself. There should be a maximum of 256 signers in a multisig wallet. In case more signers are needed, then the multisigs should be combined into a tree. 

The implementation of the Multisig Actor can be found [here](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/multisig/multisig_actor.go).

The Multisig Actor statuses can be found [here](https://github.com/filecoin-project/specs-actors/blob/master/actors/builtin/multisig/multisig_state.go).

