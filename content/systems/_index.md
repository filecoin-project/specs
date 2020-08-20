---
title: Systems
description: Filecoin Systems
weight: 2

dashboardWeight: 1
dashboardState: wip
dashboardAudit: n/a
dashboardTests: 0
---

# Systems
---

In this section we are detailing all the system components one by one in increasing level of complexity and/or interdependence to other system components. The interaction of the components between each other is only briefly discussed where appropriate, but the overall workflow is given in the Introduction section. In particular, in this section we discuss:

- Filecoin Nodes: the different types of nodes that participate in the Filecoin Network, as well as important parts and processes that these nodes run, such as the key store and IPLD store, as well as the network interface to libp2p.
- Files & Data: the data units of Filecoin, such as the Sectors and the Pieces.
- Virtual Machine: the subcomponents of the Filecoin VM, such as the actors, i.e., the smart contracts that run on the Filecoin Blockchain, and the State Tree.
- Blockchain: the main building blocks of the Filecoin blockchain, such as the structure of the Transaction and Block messages, the message pool,  as well as how nodes synchronise the blockchain when they first join the network.
- Token: the components needed for a wallet.
- Storage Mining: the details of storage mining, storage power consensus, and how storage miners prove storage (without going into details of proofs, which are discussed later).
- Markets: the storage and retrieval markets, which are primarily processes that take place off-chain, but are very important for the smooth operation of the decentralised storage market.
