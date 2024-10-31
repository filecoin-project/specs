---
title: Implementing Systems
weight: 2
dashboardWeight: 0.2
dashboardState: reliable
dashboardAudit: n/a
---

# Implementing Systems

## System Requirements

In order to make it easier to decouple functionality into systems, the Filecoin Protocol assumes
a set of functionality available to all systems. This functionality can be achieved by implementations
in a variety of ways, and should take the guidance here as a recommendation (SHOULD).

All Systems, as defined in this document, require the following:

- **Repository:**
  - **Local `IpldStore`.** Some amount of persistent local storage for data structures (small structured objects).
    Systems expect to be initialized with an IpldStore in which to store data structures they expect to persist across crashes.
  - **User Configuration Values.** A small amount of user-editable configuration values.
    These should be easy for end-users to access, view, and edit.
  - **Local, Secure `KeyStore`.** A facility to use to generate and use cryptographic keys, which MUST remain secret to the
    Filecoin Node. Systems SHOULD NOT access the keys directly, and should do so over an abstraction (ie the `KeyStore`) which
    provides the ability to Encrypt, Decrypt, Sign, SigVerify, and more.
- **Local `FileStore`.** Some amount of persistent local storage for files (large byte arrays).
  Systems expect to be initialized with a FileStore in which to store large files.
  Some systems (like Markets) may need to store and delete large volumes of smaller files (1MB - 10GB).
  Other systems (like Storage Mining) may need to store and delete large volumes of large files (1GB - 1TB).
- **Network.** Most systems need access to the network, to be able to connect to their counterparts in other Filecoin Nodes.
  Systems expect to be initialized with a `libp2p.Node` on which they can mount their own protocols.
- **Clock.** Some systems need access to current network time, some with low tolerance for drift.
  Systems expect to be initialized with a Clock from which to tell network time. Some systems (like Blockchain)
  require very little clock drift, and require _secure_ time.

For this purpose, we use the `FilecoinNode` data structure, which is passed into all systems at initialization.

## System Limitations

Further, Systems MUST abide by the following limitations:

- **Random crashes.** A Filecoin Node may crash at any moment. Systems must be secure and consistent through crashes.
  This is primarily achieved by limiting the use of persistent state, persisting such state through Ipld data structures,
  and through the use of initialization routines that check state, and perhaps correct errors.
- **Isolation.** Systems must communicate over well-defined, isolated interfaces. They must not build their critical
  functionality over a shared memory space. (Note: for performance, shared memory abstractions can be used to power
  IpldStore, FileStore, and libp2p, but the systems themselves should not require it.) This is not just an operational
  concern; it also significantly simplifies the protocol and makes it easier to understand, analyze, debug, and change.
- **No direct access to host OS Filesystem or Disk.** Systems cannot access disks directly -- they do so over the FileStore
  and IpldStore abstractions. This is to provide a high degree of portability and flexibility for end-users, especially
  storage miners and clients of large amounts of data, which need to be able to easily replace how their Filecoin Nodes
  access local storage.
- **No direct access to host OS Network stack or TCP/IP.** Systems cannot access the network directly -- they do so over the
  libp2p library. There must not be any other kind of network access. This provides a high degree of portability across
  platforms and network protocols, enabling Filecoin Nodes (and all their critical systems) to run in a wide variety of
  settings, using all kinds of protocols (eg Bluetooth, LANs, etc).
