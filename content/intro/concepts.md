---
title: Key Concepts
weight: 3
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Key Concepts

The Filecoin Protocol is a complex system that includes many different and novel concepts. It is highly recommended that users unfamiliar with the protocol go through the [Glossary](glossary) to familiarize themselves with terms that are unique to Filecoin. Here we provide a selected list of the most central concepts needed in order to proceed to the rest of the specification.

## FIL Network & Protocol

- **Node Types:** Nodes in the Filecoin network are defined in an abstract manner and are primarily identified in terms of the services they provide. The type of node, therefore, depends on which services the node provides. A basic set of services in the Filecoin network include: chain verification, storage market client, storage market provider, retrieval market client, retrieval market provider and storage mining. Nodes that extend the blockchain need to implement the chain verification and storage mining services. See [Node Types section](filecoin-nodes) for more details.

- **Files, Messages and Data:**
	- _File:_ a file represents the data that a user wants to store on the Filecoin network. Each node runs a `FileStore` which is simply an abstraction for any system or device where a user stores the data it has submitted to Filecoin. See the [File section](file) for more details.

    - _Sector:_ the default unit of storage that miners put in the network. The default value is set to 32GB, but sectors of 64GB are supported too. Sectors can contain data from multiple deals and multiple clients. Sectors are also split in "Regular Sectors", i.e., those that contain deals and "Committed Capacity" (CC), i.e., the sectors/storage that have been made available to the system, but for which a deal has not been agreed yet. See the [Sector section](sector) for more details.

    - _Piece:_ The [Piece](piece) is the main unit of account and negotiation for the data that a user wants to store on Filecoin. A piece is not a unit of storage and therefore, it can be of any size up to the size of a [sector](sector). If a piece is larger than a sector (currently set to 32GB or 64GB and chosen by the miner), then it has to be split in two (or more) pieces. A Piece has got its own _Piece CID_.

    - _Message:_ A message that includes a storage deal between a miner and a client. It is similar to what is referred to as a _transaction message_ in other blockchains. Messages are submitted to the network by clients and propagate to the network of miners.

    - _Block message:_ The Block is the main unit which helps advance the Filecoin Blockchain. Blocks pack transaction messages together. New blocks are produced upon every round completion, or _"Epoch"_, whose duration is 30sec.

    - _Tipset:_ A set of up to 5 blocks. The _Tipset_ is the unit of increment of the Filecoin blockchain. This means that in every round or epoch the Filecoin blockchain is extended by one Tipset (as opposed to one block as is common in other blockchains). A Tipset is a set of blocks that have the same parent set and have been mined in the same epoch.

    - _Message Pool:_ The message pool, `mpool` or `mempool` is a collection of messages that miners have submitted to be included in the blockchain. Every miner node maintains a collection of messages out of which it picks the ones it wants to include in the next block it intends to mine. See the [Message Pool section](message-pool) for more details.

- _Deal:_ The agreement between the client and a miner for the storage of data. A deal has its own *deal CID* that includes the identifiers of the client, the miner and details about the deal itself. See the [Storage Market section](storage-market) for more details on storage deals.

- _Sealing:_ the process that takes as input the data that a user submits to the Filecoin network and produces an encoded output that links the data with the specific node. This binding of a sealed sector to a node and the recording of this on the blockchain is what keeps the system secure, i.e., avoids colluding nodes from sharing storage space and reporting to the system more storage than they actually commit to the system. The sealing process enables much of the operation of the Proof of Replication and Proof of SpaceTime and is effectively what keeps the system safe. See the [Sector Sealing section](sector-sealing) for more details.

- **Actors:** Actors are the Ethereum equivalent of smart contracts in the Filecoin blockchain. Actors carry the logic needed in order to submit transactions, proofs and blocks, among other things, to the Filecoin blockchain. Every actor is identified by a unique address. See the [System Actors section](sysactors) for more details.

- **Virtual Machine:** The Filecoin Virtual Machine executes actors code and maintains the _state tree_, which is the latest source of truth in the Filecoin blockchain. See the [Virtual Machine section](filecoin_vm) for more details.

- **Power Table:** The table where the power (in terms of storage capacity) of each miner is kept.

- **Token:** The Filecoin token, also referred to as FIL. It is the unit of payment in the Filecoin network.

- **Storage Mining:** The system that guarantees that a storage miner can participate in: i) the Storage Market by agreeing deals with clients to store their data, and ii) the Storage Power Consensus to verify blocks, reach consensus and grow the blockchain. See the [Storage Mining System section](filecoin_mining) for more details.

- **Algorithms/Proofs:** the proofs that can verify that a specific storage miner has sealed and is actually storing the data as they promised according to the storage deal.
    - **Proof of Replication:** The algorithm that verifies that a miner has created a unique copy of some data, according to a storage deal agreed previously.
    - **Proof of SpaceTime:** The algorithm that verifies that the miner has continuously stored a unique copy of some data.

- **Faults:** A miner might fail to generate a proof out of a sector, which triggers a fault. This is a very important part of the system. A fault might be due to a technical failure (hardware or network failure), or due to malicious intentions (?), e.g., a miner not storing the sector anymore in order to allocate storage to a different sector and attempt to double its reward. There are different types of Faults defined based on whether and when a miner has reported the fault to the system. Depending on these parameters the corresponding fee/penalty is applied to the miner.

- **Markets:** There are two types of Markets in the Filecoin Network: the Storage Market and the Retrieval Market. Peers negotiate storage and retrieval deals, which include all the details on the price of the service and other important details. See the [Filecoin Markets section](filecoin_markets) for more details.

## FIL Implementation

For clarity, we refer the following types of entities to describe implementations of the  Filecoin protocol:

- **_Data structures_** are collections of semantically-tagged data members (e.g., structs, interfaces, or enums).

- **_Functions_** are computational procedures that do not depend on external state (i.e., mathematical functions, or programming language functions that do not refer to global variables).

- **_Components/Services_** are sets of functionality that are intended to be represented as single software units in the implementation structure. Depending on the choice of language and the particular component, this might correspond to a single software module, a thread or process running some main loop, a disk-backed database, or a variety of other design choices. For example, the [ChainSync](chainsync) is a component: it could be implemented as a process or thread running a single specified main loop, which waits for network messages and responds accordingly by recording and/or forwarding block data.

- **_APIs_** are the interfaces for delivering messages to components. A client's view of a given sub-protocol, such as a request to a miner node's Storage Provider component to store files in the storage market, may require the execution of a series of API requests.

- **_Nodes_** are complete software and hardware systems that interact with the protocol.  A node might be constantly running several of the above _components or services_, participating in several _subsystems_, and exposing _APIs_ locally and/or over the network, depending on the node configuration. The term _full node_ refers to a system that runs all of the above components and supports all of the APIs detailed in the spec.

- **_Subsystems_** are conceptual divisions of the entire Filecoin protocol, either in terms of complete protocols (such as the [Storage Market](storage_market) or [Retrieval Market](retrieval_market)), or in terms of functionality (such as the [VM - Virtual Machine](intro/filecoin_vm)). They do not necessarily correspond to any particular node or software component.

- **_Actors_** are virtual entities embodied in the state of the Filecoin VM. Protocol actors are analogous to participants in smart contracts; an actor carries a FIL currency balance and can interact with other actors via the operations of the VM, but does not necessarily correspond to any particular node or software component.
