---
title: FileStore
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
---

# FileStore - Local Storage for Files
---

The `FileStore` is an abstraction used to refer to any underlying system or device
that Filecoin will store its data to. It is based on Unix filesystem semantics, and
includes the notion of `Paths`. This abstraction is here in order to make sure Filecoin
implementations make it easy for end-users to replace the underlying storage system with
whatever suits their needs. The simplest version of `FileStore` is just the host operating
system's file system.

{{< embed src="filestore.id" lang="go" >}}

## Varying user needs

Filecoin user needs vary significantly, and many users -- especially miners -- will implement
complex storage architectures underneath and around Filecoin. The `FileStore` abstraction is here
to make it easy for these varying needs to be easy to satisfy. All file and sector local data
storage in the Filecoin Protocol is defined in terms of this `FileStore` interface, which makes
it easy for implementations to make swappable, and for end-users to swap out with their system
of choice.

## Implementation examples

The `FileStore` interface may be implemented by many kinds of backing data storage systems. For example:

- The host Operating System file system
- Any Unix/Posix file system
- RAID-backed file systems
- Networked of distributed file systems (NFS, HDFS, etc)
- IPFS
- Databases
- NAS systems
- Raw serial or block devices
- Raw hard drives (hdd sectors, etc)

Implementations SHOULD implement support for the host OS file system.
Implementations MAY implement support for other storage systems.
