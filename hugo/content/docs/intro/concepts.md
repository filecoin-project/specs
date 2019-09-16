---
title: "Key Concepts"
---

For clarity, we refer the following types of entities to describe implementations of the Filecoin protocol:

- **_Data structures_** are collections of semantically-tagged data members (e.g., structs, interfaces, or enums).

- **_Functions_** are computational procedures that do not depend on external state (i.e., mathematical functions,
  or programming language functions that do not refer to global variables).

- **_Components_** are sets of functionality that are intended to be represented as single software units
  in the implementation structure.
  Depending on the choice of language and the particular component, this might
  correspond to a single software module,
  a thread or process running some main loop, a disk-backed database, or a variety of other design choices.
  For example, the {{<sref block_propagator>}} is a component: it could be implemented
  as a process or thread running a single specified main loop, which waits for network messages
  and responds accordingly by recording and/or forwarding block data.

- **_APIs_** are messages that can be sent to components.
  A client's view of a given sub-protocol, such as a request to a miner node's
  {{<sref storage_provider>}} component to store files in the storage market,
  may require the execution of a series of APIs.

- **_Nodes_** are complete software and hardware systems that interact with the protocol.
  A node might be constantly running several of the above _components_, participating in several _subsystems_,
  and exposing _APIs_ locally and/or over the network,
  depending on the node configuration.
  The term _full node_ refers to a system that runs all of the above components, and supports all of the APIs detailed in the spec.

- **_Subsystems_** are conceptual divisions of the entire Filecoin protocol, either in terms of complete protocols
  (such as the {{<sref storage_market>}} or {{<sref retrieval_market>}}), or in terms of functionality
  (such as the {{<sref sys_vm>}}). They do not necessarily correspond to any particular node or software component.

- **_Actors_** are virtual entities embodied in the state of the Filecoin VM.
  Protocol actors are analogous to participants in smart contracts;
  an actor carries a FIL currency balance and can interact with other actors
  via the operations of the VM, but does not necessarily correspond to any particular node or software component.
