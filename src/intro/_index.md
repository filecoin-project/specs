---
title: Introduction
entries:
#  - specmap
- arch
- concepts
- filecoin_vm
- process
- system
---

<center><img src="./docs/intro/underconstruction.gif" height="128px" /></center>

{{% notice warning %}}
**Warning:** This draft of the Filecoin protocol specification is a work in progress.
It is intended to establish the rough overall structure of the document,
enabling experts to fill in different sections in parallel.
However, within each section, content may be out-of-order, incorrect, and/or incomplete.
The reader is advised to refer to the
[official Filecoin spec document](https://filecoin-project.github.io/specs/)
for specification and implementation questions.
{{% /notice %}}

Filecoin is a distributed storage network based on a blockchain mechanism.
Filecoin *miners* can elect to provide storage capacity for the network, and thereby
earn units of the Filecoin cryptocurrency (FIL) by periodically producing
cryptographic proofs that certify that they are providing the capacity specified.
In addition, Filecoin enables parties to exchange FIL currency
through transactions recorded in a shared ledger on the Filecoin blockchain.
Rather than using Nakamoto-style proof of work to maintain consensus on the chain, however,
Filecoin uses proof of storage itself: a miner's power in the consensus protocol
is proportional to the amount of storage it provides.

The Filecoin blockchain not only maintains the ledger for FIL transactions and
accounts, but also implements the Filecoin VM, a replicated state machine which executes
a variety of cryptographic contracts and market mechanisms among participants
on the network.
These contracts include *storage deals*, in which clients pay FIL currency to miners
in exchange for storing the specific file data that the clients request.
Via the distributed implementation of the Filecoin VM, storage deals
and other contract mechanisms recorded on the chain continue to be processed
over time, without requiring further interaction from the original parties
(such as the clients who requested the data storage).
