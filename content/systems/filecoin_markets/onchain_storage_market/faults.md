---
title: "Faults"
weight: 4
dashboardWeight: 2
dashboardState: reliable
dashboardAudit: wip
dashboardTests: 0
---

# Faults

There are two main categories of faults in the Filecoin network:

1. [Storage or Sector Faults](sector#sector-faults) that relate with the failure to store files agreed in a deal previously due to a hardware error or malicious behaviour, and
2. [Consensus Faults](expected_consensus#consensus-faults) that relate to a miner trying deviate from the protocol in order to gain more power than their storage deserves.

Please refer to the corresponding sections for more details.

Both Storage and Consensus Faults come with penalties that slash the miner's collateral. See more details on the different types of collaterals in the [Miner Collaterals](filecoin_mining#miner_collaterals).
